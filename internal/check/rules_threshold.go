package check

import (
	"fmt"
	"strings"
)

// lineThreshold: trên warn_lines là warning, trên max_lines là lỗi. Chỉ áp
// cho tài liệu Markdown trong docs/; plan và report dài là chuyện của plans/,
// tokens.json và mockup .html không phải văn bản đọc theo dòng. Loại có
// warn_lines riêng trong types.toml dùng ngưỡng đó thay cho dk.toml.
func lineThreshold(c *Context) []Finding {
	var out []Finding
	for _, m := range c.Metas {
		if m.Generated || !m.IsMarkdown() || !strings.HasPrefix(m.Rel, c.DocsDir+"/") {
			continue
		}
		warn := c.Cfg.WarnLines
		if t, ok := c.Reg[m.Type]; ok && m.HasFM && t.WarnLines > 0 {
			warn = t.WarnLines
		}
		switch {
		case m.Lines > c.Cfg.MaxLines:
			out = append(out, finding(m, "line-threshold", Error, fmt.Sprintf("%d dòng, vượt max_lines %d", m.Lines, c.Cfg.MaxLines)))
		case m.Lines > warn:
			out = append(out, finding(m, "line-threshold", Warning, fmt.Sprintf("%d dòng, vượt warn_lines %d", m.Lines, warn)))
		}
	}
	return out
}
