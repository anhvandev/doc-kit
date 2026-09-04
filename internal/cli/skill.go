package cli

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anhvandev/doc-kit/internal/skill"
	"github.com/anhvandev/doc-kit/internal/target"
)

// targetFlags là cờ chung của skill *, hook * và doctor.
type targetFlags struct {
	name   string // một tên hoặc danh sách phân tách phẩy: claude,codex
	global bool
}

func (f *targetFlags) bind(cmd *cobra.Command, withGlobal bool) {
	cmd.Flags().StringVar(&f.name, "target", "claude", "harness đích: claude, codex hoặc claude,codex")
	if withGlobal {
		cmd.Flags().BoolVar(&f.global, "global", false, "scope toàn máy ($HOME) thay vì dự án")
	}
}

// names tách danh sách --target, bỏ trùng, giữ thứ tự.
func (f targetFlags) names() []string {
	var out []string
	for _, n := range strings.Split(f.name, ",") {
		n = strings.TrimSpace(n)
		if n != "" && !slices.Contains(out, n) {
			out = append(out, n)
		}
	}
	return out
}

// resolveTargets lấy các target theo --target; scope dự án cần dk.toml, scope
// toàn máy thì không. Tên lạ là mã 2 trước khi đụng bất kỳ target nào.
func (a *app) resolveTargets(f targetFlags) ([]target.Target, error) {
	if !f.global {
		if err := a.requireProject(); err != nil {
			return nil, err
		}
	}
	return a.targetsOf(f.names(), a.root)
}

func (a *app) targetsOf(names []string, root string) ([]target.Target, error) {
	var ts []target.Target
	for _, n := range names {
		t, err := target.Get(n, root)
		if err != nil {
			return nil, fail(codeUsage, "%v", err)
		}
		ts = append(ts, t)
	}
	if len(ts) == 0 {
		return nil, fail(codeUsage, "--target trống; chọn claude, codex hoặc claude,codex")
	}
	return ts, nil
}

func newSkillCmd(a *app) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skill",
		Short: "Cài, gỡ và kiểm tra bộ skill nhúng cho harness agent",
	}
	list := &cobra.Command{
		Use:   "list",
		Short: "Liệt kê skill nhúng trong binary",
		Args:  exactArgs(0),
		RunE: func(_ *cobra.Command, _ []string) error {
			metas, err := skill.List()
			if err != nil {
				return fail(codeError, "%v", err)
			}
			if a.json {
				return a.printJSON(metas)
			}
			for _, m := range metas {
				a.printf("%s\t%s\n", m.Name, m.Description)
			}
			return nil
		},
	}
	var inst, uninst, st targetFlags
	var force bool
	install := &cobra.Command{
		Use:   "install [<tên>...]",
		Short: "Cài skill vào thư mục skill của target (mặc định mọi skill)",
		RunE: func(_ *cobra.Command, args []string) error {
			ts, err := a.resolveTargets(inst)
			if err != nil {
				return err
			}
			var all []skill.Result
			for _, t := range ts {
				res, err := skill.Install(t, args, inst.global, force, Version)
				all = append(all, res...)
				if err != nil {
					return a.printResults(all, err)
				}
			}
			return a.printResults(all, nil)
		},
	}
	inst.bind(install, true)
	install.Flags().BoolVar(&force, "force", false, "ghi đè skill đã sửa tay")
	uninstall := &cobra.Command{
		Use:   "uninstall [<tên>...]",
		Short: "Gỡ skill do dk cài (mặc định mọi skill)",
		RunE: func(_ *cobra.Command, args []string) error {
			ts, err := a.resolveTargets(uninst)
			if err != nil {
				return err
			}
			var all []skill.Result
			for _, t := range ts {
				res, err := skill.Uninstall(t, args, uninst.global)
				all = append(all, res...)
				if err != nil {
					return a.printResults(all, err)
				}
			}
			return a.printResults(all, nil)
		},
	}
	uninst.bind(uninstall, true)
	status := &cobra.Command{
		Use:   "status",
		Short: "Trạng thái cài của từng skill theo scope dự án và toàn máy",
		Args:  exactArgs(0),
		RunE: func(_ *cobra.Command, _ []string) error {
			return a.skillStatus(st)
		},
	}
	st.bind(status, false)
	cmd.AddCommand(list, install, uninstall, status)
	return cmd
}

// printResults in kết quả đã làm được rồi mới trả lỗi (nếu có).
func (a *app) printResults(res []skill.Result, err error) error {
	if a.json {
		if perr := a.printJSON(res); perr != nil {
			return perr
		}
	} else {
		for _, r := range res {
			line := fmt.Sprintf("%s\t%s\t%s", r.Name, r.Action, r.Path)
			if r.Note != "" {
				line += " (" + r.Note + ")"
			}
			a.printf("%s\n", line)
		}
	}
	if err != nil {
		return fail(codeError, "%v", err)
	}
	return nil
}

func (a *app) skillStatus(f targetFlags) error {
	var rows []skill.Row
	// Scope dự án chỉ khi đang trong dự án; toàn máy luôn xét.
	root, inProject := findProjectRoot(a.cwd)
	if inProject {
		a.root = root
	}
	ts, err := a.targetsOf(f.names(), root)
	if err != nil {
		return err
	}
	for _, t := range ts {
		if inProject {
			r, err := skill.Status(t, false, Version)
			if err != nil {
				return fail(codeError, "%v", err)
			}
			rows = append(rows, r...)
		}
		r, err := skill.Status(t, true, Version)
		if err != nil {
			return fail(codeError, "%v", err)
		}
		rows = append(rows, r...)
	}
	if a.json {
		return a.printJSON(rows)
	}
	a.printf("%s\n", strings.Join([]string{"tên", "target", "scope", "trạng thái"}, "\t"))
	for _, r := range rows {
		a.printf("%s\t%s\t%s\t%s\n", r.Name, r.Target, r.Scope, r.State)
	}
	return nil
}
