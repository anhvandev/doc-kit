package frontmatter

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestHTMLComment(t *testing.T) {
	src := "<!-- dk:\ntype: mockup\nfeature: F-001\n-->\n<!doctype html>\n<html></html>\n"
	fm, body, ok := SplitHTMLComment([]byte(src))
	if !ok || GetString(fm, "feature") != "F-001" || string(body) != "<!doctype html>\n<html></html>\n" {
		t.Fatalf("split: ok=%v body=%q", ok, body)
	}
	SetString(fm, "title", "Danh sách: rỗng")
	out, err := JoinHTMLComment(fm, body)
	if err != nil || !strings.HasPrefix(string(out), "<!-- dk:\ntype: mockup\nfeature: F-001\ntitle: 'Danh sách: rỗng'\n-->\n<!doctype") {
		t.Fatalf("join: %v\n%s", err, out)
	}
	if _, _, ok := SplitHTMLComment([]byte("<!doctype html>\n")); ok {
		t.Fatal("không có chú thích dk phải ok=false")
	}
}

func TestJSONKey(t *testing.T) {
	src := "{\n  \"color\": {\"blue\": {\"$value\": \"#00f\"}},\n  \"$dk\": {\"type\": \"design-tokens\", \"title\": \"Tokens\"},\n  \"space\": {}\n}\n"
	fm, body, ok := SplitJSONKey([]byte(src))
	if !ok || GetString(fm, "type") != "design-tokens" {
		t.Fatalf("split: ok=%v fm=%v", ok, fm)
	}
	if !json.Valid(body) || strings.Contains(string(body), "$dk") {
		t.Fatalf("body sau khi bỏ $dk: %s", body)
	}
	SetString(fm, "owner", "v")
	out, err := JoinJSONKey(fm, body)
	if err != nil || !json.Valid(out) {
		t.Fatalf("join: %v\n%s", err, out)
	}
	var m map[string]any
	_ = json.Unmarshal(out, &m)
	dk := m["$dk"].(map[string]any)
	if dk["owner"] != "v" || dk["title"] != "Tokens" || m["space"] == nil || m["color"] == nil {
		t.Fatalf("nội dung sau join: %s", out)
	}
	if !strings.HasPrefix(string(out), "{\n  \"$dk\": {\n    \"type\"") {
		t.Fatalf("$dk phải là khóa đầu, thứ tự giữ nguyên:\n%s", out)
	}
	// $dk ở cuối, không có dấu phẩy sau.
	src2 := "{\"a\": 1, \"$dk\": {\"type\": \"x\"}}"
	fm, body, ok = SplitJSONKey([]byte(src2))
	if !ok || !json.Valid(body) || string(body) != "{\"a\": 1}" {
		t.Fatalf("split cuối: ok=%v body=%q", ok, body)
	}
	if _, _, ok := SplitJSONKey([]byte("{\"a\": 1}")); ok {
		t.Fatal("thiếu $dk phải ok=false")
	}
}

func TestSplitFile(t *testing.T) {
	if _, _, ok := SplitFile("x.md", []byte("---\ntype: adr\n---\n")); !ok {
		t.Fatal("md")
	}
	if _, _, ok := SplitFile("x.html", []byte("<!-- dk:\ntype: mockup\n-->\n")); !ok {
		t.Fatal("html")
	}
	if _, _, ok := SplitFile("x.json", []byte("{\"$dk\": {}}")); !ok {
		t.Fatal("json")
	}
}
