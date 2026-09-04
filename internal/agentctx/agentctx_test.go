package agentctx

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteAndCheck(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "CLAUDE.md")
	content := []byte("## Rules\n\n- one\n")

	if st, _ := Check(path, content, "0.1.0"); st != StateMissingFile {
		t.Fatalf("chưa có file: %s", st)
	}
	if r, err := Write(path, content, "0.1.0"); err != nil || r != WriteCreated {
		t.Fatalf("tạo mới: %s %v", r, err)
	}
	if st, _ := Check(path, content, "0.1.0"); st != StateOK {
		t.Fatalf("sau tạo: %s", st)
	}
	if r, _ := Write(path, content, "0.1.0"); r != WriteUnchanged {
		t.Fatalf("chạy lại: %s", r)
	}

	// Phần ngoài khối được giữ nguyên khi thay khối.
	orig, _ := os.ReadFile(path)
	if err := os.WriteFile(path, []byte("# My project\n\nintro\n\n"+string(orig)+"\ntail\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if st, _ := Check(path, content, "0.2.0"); st != StateOutdated {
		t.Fatalf("khác phiên bản: %s", st)
	}
	if r, _ := Write(path, content, "0.2.0"); r != WriteUpdated {
		t.Fatalf("thay khối: %s", r)
	}
	got, _ := os.ReadFile(path)
	s := string(got)
	if !strings.HasPrefix(s, "# My project\n\nintro\n\n<!-- dk:agent-context start version=0.2.0 hash=") || !strings.HasSuffix(s, endMarker+"\n\ntail\n") {
		t.Fatalf("mất phần ngoài khối:\n%s", s)
	}
	if strings.Count(s, startPrefix) != 1 {
		t.Fatalf("khối bị nhân đôi:\n%s", s)
	}

	// Sửa tay bên trong khối: hash ở mốc không đổi nhưng nội dung khác.
	edited := strings.Replace(s, "- one", "- two", 1)
	if err := os.WriteFile(path, []byte(edited), 0o644); err != nil {
		t.Fatal(err)
	}
	if st, _ := Check(path, content, "0.2.0"); st != StateOutdated {
		t.Fatalf("sửa tay: %s", st)
	}
	if r, _ := Write(path, content, "0.2.0"); r != WriteUpdated {
		t.Fatalf("ghi đè sửa tay: %s", r)
	}
	if st, _ := Check(path, content, "0.2.0"); st != StateOK {
		t.Fatalf("sau ghi đè: %s", st)
	}
}

func TestAppendToExistingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "AGENTS.md")
	if err := os.WriteFile(path, []byte("# Agents\n\nnotes"), 0o644); err != nil { // không xuống dòng cuối
		t.Fatal(err)
	}
	if st, _ := Check(path, []byte("x\n"), "1"); st != StateMissing {
		t.Fatalf("chưa có khối: %s", st)
	}
	if r, err := Write(path, []byte("x\n"), "1"); err != nil || r != WriteUpdated {
		t.Fatalf("nối: %s %v", r, err)
	}
	got, _ := os.ReadFile(path)
	if want := "# Agents\n\nnotes\n\n<!-- dk:agent-context start version=1 hash=" + Hash([]byte("x\n")) + " -->\nx\n" + endMarker + "\n"; string(got) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestCRLFFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "CLAUDE.md")
	content := []byte("a\nb\n")
	if _, err := Write(path, content, "1"); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(path)
	crlf := strings.ReplaceAll("# Top\n\n"+string(b), "\n", "\r\n")
	if err := os.WriteFile(path, []byte(crlf), 0o644); err != nil {
		t.Fatal(err)
	}
	if st, _ := Check(path, content, "1"); st != StateOK {
		t.Fatalf("CRLF phải là ok: %s", st)
	}
	if r, _ := Write(path, content, "1"); r != WriteUnchanged {
		t.Fatalf("CRLF chạy lại: %s", r)
	}
	if r, _ := Write(path, []byte("c\n"), "1"); r != WriteUpdated {
		t.Fatal("CRLF thay khối")
	}
	got, _ := os.ReadFile(path)
	if strings.Count(string(got), startPrefix) != 1 || !strings.HasPrefix(string(got), "# Top\r\n\r\n<!-- dk:agent-context start") {
		t.Fatalf("CRLF nhân đôi hoặc mất đầu file:\n%q", got)
	}
}

func TestDuplicateBlocks(t *testing.T) {
	path := filepath.Join(t.TempDir(), "CLAUDE.md")
	one := Block([]byte("x\n"), "1")
	if err := os.WriteFile(path, []byte(one+"\n"+one), 0o644); err != nil {
		t.Fatal(err)
	}
	if st, _ := Check(path, []byte("x\n"), "1"); st != StateBroken {
		t.Fatalf("hai khối: %s", st)
	}
	if _, err := Write(path, []byte("x\n"), "1"); err == nil {
		t.Fatal("ghi lên file hai khối phải lỗi")
	}
}

func TestBrokenBlock(t *testing.T) {
	path := filepath.Join(t.TempDir(), "CLAUDE.md")
	if err := os.WriteFile(path, []byte("<!-- dk:agent-context start version=1 hash=abc -->\nx\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if st, _ := Check(path, []byte("x\n"), "1"); st != StateBroken {
		t.Fatalf("thiếu mốc đóng: %s", st)
	}
	if _, err := Write(path, []byte("x\n"), "1"); err == nil {
		t.Fatal("ghi lên khối hỏng phải lỗi")
	}
}
