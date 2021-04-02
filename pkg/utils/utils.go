package utils

import (
	"regexp"
	"strconv"
	"time"
)

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
