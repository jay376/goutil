package data

import (
	"bytes"
	"errors"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

type readtx struct {
	*bolt.Tx // read tx
}

// Read ...
func Read(tx *bolt.Tx) ReadTx {
	return &readtx{
		Tx: tx,
	}
}

// Get ...
func (r *readtx) Get(bn []byte, key []byte) ([]byte, error) {
	var value []byte

	bucket := r.Bucket(bn)
	if bucket == nil {
		return nil, errors.New("empty bucket: " + string(bn))
	}
	v := bucket.Get(key)
	if v == nil {
		return nil, fmt.Errorf("key:%s not exist", string(key))
	}
	value = make([]byte, len(v))
	copy(value, v)
	return value, nil
}

//UnsafeRange ...
func (r *readtx) UnsafeRange(bn []byte, start, end []byte) (KVs, error) {
	var kvs KVs
	bucket := r.Bucket(bn)
	if bucket == nil {
		return nil, errors.New("empty bucket: " + string(bn))
	}
	c := bucket.Cursor()
	var k, v []byte
	if start == nil {
		k, v = c.First()
	} else {
		k, v = c.Seek(start)
	}

	for ; k != nil; k, v = c.Next() {
		if end != nil && bytes.Compare(k, end) > 0 {
			break
		}
		copyK := make([]byte, len(k))
		copyV := make([]byte, len(v))
		copy(copyK, k)
		copy(copyV, v)
		kvs = append(kvs, KeyValue{copyK, copyV})
	}
	return kvs, nil
}

// UnsafeVisit ...
func (r *readtx) UnsafeVisit(bn []byte, start, end []byte, visit VisitFunc) error {
	bucket := r.Bucket(bn)
	if bucket == nil {
		return errors.New("empty bucket: " + string(bn))
	}
	c := bucket.Cursor()
	for k, v := c.Seek(start); k != nil; k, v = c.Next() {
		if end != nil && bytes.Compare(k, end) > 0 {
			break
		}
		visit(&KeyValue{k, v})
	}
	return nil
}

// read & write tx
type batchTx struct {
	wtx   *bolt.Tx // read & write tx
	store *boltstore
}

// Batch ...
func Batch(store *boltstore) BatchTx {
	b := &batchTx{
		store: store,
	}
	b.UnsafeCommit()
	return b
}

// Get ...
func (b *batchTx) Get(bn []byte, key []byte) ([]byte, error) {
	return Read(b.wtx).Get(bn, key)
}

// UnsafeRange ...
func (b *batchTx) UnsafeRange(bn []byte, start, end []byte) (KVs, error) {
	return Read(b.wtx).UnsafeRange(bn, start, end)
}

// UnsafeVisit ...
func (b *batchTx) UnsafeVisit(bn []byte, start, end []byte, v VisitFunc) error {
	return Read(b.wtx).UnsafeVisit(bn, start, end, v)
}

// UnsafePut ...
func (b *batchTx) UnsafePut(bn []byte, key []byte, value []byte) error {
	bucket, err := b.wtx.CreateBucketIfNotExists(bn)
	if err != nil {
		return err
	}
	return bucket.Put(key, value)
}

// UnsafeDelete ...
func (b *batchTx) UnsafeDelete(bn []byte, key []byte) error {
	if key == nil {
		return b.wtx.DeleteBucket(bn)
	}
	bucket, err := b.wtx.CreateBucketIfNotExists(bn)
	if err != nil {
		return err
	}
	return bucket.Delete(key)
}

// UnsafeCommit ...
func (b *batchTx) UnsafeCommit() error {
	if b.store.rtx.(*readtx).Tx != nil {
		if err := b.store.rtx.(*readtx).Rollback(); err != nil {
			return err
		}
	}
	if b.wtx != nil {
		if err := b.wtx.Commit(); err != nil {
			return err
		}
	}
	wtx, _ := b.store.db.Begin(true)
	b.wtx = wtx
	rtx, _ := b.store.db.Begin(false)
	b.store.rtx.(*readtx).Tx = rtx
	return nil
}
