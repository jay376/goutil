package algo

import (
	"sort"
)

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

// https://leetcode.cn/problems/reverse-nodes-in-k-group/
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

func removeNthFromEnd(head *ListNode, n int) *ListNode {
	length := 0
	for p := head; p != nil; p = p.Next {
		length++
	}
	if length < n {
		return head
	}
	moveNum := length - n
	if moveNum == 0 {
		return head.Next
	}
	p := head
	for ; moveNum > 1; moveNum-- {
		p = p.Next
	}
	p.Next = p.Next.Next
	return head
}

func maxSubArray(nums []int) int {
	max := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i-1] > 0 {
			nums[i] += nums[i-1]
		}
		if nums[i] > max {
			max = nums[i]
		}
	}
	return max
}

func minWindow(s string, t string) string {
	idxs := [256]int{0}
	left, right := 0, 0
	find, need := 0, 0
	for _, c := range t {
		if idxs[int(c)] == 0 {
			need++
		}
		idxs[int(c)] = 1
	}
	duplicate := false
	for idx, c := range s {
		if idxs[int(c)] > 1 && find == need {
			duplicate = true
		}
		if idxs[int(c)] == 1 {
			find++
			idxs[int(c)]++
			if left == 0 {
				left = idx
			} else {
				right = idx
			}
		}
	}
	if duplicate || need != find {
		return ""
	}
	return s[left : right+1]
}

func trap(height []int) int {
	if len(height) < 3 {
		return 0
	}
	var sum int
	beg, end, mid := -1, -1, -1
	for idx := 0; idx < len(height); idx++ {
		val := height[idx]
		if beg == -1 {
			if val > 0 {
				beg = idx
			}
		} else if val > 0 {
			if val < height[beg] {
				if mid == -1 {
					mid = idx
				} else if val >= height[mid] {
					mid = idx
				}
			} else {
				end = idx
			}
		}

		if end != -1 {
			for i := beg + 1; i < end; i++ {
				sum += height[beg] - height[i]
			}
			beg = end
			mid = -1
			end = -1
			idx = beg
		} else if idx == len(height)-1 && mid != -1 {
			for i := beg + 1; i < mid; i++ {
				sum += height[mid] - height[i]
			}
			beg = mid
			mid = -1
			end = -1
			idx = beg
		}
	}
	if end != -1 {
		for idx := beg + 1; idx < end; idx++ {
			sum += height[beg] - height[idx]
		}
	} else if mid != -1 {
		for idx := beg + 1; idx < mid; idx++ {
			sum += height[mid] - height[idx]
		}
	}
	return sum
}

func swapPairs(node *ListNode) *ListNode {
	if node == nil || node.Next == nil {
		return node
	}
	p := node
	phead := &ListNode{}
	var tmp, tail, ret *ListNode
	i := 0
	for p != nil {
		tmp = p
		p = p.Next
		tmp.Next = phead.Next
		phead.Next = tmp
		i++
		if tail == nil {
			tail = tmp
		}
		if i%2 == 0 {
			if ret == nil {
				ret = phead.Next
			}
			phead = tail
			tail = nil
		}
	}
	return ret
}

// https://leetcode-cn.com/problems/reverse-linked-list-ii/
func reverseBetween(head *ListNode, left int, right int) *ListNode {
	phead := &ListNode{}
	phead.Next = head
	prev := phead
	p := head
	var tail *ListNode
	for i := 0; i < right; i++ {
		node := p
		p = p.Next
		if i >= left-1 {
			if i == left-1 {
				tail = node
			}
			node.Next = prev.Next
			prev.Next = node
			if i == right-1 {
				tail.Next = p
				break
			}
		} else {
			prev = node
		}
	}

	return phead.Next
}

// https://leetcode-cn.com/problems/odd-even-linked-list/
// 输入: 2->1->3->5->6->4->7->NULL
// 输出: 2->3->6->7->1->5->4->NULL
func oddEvenList(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	odd, even := head, head.Next
	p1, p2 := odd, even
	for p1 != nil && p2 != nil {
		p1.Next = p2.Next
		if p2.Next == nil {
			break
		}
		p1 = p2.Next
		p2.Next = p1.Next
		p2 = p2.Next
	}
	p1.Next = even
	return head
}

func longestConsecutive(nums []int) int {
	counts := make(map[int]int)
	max := 0
	for _, num := range nums {
		counts[num] = 0
	}
	for _, num := range nums {
		n := 1
		for {
			if c, ok := counts[num]; ok && c < n {
				counts[num] = n
				if n > max {
					max = n
				}
				num++
				n++
			} else {
				break
			}
		}
	}
	return max
}

func quick_sort(nums []int) {
	length := len(nums)
	j := 0
	if length <= 1 {
		return
	}
	stub := nums[0]
	for i := 1; i < length; i++ {
		if nums[i] < stub {
			tmp := nums[j]
			nums[j] = nums[i]
			nums[i] = tmp
			j++
		}
	}
	if j == 0 {
		j++
	} else if j == length-1 {
		j--
	}
	quick_sort(nums[:j])
	quick_sort(nums[j:])
}
