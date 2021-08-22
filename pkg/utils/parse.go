package utils

import (
	"regexp"
	"strings"

	"github.com/weppos/publicsuffix-go/publicsuffix"
)

func ParseMetaToJson(meta string) map[string][]string {
	reg := regexp.MustCompile("(?im)^//\\s*@(.+?)($|\\s+(.+?)$)")
	list := reg.FindAllStringSubmatch(meta, -1)
	ret := make(map[string][]string)
	for _, v := range list {
		v[1] = strings.ToLower(v[1])
		if _, ok := ret[v[1]]; !ok {
			ret[v[1]] = make([]string, 0)
		}
		ret[v[1]] = append(ret[v[1]], strings.TrimSpace(v[3]))
	}
	return ret
}

func ParseMetaDomain(meta string) string {
	reg := regexp.MustCompile("(.*?://|^)(.*?)(/|$)")
	ret := reg.FindStringSubmatch(meta)
	if len(ret) == 0 || ret[2] == "" || len(ret[2]) < 3 {
		return ""
	}
	if ret[2][0] == '*' {
		ret[2] = ret[2][1:]
	}
	if ret[2][0] == '.' {
		ret[2] = ret[2][1:]
	}
	domain, err := publicsuffix.Domain(ret[2])
	if err != nil {
		return ""
	}
	if domain[0] == '*' {
		return ""
	}
	return domain
}
