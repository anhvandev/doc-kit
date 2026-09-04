package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anhvandev/doc-kit/internal/check"
)

// initProject tạo dự án có git, một brief approved và một feature-spec từ brief.
func initProject(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	if _, code := run(t, dir, "init"); code != 0 {
		t.Fatal("init")
	}
	if _, code := run(t, dir, "new", "brief", "loc-don", "--set", "owner=v", "--set", "status=approved"); code != 0 {
		t.Fatal("new brief")
	}
	briefs, _ := filepath.Glob(filepath.Join(dir, "docs", "intake", "*", "brief.md"))
	if _, code := run(t, dir, "new", "feature-spec", "bo-loc", "--from", briefs[0], "--set", "owner=v"); code != 0 {
		t.Fatal("new feature-spec")
	}
	return dir
}

func TestRenderAllIndex(t *testing.T) {
	dir := initProject(t)
	out, code := run(t, dir, "render", "--all", "--index", "--json")
	var res []renderResult
	if code != 0 || json.Unmarshal([]byte(out), &res) != nil {
		t.Fatalf("render: %s (%d)", out, code)
	}
	var spec, brief renderResult
	for _, r := range res {
		switch {
		case strings.HasPrefix(r.Src, "docs/features/"):
			spec = r
		case strings.HasSuffix(r.Src, "brief.md"):
			brief = r
		}
	}
	if !spec.Mermaid || brief.Mermaid || spec.Out != "docs/html/features/F-001-bo-loc.html" {
		t.Fatalf("kết quả render sai: %+v", res)
	}
	specHTML, _ := os.ReadFile(filepath.Join(dir, spec.Out))
	briefHTML, _ := os.ReadFile(filepath.Join(dir, brief.Out))
	if !strings.Contains(string(specHTML), "mermaid.initialize") || strings.Contains(string(briefHTML), "mermaid.initialize") {
		t.Fatal("chỉ trang có sơ đồ mới nhúng mermaid")
	}
	if strings.Contains(string(specHTML), "<link") || strings.Contains(string(specHTML), "src=\"http") {
		t.Fatal("HTML phải tự chứa, không tài nguyên ngoài")
	}
	if _, err := os.Stat(filepath.Join(dir, "docs", "html", "index.html")); err != nil {
		t.Fatal("thiếu index.html")
	}
	// Render một file, đường dẫn tương đối cwd.
	if out, code := run(t, dir, "render", "docs/features/F-001-bo-loc.md"); code != 0 || !strings.Contains(out, "docs/html/features/F-001-bo-loc.html") {
		t.Fatalf("render một file: %s (%d)", out, code)
	}
	if _, code := run(t, dir, "render"); code != 2 {
		t.Fatal("render không tham số phải mã 2")
	}
	if _, code := run(t, dir, "render", "dk.toml"); code != 2 {
		t.Fatal("render file không phải .md phải mã 2")
	}
	if _, code := run(t, dir, "render", "README.md"); code != 1 {
		t.Fatal("render file .md ngoài docs/ phải mã 1")
	}
}

func TestIndexAndChangelogSkipsGenerated(t *testing.T) {
	dir := initProject(t)
	out, code := run(t, dir, "index", "all", "--json")
	var files []string
	if code != 0 || json.Unmarshal([]byte(out), &files) != nil || len(files) != 5 {
		t.Fatalf("index all: %s (%d)", out, code)
	}
	for _, f := range files {
		b, err := os.ReadFile(filepath.Join(dir, f))
		if err != nil || !strings.HasPrefix(string(b), "---\ngenerated: true\n") {
			t.Fatalf("%s thiếu generated: true", f)
		}
	}
	feat, _ := os.ReadFile(filepath.Join(dir, "docs", "features", "README.md"))
	if !strings.Contains(string(feat), "| [F-001](F-001-bo-loc.md) | Loc don | draft | v |") || strings.Contains(string(feat), "## draft") {
		t.Fatalf("chỉ mục features sai:\n%s", feat)
	}
	// intake: một dòng một thư mục, theo trạng thái brief; thư mục chưa có brief xếp trước.
	if _, code := run(t, dir, "new", "idea", "y-tuong-moi"); code != 0 {
		t.Fatal("new idea")
	}
	if _, code := run(t, dir, "index", "intake"); code != 0 {
		t.Fatal("index intake")
	}
	in, _ := os.ReadFile(filepath.Join(dir, "docs", "intake", "README.md"))
	s := string(in)
	if strings.Count(s, "\n| [") != 2 || !strings.Contains(s, "## chưa có brief") || !strings.Contains(s, "## approved") ||
		strings.Index(s, "## chưa có brief") > strings.Index(s, "## approved") || !strings.Contains(s, "y-tuong-moi/idea.md") {
		t.Fatalf("chỉ mục intake sai:\n%s", s)
	}
	// cr: interview trong thư mục CR không phải một dòng CR.
	if _, code := run(t, dir, "new", "cr", "doi-loc"); code != 0 {
		t.Fatal("new cr")
	}
	crs, _ := filepath.Glob(filepath.Join(dir, "docs", "cr", "CR-*.md"))
	if _, code := run(t, dir, "new", "interview", "doi-loc", "--from", crs[0]); code != 0 {
		t.Fatal("new interview --from cr")
	}
	if _, code := run(t, dir, "index", "cr"); code != 0 {
		t.Fatal("index cr")
	}
	cr, _ := os.ReadFile(filepath.Join(dir, "docs", "cr", "README.md"))
	if strings.Contains(string(cr), "interview") || strings.Count(string(cr), "\n| [") != 1 {
		t.Fatalf("chỉ mục cr sai:\n%s", cr)
	}
	out, _ = run(t, dir, "changelog", "pending", "--json")
	if strings.Contains(out, "README.md") {
		t.Fatalf("changelog pending không được liệt kê chỉ mục sinh ra:\n%s", out)
	}
	if _, code := run(t, dir, "index", "bogus"); code != 2 {
		t.Fatal("index tham số lạ phải mã 2")
	}
}

func TestCheckRefsStatus(t *testing.T) {
	dir := initProject(t)
	// Dự án mới: brief có liên kết ./interview.md chưa tồn tại là lỗi thật của template.
	out, code := run(t, dir, "check", "--json")
	var fs []check.Finding
	if code != 3 || json.Unmarshal([]byte(out), &fs) != nil {
		t.Fatalf("check: %s (%d)", out, code)
	}
	// Tạo interview thì hết lỗi.
	briefs, _ := filepath.Glob(filepath.Join(dir, "docs", "intake", "*", "brief.md"))
	if _, code := run(t, dir, "new", "interview", "loc-don", "--from", briefs[0]); code != 0 {
		t.Fatal("new interview")
	}
	if out, code := run(t, dir, "check"); code != 0 {
		t.Fatalf("check sau khi có interview: %s (%d)", out, code)
	}
	// Làm lệch mã bước và vượt ngưỡng dòng.
	spec := filepath.Join(dir, "docs", "features", "F-001-bo-loc.md")
	b, _ := os.ReadFile(spec)
	s := strings.Replace(string(b), "| B5 | | |", "| B7 | | |", 1) + strings.Repeat("x\n", 900)
	_ = os.WriteFile(spec, []byte(s), 0o644)
	out, code = run(t, dir, "check", "docs/features/F-001-bo-loc.md", "--json")
	if code != 3 || json.Unmarshal([]byte(out), &fs) != nil {
		t.Fatalf("check file: %s (%d)", out, code)
	}
	rules := map[string]bool{}
	for _, f := range fs {
		if f.File != "docs/features/F-001-bo-loc.md" {
			t.Fatalf("check <file> chỉ lọc file đó, thấy %s", f.File)
		}
		rules[f.Rule+":"+f.Msg] = true
	}
	for _, want := range []string{"step-codes:mã có trong sơ đồ nhưng thiếu ở bảng hành vi: B5", "step-codes:mã có trong bảng hành vi nhưng không có trong sơ đồ: B7"} {
		if !rules[want] {
			t.Fatalf("thiếu %q trong %v", want, rules)
		}
	}
	found := false
	for k := range rules {
		if strings.HasPrefix(k, "line-threshold:") && strings.Contains(k, "vượt max_lines 800") {
			found = true
		}
	}
	if !found {
		t.Fatalf("thiếu line-threshold trong %v", rules)
	}

	if _, code := run(t, dir, "check", "docs/features/khong-co.md"); code != 1 {
		t.Fatal("check file không có trong docs/ phải mã 1, không được báo sạch")
	}
	// plans/ được quét liên kết nhưng không áp ngưỡng dòng.
	_ = os.WriteFile(filepath.Join(dir, "plans", "dai.md"), []byte("# plan\n"+strings.Repeat("x\n", 900)), 0o644)
	if out, code := run(t, dir, "check", "plans/dai.md"); code != 0 {
		t.Fatalf("plan dài không được lỗi ngưỡng dòng: %s (%d)", out, code)
	}

	out, code = run(t, dir, "refs", "docs/features/F-001-bo-loc.md", "--json")
	var r struct{ Out, In []string }
	if code != 0 || json.Unmarshal([]byte(out), &r) != nil {
		t.Fatalf("refs: %s (%d)", out, code)
	}
	if len(r.Out) != 1 || !strings.HasSuffix(r.Out[0], "/brief.md") || len(r.In) != 0 {
		t.Fatalf("refs sai: %+v", r)
	}
	out, code = run(t, dir, "refs", filepath.ToSlash(strings.TrimPrefix(briefs[0], dir+"/")))
	if code != 0 || !strings.Contains(out, "<- docs/features/F-001-bo-loc.md") {
		t.Fatalf("refs brief: %s (%d)", out, code)
	}
	if _, code := run(t, dir, "refs", "dk.toml"); code != 1 {
		t.Fatal("refs file ngoài docs/ phải mã 1")
	}

	out, code = run(t, dir, "status", "--json")
	var st statusReport
	if code != 0 || json.Unmarshal([]byte(out), &st) != nil {
		t.Fatalf("status: %s (%d)", out, code)
	}
	if st.Docs["feature-spec"]["draft"] != 1 || st.Docs["brief"]["approved"] != 1 || st.Docs["interview"]["open"] != 1 || st.OpenCR != 0 {
		t.Fatalf("status đếm sai: %+v", st)
	}
	if st.ChangelogPending == nil || *st.ChangelogPending != 3 || st.Check["error"] != 3 {
		t.Fatalf("status pending/check sai: pending=%v check=%v", st.ChangelogPending, st.Check)
	}
	if out, code := run(t, dir, "status"); code != 0 || !strings.Contains(out, "CR đang mở: 0") {
		t.Fatalf("status text: %s (%d)", out, code)
	}
}
