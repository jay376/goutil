package comm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	q := NewRingQueue(2000)
	for i := 0; i < 2000; i++ {
		assert.Equal(t, q.Put(i), nil)
	}

	for i := 0; i < 2000; i++ {
		e := q.Pop(1)
		assert.Equal(t, e[0], i)
	}
}

func TestGetFromZeroSize(t *testing.T) {
	q := NewRingQueue(0)
	for i := 0; i < 20; i++ {
		assert.Equal(t, q.Put(i), i)
	}

	for i := 0; i < 20; i++ {
		e := q.Pop(1)
		assert.Equal(t, len(e), 0)
	}
}

func BenchmarkPutGet(b *testing.B) {
	b.ReportAllocs()
	q := NewRingQueue(2000)
	for i := 0; i < b.N; i++ {
		n := 0
		for ; n < 2000; n++ {
			assert.Equal(b, q.Put(n), nil)
		}
		elements := q.Pop(n)
		for idx, value := range elements {
			assert.Equal(b, idx, value)
		}
	}
}

func TestChangeToZero(t *testing.T) {
	num := 10
	q := NewRingQueue(num)
	for i := 0; i < num; i++ {
		assert.Equal(t, q.Put(i), nil)
	}

	size := 0
	elements := q.Change(size)
	for idx, value := range elements {
		assert.Equal(t, idx, value)
	}

	elements = q.Pop(num)
	for idx, value := range elements {
		assert.Equal(t, idx+num-size, value)
	}
}

func TestChangeNormal(t *testing.T) {
	num := 10
	q := NewRingQueue(num)
	for i := 0; i < num; i++ {
		assert.Equal(t, q.Put(i), nil)
	}

	size := 5
	elements := q.Change(size)
	for idx, value := range elements {
		assert.Equal(t, idx, value)
	}

	elements = q.Pop(num)
	for idx, value := range elements {
		assert.Equal(t, idx+num-size, value)
	}
}

func TestChangeSpecial(t *testing.T) {
	num := 10
	q := NewRingQueue(num)
	for i := 0; i < num; i++ {
		assert.Equal(t, q.Put(i), nil)
	}
	q.Pop(8)
	for i := 0; i < 5; i++ {
		assert.Equal(t, q.Put(i+num), nil)
	}

	size := 5
	elements := q.Change(size)
	for idx, value := range elements {
		assert.Equal(t, idx+8, value)
	}

	elements = q.Pop(num)
	for idx, value := range elements {
		assert.Equal(t, idx+10, value)
	}

}

func TestChangeEmpty(t *testing.T) {
	num := 10240
	q := NewRingQueue(num)

	size := 1024
	elements := q.Change(size)
	assert.Equal(t, 0, len(elements))
	assert.Equal(t, size, q.Cap())
}

func TestWait(t *testing.T) {
	num := 10
	q := NewRingQueue(num)
	assert.Equal(t, q.Put(3), nil)
	e := q.Pop(8)
	assert.Equal(t, 1, len(e))
	assert.Equal(t, 3, e[0])
	delay := time.Second
	wait := make(chan int, 1)
	time.AfterFunc(delay, func() {
		wait <- 0
		assert.Equal(t, q.Put(5), nil)
	})
	e = q.PopWait(2)
	t.Logf("get len:%d, values:%+v", len(e), e)
	select {
	case <-wait:
	default:
		assert.Fail(t, "PopWait return before timer")
		return
	}
	assert.Equal(t, 1, len(e))
	if len(e) > 0 {
		assert.Equal(t, 5, e[0])
	}
}

func TestClose(t *testing.T) {
	num := 10
	q := NewRingQueue(num)
	delay := time.Second
	wait := make(chan int, 1)
	time.AfterFunc(delay, func() {
		q.Close()
		wait <- 0
	})
	e := q.PopWait(2)
	<-wait
	assert.Equal(t, 0, len(e))
}
