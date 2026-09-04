package render

import (
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
)

var (
	stepCodeRe   = regexp.MustCompile(`^B\d+[a-z]?$`)
	stepInTextRe = regexp.MustCompile(`\b(B\d+[a-z]?)\b`)
)

// stepCodesIn lấy mã bước trong một khối mermaid, giữ thứ tự xuất hiện, không trùng.
func stepCodesIn(diagram []byte) []string {
	var out []string
	seen := map[string]bool{}
	for _, m := range stepInTextRe.FindAllSubmatch(diagram, -1) {
		c := string(m[1])
		if !seen[c] {
			seen[c] = true
			out = append(out, c)
		}
	}
	return out
}

// findBehaviorTable trả về bảng hành vi: bảng đầu tiên dưới tiêu đề chứa
// "hành vi" hoặc "use case" (không phân biệt hoa thường); không có thì bảng đầu tiên có ô
// cột đầu khớp ^B\d+[a-z]?$ (hậu tố chữ thường là bước chèn giữa, ví dụ B2a).
func findBehaviorTable(doc ast.Node, src []byte) ast.Node {
	var byHeading, byCode ast.Node
	underHeading := false
	for n := doc.FirstChild(); n != nil; n = n.NextSibling() {
		switch t := n.(type) {
		case *ast.Heading:
			h := strings.ToLower(nodeText(t, src))
			underHeading = strings.Contains(h, "hành vi") || strings.Contains(h, "use case")
		case *east.Table:
			if underHeading && byHeading == nil {
				byHeading = t
			}
			if byCode == nil && hasStepColumn(t, src) {
				byCode = t
			}
		}
	}
	if byHeading != nil {
		return byHeading
	}
	return byCode
}

func hasStepColumn(t *east.Table, src []byte) bool {
	for row := t.FirstChild(); row != nil; row = row.NextSibling() {
		if _, isRow := row.(*east.TableRow); isRow && stepCodeRe.MatchString(nodeText(row.FirstChild(), src)) {
			return true
		}
	}
	return false
}

// StepCodes trả về tập mã bước trong mọi sơ đồ mermaid và trong cột đầu
// bảng hành vi của thân tài liệu. Dùng chung với `dk check`.
func StepCodes(body []byte) (diagram, table []string) {
	md := goldmark.New(goldmark.WithExtensions(extension.Table))
	doc := md.Parser().Parse(text.NewReader(body))
	seen := map[string]bool{}
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if fc, ok := n.(*ast.FencedCodeBlock); ok && entering && string(fc.Language(body)) == "mermaid" {
			for _, c := range stepCodesIn(blockText(fc, body)) {
				if !seen[c] {
					seen[c] = true
					diagram = append(diagram, c)
				}
			}
		}
		return ast.WalkContinue, nil
	})
	if t := findBehaviorTable(doc, body); t != nil {
		for row := t.FirstChild(); row != nil; row = row.NextSibling() {
			if _, isRow := row.(*east.TableRow); !isRow {
				continue
			}
			if c := nodeText(row.FirstChild(), body); stepCodeRe.MatchString(c) {
				table = append(table, c)
			}
		}
	}
	return diagram, table
}
