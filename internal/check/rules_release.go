package check

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/anhvandev/doc-kit/internal/docs"
)

// noJargon: Release brief và User guide viết cho người dùng cuối không được
// chứa thuật ngữ kỹ thuật trong danh sách [release] jargon của dk.toml
// (mặc định config.DefaultJargon). So cả từ, không phân biệt hoa thường, bỏ
// chú thích HTML và khối mã; mỗi từ báo một lần tại dòng đầu; warning vì tên
// sản phẩm có thể chứa từ đó (người bỏ từ khỏi danh sách, skill không tự xóa).
func noJargon(c *Context) []Finding {
	if len(c.Jargon) == 0 {
		return nil
	}
	var alts []string
	for _, j := range c.Jargon {
		if j = strings.TrimSpace(j); j != "" {
			alts = append(alts, regexp.QuoteMeta(j))
		}
	}
	re := regexp.MustCompile(`(?i)\b(` + strings.Join(alts, "|") + `)\b`)
	var out []Finding
	for _, m := range c.typed() {
		if m.Type != "release-brief" && m.Type != "user-guide" {
			continue
		}
		seen := map[string]bool{}
		for _, l := range bodyLines(m, true) {
			for _, w := range re.FindAllString(l.text, -1) {
				k := strings.ToLower(w)
				if seen[k] {
					continue
				}
				seen[k] = true
				out = append(out, Finding{File: m.Rel, Line: l.num, Rule: "no-jargon", Level: Warning,
					Msg: fmt.Sprintf("từ kỹ thuật %q; viết lại bằng ngôn ngữ người dùng hoặc bỏ từ khỏi [release] jargon trong dk.toml", w)})
			}
		}
	}
	return out
}

type numbered struct {
	text string
	num  int
}

// countFMLines đếm số dòng frontmatter của tài liệu (để đổi số dòng thân sang số dòng file).
func countFMLines(m docs.Meta) int { return docs.CountLines(m.Raw) - docs.CountLines(m.Body) }

// bodyLines trả về các dòng thân tài liệu kèm số dòng trong file, bỏ chú thích
// HTML (một dòng hoặc nhiều dòng); skipCode=true bỏ cả khối mã.
func bodyLines(m docs.Meta, skipCode bool) []numbered {
	var out []numbered
	fmLines := countFMLines(m)
	code, comment := false, false
	for i, l := range strings.Split(string(m.Body), "\n") {
		t := strings.TrimSpace(l)
		if skipCode && (strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~")) {
			code = !code
			continue
		}
		if code {
			continue
		}
		if comment {
			if j := strings.Index(l, "-->"); j >= 0 {
				comment = false
				l = l[j+3:]
			} else {
				continue
			}
		}
		l = commentRe.ReplaceAllString(l, "")
		if j := strings.Index(l, "<!--"); j >= 0 {
			comment = true
			l = l[:j]
		}
		if strings.TrimSpace(l) != "" {
			out = append(out, numbered{l, fmLines + i + 1})
		}
	}
	return out
}
