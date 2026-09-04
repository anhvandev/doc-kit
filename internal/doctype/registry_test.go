package doctype

import (
	"strings"
	"testing"
	"testing/fstest"

	"github.com/anhvandev/doc-kit/assets"
)

func TestLoadEmbedded(t *testing.T) {
	reg, err := Load(assets.FS)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"adr", "architecture", "backup-dr", "brief", "changelog-product", "charter", "cr", "decision-log", "deployment",
		"design-component", "design-foundations", "design-pattern", "design-tokens", "environment", "faq", "feature-spec", "glossary", "idea",
		"interview", "meeting-notes", "mockup", "monitoring", "plan", "plan-phase", "postmortem", "product-overview", "release-brief",
		"release-notes", "report", "risk-register", "roadmap", "runbook", "test-case", "test-case-table", "test-report", "testing-strategy",
		"ui-spec", "ui-test-checklist", "user-guide", "userflow", "wireframe"}
	if got := strings.Join(reg.Names(), ","); got != strings.Join(want, ",") {
		t.Fatalf("danh sách loại: %s", got)
	}
	fs := reg["feature-spec"]
	if fs.Dir != "features" || fs.IDScheme != "seq:F-{n:03}" || fs.From["brief"]["purpose"] != "outcome" {
		t.Fatalf("feature-spec đọc sai: %+v", fs)
	}
	if reg["idea"].Subdir != "{yymmdd}-{slug}" {
		t.Fatalf("idea.subdir sai: %q", reg["idea"].Subdir)
	}
	if reg["mockup"].Kind != "html" || reg["mockup"].Ext() != ".html" || reg["design-tokens"].Ext() != ".json" || reg["adr"].Ext() != ".md" || reg["test-case"].Ext() != ".feature" {
		t.Fatalf("kind đọc sai: mockup=%q tokens=%q", reg["mockup"].Kind, reg["design-tokens"].Kind)
	}
}

func TestValidate(t *testing.T) {
	cases := map[string]string{
		"thiếu dir":      "[x]\nname = \"{id}.md\"\nid = \"none\"\n",
		"id sai":         "[x]\ndir = \"a\"\nname = \"{id}.md\"\nid = \"random\"\n",
		"thiếu template": "[y]\ndir = \"a\"\nname = \"y.md\"\nid = \"none\"\n",
		"from lạ":        "[x]\ndir = \"a\"\nname = \"x.md\"\nid = \"none\"\nfrom.zzz = { title = \"title\" }\n",
		"kind lạ":        "[x]\ndir = \"a\"\nname = \"x.md\"\nid = \"none\"\nkind = \"yaml\"\n",
	}
	for name, toml := range cases {
		fsys := fstest.MapFS{
			"types.toml":     {Data: []byte(toml)},
			"templates/x.md": {Data: []byte("x")},
		}
		if _, err := Load(fsys); err == nil {
			t.Errorf("%s: phải lỗi", name)
		}
	}
}
