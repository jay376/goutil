package comm

import (
	"hash/fnv"
	"log"
	"math/rand"
	"net"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var regCompiled *regexp.Regexp = regexp.MustCompile("^[1]([3-9])[0-9]{9}$")

func init() {
	rand.Seed(time.Now().Unix())
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandStringBytes ...
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// GetIntranetIP 内网ip
func GetIntranetIP() string {
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
	}
	return ""
}

// GetLocalFreePort ...
func GetLocalFreePort() (int, error) {
	Addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", Addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// MustReplaceStringByEnv 'abc ${key:defauleValue} def' ->'abc defaultValue def'.
// The function only can be used when program init.
func MustReplaceStringByEnv(oldStr string) string {
	regStr := "\\${(\\w+?)(:(\\w+?))?}"
	regCompiled := regexp.MustCompile(regStr)

	placeholders := regCompiled.FindAllStringSubmatch(oldStr, -1)

	for i := 0; i < len(placeholders); i++ { //一级是 整个表达式的匹配 集合，二级是单个匹配内的各个括号，但二级的第一个是表达式匹配的本身。
		if len(placeholders[i]) != 4 { //[${kk:vv} kk :vv vv]
			log.Fatal("len(submatch[i]) must be 4")
		}

		key := placeholders[i][1]
		value := placeholders[i][3]

		env := os.Getenv(key)

		if env != "" {
			value = env
		}

		//待添加。获取传入的参数，其value优先级最高

		if value == "" {
			log.Fatal("key " + key + " cannot be empty")
		}

		oldStr = strings.ReplaceAll(oldStr, placeholders[i][0], value)
	}

	return oldStr
}

// ValidatePhoneNumber the phone number should match "^[1]([3-9])[0-9]{9}$"
func ValidatePhoneNumber(phoneNumber string) bool {
	//regStr := "^[1](([3][0-9])|([4][5-9])|([5][0-3,5-9])|([6][5,6])|([7][0-8])|([8][0-9])|([9][1,8,9]))[0-9]{8}$"
	//"^[1]([3-9])[0-9]{9}$"
	return regCompiled.MatchString(phoneNumber)
}

// Hash ...
// nolint
func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// Unique ...
func Unique(strs []string) []string {
	tmp := make(map[string]struct{}, len(strs))
	for _, str := range strs {
		tmp[str] = struct{}{}
	}
	ret := make([]string, len(tmp))
	idx := 0
	for key := range tmp {
		ret[idx] = key
		idx++
	}

	sort.Strings(ret)
	return ret
}

//IntsToString -
func IntsToString(nums []int64) string {
	ret := []string{}
	for _, num := range nums {
		ret = append(ret, strconv.FormatInt(num, 10))
	}
	return strings.Join(ret, ",")
}
