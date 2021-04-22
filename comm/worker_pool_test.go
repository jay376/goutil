//nolint
package comm

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkerPool(t *testing.T) {
	wp := NewWorkerPool(100)
	wp.Start()

	var buf []int
	var mu sync.Mutex
	num := 100000
	for i := 0; i < num; i++ {
		n := i
		task := func() {
			mu.Lock()
			buf = append(buf, n)
			mu.Unlock()
		}
		wp.Schedule(task, i)
	}

	time.Sleep(5 * time.Second)
	mu.Lock()
	assert.Equal(t, num, len(buf))
	mu.Unlock()
	numSet := make(map[int]struct{})
	for _, num := range buf {
		numSet[num] = struct{}{}
	}

	assert.Equal(t, num, len(numSet))
}
