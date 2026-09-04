package render

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// nodeRenderer ghi đè fenced code (mermaid) và ô bảng (id mã bước).
type nodeRenderer struct {
	hasMermaid    bool
	behaviorTable ast.Node // bảng hành vi, nil nếu không có
}

func (r *nodeRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindFencedCodeBlock, r.renderFenced)
	reg.Register(east.KindTableCell, r.renderCell)
}

func (r *nodeRenderer) renderFenced(w util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.FencedCodeBlock)
	lang := string(n.Language(src))
	body := blockText(n, src)
	if lang != "mermaid" {
		_, _ = w.WriteString("<pre><code")
		if lang != "" {
			fmt.Fprintf(w, ` class="language-%s"`, util.EscapeHTML([]byte(lang)))
		}
		_, _ = w.WriteString(">")
		_, _ = w.Write(util.EscapeHTML(body))
		_, _ = w.WriteString("</code></pre>\n")
		return ast.WalkSkipChildren, nil
	}
	r.hasMermaid = true
	_, _ = w.WriteString(`<pre class="mermaid">`)
	_, _ = w.Write(util.EscapeHTML(body))
	_, _ = w.WriteString("</pre>\n")
	if codes := stepCodesIn(body); len(codes) > 0 {
		_, _ = w.WriteString(`<p class="steps">Bước:`)
		for _, c := range codes {
			fmt.Fprintf(w, ` <a href="#step-%s">%s</a>`, c, c)
		}
		_, _ = w.WriteString("</p>\n")
	}
	return ast.WalkSkipChildren, nil
}

func blockText(n *ast.FencedCodeBlock, src []byte) []byte {
	var sb strings.Builder
	for i := 0; i < n.Lines().Len(); i++ {
		seg := n.Lines().At(i)
		sb.Write(seg.Value(src))
	}
	return []byte(sb.String())
}

// renderCell như goldmark mặc định, thêm id="step-B1" cho ô cột đầu của bảng hành vi.
func (r *nodeRenderer) renderCell(w util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*east.TableCell)
	tag := "td"
	if _, ok := n.Parent().(*east.TableHeader); ok {
		tag = "th"
	}
	if !entering {
		fmt.Fprintf(w, "</%s>\n", tag)
		return ast.WalkContinue, nil
	}
	fmt.Fprintf(w, "<%s", tag)
	if n.Alignment != east.AlignNone {
		fmt.Fprintf(w, ` style="text-align:%s"`, n.Alignment.String())
	}
	if tag == "td" && r.behaviorTable != nil && n.Parent().Parent() == r.behaviorTable && n.PreviousSibling() == nil {
		if code := nodeText(n, src); stepCodeRe.MatchString(code) {
			fmt.Fprintf(w, ` id="step-%s"`, code)
		}
	}
	_, _ = w.WriteString(">")
	return ast.WalkContinue, nil
}
