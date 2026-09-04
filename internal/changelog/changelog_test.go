package changelog

import (
	"strings"
	"testing"
	"time"
)

var t0 = time.Date(2026, 9, 3, 14, 5, 0, 0, time.Local)

func entry(t time.Time, sum, src string) Entry {
	return Entry{Time: t, Path: "features/F-001-x.md", Added: 3, Deleted: 1, Summary: sum, Source: src}
}

func TestParseEmptyAndMultiDay(t *testing.T) {
	f, err := Parse(nil)
	if err != nil || len(f.Days) != 0 {
		t.Fatal("file rỗng phải cho kết quả rỗng")
	}
	in := Title + "\n\n## 2026-09-03\n\n- 14:10 | features/F-012.md | +18 −4 | Thêm ngoại lệ B4 | CR-1\n- 14:05 | design/mockups/F-012-B3.html | mới, 210 dòng | Mockup | CR-1\n\n## 2026-09-02\n\n- 09:00 | intake/x/idea.md | không git, 40 dòng | chưa tóm tắt | -\n"
	f, err = Parse([]byte(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(f.Days) != 2 || len(f.Days[0].Entries) != 2 || f.Days[0].Entries[0].Added != 18 || !f.Days[0].Entries[1].New || f.Days[0].Entries[1].Lines != 210 {
		t.Fatalf("parse sai: %+v", f)
	}
	if e := f.Days[1].Entries[0]; !e.NoGit || e.Source != "" {
		t.Fatalf("no git parse sai: %+v", e)
	}
	if string(f.Format()) != in {
		t.Fatalf("Format không giữ định dạng:\n%s", f.Format())
	}
}

func TestMerge(t *testing.T) {
	win := 10 * time.Minute
	f := &File{}
	f.Add(entry(t0, "", "CR-1"), win)
	if got := f.Add(entry(t0.Add(5*time.Minute), "Sửa B2", "CR-1"), win); got.Summary != "Sửa B2" {
		t.Fatalf("Add phải trả mục đã gộp: %+v", got)
	}
	if got := f.Add(entry(t0.Add(6*time.Minute), "", "CR-1"), win); got.Summary != "Sửa B2" {
		t.Fatalf("tóm tắt rỗng phải giữ tóm tắt cũ và trả về đúng: %+v", got)
	}

	if n := len(f.Days[0].Entries); n != 1 {
		t.Fatalf("cách 5 phút cùng file cùng nguồn phải gộp, có %d", n)
	}
	if e := f.Days[0].Entries[0]; e.Summary != "Sửa B2" || e.Time != t0.Add(6*time.Minute) {
		t.Fatalf("mục gộp phải mang giờ và tóm tắt mới: %+v", e)
	}

	f.Add(entry(t0.Add(16*time.Minute), "Sửa B3", "CR-1"), win)
	if n := len(f.Days[0].Entries); n != 2 {
		t.Fatalf("cách 11 phút phải thành 2 dòng, có %d", n)
	}
	f.Add(entry(t0.Add(17*time.Minute), "Khác nguồn", "CR-2"), win)
	if n := len(f.Days[0].Entries); n != 3 {
		t.Fatalf("khác nguồn phải thành dòng mới, có %d", n)
	}
	if f.Days[0].Entries[0].Source != "CR-2" {
		t.Fatal("mục mới nhất phải ở trên")
	}

	n := len(f.Days[0].Entries)
	if f.Add(entry(t0.Add(-time.Minute), "Cũ hơn", "CR-2"), win); len(f.Days[0].Entries) != n+1 {
		t.Fatal("mục cũ hơn dòng đầu không được gộp")
	}
	f.Add(entry(t0.Add(48*time.Hour), "Ngày sau", ""), win)
	if f.Days[0].Date != "2026-09-05" || len(f.Days) != 2 {
		t.Fatalf("ngày mới phải lên đầu: %+v", f.Days)
	}
	if !strings.Contains(string(f.Format()), "| Ngày sau | -\n") {
		t.Fatalf("nguồn rỗng ghi '-':\n%s", f.Format())
	}
}

func TestSince(t *testing.T) {
	f := &File{}
	f.Add(entry(t0, "a", ""), 0)
	if !f.Since(t0.Add(30 * time.Second))["features/F-001-x.md"] {
		t.Fatal("mốc cùng phút phải tính")
	}
	if f.Since(t0.Add(2 * time.Minute))["features/F-001-x.md"] {
		t.Fatal("mốc sau mục thì không tính")
	}
}

func TestAddReplacesPlaceholderAcrossSource(t *testing.T) {
	f := &File{}
	now := time.Date(2026, 9, 3, 10, 0, 0, 0, time.Local)
	f.Add(Entry{Time: now, Path: "a.md", Summary: NoSummary, Source: "trực-tiếp", Added: 1}, MergeWindow)
	e := f.Add(Entry{Time: now.Add(time.Minute), Path: "a.md", Summary: "thật", Source: "CR-1", Added: 2}, MergeWindow)
	if len(f.Days[0].Entries) != 1 || e.Summary != "thật" || e.Source != "CR-1" {
		t.Fatalf("%+v", f.Days[0].Entries)
	}
	// Mục thật không gộp với mục thật khác nguồn.
	f.Add(Entry{Time: now.Add(2 * time.Minute), Path: "a.md", Summary: "khác", Source: "CR-2"}, MergeWindow)
	if len(f.Days[0].Entries) != 2 {
		t.Fatal("mục thật khác nguồn không được gộp")
	}
}
