package cli

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/anhvandev/doc-kit/internal/docs"
	"github.com/anhvandev/doc-kit/internal/refs"
)

func newRefsCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "refs <file>",
		Short: "In liên kết đi và đến của một tài liệu trong docs/ hoặc plans/",
		Args:  exactArgs(1),
		RunE:  func(_ *cobra.Command, args []string) error { return a.runRefs(args[0]) },
	}
}

func (a *app) runRefs(file string) error {
	if err := a.requireProject(); err != nil {
		return err
	}
	if !filepath.IsAbs(file) {
		file = filepath.Join(a.cwd, file)
	}
	rel := a.relRoot(filepath.Clean(file))
	metas, err := docs.Scan(a.root, a.cfg.DocsDir, a.cfg.PlansDir)
	if err != nil {
		return fail(codeError, "%v", err)
	}
	if !scanned(metas, rel) {
		return fail(codeError, "%s không nằm trong %s/ hoặc %s/", rel, a.cfg.DocsDir, a.cfg.PlansDir)
	}
	g := refs.Build(a.root, metas)
	out, in := g.Out[rel], g.In[rel]
	if out == nil {
		out = []string{}
	}
	if in == nil {
		in = []string{}
	}
	if a.json {
		return a.printJSON(map[string]any{"file": rel, "out": out, "in": in})
	}
	a.printf("%s\n", rel)
	a.printf("Đi (%d):\n", len(out))
	for _, p := range out {
		a.printf("  -> %s\n", p)
	}
	a.printf("Đến (%d):\n", len(in))
	for _, p := range in {
		a.printf("  <- %s\n", p)
	}
	return nil
}
