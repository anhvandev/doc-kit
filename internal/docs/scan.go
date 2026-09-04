package docs

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/anhvandev/doc-kit/internal/frontmatter"
)

// Meta là frontmatter và thân của một file Markdown đã quét.
type Meta struct {
	Path      string     `json:"-"`    // tuyệt đối
	Rel       string     `json:"path"` // tương đối gốc dự án, dạng slash
	ID        string     `json:"id"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Status    string     `json:"status"`
	Owner     string     `json:"owner"`
	Created   string     `json:"created"`
	Updated   string     `json:"updated"`
	Source    string     `json:"source"`
	Generated bool       `json:"generated"`
	HasFM     bool       `json:"-"` // có frontmatter hợp lệ
	Lines     int        `json:"lines"`
	FM        *yaml.Node `json:"-"`
	Body      []byte     `json:"-"`
	Raw       []byte     `json:"-"`
}

// Scan đọc mọi file .md, .html, .json, .feature trong các thư mục dirs (tương
// đối root), bỏ qua html/ và file ẩn. Kết quả sắp theo Rel. File .html và
// .json là tài liệu họ Design (mockup, tokens), .feature là test case Gherkin;
// bản sinh như tokens.css không quét.
func Scan(root string, dirs ...string) ([]Meta, error) {
	var metas []Meta
	for _, d := range dirs {
		base := filepath.Join(root, filepath.FromSlash(d))
		err := filepath.WalkDir(base, func(p string, e fs.DirEntry, err error) error {
			if err != nil {
				if p == base && os.IsNotExist(err) {
					return filepath.SkipDir
				}
				return err
			}
			name := e.Name()
			if e.IsDir() {
				if p != base && (strings.HasPrefix(name, ".") || (name == "html" && filepath.Dir(p) == base)) {
					return filepath.SkipDir
				}
				return nil
			}
			if strings.HasPrefix(name, ".") || !scannedExt[strings.ToLower(filepath.Ext(name))] {
				return nil
			}
			m, err := Read(root, p)
			if err != nil {
				return err
			}
			metas = append(metas, m)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	sort.Slice(metas, func(i, j int) bool { return metas[i].Rel < metas[j].Rel })
	return metas, nil
}

var scannedExt = map[string]bool{".md": true, ".html": true, ".json": true, ".feature": true}

// IsMarkdown báo Meta là file Markdown (render được sang HTML).
func (m Meta) IsMarkdown() bool { return strings.EqualFold(filepath.Ext(m.Rel), ".md") }

// Read đọc một file thành Meta; metadata tách theo đuôi file.
func Read(root, path string) (Meta, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Meta{}, err
	}
	rel, err := filepath.Rel(root, path)
	if err != nil {
		rel = path
	}
	m := Meta{Path: path, Rel: filepath.ToSlash(rel), Raw: b, Lines: CountLines(b)}
	m.FM, m.Body, m.HasFM = frontmatter.SplitFile(path, b)
	if m.HasFM {
		get := func(k string) string { return frontmatter.GetString(m.FM, k) }
		m.ID, m.Type, m.Title, m.Status = get("id"), get("type"), get("title"), get("status")
		m.Owner, m.Created, m.Updated, m.Source = get("owner"), get("created"), get("updated"), get("source")
		m.Generated = get("generated") == "true"
	}
	return m, nil
}

// CountLines đếm dòng, dòng cuối không có \n vẫn tính.
func CountLines(b []byte) int {
	n := bytes.Count(b, []byte("\n"))
	if len(b) > 0 && b[len(b)-1] != '\n' {
		n++
	}
	return n
}

// ByID lập bảng id sang Meta; id rỗng bỏ qua.
func ByID(metas []Meta) map[string]Meta {
	m := map[string]Meta{}
	for _, d := range metas {
		if d.ID != "" {
			m[d.ID] = d
		}
	}
	return m
}

// Resolve tìm tài liệu theo tham chiếu: id, hoặc "<thư mục>/<file>" với
// tài liệu không có id (intake).
func Resolve(metas []Meta, ref string) (Meta, bool) {
	if ref == "" {
		return Meta{}, false
	}
	for _, d := range metas {
		if d.ID == ref || (strings.Contains(ref, "/") && strings.HasSuffix(d.Rel, "/"+ref)) {
			return d, true
		}
	}
	return Meta{}, false
}
