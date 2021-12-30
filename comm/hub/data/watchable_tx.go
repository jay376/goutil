package data

// OpType ...
type OpType int32

const (
	// PUT ...
	PUT OpType = 0
	// DELETE ...
	DELETE OpType = 1
)

// Event ...
type Event struct {
	Key      []byte
	Value    []byte
	Op       OpType
	Revision uint64
}

type watchableTx struct {
	BatchTx
	changes []*Event
	wstore  *watchableStore
}

// UnsafePut ...
func (w *watchableTx) UnsafePut(bn []byte, key []byte, value []byte) (e error) {
	if e = w.BatchTx.UnsafePut(bn, key, value); e == nil {
		w.wstore.revision++
		w.changes = append(w.changes, &Event{
			Key:      key,
			Value:    value,
			Op:       PUT,
			Revision: w.wstore.revision,
		})
	}
	return
}

// UnsafeDelete ...
func (w *watchableTx) UnsafeDelete(bn []byte, key []byte) (e error) {
	if e = w.BatchTx.UnsafeDelete(bn, key); e == nil {
		w.wstore.revision++
		w.changes = append(w.changes, &Event{
			Key:      key,
			Op:       DELETE,
			Revision: w.wstore.revision,
		})
	}
	return
}

// UnsafeCommit ...
func (w *watchableTx) UnsafeCommit() (e error) {
	if e = w.BatchTx.UnsafeCommit(); e == nil {
		w.wstore.diff <- w.changes
	} else {
		panic(e)
	}
	w.changes = make([]*Event, 0, 100)
	return
}
