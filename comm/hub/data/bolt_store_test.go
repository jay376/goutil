package data

import (
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	// "time"
)

var store = newBoltStore("./tmp")
var gkeys [][]byte

func TestPutRange(t *testing.T) {
	kvs := KVs{
		{[]byte("123"), []byte("222")},
		{[]byte("124"), []byte("wdf")},
		{[]byte("135"), []byte("233")},
		{[]byte("333"), []byte("324")},
		{[]byte("384"), []byte("23423")},
	}
	store.Clear()
	for _, item := range kvs {
		if e := store.Put(item.Key, item.Value); e != nil {
			t.Error(e)
		}
	}

	// time.Sleep(defaultCommitInteval)
	store.Commit()
	// range all
	if rkvs, _, e := store.Range(nil, nil); e != nil {
		t.Error(e)
	} else {
		assert.Len(t, rkvs, len(kvs))
	}

	if rkvs, _, e := store.Range(kvs[0].Key, kvs[len(kvs)-1].Key); e != nil {
		t.Error(e)
	} else {
		t.Log(len(rkvs))
		assert.Equal(t, true, reflect.DeepEqual(rkvs, kvs))
	}
	num := 0
	visit := func(kv *KeyValue) {
		num++
	}
	store.Visit(nil, nil, visit)
	assert.Equal(t, num, len(kvs))
	assert.Equal(t, nil, store.Del([]byte("384")))
}

func BenchmarkPut(b *testing.B) {
	store.Clear()
	b.ResetTimer()
	keys := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = make([]byte, 64)
		rand.Read(keys[i])
		store.Put(keys[i], keys[i])
		store.Commit()
	}
	// time.Sleep(defaultCommitInteval)
}

func BenchmarkBatchPut(b *testing.B) {
	store.Clear()
	b.ResetTimer()
	keys := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = make([]byte, 64)
		rand.Read(keys[i])
		store.Put(keys[i], keys[i])
	}
	if rkvs, _, e := store.Range([]byte{0}, []byte{255}); e != nil {
		b.Error(e)
	} else {
		assert.Len(b, rkvs, len(keys))
	}
	// time.Sleep(defaultCommitInteval)
	store.Commit()
	gkeys = keys
}

func BenchmarkGet(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx := i % len(gkeys)
		value, _ := store.Get(gkeys[idx])
		assert.Equal(b, value, gkeys[idx])
	}
}
