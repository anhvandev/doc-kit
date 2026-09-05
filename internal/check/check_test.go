package check

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/anhvandev/doc-kit/assets"
	"github.com/anhvandev/doc-kit/internal/config"
	"github.com/anhvandev/doc-kit/internal/docs"
	"github.com/anhvandev/doc-kit/internal/doctype"
)

func runFixture(t *testing.T, strict bool) []Finding {
	t.Helper()
	root, _ := filepath.Abs("testdata")
	reg, err := doctype.Load(assets.FS)
	if err != nil {
		t.Fatal(err)
	}
	// Loại chưa có trong types.toml (test) vẫn được quét; check bỏ qua quy tắc theo loại.
	metas, err := docs.Scan(root, "docs", "plans")
	if err != nil {
		t.Fatal(err)
	}
	return Run(&Context{Root: root, DocsDir: "docs", Metas: metas, Reg: reg, Cfg: config.Check{WarnLines: 500, MaxLines: 800}, Strict: strict})
}

// has báo có finding khớp file, quy tắc, mức và chứa msg.
func has(fs []Finding, file, rule, level, msg string) bool {
	for _, f := range fs {
		if strings.HasSuffix(f.File, file) && f.Rule == rule && f.Level == level && strings.Contains(f.Msg, msg) {
			return true
		}
	}
	return false
}

func TestRules(t *testing.T) {
	fs := runFixture(t, false)
	cases := []struct{ file, rule, level, msg string }{
		{"features/F-002-b.md", "frontmatter-required", Error, "thiếu trường owner"},
		{"intake/260901-x/idea.md", "frontmatter-required", Warning, `created_by là "tay"`},
		{"features/F-002-b.md", "status-valid", Error, `status "banana"`},
		{"intake/260901-x/brief.md", "link-broken", Error, "./khong-co.md"},
		{"features/F-001-a.md", "link-broken", Error, "trỏ vào docs/html/"},
		{"features/F-001-a.md", "step-codes", Error, "thiếu ở bảng hành vi: B3"},
		{"features/F-001-a.md", "step-codes", Error, "không có trong sơ đồ: B4"},
		{"features/F-002-b.md", "backlink", Error, `source "CR-999999-khong-co"`},
		{"cr/CR-260901-chet.md", "backlink", Warning, "tài liệu chết"},
		{"features/F-001-a.md", "spec-has-test", Warning, "docs/test/"},
		{"features/F-003-dai.md", "line-threshold", Error, "vượt max_lines 800"},
		{"features/F-004-vua.md", "line-threshold", Warning, "vượt warn_lines 500"},
		{"overview/product-overview.md", "line-threshold", Warning, "160 dòng, vượt warn_lines 150"},
		{"features/F-002-b.md", "spec-section-order", Error, "thiếu mục 2, 3, 5, 6, 7, 9, 10"},
		{"features/F-006-crud.md", "spec-section-order", Error, "tiêu đề cấp 2 phải là 1, 2, 4, 5, 6, 7, 8, 9, 10 theo thứ tự; thiếu mục 8; mục lạ 11; mục lặp 7"},
		{"features/F-005-cr.md", "cr-approval-order", Error, "CR-260901-cho còn review (cập nhật 2026-09-01 10:00)"},
		{"features/F-002-b.md", "glossary-term", Warning, "thuật ngữ **kho** chưa có trong docs/overview/glossary-van-chuyen.md, docs/overview/glossary.md"},
		{"design/flows/F-001-flow.md", "userflow-steps", Error, "mã bước không có trong F-001: B9"},
		{"design/flows/F-404-flow.md", "userflow-steps", Error, `feature "F-404" không trỏ đến Feature Spec nào`},
		{"design/mockups/F-001-B1.html", "mockup-tokens", Error, `giá trị gõ tay "#333"`},
		{"design/mockups/F-001-B1.html", "mockup-tokens", Error, `giá trị gõ tay "12px"`},
		{"design/mockups/F-001-B1.html", "mockup-tokens", Error, `giá trị gõ tay "4px"`},
		{"design/mockups/F-001-B1.html", "mockup-tokens", Error, `giá trị gõ tay "#ABCDEF"`},
		{"design/mockups/khong-metadata.html", "mockup-tokens", Error, "thiếu khối <!-- dk:"},
		{"design/mockups/khong-metadata.html", "mockup-tokens", Error, `giá trị gõ tay "#444"`},    // style='' nháy đơn vẫn bị lint
		{"reports/report-260901-0900-thieu.md", "report-evidence", Warning, "không có bằng chứng"}, // liên kết .md không tính
		{"test/hong.feature", "frontmatter-required", Error, "khối `# dk:`"},
	}
	for _, c := range cases {
		if !has(fs, c.file, c.rule, c.level, c.msg) {
			t.Errorf("thiếu finding %s %s %s %q\ncó: %v", c.file, c.rule, c.level, c.msg, fs)
		}
	}
	// Không báo giả: brief được F-001 trỏ về; F-002 có test; README generated bỏ qua; link hợp lệ.
	for _, bad := range []struct{ file, rule string }{
		{"features/F-005-cr.md", "spec-section-order"},  // has_ui: false, bỏ mục 5
		{"features/F-005-cr.md", "step-codes"},          // B2a hậu tố khớp sơ đồ và bảng
		{"features/F-006-crud.md", "cr-approval-order"}, // spec cũ hơn CR
		{"features/F-002-b.md", "cr-approval-order"},    // source không phải CR tồn tại
		{"intake/260901-x/brief.md", "backlink"},
		{"features/F-002-b.md", "spec-has-test"},
		{"features/F-005-cr.md", "spec-has-test"}, // .feature với metadata # dk: source F-005
		{"reports/report-260901-0901-du.md", "report-evidence"},
		{"features/README.md", "link-broken"},
		{"features/F-001-a.md", "backlink"},
		{"cr/CR-260901-song.md", "backlink"},
		{"design/mockups/F-001-B2.html", "mockup-tokens"}, // external: Figma, không lint
		{"design/flows/F-001-flow.md", "step-codes"},      // userflow chỉ chịu userflow-steps
	} {
		for _, f := range fs {
			if strings.HasSuffix(f.File, bad.file) && f.Rule == bad.rule {
				t.Errorf("báo giả: %+v", f)
			}
		}
	}
	// Dòng của liên kết hỏng đúng: dòng 10 của brief.md (frontmatter 8 dòng + tiêu đề).
	for _, f := range fs {
		if f.Rule == "link-broken" && strings.HasSuffix(f.File, "brief.md") && f.Line != 10 {
			t.Errorf("dòng liên kết hỏng = %d, muốn 10", f.Line)
		}
	}
	// px trong prelude @media không thay bằng biến được nên không báo.
	for _, f := range fs {
		if f.Rule == "mockup-tokens" && strings.Contains(f.Msg, "768px") {
			t.Errorf("báo giả px trong @media: %+v", f)
		}
		if f.Rule == "line-threshold" && !strings.HasSuffix(f.File, ".md") {
			t.Errorf("line-threshold chỉ áp cho Markdown: %+v", f)
		}
	}
	// mockup-tokens đúng dòng: #333 ở dòng 17 (chú thích 12 dòng + 5), style="" ở dòng 21.
	for _, f := range fs {
		if f.Rule == "mockup-tokens" && strings.Contains(f.Msg, `"#333"`) && f.Line != 17 {
			t.Errorf("dòng mockup-tokens = %d, muốn 17", f.Line)
		}
		if f.Rule == "mockup-tokens" && strings.Contains(f.Msg, `"4px"`) && f.Line != 21 {
			t.Errorf("dòng mockup-tokens style= = %d, muốn 21", f.Line)
		}
	}
	if e, w := Count(fs); e == 0 || w == 0 {
		t.Fatalf("Count: %d lỗi %d cảnh báo", e, w)
	}
	// glossary-term: thuật ngữ đã định nghĩa, Given/When/Then ở mục 8 và chữ in đậm trong khối mã không báo; mỗi thuật ngữ một lần (gộp hoa thường), đúng dòng.
	n := 0
	for _, f := range fs {
		if f.Rule == "glossary-term" {
			n++
			if strings.Contains(f.Msg, "đơn hàng") || strings.Contains(f.Msg, "Given") || strings.Contains(f.Msg, "macro") || strings.Contains(f.Msg, "vận chuyển") || f.Line != 15 {
				t.Errorf("glossary-term sai: %+v", f)
			}
		}
	}
	if n != 1 {
		t.Errorf("glossary-term: %d finding, muốn 1", n)
	}
	// File không có frontmatter hoặc loại lạ không sinh finding ngoài ngưỡng dòng.
	for _, f := range fs {
		if strings.HasSuffix(f.File, "T-001.md") {
			t.Errorf("loại chưa có trong types.toml không được báo: %+v", f)
		}
	}
}

func TestStrict(t *testing.T) {
	if !has(runFixture(t, true), "idea.md", "frontmatter-required", Error, "created_by") {
		t.Fatal("--strict phải nâng created_by thành lỗi")
	}
}

func TestSorted(t *testing.T) {
	fs := runFixture(t, false)
	for i := 1; i < len(fs); i++ {
		if fs[i-1].File > fs[i].File {
			t.Fatalf("chưa sắp theo file: %s > %s", fs[i-1].File, fs[i].File)
		}
	}
}
