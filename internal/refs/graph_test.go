package refs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anhvandev/doc-kit/internal/docs"
)

func TestBuild(t *testing.T) {
	root := t.TempDir()
	w := func(rel, content string) {
		p := filepath.Join(root, filepath.FromSlash(rel))
		_ = os.MkdirAll(filepath.Dir(p), 0o755)
		_ = os.WriteFile(p, []byte(content), 0o644)
	}
	w("docs/cr/CR-260901-a.md", "---\nid: CR-260901-a\ntype: cr\n---\n# A\nXem [B](../features/F-001-b.md) và [ngoài](https://x/y.md) và [self](CR-260901-a.md).\n")
	w("docs/cr/CR-260901-a-b.md", "---\nid: CR-260901-a-b\ntype: cr\n---\n# AB\n")
	w("docs/features/F-001-b.md", "---\nid: F-001\ntype: feature-spec\nsource: CR-260901-a\n---\n# B\nNhắc CR-260901-a-b (id dài hơn), không phải CR-260901-a.\n")
	w("docs/features/README.md", "---\ngenerated: true\n---\n[B](F-001-b.md)\n")
	w("plans/p.md", "# plan\n[spec](../docs/features/F-001-b.md)\n")
	metas, err := docs.Scan(root, "docs", "plans")
	if err != nil {
		t.Fatal(err)
	}
	g := Build(root, metas)
	a, b := "docs/cr/CR-260901-a.md", "docs/features/F-001-b.md"
	check := func(got []string, want string) {
		t.Helper()
		if strings.Join(got, ",") != want {
			t.Errorf("được %v, muốn %s", got, want)
		}
	}
	check(g.Out[a], b)
	check(g.In[a], b)
	check(g.Out[b], "docs/cr/CR-260901-a-b.md,"+a)
	check(g.In[b], a+",plans/p.md")
	if _, ok := g.Out["docs/features/README.md"]; ok {
		t.Error("file generated không vào đồ thị")
	}
}
