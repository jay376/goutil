package comm

import (
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8}
	nnums := nums[:0]
	for _, num := range nums {
		if num%2 == 0 {
			nnums = append(nnums, num*2)
		}
	}
	fmt.Println(nums)
	fmt.Println(nnums)
	fmt.Println(nums[1:4])
}
