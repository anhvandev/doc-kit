package render

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/anhvandev/doc-kit/internal/docs"
)

var update = flag.Bool("update", false, "ghi lại file .golden")

func renderFixture(t *testing.T, name string) Page {
	t.Helper()
	root := t.TempDir()
	docsDir := filepath.Join(root, "docs")
	if err := os.MkdirAll(filepath.Join(docsDir, "features"), 0o755); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(docsDir, "features", name)
	if err := os.WriteFile(src, b, 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := docs.Read(root, src)
	if err != nil {
		t.Fatal(err)
	}
	p, err := Options{DocsDir: docsDir, OutDir: filepath.Join(docsDir, "html")}.Render(m)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(docsDir, "html", "features", strings.TrimSuffix(name, ".md")+".html"); p.Out != want {
		t.Fatalf("Out = %s, muốn %s", p.Out, want)
	}
	return p
}

func TestGoldenPlain(t *testing.T) {
	p := renderFixture(t, "plain.md")
	if p.HasMermaid || bytes.Contains(p.HTML, []byte("mermaid.initialize")) {
		t.Fatal("file không có mermaid không được nhúng script mermaid")
	}
	golden := filepath.Join("testdata", "plain.golden.html")
	if *update {
		if err := os.WriteFile(golden, p.HTML, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(p.HTML, want) {
		t.Fatalf("HTML lệch golden; chạy go test ./internal/render -update để cập nhật\n%s", p.HTML)
	}
}

func TestSpecMermaidStepsLinks(t *testing.T) {
	p := renderFixture(t, "spec.md")
	h := string(p.HTML)
	if !p.HasMermaid || !strings.Contains(h, `mermaid.initialize({startOnLoad: true`) || !strings.Contains(h, `globalThis["mermaid"]`) {
		t.Fatal("thiếu script mermaid nhúng")
	}
	if !strings.Contains(h, `<pre class="mermaid">flowchart TD`) || strings.Contains(h, `<code class="language-mermaid">`) {
		t.Fatal("khối mermaid phải ra <pre class=\"mermaid\">")
	}
	if want := `<p class="steps">Bước: <a href="#step-B1">B1</a> <a href="#step-B2">B2</a> <a href="#step-B3">B3</a></p>`; !strings.Contains(h, want) {
		t.Fatalf("dòng Bước sai:\n%s", h)
	}
	// id chỉ gắn ở bảng hành vi (bảng đầu tiên có cột đầu B\d+), mỗi mã một lần.
	for id, want := range map[string]int{"step-B1": 1, "step-B2": 1, "step-B9": 1} {
		if n := strings.Count(h, `id="`+id+`"`); n != want {
			t.Fatalf("id %s xuất hiện %d lần, muốn %d", id, n, want)
		}
	}
	if !strings.Contains(h, "<td id=\"step-B1\">B1</td>\n<td>mở</td>") || strings.Contains(h, "<td id=\"step-B1\">B1</td>\n<td>Danh sách</td>") {
		t.Fatal("bảng hành vi phải là bảng dưới tiêu đề chứa 'hành vi', không phải bảng Giao diện đứng trước")
	}
	if !strings.Contains(h, "Knut Sveidqvist") {
		t.Fatal("trang có mermaid phải kèm thông báo giấy phép MIT")
	}
	// Liên kết: .md trong docs đổi .html giữ neo; ngoài docs trỏ file gốc; http và neo giữ nguyên.
	for _, want := range []string{`href="../adr/ADR-0001-x.html#muc-2"`, `href="../../../plans/p.md"`, `href="https://example.com/a.md"`, `href="#x"`} {
		if !strings.Contains(h, want) {
			t.Fatalf("thiếu %s", want)
		}
	}
	// HTML thô bị escape, fenced code thường giữ nguyên.
	if strings.Contains(h, "<script>alert") || !strings.Contains(h, `<code class="language-go">fmt.Println(&quot;&lt;b&gt;&quot;)`) {
		t.Fatalf("escape sai:\n%s", h)
	}
	if !strings.Contains(h, "<title>Bộ lọc: đơn hàng</title>") || !strings.Contains(h, `<td class="id">F-001</td>`) || strings.Contains(h, "acceptance") {
		t.Fatal("metadata sai: cần title, id; bỏ trường rỗng")
	}
}

func TestTOC(t *testing.T) {
	root := t.TempDir()
	docsDir := filepath.Join(root, "docs")
	if err := os.MkdirAll(docsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(name, body string) docs.Meta {
		src := filepath.Join(docsDir, name)
		if err := os.WriteFile(src, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		m, err := docs.Read(root, src)
		if err != nil {
			t.Fatal(err)
		}
		return m
	}
	opts := Options{DocsDir: docsDir, OutDir: filepath.Join(docsDir, "html")}

	// Đủ tiêu đề: mục lục cấp 2 và 3 theo thứ tự, id trùng goldmark, bỏ cấp 1 và 4.
	long := "# Tiêu đề\n\n## 1. Mục đích\n\n### 1.1 Chi tiết `mã`\n\n#### Bỏ qua\n\n## 2. Phạm vi\n\n## 2. Phạm vi\n"
	p, err := opts.Render(write("long.md", long))
	if err != nil {
		t.Fatal(err)
	}
	h := string(p.HTML)
	want := "<nav class=\"toc\" aria-label=\"Mục lục\">\n<ul>\n" +
		"<li class=\"l2\"><a href=\"#1-mc-ch\">1. Mục đích</a></li>\n" +
		"<li class=\"l3\"><a href=\"#11-chi-tit-m\">1.1 Chi tiết mã</a></li>\n" +
		"<li class=\"l2\"><a href=\"#2-phm-vi\">2. Phạm vi</a></li>\n" +
		"<li class=\"l2\"><a href=\"#2-phm-vi-1\">2. Phạm vi</a></li>\n</ul>\n</nav>"
	if !strings.Contains(h, want) {
		t.Fatalf("mục lục sai:\n%s", h)
	}
	if strings.Contains(h, "Bỏ qua</a>") || strings.Contains(h, "Tiêu đề</a>") {
		t.Fatal("mục lục không được có cấp 1 hoặc cấp 4")
	}
	for _, id := range []string{`id="1-mc-ch"`, `id="2-phm-vi-1"`} {
		if !strings.Contains(h, id) {
			t.Fatalf("thân thiếu neo %s", id)
		}
	}
	// Không frontmatter vẫn có aside vì có mục lục.
	if !strings.Contains(h, "<aside class=\"meta\">\n<nav") {
		t.Fatal("aside phải có khi chỉ có mục lục, không có bảng metadata")
	}

	// Dưới 3 tiêu đề: không mục lục.
	p, err = opts.Render(write("short.md", "# T\n\n## A\n\n## B\n"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(p.HTML), `class="toc"`) {
		t.Fatal("dưới 3 tiêu đề không được có mục lục")
	}
}

func TestStepCodes(t *testing.T) {
	b, _ := os.ReadFile(filepath.Join("testdata", "spec.md"))
	_, body, _ := splitFM(b)
	d, tb := StepCodes(body)
	if strings.Join(d, ",") != "B1,B2,B3" || strings.Join(tb, ",") != "B1,B2,B9" {
		t.Fatalf("diagram=%v table=%v", d, tb)
	}
	// Mã có hậu tố chữ thường (bước chèn giữa) được nhận ở cả sơ đồ và bảng.
	d, tb = StepCodes([]byte("```mermaid\nflowchart TD\n  B1 --> B2a --> B2\n```\n\n| Mã | Hành vi |\n|---|---|\n| B1 | |\n| B2a | |\n| B2 | |\n"))
	if strings.Join(d, ",") != "B1,B2a,B2" || strings.Join(tb, ",") != "B1,B2a,B2" {
		t.Fatalf("hậu tố: diagram=%v table=%v", d, tb)
	}
}

func splitFM(b []byte) ([]byte, []byte, bool) {
	parts := bytes.SplitN(b, []byte("---\n"), 3)
	if len(parts) < 3 {
		return nil, b, false
	}
	return parts[1], parts[2], true
}

func TestIndex(t *testing.T) {
	root := t.TempDir()
	docsDir := filepath.Join(root, "docs")
	mk := func(rel, fm string) docs.Meta {
		p := filepath.Join(docsDir, filepath.FromSlash(rel))
		_ = os.MkdirAll(filepath.Dir(p), 0o755)
		_ = os.WriteFile(p, []byte("---\n"+fm+"\n---\n# x\n"), 0o644)
		m, err := docs.Read(root, p)
		if err != nil {
			t.Fatal(err)
		}
		return m
	}
	metas := []docs.Meta{
		mk("features/F-002-b.md", "id: F-002\ntype: feature-spec\ntitle: B\nstatus: draft"),
		mk("features/F-001-a.md", "id: F-001\ntype: feature-spec\ntitle: A\nstatus: review\nowner: v"),
		mk("adr/ADR-0001-x.md", "id: ADR-0001\ntype: adr\ntitle: X\nstatus: accepted"),
		mk("features/README.md", "generated: true"),
	}
	b, err := Options{DocsDir: docsDir, OutDir: filepath.Join(docsDir, "html")}.Index("Chỉ mục", metas)
	if err != nil {
		t.Fatal(err)
	}
	h := string(b)
	rows := regexp.MustCompile(`<tr><td class="id">`).FindAllString(h, -1)
	if len(rows) != 3 {
		t.Fatalf("muốn 3 dòng, được %d\n%s", len(rows), h)
	}
	if strings.Index(h, `<h2 id="adr">`) > strings.Index(h, `<h2 id="feature-spec">`) {
		t.Fatal("nhóm phải theo thứ tự chữ cái của type")
	}
	if strings.Index(h, "F-002") < strings.Index(h, "F-001") || !strings.Contains(h, `href="features/F-001-a.html"`) {
		t.Fatalf("thứ tự hoặc liên kết sai:\n%s", h)
	}
	if strings.Contains(h, "mermaid.initialize") {
		t.Fatal("index không nhúng mermaid")
	}
}

// TestStyleRules kiểm style.css theo html-style.md.
func TestStyleRules(t *testing.T) {
	css, err := os.ReadFile(filepath.Join("..", "..", "assets", "html", "style.css"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(css)
	for _, bad := range []string{"background-clip: text", "background-clip:text", "@import", "url(http"} {
		if strings.Contains(s, bad) {
			t.Fatalf("style.css chứa %q, vi phạm html-style.md", bad)
		}
	}
	if regexp.MustCompile(`(?m)^\s*color\s*:[^;]*gradient`).MatchString(s) {
		t.Fatal("style.css dùng gradient trong color")
	}
	if !strings.Contains(s, "prefers-reduced-motion") || !strings.Contains(s, "@media print") {
		t.Fatal("style.css thiếu prefers-reduced-motion hoặc @media print")
	}
}
