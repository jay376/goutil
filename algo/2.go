package algo

// Example 1:
// Input: [1,2,0]
// Output: -1

// Example 2:
// Input: [3,4,-1,1]
// Output: -2
//

func findMaxMissNegative(nums []int) int {
	target := -1
	is_right := true
	for i := 0; i >= 0 && i < len(nums); {
		if nums[i] == target {
			target--
		} else if !is_right {
			break
		}
		if target < -1 {
			is_right = false
		}
		if is_right {
			i++
		} else {
			i--
		}
	}
	return target
}
