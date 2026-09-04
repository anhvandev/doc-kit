package skill

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"

	"github.com/anhvandev/doc-kit/internal/frontmatter"
)

// Hash băm sha256 mọi file theo tên đã sắp; SKILL.md được bỏ khối metadata
// trước khi băm để dấu vết cài không đổi hash.
func Hash(files map[string][]byte) string {
	names := make([]string, 0, len(files))
	for n := range files {
		names = append(names, n)
	}
	sort.Strings(names)
	h := sha256.New()
	for _, n := range names {
		b := files[n]
		if n == skillFile {
			b = stripMeta(b)
		}
		h.Write([]byte(n))
		h.Write([]byte{0})
		h.Write(b)
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}

// stripMeta bỏ khóa metadata khỏi frontmatter và ghi lại theo dạng chuẩn để
// bản nguồn và bản đã cài băm giống nhau.
func stripMeta(b []byte) []byte {
	fm, body, ok := frontmatter.Split(b)
	if !ok {
		return b
	}
	for i := 0; i+1 < len(fm.Content); i += 2 {
		if fm.Content[i].Value == metaKey {
			fm.Content = append(fm.Content[:i], fm.Content[i+2:]...)
			break
		}
	}
	out, err := frontmatter.Join(fm, body)
	if err != nil {
		return b
	}
	return out
}
