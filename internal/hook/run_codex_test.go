package hook

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anhvandev/doc-kit/internal/changelog"
)

// codexPayload đọc payload mẫu apply_patch, thay __ROOT__ bằng dự án tạm.
func codexPayload(t *testing.T, root string) string {
	t.Helper()
	b, err := os.ReadFile("testdata/codex-pre-tool-use-apply-patch.json")
	if err != nil {
		t.Fatal(err)
	}
	return strings.ReplaceAll(string(b), "__ROOT__", filepath.ToSlash(root))
}

func TestPatchFiles(t *testing.T) {
	got := patchFiles("*** Begin Patch\n*** Add File: a.md\n+x\n*** Update File: b/c.md\n*** Move to: docs/e.md\n@@\n*** Delete File: d.md\n*** End Patch")
	if strings.Join(got, ",") != "a.md,b/c.md,docs/e.md" {
		t.Fatal(got)
	}
}

func TestCodexApplyPatchDeniesNewDoc(t *testing.T) {
	root := project(t)
	os.WriteFile(filepath.Join(root, "docs", "features", "co.md"), []byte("---\ntitle: có\n---\n# có\n"), 0o644)
	out, errOut := call(t, PreWrite, codexPayload(t, root))
	if strings.Count(out, `"permissionDecision":"deny"`) != 1 || errOut != "" {
		t.Fatalf("phải từ chối đúng một lần: %q %q", out, errOut)
	}
	// Đổi tên vào docs/ là tạo tài liệu mới: từ chối.
	mv := strings.ReplaceAll(codexPayload(t, root), "*** Add File: docs/features/moi.md\n+# Mới\n", "*** Update File: ghi-chu.md\n*** Move to: docs/features/chuyen.md\n")
	if out, _ := call(t, PreWrite, mv); !strings.Contains(out, "deny") {
		t.Fatalf("move vào docs/ phải bị từ chối: %q", out)
	}
	// Patch chỉ sửa file có sẵn: im lặng.
	upd := strings.ReplaceAll(codexPayload(t, root), "*** Add File: docs/features/moi.md\\n+# Mới\\n", "")
	if out, errOut := call(t, PreWrite, upd); out != "" || errOut != "" {
		t.Fatalf("update phải im lặng: %q %q", out, errOut)
	}
}

func TestCodexApplyPatchRecordsChangelog(t *testing.T) {
	root := project(t)
	os.WriteFile(filepath.Join(root, "docs", "features", "co.md"), []byte("---\ntitle: có\n---\n# có sửa\n"), 0o644)
	os.WriteFile(filepath.Join(root, "docs", "features", "moi.md"), []byte("# Mới\n"), 0o644)
	post := strings.ReplaceAll(codexPayload(t, root), "PreToolUse", "PostToolUse")
	if out, errOut := call(t, PostEdit, post); out != "" || errOut != "" {
		t.Fatalf("%q %q", out, errOut)
	}
	cl, _ := os.ReadFile(filepath.Join(root, "docs", changelog.FileName))
	for _, f := range []string{"features/moi.md", "features/co.md"} {
		if !strings.Contains(string(cl), "| "+f+" |") {
			t.Fatalf("thiếu %s trong changelog:\n%s", f, cl)
		}
	}
}
