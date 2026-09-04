package cli

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anhvandev/doc-kit/internal/docs"
	"github.com/anhvandev/doc-kit/internal/render"
)

// IndexTitle là tiêu đề trang docs/html/index.html.
const IndexTitle = "Chỉ mục tài liệu"

func newRenderCmd(a *app) *cobra.Command {
	var all, index bool
	cmd := &cobra.Command{
		Use:   "render [<file>]",
		Short: "Đổi Markdown trong docs/ sang HTML một file tự chứa trong docs/html/",
		Args: func(cmd *cobra.Command, args []string) error {
			if (len(args) == 1) == all {
				return fail(codeUsage, "render cần đúng một file hoặc --all")
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			var file string
			if len(args) == 1 {
				file = args[0]
			}
			return a.runRender(file, index)
		},
	}
	cmd.Flags().BoolVar(&all, "all", false, "render mọi file .md trong docs/")
	cmd.Flags().BoolVar(&index, "index", false, "sinh thêm docs/html/index.html từ frontmatter mọi file")
	return cmd
}

type renderResult struct {
	Src     string `json:"src"`
	Out     string `json:"out"`
	Mermaid bool   `json:"mermaid"`
}

func (a *app) renderOpts() render.Options {
	return render.Options{DocsDir: a.docsDir(), OutDir: filepath.Join(a.docsDir(), "html")}
}

func (a *app) runRender(file string, index bool) error {
	if err := a.requireProject(); err != nil {
		return err
	}
	opts := a.renderOpts()
	var targets []docs.Meta
	if file != "" {
		if !strings.EqualFold(filepath.Ext(file), ".md") {
			return fail(codeUsage, "render chỉ nhận file .md, nhận %s", file)
		}
		if !filepath.IsAbs(file) {
			file = filepath.Join(a.cwd, file)
		}
		m, err := docs.Read(a.root, filepath.Clean(file))
		if err != nil {
			return fail(codeError, "%v", err)
		}
		if _, err := opts.OutPath(m.Path); err != nil {
			return fail(codeError, "%s nằm ngoài %s/", a.relRoot(m.Path), a.cfg.DocsDir)
		}
		targets = []docs.Meta{m}
	}
	var all []docs.Meta
	if file == "" || index {
		var err error
		if all, err = docs.Scan(a.root, a.cfg.DocsDir); err != nil {
			return fail(codeError, "%v", err)
		}
		all = slices.DeleteFunc(all, func(m docs.Meta) bool { return !m.IsMarkdown() })
		if file == "" {
			targets = all
		}
	}
	results := []renderResult{}
	for _, m := range targets {
		p, err := opts.Render(m)
		if err != nil {
			return fail(codeError, "%s: %v", m.Rel, err)
		}
		if err := writeFile(p.Out, p.HTML); err != nil {
			return fail(codeError, "%v", err)
		}
		results = append(results, renderResult{Src: m.Rel, Out: a.relRoot(p.Out), Mermaid: p.HasMermaid})
	}
	if index {
		b, err := opts.Index(IndexTitle, all)
		if err != nil {
			return fail(codeError, "%v", err)
		}
		out := filepath.Join(opts.OutDir, "index.html")
		if err := writeFile(out, b); err != nil {
			return fail(codeError, "%v", err)
		}
		results = append(results, renderResult{Src: "", Out: a.relRoot(out)})
	}
	if a.json {
		return a.printJSON(results)
	}
	for _, r := range results {
		a.printf("%s\n", r.Out)
	}
	return nil
}

func writeFile(path string, b []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
