package cli

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/anhvandev/doc-kit/internal/check"
	"github.com/anhvandev/doc-kit/internal/docs"
)

func newCheckCmd(a *app) *cobra.Command {
	var strict bool
	cmd := &cobra.Command{
		Use:   "check [<file>]",
		Short: "Kiểm tra frontmatter, liên kết, mã bước, liên kết ngược, ngưỡng dòng; mã thoát 3 khi có lỗi",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return fail(codeUsage, "check nhận tối đa một file")
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			var file string
			if len(args) == 1 {
				file = args[0]
			}
			return a.runCheck(file, strict)
		},
	}
	cmd.Flags().BoolVar(&strict, "strict", false, "created_by khác dk là lỗi thay vì warning")
	return cmd
}

// runChecks chạy mọi quy tắc trên docs/ và plans/; trả cả tập file đã quét.
func (a *app) runChecks(strict bool) ([]check.Finding, []docs.Meta, error) {
	metas, err := docs.Scan(a.root, a.cfg.DocsDir, a.cfg.PlansDir)
	if err != nil {
		return nil, nil, fail(codeError, "%v", err)
	}
	ctx := &check.Context{Root: a.root, DocsDir: a.cfg.DocsDir, Metas: metas, Reg: a.reg, Cfg: a.cfg.Check, Jargon: a.cfg.Release.Jargon, Strict: strict}
	return check.Run(ctx), metas, nil
}

// scanned báo rel có trong tập file đã quét.
func scanned(metas []docs.Meta, rel string) bool {
	for _, m := range metas {
		if m.Rel == rel {
			return true
		}
	}
	return false
}

func (a *app) runCheck(file string, strict bool) error {
	if err := a.requireProject(); err != nil {
		return err
	}
	findings, metas, err := a.runChecks(strict)
	if err != nil {
		return err
	}
	if file != "" {
		if !filepath.IsAbs(file) {
			file = filepath.Join(a.cwd, file)
		}
		rel := a.relRoot(filepath.Clean(file))
		if fi, err := os.Stat(file); err == nil && fi.IsDir() {
			return fail(codeUsage, "%s là thư mục; `dk check` nhận một file hoặc không tham số", rel)
		}
		if !scanned(metas, rel) {
			return fail(codeError, "%s không nằm trong %s/ hoặc %s/", rel, a.cfg.DocsDir, a.cfg.PlansDir)
		}
		kept := findings[:0]
		for _, f := range findings {
			if f.File == rel {
				kept = append(kept, f)
			}
		}
		findings = kept
	}
	if findings == nil {
		findings = []check.Finding{}
	}
	errs, warns := check.Count(findings)
	if a.json {
		if err := a.printJSON(findings); err != nil {
			return err
		}
	} else {
		for _, f := range findings {
			if f.Line > 0 {
				a.printf("%s:%d: %s %s: %s\n", f.File, f.Line, f.Level, f.Rule, f.Msg)
			} else {
				a.printf("%s: %s %s: %s\n", f.File, f.Level, f.Rule, f.Msg)
			}
		}
		a.printf("%d lỗi, %d cảnh báo\n", errs, warns)
	}
	if errs > 0 {
		return fail(codeCheck, "dk check: %d lỗi", errs)
	}
	return nil
}
