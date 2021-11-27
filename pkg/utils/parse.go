package utils

import (
	"errors"
	"regexp"
	"strings"

	"github.com/scriptscat/scriptlist/internal/domain/script/entity"
	"github.com/weppos/publicsuffix-go/publicsuffix"
)

func GetCodeMeta(code string) (string, string, int, error) {
	reg := regexp.MustCompile("\\/\\/\\s*==UserScript==([\\s\\S]+?)\\/\\/\\s*==\\/UserScript==")
	ret := reg.FindString(code)
	scriptType := entity.USERSCRIPT_TYPE
	if ret == "" {
		reg = regexp.MustCompile("\\/\\/\\s*==UserScript==([\\s\\S]+?)\\/\\/\\s*==\\/UserScript==")
		ret = reg.FindString(code)
		if ret == "" {
			return "", "", 0, errors.New("错误的格式")
		}
		scriptType = entity.SUBSCRIBE_TYPE
	}
	// 处理
	reg2 := regexp.MustCompile("(?im)^//\\s*@(updateurl|downloadurl)($|\\s+(.+?)$)\\s+")
	ret = reg2.ReplaceAllString(ret, "")
	code = reg.ReplaceAllString(code, ret)
	return code, ret, scriptType, nil
}

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
