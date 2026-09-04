package cli

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/anhvandev/doc-kit/internal/changelog"
	"github.com/anhvandev/doc-kit/internal/gitx"
)

func newChangelogCmd(a *app) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "changelog",
		Short: "Ghi và kiểm tra changelog tài liệu",
	}
	var summary, source string
	add := &cobra.Command{
		Use:   "add <file>",
		Short: "Ghi một dòng thay đổi cho file trong docs/",
		Args:  exactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.changelogAdd(args[0], summary, source)
		},
	}
	add.Flags().StringVar(&summary, "summary", "", "tóm tắt thay đổi")
	add.Flags().StringVar(&source, "source", "", "brief hoặc CR dẫn đến thay đổi")
	pending := &cobra.Command{
		Use:   "pending",
		Short: "Liệt kê file trong docs/ đã đổi mà chưa có dòng changelog; mã thoát 1 nếu còn",
		Args:  exactArgs(0),
		RunE:  func(_ *cobra.Command, _ []string) error { return a.changelogPending() },
	}
	cmd.AddCommand(add, pending)
	return cmd
}

func (a *app) changelogPath() string {
	return filepath.Join(a.docsDir(), ChangelogFile)
}

func (a *app) loadChangelog() (*changelog.File, error) {
	f, err := changelog.Load(a.docsDir())
	if err != nil {
		return nil, fail(codeError, "%s: %v", a.relRoot(a.changelogPath()), err)
	}
	return f, nil
}

// docRel chuẩn hóa đường dẫn về tương đối docs/ và từ chối file không ghi changelog.
func (a *app) docRel(p string) (abs, rel string, err error) {
	if !filepath.IsAbs(p) {
		p = filepath.Join(a.cwd, p)
	}
	abs = filepath.Clean(p)
	r, err := filepath.Rel(a.docsDir(), abs)
	if err != nil || strings.HasPrefix(r, "..") {
		return "", "", fail(codeError, "%s nằm ngoài %s/", a.relRoot(abs), a.cfg.DocsDir)
	}
	rel = filepath.ToSlash(r)
	if !a.tracksChangelog(rel) {
		return "", "", fail(codeError, "%s không ghi changelog (html/, file ẩn như .gitkeep, chỉ mục sinh ra hoặc chính changelog)", a.relRoot(abs))
	}
	return abs, rel, nil
}

// tracksChangelog báo file (tương đối docs/) có cần dòng changelog không.
func (a *app) tracksChangelog(rel string) bool {
	return changelog.Tracks(a.docsDir(), rel)
}

func (a *app) changelogAdd(file, summary, source string) error {
	if err := a.requireProject(); err != nil {
		return err
	}
	for name, v := range map[string]string{"--summary": summary, "--source": source} {
		if strings.ContainsAny(v, "\r\n") || strings.Contains(v, " | ") {
			return fail(codeUsage, "%s không được chứa xuống dòng hoặc ' | '", name)
		}
	}
	_, rel, err := a.docRel(file)
	if err != nil {
		return err
	}
	e, err := changelog.Record(a.root, a.docsDir(), rel, summary, source, time.Now(), true)
	if err != nil {
		return fail(codeError, "%v", err)
	}
	if a.json {
		return a.printJSON(map[string]any{"path": rel, "added": e.Added, "deleted": e.Deleted,
			"new": e.New, "no_git": e.NoGit, "lines": e.Lines, "summary": e.Summary, "source": e.Source})
	}
	a.printf("%s\n", e.String())
	return nil
}

// pendingDocs liệt kê file trong docs/ đã đổi mà chưa có dòng changelog.
func (a *app) pendingDocs() ([]string, error) {
	if !gitx.IsRepo(a.root) {
		return nil, fail(codeError, "%s không nằm trong repo git; `changelog pending` cần git", a.root)
	}
	gitRoot, err := gitx.Root(a.root)
	if err != nil {
		return nil, fail(codeError, "%v", err)
	}
	changed, err := gitx.ChangedDocs(gitRoot, a.docsDir())
	if err != nil {
		return nil, fail(codeError, "%v", err)
	}
	head, err := gitx.HeadTime(gitRoot)
	if err != nil {
		return nil, fail(codeError, "%v", err)
	}
	f, err := a.loadChangelog()
	if err != nil {
		return nil, err
	}
	logged := f.Since(head)
	pending := []string{}
	for _, p := range changed {
		abs := filepath.Join(gitRoot, filepath.FromSlash(p))
		r, err := filepath.Rel(a.docsDir(), abs)
		if err != nil || strings.HasPrefix(r, "..") {
			continue
		}
		rel := filepath.ToSlash(r)
		if !a.tracksChangelog(rel) || logged[rel] {
			continue
		}
		pending = append(pending, rel)
	}
	return pending, nil
}

func (a *app) changelogPending() error {
	if err := a.requireProject(); err != nil {
		return err
	}
	pending, err := a.pendingDocs()
	if err != nil {
		return err
	}
	if a.json {
		if err := a.printJSON(pending); err != nil {
			return err
		}
	} else {
		for _, p := range pending {
			a.printf("%s\n", p)
		}
	}
	if len(pending) > 0 {
		return fail(codeError, "%d file trong %s/ đổi mà chưa có dòng changelog; chạy `dk changelog add <file> --summary ...`", len(pending), a.cfg.DocsDir)
	}
	return nil
}
