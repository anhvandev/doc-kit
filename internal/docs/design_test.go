package docs

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anhvandev/doc-kit/internal/frontmatter"
)

func readFMFile(t *testing.T, path string) map[string]any {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	fm, _, ok := frontmatter.SplitFile(path, b)
	if !ok {
		t.Fatalf("%s: không có metadata", path)
	}
	return frontmatter.Map(fm)
}

func TestDesignKinds(t *testing.T) {
	reg, o := setup(t)
	// design-tokens: JSON với $dk là khóa đầu, --set ghi vào $dk, phần tokens giữ nguyên.
	o.Set = map[string]string{"title": "Tokens: dự án", "owner": "an"}
	r, err := New(reg, "design-tokens", "tokens", o)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(r.Path) != "tokens.json" || filepath.Base(filepath.Dir(r.Path)) != "tokens" {
		t.Fatalf("đường dẫn: %s", r.Path)
	}
	b, _ := os.ReadFile(r.Path)
	if !json.Valid(b) || !strings.HasPrefix(string(b), "{\n  \"$dk\": {") {
		t.Fatalf("tokens.json không hợp lệ:\n%s", b)
	}
	fm := readFMFile(t, r.Path)
	if fm["title"] != "Tokens: dự án" || fm["type"] != "design-tokens" || fm["created_by"] != "dk" {
		t.Fatalf("$dk: %v", fm)
	}
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	if m["color"] == nil || m["space"] == nil {
		t.Fatal("thiếu nhóm token")
	}

	// mockup: HTML với chú thích dk, tên theo feature và step; thiếu --set là lỗi rõ.
	o.Set = map[string]string{"feature": "F-012", "step": "B3", "title": "Danh sách rỗng"}
	r, err = New(reg, "mockup", "danh-sach", o)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(r.Path) != "F-012-B3.html" {
		t.Fatalf("tên mockup: %s", r.Path)
	}
	b, _ = os.ReadFile(r.Path)
	if !strings.HasPrefix(string(b), "<!-- dk:\ntype: mockup\n") || !strings.Contains(string(b), `href="../tokens/tokens.css"`) || strings.Contains(string(b), "{{") {
		t.Fatalf("mockup.html:\n%s", b)
	}
	fm = readFMFile(t, r.Path)
	if fm["feature"] != "F-012" || fm["step"] != "B3" || fm["external"] != "" {
		t.Fatalf("metadata mockup: %v", fm)
	}
	o.Set = map[string]string{"feature": "F-012"}
	if _, err := New(reg, "mockup", "x", o); err == nil || !strings.Contains(err.Error(), "--set step=") {
		t.Fatalf("thiếu step phải lỗi: %v", err)
	}
	o.Set = map[string]string{"feature": "../esc", "step": "B3"}
	if _, err := New(reg, "mockup", "x", o); err == nil || !strings.Contains(err.Error(), "feature") {
		t.Fatalf("feature có / phải lỗi: %v", err)
	}
	o.Set = map[string]string{"feature": "F-012", "step": "b3"}
	if _, err := New(reg, "mockup", "x", o); err == nil || !strings.Contains(err.Error(), "step") {
		t.Fatalf("step sai dạng phải lỗi: %v", err)
	}

	// mockup external: chỉ liên kết và ảnh, không tokens.css.
	o.Set = map[string]string{"feature": "F-012", "step": "B4", "external": "https://figma.com/file/x"}
	r, err = New(reg, "mockup", "figma", o)
	if err != nil {
		t.Fatal(err)
	}
	b, _ = os.ReadFile(r.Path)
	if strings.Contains(string(b), "tokens.css") || !strings.Contains(string(b), `href="https://figma.com/file/x"`) || !strings.Contains(string(b), `src="F-012-B4.png"`) {
		t.Fatalf("mockup external:\n%s", b)
	}
	if readFMFile(t, r.Path)["external"] != "https://figma.com/file/x" {
		t.Fatal("external chưa vào metadata")
	}

	// design-component: thư mục theo layer; layer lạ là lỗi.
	o.Set = map[string]string{"layer": "atom"}
	r, err = New(reg, "design-component", "button", o)
	if err != nil {
		t.Fatal(err)
	}
	if rel, _ := filepath.Rel(o.DocsDir, r.Path); filepath.ToSlash(rel) != "design/atoms/button.md" {
		t.Fatalf("đường dẫn component: %s", rel)
	}
	if fm := readFM(t, r.Path); fm["layer"] != "atom" {
		t.Fatalf("layer: %v", fm)
	}
	o.Set = map[string]string{"layer": "page"}
	if _, err := New(reg, "design-component", "x", o); err == nil || !strings.Contains(err.Error(), "layer") {
		t.Fatalf("layer lạ phải lỗi: %v", err)
	}

	// userflow --from feature-spec: feature chép từ id, tên file theo feature.
	o.Set = nil
	spec, err := New(reg, "feature-spec", "bo-loc", o)
	if err != nil {
		t.Fatal(err)
	}
	o.From = spec.Path
	r, err = New(reg, "userflow", "bo-loc", o)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(r.Path) != "F-001-flow.md" {
		t.Fatalf("tên userflow: %s", r.Path)
	}
	if fm := readFM(t, r.Path); fm["feature"] != "F-001" || fm["source"] != "F-001" || fm["title"] != "Bo loc" {
		t.Fatalf("userflow từ spec: %v", fm)
	}
	// wireframe --from userflow: feature chép, step từ --set.
	o.From, o.Set = r.Path, map[string]string{"step": "B2"}
	r, err = New(reg, "wireframe", "bo-loc", o)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(r.Path) != "F-001-B2.md" || readFM(t, r.Path)["source"] != "flows/F-001-flow.md" {
		t.Fatalf("wireframe: %s %v", r.Path, readFM(t, r.Path))
	}
}

func TestScanKinds(t *testing.T) {
	reg, o := setup(t)
	o.Set = map[string]string{"feature": "F-001", "step": "B1"}
	if _, err := New(reg, "mockup", "x", o); err != nil {
		t.Fatal(err)
	}
	o.Set = nil
	if _, err := New(reg, "design-tokens", "x", o); err != nil {
		t.Fatal(err)
	}
	css := filepath.Join(o.DocsDir, "design", "tokens", "tokens.css")
	_ = os.WriteFile(css, []byte(":root{}"), 0o644)
	metas, err := Scan(filepath.Dir(o.DocsDir), "docs")
	if err != nil {
		t.Fatal(err)
	}
	var rels []string
	for _, m := range metas {
		if !m.HasFM || m.Type == "" {
			t.Errorf("%s: không đọc được metadata", m.Rel)
		}
		rels = append(rels, m.Rel)
	}
	if strings.Join(rels, ",") != "docs/design/mockups/F-001-B1.html,docs/design/tokens/tokens.json" {
		t.Fatalf("Scan: %v", rels)
	}
	if metas[0].IsMarkdown() {
		t.Fatal("mockup không phải Markdown")
	}
}
