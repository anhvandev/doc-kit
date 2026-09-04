package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anhvandev/doc-kit/assets"
	"github.com/anhvandev/doc-kit/internal/config"
	"github.com/anhvandev/doc-kit/internal/docs"
	"github.com/anhvandev/doc-kit/internal/doctype"
)

func ruleFindings(t *testing.T, root, rule string) []Finding {
	t.Helper()
	reg, err := doctype.Load(assets.FS)
	if err != nil {
		t.Fatal(err)
	}
	metas, err := docs.Scan(root, "docs")
	if err != nil {
		t.Fatal(err)
	}
	var out []Finding
	for _, f := range Run(&Context{Root: root, DocsDir: "docs", Metas: metas, Reg: reg, Cfg: config.Check{WarnLines: 500, MaxLines: 800}, Jargon: config.DefaultJargon}) {
		if f.Rule == rule {
			out = append(out, f)
		}
	}
	return out
}

func writeDoc(t *testing.T, root, rel, content string) {
	t.Helper()
	p := filepath.Join(root, filepath.FromSlash(rel))
	os.MkdirAll(filepath.Dir(p), 0o755)
	os.WriteFile(p, []byte(content), 0o644)
}

func TestNoJargon(t *testing.T) {
	root := t.TempDir()
	head := "---\ntype: release-brief\ntitle: x\nstatus: draft\nowner: v\ncreated: 2026-09-01\nsource: F-001\ncreated_by: dk\nfeature: F-001\nkind: feature\n---\n"
	writeDoc(t, root, "docs/release/briefs/F-001.md", head+"# x\n\n<!-- gợi ý: API trong chú thích không tính -->\n\nGọi endpoint để lấy dữ liệu.\n\nLại nhắc Endpoint lần hai và Database.\n\n```\nJSON trong khối mã không tính\n```\n")
	fs := ruleFindings(t, root, "no-jargon")
	if len(fs) != 2 || fs[0].Level != Warning || fs[0].Line != 16 || !strings.Contains(fs[0].Msg, `"endpoint"`) || fs[1].Line != 18 || !strings.Contains(fs[1].Msg, `"Database"`) {
		t.Fatalf("no-jargon: %+v", fs)
	}
	// Loại khác không chịu rule.
	writeDoc(t, root, "docs/features/F-001-x.md", strings.Replace(head, "type: release-brief", "type: feature-spec", 1)+"# x\n\nAPI endpoint database\n")
	if fs := ruleFindings(t, root, "no-jargon"); len(fs) != 2 {
		t.Fatalf("feature-spec không chịu no-jargon: %+v", fs)
	}
}

func TestEnvNoSecret(t *testing.T) {
	root := t.TempDir()
	head := "---\ntype: environment\ntitle: x\nstatus: draft\nowner: v\ncreated: 2026-09-01\ncreated_by: dk\n---\n"
	writeDoc(t, root, "docs/ops/environment.md", head+"# x\n\n<!-- gợi ý: SECRET=abc trong chú thích không tính -->\n\n```\nDB_PASSWORD=abc\nDB_HOST=<host của vault>\nexport API_KEY=\"sk-123\"\nEMPTY=\nTRICK=<a>thật<b>\n```\n\n| DB_PASSWORD | mật khẩu | có |\n")
	fs := ruleFindings(t, root, "env-no-secret")
	if len(fs) != 3 || fs[0].Level != Error || fs[0].Line != 14 || !strings.Contains(fs[0].Msg, "DB_PASSWORD") || fs[1].Line != 16 || !strings.Contains(fs[1].Msg, "API_KEY") || fs[2].Line != 18 || !strings.Contains(fs[2].Msg, "TRICK") {
		t.Fatalf("env-no-secret: %+v", fs)
	}
}
