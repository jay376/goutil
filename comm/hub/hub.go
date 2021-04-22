package hub

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Op ...
type Op int

const (
	// Put ...
	Put Op = 0
	// Del ...
	Del Op = 1
)

// Event ...
type Event struct {
	Key      string
	Value    []byte
	PreValue []byte
	Op       Op
	Revison  uint64
}

// Watcher ...
type Watcher struct {
	Key         string
	WithPreffix bool
	Events      chan []*Event
	Snapshot    uint64
}

// Hub ...
type Hub struct {
	sMu      sync.RWMutex
	st       Store
	evs      []*Event
	wMu      sync.Mutex
	watchers map[*Watcher]struct{}
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewHub ...
func NewHub(ctx context.Context) *Hub {
	c, cancel := context.WithCancel(ctx)
	h := &Hub{
		st:       newBtreeStore(),
		evs:      make([]*Event, 0, 100),
		watchers: make(map[*Watcher]struct{}),
	}
	h.ctx, h.cancel = c, cancel
	go h.notify()
	return h
}

// Stop ...
func (h *Hub) Stop() {
	h.cancel()
	h.wMu.Lock()
	defer h.wMu.Unlock()
	for watcher := range h.watchers {
		close(watcher.Events)
	}
	h.watchers = nil
}

// Put ...
func (h *Hub) Put(key string, value []byte) error {
	h.sMu.Lock()
	defer h.sMu.Unlock()

	var old []byte
	if v, _, err := h.st.Get(key); err == nil {
		old = v
		if bytes.Equal(old, value) {
			return nil
		}
	} else {
		return err
	}

	snapshot, err := h.st.Put(key, value)
	if err != nil {
		return err
	}
	event := &Event{
		Key:     key,
		Value:   value,
		Op:      Put,
		Revison: snapshot,
	}
	if old != nil {
		event.PreValue = old
	}
	h.evs = append(h.evs, event)
	return nil
}

// Del ...
func (h *Hub) Del(key string) error {
	h.sMu.Lock()
	defer h.sMu.Unlock()

	var old []byte
	v, _, err := h.st.Get(key)
	switch {
	case err != nil:
		return err
	case v == nil:
		return nil
	default:
		old = v
	}

	snapshot, err := h.st.Del(key)
	if err != nil {
		return err
	}
	event := &Event{
		Key:      key,
		PreValue: old,
		Op:       Del,
		Revison:  snapshot,
	}
	h.evs = append(h.evs, event)
	return nil
}

func getPrefix(key []byte) []byte {
	end := make([]byte, len(key))
	copy(end, key)
	for i := len(end) - 1; i >= 0; i-- {
		if end[i] < 0xff {
			end[i]++
			end = end[:i+1]
			return end
		}
	}
	return nil
}

// WatchWithPreffix ...
func (h *Hub) WatchWithPreffix(key string) (*Watcher, error) {
	end := string(getPrefix([]byte(key)))
	h.sMu.RLock()
	kvs, snaphost, err := h.st.Range(key, end)
	if err != nil {
		h.sMu.RUnlock()
		return nil, err
	}
	h.sMu.RUnlock()
	w := &Watcher{
		Key:         key,
		WithPreffix: true,
		Events:      make(chan []*Event, 10),
		Snapshot:    snaphost,
	}
	evs := make([]*Event, len(kvs))
	for idx, kv := range kvs {
		evs[idx] = &Event{
			Key:     kv.Key,
			Value:   kv.Value,
			Op:      Put,
			Revison: snaphost,
		}
	}
	if len(evs) > 0 {
		w.Events <- evs
	}
	h.wMu.Lock()
	h.watchers[w] = struct{}{}
	h.wMu.Unlock()
	return w, nil
}

// UnWatch ...
func (h *Hub) UnWatch(w *Watcher) {
	h.wMu.Lock()
	delete(h.watchers, w)
	close(w.Events)
	h.wMu.Unlock()
}

func (h *Hub) notify() {
	ticker := time.NewTicker(10 * time.Millisecond)
	for {
		select {
		case <-h.ctx.Done():
			return
		case <-ticker.C:
		}
		var evs []*Event
		h.sMu.Lock()
		if len(h.evs) == 0 {
			h.sMu.Unlock()
			continue
		}
		evs = h.evs
		h.evs = make([]*Event, 0, 100)
		h.sMu.Unlock()

		h.wMu.Lock()
		for watcher := range h.watchers {
			retEvs := make([]*Event, 0, 10)
			for _, event := range evs {
				if watcher.WithPreffix {
					if strings.HasPrefix(event.Key, watcher.Key) && event.Revison > watcher.Snapshot {
						retEvs = append(retEvs, event)
					}
				} else {
					if event.Key == watcher.Key && event.Revison > watcher.Snapshot {
						retEvs = append(retEvs, event)
					}
				}
			}

			if len(retEvs) == 0 {
				continue
			}
			select {
			case watcher.Events <- retEvs:
			default:
				zap.L().Warn("notify channel is full", zap.Any("watcher", watcher))
			}
		}
		h.wMu.Unlock()
	}
}
