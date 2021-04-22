package comm

import (
	"container/heap"
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Cleanup ...
type Cleanup func(key string, value interface{})

// Item ...
type Item struct {
	key     string
	value   interface{}
	time    int64
	index   int
	cleanup Cleanup
}

// Key ...
func (i *Item) Key() string {
	return i.key
}

// Value ...
func (i *Item) Value() interface{} {
	return i.value
}

// ItemQueue ...
type ItemQueue []*Item

// Len ...
func (pq ItemQueue) Len() int {
	return len(pq)
}

// Less ...
func (pq ItemQueue) Less(i, j int) bool {
	return pq[i].time < pq[j].time
}

// Swap ...
func (pq ItemQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push ...
func (pq *ItemQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

// Pop ...
func (pq *ItemQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// ExpiredKv ...
type ExpiredKv struct {
	sync.RWMutex
	kv     map[string]*Item
	queue  ItemQueue
	ctx    context.Context
	cancel context.CancelFunc
}

// NewExpireKv ...
func NewExpireKv(ctx context.Context) *ExpiredKv {
	ekv := &ExpiredKv{
		kv:    make(map[string]*Item),
		queue: make(ItemQueue, 0, 100),
	}
	ekv.ctx, ekv.cancel = context.WithCancel(ctx)
	go ekv.start()
	return ekv
}

// Put key value, expireTime is unixtimenano stamp
// if expireTime = 0, key will not expire
// cleanup will be called if cleanup is not nil && key is expired
func (e *ExpiredKv) Put(key string, value interface{}, expireTime int64, clean Cleanup) {
	e.Lock()
	if old, ok := e.kv[key]; ok {
		old.value = value
		if expireTime == 0 {
			// Not need expire, so delete from heap
			heap.Remove(&e.queue, old.index)
			old.cleanup = nil
			old.time = 0
		} else {
			// Update expire time
			old.time = expireTime
			heap.Fix(&e.queue, old.index)
			zap.L().Debug("put update", zap.String("key", key))
		}
	} else {
		item := &Item{
			key:     key,
			value:   value,
			time:    expireTime,
			cleanup: clean,
		}
		if expireTime != 0 {
			heap.Push(&e.queue, item)
		}
		e.kv[key] = item
		zap.L().Debug("put", zap.String("key", key))
	}
	e.Unlock()
}

// PutNoLock ...
func (e *ExpiredKv) PutNoLock(key string, value interface{}, expireTime int64, clean Cleanup) {
	if old, ok := e.kv[key]; ok {
		old.value = value
		if expireTime == 0 {
			// Not need expire, so delete from heap
			heap.Remove(&e.queue, old.index)
			old.cleanup = nil
			old.time = 0
		} else {
			// Update expire time
			old.time = expireTime
			heap.Fix(&e.queue, old.index)
			zap.L().Debug("put update", zap.String("key", key))
		}
	} else {
		item := &Item{
			key:     key,
			value:   value,
			time:    expireTime,
			cleanup: clean,
		}
		if expireTime != 0 {
			heap.Push(&e.queue, item)
		}
		e.kv[key] = item
		zap.L().Debug("put", zap.String("key", key))
	}
}

// Get ...
func (e *ExpiredKv) Get(key string) interface{} {
	e.RLock()
	defer e.RUnlock()
	if item, ok := e.kv[key]; ok {
		zap.L().Debug("get ok", zap.String("key", key))
		return item.value
	}
	zap.L().Debug("get fail", zap.String("key", key))
	return nil
}

// GetNoLock ...
func (e *ExpiredKv) GetNoLock(key string) interface{} {
	if item, ok := e.kv[key]; ok {
		zap.L().Debug("get ok", zap.String("key", key))
		return item.value
	}
	zap.L().Debug("get fail", zap.String("key", key))
	return nil
}

// Del ...
func (e *ExpiredKv) Del(key string) (value interface{}) {
	e.Lock()
	if v, ok := e.kv[key]; ok {
		value = v.value
		if v.time != 0 {
			heap.Remove(&e.queue, v.index)
			v.cleanup = nil
			v.time = 0
		}
	}
	delete(e.kv, key)
	e.Unlock()
	zap.L().Debug("del", zap.String("key", key))
	return
}

// DelNoLock ...
func (e *ExpiredKv) DelNoLock(key string) (value interface{}) {
	if v, ok := e.kv[key]; ok {
		value = v.value
		if v.time != 0 {
			heap.Remove(&e.queue, v.index)
			v.cleanup = nil
			v.time = 0
		}
	}
	delete(e.kv, key)
	zap.L().Debug("del", zap.String("key", key))
	return
}

func runClean(item *Item) {
	gap := time.Now().UnixNano() - item.time
	if gap > int64(10*time.Second) {
		zap.L().Error("expire_clean time out 10s "+item.key, zap.Int64("gap", gap))
	}
	zap.L().Debug("expire_clean", zap.String("key", item.key))
	item.cleanup(item.key, item.value)
}

// loop handle expired key
func (e *ExpiredKv) start() {
	for {
		select {
		case <-time.After(10 * time.Millisecond):
		case <-e.ctx.Done():
			return
		}
		e.Lock()
		num := 0
		for {
			if e.queue.Len() == 0 || time.Now().UnixNano() < e.queue[0].time {
				break
			}
			num++
			item := heap.Pop(&e.queue).(*Item)
			if _, ok := e.kv[item.key]; ok && item.cleanup != nil {
				go runClean(item)
				// cleanC <- item
			}
			if item.time != 0 {
				delete(e.kv, item.key)
				zap.L().Debug("del", zap.String("key", item.key))
			}
		}
		// if e.queue.Len() > 0 {
		// 	expireAt := time.Unix(e.queue[0].time/1000000000, e.queue[0].time%1000000000)
		// 	zap.L().Debug("ekv clean ", zap.Int("keys", num), zap.Int("len(kvs", len(e.kv)), zap.Int("len(queue)", e.queue.Len()), zap.Time("expireAt", expireAt), zap.Time("now", time.Now()))

		// } else {
		// 	zap.L().Debug("ekv clean ", zap.Int("keys", num), zap.Int("len(kvs", len(e.kv)), zap.Int("len(queue)", e.queue.Len()))
		// }
		e.Unlock()
	}
}

// Stop ...
func (e *ExpiredKv) Stop() {
	e.cancel()
}

// Len ...
func (e *ExpiredKv) Len() int {
	e.Lock()
	defer e.Unlock()
	return len(e.kv)
}

// ForceDump ...
func (e *ExpiredKv) ForceDump(dump Cleanup) {
	e.Lock()
	defer e.Unlock()
	for {
		if e.queue.Len() == 0 {
			break
		}
		item := heap.Pop(&e.queue).(*Item)
		if _, ok := e.kv[item.key]; ok {
			dump(item.key, item.value)
		}
	}
}
