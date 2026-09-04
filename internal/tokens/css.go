package tokens

import (
	"fmt"
	"sort"
	"strings"
)

// Header là dòng đầu của tokens.css; changelog và Scan nhận ra bản sinh nhờ nó.
const Header = "/* generated: dk tokens css */"

// CSS sinh `:root { --a-b: v; }` cho mọi token theo thứ tự file, rồi một khối
// `[data-theme="<tên>"]` cho mỗi theme khai ở $extensions.dk.theme. Trả về
// số biến đã ghi.
func (s *Set) CSS() ([]byte, int, error) {
	var b strings.Builder
	b.WriteString(Header + "\n")
	b.WriteString("/* Sinh từ tokens.json; không sửa tay, chạy lại `dk tokens css`. */\n")
	b.WriteString(":root {\n")
	themes := map[string][]string{}
	var themeNames []string
	for _, t := range s.List {
		v, err := s.Resolve(t)
		if err != nil {
			return nil, 0, err
		}
		fmt.Fprintf(&b, "  %s: %s;\n", t.Name(), v)
		for name, alt := range t.Theme {
			tv, err := s.resolve(alt, t.Type, []string{strings.Join(t.Path, ".") + "@" + name})
			if err != nil {
				return nil, 0, err
			}
			if _, seen := themes[name]; !seen {
				themeNames = append(themeNames, name)
			}
			themes[name] = append(themes[name], fmt.Sprintf("  %s: %s;", t.Name(), tv))
		}
	}
	b.WriteString("}\n")
	sort.Strings(themeNames)
	for _, name := range themeNames {
		fmt.Fprintf(&b, "\n[data-theme=%q] {\n%s\n}\n", name, strings.Join(themes[name], "\n"))
	}
	return []byte(b.String()), len(s.List), nil
}
