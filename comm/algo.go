package comm

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
