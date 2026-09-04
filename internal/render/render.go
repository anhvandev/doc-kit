// Package render đổi tài liệu Markdown sang HTML một file tự chứa: CSS nhúng,
// Mermaid nhúng khi cần, metadata từ frontmatter, liên kết .md đổi thành .html.
package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"gopkg.in/yaml.v3"

	"github.com/anhvandev/doc-kit/assets"
	"github.com/anhvandev/doc-kit/internal/docs"
)

// Options là thư mục nguồn và đích, đều tuyệt đối.
type Options struct {
	DocsDir string
	OutDir  string // thường là <DocsDir>/html
}

// Page là kết quả render một file.
type Page struct {
	Out        string `json:"out"` // đường dẫn HTML tuyệt đối
	HasMermaid bool   `json:"mermaid"`
	HTML       []byte `json:"-"`
}

type metaRow struct{ Key, Value string }

// tocItem là một dòng mục lục: tiêu đề cấp 2 hoặc 3 và id neo goldmark đã gán.
type tocItem struct {
	Level int
	ID    string
	Text  string
}

// tocMin là số tiêu đề tối thiểu để hiện mục lục; ít hơn thì trang tự đủ ngắn.
const tocMin = 3

type pageData struct {
	Title      string
	Meta       []metaRow
	TOC        []tocItem
	Body       template.HTML
	CSS        template.CSS
	HasMermaid bool
	Mermaid    template.JS
	Index      bool
}

var (
	loadOnce sync.Once
	pageTmpl *template.Template
	cssText  template.CSS
	jsText   template.JS
	loadErr  error
)

func load() error {
	loadOnce.Do(func() {
		var b []byte
		if b, loadErr = fs.ReadFile(assets.FS, "html/page.html"); loadErr != nil {
			return
		}
		if pageTmpl, loadErr = template.New("page").Parse(string(b)); loadErr != nil {
			return
		}
		if b, loadErr = fs.ReadFile(assets.FS, "html/style.css"); loadErr != nil {
			return
		}
		cssText = template.CSS(b)
		if b, loadErr = fs.ReadFile(assets.FS, "html/mermaid.min.js"); loadErr != nil {
			return
		}
		// Bản IIFE không chứa "</script"; thay phòng khi nâng phiên bản.
		js := strings.ReplaceAll(string(b), "</script", "<\\/script")
		if b, loadErr = fs.ReadFile(assets.FS, "html/MERMAID-LICENSE.txt"); loadErr != nil {
			return
		}
		// html/template bỏ chú thích HTML, nên giấy phép MIT đi kèm dưới dạng chú thích JS.
		jsText = template.JS("/*! mermaid.min.js https://github.com/mermaid-js/mermaid\n" + strings.TrimSpace(string(b)) + "\n*/\n" + js)
	})
	return loadErr
}

// OutPath trả về đường dẫn HTML đích cho file Markdown src (tuyệt đối, trong DocsDir).
func (o Options) OutPath(src string) (string, error) {
	rel, err := filepath.Rel(o.DocsDir, src)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("%s nằm ngoài %s", src, o.DocsDir)
	}
	return filepath.Join(o.OutDir, strings.TrimSuffix(rel, filepath.Ext(rel))+".html"), nil
}

// Render đổi một tài liệu sang HTML; không ghi file.
func (o Options) Render(m docs.Meta) (Page, error) {
	if err := load(); err != nil {
		return Page{}, err
	}
	out, err := o.OutPath(m.Path)
	if err != nil {
		return Page{}, err
	}
	nr := &nodeRenderer{}
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID(), parser.WithASTTransformers(
			util.Prioritized(&linkRewriter{opts: o, srcDir: filepath.Dir(m.Path), outDir: filepath.Dir(out)}, 100))),
		goldmark.WithRendererOptions(renderer.WithNodeRenderers(util.Prioritized(nr, 100))),
	)
	doc := md.Parser().Parse(text.NewReader(m.Body))
	nr.behaviorTable = findBehaviorTable(doc, m.Body)
	var body bytes.Buffer
	if err := md.Renderer().Render(&body, m.Body, doc); err != nil {
		return Page{}, err
	}
	title := m.Title
	if title == "" {
		title = firstHeading(doc, m.Body)
	}
	if title == "" {
		title = filepath.Base(m.Path)
	}
	data := pageData{Title: title, Meta: metaRows(m.FM), TOC: toc(doc, m.Body), Body: template.HTML(body.String()), CSS: cssText, HasMermaid: nr.hasMermaid}
	if nr.hasMermaid {
		data.Mermaid = jsText
	}
	var buf bytes.Buffer
	if err := pageTmpl.Execute(&buf, data); err != nil {
		return Page{}, err
	}
	return Page{Out: out, HasMermaid: nr.hasMermaid, HTML: buf.Bytes()}, nil
}

// metaRows lấy frontmatter theo đúng thứ tự khóa; giá trị không phải scalar in dạng YAML một dòng.
func metaRows(fm *yaml.Node) []metaRow {
	if fm == nil {
		return nil
	}
	var rows []metaRow
	for i := 0; i+1 < len(fm.Content); i += 2 {
		v := fm.Content[i+1]
		var s string
		if v.Kind == yaml.ScalarNode {
			s = v.Value
		} else {
			b, _ := yaml.Marshal(v)
			s = strings.Join(strings.Fields(string(b)), " ")
		}
		if s == "" || s == "[]" || s == "{}" {
			continue
		}
		rows = append(rows, metaRow{Key: fm.Content[i].Value, Value: s})
	}
	return rows
}

// toc gom tiêu đề cấp 2 và 3 ở mức đầu tài liệu (cấp 1 là tiêu đề trang, cấp
// 4 trở xuống quá chi tiết); dưới tocMin mục trả về nil để không hiện mục lục.
func toc(doc ast.Node, src []byte) []tocItem {
	var items []tocItem
	for n := doc.FirstChild(); n != nil; n = n.NextSibling() {
		h, ok := n.(*ast.Heading)
		if !ok || h.Level < 2 || h.Level > 3 {
			continue
		}
		id, _ := h.AttributeString("id")
		idb, _ := id.([]byte)
		items = append(items, tocItem{Level: h.Level, ID: string(idb), Text: nodeText(h, src)})
	}
	if len(items) < tocMin {
		return nil
	}
	return items
}

func firstHeading(doc ast.Node, src []byte) string {
	for n := doc.FirstChild(); n != nil; n = n.NextSibling() {
		if h, ok := n.(*ast.Heading); ok && h.Level == 1 {
			return nodeText(h, src)
		}
	}
	return ""
}

// nodeText ghép mọi đoạn chữ trong nút.
func nodeText(n ast.Node, src []byte) string {
	var sb strings.Builder
	_ = ast.Walk(n, func(c ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch t := c.(type) {
		case *ast.Text:
			sb.Write(t.Segment.Value(src))
		case *ast.String:
			sb.Write(t.Value)
		}
		return ast.WalkContinue, nil
	})
	return strings.TrimSpace(sb.String())
}

// linkRewriter đổi liên kết tương đối .md thành đường dẫn HTML tương đối
// từ vị trí trang đích; liên kết ra ngoài DocsDir trỏ thẳng đến file gốc.
type linkRewriter struct {
	opts           Options
	srcDir, outDir string
}

func (l *linkRewriter) Transform(doc *ast.Document, _ text.Reader, _ parser.Context) {
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if link, ok := n.(*ast.Link); ok && entering {
			link.Destination = []byte(l.rewrite(string(link.Destination)))
		}
		return ast.WalkContinue, nil
	})
}

func (l *linkRewriter) rewrite(dest string) string {
	if !IsRelativeLink(dest) {
		return dest
	}
	path, frag, _ := strings.Cut(dest, "#")
	if frag != "" {
		frag = "#" + frag
	}
	if !strings.HasSuffix(strings.ToLower(path), ".md") {
		return dest
	}
	target := filepath.Clean(filepath.Join(l.srcDir, filepath.FromSlash(path)))
	if out, err := l.opts.OutPath(target); err == nil {
		target = out
	}
	rel, err := filepath.Rel(l.outDir, target)
	if err != nil {
		return dest
	}
	return filepath.ToSlash(rel) + frag
}

// IsRelativeLink báo đích liên kết là đường dẫn tương đối trong repo
// (không scheme, không tuyệt đối, không chỉ neo).
func IsRelativeLink(dest string) bool {
	if dest == "" || strings.HasPrefix(dest, "#") || strings.HasPrefix(dest, "/") {
		return false
	}
	if i := strings.Index(dest, ":"); i >= 0 && !strings.Contains(dest[:i], "/") {
		return false // có scheme: http:, mailto:
	}
	return true
}
