package march

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLength(t *testing.T) {
	assert.Equal(t, 3, lengthOfLongestSubstring("abcabcbb"))

}

func TestKth(t *testing.T) {
	nums := []int{3, 3, 1, 2, 4, 5, 5, 6}
	// nums := []int{5, 6}
	assert.Equal(t, 4, findKthLargest(nums, 4))

}

func TestMid(t *testing.T) {
	nums1 := []int{0, 0, 0, 0, 0}
	nums2 := []int{-1, 0, 0, 0, 0, 0, 1}
	// nums := []int{5, 6}

	assert.Equal(t, 0, findK(nums1, nums2, 6))
	// assert.Equal(t, 0, findK(nums1, nums2, 7))
}

func TestRepeat(t *testing.T) {
	nums := []int{1, 1, 1, 2, 2, 3}
	length, newNums := clearRepeat(nums)
	assert.Equal(t, length, 5)
	fmt.Println(newNums)

	nums2 := []int{0, 0, 1, 1, 1, 2, 3, 3}
	length, newNums = clearRepeat(nums2)
	assert.Equal(t, length, 7)
	fmt.Println(newNums)
}

func TestPermute(t *testing.T) {
	nums := []int{1, 2, 3}
	ret := permute(nums)
	fmt.Println(ret)
}

func TestCommu(t *testing.T) {
	commu()
}
