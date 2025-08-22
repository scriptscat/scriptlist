package subscribe

import (
	"testing"
)

func Test_script_parseMatchDomain(t *testing.T) {
	type args struct {
		meta string
	}
	tests := []struct {
		name  string
		args  args
		want1 string
		want2 string
	}{
		{"case1", args{"https://www.baidu.com"}, "baidu.com", "www.baidu.com"},
		{"case2", args{"*://*"}, "*", "*"},
		{"case3", args{"https://baidu.com"}, "baidu.com", "baidu.com"},
		{"case4", args{"https://*.baidu.com/"}, "baidu.com", "baidu.com"},
		{"case5", args{"https://*baidu.com/"}, "baidu.com", "baidu.com"},
		{"case6", args{"*"}, "*", "*"},
		{"case7", args{"https://www.sub.baidu.com/"}, "baidu.com", "www.sub.baidu.com"},
		{"case8", args{"https://bbs.tampermonkey.net.cn/"}, "tampermonkey.net.cn", "bbs.tampermonkey.net.cn"},
		{"case9", args{"https://go.dev/"}, "go.dev", "go.dev"},
		{"case10", args{"https://mp.weixin.qq.com/wx*"}, "qq.com", "mp.weixin.qq.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Script{}
			if got1, got2 := s.parseMatchDomain(tt.args.meta); got1 != tt.want1 || got2 != tt.want2 {
				t.Errorf("parseMatchDomain() = %v, %v, want %v, %v", got1, got2, tt.want1, tt.want2)
			}
		})
	}
}
