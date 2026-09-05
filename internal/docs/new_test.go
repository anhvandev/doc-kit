package docs

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/anhvandev/doc-kit/assets"
	"github.com/anhvandev/doc-kit/internal/doctype"
	"github.com/anhvandev/doc-kit/internal/frontmatter"
)

var now = time.Date(2026, 9, 3, 14, 5, 0, 0, time.Local)

func setup(t *testing.T) (doctype.Registry, Options) {
	t.Helper()
	reg, err := doctype.Load(assets.FS)
	if err != nil {
		t.Fatal(err)
	}
	return reg, Options{DocsDir: filepath.Join(t.TempDir(), "docs"), Owner: "an", Version: "test", Now: now}
}

func readFM(t *testing.T, path string) map[string]any {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	fm, _, ok := frontmatter.Split(b)
	if !ok {
		t.Fatalf("%s: không có frontmatter", path)
	}
	return frontmatter.Map(fm)
}

func TestSeqID(t *testing.T) {
	reg, o := setup(t)
	r1, err := New(reg, "feature-spec", "bo-loc-don-hang", o)
	if err != nil {
		t.Fatal(err)
	}
	if r1.ID != "F-001" || filepath.Base(r1.Path) != "F-001-bo-loc-don-hang.md" {
		t.Fatalf("lần 1: %+v", r1)
	}
	fm := readFM(t, r1.Path)
	for _, k := range []string{"id", "type", "title", "status", "owner", "created", "updated", "source", "created_by", "dk_version"} {
		if _, ok := fm[k]; !ok {
			t.Errorf("thiếu trường %s", k)
		}
	}
	if fm["title"] != "Bo loc don hang" || fm["owner"] != "an" || fm["created_by"] != "dk" || fm["dk_version"] != "test" {
		t.Fatalf("frontmatter sai: %v", fm)
	}
	r2, _ := New(reg, "feature-spec", "xuat-excel", o)
	if r2.ID != "F-002" {
		t.Fatalf("lần 2 phải F-002, được %s", r2.ID)
	}
	r3, _ := New(reg, "feature-spec", "xuat-excel", o)
	if r3.ID != "F-003" {
		t.Fatalf("seq luôn tăng, trùng slug vẫn ra id mới: %s", r3.ID)
	}
}

func TestIDPrefixAndSet(t *testing.T) {
	reg, o := setup(t)
	o.IDPrefix = "SHOP-"
	o.Set = map[string]string{"title": "Sửa lỗi: mất đơn", "owner": "bình"}
	r, err := New(reg, "feature-spec", "bo-loc", o)
	if err != nil || r.ID != "SHOP-F-001" {
		t.Fatalf("%+v %v", r, err)
	}
	r2, _ := New(reg, "feature-spec", "khac", o)
	if r2.ID != "SHOP-F-002" {
		t.Fatalf("tiền tố phải đếm tiếp: %s", r2.ID)
	}
	fm := readFM(t, r.Path)
	if fm["title"] != "Sửa lỗi: mất đơn" || fm["owner"] != "bình" {
		t.Fatalf("--set không áp dụng hoặc dấu hai chấm làm hỏng frontmatter: %v", fm)
	}
	if fm["format"] != "spec" || fm["has_ui"] != true {
		t.Fatalf("mặc định format=spec, has_ui=true (bool): %v", fm)
	}
	// format và has_ui đổi biến thể thân; has_ui ghi kiểu bool; format lạ là lỗi.
	o.Set = map[string]string{"format": "crud", "has_ui": "false"}
	r3, err := New(reg, "feature-spec", "crud", o)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(r3.Path)
	fm = readFM(t, r3.Path)
	if fm["format"] != "crud" || fm["has_ui"] != false || bytes.Contains(b, []byte("## 3.")) || bytes.Contains(b, []byte("## 5.")) || !bytes.Contains(b, []byte("## 4. Bảng field")) {
		t.Fatalf("biến thể crud, has_ui=false sai:\n%s", b)
	}
	o.Set = map[string]string{"has_ui": "no"}
	if _, err := New(reg, "feature-spec", "hasui", o); err == nil || !strings.Contains(err.Error(), "has_ui") {
		t.Fatalf("has_ui lạ phải lỗi: %v", err)
	}
	o.Set = map[string]string{"format": "xyz"}
	if _, err := New(reg, "feature-spec", "xyz", o); err == nil || !strings.Contains(err.Error(), "format") {
		t.Fatalf("format lạ phải lỗi: %v", err)
	}
}

func TestIntakeChainAndFrom(t *testing.T) {
	reg, o := setup(t)
	idea, err := New(reg, "idea", "bo-loc-don-hang", o)
	if err != nil {
		t.Fatal(err)
	}
	wantDir := filepath.Join(o.DocsDir, "intake", "260903-bo-loc-don-hang")
	if filepath.Dir(idea.Path) != wantDir || idea.ID != "" {
		t.Fatalf("idea: %+v", idea)
	}
	o.From = idea.Path
	o.Now = now.Add(48 * time.Hour) // ngày khác vẫn vào cùng thư mục nguồn
	iv, err := New(reg, "interview", "bo-loc-don-hang", o)
	if err != nil || filepath.Dir(iv.Path) != wantDir {
		t.Fatalf("interview: %+v %v", iv, err)
	}
	br, err := New(reg, "brief", "bo-loc-don-hang", o)
	if err != nil || filepath.Dir(br.Path) != wantDir {
		t.Fatalf("brief: %+v %v", br, err)
	}
	ents, _ := os.ReadDir(wantDir)
	if len(ents) != 3 {
		t.Fatalf("thư mục intake phải có 3 file, có %d", len(ents))
	}
	if fm := readFM(t, br.Path); fm["source"] != "260903-bo-loc-don-hang/idea.md" {
		t.Fatalf("source của brief: %v", fm["source"])
	}

	// brief -> feature-spec chép title, purpose <- outcome, acceptance
	b, _ := os.ReadFile(br.Path)
	s := strings.Replace(string(b), `outcome: ""`, "outcome: Lọc đơn theo trạng thái", 1)
	s = strings.Replace(s, "acceptance: []", "acceptance:\n  - Lọc dưới 1 giây\n  - Giữ bộ lọc khi tải lại", 1)
	s = strings.Replace(s, "title: Bo loc don hang", "title: Bộ lọc đơn hàng", 1)
	os.WriteFile(br.Path, []byte(s), 0o644)
	o.From = br.Path
	fs, err := New(reg, "feature-spec", "bo-loc-don-hang", o)
	if err != nil {
		t.Fatal(err)
	}
	fm := readFM(t, fs.Path)
	acc, _ := fm["acceptance"].([]any)
	if fm["title"] != "Bộ lọc đơn hàng" || fm["purpose"] != "Lọc đơn theo trạng thái" || len(acc) != 2 {
		t.Fatalf("--from brief chép sai: %v", fm)
	}
	if fm["source"] != "260903-bo-loc-don-hang/brief.md" {
		t.Fatalf("source sai: %v", fm["source"])
	}
}

func TestDateID(t *testing.T) {
	reg, o := setup(t)
	r, err := New(reg, "cr", "doi-cach-loc", o)
	if err != nil || r.ID != "CR-260903-doi-cach-loc" || filepath.Base(r.Path) != "CR-260903-doi-cach-loc.md" {
		t.Fatalf("%+v %v", r, err)
	}
	o.From = r.Path
	br, err := New(reg, "brief", "doi-cach-loc", o)
	if err != nil || filepath.Dir(br.Path) != filepath.Join(o.DocsDir, "intake", "260903-doi-cach-loc") {
		t.Fatalf("--from ngoài intake phải tạo subdir riêng: %+v %v", br, err)
	}
	iv, err := New(reg, "interview", "doi-cach-loc", o)
	if err != nil || filepath.Dir(iv.Path) != filepath.Join(o.DocsDir, "cr", "CR-260903-doi-cach-loc") {
		t.Fatalf("interview --from cr phải vào thư mục cùng tên CR: %+v %v", iv, err)
	}
	if fm := readFM(t, iv.Path); fm["source"] != "CR-260903-doi-cach-loc" {
		t.Fatalf("source interview sai: %v", fm["source"])
	}
	o.From = ""
	if _, err := New(reg, "cr", "Sai Slug", o); err == nil {
		t.Fatal("slug sai phải lỗi")
	}
	if _, err := New(reg, "cr", "doi-cach-loc", o); !errors.Is(err, ErrExists) {
		t.Fatalf("file tồn tại phải ErrExists, được %v", err)
	}
	o.Force = true
	if _, err := New(reg, "cr", "doi-cach-loc", o); err != nil {
		t.Fatalf("--force phải ghi đè: %v", err)
	}
}
