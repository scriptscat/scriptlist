package utils

import (
	"fmt"
	"testing"
)

func TestGetCodeMeta(t *testing.T) {
	meta := `// ==UserScript==
// @name         bilibili自动签到
// @description  一键查看挂载到window上的非原生属性，并注入一个$searchKey函数搜索属性名。
// @namespace    wyz
// @version      1.1.3
// @author       wyz
// @crontab * * once * *
// @debug
// @grant GM_xmlhttpRequest
// @grant GM_notification
// @connect api.bilibili.com
// @connect api.live.bilibili.com
// @cloudCat
// @downloadURL http://www.test.com
// @exportCookie domain=api.bilibili.com
// @exportCookie domain=api.live.bilibili.com
// @updateURL http://www.test.com
// ==/UserScript==
aaa`
	fmt.Println(GetCodeMeta(meta))
}
