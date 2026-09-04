// Package tokens đọc file Design Tokens theo khung W3C (nhóm lồng, $value,
// $type, alias {a.b.c}) và sinh CSS variables. Chỉ hai loại được chuẩn hóa:
// color và dimension; loại khác chép nguyên giá trị chuỗi hoặc số.
package tokens

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Token là một giá trị đã đọc, giữ thứ tự xuất hiện trong file.
type Token struct {
	Path  []string // ví dụ color, action, primary
	Type  string   // $type, kế thừa từ nhóm cha
	Value any      // chuỗi, số, hoặc map cho giá trị hợp thành
	Theme map[string]any
}

// Name trả về tên biến CSS: đường dẫn nối gạch ngang, ví dụ --color-action-primary.
func (t Token) Name() string { return "--" + strings.Join(t.Path, "-") }

// Set là danh sách token theo thứ tự file và bảng tra theo đường dẫn.
type Set struct {
	List  []Token
	byRef map[string]*Token
}

// Parse đọc JSON tokens; bỏ qua khóa `$dk` (metadata của dk) và các khóa `$...` khác ở nhóm.
func Parse(b []byte) (*Set, error) {
	root, err := parseOrdered(b)
	if err != nil {
		return nil, fmt.Errorf("tokens.json: %w", err)
	}
	if root.kind != kindObject {
		return nil, fmt.Errorf("tokens.json: cấp cao nhất phải là object")
	}
	s := &Set{byRef: map[string]*Token{}}
	if err := s.walk(root, nil, ""); err != nil {
		return nil, err
	}
	for i := range s.List {
		s.byRef[strings.Join(s.List[i].Path, ".")] = &s.List[i]
	}
	return s, nil
}

func (s *Set) walk(n *node, path []string, typ string) error {
	if t := n.get("$type"); t != nil && t.kind == kindString {
		typ = t.str
	}
	if v := n.get("$value"); v != nil {
		tok := Token{Path: append([]string{}, path...), Type: typ, Value: v.plain()}
		if ext := n.get("$extensions"); ext != nil {
			if dk := ext.get("dk"); dk != nil {
				if th := dk.get("theme"); th != nil && th.kind == kindObject {
					tok.Theme = map[string]any{}
					for _, k := range th.keys {
						tok.Theme[k] = th.obj[k].plain()
					}
				}
			}
		}
		s.List = append(s.List, tok)
		return nil
	}
	for _, k := range n.keys {
		if strings.HasPrefix(k, "$") {
			continue
		}
		child := n.obj[k]
		if child.kind != kindObject {
			return fmt.Errorf("tokens.json: %s không phải nhóm hay token (thiếu $value)", strings.Join(append(path, k), "."))
		}
		if err := s.walk(child, append(path, k), typ); err != nil {
			return err
		}
	}
	return nil
}

var aliasRe = regexp.MustCompile(`\{([A-Za-z0-9_.-]+)\}`)

// Resolve trả về giá trị CSS của token đã giải alias (kể cả alias lồng trong
// chuỗi như "0 1px 2px {color.shadow}") và chuẩn hóa theo loại; alias vòng
// hoặc trỏ token không có là lỗi.
func (s *Set) Resolve(t Token) (string, error) {
	return s.resolve(t.Value, t.Type, []string{strings.Join(t.Path, ".")})
}

func (s *Set) resolve(v any, typ string, chain []string) (string, error) {
	switch x := v.(type) {
	case string:
		var err error
		out := aliasRe.ReplaceAllStringFunc(x, func(m string) string {
			if err != nil {
				return m
			}
			ref := m[1 : len(m)-1]
			for _, c := range chain {
				if c == ref {
					err = fmt.Errorf("alias vòng: %s -> %s", strings.Join(chain, " -> "), ref)
					return m
				}
			}
			target, ok := s.byRef[ref]
			if !ok {
				err = fmt.Errorf("%s: alias {%s} không trỏ đến token nào", chain[0], ref)
				return m
			}
			r, e := s.resolve(target.Value, target.Type, append(append([]string{}, chain...), ref))
			if e != nil {
				err = e
				return m
			}
			return r
		})
		return out, err
	case json.Number:
		if typ == "dimension" {
			return x.String() + "px", nil
		}
		return x.String(), nil
	case bool:
		return fmt.Sprint(x), nil
	case map[string]any:
		return composite(x, typ, chain[0])
	}
	return "", fmt.Errorf("%s: giá trị %v không hỗ trợ", chain[0], v)
}

// composite chuẩn hóa dimension {value, unit} và color {hex}; loại khác lỗi.
func composite(m map[string]any, typ, name string) (string, error) {
	switch typ {
	case "dimension":
		val, unit := m["value"], m["unit"]
		if val == nil {
			return "", fmt.Errorf("%s: dimension thiếu value", name)
		}
		if unit == nil {
			unit = "px"
		}
		return fmt.Sprintf("%v%v", val, unit), nil
	case "color":
		if hex, ok := m["hex"].(string); ok {
			return hex, nil
		}
		return "", fmt.Errorf("%s: color dạng object cần trường hex", name)
	}
	return "", fmt.Errorf("%s: loại %q dạng object không hỗ trợ; ghi thành chuỗi", name, typ)
}
