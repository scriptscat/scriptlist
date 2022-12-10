package script

import (
	"context"
	"regexp"
	"strings"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
)

// 解析脚本的元数据
func parseCodeMeta(ctx context.Context, scriptCode string) (string, string, error) {
	reg := regexp.MustCompile(`\/\/\s*==UserScript==([\s\S]+?)\/\/\s*==\/UserScript==`)
	ret := reg.FindString(scriptCode)
	if ret == "" {
		return "", "", i18n.NewError(ctx, code.ScriptParseFailed)
	}
	// 移除updateurl和downloadurl
	reg2 := regexp.MustCompile(`(?im)^//\s*@(updateurl|downloadurl)($|\s+(.+?)$)\s+`)
	ret = reg2.ReplaceAllString(ret, "")
	scriptCode = reg.ReplaceAllLiteralString(scriptCode, ret)
	return scriptCode, ret, nil
}

func parseMetaToJson(meta string) map[string][]string {
	reg := regexp.MustCompile(`(?im)^//\s*@(.+?)($|\s+(.+?)$)`)
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
