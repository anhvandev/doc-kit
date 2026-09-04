package docs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const specBody = `
# F-001: Bộ lọc

## 2. Mục đích

## 3. Tác nhân và điều kiện tiên quyết

- Tác nhân: nhân viên kho đã đăng nhập
- Điều kiện tiên quyết: có ít nhất một đơn hàng
- Ghi chú:

## 5. Hành vi theo mã bước

| Mã | Hành động | Phản hồi |
|---|---|---|
| B1 | mở | hiện |

## 6. Giao diện

| Mã bước | Mockup | Trạng thái hiển thị |
|---|---|---|
| B1 | [B1](../design/mockups/F-001-B1.html) | normal |
| B2a | chưa có, xem họ Design | |

## 9. Tiêu chí chấp nhận

- AC1. **Given** đang ở danh sách **When** chọn trạng thái "đã giao" **Then** chỉ còn đơn đã giao
- AC2. Given không có đơn When lọc Then hiện "Không có đơn"
- AC3. **Given** có 3 đơn, **When** lọc theo ngày **Then** hiện 3 đơn
- AC4. Kết quả đúng là được
- Ghi chú không phải AC
5. AC5. **Given** a **When** b **Then** c
| AC6 | bảng | x |
- Đã bật cờ tính năng: beta

## 10. Dữ liệu
`

func TestExtractSpec(t *testing.T) {
	ex := ExtractSpec([]byte(specBody))
	if len(ex.Scenarios) != 6 {
		t.Fatalf("muốn 6 scenario, có %d: %+v", len(ex.Scenarios), ex.Scenarios)
	}
	// Danh sách đánh số và dòng bảng không bị bỏ im lặng: giữ Raw.
	if s := ex.Scenarios[4]; s.Code != "AC5" || s.Raw != "5. AC5. **Given** a **When** b **Then** c" {
		t.Fatalf("AC5 đánh số phải giữ Raw: %+v", s)
	}
	if s := ex.Scenarios[5]; s.Code != "AC6" || !strings.HasPrefix(s.Raw, "AC6 | bảng") {
		t.Fatalf("AC6 bảng phải giữ Raw: %+v", s)
	}
	s := ex.Scenarios[0]
	if s.Code != "AC1" || s.Given != "đang ở danh sách" || s.When != `chọn trạng thái "đã giao"` || s.Then != "chỉ còn đơn đã giao" || s.Raw != "" {
		t.Fatalf("AC1 tách sai: %+v", s)
	}
	if s := ex.Scenarios[1]; s.Given != "không có đơn" || s.Then != `hiện "Không có đơn"` {
		t.Fatalf("AC2 không in đậm tách sai: %+v", s)
	}
	if s := ex.Scenarios[2]; s.Given != "có 3 đơn" {
		t.Fatalf("AC3 bỏ dấu phẩy cuối sai: %+v", s)
	}
	if s := ex.Scenarios[3]; s.Code != "AC4" || s.Raw != "Kết quả đúng là được" || s.Given != "" {
		t.Fatalf("AC4 lệch khung phải giữ Raw: %+v", s)
	}
	if len(ex.Background) != 2 || ex.Background[0] != "nhân viên kho đã đăng nhập" || ex.Background[1] != "có ít nhất một đơn hàng" {
		t.Fatalf("Background: %v", ex.Background)
	}
	if len(ex.Steps) != 2 || ex.Steps[0].Code != "B1" || ex.Steps[0].Mockup != "[B1](../design/mockups/F-001-B1.html)" || ex.Steps[1].Code != "B2a" {
		t.Fatalf("Steps: %+v", ex.Steps)
	}
}

func TestExtractSpecGherkin(t *testing.T) {
	body := "## 9. Tiêu chí\n\n```gherkin\nFeature: x\n\n  Scenario: AC1 lọc\n    Given có đơn\n    And đã đăng nhập\n    When lọc\n    Then thấy\n\n  Scenario: AC2 rỗng\n    Given không\n    When lọc\n    Then trống\n\n  Scenario Outline: AC3 nhiều\n    Given <n>\n    When x\n    Then y\n\n    Examples:\n      | n |\n      | 1 |\n```\n\n## 10. Dữ liệu\n"
	ex := ExtractSpec([]byte(body))
	if len(ex.Scenarios) != 3 || ex.Scenarios[0].Given != "có đơn và đã đăng nhập" || ex.Scenarios[0].Title != "lọc" || ex.Scenarios[1].Then != "trống" {
		t.Fatalf("gherkin: %+v", ex.Scenarios)
	}
	if s := ex.Scenarios[2]; s.Code != "AC3" || !strings.HasPrefix(s.Raw, "Scenario Outline: nhiều") || s.Given != "" {
		t.Fatalf("Scenario Outline phải giữ Raw: %+v", s)
	}
}

// test-case --from spec: một Scenario mỗi AC với tag @F-xxx @ACn; metadata # dk: đọc lại được.
func TestNewTestCaseFromSpec(t *testing.T) {
	reg, o := setup(t)
	spec := filepath.Join(o.DocsDir, "features", "F-001-bo-loc.md")
	os.MkdirAll(filepath.Dir(spec), 0o755)
	os.WriteFile(spec, []byte("---\nid: F-001\ntype: feature-spec\ntitle: Bộ lọc\nstatus: approved\n---\n"+specBody), 0o644)
	o.From = spec
	r, err := New(reg, "test-case", "bo-loc", o)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(r.Path) != "F-001-cases.feature" || filepath.Base(filepath.Dir(r.Path)) != "test" {
		t.Fatalf("đường dẫn: %s", r.Path)
	}
	b, _ := os.ReadFile(r.Path)
	out := string(b)
	if r.Scenarios != 3 || r.Unparsed != 3 {
		t.Fatalf("đếm: %d tách được, %d chưa", r.Scenarios, r.Unparsed)
	}
	if !strings.HasPrefix(out, "# dk:\n# type: test-case\n") || strings.Count(out, "\n  Scenario: AC") != 6 || !strings.Contains(out, "\n  Scenario: AC1\n") || strings.Count(out, "@F-001 @AC") != 6 {
		t.Fatalf("feature sinh sai:\n%s", out)
	}
	if !strings.Contains(out, "\nFeature: Bộ lọc\n") || !strings.Contains(out, "  Background:\n    Given nhân viên kho đã đăng nhập\n    And có ít nhất một đơn hàng\n") || !strings.Contains(out, "# chưa tách được: Kết quả đúng là được\n    Given TODO") {
		t.Fatalf("thiếu Feature, Background hoặc ghi chú AC lệch:\n%s", out)
	}
	fm := readFMFile(t, r.Path)
	if fm["source"] != "F-001" || fm["feature"] != "F-001" || fm["title"] != "Bộ lọc" || fm["owner"] != "an" {
		t.Fatalf("metadata: %v", fm)
	}
	// Bảng và checklist UI từ cùng spec.
	r, err = New(reg, "test-case-table", "bo-loc", o)
	if err != nil {
		t.Fatal(err)
	}
	b, _ = os.ReadFile(r.Path)
	if filepath.Base(r.Path) != "F-001-cases.md" || strings.Count(string(b), "\n| AC") != 6 || !strings.Contains(string(b), `chưa tách được: AC6 \| bảng \| x`) {
		t.Fatalf("bảng sinh sai:\n%s", b)
	}
	r, err = New(reg, "ui-test-checklist", "bo-loc", o)
	if err != nil {
		t.Fatal(err)
	}
	b, _ = os.ReadFile(r.Path)
	if filepath.Base(r.Path) != "F-001-ui.md" || !strings.Contains(string(b), "- [ ] B1: khớp mockup [B1](../design/mockups/F-001-B1.html)\n- [ ] B2a: khớp mockup chưa có, xem họ Design") {
		t.Fatalf("checklist sinh sai:\n%s", b)
	}
	// Không --from: một Scenario trống, cần --set feature.
	o.From = ""
	if _, err := New(reg, "test-case", "x", o); err == nil || !strings.Contains(err.Error(), "--set feature") {
		t.Fatalf("thiếu feature phải lỗi: %v", err)
	}
	o.Set = map[string]string{"feature": "F-009"}
	r, err = New(reg, "test-case", "x", o)
	if err != nil {
		t.Fatal(err)
	}
	b, _ = os.ReadFile(r.Path)
	if strings.Count(string(b), "Scenario:") != 1 || !strings.Contains(string(b), "@F-009 @AC1") {
		t.Fatalf("khung trống sai:\n%s", b)
	}
}

// plan tạo thư mục {yymmdd}-{hhmm}-{slug}/plan.md trong plans/; plan-phase đếm
// số trong thư mục --in; report vào <in>/reports/; thiếu --in là lỗi.
func TestNewPlanPhaseReport(t *testing.T) {
	reg, o := setup(t)
	o.PlansDir = filepath.Join(filepath.Dir(o.DocsDir), "plans")
	o.IDPrefix = "X"
	r, err := New(reg, "plan", "dot-mot", o)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(o.PlansDir, "260903-1405-dot-mot", "plan.md")
	if r.Path != want {
		t.Fatalf("plan: %s, muốn %s", r.Path, want)
	}
	if _, err := New(reg, "plan-phase", "khoi-tao", o); err == nil || !strings.Contains(err.Error(), "--in") {
		t.Fatalf("plan-phase thiếu --in phải lỗi: %v", err)
	}
	o.In = filepath.Join(t.TempDir(), "ngoai")
	os.MkdirAll(o.In, 0o755)
	if _, err := New(reg, "plan-phase", "x", o); err == nil || !strings.Contains(err.Error(), "phải nằm trong") {
		t.Fatalf("--in ngoài plans_dir phải lỗi: %v", err)
	}
	o.In = filepath.Join(o.PlansDir, "khong-co")
	if _, err := New(reg, "report", "x", o); err == nil || !strings.Contains(err.Error(), "không phải thư mục có sẵn") {
		t.Fatalf("--in thư mục chưa có phải lỗi: %v", err)
	}
	if _, err := New(reg, "adr", "x", Options{DocsDir: o.DocsDir, In: o.PlansDir, Now: o.Now}); err == nil || !strings.Contains(err.Error(), "không dùng --in") {
		t.Fatalf("--in với loại không dùng phải lỗi: %v", err)
	}
	o.In = filepath.Dir(want)
	for i, slug := range []string{"khoi-tao", "hoan-thien"} {
		r, err = New(reg, "plan-phase", slug, o)
		if err != nil {
			t.Fatal(err)
		}
		wantName := []string{"phase-01-khoi-tao.md", "phase-02-hoan-thien.md"}[i]
		if filepath.Base(r.Path) != wantName || filepath.Dir(r.Path) != o.In || r.ID != strings.TrimSuffix(wantName, "-"+slug+".md") {
			t.Fatalf("phase %d: %s id %s (id_prefix không áp cho phase)", i+1, r.Path, r.ID)
		}
	}
	r, err = New(reg, "report", "phase-01-run", o)
	if err != nil {
		t.Fatal(err)
	}
	if r.Path != filepath.Join(o.In, "reports", "report-260903-1405-phase-01-run.md") {
		t.Fatalf("report: %s", r.Path)
	}
	o.In = o.PlansDir
	r, _ = New(reg, "report", "tong-hop", o)
	if r.Path != filepath.Join(o.PlansDir, "reports", "report-260903-1405-tong-hop.md") {
		t.Fatalf("report ngoài plan: %s", r.Path)
	}
}

// decision-log --append nối dòng, giữ nguyên nội dung cũ, tạo file khi chưa có.
func TestDecisionLogAppend(t *testing.T) {
	reg, o := setup(t)
	o.Append = "Dùng nút màu chính cho Lưu | an | F-001"
	r, err := New(reg, "decision-log", "decision-log", o)
	if err != nil {
		t.Fatal(err)
	}
	if r.Appended {
		t.Fatal("lần đầu là tạo file, không phải nối")
	}
	b, _ := os.ReadFile(r.Path)
	first := string(b)
	if !strings.HasSuffix(first, "\n- 2026-09-03 | Dùng nút màu chính cho Lưu | an | F-001\n") || !strings.Contains(first, "type: decision-log") {
		t.Fatalf("append lần đầu:\n%s", first)
	}
	o.Append = "Bỏ cột Ghi chú | an | -"
	o.Now = o.Now.Add(24 * time.Hour)
	if r2, err := New(reg, "decision-log", "decision-log", o); err != nil || !r2.Appended {
		t.Fatalf("lần hai phải là nối: %+v %v", r2, err)
	}
	b, _ = os.ReadFile(r.Path)
	if !strings.HasPrefix(string(b), strings.Replace(first, "updated: 2026-09-03 14:05", "updated: 2026-09-04 14:05", 1)) || !strings.HasSuffix(string(b), "\n- 2026-09-04 | Bỏ cột Ghi chú | an | -\n") {
		t.Fatalf("append lần hai phá nội dung cũ:\n%s", b)
	}
	if _, err := New(reg, "adr", "x", Options{DocsDir: o.DocsDir, Append: "x", Now: o.Now}); err == nil || !strings.Contains(err.Error(), "decision-log") {
		t.Fatalf("--append loại khác phải lỗi: %v", err)
	}
}
