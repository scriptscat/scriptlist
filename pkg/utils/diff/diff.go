package diff

import (
	"bytes"
	"fmt"
	"html"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func Diff(src, dst string) string {
	dmp := diffmatchpatch.New()
	fileAdmp, fileBdmp, dmpStrings := dmp.DiffLinesToChars(src, dst)
	diffs := dmp.DiffMain(fileAdmp, fileBdmp, false)
	diffs = dmp.DiffCharsToLines(diffs, dmpStrings)
	diffs = dmp.DiffCleanupSemantic(diffs)

	var buff bytes.Buffer
	for i := 0; i < len(diffs); i++ {
		diff := diffs[i]
		text := strings.Replace(html.EscapeString(diff.Text), "\n", "<br>", -1)
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			_, _ = buff.WriteString("<ins style=\"background:#e6ffe6;\">")
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString("</ins>")
		case diffmatchpatch.DiffDelete:
			// 一直遍历到下一次Equal
			src := ""
			dst := ""
			flag := false
			for n := i; n < len(diffs); n++ {
				diff := diffs[n]
				switch diff.Type {
				case diffmatchpatch.DiffInsert:
					dst += diff.Text
				case diffmatchpatch.DiffDelete:
					src += diff.Text
				case diffmatchpatch.DiffEqual:
					flag = true
				}
				if flag {
					i = n - 1
					break
				}
			}
			diffs := dmp.DiffMain(src, dst, false)
			del := ""
			ins := ""
			for _, diff := range diffs {
				switch diff.Type {
				case diffmatchpatch.DiffInsert:
					ins += fmt.Sprintf("<strong>%s</strong>", strings.Replace(html.EscapeString(diff.Text), "\n", "<br>", -1))
				case diffmatchpatch.DiffDelete:
					del += fmt.Sprintf("<strong>%s</strong>", strings.Replace(html.EscapeString(diff.Text), "\n", "<br>", -1))
				case diffmatchpatch.DiffEqual:
					del += strings.Replace(html.EscapeString(diff.Text), "\n", "<br>", -1)
					ins += strings.Replace(html.EscapeString(diff.Text), "\n", "<br>", -1)
				}
			}

			_, _ = buff.WriteString(fmt.Sprintf("<del style=\"background:#ffe6e6;\">%s</del>", del))
			_, _ = buff.WriteString(fmt.Sprintf("<ins style=\"background:#e6ffe6;\">%s</ins>", ins))
		case diffmatchpatch.DiffEqual:
			_, _ = buff.WriteString("<span>")
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString("</span>")
		}
	}
	return buff.String()
}
