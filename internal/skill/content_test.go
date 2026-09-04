package skill

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/anhvandev/doc-kit/internal/frontmatter"
)

// Nội dung skill nhúng phải trung lập target và ngắn: frontmatter hợp lệ với
// name trùng thư mục, mọi file Markdown dưới 300 dòng, không nhắc target hay
// tên tool của một agent cụ thể, mô tả kích hoạt không trùng nhau.
const maxSkillLines = 300

var forbidden = regexp.MustCompile(`ak-|Claude Code|\.claude/|\.codex/|\bEdit\b|\bWrite\b|\bMultiEdit\b`)

func TestSkillContent(t *testing.T) {
	list, err := List()
	if err != nil || len(list) == 0 {
		t.Fatalf("List: %v (%d)", err, len(list))
	}
	descs := map[string]string{}
	for _, m := range list {
		files, err := Files(m.Name)
		if err != nil {
			t.Fatal(err)
		}
		fm, _, ok := frontmatter.Split(files["SKILL.md"])
		if !ok || frontmatter.GetString(fm, "name") != m.Name || m.Description == "" {
			t.Errorf("%s: frontmatter SKILL.md cần name trùng thư mục và description", m.Name)
		}
		if _, ok := files["references/rules.md"]; !ok {
			t.Errorf("%s: thiếu references/rules.md", m.Name)
		}
		descs[m.Name] = m.Description
		for name, b := range files {
			if !strings.HasSuffix(name, ".md") {
				continue
			}
			if n := len(bytes.Split(bytes.TrimRight(b, "\n"), []byte("\n"))); n >= maxSkillLines {
				t.Errorf("%s/%s: %d dòng, phải dưới %d", m.Name, name, n, maxSkillLines)
			}
			if loc := forbidden.FindIndex(b); loc != nil {
				t.Errorf("%s/%s: chứa từ cấm %q", m.Name, name, b[loc[0]:loc[1]])
			}
		}
	}
	for a, da := range descs {
		for b, db := range descs {
			if a < b {
				if tri := sharedTrigram(da, db); tri != "" {
					t.Errorf("mô tả %s và %s trùng 3 từ liên tiếp: %q", a, b, tri)
				}
			}
		}
	}
}

// sharedTrigram tìm 3 từ liên tiếp có trong cả hai mô tả. Câu phủ định
// "Không dùng cho ..." ở cuối bị bỏ qua vì nó cố ý nhắc phạm vi của skill kia.
func sharedTrigram(a, b string) string {
	seen := map[string]bool{}
	for _, tri := range trigrams(a) {
		seen[tri] = true
	}
	for _, tri := range trigrams(b) {
		if seen[tri] {
			return tri
		}
	}
	return ""
}

var wordRe = regexp.MustCompile(`[\p{L}\p{N}]+`)

func trigrams(s string) []string {
	if i := strings.Index(s, "Không dùng"); i >= 0 {
		s = s[:i]
	}
	words := wordRe.FindAllString(strings.ToLower(s), -1)
	var out []string
	for i := 0; i+3 <= len(words); i++ {
		out = append(out, fmt.Sprintf("%s %s %s", words[i], words[i+1], words[i+2]))
	}
	return out
}
