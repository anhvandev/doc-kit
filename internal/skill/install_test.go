package skill

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anhvandev/doc-kit/internal/target"
)

func TestHash(t *testing.T) {
	src, err := Files("doc-cr")
	if err != nil {
		t.Fatal(err)
	}
	h := Hash(src)
	withMeta := map[string][]byte{}
	for k, v := range src {
		withMeta[k] = v
	}
	md, err := withTrace(src["SKILL.md"], "1.2.3", h)
	if err != nil {
		t.Fatal(err)
	}
	withMeta["SKILL.md"] = md
	if Hash(withMeta) != h {
		t.Fatal("thêm khối metadata đổi hash")
	}
	if tr := readTrace(withMeta); tr.By != "dk" || tr.Version != "1.2.3" || tr.Hash != h {
		t.Fatalf("dấu vết sai: %+v", tr)
	}
	changed := map[string][]byte{}
	for k, v := range src {
		changed[k] = v
	}
	changed["references/rules.md"] = append([]byte("x"), src["references/rules.md"]...)
	if Hash(changed) == h {
		t.Fatal("đổi rules.md không đổi hash")
	}
}

func newTarget(t *testing.T) (target.Target, string) {
	t.Helper()
	root := t.TempDir()
	t.Setenv("HOME", filepath.Join(root, "home"))
	t.Setenv("USERPROFILE", filepath.Join(root, "home"))
	tg, err := target.Get("claude", root)
	if err != nil {
		t.Fatal(err)
	}
	return tg, root
}

func TestInstallLifecycle(t *testing.T) {
	tg, root := newTarget(t)
	dest := filepath.Join(root, ".claude", "skills", "doc-cr")
	res, err := Install(tg, nil, false, false, "0.1.0")
	if err != nil || len(res) != 12 || resOf(res, "doc-cr").Action != "đã cài" || res[11].Action != "đã cài" {
		t.Fatalf("cài mới: %+v %v", res, err)
	}
	b, _ := os.ReadFile(filepath.Join(dest, "SKILL.md"))
	if !strings.Contains(string(b), "dk_installed_by: dk") || !strings.Contains(string(b), "dk_hash:") {
		t.Fatalf("thiếu dấu vết:\n%s", b)
	}
	if _, err := os.Stat(filepath.Join(dest, "references", "rules.md")); err != nil {
		t.Fatal("thiếu references/rules.md")
	}
	rows, _ := Status(tg, false, "0.1.0")
	if rowOf(rows, "doc-cr").State != StateCurrent {
		t.Fatalf("trạng thái: %+v", rows)
	}

	res, err = Install(tg, []string{"doc-cr"}, false, false, "0.1.0")
	if err != nil || resOf(res, "doc-cr").Action != "không đổi" {
		t.Fatalf("cài lại: %+v %v", res, err)
	}
	rows, _ = Status(tg, false, "0.2.0")
	if rowOf(rows, "doc-cr").State != "cũ (v0.1.0)" {
		t.Fatalf("phiên bản cũ: %+v", rows)
	}
	res, err = Install(tg, nil, false, false, "0.2.0")
	if err != nil || resOf(res, "doc-cr").Action != "đã cập nhật" {
		t.Fatalf("nâng cấp: %+v %v", res, err)
	}

	rules := filepath.Join(dest, "references", "rules.md")
	os.WriteFile(rules, []byte("sửa tay\n"), 0o644)
	rows, _ = Status(tg, false, "0.2.0")
	if rowOf(rows, "doc-cr").State != StateModified {
		t.Fatalf("đã sửa tay: %+v", rows)
	}
	if _, err := Install(tg, nil, false, false, "0.2.0"); err == nil || !strings.Contains(err.Error(), "đã sửa tay") {
		t.Fatalf("phải từ chối skill sửa tay, được %v", err)
	}
	res, err = Install(tg, nil, false, true, "0.2.0")
	if err != nil || resOf(res, "doc-cr").Action != "đã cập nhật" {
		t.Fatalf("--force: %+v %v", res, err)
	}
	if b, _ := os.ReadFile(rules); string(b) == "sửa tay\n" {
		t.Fatal("--force không ghi đè")
	}

	res, err = Uninstall(tg, nil, false)
	if err != nil || resOf(res, "doc-cr").Action != "đã gỡ" {
		t.Fatalf("gỡ: %+v %v", res, err)
	}
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		t.Fatal("chưa xóa thư mục")
	}
	res, _ = Uninstall(tg, nil, false)
	if resOf(res, "doc-cr").Action != "bỏ qua" || res[0].Note != "chưa cài" {
		t.Fatalf("gỡ lần hai: %+v", res)
	}
}

func TestForeignSkillUntouched(t *testing.T) {
	tg, root := newTarget(t)
	dest := filepath.Join(root, ".claude", "skills", "doc-cr")
	os.MkdirAll(dest, 0o755)
	own := []byte("---\nname: doc-cr\ndescription: của tôi\n---\nnội dung riêng\n")
	os.WriteFile(filepath.Join(dest, "SKILL.md"), own, 0o644)
	if _, err := Install(tg, nil, false, true, "0.1.0"); err == nil || !strings.Contains(err.Error(), "không do dk") {
		t.Fatalf("phải từ chối, được %v", err)
	}
	rows, _ := Status(tg, false, "0.1.0")
	if rowOf(rows, "doc-cr").State != StateForeign {
		t.Fatalf("trạng thái: %+v", rows)
	}
	res, err := Uninstall(tg, nil, false)
	if err != nil || resOf(res, "doc-cr").Action != "bỏ qua" {
		t.Fatalf("gỡ: %+v %v", res, err)
	}
	if b, _ := os.ReadFile(filepath.Join(dest, "SKILL.md")); string(b) != string(own) {
		t.Fatal("đụng skill người dùng")
	}
}

func TestGlobalScope(t *testing.T) {
	tg, root := newTarget(t)
	if _, err := Install(tg, nil, true, false, "0.1.0"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "home", ".claude", "skills", "doc-cr", "SKILL.md")); err != nil {
		t.Fatal("không cài vào $HOME/.claude/skills")
	}
	if _, err := os.Stat(filepath.Join(root, ".claude")); !os.IsNotExist(err) {
		t.Fatal("--global đụng scope dự án")
	}
	if _, err := Install(tg, []string{"khong-co"}, true, false, "0.1.0"); err == nil {
		t.Fatal("skill không có phải lỗi")
	}
	for _, bad := range []string{"../x", "a/b", ".", ""} {
		if _, err := Uninstall(tg, []string{bad}, true); err == nil {
			t.Fatalf("tên %q phải bị từ chối", bad)
		}
	}
}

func TestInterruptedUpdateRecoverable(t *testing.T) {
	tg, root := newTarget(t)
	dest := filepath.Join(root, ".claude", "skills", "doc-cr")
	if _, err := Install(tg, nil, false, false, "0.1.0"); err != nil {
		t.Fatal(err)
	}
	// Giả lập đứt giữa chừng: thư mục tạm còn sót phải không đụng đích.
	os.MkdirAll(filepath.Join(root, ".claude", "skills", ".doc-cr-dk-x"), 0o755)
	res, err := Install(tg, nil, false, false, "0.1.0")
	if err != nil || resOf(res, "doc-cr").Action != "không đổi" {
		t.Fatalf("%+v %v", res, err)
	}
	if _, err := os.Stat(filepath.Join(dest, "SKILL.md")); err != nil {
		t.Fatal("đích bị đụng")
	}
}

func rowOf(rows []Row, name string) Row {
	for _, r := range rows {
		if r.Name == name {
			return r
		}
	}
	return Row{}
}

func resOf(res []Result, name string) Result {
	for _, r := range res {
		if r.Name == name {
			return r
		}
	}
	return Result{}
}
