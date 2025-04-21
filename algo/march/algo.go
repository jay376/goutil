package march

import (
	"fmt"
	"math"
)

/*
*
  - Definition for singly-linked list.
*/
type ListNode struct {
	Val  int
	Next *ListNode
}

// https://leetcode.cn/problems/reverse-nodes-in-k-group/
func reverseKGroup(head *ListNode, k int) *ListNode {
	if head == nil || k <= 1 {
		return head
	}

	p, pk := head, head
	ret := &ListNode{}
	phead := &ret.Next
	for i := 1; pk != nil; i++ {
		if i%k == 0 {
			tmp := &ListNode{}
			tmpHead := &p.Next
			for p != pk {
				node := p
				p = p.Next
				node.Next = tmp.Next
				tmp.Next = node
			}
			pk = p.Next
			p.Next = tmp.Next
			tmp.Next = p
			*phead = tmp.Next
			phead = tmpHead
			p = pk
		} else {
			pk = pk.Next
		}
	}
	*phead = p
	return ret.Next
}

// https://leetcode.cn/problems/qJnOS7/
// 输入：text1 = "abcde", text2 = "ace"
// 输出：3
// 解释：最长公共子序列是 "ace" ，它的长度为 3 。

func longestCommonSubsequence(text1 string, text2 string) int {
	len1 := len(text1)
	len2 := len(text2)
	s := make([][]int, len1+1)
	for i := range s {
		s[i] = make([]int, len2+1)
	}
	// s := [len1 + 1][len2 + 1]int{}
	for i := len1 - 1; i >= 0; i-- {
		for j := len2 - 1; j >= 0; j-- {
			if text1[i] == text2[j] {
				s[i][j] = s[i+1][j+1] + 1
			} else {
				s[i][j] = max(s[i+1][j], s[i][j+1])
			}
		}
	}
	return s[0][0]
}

// https://leetcode.cn/problems/longest-substring-without-repeating-characters/
// 输入: s = "pwwkew"
// 输出: 3
// 解释: 因为无重复字符的最长子串是 "wke"，所以其长度为 3。
//
//	请注意，你的答案必须是 子串 的长度，"pwke" 是一个子序列，不是子串。
func lengthOfLongestSubstring(s string) int {
	strs := [256]int{}
	for i := 0; i < 256; i++ {
		strs[i] = -1
	}
	ret := 0
	i := 0
	j := 0
	for ; j < len(s); j++ {
		if strs[s[j]] >= i {
			ret = max(ret, j-i)
			i = strs[s[j]] + 1
		}
		strs[s[j]] = j
	}
	return max(ret, j-i)
}

// https://leetcode.cn/problems/add-two-numbers/
// 输入：l1 = [9,9,9,9,9,9,9], l2 = [9,9,9,9]
// 输出：[8,9,9,9,0,0,0,1]
func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	head := &ListNode{}
	pnext := &head.Next
	c := 0
	p1, p2 := l1, l2
	for p1 != nil || p2 != nil {
		val := c
		if p1 != nil {
			val += p1.Val
			p1 = p1.Next
		}

		if p2 != nil {
			val += p2.Val
			p2 = p2.Next
		}
		c = val / 10
		*pnext = &ListNode{val % 10, nil}
		pnext = &(*pnext).Next
	}
	if c > 0 {
		*pnext = &ListNode{c, nil}
	}
	return head.Next
}

// https://leetcode.cn/problems/kth-largest-element-in-an-array/?envType=problem-list-v2&envId=array
// 输入: [3,2,1,5,6,4], k = 2
// 输出: 5
func findKthLargest(nums []int, k int) int {
	for {
		length := len(nums)
		if length == 1 {
			return nums[0]
		}

		pivot := nums[length-1]
		i := 0
		for j := i; j < length; j++ {
			if nums[j] > pivot {
				nums[i], nums[j] = nums[j], nums[i]
				i++
			}
		}
		eq := i
		for j := eq; j < length; j++ {
			if nums[j] == pivot {
				nums[eq], nums[j] = nums[j], nums[eq]
				eq++
			}
		}
		if i >= k {
			nums = nums[:i]
		} else if eq >= k {
			return pivot
		}

		if eq < k {
			nums = nums[eq:]
			k -= eq
		}
		// fmt.Println(nums)
		// fmt.Println(k)

	}
}

// https://leetcode.cn/problems/longest-increasing-subsequence/?envType=problem-list-v2&envId=array
// 输入：nums = [10,9,2,5,3,7,101,18]
// 输出：4
// 解释：最长递增子序列是 [2,3,7,101]，因此长度为 4 。
func lengthOfLIS(nums []int) int {
	dp := make([]int, len(nums))
	ret, length := 0, len(nums)
	if length == 0 {
		return 0
	}
	for i := length - 1; i >= 0; i-- {
		tmp := 1
		for j := i + 1; j < length; j++ {
			if nums[i] < nums[j] {
				tmp = max(tmp, 1+dp[j])
			}
		}
		dp[i] = tmp
		ret = max(ret, tmp)
	}

	return ret
}

// https://leetcode.cn/problems/median-of-two-sorted-arrays/?envType=problem-list-v2&envId=array
// 输入：nums1 = [1,2], nums2 = [3,4]
// 输出：2.50000
// 解释：合并数组 = [1,2,3,4] ，中位数 (2 + 3) / 2 = 2.5
func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	l1, l2 := len(nums1), len(nums2)
	total := l1 + l2
	m := total / 2

	if total%2 == 1 {
		return float64(findK(nums1, nums2, m+1))
	}
	return float64(findK(nums1, nums2, m)+findK(nums1, nums2, m+1)) / 2.0
}

func findK(nums1 []int, nums2 []int, k int) int {
	for {
		l1, l2 := len(nums1), len(nums2)
		if l1 == 0 {
			return nums2[k-1]
		}
		if l2 == 0 {
			return nums1[k-1]
		}
		if k == 1 {
			return min(nums1[0], nums2[0])
		}

		m1 := min(k/2, l1) - 1
		m2 := min(k/2, l2) - 1

		if nums1[m1] <= nums2[m2] {
			nums1 = nums1[m1+1:]
			k -= (m1 + 1)
		} else {
			nums2 = nums2[m2+1:]
			k -= (m2 + 1)
		}
		// fmt.Println(nums1)
		// fmt.Println(nums2)
		// fmt.Println(k)
	}
}

// https://leetcode.cn/problems/first-missing-positive/?envType=problem-list-v2&envId=array
// 输入：nums = [3,4,-1,1]
// 输出：2
// 解释：1 在数组中，但 2 没有。
func firstMissingPositive(nums []int) int {
	length := len(nums)
	for i := 0; i < length; {
		val := nums[i]
		if val > 0 && val < length && val != i+1 && nums[val-1] != val {
			nums[i], nums[val-1] = nums[val-1], nums[i]
		} else {
			i++
		}
	}
	ret := -1
	for i, num := range nums {
		if num != i+1 {
			ret = i + 1
			break
		}
	}
	if ret == -1 {
		ret = length + 1
	}
	return ret
}

// https://leetcode.cn/problems/longest-consecutive-sequence/?envType=problem-list-v2&envId=array
// 输入：nums = [100,4,200,1,3,2]
// 输出：4
// 解释：最长数字连续序列是 [1, 2, 3, 4]。它的长度为 4。
// func longestConsecutive(nums []int) int {

// }

// // https://leetcode.cn/problems/n-queens/description/?envType=problem-list-v2&envId=array
// // 输入：n = 4
// // 输出：[[".Q..","...Q","Q...","..Q."],["..Q.","Q...","...Q",".Q.."]]
// func solveNQueens(n int) []int {
// 	s := make([]int, n)
// 	step := 1
// 	i := 0
// 	for {
// 		if i == 0 && step = 0 {
// 			break
// 		}
// 		if step == 1 {

// 		}
// 	}
// }

// https://leetcode.cn/problems/container-with-most-water/?envType=study-plan-v2&envId=top-100-liked
// 输入：[1,8,6,2,5,4,8,3,7]
// 输出：49
func maxArea(height []int) int {
	ret := 0
	for i, j := 0, len(height)-1; i < j; {
		ret = max(ret, (j-i)*min(height[i], height[j]))
		if height[i] < height[j] {
			i++
		} else {
			j--
		}
	}

	return ret
}

func longestConsecutive(nums []int) int {
	cnts := make(map[int]int)
	ret := 0
	for _, num := range nums {
		cnts[num] = 1
	}

	for _, num := range nums {
		i := 1
		for {
			cnts[num+i]--
			v := cnts[num+i]
			if v != 0 {
				break
			}
			i++
		}
		j := 1
		for {
			cnts[num-j]--
			v := cnts[num-j]
			if v != 0 {
				break
			}
			j++
		}
		ret = max(ret, i+j-1)
	}
	return ret
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// https://leetcode.cn/problems/binary-tree-maximum-path-sum/?envType=study-plan-v2&envId=top-100-liked
func maxPathSum(root *TreeNode) int {
	ret := math.MinInt
	var maxPath func(node *TreeNode) int
	maxPath = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		lgain := max(maxPath(node.Left), 0)
		rgain := max(maxPath(node.Right), 0)
		ret = max(ret, node.Val+lgain+rgain)
		return max(lgain, rgain) + node.Val
	}

	maxPath(root)
	return ret
}

// https://leetcode.cn/problems/linked-list-cycle-ii/?envType=study-plan-v2&envId=top-100-liked
func detectCycle(head *ListNode) *ListNode {
	n := 0
	if head == nil {
		return nil
	}

	p1, p2 := head, head
	for {
		if p1.Next == nil || p2.Next == nil || p2.Next.Next == nil {
			return nil
		}
		p1 = p1.Next
		p2 = p2.Next.Next
		n++
		if p1 == p2 {
			break
		}
	}
	p1 = head
	for p1 != p2 {
		p1 = p1.Next
		p2 = p2.Next
	}
	return p1
}

// 回退法
// https://leetcode.cn/problems/permutations/?envType=study-plan-v2&envId=top-100-liked
// 输入：nums = [1,2,3]
// 输出：[[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]
func permute(nums []int) [][]int {
	locs := make([]int, len(nums))
	ret := make([][]int, 0, len(nums))
	cur, n := 0, len(nums)
	numSet := make(map[int]struct{}, n)
	if n == 1 {
		ret = append(ret, nums)
		return ret
	}
	for {

		if cur == 0 && locs[cur] >= n {
			break
		}
		if locs[cur] >= n {
			locs[cur] = 0
			cur--
			delete(numSet, locs[cur])
			locs[cur]++
			continue
		}

		if _, ok := numSet[locs[cur]]; ok {
			locs[cur]++
			continue
		}
		numSet[locs[cur]] = struct{}{}
		if cur == n-1 {
			tmp := make([]int, n)
			fmt.Println(locs)
			for i, v := range locs {
				tmp[i] = nums[v]
			}
			ret = append(ret, tmp)
			delete(numSet, locs[cur])
			locs[cur] = 0
			cur--
			delete(numSet, locs[cur])
			locs[cur]++
		} else {
			cur++
		}
	}
	return ret
}
