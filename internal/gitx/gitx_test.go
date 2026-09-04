package gitx

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %s", args, out)
	}
}

func TestNumStat(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	p := filepath.Join(dir, "a.md")
	os.WriteFile(p, []byte("1\n2\n3\n"), 0o644)
	git(t, dir, "add", "a.md")
	git(t, dir, "commit", "-q", "-m", "init")

	os.WriteFile(p, []byte("1\n2x\n3\n4\n5\n6\n"), 0o644) // sửa 1 (+1 -1), thêm 3 -> +4 -1
	st, err := NumStat(dir, p)
	if err != nil || st != (Stat{Added: 4, Deleted: 1, Tracked: true}) {
		t.Fatalf("numstat = %+v, err=%v", st, err)
	}

	git(t, dir, "add", "a.md") // đã stage vẫn đếm
	st, _ = NumStat(dir, "a.md")
	if st.Added != 4 || st.Deleted != 1 {
		t.Fatalf("staged numstat = %+v", st)
	}

	os.WriteFile(filepath.Join(dir, "new.md"), []byte("x\n"), 0o644)
	st, _ = NumStat(dir, "new.md")
	if st.Tracked || st.NoGit {
		t.Fatalf("untracked = %+v", st)
	}

	os.WriteFile(filepath.Join(dir, "bộ lọc đơn.md"), []byte("x\n"), 0o644)
	changed, _ := ChangedDocs(dir, ".")
	if len(changed) != 3 || changed[1] != "bộ lọc đơn.md" {
		t.Fatalf("changed phải giữ nguyên dấu và khoảng trắng: %q", changed)
	}
	if ht, _ := HeadTime(dir); ht.IsZero() {
		t.Fatal("HeadTime phải khác zero sau commit")
	}
}

func TestNoGit(t *testing.T) {
	dir := t.TempDir()
	st, err := NumStat(dir, "x.md")
	if err != nil || !st.NoGit {
		t.Fatalf("ngoài repo phải NoGit: %+v %v", st, err)
	}
	if IsRepo(dir) {
		t.Fatal("IsRepo phải false")
	}
}
