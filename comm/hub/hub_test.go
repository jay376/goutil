package hub

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBtreeStorage(t *testing.T) {
	kvs := map[string][]byte{
		"1":      []byte("sdfjsdfj"),
		"sdfsd":  []byte("sdfjsdfj"),
		"swerwe": []byte("werw222"),
		"werwer": []byte("sdfjsdfj"),
		"wwww":   []byte("werwe"),
		"xxxxx":  []byte("xxxx"),
	}

	st := newBtreeStore()
	for k, v := range kvs {
		_, err := st.Put(k, v)
		assert.Nil(t, err)
	}

	for k, v := range kvs {
		value, _, _ := st.Get(k)
		assert.Equal(t, value, v)
	}

	rangeKvs, _, _ := st.Range("w", "wyz")
	for _, kv := range rangeKvs {
		assert.Equal(t, kv.Value, kvs[kv.Key])
	}
}

func TestHub(t *testing.T) {
	hub := NewHub(context.Background())
	var num int32
	seed := 100
	gorouteNum := 100
	var wg sync.WaitGroup
	generater := func(n int) {
		for i := 0; i < seed; i++ {
			num := atomic.AddInt32(&num, 1)
			path := fmt.Sprintf("/%v/value/%v", n, num)
			time.Sleep(5 * time.Millisecond)
			assert.Nil(t, hub.Put(path, []byte(strconv.Itoa(int(num)))))
		}
		wg.Done()
	}

	valueMap := make(map[string]string)
	var mu sync.Mutex
	for i := 1; i <= gorouteNum*seed; i++ {
		value := strconv.Itoa(i)
		mu.Lock()
		valueMap[value] = value
		mu.Unlock()
	}

	watch := func(n int) {
		path := fmt.Sprintf("/%v/value", n)
		t.Log(path)
		watcher, err := hub.WatchWithPreffix(path)
		assert.Nil(t, err)
		num := 0
		for {
			evs := <-watcher.Events
			t.Log(len(evs))
			for _, ev := range evs {
				num++
				assert.True(t, strings.HasPrefix(ev.Key, path))
				assert.Equal(t, Put, ev.Op)
				mu.Lock()
				delete(valueMap, string(ev.Value))
				mu.Unlock()
			}
			if num == seed {
				break
			}
		}
		wg.Done()
	}
	wg.Add(gorouteNum * 2)
	for i := 0; i < gorouteNum; i++ {
		go generater(i)
		go watch(i)
	}
	wg.Wait()
	assert.Equal(t, 0, len(valueMap))
}
