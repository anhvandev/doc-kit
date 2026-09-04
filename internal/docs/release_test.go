package docs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/anhvandev/doc-kit/assets"
	"github.com/anhvandev/doc-kit/internal/doctype"
	"github.com/anhvandev/doc-kit/internal/frontmatter"
)

const releaseSpecBody = `# F-001: Bộ lọc

## 2. Mục đích và giá trị

<!-- gợi ý: bỏ -->

Người dùng lọc đơn theo trạng thái.

## 3. Tác nhân và điều kiện tiên quyết

- Tác nhân: nhân viên bán hàng
- Điều kiện tiên quyết: đã đăng nhập

## 4. Sơ đồ luồng chính

## 5. Hành vi theo mã bước

| Mã | Hành động của tác nhân | Phản hồi quan sát được của hệ thống |
|---|---|---|
| B1 | Mở danh sách đơn | Hiển thị bộ lọc |
| B2 | Chọn trạng thái | Danh sách cập nhật |

## 6. Giao diện

| Mã bước | Mockup | Trạng thái hiển thị |
|---|---|---|
| B1 | [B1](../design/mockups/F-001-B1.html) | normal |

## 11. Ngoài phạm vi

- Lọc theo khách hàng
`

func TestExtractRelease(t *testing.T) {
	ex := ExtractRelease([]byte(releaseSpecBody), "/p/docs/features", "/p/docs/release/briefs")
	if strings.Join(ex.Purpose, "|") != "Người dùng lọc đơn theo trạng thái." ||
		strings.Join(ex.Actors, "|") != "nhân viên bán hàng|đã đăng nhập" ||
		strings.Join(ex.Actions, "|") != "Mở danh sách đơn|Chọn trạng thái" ||
		strings.Join(ex.Limits, "|") != "Lọc theo khách hàng" {
		t.Fatalf("extract: %+v", ex)
	}
	if len(ex.Screens) != 1 || ex.Screens[0].Mockup != "[B1](../../design/mockups/F-001-B1.html)" {
		t.Fatalf("liên kết mockup phải đổi gốc sang release/briefs: %+v", ex.Screens)
	}
	// Ô mockup không phải liên kết (chữ mẫu của template) bị bỏ.
	if ex := ExtractRelease([]byte("## 6. Giao diện\n\n| Mã | Mockup | T |\n|---|---|---|\n| B1 | chưa có, xem họ Design | |\n"), "/a", "/b"); len(ex.Screens) != 0 {
		t.Fatalf("ô không liên kết phải bỏ: %+v", ex.Screens)
	}
	got := rebaseLinks("[a](../x.html#neo) [b](/abs/y.html) [c](https://e.com/z) [d](mailto:a@b.c)", "/p/docs/features", "/p/docs/release/briefs")
	if got != "[a](../../x.html#neo) [b](/abs/y.html) [c](https://e.com/z) [d](mailto:a@b.c)" {
		t.Fatalf("rebaseLinks: %s", got)
	}
}

// --collect gom brief ready chưa có released_in theo nhóm kind, ghi released_in
// vào từng brief; brief draft hoặc đã phát hành bị bỏ qua.
func TestCollectReleaseNotes(t *testing.T) {
	reg, _ := doctype.Load(assets.FS)
	root := t.TempDir()
	docsDir := filepath.Join(root, "docs")
	now := time.Date(2026, 9, 4, 9, 0, 0, 0, time.Local)
	mk := func(feature, kind, status, released string) string {
		os.MkdirAll(filepath.Join(docsDir, "release", "briefs"), 0o755)
		p := filepath.Join(docsDir, "release", "briefs", feature+".md")
		fm := "---\ntype: release-brief\ntitle: Tính năng " + feature + "\nstatus: " + status + "\nowner: v\ncreated: 2026-09-01\nupdated: 2026-09-01 10:00\nsource: " + feature + "\ncreated_by: dk\nfeature: " + feature + "\nkind: " + kind + "\nreleased_in: \"" + released + "\"\n---\n\n# x\n"
		os.WriteFile(p, []byte(fm), 0o644)
		return p
	}
	f1 := mk("F-001", "feature", "ready", "")
	f2 := mk("F-002", "fix", "ready", "")
	f3 := mk("F-003", "feature", "draft", "")
	f4 := mk("F-004", "feature", "ready", "v0.9.0")

	if _, err := New(reg, "release-brief", "x", Options{DocsDir: docsDir, Collect: "v1.0.0", Now: now}); err == nil || !strings.Contains(err.Error(), "--collect chỉ dùng cho release-notes") {
		t.Fatalf("--collect loại khác: %v", err)
	}
	res, err := New(reg, "release-notes", "v1-0-0", Options{DocsDir: docsDir, Collect: "v1.0.0", Now: now, Owner: "v"})
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(res.Path) != "v1.0.0.md" || len(res.Released) != 2 || res.Released[0] != f1 || res.Released[1] != f2 {
		t.Fatalf("kết quả: %+v", res)
	}
	b, _ := os.ReadFile(res.Path)
	s := string(b)
	if !strings.Contains(s, "version: \"v1.0.0\"") || !strings.Contains(s, "## 2. Mới\n\n- [Tính năng F-001](briefs/F-001.md)\n") ||
		!strings.Contains(s, "## 3. Sửa lỗi\n\n- [Tính năng F-002](briefs/F-002.md)\n") || strings.Contains(s, "F-003") || strings.Contains(s, "F-004") {
		t.Fatalf("release notes:\n%s", s)
	}
	for p, want := range map[string]string{f1: "v1.0.0", f2: "v1.0.0", f3: "", f4: "v0.9.0"} {
		b, _ := os.ReadFile(p)
		fm, _, _ := frontmatter.Split(b)
		if got := frontmatter.GetString(fm, "released_in"); got != want {
			t.Errorf("%s released_in %q, muốn %q", filepath.Base(p), got, want)
		}
	}
	// Lần hai: không còn brief nào ready chưa phát hành.
	if _, err := New(reg, "release-notes", "v1-1-0", Options{DocsDir: docsDir, Collect: "v1.1.0", Now: now}); err == nil || !strings.Contains(err.Error(), "không có Release brief nào") {
		t.Fatalf("collect rỗng phải lỗi: %v", err)
	}
	// Sinh lại cùng phiên bản với --force: giữ brief đã có released_in v1.0.0, gom thêm brief mới ready, chỉ ghi released_in cho brief mới.
	f5 := mk("F-005", "feature", "ready", "")
	res, err = New(reg, "release-notes", "v1-0-0", Options{DocsDir: docsDir, Collect: "v1.0.0", Now: now, Force: true})
	if err != nil || len(res.Released) != 1 || res.Released[0] != f5 {
		t.Fatalf("collect --force: %+v %v", res, err)
	}
	b, _ = os.ReadFile(res.Path)
	if !strings.Contains(string(b), "F-001.md") || !strings.Contains(string(b), "F-002.md") || !strings.Contains(string(b), "F-005.md") || strings.Contains(string(b), "F-004") {
		t.Fatalf("notes --force:\n%s", b)
	}
	if _, err := New(reg, "release-notes", "v1-0-0", Options{DocsDir: docsDir, Collect: "v1.0.0", Now: now}); err == nil || !strings.Contains(err.Error(), "đã tồn tại") {
		t.Fatalf("collect trùng file không --force phải lỗi: %v", err)
	}
	if VersionSlug("v1.0.0-RC.1") != "v1-0-0-rc-1" {
		t.Fatal("VersionSlug")
	}
}

func TestPostmortemWithin48h(t *testing.T) {
	reg, _ := doctype.Load(assets.FS)
	docsDir := filepath.Join(t.TempDir(), "docs")
	now := time.Date(2026, 9, 4, 9, 0, 0, 0, time.Local)
	for _, c := range []struct {
		at   string
		want string
	}{{"2026-09-03 10:00", "written_within_48h: true"}, {"2026-09-01 10:00", "written_within_48h: false"}, {"", "written_within_48h: false"}} {
		res, err := New(reg, "postmortem", "mat-db", Options{DocsDir: docsDir, Now: now, Force: true, Set: map[string]string{"incident_at": c.at}})
		if err != nil {
			t.Fatal(err)
		}
		b, _ := os.ReadFile(res.Path)
		if !strings.Contains(string(b), c.want) {
			t.Errorf("incident_at=%q: thiếu %s", c.at, c.want)
		}
	}
}
