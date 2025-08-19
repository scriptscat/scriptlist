package script_entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCode_SRI(t *testing.T) {
	tests := []struct {
		name     string
		code     *Code
		expected string
	}{
		{
			name:     "空代码应返回空字符串",
			code:     &Code{Code: ""},
			expected: "",
		},
		{
			name:     "nil对象应返回空字符串",
			code:     nil,
			expected: "",
		},
		{
			name: "正常代码应返回正确的SHA-512 SRI",
			code: &Code{
				Code: `// ==UserScript==
// @name         Test Script
// @version      1.0.0
// @description  A test script
// ==/UserScript==

console.log('Hello World!');`,
			},
			expected: "sha384-acevngNzZhpc4hArTtbtrv+rpuNu7BTAJzqM8k6pSm+4ljNd3iWWenkRSXIcsR37",
		},
		{
			name: "简单代码字符串",
			code: &Code{
				Code: "console.log('test');",
			},
			expected: "sha384-xvs5/LKScz0YatxcyoqdjZ+pPwaZ2U0z+xZNsaS6SetrbGsfUogeVjbwWIODxNMU",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.code.SRI()

			// 对于空字符串和nil的情况
			if tt.expected == "" && (tt.code == nil || tt.code.Code == "") {
				assert.Equal(t, "", result)
				return
			}

			// 对于有实际代码的情况，验证SRI格式和内容
			if tt.code != nil && tt.code.Code != "" {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
