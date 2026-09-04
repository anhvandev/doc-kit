package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/anhvandev/doc-kit/internal/skill"
)

// TestInstallEndToEnd chạy init, skill install, hook install rồi gỡ ngược trong
// repo tạm với HOME giả; không cần Claude Code.
func TestInstallEndToEnd(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", filepath.Join(dir, "home"))
	t.Setenv("USERPROFILE", filepath.Join(dir, "home"))
	git(t, dir, "init", "-q")

	out, code := run(t, dir, "init")
	if code != 0 || !strings.Contains(out, "Đã cài pre-commit") || !strings.Contains(out, "dk skill install") {
		t.Fatalf("init: %s (%d)", out, code)
	}
	if _, err := os.Stat(filepath.Join(dir, ".claude")); !os.IsNotExist(err) {
		t.Fatal("init đụng .claude/")
	}
	hook := filepath.Join(dir, ".git", "hooks", "pre-commit")
	if st, err := os.Stat(hook); err != nil || (runtime.GOOS != "windows" && st.Mode()&0o111 == 0) {
		t.Fatal("pre-commit thiếu hoặc không chạy được")
	}
	// Không có dk trên PATH: pre-commit cảnh báo và cho qua (cần sh; Windows
	// chạy hook qua sh của Git for Windows, không kiểm ở đây).
	if runtime.GOOS != "windows" {
		cmd := exec.Command("sh", hook)
		cmd.Dir = dir
		cmd.Env = []string{"PATH=/nonexistent"}
		hookOut, err := cmd.CombinedOutput()
		if err != nil || !strings.Contains(string(hookOut), "dk chưa cài") {
			t.Fatalf("pre-commit thiếu dk: %s %v", hookOut, err)
		}
	}
	// init lần hai với --force: pre-commit đã có nên in đoạn cần thêm.
	if out, code = run(t, dir, "init", "--force"); code != 0 || !strings.Contains(out, "đã có; thêm đoạn sau") {
		t.Fatalf("init --force: %s (%d)", out, code)
	}

	if out, code = run(t, dir, "skill", "list"); code != 0 || !strings.HasPrefix(out, "doc-adr\t") {
		t.Fatalf("skill list: %s", out)
	}
	if out, code = run(t, dir, "skill", "install"); code != 0 || !strings.Contains(out, "đã cài") {
		t.Fatalf("skill install: %s (%d)", out, code)
	}
	md, _ := os.ReadFile(filepath.Join(dir, ".claude", "skills", "doc-cr", "SKILL.md"))
	if !strings.Contains(string(md), "dk_installed_by: dk") || !strings.Contains(string(md), "dk_version: dev") {
		t.Fatalf("thiếu metadata dk_*:\n%s", md)
	}
	if out, code = run(t, dir, "hook", "install"); code != 0 {
		t.Fatalf("hook install: %s (%d)", out, code)
	}
	settings := filepath.Join(dir, ".claude", "settings.json")
	sb, _ := os.ReadFile(settings)
	if strings.Count(string(sb), "dk hook run") != 2 {
		t.Fatalf("settings.json:\n%s", sb)
	}
	out, code = run(t, dir, "skill", "status", "--json")
	var rows []skill.Row
	if code != 0 || json.Unmarshal([]byte(out), &rows) != nil || len(rows) != 24 ||
		rows[0].Scope != "dự án" || rows[0].State != skill.StateCurrent || rows[10].State != skill.StateCurrent ||
		rows[12].Scope != "toàn máy" || rows[12].State != skill.StateMissing {
		t.Fatalf("skill status: %s", out)
	}

	// hook run qua CLI: Write file mới trong docs/ bị từ chối.
	stdin = strings.NewReader(`{"tool_name":"Write","tool_input":{"file_path":"` + filepath.Join(dir, "docs", "features", "thu.md") + `"}}`)
	if out, code = run(t, dir, "hook", "run", "pre-write"); code != 0 || !strings.Contains(out, "deny") {
		t.Fatalf("hook run: %s (%d)", out, code)
	}
	if _, code = run(t, dir, "hook", "run"); code != 2 {
		t.Fatal("hook run thiếu tham số phải mã 2")
	}
	var help bytes.Buffer
	_ = Run([]string{"hook", "--help"}, &help)
	if strings.Contains(help.String(), "run <pre-write") {
		t.Fatal("hook run phải ẩn khỏi help")
	}

	if out, code = run(t, dir, "skill", "uninstall"); code != 0 || !strings.Contains(out, "đã gỡ") {
		t.Fatalf("skill uninstall: %s (%d)", out, code)
	}
	if out, code = run(t, dir, "hook", "uninstall"); code != 0 {
		t.Fatalf("hook uninstall: %s (%d)", out, code)
	}
	if rest, _ := os.ReadDir(filepath.Join(dir, ".claude")); len(rest) != 0 {
		t.Fatalf(".claude/ sau gỡ chưa rỗng: %v", rest)
	}
	if _, code = run(t, dir, "skill", "install", "--target", "x"); code != 2 {
		t.Fatal("target lạ phải mã 2")
	}

	// Hai target trong một lệnh; Codex có nhắc trust; gỡ sạch cả hai.
	if out, code = run(t, dir, "skill", "install", "--target", "claude,codex"); code != 0 || strings.Count(out, "doc-cr\t") != 2 {
		t.Fatalf("skill install hai target: %s (%d)", out, code)
	}
	if _, err := os.Stat(filepath.Join(dir, ".codex", "skills", "doc-cr", "SKILL.md")); err != nil {
		t.Fatal("chưa cài skill vào .codex/skills")
	}
	if out, code = run(t, dir, "hook", "install", "--target", "claude,codex"); code != 0 || !strings.Contains(out, "đã cài 4 hook") || !strings.Contains(out, "/hooks") {
		t.Fatalf("hook install hai target: %s (%d)", out, code)
	}
	hj, _ := os.ReadFile(filepath.Join(dir, ".codex", "hooks.json"))
	if strings.Count(string(hj), "dk hook run") != 2 || !strings.Contains(string(hj), `"matcher": "apply_patch"`) {
		t.Fatalf("hooks.json: %s", hj)
	}
	out, _ = run(t, dir, "skill", "status", "--target", "claude,codex", "--json")
	if strings.Count(out, `"target": "codex"`) == 0 || strings.Count(out, `"target": "claude"`) == 0 {
		t.Fatalf("skill status hai target: %s", out)
	}
	for _, args := range [][]string{{"skill", "uninstall", "--target", "claude,codex"}, {"hook", "uninstall", "--target", "claude,codex"}} {
		if out, code = run(t, dir, args...); code != 0 {
			t.Fatalf("%v: %s (%d)", args, out, code)
		}
	}
	for _, d := range []string{".claude", ".codex"} {
		if rest, _ := os.ReadDir(filepath.Join(dir, d)); len(rest) != 0 {
			t.Fatalf("%s/ sau gỡ chưa rỗng: %v", d, rest)
		}
	}
}

func TestInitWithoutGit(t *testing.T) {
	dir := t.TempDir()
	out, code := run(t, dir, "init", "--json")
	var res map[string]any
	if code != 0 || json.Unmarshal([]byte(out), &res) != nil {
		t.Fatalf("init: %s (%d)", out, code)
	}
	if pc, _ := res["pre_commit"].(map[string]any); pc["status"] != "no-git" {
		t.Fatalf("pre_commit: %v", res["pre_commit"])
	}
}
