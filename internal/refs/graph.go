// Package refs xây đồ thị liên kết giữa tài liệu: liên kết Markdown tương đối,
// trường source trong frontmatter và mã tài liệu nhắc trong thân.
package refs

import (
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/anhvandev/doc-kit/internal/docs"
	"github.com/anhvandev/doc-kit/internal/render"
)

// Graph là đồ thị có hướng; khóa là đường dẫn tương đối gốc dự án.
type Graph struct {
	Out map[string][]string `json:"out"` // file trỏ đi đâu
	In  map[string][]string `json:"in"`  // ai trỏ đến file
}

// LinkRe bắt liên kết Markdown [..](đích), kể cả ảnh; nhóm 1 là đích.
var LinkRe = regexp.MustCompile(`\]\(([^)\s]+)(?:\s+"[^"]*")?\)`)

// Build xây đồ thị từ tập tài liệu đã quét; root là gốc dự án tuyệt đối.
func Build(root string, metas []docs.Meta) Graph {
	g := Graph{Out: map[string][]string{}, In: map[string][]string{}}
	byPath := map[string]bool{}
	for _, m := range metas {
		byPath[m.Rel] = true
	}
	idRe := idRegexp(metas)
	for _, m := range metas {
		if m.Generated || filepath.Base(m.Rel) == "CHANGELOG-DOCS.md" {
			continue // chỉ mục sinh ra và changelog trỏ đến mọi file, không mang nghĩa
		}
		targets := map[string]bool{}
		for _, lm := range LinkRe.FindAllSubmatch(m.Body, -1) {
			dest := string(lm[1])
			if !render.IsRelativeLink(dest) {
				continue
			}
			dest, _, _ = strings.Cut(dest, "#")
			abs := filepath.Clean(filepath.Join(filepath.Dir(m.Path), filepath.FromSlash(dest)))
			rel, err := filepath.Rel(root, abs)
			if err != nil {
				continue
			}
			if rel = filepath.ToSlash(rel); byPath[rel] && rel != m.Rel {
				targets[rel] = true
			}
		}
		if src, ok := docs.Resolve(metas, m.Source); ok && src.Rel != m.Rel {
			targets[src.Rel] = true
		}
		if idRe != nil {
			for _, im := range idRe.FindAllString(string(m.Body), -1) {
				if t, ok := docs.Resolve(metas, im); ok && t.Rel != m.Rel {
					targets[t.Rel] = true
				}
			}
		}
		for t := range targets {
			g.Out[m.Rel] = append(g.Out[m.Rel], t)
			g.In[t] = append(g.In[t], m.Rel)
		}
	}
	for _, mm := range []map[string][]string{g.Out, g.In} {
		for k := range mm {
			sort.Strings(mm[k])
		}
	}
	return g
}

// idRegexp ghép mọi id đã biết thành một regexp, id dài xếp trước để
// CR-260910-loc-theo không bị bắt ngắn thành CR-260910-loc.
func idRegexp(metas []docs.Meta) *regexp.Regexp {
	var ids []string
	for _, m := range metas {
		if m.ID != "" {
			ids = append(ids, regexp.QuoteMeta(m.ID))
		}
	}
	if len(ids) == 0 {
		return nil
	}
	sort.Slice(ids, func(i, j int) bool { return len(ids[i]) > len(ids[j]) })
	return regexp.MustCompile(`\b(?:` + strings.Join(ids, "|") + `)\b`)
}
