// Package doctype đọc bảng loại tài liệu từ assets/types.toml.
package doctype

import (
	"fmt"
	"io/fs"
	"regexp"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// Type mô tả một loại tài liệu.
type Type struct {
	Name        string                       `toml:"-"`
	Dir         string                       `toml:"dir"`
	Subdir      string                       `toml:"subdir"`
	NamePattern string                       `toml:"name"`
	IDScheme    string                       `toml:"id"`
	Description string                       `toml:"description"`
	Required    []string                     `toml:"required"`
	Statuses    []string                     `toml:"statuses"`
	Final       []string                     `toml:"final"`
	From        map[string]map[string]string `toml:"from"`
	// BesideSource: khi --from là file ngoài thư mục loại (ví dụ cr/CR-x.md),
	// đặt file mới vào thư mục cùng tên nguồn (cr/CR-x/) thay vì subdir mới.
	BesideSource bool `toml:"beside_source"`
	// WarnLines: ngưỡng cảnh báo số dòng riêng của loại; 0 dùng [check] warn_lines của dk.toml.
	WarnLines int `toml:"warn_lines"`
	// Kind: định dạng file: "" hoặc "md" (Markdown, frontmatter YAML), "html"
	// (metadata trong chú thích <!-- dk: --> đầu file), "json" (khóa $dk),
	// "feature" (Gherkin, metadata trong khối chú thích # dk: đầu file).
	Kind string `toml:"kind"`
}

// Ext trả về đuôi file của loại, kể cả dấu chấm.
func (t Type) Ext() string {
	if t.Kind == "" {
		return ".md"
	}
	return "." + t.Kind
}

// Registry là tập loại tài liệu, tra theo tên.
type Registry map[string]Type

var idSchemeRe = regexp.MustCompile(`^(none|seq:[^{}]*\{n(?::0\d+)?\}[^{}]*|date:[^{}]*\{yymmdd\}[^{}]*(?:\{slug\}[^{}]*)?)$`)

// Load đọc types.toml và kiểm tra: dir không rỗng, id hợp lệ, template tồn tại.
func Load(fsys fs.FS) (Registry, error) {
	b, err := fs.ReadFile(fsys, "types.toml")
	if err != nil {
		return nil, fmt.Errorf("đọc types.toml: %w", err)
	}
	raw := map[string]Type{}
	if err := toml.Unmarshal(b, &raw); err != nil {
		return nil, fmt.Errorf("phân tích types.toml: %w", err)
	}
	reg := Registry{}
	for name, t := range raw {
		t.Name = name
		if err := validate(fsys, t); err != nil {
			return nil, fmt.Errorf("loại %q: %w", name, err)
		}
		reg[name] = t
	}
	for name, t := range reg {
		for src := range t.From {
			if _, ok := reg[src]; !ok {
				return nil, fmt.Errorf("loại %q: from.%s trỏ loại không tồn tại", name, src)
			}
		}
	}
	return reg, nil
}

func validate(fsys fs.FS, t Type) error {
	if strings.TrimSpace(t.Dir) == "" {
		return fmt.Errorf("thiếu dir")
	}
	if strings.TrimSpace(t.NamePattern) == "" {
		return fmt.Errorf("thiếu name")
	}
	if !idSchemeRe.MatchString(t.IDScheme) {
		return fmt.Errorf("id %q không hợp lệ", t.IDScheme)
	}
	if t.IDScheme != "none" && !strings.Contains(t.NamePattern, "{id}") {
		return fmt.Errorf("name phải chứa {id} khi id khác none")
	}
	if t.Kind != "" && t.Kind != "md" && t.Kind != "html" && t.Kind != "json" && t.Kind != "feature" {
		return fmt.Errorf("kind %q không hỗ trợ; chỉ md, html, json, feature", t.Kind)
	}
	if _, err := fs.Stat(fsys, "templates/"+t.Name+t.Ext()); err != nil {
		return fmt.Errorf("thiếu template templates/%s%s", t.Name, t.Ext())
	}
	return nil
}

// Names trả về tên loại theo thứ tự chữ cái.
func (r Registry) Names() []string {
	names := make([]string, 0, len(r))
	for n := range r {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// Get tra loại; lỗi nếu không có.
func (r Registry) Get(name string) (Type, error) {
	t, ok := r[name]
	if !ok {
		return Type{}, fmt.Errorf("không có loại tài liệu %q; xem `dk template list`", name)
	}
	return t, nil
}
