package data

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

var wstore = NewWatchableStore("./wtmp")
var wgkeys [][]byte
var sum = 0
var expectSum = 0

func BenchmarkWatchPut(b *testing.B) {
	go func() {
		// for {
		// 	select {
		// 	case d := <-wstore.Watch():
		// 		sum += len(d)
		// 	}
		// }
		for d := range wstore.Watch() {
			sum += len(d)
		}
	}()

	b.ResetTimer()
	keys := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = make([]byte, 64)
		rand.Read(keys[i])
		wstore.Put(keys[i], keys[i])
		// wstore.Commit()
	}
	expectSum += b.N
}

func BenchmarkWatchBatchPut(b *testing.B) {
	b.ResetTimer()
	keys := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = make([]byte, 64)
		rand.Read(keys[i])
		wstore.Put(keys[i], keys[i])
	}
	expectSum += b.N
	// wstore.Commit()
	wgkeys = keys
}

func BenchmarkWatchGet(b *testing.B) {
	wstore.(*watchableStore).Commit()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx := i % len(wgkeys)
		value, _ := wstore.Get(wgkeys[idx])
		assert.Equal(b, value, wgkeys[idx])
	}
}

func BenchmarkNotifyDiff(b *testing.B) {
	// judge put num with notify diff num
	assert.Equal(b, expectSum, sum)
}
