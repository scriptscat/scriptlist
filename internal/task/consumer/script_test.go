package consumer

import (
	"testing"
)

func Test_script_parseMatchDomain(t *testing.T) {
	type args struct {
		meta string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"case1", args{"https://www.baidu.com"}, "baidu.com"},
		{"case2", args{"*://*"}, "*"},
		{"case3", args{"https://baidu.com"}, "baidu.com"},
		{"case4", args{"https://*.baidu.com/"}, "baidu.com"},
		{"case5", args{"https://*baidu.com/"}, "baidu.com"},
		{"case6", args{"*"}, "*"},
		{"case7", args{"https://www.sub.baidu.com/"}, "baidu.com"},
		{"case8", args{"https://bbs.tampermonkey.net.cn/"}, "tampermonkey.net.cn"},
		{"case9", args{"https://go.dev/"}, "go.dev"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &script{}
			if got := s.parseMatchDomain(tt.args.meta); got != tt.want {
				t.Errorf("parseMetaDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
