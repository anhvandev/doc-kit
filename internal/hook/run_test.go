package hook

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anhvandev/doc-kit/internal/changelog"
	"github.com/anhvandev/doc-kit/internal/config"
)

func project(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := config.Write(filepath.Join(root, config.FileName), config.Default("t")); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"docs/features", "docs/html"} {
		os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	os.WriteFile(filepath.Join(root, "docs", changelog.FileName), []byte(changelog.Title+"\n"), 0o644)
	return root
}

func call(t *testing.T, event, payload string) (out, errOut string) {
	t.Helper()
	var o, e bytes.Buffer
	if code := Run(event, strings.NewReader(payload), &o, &e); code != 0 {
		t.Fatalf("mã thoát %d", code)
	}
	return o.String(), e.String()
}

func writePayload(root, tool, rel string) string {
	return fmt.Sprintf(`{"tool_name":%q,"tool_input":{"file_path":%q},"cwd":%q}`, tool, filepath.Join(root, rel), root)
}

func TestPreWriteDeniesNewDoc(t *testing.T) {
	root := project(t)
	out, errOut := call(t, PreWrite, writePayload(root, "Write", "docs/features/thu.md"))
	if !strings.Contains(out, `"permissionDecision":"deny"`) || !strings.Contains(out, "dk new") {
		t.Fatalf("phải từ chối: %q %q", out, errOut)
	}
	// Đường dẫn tương đối cwd cũng bị bắt.
	out, _ = call(t, PreWrite, fmt.Sprintf(`{"tool_input":{"file_path":"docs/x.md"},"cwd":%q}`, root))
	if !strings.Contains(out, "deny") {
		t.Fatal("đường dẫn tương đối không bị bắt")
	}
}

func TestPreWriteSilentCases(t *testing.T) {
	root := project(t)
	existing := filepath.Join(root, "docs", "features", "co.md")
	os.WriteFile(existing, []byte("# có\n"), 0o644)
	cases := map[string]string{
		"file đã tồn tại": writePayload(root, "Write", "docs/features/co.md"),
		"html":            writePayload(root, "Write", "docs/html/a.md"),
		"ngoài docs":      writePayload(root, "Write", "README.md"),
		"không .md":       writePayload(root, "Write", "docs/a.txt"),
		"ngoài dự án":     `{"tool_input":{"file_path":"/nowhere/docs/a.md"},"cwd":"/nowhere"}`,
		"thiếu file_path": `{"tool_name":"Bash","tool_input":{"command":"ls"}}`,
	}
	for name, p := range cases {
		if out, errOut := call(t, PreWrite, p); out != "" || errOut != "" {
			t.Fatalf("%s: phải im lặng, được %q %q", name, out, errOut)
		}
	}
}

func TestPostEditRecordsChangelog(t *testing.T) {
	root := project(t)
	doc := filepath.Join(root, "docs", "x.md")
	os.WriteFile(doc, []byte("---\ntitle: x\n---\n# x\n"), 0o644)
	out, errOut := call(t, PostEdit, writePayload(root, "Edit", "docs/x.md"))
	if out != "" || errOut != "" {
		t.Fatalf("%q %q", out, errOut)
	}
	cl, _ := os.ReadFile(filepath.Join(root, "docs", changelog.FileName))
	if !strings.Contains(string(cl), "| x.md | không git, 4 dòng | chưa tóm tắt | trực-tiếp") {
		t.Fatalf("changelog:\n%s", cl)
	}
	// Hook không sửa file agent vừa ghi (tránh "file modified since read").
	if b, _ := os.ReadFile(doc); strings.Contains(string(b), "updated:") {
		t.Fatal("hook không được bump updated")
	}

	os.WriteFile(filepath.Join(root, "docs", "html", "a.html"), []byte("<p>"), 0o644)
	call(t, PostEdit, writePayload(root, "Edit", "docs/html/a.html"))
	gen := filepath.Join(root, "docs", "features", "README.md")
	os.WriteFile(gen, []byte("---\ngenerated: true\n---\n"), 0o644)
	call(t, PostEdit, writePayload(root, "Edit", "docs/features/README.md"))
	cl2, _ := os.ReadFile(filepath.Join(root, "docs", changelog.FileName))
	if string(cl2) != string(cl) {
		t.Fatalf("html hoặc generated bị ghi:\n%s", cl2)
	}
}

func TestFailOpen(t *testing.T) {
	out, errOut := call(t, PreWrite, "không phải json")
	if out != "" || !strings.Contains(errOut, "dk hook pre-write") {
		t.Fatalf("stdin rác: %q %q", out, errOut)
	}
	root := project(t)
	_, errOut = call(t, "la", writePayload(root, "Write", "docs/a.md"))
	if !strings.Contains(errOut, "sự kiện lạ") {
		t.Fatal(errOut)
	}
}
