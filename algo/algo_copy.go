package algo

import (
	"fmt"
	"strconv"
)

// https://leetcode.cn/problems/longest-substring-without-repeating-characters/
// 输入: s = "abcabcbb"
// 输出: 3
// 解释: 因为无重复字符的最长子串是 "abc"，所以其长度为 3。
func lengthOfLongestSubstring_1(s string) int {
	pos := [256]int{0}
	max := 0
	left := 0
	for i := 0; i < len(s); i++ {
		cur_max := 0
		last_pos := pos[s[i]]
		pos[s[i]] = i + 1
		if last_pos != 0 {
			for j := left; j < last_pos-1; j++ {
				pos[s[j]] = 0
			}
			left = last_pos
		}
		cur_max = i - left + 1
		if cur_max > max {
			max = cur_max
		}
	}

	return max
}

func firstMissingPositive_2(nums []int) int {
	length := len(nums)
	for idx := 0; idx < length; {
		value := nums[idx]
		if value == idx+1 {
			idx++
			continue
		}
		if value > length || value <= 0 {
			nums[idx] = 0
			idx++
		} else {
			tmp := nums[value-1]
			nums[value-1] = value
			nums[idx] = tmp
		}
	}
	miss := 0
	for idx, value := range nums {
		if idx != value-1 {
			miss = idx + 1
		}
	}
	return miss
}

// https://leetcode.cn/problems/longest-increasing-subsequence/
// 输入：nums = [10,9,2,5,3,7,101,18]
// 输出：4
// 解释：最长递增子序列是 [2,3,7,101]，因此长度为 4
func lengthOfLIS(nums []int) int {
	length := len(nums)
	maxs := make([]int, length, length)
	maxs[0] = 1
	max := 0
	for i := 0; i < length; i++ {
		tmp_max := 0
		for j := 0; j < i; j++ {
			if nums[j] < nums[i] && maxs[j] > tmp_max {
				tmp_max = maxs[j]
			}
		}
		maxs[i] = tmp_max + 1
		if maxs[i] > max {
			max = maxs[i]
		}
	}

	return max
}

// https://leetcode-cn.com/problems/odd-even-linked-list/
// 输入: 2->1->3->5->6->4->7->NULL
// 输出: 2->3->6->7->1->5->4->NULL
func oddEvenList_1(head *ListNode) *ListNode {
	num := 1
	odd, even := &ListNode{}, &ListNode{}
	podd, peven := &odd.Next, &even.Next
	for p := head; p != nil; p = p.Next {
		if num%2 == 1 {
			*podd = p
			podd = &p.Next
		} else {
			*peven = p
			peven = &p.Next
		}
		num++
	}
	*podd = even.Next
	*peven = nil
	return odd.Next
}

// https://leetcode.cn/problems/swap-nodes-in-pairs/submissions/
// Input: head = [1,2,3,4]
// Output: [2,1,4,3]
func swapPairs_1(node *ListNode) *ListNode {
	if node == nil || node.Next == nil {
		return node
	}
	tmphead := &ListNode{0, node}
	phead := node.Next
	for node != nil && node.Next != nil {
		p1 := node.Next
		p2 := node.Next.Next
		node.Next = p2
		p1.Next = node
		tmphead.Next = p1
		tmphead = node
		node = p2
	}

	return phead
}

// https://leetcode.cn/problems/reverse-nodes-in-k-group/
func reverseKGroup_1(head *ListNode, k int) *ListNode {
	length := 0
	for p := head; p != nil; p = p.Next {
		length++
	}
	if length < k || k <= 1 {
		return head
	}
	phead := &ListNode{}
	var first, ret *ListNode
	num := 0
	for p := head; p != nil; {
		if first == nil {
			first = p
		}
		next := p.Next
		p.Next = phead.Next
		phead.Next = p
		num++
		if num%k == 0 {
			if num == k {
				ret = phead
			}
			if length-num < k {
				first.Next = next
				break
			}
			phead = first
			first = nil

		}
		p = next
	}
	return ret.Next
}

// https://leetcode.cn/problems/longest-consecutive-sequence/?fileGuid=jhktdryWkkkCTJkv
// 输入：nums = [100,4,200,1,3,2]
// 输出：4
// 解释：最长数字连续序列是 [1, 2, 3, 4]。它的长度为 4
func longestConsecutive_1(nums []int) int {
	kv := make(map[int]int)
	for _, num := range nums {
		kv[num] = 1
	}
	max := 0
	for num := range kv {
		tmp_max := 1
		delete(kv, num)
		n := num + 1
		for {
			if _, ok := kv[n]; ok {
				delete(kv, n)
				n++
				tmp_max++
			} else {
				break
			}
		}
		n = num - 1
		for {
			if _, ok := kv[n]; ok {
				delete(kv, n)
				n--
				tmp_max++
			} else {
				break
			}
		}
		if tmp_max > max {
			max = tmp_max
		}
	}
	return max
}

// https://leetcode-cn.com/problems/reverse-linked-list-ii/
func reverseBetween_1(head *ListNode, left int, right int) *ListNode {
	phead := &ListNode{0, head}
	num := 0
	var prev, tail *ListNode
	for p := phead; p != nil; {
		num++
		next := p.Next
		if num < left+1 {
			prev = p
		}
		if num == left+1 {
			p.Next = nil
			tail = p
		}
		if num > left+1 && num <= right+1 {
			tmp := prev.Next
			prev.Next = p
			p.Next = tmp
		}
		if num == right+1 {
			tail.Next = next
			break
		}
		p = next
	}

	return phead.Next
}

// https://leetcode.cn/problems/reverse-nodes-in-k-group/
func reverseKGroup_2(head *ListNode, k int) *ListNode {
	length := 0
	for p := head; p != nil; p = p.Next {
		length++
	}
	if length < k {
		return head
	}
	num := 0
	var tail, ret *ListNode
	phead := &ListNode{}
	for p := head; p != nil; {
		num++
		next := p.Next
		if tail == nil {
			tail = p
		}
		tmp := phead.Next
		phead.Next = p
		p.Next = tmp
		if num%k == 0 {
			if num == k {
				ret = phead.Next
			}
			phead = tail
			tail = nil
			if length-num < k {
				phead.Next = next
				break
			}
		}
		p = next
	}

	return ret
}

// https://leetcode.cn/problems/minimum-window-substring/submissions/
// 输入：s = "ADOBECODEBANC", t = "ABC"
// 输出："BANC"
// 滑动窗口
func minWindow_1(s string, t string) string {
	needs := make(map[byte]int)
	for i := 0; i < len(t); i++ {
		needs[t[i]]++
	}
	min := len(s)
	l := 0
	ret := ""
	need_cnt := len(t)
	for r := 0; r < len(s); r++ {
		c := s[r]
		if needs[c] > 0 {
			need_cnt--
		}
		needs[c]--
		if need_cnt == 0 {
			for ; l < r && needs[s[l]] < 0; l++ {
				needs[s[l]]++
			}
			if r-l+1 <= min {
				min = r - l + 1
				ret = s[l : r+1]
			}
			fmt.Println(l, r)
			need_cnt++
			needs[s[l]]++
			l++
		}
	}
	return ret
}

func trap_1(height []int) (ans int) {
	left, right := 0, len(height)-1
	leftMax, rightMax := 0, 0
	for left < right {
		leftMax = max(leftMax, height[left])
		rightMax = max(rightMax, height[right])
		if height[left] < height[right] {
			ans += leftMax - height[left]
			left++
		} else {
			ans += rightMax - height[right]
			right--
		}
	}
	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func match(s string) bool {
	if s[0] == '0' && len(s) > 1 {
		return false
	}
	n, _ := strconv.Atoi(s)
	return n <= 255
}

// 输入：s = "25525511135"
// 输出：["255.255.11.135","255.255.111.35"]
func restoreIpAddresses(s string) []string {
	res := []string{}
	pos := [4]int{0, -1, -1, -1}
	back := false
	i := 1
	for {
		if i == 0 {
			break
		}
		if back {
			pos[i]++

		} else {
			pos[i] = pos[i-1] + 1
		}
		back = false
		if pos[i] > len(s)-1 {
			i--
			back = true
			continue
		}
		if match(s[pos[i-1]:pos[i]]) {
			if i == 3 {
				if match(s[pos[i]:]) {
					tmp := fmt.Sprintf("%s.%s.%s.%s", s[0:pos[1]], s[pos[1]:pos[2]], s[pos[2]:pos[3]], s[pos[3]:])
					res = append(res, tmp)
				}
				back = true
			} else {
				i++
			}
			continue
		} else {
			i--
			back = true
			continue
		}
	}
	return res
}
