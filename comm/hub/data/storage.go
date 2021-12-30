package data

// KeyValue ...
type KeyValue struct {
	Key   []byte
	Value []byte
}

// KVs ...
type KVs []KeyValue

// VisitFunc ...
type VisitFunc func(kv *KeyValue)

// Store ...
type Store interface {
	// Range [start, end) ...
	Range(start, end []byte) (KVs, uint64, error)

	// Get single key value
	Get(key []byte) ([]byte, error)

	// Del single key, default batch commit with delay some ms < 100ms
	// , if need take effective right now, so call Commit
	Del(key []byte) error

	// Put single key value, default batch commit with delay some ms < 100ms
	// , if need take effective right now, so call Commit
	Put(key []byte, value []byte) error

	// Visit all kv
	Visit(start, end []byte, v VisitFunc) error

	// Close ...
	Close() error

	// Clear ...
	Clear() error
}
