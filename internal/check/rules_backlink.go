package check

import (
	"fmt"
	"strings"

	"github.com/anhvandev/doc-kit/internal/docs"
)

// backlink: loại có source trong required phải trỏ đến tài liệu tồn tại;
// tài liệu trạng thái đã chốt (final trong types.toml) mà không ai trỏ về là
// warning "tài liệu chết".
func backlink(c *Context) []Finding {
	var out []Finding
	referenced := map[string]bool{}
	for _, m := range c.Metas {
		if !m.HasFM || m.Source == "" {
			continue // kể cả loại chưa có trong types.toml, ví dụ file test
		}
		if src, ok := docs.Resolve(c.Metas, m.Source); ok {
			referenced[src.Rel] = true
		}
	}
	for _, m := range c.typed() {
		t := c.Reg[m.Type]
		if contains(t.Required, "source") {
			if m.Source == "" {
				continue // frontmatter-required đã báo
			}
			if _, ok := docs.Resolve(c.Metas, m.Source); !ok {
				out = append(out, finding(m, "backlink", Error, fmt.Sprintf("source %q không trỏ đến tài liệu nào", m.Source)))
			}
		}
		if len(t.Final) > 0 && contains(t.Final, m.Status) && !referenced[m.Rel] {
			out = append(out, finding(m, "backlink", Warning, fmt.Sprintf("tài liệu chết: status %s nhưng không tài liệu nào có source trỏ về", m.Status)))
		}
	}
	return out
}

// specHasTest: feature-spec không còn draft phải có test case trong docs/test/
// có source trỏ về: file .feature (metadata # dk:) hoặc .md loại test-case-table;
// file loại chưa có trong types.toml nằm trong docs/test/ vẫn được tính.
// Checklist giao diện và test report không thay được test case.
func specHasTest(c *Context) []Finding {
	var out []Finding
	testDir := c.DocsDir + "/test/"
	tested := map[string]bool{}
	for _, m := range c.Metas {
		if !m.HasFM || !strings.HasPrefix(m.Rel, testDir) || m.Source == "" {
			continue
		}
		if _, known := c.Reg[m.Type]; !known || strings.HasPrefix(m.Type, "test-case") {
			tested[m.Source] = true
		}
	}
	for _, m := range c.typed() {
		if m.Type == "feature-spec" && m.Status != "draft" && m.ID != "" && !tested[m.ID] {
			out = append(out, finding(m, "spec-has-test", Warning, fmt.Sprintf("status %s nhưng chưa có test case trong %s có source: %s (dk new test-case --from)", m.Status, testDir, m.ID)))
		}
	}
	return out
}
