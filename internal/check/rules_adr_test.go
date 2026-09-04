package check

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anhvandev/doc-kit/assets"
	"github.com/anhvandev/doc-kit/internal/config"
	"github.com/anhvandev/doc-kit/internal/docs"
	"github.com/anhvandev/doc-kit/internal/doctype"
)

func gitRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %s", args, out)
	}
}

const adrHead = "---\nid: ADR-0001\ntype: adr\ntitle: Dùng Postgres\nstatus: %s\nowner: v\ncreated: 2026-09-01\ncreated_by: dk\nsupersedes: \"\"\nsuperseded_by: \"\"\n---\n# ADR-0001\n\n## 2. Bối cảnh\n\nCần cơ sở dữ liệu.  \n"

func adrFindings(t *testing.T, root string) []Finding {
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
	for _, f := range Run(&Context{Root: root, DocsDir: "docs", Metas: metas, Reg: reg, Cfg: config.Check{WarnLines: 500, MaxLines: 800}}) {
		if f.Rule == "adr-immutable" {
			out = append(out, f)
		}
	}
	return out
}

func TestADRImmutable(t *testing.T) {
	root := t.TempDir()
	adr := filepath.Join(root, "docs", "adr", "ADR-0001-postgres.md")
	os.MkdirAll(filepath.Dir(adr), 0o755)
	write := func(status, body string) {
		os.WriteFile(adr, []byte(strings.Replace(adrHead, "%s", status, 1)+body), 0o644)
	}
	// Không có git: bỏ qua.
	write("accepted", "")
	if fs := adrFindings(t, root); len(fs) != 0 {
		t.Fatalf("không git phải bỏ qua: %v", fs)
	}
	gitRun(t, root, "init", "-q")
	// Chưa commit: bỏ qua.
	if fs := adrFindings(t, root); len(fs) != 0 {
		t.Fatalf("chưa có HEAD phải bỏ qua: %v", fs)
	}
	// HEAD là proposed: sửa thân rồi chốt được.
	write("proposed", "")
	gitRun(t, root, "add", ".")
	gitRun(t, root, "commit", "-q", "-m", "adr")
	write("accepted", "Thêm dòng khi chốt.\n")
	if fs := adrFindings(t, root); len(fs) != 0 {
		t.Fatalf("HEAD proposed phải cho sửa thân: %v", fs)
	}
	gitRun(t, root, "commit", "-q", "-am", "chốt")
	// HEAD accepted: chỉ frontmatter đổi (superseded), khoảng trắng cuối dòng và CRLF không tính.
	crlf := strings.ReplaceAll(strings.Replace(adrHead, "%s", "superseded", 1)+"Thêm dòng khi chốt.\n", "  \n", "\r\n")
	os.WriteFile(adr, []byte(crlf), 0o644)
	if fs := adrFindings(t, root); len(fs) != 0 {
		t.Fatalf("đổi frontmatter và định dạng không được báo: %v", fs)
	}
	// Đổi thân: lỗi.
	write("accepted", "Thêm dòng khi chốt. Sửa thêm.\n")
	fs := adrFindings(t, root)
	if len(fs) != 1 || fs[0].Level != Error || !strings.Contains(fs[0].Msg, "supersedes") {
		t.Fatalf("thân ADR đã chốt đổi phải lỗi: %v", fs)
	}
}
