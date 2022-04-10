package utils

import (
	"encoding/json"
	"math/rand"
	"regexp"
	"strconv"
	"time"
	"unsafe"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func ErrFunc(funcs ...func() error) error {
	for _, v := range funcs {
		if err := v(); err != nil {
			return err
		}
	}
	return nil
}

func Errs(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
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

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandString(n int, stype int) string {
	b := make([]byte, n)
	l := 10 + (stype * 24)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < l {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func MarshalJson(h interface{}) string {
	ret, _ := json.Marshal(h)
	return string(ret)
}

func MarshalJsonByte(h interface{}) []byte {
	ret, _ := json.Marshal(h)
	return ret
}
