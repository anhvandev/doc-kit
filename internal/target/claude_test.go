package target

import (
	"os"
	"path/filepath"
	"testing"
)

var entries = []HookEntry{
	{Event: "PreToolUse", Matcher: "Write", Command: "dk hook run pre-write"},
	{Event: "PostToolUse", Matcher: "Edit|Write", Command: "dk hook run post-edit"},
}

func TestRegistry(t *testing.T) {
	for _, n := range Names {
		if tg, err := Get(n, ""); err != nil || tg.Name() != n {
			t.Fatalf("%s: %v", n, err)
		}
	}
	if _, err := Get("x", ""); err == nil {
		t.Fatal("target lạ phải lỗi")
	}
}

func TestSkillDir(t *testing.T) {
	t.Setenv("HOME", "/tmp/h")
	t.Setenv("USERPROFILE", "/tmp/h")
	c := &Claude{Root: "/p"}
	if d, _ := c.SkillDir(false); d != filepath.Join("/p", ".claude", "skills") {
		t.Fatal(d)
	}
	if d, _ := c.SkillDir(true); d != filepath.Join("/tmp/h", ".claude", "skills") {
		t.Fatal(d)
	}
}

func TestHooksGolden(t *testing.T) {
	root := t.TempDir()
	c := &Claude{Root: root}
	path := filepath.Join(root, ".claude", "settings.json")
	os.MkdirAll(filepath.Dir(path), 0o755)
	before, _ := os.ReadFile("testdata/settings.before.json")
	after, _ := os.ReadFile("testdata/settings.after.json")
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
	if err := c.UninstallHooks(false); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != string(before) {
		t.Fatalf("sau gỡ:\n%s\nmuốn:\n%s", got, before)
	}
}

func TestHooksFreshFile(t *testing.T) {
	root := t.TempDir()
	c := &Claude{Root: root}
	path := filepath.Join(root, ".claude", "settings.json")
	if err := c.InstallHooks(false, entries); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("chưa tạo settings.json")
	}
	if err := c.UninstallHooks(false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("file chỉ do dk tạo phải bị xóa khi rỗng")
	}
	if err := c.UninstallHooks(false); err != nil {
		t.Fatal("gỡ khi không có file phải im lặng")
	}
}

func TestUninstallKeepsUserHookInSameGroup(t *testing.T) {
	root := t.TempDir()
	c := &Claude{Root: root}
	path := filepath.Join(root, ".claude", "settings.json")
	os.MkdirAll(filepath.Dir(path), 0o755)
	mixed := `{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Write",
        "hooks": [
          {
            "type": "command",
            "command": "dk hook run pre-write"
          },
          {
            "type": "command",
            "command": "my-own-guard.sh"
          }
        ]
      }
    ]
  }
}
`
	want := `{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Write",
        "hooks": [
          {
            "type": "command",
            "command": "my-own-guard.sh"
          }
        ]
      }
    ]
  }
}
`
	os.WriteFile(path, []byte(mixed), 0o644)
	if err := c.UninstallHooks(false); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != want {
		t.Fatalf("sau gỡ:\n%s", got)
	}
}

func TestNullHooksTreatedAsEmpty(t *testing.T) {
	root := t.TempDir()
	c := &Claude{Root: root}
	path := filepath.Join(root, ".claude", "settings.json")
	os.MkdirAll(filepath.Dir(path), 0o755)
	os.WriteFile(path, []byte("{\"hooks\": null}\n"), 0o644)
	if err := c.InstallHooks(false, entries); err != nil {
		t.Fatal(err)
	}
}
