package comm

import "sort"

func lengthOfLongestSubstring(s string) int {
	length := len(s)
	strmp := [256]int{}
	for i := 0; i < 256; i++ {
		strmp[i] = -1
	}
	max := 0
	l := 0
	r := l
	for ; r < length; r++ {
		last_pos := strmp[int(s[r])]
		if last_pos > -1 && last_pos >= l {
			if r-l > max {
				max = r - l
			}
			l = last_pos + 1
		}
		strmp[int(s[r])] = r
	}
	if r-l > max {
		max = r - l
	}
	return max
}

func firstMissingPositive(nums []int) int {
	length := len(nums)
	for idx, val := range nums {
		for idx+1 != val && val > 0 && val <= length {
			tmp := nums[val-1]
			idx = val - 1
			nums[val-1] = val
			val = tmp
		}
	}

	for idx, val := range nums {
		if idx+1 != val {
			return idx + 1
		}
	}
	return length + 1
}

type ListNode struct {
	Val  int
	Next *ListNode
}

func reverseKGroup(head *ListNode, k int) *ListNode {
	var ret, tail *ListNode
	var front = &ListNode{}
	p := head
	length := 0

	for ; p != nil; p = p.Next {
		length++
	}
	if length < k {
		return head
	}

	p = head
	for idx := 1; idx < length+1; idx++ {
		if tail == nil {
			tail = p
		}
		next := p.Next
		p.Next = front.Next
		front.Next = p
		p = next
		if idx%k == 0 {
			if idx == k {
				ret = front.Next
			}
			if tail != nil {
				front = tail
				tail = p
			}
			if length-idx < k {
				front.Next = p
				break
			}
		}
	}
	return ret
}

func threeSum(ns []int) [][]int {
	sort.Ints(ns)
	res := [][]int{}
	for i := 0; i < len(ns)-1; i++ {
		if ns[i] > 0 {
			break
		}
		if i > 0 && ns[i] == ns[i-1] {
			continue
		}
		for l, r := i+1, len(ns)-1; l < r; {
			if l > i+1 && ns[l] == ns[l-1] {
				l++
				continue
			}
			if r < len(ns)-1 && ns[r] == ns[r+1] {
				r--
				continue
			}

			sum := ns[l] + ns[r] + ns[i]
			switch {
			case sum == 0:
				res = append(res, []int{ns[l], ns[i], ns[r]})
				l++
				r--
			case sum < 0:
				l++
			case sum > 0:
				r--
			}
		}
	}

	return res
}

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	head := &ListNode{}
	c := 0
	r1 := head
	for l1 != nil || l2 != nil {
		value := c
		if l1 != nil {
			value = value + l1.Val
			l1 = l1.Next
		}

		if l2 != nil {
			value = value + l2.Val
			l2 = l2.Next
		}
		node := &ListNode{
			Val: value % 10,
		}
		c = value / 10
		r1.Next = node
		r1 = node
	}
	if c > 0 {
		r1.Next = &ListNode{
			Val: c,
		}
	}
	return head.Next
}
