package utils

import (
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func StringToInt(i string) int {
	ret, _ := strconv.Atoi(i)
	return ret
}

func StringToInt64(i string) int64 {
	ret, _ := strconv.ParseInt(i, 10, 64)
	return ret
}

func StringToTime(layout string, s string) *time.Time {
	t, err := time.ParseInLocation(layout, s, time.Local)
	if err != nil {
		return nil
	}
	return &t
}

func RegexMatch(content string, command string) []string {
	reg := regexp.MustCompile(command)
	return reg.FindStringSubmatch(content)
}

func StringReverse(s string) string {
	a := []rune(s)
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return string(a)
}

var str = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func GetRandomString(n int) string {
	var result []byte
	for i := 0; i < n; i++ {
		result = append(result, str[rand.Intn(62)])
	}
	return string(result)
}
