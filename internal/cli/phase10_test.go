package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/anhvandev/doc-kit/assets"
	"github.com/anhvandev/doc-kit/internal/agentctx"
	"github.com/anhvandev/doc-kit/internal/skill"
)

func TestInitAgentContext(t *testing.T) {
	dir := t.TempDir()
	// CLAUDE.md đã có nội dung riêng; AGENTS.md chưa có.
	if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("# Dự án\n\nghi chú riêng\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := run(t, dir, "init", "--agent-context", "--json")
	if code != 0 {
		t.Fatalf("mã %d: %s", code, out)
	}
	var res []agentContextResult
	if json.Unmarshal([]byte(out), &res) != nil || len(res) != 2 ||
		res[0].File != "CLAUDE.md" || res[0].Result != agentctx.WriteUpdated ||
		res[1].File != "AGENTS.md" || res[1].Result != agentctx.WriteCreated {
		t.Fatalf("%s", out)
	}
	if _, err := os.Stat(filepath.Join(dir, "dk.toml")); err == nil {
		t.Fatal("--agent-context không được tạo dk.toml")
	}

	b, _ := os.ReadFile(filepath.Join(dir, "CLAUDE.md"))
	s := string(b)
	if !strings.HasPrefix(s, "# Dự án\n\nghi chú riêng\n\n<!-- dk:agent-context start version=dev hash=") || !strings.HasSuffix(s, "<!-- dk:agent-context end -->\n") {
		t.Fatalf("CLAUDE.md:\n%s", s)
	}
	content, _ := assets.FS.ReadFile("agent-context.md")
	if n := len(strings.Split(strings.TrimRight(string(content), "\n"), "\n")); n >= 120 {
		t.Fatalf("khối agent context %d dòng, phải dưới 120", n)
	}
	metas, _ := skill.List()
	for _, m := range metas {
		if !strings.Contains(s, "`"+m.Name+"`") {
			t.Errorf("thiếu skill %s trong bảng", m.Name)
		}
	}

	// Chạy lại: không đổi, không nhân đôi khối.
	out, _ = run(t, dir, "init", "--agent-context")
	if out != filepath.Join(dir, "CLAUDE.md")+": unchanged\n"+filepath.Join(dir, "AGENTS.md")+": unchanged\n" {
		t.Fatalf("chạy lại: %s", out)
	}
	b2, _ := os.ReadFile(filepath.Join(dir, "CLAUDE.md"))
	if string(b2) != s {
		t.Fatal("chạy lại làm đổi CLAUDE.md")
	}

	// Đã có dk.toml: chạy từ thư mục con vẫn ghi ở gốc dự án, khớp nơi doctor kiểm.
	if _, code := run(t, dir, "init"); code != 0 {
		t.Fatal("init")
	}
	sub := filepath.Join(dir, "docs", "features")
	if out, code := run(t, sub, "init", "--agent-context"); code != 0 || !strings.Contains(out, filepath.Join(dir, "CLAUDE.md")+": unchanged") {
		t.Fatalf("từ thư mục con: %s", out)
	}
	if _, err := os.Stat(filepath.Join(sub, "CLAUDE.md")); err == nil {
		t.Fatal("ghi nhầm vào thư mục con")
	}
}

func TestSelfCheck(t *testing.T) {
	out, code := run(t, t.TempDir(), "self-check", "--json")
	var res selfCheckResult
	if code != 0 || json.Unmarshal([]byte(out), &res) != nil {
		t.Fatalf("self-check: %s (%d)", out, code)
	}
	metas, _ := skill.List()
	if res.Version != "dev" || res.Skills != len(metas) || res.Templates < 41 || len(res.EmbedHash) != 64 || len(res.Errors) != 0 {
		t.Fatalf("%+v", res)
	}
	if out, code = run(t, t.TempDir(), "self-check"); code != 0 || !strings.Contains(out, "embed sha256: "+res.EmbedHash) {
		t.Fatalf("self-check text: %s", out)
	}
}

// rowOfDoctor tìm dòng theo tên mục.
func rowOfDoctor(rows []doctorRow, item string) doctorRow {
	for _, r := range rows {
		if r.Item == item {
			return r
		}
	}
	return doctorRow{Item: item, Status: "không có dòng"}
}

func doctor(t *testing.T, dir string, extra ...string) ([]doctorRow, int) {
	t.Helper()
	out, code := run(t, dir, append([]string{"doctor", "--json"}, extra...)...)
	var rows []doctorRow
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		t.Fatalf("doctor: %s (%d)", out, code)
	}
	return rows, code
}

func TestDoctor(t *testing.T) {
	// dk giả trên PATH để mục "dk trên PATH" đạt.
	fake := t.TempDir()
	name := "dk"
	if runtime.GOOS == "windows" { // LookPath trên Windows cần đuôi trong PATHEXT
		name = "dk.bat"
	}
	if err := os.WriteFile(filepath.Join(fake, name), []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", fake+string(os.PathListSeparator)+os.Getenv("PATH"))
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	dir := t.TempDir()
	rows, code := doctor(t, dir)
	if code != 3 || rows[0].Item != "dk.toml" || rows[0].OK || !strings.Contains(rows[0].Fix, "dk init") {
		t.Fatalf("chưa init: %+v (%d)", rows, code)
	}

	git(t, dir, "init", "-q")
	for _, args := range [][]string{{"init"}, {"init", "--agent-context"}, {"skill", "install"}, {"hook", "install"}} {
		if out, code := run(t, dir, args...); code != 0 {
			t.Fatalf("%v: %s", args, out)
		}
	}
	rows, code = doctor(t, dir)
	if code != 0 {
		t.Fatalf("đủ cài đặt vẫn lỗi: %+v", rows)
	}
	for _, item := range []string{"dk.toml", "docs/", "dk trên PATH", "git", "pre-commit", "agent context (CLAUDE.md)", "agent context (AGENTS.md)", "skill (claude, dự án)", "hook (claude, dự án)"} {
		if r := rowOfDoctor(rows, item); !r.OK {
			t.Errorf("%s: %+v", item, r)
		}
	}

	// Thiếu pre-commit, skill sửa tay, skill cũ phiên bản, hook gỡ.
	hooks := git(t, dir, "rev-parse", "--git-path", "hooks")
	if err := os.Remove(filepath.Join(dir, hooks, "pre-commit")); err != nil {
		t.Fatal(err)
	}
	adr := filepath.Join(dir, ".claude", "skills", "doc-adr", "SKILL.md")
	b, _ := os.ReadFile(adr)
	if err := os.WriteFile(adr, append(b, []byte("\nsửa tay\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	cr := filepath.Join(dir, ".claude", "skills", "doc-cr", "SKILL.md")
	b, _ = os.ReadFile(cr)
	if err := os.WriteFile(cr, []byte(strings.Replace(string(b), "dk_version: dev", "dk_version: 0.0.1", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, code := run(t, dir, "hook", "uninstall"); code != 0 {
		t.Fatal("hook uninstall")
	}
	agents := filepath.Join(dir, "AGENTS.md")
	b, _ = os.ReadFile(agents)
	if err := os.WriteFile(agents, []byte(strings.Replace(string(b), "## Working rules", "## Rules", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(dir, "CLAUDE.md")); err != nil {
		t.Fatal(err)
	}
	rows, code = doctor(t, dir)
	if code != 3 {
		t.Fatalf("phải mã 3: %+v", rows)
	}
	cases := map[string]string{
		"pre-commit":                "dk init --force",
		"skill doc-adr (claude)":    "--force",
		"skill doc-cr (claude)":     "dk skill install doc-cr",
		"hook (claude, dự án)":      "dk hook install",
		"agent context (AGENTS.md)": "dk init --agent-context",
		"agent context (CLAUDE.md)": "dk init --agent-context",
	}
	for item, fix := range cases {
		r := rowOfDoctor(rows, item)
		if r.OK || !strings.Contains(r.Fix, fix) {
			t.Errorf("%s: %+v, cần cách sửa chứa %q", item, r, fix)
		}
	}
	if r := rowOfDoctor(rows, "skill doc-cr (claude)"); !strings.Contains(r.Status, "cũ (v0.0.1)") {
		t.Errorf("doc-cr: %+v", r)
	}
	if r := rowOfDoctor(rows, "skill (claude, dự án)"); r.OK || !strings.Contains(r.Status, "10/12") {
		t.Errorf("tổng skill: %+v", r)
	}

	// Pre-commit có nhưng không gọi dk.
	if err := os.WriteFile(filepath.Join(dir, hooks, "pre-commit"), []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	rows, _ = doctor(t, dir)
	if r := rowOfDoctor(rows, "pre-commit"); r.OK || !strings.Contains(r.Status, "không gọi") {
		t.Errorf("pre-commit lạ: %+v", r)
	}

	// Scope toàn máy: HOME tạm chưa cài gì; target lạ vẫn in bảng.
	rows, code = doctor(t, dir, "--global")
	if r := rowOfDoctor(rows, "hook (claude, toàn máy)"); code != 3 || r.OK || !strings.Contains(r.Fix, "--global") {
		t.Errorf("doctor --global: %+v", r)
	}
	if out, code := run(t, dir, "doctor", "--target", "x"); code != 3 || !strings.Contains(out, "target x") {
		t.Errorf("doctor --target lạ: %s (%d)", out, code)
	}
	// Codex chưa cài gì: thiếu skill và hook, cách sửa có --target codex.
	rows, code = doctor(t, dir, "--target", "codex")
	if r := rowOfDoctor(rows, "hook (codex, dự án)"); code != 3 || r.OK || !strings.Contains(r.Fix, "--target codex") || !strings.Contains(r.Fix, "/hooks") {
		t.Errorf("doctor --target codex: %+v", r)
	}
	if rowOfDoctor(rows, "hook (claude, dự án)").Status != "không có dòng" {
		t.Error("--target codex không được kiểm claude")
	}
	// Có .codex/ thì doctor mặc định kiểm cả hai target.
	if out, code := run(t, dir, "hook", "install", "--target", "codex"); code != 0 {
		t.Fatalf("hook install codex: %s", out)
	}
	rows, _ = doctor(t, dir)
	if r := rowOfDoctor(rows, "hook (codex, dự án)"); !r.OK {
		t.Errorf("doctor tự kiểm codex khi có .codex/: %+v", r)
	}
	if rowOfDoctor(rows, "hook (claude, dự án)").Status == "không có dòng" {
		t.Error("doctor mặc định vẫn phải kiểm claude")
	}
	// Gỡ hết Codex: .codex/ rỗng còn lại không làm doctor kiểm codex nữa.
	if out, code := run(t, dir, "hook", "uninstall", "--target", "codex"); code != 0 {
		t.Fatalf("hook uninstall codex: %s", out)
	}
	rows, _ = doctor(t, dir)
	if rowOfDoctor(rows, "hook (codex, dự án)").Status != "không có dòng" {
		t.Error(".codex/ rỗng không được kích hoạt kiểm codex")
	}

	// Bảng chữ.
	out, code := run(t, dir, "doctor")
	if code != 3 || !strings.Contains(out, "mục") || !strings.Contains(out, "!! hook (claude, dự án)") {
		t.Fatalf("doctor text: %s", out)
	}
}
