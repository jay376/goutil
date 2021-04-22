//nolint
package comm

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var ekv = NewExpireKv(context.Background())

func TestPut(t *testing.T) {
	cleanup := func(key string, value interface{}) {
		t.Log(fmt.Sprintf("cleanup item:%+v", value))
	}
	now := time.Now().UnixNano()
	ekv.Put("1", 1, now+int64(3*time.Second), cleanup)
	assert.Equal(t, 1, ekv.Get("1"))
	ekv.Put("2", 2, now+int64(3*time.Second), cleanup)
	assert.Equal(t, 2, ekv.Get("2"))

	assert.Equal(t, nil, ekv.Get("3"))

	time.Sleep(2 * time.Second)
	now = time.Now().UnixNano()
	ekv.Put("2", 22, now+int64(5*time.Second), cleanup)

	time.Sleep(2 * time.Second)
	time.Sleep(2 * time.Second)

	// test expired
	assert.Equal(t, nil, ekv.Get("1"))
	assert.Equal(t, 22, ekv.Get("2"))

	// update not expired
	ekv.Put("2", 32, 0, nil)
	time.Sleep(2 * time.Second)
	assert.Equal(t, 32, ekv.Get("2"))
}

func TestDel(t *testing.T) {
	now := time.Now().UnixNano()
	var ekv = NewExpireKv(context.Background())
	ekv.Put("1", 1, now+int64(5*time.Second), nil)
	assert.Equal(t, 1, ekv.Get("1"))
	ekv.Put("2", 2, now+int64(5*time.Second), nil)
	assert.Equal(t, 2, ekv.Get("2"))

	ekv.Del("1")
	ekv.Del("2")
	assert.Equal(t, nil, ekv.Get("1"))
	assert.Equal(t, nil, ekv.Get("2"))
	tf := func() {
		now := time.Now().UnixNano()
		var mu sync.Mutex
		var cleanValue []int
		cleanup := func(key string, value interface{}) {
			mu.Lock()
			cleanValue = append(cleanValue, value.(int))
			mu.Unlock()
		}

		num := 100000
		for i := 0; i < num; i++ {
			num := i
			go func() {
				key := strconv.Itoa(num)
				ekv.Put(key, num, now+int64(3*time.Second+time.Duration(num)), cleanup)
				assert.Equal(t, ekv.Get(key).(int), num)
			}()
		}
		time.Sleep(2 * time.Second)

		assert.Equal(t, ekv.Len(), num)
		assert.Equal(t, ekv.queue.Len(), num)

		realDelNum := 0
		for i := 0; i < 100; i++ {
			delNum := rand.Intn(num) % 100
			delKey := strconv.Itoa(delNum)
			if v := ekv.Del(delKey); v != nil {
				assert.Equal(t, delNum, v.(int))
				realDelNum++
			}
		}
		assert.Equal(t, realDelNum, num-ekv.Len())

		time.Sleep(3 * time.Second)
		mu.Lock()
		assert.Equal(t, 0, ekv.Len())
		assert.Equal(t, 0, ekv.queue.Len())
		mu.Unlock()
		assert.Equal(t, num-realDelNum, len(cleanValue))
	}

	for i := 0; i < 5; i++ {
		tf()
	}
}
