package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Chuỗi lệnh của doc-plan-report và doc-test: plan, phase --in, report --in
// với report-evidence, decision-log --append, test-case --from spec làm
// spec-has-test hết cảnh báo.
func TestPlanReportAndTestCase(t *testing.T) {
	dir := initProject(t)
	out, code := run(t, dir, "new", "plan", "dot-mot", "--set", "owner=v")
	if code != 0 || !strings.Contains(out, "plans/") || !strings.HasSuffix(strings.TrimSpace(out), "-dot-mot/plan.md") {
		t.Fatalf("new plan: %s (%d)", out, code)
	}
	plans, _ := filepath.Glob(filepath.Join(dir, "plans", "*-dot-mot"))
	if _, code := run(t, dir, "new", "plan-phase", "x"); code != 1 {
		t.Fatalf("plan-phase thiếu --in phải mã 1, được %d", code)
	}
	rel := "plans/" + filepath.Base(plans[0])
	if out, code := run(t, dir, "new", "plan-phase", "khoi-tao", "--in", rel); code != 0 || !strings.Contains(out, "phase-01-khoi-tao.md (id phase-01)") {
		t.Fatalf("plan-phase: %s (%d)", out, code)
	}
	if out, code := run(t, dir, "new", "report", "phase-01-run", "--in", rel, "--set", "owner=v"); code != 0 || !strings.Contains(out, rel+"/reports/report-") {
		t.Fatalf("report: %s (%d)", out, code)
	}
	reports, _ := filepath.Glob(filepath.Join(plans[0], "reports", "*.md"))
	if out, code := run(t, dir, "check", reports[0]); code != 0 || !strings.Contains(out, "report-evidence") {
		t.Fatalf("report trống phải cảnh báo report-evidence: %s (%d)", out, code)
	}
	b, _ := os.ReadFile(reports[0])
	os.WriteFile(reports[0], []byte(strings.Replace(string(b), "- Commit: ", "- Commit: commit deadbeef1", 1)), 0o644)
	if out, code := run(t, dir, "check", reports[0]); code != 0 || strings.Contains(out, "report-evidence") {
		t.Fatalf("report có commit không được cảnh báo: %s (%d)", out, code)
	}
	if _, code := run(t, dir, "changelog", "add", reports[0], "--summary", "x"); code == 0 {
		t.Fatal("changelog add cho file trong plans/ phải bị từ chối")
	}

	if out, code := run(t, dir, "new", "decision-log", "--append", "Nút Lưu | v | F-001", "--set", "owner=v"); code != 0 || !strings.Contains(out, "và nối dòng") {
		t.Fatalf("decision-log lần đầu: %s (%d)", out, code)
	}
	if out, code := run(t, dir, "new", "decision-log", "--append", "Bỏ cột | v | -"); code != 0 || !strings.Contains(out, "Đã nối dòng") {
		t.Fatalf("decision-log lần hai: %s (%d)", out, code)
	}
	if _, code := run(t, dir, "new", "decision-log"); code != 2 {
		t.Fatalf("thiếu slug không --append phải mã 2, được %d", code)
	}
	log, _ := os.ReadFile(filepath.Join(dir, "docs", "plan", "decision-log.md"))
	if strings.Count(string(log), "\n- 20") != 2 {
		t.Fatalf("decision-log:\n%s", log)
	}

	spec := filepath.Join(dir, "docs", "features", "F-001-bo-loc.md")
	b, _ = os.ReadFile(spec)
	s := strings.Replace(string(b), "status: draft", "status: approved", 1)
	s = strings.Replace(s, "- AC1. **Given** ... **When** ... **Then** ...", "- AC1. **Given** a **When** b **Then** c\n- AC2. lệch khung", 1)
	os.WriteFile(spec, []byte(s), 0o644)
	if out, code := run(t, dir, "check", spec); code != 0 || !strings.Contains(out, "spec-has-test") {
		t.Fatalf("spec approved chưa có test phải cảnh báo: %s (%d)", out, code)
	}
	out, code = run(t, dir, "new", "test-case", "bo-loc", "--from", spec, "--set", "owner=v")
	if code != 0 || !strings.Contains(out, "docs/test/F-001-cases.feature: 1 Scenario, 1 dòng AC chưa tách được") {
		t.Fatalf("test-case: %s (%d)", out, code)
	}
	if out, code := run(t, dir, "check", spec); code != 0 || strings.Contains(out, "spec-has-test") {
		t.Fatalf("spec đã có .feature không được cảnh báo: %s (%d)", out, code)
	}
	if out, code := run(t, dir, "changelog", "add", "docs/test/F-001-cases.feature", "--summary", "test"); code != 0 || !strings.Contains(out, "test/F-001-cases.feature | mới,") {
		t.Fatalf("changelog add .feature: %s (%d)", out, code)
	}
	if out, code := run(t, dir, "check", "docs/test/F-001-cases.feature"); code != 0 {
		t.Fatalf("check .feature: %s (%d)", out, code)
	}
}
