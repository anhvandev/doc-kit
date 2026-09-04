package check

import (
	"bytes"
	"strings"

	"github.com/anhvandev/doc-kit/internal/frontmatter"
	"github.com/anhvandev/doc-kit/internal/gitx"
)

// adrFrozen là các trạng thái mà thân ADR không được đổi nữa.
var adrFrozen = []string{"accepted", "superseded", "deprecated"}

// adrImmutable: ADR đã chốt ở HEAD (accepted, superseded, deprecated) mà thân
// khác với HEAD là lỗi; frontmatter được đổi (status, superseded_by). So sau
// khi chuẩn hóa xuống dòng và cắt khoảng trắng cuối dòng. Không có git thì bỏ qua.
func adrImmutable(c *Context) []Finding {
	var out []Finding
	for _, m := range c.typed() {
		if m.Type != "adr" {
			continue
		}
		head, ok, err := gitx.HeadFile(c.Root, m.Path)
		if err != nil {
			out = append(out, finding(m, "adr-immutable", Warning, "không đọc được bản HEAD: "+err.Error()))
			continue
		}
		if !ok {
			continue
		}
		hfm, hbody, has := frontmatter.Split(head)
		if !has || !contains(adrFrozen, frontmatter.GetString(hfm, "status")) {
			continue
		}
		if !bytes.Equal(normalize(hbody), normalize(m.Body)) {
			out = append(out, finding(m, "adr-immutable", Error, "thân ADR đã chốt bị đổi so với HEAD; quyết định mới ghi ADR mới với supersedes"))
		}
	}
	return out
}

// normalize đưa CRLF về LF, cắt khoảng trắng cuối mỗi dòng và dòng trống cuối file.
func normalize(b []byte) []byte {
	s := strings.ReplaceAll(string(b), "\r\n", "\n")
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimRight(l, " \t")
	}
	return []byte(strings.TrimRight(strings.Join(lines, "\n"), "\n"))
}
