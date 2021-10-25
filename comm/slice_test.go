package comm

import (
	"errors"
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

func test() (int, error) {
	return 1, errors.New("hello")
}

func TestError(t *testing.T) {
	a, err := test()
	if err != nil {
		fmt.Println(err)
	}
	b, err := test()
	fmt.Printf("%v, %v, %v\n", a, b, err)

	var c int
	fmt.Printf("%v\n", &c)
	c, e := test()
	fmt.Printf("%v, %v\n", &c, e)
}
