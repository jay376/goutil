package hub

// KeyValue ...
type KeyValue struct {
	Key   string
	Value []byte
}

// KVs ...
type KVs []*KeyValue

// Store ...
type Store interface {
	// Range [start, end) ...
	Range(start, end string) (KVs, uint64, error)

	// Get single key value
	Get(key string) ([]byte, uint64, error)

	// Del single key, default batch commit with delay some ms < 100ms
	// , if need take effective right now, so call Commit
	Del(key string) (uint64, error)

	// Put single key value, default batch commit with delay some ms < 100ms
	// , if need take effective right now, so call Commit
	Put(key string, value []byte) (uint64, error)
}
