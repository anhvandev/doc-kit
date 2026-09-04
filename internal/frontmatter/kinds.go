package frontmatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Tài liệu không phải Markdown giữ metadata theo ba định dạng, không thêm
// định dạng khác: HTML đặt YAML trong chú thích `<!-- dk:` ... `-->` đầu file;
// JSON đặt object dưới khóa `$dk` ở cấp cao nhất; Gherkin (.feature) đặt YAML
// trong khối chú thích liền nhau đầu file, mở bằng dòng `# dk:`, mỗi dòng YAML
// mang tiền tố `# `, kết thúc ở dòng đầu tiên không bắt đầu bằng `#` (hoặc hết file).

const (
	htmlOpen    = "<!-- dk:"
	htmlEnd     = "-->"
	jsonKey     = "$dk"
	featureOpen = "# dk:"
)

// SplitFile chọn cách tách theo đuôi file: .html, .json, .feature, còn lại là Markdown.
func SplitFile(name string, b []byte) (fm *yaml.Node, body []byte, ok bool) {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".html":
		return SplitHTMLComment(b)
	case ".json":
		return SplitJSONKey(b)
	case ".feature":
		return SplitFeatureComment(b)
	}
	return Split(b)
}

// JoinFile ghép lại theo đuôi file, đối xứng với SplitFile.
func JoinFile(name string, fm *yaml.Node, body []byte) ([]byte, error) {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".html":
		return JoinHTMLComment(fm, body)
	case ".json":
		return JoinJSONKey(fm, body)
	case ".feature":
		return JoinFeatureComment(fm, body)
	}
	return Join(fm, body)
}

// SplitHTMLComment tách YAML trong chú thích `<!-- dk:` ... `-->` ở đầu file HTML.
func SplitHTMLComment(b []byte) (fm *yaml.Node, body []byte, ok bool) {
	nl := newline(b)
	first, rest, found := bytes.Cut(b, []byte(nl))
	if !found || strings.TrimRight(string(first), "\r") != htmlOpen {
		return nil, b, false
	}
	lines := bytes.Split(rest, []byte(nl))
	for i, ln := range lines {
		if strings.TrimRight(string(ln), "\r") != htmlEnd {
			continue
		}
		raw := bytes.Join(lines[:i], []byte("\n"))
		body = bytes.Join(lines[i+1:], []byte(nl))
		fm, ok = parseMapping(raw)
		if !ok {
			return nil, b, false
		}
		return fm, body, true
	}
	return nil, b, false
}

// JoinHTMLComment ghép YAML vào chú thích đầu file HTML.
func JoinHTMLComment(fm *yaml.Node, body []byte) ([]byte, error) {
	nl := newline(body)
	y, err := encodeYAML(fm, nl)
	if err != nil {
		return nil, err
	}
	out := htmlOpen + nl + y + htmlEnd + nl
	return append([]byte(out), body...), nil
}

// SplitFeatureComment tách YAML trong khối chú thích `# dk:` đầu file Gherkin.
// Thân là mọi dòng từ dòng đầu tiên không phải chú thích `#`.
func SplitFeatureComment(b []byte) (fm *yaml.Node, body []byte, ok bool) {
	nl := newline(b)
	first, rest, found := bytes.Cut(b, []byte(nl))
	if !found || strings.TrimRight(string(first), "\r") != featureOpen {
		return nil, b, false
	}
	lines := bytes.Split(rest, []byte(nl))
	var raw [][]byte
	for i, ln := range lines {
		l := strings.TrimRight(string(ln), "\r")
		if !strings.HasPrefix(l, "#") {
			body = bytes.Join(lines[i:], []byte(nl))
			fm, ok = parseMapping(bytes.Join(raw, []byte("\n")))
			if !ok || len(fm.Content) == 0 {
				return nil, b, false // khối rỗng coi như thiếu metadata
			}
			return fm, body, true
		}
		raw = append(raw, []byte(strings.TrimPrefix(strings.TrimPrefix(l, "#"), " ")))
	}
	// Hết file khi vẫn trong khối: file chỉ có metadata, thân rỗng.
	fm, ok = parseMapping(bytes.Join(raw, []byte("\n")))
	if !ok || len(fm.Content) == 0 {
		return nil, b, false
	}
	return fm, nil, true
}

// JoinFeatureComment ghép YAML thành khối chú thích `# dk:` đầu file Gherkin.
// Thân bắt đầu bằng `#` sẽ được chèn một dòng trống để không bị nuốt vào khối.
func JoinFeatureComment(fm *yaml.Node, body []byte) ([]byte, error) {
	nl := newline(body)
	y, err := encodeYAML(fm, nl)
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	out.WriteString(featureOpen + nl)
	for _, l := range strings.Split(strings.TrimSuffix(y, nl), nl) {
		if l == "" {
			out.WriteString("#" + nl)
			continue
		}
		out.WriteString("# " + l + nl)
	}
	if len(body) == 0 || body[0] == '#' {
		out.WriteString(nl)
	}
	out.Write(body)
	return out.Bytes(), nil
}

// SplitJSONKey tách object `$dk` khỏi một file JSON. body là phần JSON còn
// lại sau khi bỏ khóa `$dk`, giữ nguyên định dạng các khóa khác.
func SplitJSONKey(b []byte) (fm *yaml.Node, body []byte, ok bool) {
	start, end, raw, found := findJSONKey(b, jsonKey)
	if !found {
		return nil, b, false
	}
	fm, ok = parseMapping(raw)
	if !ok {
		return nil, b, false
	}
	// Bỏ cả dấu phẩy đi kèm: sau giá trị nếu còn khóa khác, không thì trước khóa.
	after := end
	for after < len(b) && (b[after] == ' ' || b[after] == '\t' || b[after] == '\r' || b[after] == '\n') {
		after++
	}
	if after < len(b) && b[after] == ',' {
		after++
		for after < len(b) && (b[after] == ' ' || b[after] == '\t' || b[after] == '\r' || b[after] == '\n') {
			after++
		}
		return fm, append(append([]byte{}, b[:start]...), b[after:]...), true
	}
	before := start
	for before > 0 && (b[before-1] == ' ' || b[before-1] == '\t' || b[before-1] == '\r' || b[before-1] == '\n') {
		before--
	}
	if before > 0 && b[before-1] == ',' {
		before--
	}
	return fm, append(append([]byte{}, b[:before]...), b[end:]...), true
}

// JoinJSONKey đặt `$dk` làm khóa đầu tiên của object JSON trong body.
func JoinJSONKey(fm *yaml.Node, body []byte) ([]byte, error) {
	if fm == nil || fm.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("ghi $dk: frontmatter phải là mapping")
	}
	// encoding/json sắp map theo chữ cái nên ghi từng cặp theo thứ tự nút.
	meta := orderedJSON(fm)
	i := bytes.IndexByte(body, '{')
	if i < 0 {
		return nil, fmt.Errorf("ghi $dk: body không phải object JSON")
	}
	nl := newline(body)
	rest := bytes.TrimLeft(body[i+1:], " \t\r\n")
	sep := "," + nl + "  "
	if len(rest) > 0 && rest[0] == '}' {
		sep = nl
	}
	var out bytes.Buffer
	out.Write(body[:i+1])
	out.WriteString(nl + "  \"" + jsonKey + "\": ")
	out.Write(meta)
	out.WriteString(sep)
	out.Write(rest)
	return out.Bytes(), nil
}

// orderedJSON ghi mapping theo thứ tự khóa, thụt 2 lề bên trong `$dk`.
func orderedJSON(fm *yaml.Node) []byte {
	var out bytes.Buffer
	out.WriteString("{")
	for i := 0; i+1 < len(fm.Content); i += 2 {
		var v any
		_ = fm.Content[i+1].Decode(&v)
		val, _ := json.Marshal(v)
		if i > 0 {
			out.WriteString(",")
		}
		key, _ := json.Marshal(fm.Content[i].Value)
		out.WriteString("\n    " + string(key) + ": " + string(val))
	}
	if len(fm.Content) > 0 {
		out.WriteString("\n  ")
	}
	out.WriteString("}")
	return out.Bytes()
}

// findJSONKey tìm khóa key ở cấp cao nhất của object JSON; trả về vị trí bắt
// đầu của khóa, vị trí kết thúc của giá trị và giá trị thô.
func findJSONKey(b []byte, key string) (start, end int, raw []byte, ok bool) {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	tok, err := dec.Token()
	if err != nil || tok != json.Delim('{') {
		return 0, 0, nil, false
	}
	prev := int(dec.InputOffset())
	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			return 0, 0, nil, false
		}
		k, _ := tok.(string)
		keyEnd := int(dec.InputOffset())
		if k == key {
			start = prev + bytes.IndexByte(b[prev:keyEnd], '"')
			var rv json.RawMessage
			if err := dec.Decode(&rv); err != nil {
				return 0, 0, nil, false
			}
			return start, int(dec.InputOffset()), rv, true
		}
		var skip json.RawMessage
		if err := dec.Decode(&skip); err != nil {
			return 0, 0, nil, false
		}
		prev = int(dec.InputOffset())
	}
	return 0, 0, nil, false
}

// parseMapping đọc YAML (hoặc JSON, là tập con) thành mapping node.
func parseMapping(raw []byte) (*yaml.Node, bool) {
	var doc yaml.Node
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return nil, false
	}
	if len(doc.Content) == 0 {
		return emptyMapping(), true
	}
	if doc.Content[0].Kind != yaml.MappingNode {
		return nil, false
	}
	return doc.Content[0], true
}

func encodeYAML(fm *yaml.Node, nl string) (string, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(fm); err != nil {
		return "", fmt.Errorf("ghi frontmatter: %w", err)
	}
	if err := enc.Close(); err != nil {
		return "", err
	}
	y := buf.String()
	if len(fm.Content) == 0 {
		y = ""
	}
	if nl == "\r\n" {
		y = strings.ReplaceAll(y, "\n", nl)
	}
	return y, nil
}
