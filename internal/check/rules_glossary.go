package check

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/anhvandev/doc-kit/internal/docs"
)

var (
	boldRe    = regexp.MustCompile(`\*\*([^*\n]+)\*\*`)
	sectionRe = regexp.MustCompile(`^## (\d+)\.`)
)

// glossarySections là các mục của Feature Spec được quét thuật ngữ in đậm:
// 1 Mục đích, 4 Hành vi, 7 Quy tắc nghiệp vụ. Mục 8 có **Given/When/Then** không phải thuật ngữ.
var glossarySections = map[string]bool{"1": true, "4": true, "7": true}

// glossaryTerm: thuật ngữ in đậm lần đầu trong mục 1, 4, 7 của Feature Spec
// phải có ở cột đầu bảng của một file Glossary (mọi file loại glossary được
// gộp, để tách Glossary theo miền vẫn đúng); thiếu là warning. Chưa có Glossary
// thì bỏ qua.
func glossaryTerm(c *Context) []Finding {
	defined := map[string]bool{}
	var files []string
	for _, m := range c.Metas {
		if m.HasFM && m.Type == "glossary" {
			files = append(files, m.Rel)
			for t := range glossaryTerms(m.Body) {
				defined[t] = true
			}
		}
	}
	if len(files) == 0 {
		return nil
	}
	where := strings.Join(files, ", ")
	var out []Finding
	for _, m := range c.typed() {
		if m.Type != "feature-spec" {
			continue
		}
		for _, b := range firstBold(m) {
			if !defined[strings.ToLower(b.term)] {
				out = append(out, Finding{File: m.Rel, Line: b.line, Rule: "glossary-term", Level: Warning,
					Msg: fmt.Sprintf("thuật ngữ **%s** chưa có trong %s", b.term, where)})
			}
		}
	}
	return out
}

// glossaryTerms lấy cột đầu của mọi dòng bảng trong Glossary, chữ thường.
func glossaryTerms(body []byte) map[string]bool {
	terms := map[string]bool{}
	for _, l := range bytes.Split(body, []byte("\n")) {
		s := strings.TrimSpace(string(l))
		if !strings.HasPrefix(s, "|") {
			continue
		}
		cells := strings.Split(strings.Trim(s, "|"), "|")
		t := strings.TrimSpace(cells[0])
		if t == "" || t == "Thuật ngữ" || strings.Trim(t, "-: ") == "" {
			continue
		}
		terms[strings.ToLower(t)] = true
	}
	return terms
}

type boldTerm struct {
	term string
	line int
}

// firstBold trả về mỗi thuật ngữ in đậm (gộp hoa thường) và dòng xuất hiện
// đầu tiên trong các mục glossarySections của thân tài liệu, theo thứ tự xuất
// hiện; bỏ qua khối mã.
func firstBold(m docs.Meta) []boldTerm {
	var out []boldTerm
	seen := map[string]bool{}
	fmLines := docs.CountLines(m.Raw) - docs.CountLines(m.Body)
	in, code := false, false
	for i, l := range strings.Split(string(m.Body), "\n") {
		if strings.HasPrefix(strings.TrimSpace(l), "```") {
			code = !code
			continue
		}
		if code {
			continue
		}
		if sm := sectionRe.FindStringSubmatch(l); sm != nil {
			in = glossarySections[sm[1]]
			continue
		}
		if !in {
			continue
		}
		for _, sm := range boldRe.FindAllStringSubmatch(l, -1) {
			term := strings.TrimSpace(strings.TrimSuffix(sm[1], ":"))
			key := strings.ToLower(term)
			if term == "" || seen[key] {
				continue
			}
			seen[key] = true
			out = append(out, boldTerm{term, fmLines + i + 1})
		}
	}
	return out
}
