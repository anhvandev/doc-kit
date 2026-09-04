package check

import (
	"fmt"
	"regexp"
	"strings"
)

// envAssignRe khớp dòng KEY=value (kể cả trong khối mã, có hoặc không "export").
var envAssignRe = regexp.MustCompile(`^\s*(?:export\s+)?([A-Z][A-Z0-9_]*)\s*=\s*(.*?)\s*$`)

// placeholderRe khớp trọn giá trị dạng <mô tả>.
var placeholderRe = regexp.MustCompile(`^<[^<>]*>$`)

// envNoSecret: Environment chỉ ghi tên biến và ý nghĩa; dòng KEY=value với
// value không rỗng và không phải placeholder <...> là lỗi. Quét cả khối mã
// (mẫu file cấu hình nằm đó), bỏ chú thích HTML (gợi ý của template). Lớp phụ,
// không thay công cụ quét secret của dự án (secret không theo mẫu KEY=value không bắt).
func envNoSecret(c *Context) []Finding {
	var out []Finding
	for _, m := range c.typed() {
		if m.Type != "environment" {
			continue
		}
		for _, l := range bodyLines(m, false) {
			sm := envAssignRe.FindStringSubmatch(l.text)
			if sm == nil {
				continue
			}
			v := strings.Trim(sm[2], `"'`)
			if v == "" || placeholderRe.MatchString(v) {
				continue
			}
			out = append(out, Finding{File: m.Rel, Line: l.num, Rule: "env-no-secret", Level: Error,
				Msg: fmt.Sprintf("%s có giá trị thật; chỉ ghi placeholder dạng %s=<mô tả>", sm[1], sm[1])})
		}
	}
	return out
}
