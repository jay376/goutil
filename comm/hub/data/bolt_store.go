package data

import (
	"sync"
	"time"

	bolt "go.etcd.io/bbolt"
)

var defaultBucket = []byte("DATA")
var defaultCommitInteval = 100 * time.Millisecond
var defaultMaxPending = 1000

// ReadTx ...
type ReadTx interface {
	Get(bucket []byte, key []byte) ([]byte, error)
	UnsafeRange(bucket []byte, start, end []byte) (KVs, error)
	UnsafeVisit(bucket []byte, start, end []byte, v VisitFunc) error
}

// BatchTx ...
type BatchTx interface {
	ReadTx
	UnsafePut(bucket []byte, key []byte, value []byte) error
	UnsafeDelete(bucket []byte, key []byte) error
	// Commit commits a previous tx and begins a new writable one.
	UnsafeCommit() error
}

type boltstore struct {
	db     *bolt.DB
	bucket []byte
	mu     sync.RWMutex // protect commit
	wmu    sync.Mutex   // protect write

	rtx              ReadTx  // read tx
	btx              BatchTx // read write tx
	commitCh         chan struct{}
	pending          int
	stopc            chan struct{}
	revision         uint64
	commitedRevision uint64
}

func newBoltStore(path string) *boltstore {
	db, err := bolt.Open(path, 0600, nil)

	if err != nil {
		panic(err)
	}

	b := &boltstore{
		db:       db,
		bucket:   defaultBucket,
		commitCh: make(chan struct{}, 1),
		pending:  0,
		rtx:      &readtx{},
		stopc:    make(chan struct{}, 1),
	}
	b.btx = Batch(b)
	go b.run()
	return b
}

// Close ...
func (b *boltstore) Close() error {
	b.stopc <- struct{}{}
	return nil
}

func (b *boltstore) run() {
	t := time.NewTicker(defaultCommitInteval)
	defer t.Stop()
loop:
	for {
		select {
		case <-t.C:
			b.Commit()
		case <-b.commitCh:
			b.Commit()
		case <-b.stopc:
			break loop
		}
	}
	b.Commit()
	b.db.Close()
}

func (b *boltstore) Commit() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.pending == 0 {
		return nil
	}
	if err := b.btx.UnsafeCommit(); err != nil {
		return err
	}
	b.commitedRevision = b.revision
	b.pending = 0
	return nil
}

// Get ...
func (b *boltstore) Get(key []byte) ([]byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.rtx.Get(b.bucket, key)
}

// Range ...
func (b *boltstore) Range(start, end []byte) (KVs, uint64, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	kvs, e := b.rtx.UnsafeRange(b.bucket, start, end)
	return kvs, b.commitedRevision, e
}

// UnsafeVisit ...
func (b *boltstore) Visit(start, end []byte, v VisitFunc) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.rtx.UnsafeVisit(b.bucket, start, end, v)
}

// Clear all data
func (b *boltstore) Clear() error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	b.wmu.Lock()
	// only one writer
	defer b.wmu.Unlock()

	return b.btx.UnsafeDelete(b.bucket, nil)
}

// Put single key value
func (b *boltstore) Put(key []byte, value []byte) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	b.wmu.Lock()
	// only one writer
	defer b.wmu.Unlock()
	if b.pending >= defaultMaxPending {
		select {
		case b.commitCh <- struct{}{}:
		default:
		}
		b.pending = 0
	}
	b.pending++
	return b.btx.UnsafePut(b.bucket, key, value)
}

// Del single key value
func (b *boltstore) Del(key []byte) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	b.wmu.Lock()
	// only one writer
	defer b.wmu.Unlock()

	if b.pending >= defaultMaxPending {
		select {
		case b.commitCh <- struct{}{}:
		default:
		}
		b.pending = 0
	}
	b.pending++
	return b.btx.UnsafeDelete(b.bucket, key)
}
