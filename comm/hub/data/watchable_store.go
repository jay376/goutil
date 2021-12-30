package data

// WatchableStore ...
type WatchableStore interface {
	Store
	Watch() chan []*Event
}

type watchableStore struct {
	*boltstore
	diff chan []*Event
}

// NewWatchableStore ...
func NewWatchableStore(path string) WatchableStore {
	boltstore := newBoltStore(path)
	ws := &watchableStore{
		boltstore: boltstore,
		diff:      make(chan []*Event, 100),
	}

	watchtx := &watchableTx{
		BatchTx: boltstore.btx,
		changes: make([]*Event, 0, 100),
		wstore:  ws,
	}
	boltstore.btx = watchtx
	return ws
}

// Watch ...
func (w *watchableStore) Watch() chan []*Event {
	return w.diff
}
