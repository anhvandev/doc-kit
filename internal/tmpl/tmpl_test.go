package tmpl

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/anhvandev/doc-kit/assets"
	"github.com/anhvandev/doc-kit/internal/doctype"
	"github.com/anhvandev/doc-kit/internal/frontmatter"
)

func TestRenderAll(t *testing.T) {
	reg, err := doctype.Load(assets.FS)
	if err != nil {
		t.Fatal(err)
	}
	data := Data{ID: "X-001", Type: "t", Slug: "bo-loc", Title: "Bộ lọc", Owner: "an",
		Created: "2026-09-03", Updated: "2026-09-03 14:05", Source: "CR-260903-x", DKVersion: "dev",
		Feature: "F-001", Step: "B1", Layer: "atom"}
	for _, name := range reg.Names() {
		out, err := Render(name, data)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if bytes.Contains(out, []byte("{{")) {
			t.Errorf("%s: còn dấu {{ chưa render", name)
		}
		fm, _, ok := frontmatter.SplitFile("x"+reg[name].Ext(), out)
		if !ok {
			t.Fatalf("%s: frontmatter không hợp lệ", name)
		}
		for _, req := range reg[name].Required {
			if _, has := frontmatter.Get(fm, req); !has {
				t.Errorf("%s: thiếu trường bắt buộc %s", name, req)
			}
		}
		if frontmatter.GetString(fm, "type") != name {
			t.Errorf("%s: type sai", name)
		}
	}
}

func TestFeatureSpecSections(t *testing.T) {
	out, _ := Render("feature-spec", Data{Format: "spec", HasUI: true})
	if n := strings.Count(string(out), "\n## "); n != 10 {
		t.Fatalf("feature-spec cần 10 mục thân (11 kể frontmatter), có %d", n)
	}
	if !strings.Contains(string(out), "```mermaid") || !strings.Contains(string(out), "B1[") {
		t.Fatal("mục 3 phải có khối mermaid với nút B1")
	}
	// Biến thể format và has_ui: tập số mục cấp 2 theo đúng quy tắc spec-section-order.
	for _, c := range []struct {
		format string
		hasUI  bool
		want   string
		marker string
	}{
		{"spec", true, "1,2,3,4,5,6,7,8,9,10", "flowchart TD"},
		{"spec", false, "1,2,3,4,6,7,8,9,10", "flowchart TD"},
		{"use-case", true, "1,2,3,4,5,6,7,8,9,10", "sequenceDiagram"},
		{"story", true, "1,2,3,4,5,6,7,8,9,10", "```gherkin"},
		{"crud", true, "1,2,4,5,6,7,8,9,10", "## 4. Bảng field và quyền"},
		{"state", true, "1,2,3,4,5,6,7,8,9,10", "stateDiagram-v2"},
	} {
		out, err := Render("feature-spec", Data{Format: c.format, HasUI: c.hasUI})
		if err != nil {
			t.Fatal(err)
		}
		var nums []string
		for _, m := range regexp.MustCompile(`(?m)^## (\d+)\.`).FindAllStringSubmatch(string(out), -1) {
			nums = append(nums, m[1])
		}
		if got := strings.Join(nums, ","); got != c.want {
			t.Errorf("format=%s has_ui=%v: mục %s, muốn %s", c.format, c.hasUI, got, c.want)
		}
		if !strings.Contains(string(out), c.marker) {
			t.Errorf("format=%s: thiếu %q", c.format, c.marker)
		}
		if c.format == "use-case" && strings.Contains(string(out), "flowchart TD") || c.format == "crud" && strings.Contains(string(out), "mermaid") || c.format == "story" && strings.Contains(string(out), "**Given**") {
			t.Errorf("format=%s: còn nội dung của biến thể spec", c.format)
		}
	}
	cr, _ := Render("cr", Data{})
	if n := strings.Count(string(cr), "\n## "); n != 6 {
		t.Fatalf("cr cần 6 mục thân (7 kể frontmatter), có %d", n)
	}
	if strings.Count(string(cr), "| Có / Không |") != 6 {
		t.Fatal("bảng tác động CR cần 6 dòng")
	}
}
