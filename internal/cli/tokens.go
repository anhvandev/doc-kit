package cli

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/anhvandev/doc-kit/internal/tokens"
)

const (
	tokensJSON = "design/tokens/tokens.json"
	tokensCSS  = "design/tokens/tokens.css"
)

func newTokensCmd(a *app) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens",
		Short: "Design tokens: sinh CSS variables từ tokens.json",
	}
	var in, out string
	css := &cobra.Command{
		Use:   "css",
		Short: "Sinh tokens.css (:root và [data-theme]) từ tokens.json khung W3C Design Tokens",
		Args:  exactArgs(0),
		RunE:  func(_ *cobra.Command, _ []string) error { return a.tokensCSS(in, out) },
	}
	css.Flags().StringVar(&in, "in", "", "file tokens.json (mặc định docs/"+tokensJSON+")")
	css.Flags().StringVar(&out, "out", "", "file CSS đích (mặc định docs/"+tokensCSS+")")
	cmd.AddCommand(css)
	return cmd
}

func (a *app) tokensCSS(in, out string) error {
	if err := a.requireProject(); err != nil {
		return err
	}
	in = a.resolve(in, filepath.Join(a.docsDir(), filepath.FromSlash(tokensJSON)))
	out = a.resolve(out, filepath.Join(a.docsDir(), filepath.FromSlash(tokensCSS)))
	b, err := os.ReadFile(in)
	if err != nil {
		return fail(codeError, "%v", err)
	}
	set, err := tokens.Parse(b)
	if err != nil {
		return fail(codeError, "%v", err)
	}
	css, n, err := set.CSS()
	if err != nil {
		return fail(codeError, "%s: %v", a.relRoot(in), err)
	}
	if err := writeFile(out, css); err != nil {
		return fail(codeError, "%v", err)
	}
	if a.json {
		return a.printJSON(map[string]any{"in": a.relRoot(in), "out": a.relRoot(out), "variables": n})
	}
	a.printf("Đã sinh %s (%d biến) từ %s\n", a.relRoot(out), n, a.relRoot(in))
	return nil
}

// resolve trả về p tuyệt đối theo cwd, hoặc def khi p rỗng.
func (a *app) resolve(p, def string) string {
	if p == "" {
		return def
	}
	if !filepath.IsAbs(p) {
		return filepath.Join(a.cwd, p)
	}
	return p
}
