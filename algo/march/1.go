// byd
package march

func clearRepeat(nums []int) (int, []int) {
	if len(nums) <= 2 {
		return len(nums), nums
	}

	cnt := 1
	pos := 1

	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[pos-1] {
			nums[pos] = nums[i]
			pos++
			cnt = 1
		} else {
			cnt++
			if cnt == 2 {
				nums[pos] = nums[i]
				pos++
			}
		}
	}

	return pos, nums[:pos]
}
