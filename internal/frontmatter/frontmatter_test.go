package frontmatter

import (
	"bytes"
	"testing"
)

const doc = "---\nid: F-001\ntitle: Bộ lọc\ncreated: 2026-09-03\n---\n\n# Thân\n\ndòng {{ giữ nguyên }}\n"

func TestSplitJoinRoundTrip(t *testing.T) {
	fm, body, ok := Split([]byte(doc))
	if !ok {
		t.Fatal("phải tách được frontmatter")
	}
	if GetString(fm, "id") != "F-001" || GetString(fm, "created") != "2026-09-03" {
		t.Fatalf("đọc sai: %v", Map(fm))
	}
	if string(body) != "\n# Thân\n\ndòng {{ giữ nguyên }}\n" {
		t.Fatalf("thân sai: %q", body)
	}
	out, err := Join(fm, body)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != doc {
		t.Fatalf("round-trip đổi nội dung:\n%s", out)
	}
}

func TestSetKeepsOrder(t *testing.T) {
	fm, body, _ := Split([]byte(doc))
	SetString(fm, "title", "Bộ lọc đơn hàng")
	SetString(fm, "updated", "2026-09-03 14:05")
	out, _ := Join(fm, body)
	want := "---\nid: F-001\ntitle: Bộ lọc đơn hàng\ncreated: 2026-09-03\nupdated: 2026-09-03 14:05\n---\n"
	if !bytes.HasPrefix(out, []byte(want)) {
		t.Fatalf("thứ tự khóa sai:\n%s", out)
	}
}

func TestNoFrontmatter(t *testing.T) {
	in := []byte("# Chỉ có thân\n")
	fm, body, ok := Split(in)
	if ok || fm != nil || !bytes.Equal(body, in) {
		t.Fatal("không có frontmatter thì ok=false và body nguyên")
	}
	_, _, ok = Split([]byte("---\nid: 1\nkhông đóng\n"))
	if ok {
		t.Fatal("frontmatter không đóng phải ok=false")
	}
}

func TestCRLF(t *testing.T) {
	in := []byte("---\r\nid: X\r\n---\r\nthân\r\n")
	fm, body, ok := Split(in)
	if !ok || GetString(fm, "id") != "X" || string(body) != "thân\r\n" {
		t.Fatalf("CRLF tách sai: ok=%v body=%q", ok, body)
	}
	out, _ := Join(fm, body)
	if !bytes.Equal(out, in) {
		t.Fatalf("CRLF phải giữ nguyên:\n%q", out)
	}
}

func TestEmptyFrontmatter(t *testing.T) {
	fm, body, ok := Split([]byte("---\n---\nthân\n"))
	if !ok || len(fm.Content) != 0 || string(body) != "thân\n" {
		t.Fatal("frontmatter rỗng phải hợp lệ")
	}
	SetString(fm, "id", "A")
	out, _ := Join(fm, body)
	if string(out) != "---\nid: A\n---\nthân\n" {
		t.Fatalf("ghi lại sai:\n%s", out)
	}
}
