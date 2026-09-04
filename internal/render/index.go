package render

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"path/filepath"
	"sort"
	"strings"

	"github.com/anhvandev/doc-kit/internal/docs"
)

// Index sinh trang docs/html/index.html: nhóm theo type, bỏ file sinh ra
// và file không có frontmatter.
func (o Options) Index(title string, metas []docs.Meta) ([]byte, error) {
	if err := load(); err != nil {
		return nil, err
	}
	groups := map[string][]docs.Meta{}
	for _, m := range metas {
		if m.Generated || !m.HasFM || m.Type == "" {
			continue
		}
		groups[m.Type] = append(groups[m.Type], m)
	}
	types := make([]string, 0, len(groups))
	for t := range groups {
		types = append(types, t)
	}
	sort.Strings(types)
	for _, t := range types {
		g := groups[t]
		sort.Slice(g, func(i, j int) bool { return g[i].Rel < g[j].Rel })
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "<h1>%s</h1>\n", html.EscapeString(title))
	for _, t := range types {
		fmt.Fprintf(&sb, "<h2 id=\"%s\">%s</h2>\n<table>\n<thead><tr><th>id</th><th>title</th><th>status</th><th>owner</th><th>updated</th></tr></thead>\n<tbody>\n", html.EscapeString(t), html.EscapeString(t))
		for _, m := range groups[t] {
			out, err := o.OutPath(m.Path)
			if err != nil {
				return nil, err
			}
			rel, _ := filepath.Rel(o.OutDir, out)
			href := html.EscapeString(filepath.ToSlash(rel))
			id := m.ID
			if id == "" {
				id = filepath.Base(filepath.Dir(m.Path)) + "/" + filepath.Base(m.Path)
			}
			fmt.Fprintf(&sb, "<tr><td class=\"id\"><a href=\"%s\">%s</a></td><td><a href=\"%s\">%s</a></td><td><span class=\"status\">%s</span></td><td>%s</td><td>%s</td></tr>\n",
				href, html.EscapeString(id), href, html.EscapeString(m.Title), html.EscapeString(m.Status), html.EscapeString(m.Owner), html.EscapeString(m.Updated))
		}
		sb.WriteString("</tbody>\n</table>\n")
	}
	var buf bytes.Buffer
	err := pageTmpl.Execute(&buf, pageData{Title: title, Body: template.HTML(sb.String()), CSS: cssText, Index: true})
	return buf.Bytes(), err
}
