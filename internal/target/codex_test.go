package target

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCodexSkillDir(t *testing.T) {
	t.Setenv("HOME", "/tmp/h")
	t.Setenv("USERPROFILE", "/tmp/h")
	t.Setenv("CODEX_HOME", "")
	c := &Codex{Root: "/p"}
	if d, _ := c.SkillDir(false); d != filepath.Join("/p", ".codex", "skills") {
		t.Fatal(d)
	}
	if d, _ := c.SkillDir(true); d != filepath.Join("/tmp/h", ".codex", "skills") {
		t.Fatal(d)
	}
	t.Setenv("CODEX_HOME", "/cx")
	if d, _ := c.SkillDir(true); d != filepath.Join("/cx", "skills") {
		t.Fatal("CODEX_HOME phải được ưu tiên: " + d)
	}
	if p, _ := c.HooksPath(true); p != filepath.Join("/cx", "hooks.json") {
		t.Fatal(p)
	}
	if _, err := (&Codex{}).SkillDir(false); err == nil {
		t.Fatal("scope dự án không gốc phải lỗi")
	}
}

func TestCodexHooksGolden(t *testing.T) {
	root := t.TempDir()
	c := &Codex{Root: root}
	path := filepath.Join(root, ".codex", "hooks.json")
	os.MkdirAll(filepath.Dir(path), 0o755)
	before, _ := os.ReadFile("testdata/hooks.before.json")
	after, _ := os.ReadFile("testdata/hooks.after.json")
	os.WriteFile(path, before, 0o644)

	for i := 0; i < 2; i++ { // lần hai idempotent
		if err := c.InstallHooks(false, entries); err != nil {
			t.Fatal(err)
		}
		got, _ := os.ReadFile(path)
		if string(got) != string(after) {
			t.Fatalf("lần %d, sau cài:\n%s\nmuốn:\n%s", i+1, got, after)
		}
	}
	cmds, err := c.InstalledHooks(false)
	if err != nil || len(cmds) != 2 {
		t.Fatalf("installed: %v %v", cmds, err)
	}
	if err := c.UninstallHooks(false); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != string(before) {
		t.Fatalf("sau gỡ:\n%s\nmuốn:\n%s", got, before)
	}
}

func TestCodexHooksFreshFile(t *testing.T) {
	root := t.TempDir()
	c := &Codex{Root: root}
	path := filepath.Join(root, ".codex", "hooks.json")
	if err := c.InstallHooks(false, entries); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("chưa tạo hooks.json")
	}
	if err := c.UninstallHooks(false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("file chỉ do dk tạo phải bị xóa khi rỗng")
	}
}
