package check

import (
	"fmt"
	"strings"

	"github.com/anhvandev/doc-kit/internal/frontmatter"
)

// frontmatterRequired: thiếu trường theo types.toml là lỗi; created_by
// không phải dk là warning, --strict nâng thành lỗi. File .feature trong docs/
// không đọc được khối # dk: là lỗi (như .html thiếu <!-- dk: --> ở mockup-tokens).
func frontmatterRequired(c *Context) []Finding {
	var out []Finding
	for _, m := range c.Metas {
		if !m.HasFM && strings.HasSuffix(m.Rel, ".feature") {
			out = append(out, finding(m, "frontmatter-required", Error, "thiếu hoặc hỏng khối `# dk:` đầu file (chú thích khác phải cách khối một dòng trống)"))
		}
	}
	for _, m := range c.typed() {
		t := c.Reg[m.Type]
		var missing []string
		for _, k := range t.Required {
			if frontmatter.GetString(m.FM, k) == "" {
				missing = append(missing, k)
			}
		}
		if len(missing) > 0 {
			out = append(out, finding(m, "frontmatter-required", Error, "thiếu trường "+strings.Join(missing, ", ")))
		}
		if by := frontmatter.GetString(m.FM, "created_by"); by != "dk" {
			level := Warning
			if c.Strict {
				level = Error
			}
			out = append(out, finding(m, "frontmatter-required", level, fmt.Sprintf("created_by là %q, không phải dk", by)))
		}
	}
	return out
}

// statusValid: status phải thuộc statuses của loại.
func statusValid(c *Context) []Finding {
	var out []Finding
	for _, m := range c.typed() {
		t := c.Reg[m.Type]
		if m.Status == "" || len(t.Statuses) == 0 {
			continue
		}
		if !contains(t.Statuses, m.Status) {
			out = append(out, finding(m, "status-valid", Error, fmt.Sprintf("status %q không thuộc %s", m.Status, strings.Join(t.Statuses, ", "))))
		}
	}
	return out
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
