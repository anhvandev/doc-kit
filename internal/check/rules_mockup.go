package check

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/anhvandev/doc-kit/internal/frontmatter"
)

var (
	styleBlockRe = regexp.MustCompile(`(?is)<style[^>]*>.*?</style>`)
	styleAttrRe  = regexp.MustCompile(`(?is)\bstyle\s*=\s*(?:"[^"]*"|'[^']*')`)
	atRuleRe     = regexp.MustCompile(`@[a-zA-Z-]+[^{;]*`) // prelude của @media, @container...: px ở đây không thay bằng biến được
	hexRe        = regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)
	pxRe         = regexp.MustCompile(`\b\d+(?:\.\d+)?px\b`)
)

// mockupTokens: mockup HTML không có giá trị màu hex hay khoảng cách px gõ tay
// trong <style> hoặc style=""; mọi giá trị phải qua biến của tokens.css.
// File .html trong docs/ thiếu khối <!-- dk: --> là lỗi và vẫn bị lint, để
// không lách quy tắc bằng cách bỏ metadata. Mockup có external (Figma) không
// phải HTML do dk quản, bỏ qua.
func mockupTokens(c *Context) []Finding {
	var out []Finding
	for _, m := range c.Metas {
		if !strings.HasSuffix(strings.ToLower(m.Rel), ".html") || !strings.HasPrefix(m.Rel, c.DocsDir+"/") {
			continue
		}
		if !m.HasFM {
			out = append(out, finding(m, "mockup-tokens", Error, "thiếu khối <!-- dk: ... --> đầu file; tạo mockup bằng dk new mockup"))
		} else if m.Type != "mockup" || frontmatter.GetString(m.FM, "external") != "" {
			continue
		}
		offset := len(m.Raw) - len(m.Body)
		spans := append(styleBlockRe.FindAllIndex(m.Body, -1), styleAttrRe.FindAllIndex(m.Body, -1)...)
		for _, sp := range spans {
			seg := atRuleRe.ReplaceAllFunc(m.Body[sp[0]:sp[1]], func(b []byte) []byte { return bytes.Repeat([]byte(" "), len(b)) })
			for _, re := range []*regexp.Regexp{hexRe, pxRe} {
				for _, loc := range re.FindAllIndex(seg, -1) {
					line := bytes.Count(m.Raw[:offset+sp[0]+loc[0]], []byte("\n")) + 1
					out = append(out, Finding{File: m.Rel, Line: line, Rule: "mockup-tokens", Level: Error,
						Msg: fmt.Sprintf("giá trị gõ tay %q; dùng biến từ tokens.css", seg[loc[0]:loc[1]])})
				}
			}
		}
	}
	return out
}
