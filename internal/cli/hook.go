package cli

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/anhvandev/doc-kit/internal/hook"
	"github.com/anhvandev/doc-kit/internal/target"
)

// stdin là nguồn stdin của `hook run`; test thay bằng buffer.
var stdin io.Reader = os.Stdin

func newHookCmd(a *app) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hook",
		Short: "Cài và gỡ hook agent chặn tạo tài liệu tay, tự ghi changelog",
	}
	var inst, uninst targetFlags
	install := &cobra.Command{
		Use:   "install",
		Short: "Thêm hook dk vào cấu hình của target",
		Args:  exactArgs(0),
		RunE: func(_ *cobra.Command, _ []string) error {
			ts, err := a.resolveTargets(inst)
			if err != nil {
				return err
			}
			note := ""
			for _, t := range ts {
				if err := hook.Install(t, inst.global); err != nil {
					return fail(codeError, "%s: %v", t.Name(), err)
				}
				if t.Name() == "codex" {
					note = target.CodexTrustNote
				}
			}
			return a.printHookResult("đã cài", len(hook.Entries())*len(ts), note)
		},
	}
	inst.bind(install, true)
	uninstall := &cobra.Command{
		Use:   "uninstall",
		Short: "Gỡ hook dk khỏi cấu hình của target, giữ nguyên hook khác",
		Args:  exactArgs(0),
		RunE: func(_ *cobra.Command, _ []string) error {
			ts, err := a.resolveTargets(uninst)
			if err != nil {
				return err
			}
			for _, t := range ts {
				if err := hook.Uninstall(t, uninst.global); err != nil {
					return fail(codeError, "%s: %v", t.Name(), err)
				}
			}
			return a.printHookResult("đã gỡ", len(hook.Entries())*len(ts), "")
		},
	}
	uninst.bind(uninstall, true)
	run := &cobra.Command{
		Use:    "run <pre-write|post-edit>",
		Short:  "Do harness gọi; đọc JSON stdin, luôn thoát 0",
		Hidden: true,
		Args:   exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hook.Run(args[0], stdin, a.out, cmd.ErrOrStderr())
			return nil
		},
	}
	cmd.AddCommand(install, uninstall, run)
	return cmd
}

// printHookResult in kết quả; note là nhắc thêm (trust hook Codex), rỗng thì bỏ.
func (a *app) printHookResult(action string, n int, note string) error {
	if a.json {
		out := map[string]any{"action": action, "hooks": n}
		if note != "" {
			out["note"] = note
		}
		return a.printJSON(out)
	}
	a.printf("%s %d hook\n", action, n)
	if note != "" {
		a.printf("%s\n", note)
	}
	return nil
}
