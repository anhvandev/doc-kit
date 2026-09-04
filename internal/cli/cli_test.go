package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

// run chạy CLI trong thư mục dir, trả stdout và mã thoát.
func run(t *testing.T, dir string, args ...string) (string, int) {
	t.Helper()
	var out bytes.Buffer
	err := Run(append([]string{"--cwd", dir}, args...), &out)
	if err == nil {
		return out.String(), 0
	}
	var ee *ExitError
	if errors.As(err, &ee) {
		t.Logf("dk %s: %s", strings.Join(args, " "), ee.Msg)
		return out.String(), ee.Code
	}
	t.Fatalf("lỗi không có mã thoát: %v", err)
	return "", -1
}

func git(t *testing.T, dir string, args ...string) string {
	return gitEnv(t, dir, nil, args...)
}

func gitEnv(t *testing.T, dir string, env []string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
	cmd.Env = append(cmd.Env, env...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %s", args, out)
	}
	return strings.TrimSpace(string(out))
}

func TestVersion(t *testing.T) {
	out, code := run(t, t.TempDir(), "--version")
	if code != 0 || strings.TrimSpace(out) != "dev" {
		t.Fatalf("--version: %q (%d)", out, code)
	}
}

func TestUsageErrors(t *testing.T) {
	dir := t.TempDir()
	if _, code := run(t, dir, "new"); code != 2 {
		t.Fatalf("thiếu tham số phải mã 2, được %d", code)
	}
	if _, code := run(t, dir, "--nope"); code != 2 {
		t.Fatalf("cờ lạ phải mã 2, được %d", code)
	}
	if _, code := run(t, dir, "bogus"); code != 2 {
		t.Fatalf("lệnh lạ phải mã 2, được %d", code)
	}
	if _, code := run(t, dir, "new", "cr", "x"); code != 1 {
		t.Fatalf("chưa init phải mã 1, được %d", code)
	}
}

func TestTemplate(t *testing.T) {
	dir := t.TempDir()
	out, code := run(t, dir, "template", "list", "--json")
	var infos []typeInfo
	if code != 0 || json.Unmarshal([]byte(out), &infos) != nil || len(infos) != 41 {
		t.Fatalf("template list: %s (%d)", out, code)
	}
	out, code = run(t, dir, "template", "show", "feature-spec")
	if code != 0 || strings.Count(out, "\n## ") != 10 || !strings.Contains(out, "```mermaid") {
		t.Fatalf("template show feature-spec thiếu mục: %d", strings.Count(out, "\n## "))
	}
	if _, code = run(t, dir, "template", "show", "nope"); code != 1 {
		t.Fatalf("loại lạ phải mã 1, được %d", code)
	}
}

func TestEndToEnd(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init", "-q")

	out, code := run(t, dir, "init")
	if code != 0 || !strings.Contains(out, "dk skill install") {
		t.Fatalf("init: %s (%d)", out, code)
	}
	for _, p := range []string{"dk.toml", "docs/CHANGELOG-DOCS.md", "docs/features/.gitkeep", "docs/design/mockups/.gitkeep", "docs/html/.gitkeep", "plans/.gitkeep"} {
		if _, err := os.Stat(filepath.Join(dir, p)); err != nil {
			t.Errorf("init thiếu %s", p)
		}
	}
	if _, code = run(t, dir, "init"); code != 1 {
		t.Fatalf("init lần hai phải mã 1, được %d", code)
	}
	os.WriteFile(filepath.Join(dir, "dk.toml"), []byte("project_name = \"shop\"\nid_prefix = \"SHOP-\"\ndocs_dir = \"docs\"\n"), 0o644)
	if _, code = run(t, dir, "init", "--force"); code != 0 {
		t.Fatalf("init --force phải qua, được %d", code)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, "dk.toml")); !strings.Contains(string(b), "id_prefix = 'SHOP-'") {
		t.Fatalf("init --force phải giữ cấu hình cũ:\n%s", b)
	}
	nested := filepath.Join(dir, "docs", "features")
	if _, code = run(t, nested, "init"); code != 1 {
		t.Fatal("init trong dự án có sẵn phải bị từ chối")
	}

	// .gitkeep do init tạo không bị tính pending
	if out, code = run(t, dir, "changelog", "pending"); code != 0 || strings.TrimSpace(out) != "" {
		t.Fatalf("pending sau init phải sạch: %q (%d)", out, code)
	}
	if _, code = run(t, dir, "changelog", "add", "docs/adr/.gitkeep"); code != 1 {
		t.Fatal("add .gitkeep phải bị từ chối")
	}
	git(t, dir, "add", ".")
	git(t, dir, "commit", "-q", "-m", "init")

	// dk new từ thư mục con vẫn tìm được gốc
	sub := filepath.Join(dir, "docs", "features")
	out, code = run(t, sub, "new", "feature-spec", "bo-loc-don-hang", "--json")
	var res map[string]string
	if code != 0 || json.Unmarshal([]byte(out), &res) != nil {
		t.Fatalf("new: %s (%d)", out, code)
	}
	if _, code := run(t, dir, "new", "cr", "Sai Slug"); code != 2 {
		t.Fatalf("slug sai phải mã 2, được %d", code)
	}
	if res["id"] != "SHOP-F-001" || res["path"] != "docs/features/SHOP-F-001-bo-loc-don-hang.md" {
		t.Fatalf("new: %v", res)
	}
	spec := filepath.Join(dir, res["path"])
	if out, _ = run(t, dir, "new", "feature-spec", "khac"); !strings.Contains(out, "SHOP-F-002") {
		t.Fatalf("lần hai phải F-002: %s", out)
	}

	// file mới, chưa commit: changelog ghi "mới, N dòng"; pending sạch với file này
	out, code = run(t, dir, "changelog", "add", res["path"], "--summary", "Khung spec", "--source", "brief-x")
	if code != 0 || !strings.Contains(out, "| features/SHOP-F-001-bo-loc-don-hang.md | mới, ") {
		t.Fatalf("changelog add file mới: %s (%d)", out, code)
	}
	out, code = run(t, dir, "changelog", "pending")
	if code != 1 || strings.TrimSpace(out) != "features/SHOP-F-002-khac.md" {
		t.Fatalf("pending phải liệt kê F-002: %q (%d)", out, code)
	}
	run(t, dir, "changelog", "add", "docs/features/SHOP-F-002-khac.md")
	if out, code = run(t, dir, "changelog", "pending"); code != 0 || strings.TrimSpace(out) != "" {
		t.Fatalf("pending phải sạch: %q (%d)", out, code)
	}

	git(t, dir, "add", ".")
	git(t, dir, "commit", "-q", "-m", "spec")

	// sửa 2 dòng, add, số dòng khớp git diff --numstat
	b, _ := os.ReadFile(spec)
	s := strings.Replace(string(b), "- Tác nhân: ", "- Tác nhân: Nhân viên bán hàng", 1)
	s = strings.Replace(s, "- Điều kiện tiên quyết: ", "- Điều kiện tiên quyết: Đã đăng nhập", 1)
	s = regexp.MustCompile(`(?m)^updated: .*$`).ReplaceAllString(s, "updated: 2000-01-01 00:00")
	os.WriteFile(spec, []byte(s), 0o644)
	if _, code = run(t, dir, "changelog", "add", spec, "--summary", "a | b"); code != 2 {
		t.Fatal("tóm tắt chứa ' | ' phải mã 2")
	}
	if _, code = run(t, dir, "changelog", "add", spec, "--source", "x\ny"); code != 2 {
		t.Fatal("nguồn chứa xuống dòng phải mã 2")
	}
	out, code = run(t, dir, "changelog", "add", spec, "--summary", "Điền tác nhân", "--source", "brief-x")
	if code != 0 {
		t.Fatalf("changelog add: %s (%d)", out, code)
	}
	numstat := strings.Fields(git(t, dir, "diff", "--numstat", "HEAD", "--", res["path"]))
	want := "| +" + numstat[0] + " −" + numstat[1] + " |"
	if !strings.Contains(out, want) {
		t.Fatalf("số dòng %q không khớp git %v", out, numstat)
	}
	cl, _ := os.ReadFile(filepath.Join(dir, "docs", "CHANGELOG-DOCS.md"))
	if !strings.HasPrefix(string(cl), "# Changelog tài liệu\n\n## ") || !strings.Contains(string(cl), want+" Điền tác nhân | brief-x") {
		t.Fatalf("changelog sai:\n%s", cl)
	}
	b2, _ := os.ReadFile(spec)
	if !strings.Contains(string(b2), "\nupdated: ") || strings.Contains(string(b2), "updated: 2000-01-01") {
		t.Fatal("updated: phải được bump")
	}
	if out, code = run(t, dir, "changelog", "pending"); code != 0 || strings.TrimSpace(out) != "" {
		t.Fatalf("pending sau add phải sạch: %q (%d)", out, code)
	}

	// gộp trong 10 phút: thêm lần nữa cùng file cùng nguồn -> vẫn 1 dòng, hai tóm tắt
	// thật nối bằng "; "; lần không có tóm tắt giữ nguyên tóm tắt đã gộp.
	out, _ = run(t, dir, "changelog", "add", spec, "--summary", "Điền tác nhân và điều kiện", "--source", "brief-x")
	out2, _ := run(t, dir, "changelog", "add", spec, "--source", "brief-x")
	if !strings.Contains(out2, "| Điền tác nhân; Điền tác nhân và điều kiện |") {
		t.Fatalf("stdout phải in mục đã gộp (nối tóm tắt cũ): %s", out2)
	}
	cl, _ = os.ReadFile(filepath.Join(dir, "docs", "CHANGELOG-DOCS.md"))
	if n := strings.Count(string(cl), "features/SHOP-F-001-bo-loc-don-hang.md"); n != 2 {
		t.Fatalf("phải gộp thành 1 dòng cho lần sửa này (cộng 1 dòng 'mới' lúc tạo), có %d:\n%s", n, cl)
	}
	if !strings.Contains(string(cl), "Điền tác nhân và điều kiện") || strings.Contains(string(cl), "| Điền tác nhân |") {
		t.Fatalf("tóm tắt mới phải thay tóm tắt cũ:\n%s", cl)
	}

	// sửa tiếp không add: pending liệt kê, mã 1. Changelog chỉ lưu phút nên commit
	// cùng phút với dòng vừa ghi sẽ không phân biệt được; ép mốc commit sang phút sau.
	git(t, dir, "add", ".")
	future := time.Now().Add(time.Minute).Format(time.RFC3339)
	gitEnv(t, dir, []string{"GIT_COMMITTER_DATE=" + future}, "commit", "-q", "-m", "fill")
	os.WriteFile(spec, append(b2, []byte("\nthêm dòng\n")...), 0o644)
	out, code = run(t, dir, "changelog", "pending", "--json")
	if code != 1 || !strings.Contains(out, "features/SHOP-F-001-bo-loc-don-hang.md") {
		t.Fatalf("pending phải báo F-001: %s (%d)", out, code)
	}

	// từ chối file ngoài docs/, trong html/, và chỉ mục sinh ra
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("x"), 0o644)
	if _, code = run(t, dir, "changelog", "add", "README.md"); code != 1 {
		t.Fatal("file ngoài docs/ phải bị từ chối")
	}
	os.WriteFile(filepath.Join(dir, "docs", "html", "a.html"), []byte("x"), 0o644)
	if _, code = run(t, dir, "changelog", "add", "docs/html/a.html"); code != 1 {
		t.Fatal("docs/html/ phải bị từ chối")
	}
	os.WriteFile(filepath.Join(dir, "docs", "index.md"), []byte("---\ngenerated: true\n---\n"), 0o644)
	if _, code = run(t, dir, "changelog", "add", "docs/index.md"); code != 1 {
		t.Fatal("file generated phải bị từ chối")
	}
	if out, _ = run(t, dir, "changelog", "pending"); strings.Contains(out, "index.md") || strings.Contains(out, "html/") {
		t.Fatalf("pending không được liệt kê html/ hay generated: %s", out)
	}
}

func TestIntakeFrom(t *testing.T) {
	dir := t.TempDir()
	run(t, dir, "init")
	out, code := run(t, dir, "new", "idea", "bo-loc", "--set", "title=Bộ lọc đơn hàng", "--json")
	var idea map[string]string
	if code != 0 || json.Unmarshal([]byte(out), &idea) != nil {
		t.Fatalf("new idea: %s", out)
	}
	out, code = run(t, dir, "new", "brief", "bo-loc", "--from", idea["path"], "--json")
	var brief map[string]string
	if code != 0 || json.Unmarshal([]byte(out), &brief) != nil {
		t.Fatalf("new brief --from: %s", out)
	}
	if filepath.Dir(brief["path"]) != filepath.Dir(idea["path"]) {
		t.Fatalf("brief phải cùng thư mục idea: %v", brief)
	}
	b, _ := os.ReadFile(filepath.Join(dir, brief["path"]))
	if !strings.Contains(string(b), "title: Bộ lọc đơn hàng") || !strings.Contains(string(b), "source: "+filepath.Base(filepath.Dir(idea["path"]))+"/idea.md") {
		t.Fatalf("brief chép trường sai:\n%s", b)
	}
	// không git: changelog ghi "không git"
	out, code = run(t, dir, "changelog", "add", brief["path"], "--summary", "Nháp brief")
	if code != 0 || !strings.Contains(out, "| không git, ") {
		t.Fatalf("không git: %s (%d)", out, code)
	}
	if _, code = run(t, dir, "changelog", "pending"); code != 1 {
		t.Fatal("pending ngoài git phải mã 1")
	}
}
