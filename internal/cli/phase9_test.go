package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Chuỗi lệnh của doc-release và doc-ops: release-brief --from spec, no-jargon,
// release-notes --collect ghi released_in và changelog, index user-guide,
// faq --append, environment với env-no-secret, postmortem, backup-dr trong status,
// Tầng 9 chỉ dk new.
func TestReleaseAndOps(t *testing.T) {
	dir := initProject(t)
	spec := filepath.Join(dir, "docs", "features", "F-001-bo-loc.md")
	b, _ := os.ReadFile(spec)
	s := strings.Replace(string(b), "status: draft", "status: implemented", 1)
	s = strings.Replace(s, "| B1 | | |", "| B1 | Mở màn hình endpoint | Hiển thị |", 1)
	s = strings.Replace(s, "## 11. Ngoài phạm vi\n\n<!-- gợi ý: thứ cố ý không làm trong tính năng này; lấy từ mục 3 của brief -->\n\n- ", "## 11. Ngoài phạm vi\n\n- Lọc theo khách", 1)
	os.WriteFile(spec, []byte(s), 0o644)

	out, code := run(t, dir, "new", "release-brief", "bo-loc", "--from", spec, "--set", "owner=v")
	if code != 0 || !strings.Contains(out, "docs/release/briefs/F-001.md") {
		t.Fatalf("release-brief: %s (%d)", out, code)
	}
	brief := filepath.Join(dir, "docs", "release", "briefs", "F-001.md")
	b, _ = os.ReadFile(brief)
	if !strings.Contains(string(b), "feature: F-001") || !strings.Contains(string(b), "source: F-001") || !strings.Contains(string(b), "1. Mở màn hình endpoint") || !strings.Contains(string(b), "- Lọc theo khách") {
		t.Fatalf("brief:\n%s", b)
	}
	if out, code := run(t, dir, "check", brief); code != 0 || !strings.Contains(out, "no-jargon") || !strings.Contains(out, `"endpoint"`) {
		t.Fatalf("brief chứa endpoint phải cảnh báo no-jargon: %s (%d)", out, code)
	}
	if _, code := run(t, dir, "new", "release-brief", "khac", "--set", "owner=v", "--set", "feature=F-002", "--set", "kind=fix", "--set", "status=ready"); code != 0 {
		t.Fatal("brief thứ hai")
	}
	// --collect: brief 1 chưa ready nên chỉ gom brief 2; sau khi ready cả hai thì gom 1 (brief 2 đã có released_in).
	if _, code := run(t, dir, "new", "release-notes", "--collect", "v1.0.0", "--set", "owner=v"); code != 0 {
		t.Fatal("collect lần 1")
	}
	os.WriteFile(brief, []byte(strings.Replace(string(b), "status: draft", "status: ready", 1)), 0o644)
	out, code = run(t, dir, "new", "release-notes", "--collect", "v1.1.0", "--set", "owner=v", "--json")
	var res struct {
		Path     string
		Released []string
	}
	if code != 0 || json.Unmarshal([]byte(out), &res) != nil || res.Path != "docs/release/v1.1.0.md" || len(res.Released) != 1 || res.Released[0] != "docs/release/briefs/F-001.md" {
		t.Fatalf("collect lần 2: %s (%d)", out, code)
	}
	notes, _ := os.ReadFile(filepath.Join(dir, res.Path))
	if !strings.Contains(string(notes), "- [Loc don](briefs/F-001.md)") || strings.Contains(string(notes), "F-002") {
		t.Fatalf("notes:\n%s", notes)
	}
	b, _ = os.ReadFile(brief)
	if !strings.Contains(string(b), "released_in: v1.1.0") {
		t.Fatalf("released_in chưa ghi:\n%s", b)
	}
	log, _ := os.ReadFile(filepath.Join(dir, "docs", "CHANGELOG-DOCS.md"))
	if !strings.Contains(string(log), "release/briefs/F-001.md | ") || !strings.Contains(string(log), "Phát hành trong v1.1.0 | v1.1.0") || !strings.Contains(string(log), "Phát hành trong v1.0.0 | v1.0.0") {
		t.Fatalf("changelog:\n%s", log)
	}
	if _, code := run(t, dir, "new", "release-notes", "--collect", "v1.2.0"); code != 1 {
		t.Fatalf("collect rỗng phải mã 1, được %d", code)
	}
	if _, code := run(t, dir, "new", "adr", "x", "--collect", "v1"); code != 1 {
		t.Fatalf("--collect loại khác phải mã 1, được %d", code)
	}

	// User guide theo nhiệm vụ và mục lục.
	if _, code := run(t, dir, "new", "user-guide", "loc-don", "--from", brief, "--set", "owner=v", "--set", "task=Bán hàng"); code != 0 {
		t.Fatal("user-guide")
	}
	if _, code := run(t, dir, "new", "user-guide", "xuat-file", "--set", "owner=v"); code != 0 {
		t.Fatal("user-guide 2")
	}
	if out, code := run(t, dir, "index", "user-guide"); code != 0 || !strings.Contains(out, "docs/release/guide/README.md") {
		t.Fatalf("index user-guide: %s (%d)", out, code)
	}
	idx, _ := os.ReadFile(filepath.Join(dir, "docs", "release", "guide", "README.md"))
	if !strings.Contains(string(idx), "## Bán hàng\n\n| Trang | Trạng thái | Cập nhật |\n|---|---|---|\n| [Loc don](loc-don.md) | draft |") || !strings.Contains(string(idx), "## Chưa phân nhóm") || strings.Index(string(idx), "Bán hàng") > strings.Index(string(idx), "Chưa phân nhóm") {
		t.Fatalf("mục lục guide:\n%s", idx)
	}

	// FAQ nối dòng.
	if out, code := run(t, dir, "new", "faq", "--append", "Quên mật khẩu? | Dùng nút Quên mật khẩu | guide/loc-don.md", "--set", "owner=v"); code != 0 || !strings.Contains(out, "và nối dòng") {
		t.Fatalf("faq lần đầu: %s (%d)", out, code)
	}
	if out, code := run(t, dir, "new", "faq", "--append", "Câu hai | Trả lời | -"); code != 0 || !strings.Contains(out, "Đã nối dòng") {
		t.Fatalf("faq lần hai: %s (%d)", out, code)
	}
	faq, _ := os.ReadFile(filepath.Join(dir, "docs", "release", "faq.md"))
	if strings.Count(string(faq), "\n- 20") != 2 {
		t.Fatalf("faq:\n%s", faq)
	}

	// Environment: secret thật bị chặn, placeholder qua.
	if _, code := run(t, dir, "new", "environment", "environment", "--set", "owner=v"); code != 0 {
		t.Fatal("environment")
	}
	env := filepath.Join(dir, "docs", "ops", "environment.md")
	b, _ = os.ReadFile(env)
	os.WriteFile(env, []byte(strings.Replace(string(b), "```\n```", "```\nDB_PASSWORD=abc\n```", 1)), 0o644)
	if out, code := run(t, dir, "check", env); code != 3 || !strings.Contains(out, "error env-no-secret: DB_PASSWORD") {
		t.Fatalf("secret thật phải lỗi mã 3: %s (%d)", out, code)
	}
	os.WriteFile(env, []byte(strings.Replace(string(b), "```\n```", "```\nDB_PASSWORD=<secret>\n```", 1)), 0o644)
	if out, code := run(t, dir, "check", env); code != 0 || strings.Contains(out, "env-no-secret") {
		t.Fatalf("placeholder phải qua: %s (%d)", out, code)
	}

	// Postmortem quá 48 giờ, backup-dr chưa diễn tập trong status, Tầng 9.
	out, code = run(t, dir, "new", "postmortem", "mat-db", "--set", "owner=v", "--set", "incident_at=2026-01-01 10:00")
	if code != 0 {
		t.Fatalf("postmortem: %s", out)
	}
	pms, _ := filepath.Glob(filepath.Join(dir, "docs", "ops", "postmortems", "*-mat-db.md"))
	b, _ = os.ReadFile(pms[0])
	if !strings.Contains(string(b), "written_within_48h: false") || !strings.Contains(string(b), "incident_at: 2026-01-01 10:00") {
		t.Fatalf("postmortem:\n%s", b)
	}
	if _, code := run(t, dir, "new", "backup-dr", "backup-dr", "--set", "owner=v"); code != 0 {
		t.Fatal("backup-dr")
	}
	out, code = run(t, dir, "status", "--json")
	var st struct {
		DROverdue []string `json:"dr_overdue"`
	}
	if code != 0 || json.Unmarshal([]byte(out), &st) != nil || len(st.DROverdue) != 1 || st.DROverdue[0] != "docs/ops/backup-dr.md" {
		t.Fatalf("status dr_overdue: %s (%d)", out, code)
	}
	if out, code := run(t, dir, "status"); code != 0 || !strings.Contains(out, "DR chưa diễn tập quá 6 tháng: docs/ops/backup-dr.md") {
		t.Fatalf("status text: %s", out)
	}
	dr := filepath.Join(dir, "docs", "ops", "backup-dr.md")
	b, _ = os.ReadFile(dr)
	os.WriteFile(dr, []byte(strings.Replace(string(b), `last_drill: ""`, "last_drill: 15/08/2026", 1)), 0o644)
	if out, code := run(t, dir, "status"); code != 0 || !strings.Contains(out, "DR last_drill sai định dạng (cần yyyy-mm-dd): docs/ops/backup-dr.md") {
		t.Fatalf("status sai định dạng: %s", out)
	}
	// Các loại còn lại của Tầng 8 và 9 tạo từ template và qua check sạch.
	for _, c := range [][]string{
		{"runbook", "db-day", "docs/ops/runbooks/db-day.md"},
		{"deployment", "deployment", "docs/ops/deployment.md"},
		{"monitoring", "monitoring", "docs/ops/monitoring.md"},
		{"charter", "charter", "docs/governance/charter.md"},
		{"risk-register", "risk-register", "docs/governance/risk-register.md"},
		{"meeting-notes", "kickoff", ""},
	} {
		out, code := run(t, dir, "new", c[0], c[1], "--set", "owner=v")
		if code != 0 || (c[2] != "" && !strings.Contains(out, c[2])) {
			t.Fatalf("new %s: %s (%d)", c[0], out, code)
		}
		file := strings.TrimPrefix(strings.TrimSpace(out), "Đã tạo ")
		if out, code := run(t, dir, "check", filepath.Join(dir, file)); code != 0 || strings.Contains(out, "error") {
			t.Fatalf("check %s: %s (%d)", file, out, code)
		}
	}
	if m, _ := filepath.Glob(filepath.Join(dir, "docs", "governance", "meetings", "*-kickoff.md")); len(m) != 1 {
		t.Fatal("meeting-notes chưa tạo theo {yymmdd}-{slug}")
	}
}
