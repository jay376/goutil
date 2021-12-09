package algo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThreeSum(t *testing.T) {
	ns := []int{-2, 0, 0, 2, 2}
	fmt.Println(threeSum(ns))
}

func TestTrap(t *testing.T) {
	ns := []int{2, 6, 3, 8, 2, 7, 2, 5, 0}
	// ns := []int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}
	fmt.Println(trap(ns))
}

func TestLRU(t *testing.T) {
	cache := Constructor(2)
	cache.Put(1, 1)
	cache.Dump()
	cache.Put(2, 2)
	cache.Dump()
	cache.Get(1)
	cache.Dump()
	cache.Put(3, 3)
	cache.Dump()
	assert.Equal(t, -1, cache.Get(2))
	cache.Put(4, 4)
	assert.Equal(t, -1, cache.Get(1))
}

// ["LRUCache","put","put","get","put","get","put","get","get","get"]
// [[2],[1,1],[2,2],[1],[3,3],[2],[4,4],[1],[3],[4]]
