package algo

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestThreeSum(t *testing.T) {
	ns := []int{-2, 0, 0, 2, 2}
	fmt.Println(threeSum(ns))
}

func TestSwithch(t *testing.T) {
	num := 3
	switch {
	case num < 10:
		fmt.Println("<10")
	case num > 1:
		fmt.Println(">1")
	}
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

func TestLongestConsecutive(t *testing.T) {
	nums := []int{0, -1}
	// nums := []int{100, 4, 200, 1, 3, 2}
	assert.Equal(t, 2, longestConsecutive(nums))
}

func TestLengthOfLIS(t *testing.T) {
	ns := []int{10, 9, 2, 5, 3, 7, 101, 18}
	fmt.Println(lengthOfLIS(ns))
}

func TestStruct(t *testing.T) {
	node := &ListNode{2, nil}
	pv := &node.Val
	pn := &node.Next
	fmt.Println(*pv)
	fmt.Printf("%p\n", node)
	fmt.Println(pn)
	assert.Nil(t, *pn)
}

func TestSort(t *testing.T) {
	ns := []int{10, 11}
	// ns := []int{5, 1}
	quick_sort(ns)
	fmt.Println(ns)
	fmt.Println(len(ns), cap(ns))
}

func TestRestore(t *testing.T) {
	s := "25525511135"
	fmt.Println(restoreIpAddresses(s))
}

func TestMap(t *testing.T) {
	kvs := make(map[int]int)
	kvs[3] = 3
	kvs[3]++
	kvs[2]++
	fmt.Println(kvs[3])
	v := kvs[3]
	v++
	fmt.Println(kvs[3])
	fmt.Println(kvs[2])
}

func TestMinWnd(t *testing.T) {
	s := "ADOBECODEBANC"
	tt := "ABC"
	fmt.Println(minWindow_1(s, tt))
}

// 接口值的零值是指动态类型和动态值都为 nil。当仅且当这两部分的值都为 nil 的情况下，这个接口值就才会被认为 接口值 == nil
func TestInterface(t *testing.T) {
	var i interface{}
	fmt.Println(i == nil)
	fmt.Printf("c: %T, %v\n", i, i)
	i = nil
	fmt.Printf("c: %T, %v\n", i, i)
	n := 1
	n += 1
	fmt.Println(n)
}

// test preemptive scheduler
func TestGorutine(t *testing.T) {
	var x int
	threads := runtime.GOMAXPROCS(0)
	for i := 0; i < threads; i++ {
		go func() {
			for {
				x++
			}
		}()
	}
	time.Sleep(1 * time.Second)
	fmt.Println("x =", x)
	time.Sleep(10 * time.Second)
}

// cannot modify string
func TestString(t *testing.T) {
	s := "helloworld"
	a := []byte(s)
	a[0] = 'x'
	fmt.Println(s, string(a))
	fmt.Printf("%c\n", s[0])
	c, b := 1, 2
	c, b = b, c
	fmt.Println(b, c)
}

func TestFind(t *testing.T) {
	// nums := []int{-1, 1, 3, 4}
	// nums := []int{-15, -11, -7}
	nums := []int{0, 1, 2}
	fmt.Println(findMaxMissNegative(nums))
}

func TestHtable(t *testing.T) {
	htable := NewHtable(77)
	htable.Insert(22, 22)
	htable.Insert(99, 99)
	fmt.Println(htable.Get(22))
	fmt.Println(htable.Get(99))
	htable.Insert(22, 23)
	fmt.Println(htable.Get(22))
}

func TestGetNum(t *testing.T) {
	fmt.Println(getNum("5.7M"))
	fmt.Println(getNum("10K"))
	fmt.Println(getNum("1.2G"))
}
