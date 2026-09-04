package check

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/anhvandev/doc-kit/internal/refs"
	"github.com/anhvandev/doc-kit/internal/render"
)

// linkBroken: liên kết tương đối tới file không tồn tại; liên kết .md vào docs/html/.
func linkBroken(c *Context) []Finding {
	var out []Finding
	htmlDir := filepath.Join(c.Root, filepath.FromSlash(c.DocsDir), "html")
	for _, m := range c.Metas {
		if m.Generated {
			continue
		}
		offset := len(m.Raw) - len(m.Body) // thân bắt đầu sau frontmatter
		for _, loc := range refs.LinkRe.FindAllSubmatchIndex(m.Body, -1) {
			dest := string(m.Body[loc[2]:loc[3]])
			if !render.IsRelativeLink(dest) {
				continue
			}
			path, _, _ := strings.Cut(dest, "#")
			if path == "" {
				continue
			}
			line := bytes.Count(m.Raw[:offset+loc[0]], []byte("\n")) + 1
			abs := filepath.Clean(filepath.Join(filepath.Dir(m.Path), filepath.FromSlash(path)))
			f := Finding{File: m.Rel, Line: line, Rule: "link-broken", Level: Error}
			if strings.HasSuffix(strings.ToLower(path), ".md") && inside(htmlDir, abs) {
				f.Msg = "liên kết .md trỏ vào " + c.DocsDir + "/html/: " + dest
				out = append(out, f)
				continue
			}
			if _, err := os.Stat(abs); err != nil {
				f.Msg = "liên kết tới file không tồn tại: " + dest
				out = append(out, f)
			}
		}
	}
	return out
}

func inside(parent, child string) bool {
	rel, err := filepath.Rel(parent, child)
	return err == nil && !strings.HasPrefix(rel, "..")
}
