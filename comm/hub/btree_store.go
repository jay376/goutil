package hub

import (
	"strings"
	"sync"

	"github.com/google/btree"
)

// Less ...
func (i *KeyValue) Less(than btree.Item) bool {
	return strings.Compare(i.Key, than.(*KeyValue).Key) < 0
}

// BtreeStore ...
type BtreeStore struct {
	tree     *btree.BTree
	snapshot uint64
	sync.RWMutex
}

func newBtreeStore() Store {
	return &BtreeStore{
		tree: btree.New(32),
	}
}

// Put ...
func (b *BtreeStore) Put(key string, value []byte) (uint64, error) {
	b.Lock()
	defer b.Unlock()
	b.tree.ReplaceOrInsert(&KeyValue{Key: key, Value: value})
	b.snapshot++
	return b.snapshot, nil
}

// Del ...
func (b *BtreeStore) Del(key string) (uint64, error) {
	b.Lock()
	defer b.Unlock()
	b.tree.Delete(&KeyValue{Key: key})
	b.snapshot++
	return b.snapshot, nil
}

// Get ...
func (b *BtreeStore) Get(key string) (value []byte, snapshot uint64, err error) {
	b.RLock()
	defer b.RUnlock()
	v := b.tree.Get(&KeyValue{Key: key})
	if v != nil {
		value = v.(*KeyValue).Value
		snapshot = b.snapshot
	}
	return
}

// Range ...
func (b *BtreeStore) Range(start, end string) (KVs, uint64, error) {
	starti, endi := &KeyValue{Key: start}, &KeyValue{Key: end}
	kvs := make([]*KeyValue, 0, 10)
	b.tree.AscendGreaterOrEqual(starti, func(item btree.Item) bool {
		if len(endi.Key) > 0 && !item.Less(endi) {
			return false
		}
		kvs = append(kvs, item.(*KeyValue))
		return true
	})

	return kvs, b.snapshot, nil
}
