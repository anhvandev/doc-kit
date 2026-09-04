// Package frontmatter tách, đọc, sửa và ghi lại YAML frontmatter của tài liệu
// Markdown. Dùng yaml.Node để giữ nguyên thứ tự khóa và không đổi thân.
package frontmatter

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const fence = "---"

// Split tách frontmatter khỏi thân. Trả về ok=false khi file không có
// frontmatter hợp lệ; khi đó body là toàn bộ nội dung.
func Split(b []byte) (fm *yaml.Node, body []byte, ok bool) {
	nl := newline(b)
	first, rest, found := bytes.Cut(b, []byte(nl))
	if !found || strings.TrimRight(string(first), "\r") != fence {
		return nil, b, false
	}
	lines := bytes.Split(rest, []byte(nl))
	for i, ln := range lines {
		if strings.TrimRight(string(ln), "\r") != fence {
			continue
		}
		raw := bytes.Join(lines[:i], []byte("\n"))
		body = bytes.Join(lines[i+1:], []byte(nl))
		var doc yaml.Node
		if err := yaml.Unmarshal(raw, &doc); err != nil {
			return nil, b, false
		}
		if len(doc.Content) == 0 {
			return emptyMapping(), body, true
		}
		if doc.Content[0].Kind != yaml.MappingNode {
			return nil, b, false
		}
		return doc.Content[0], body, true
	}
	return nil, b, false
}

// Join ghép frontmatter và thân. Dùng CRLF khi thân dùng CRLF.
func Join(fm *yaml.Node, body []byte) ([]byte, error) {
	nl := newline(body)
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(fm); err != nil {
		return nil, fmt.Errorf("ghi frontmatter: %w", err)
	}
	if err := enc.Close(); err != nil {
		return nil, err
	}
	y := buf.String()
	if len(fm.Content) == 0 {
		y = ""
	}
	if nl == "\r\n" {
		y = strings.ReplaceAll(y, "\n", nl)
	}
	out := fence + nl + y + fence + nl
	return append([]byte(out), body...), nil
}

// Get trả về nút giá trị của khóa.
func Get(fm *yaml.Node, key string) (*yaml.Node, bool) {
	if fm == nil {
		return nil, false
	}
	for i := 0; i+1 < len(fm.Content); i += 2 {
		if fm.Content[i].Value == key {
			return fm.Content[i+1], true
		}
	}
	return nil, false
}

// GetString trả về giá trị chuỗi của khóa, rỗng nếu thiếu hoặc không phải scalar.
func GetString(fm *yaml.Node, key string) string {
	v, ok := Get(fm, key)
	if !ok || v.Kind != yaml.ScalarNode {
		return ""
	}
	return v.Value
}

// Set thay giá trị của khóa; khóa chưa có thì thêm vào cuối.
func Set(fm *yaml.Node, key string, value *yaml.Node) {
	for i := 0; i+1 < len(fm.Content); i += 2 {
		if fm.Content[i].Value == key {
			fm.Content[i+1] = value
			return
		}
	}
	fm.Content = append(fm.Content, scalar(key), value)
}

// SetString đặt khóa thành chuỗi.
func SetString(fm *yaml.Node, key, val string) {
	Set(fm, key, scalar(val))
}

// SetBool đặt khóa thành bool (không quote).
func SetBool(fm *yaml.Node, key string, val bool) {
	Set(fm, key, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: strconv.FormatBool(val)})
}

// Map chuyển frontmatter thành map để đọc tiện.
func Map(fm *yaml.Node) map[string]any {
	m := map[string]any{}
	if fm != nil {
		_ = fm.Decode(&m)
	}
	return m
}

func scalar(v string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: v}
}

func emptyMapping() *yaml.Node {
	return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
}

func newline(b []byte) string {
	i := bytes.IndexByte(b, '\n')
	if i > 0 && b[i-1] == '\r' {
		return "\r\n"
	}
	return "\n"
}
