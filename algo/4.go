package algo

import (
	"strconv"
)

func getNum(s string) int {
	length := len(s)
	if length == 0 {
		return -1
	}
	tail := s[len(s)-1]
	str := ""
	if tail < 'A' || tail > 'Z' && length >= 2 {
		str = s[0:length]
	}
	if tail > 'A' && tail < 'Z' {
		str = s[0 : length-1]
	}
	if str == "" {
		return -1
	}
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return -1
	}
	var ret float64
	max := 0xFFFFFFFF
	if f > float64(max) {
		return -1
	}
	switch s[len(s)-1] {
	case 'K':
		ret = f * 1024
	case 'M':
		ret = f * 1024 * 1024
	case 'G':
		ret = f * 1024 * 1024 * 1024
	default:
		ret = f
	}

	if ret > float64(max) {
		return -1
	}

	return int(ret)
}
