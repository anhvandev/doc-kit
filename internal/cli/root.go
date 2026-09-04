// Package cli định nghĩa các lệnh của dk trên cobra.
package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anhvandev/doc-kit/assets"
	"github.com/anhvandev/doc-kit/internal/config"
	"github.com/anhvandev/doc-kit/internal/doctype"
)

// Version được gán lúc build qua -ldflags -X; `go install <module>@vX` không
// có ldflags nên lấy phiên bản module từ build info.
var Version = "dev"

func init() {
	if Version != "dev" {
		return
	}
	if bi, ok := debug.ReadBuildInfo(); ok && bi.Main.Version != "" && bi.Main.Version != "(devel)" {
		Version = strings.TrimPrefix(bi.Main.Version, "v")
	}
}

// Mã thoát: 0 thành công, 1 lỗi I/O hoặc logic, 2 sai cờ, 3 kiểm tra không qua.
const (
	codeOK    = 0
	codeError = 1 // lỗi I/O hoặc logic
	codeUsage = 2 // sai cờ hoặc tham số
	codeCheck = 3 // kiểm tra không qua
)

// ExitError mang mã thoát và thông điệp cho main.
type ExitError struct {
	Code int
	Msg  string
}

func (e *ExitError) Error() string { return e.Msg }

func fail(code int, format string, a ...any) error {
	return &ExitError{Code: code, Msg: fmt.Sprintf(format, a...)}
}

// app giữ trạng thái chung của một lần chạy.
type app struct {
	json bool
	cwd  string
	out  io.Writer
	reg  doctype.Registry

	root string // gốc dự án (thư mục chứa dk.toml), rỗng khi chưa tìm
	cfg  config.Config
}

// Execute chạy CLI với os.Args và stdout.
func Execute() error {
	return Run(os.Args[1:], os.Stdout)
}

// Run chạy CLI với args và writer cho trước; dùng cho test.
func Run(args []string, out io.Writer) error {
	reg, err := doctype.Load(assets.FS)
	if err != nil {
		return fail(codeError, "bảng loại tài liệu nhúng lỗi: %v", err)
	}
	a := &app{out: out, reg: reg}
	rootCmd := newRootCmd(a)
	rootCmd.SetArgs(args)
	rootCmd.SetOut(out)
	rootCmd.SetErr(os.Stderr)
	err = rootCmd.Execute()
	var ee *ExitError
	if err != nil && !errors.As(err, &ee) {
		if strings.HasPrefix(err.Error(), "unknown command") {
			return fail(codeUsage, "%v", err)
		}
		return fail(codeError, "%v", err)
	}
	return err
}

func newRootCmd(a *app) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "dk",
		Short:         "Tạo và duy trì tài liệu dự án từ template nhúng",
		Version:       Version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.SetVersionTemplate("{{.Version}}\n")
	cmd.PersistentFlags().BoolVar(&a.json, "json", false, "in kết quả dạng JSON")
	cmd.PersistentFlags().StringVar(&a.cwd, "cwd", "", "thư mục làm việc (mặc định thư mục hiện tại)")
	cmd.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return fail(codeUsage, "%v", err)
	})
	cmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		if a.cwd == "" {
			wd, err := os.Getwd()
			if err != nil {
				return fail(codeError, "%v", err)
			}
			a.cwd = wd
		}
		abs, err := filepath.Abs(a.cwd)
		if err != nil {
			return fail(codeError, "%v", err)
		}
		a.cwd = abs
		return nil
	}
	cmd.AddCommand(newInitCmd(a), newTemplateCmd(a), newNewCmd(a), newChangelogCmd(a),
		newRenderCmd(a), newIndexCmd(a), newCheckCmd(a), newRefsCmd(a), newStatusCmd(a),
		newSkillCmd(a), newHookCmd(a), newTokensCmd(a), newDoctorCmd(a), newSelfCheckCmd(a))
	return cmd
}

// exactArgs như cobra.ExactArgs nhưng trả mã thoát 2.
func exactArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != n {
			return fail(codeUsage, "lệnh %s cần %d tham số, nhận %d", cmd.CommandPath(), n, len(args))
		}
		return nil
	}
}

// findProjectRoot đi lên từ dir đến thư mục có dk.toml.
func findProjectRoot(dir string) (string, bool) {
	for {
		if _, err := os.Stat(filepath.Join(dir, config.FileName)); err == nil {
			return dir, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}

// requireProject nạp gốc dự án và dk.toml; lỗi nếu chưa `dk init`.
func (a *app) requireProject() error {
	root, ok := findProjectRoot(a.cwd)
	if !ok {
		return fail(codeError, "không tìm thấy %s từ %s trở lên; chạy `dk init` trước", config.FileName, a.cwd)
	}
	cfg, err := config.Load(filepath.Join(root, config.FileName))
	if err != nil {
		return fail(codeError, "%v", err)
	}
	a.root, a.cfg = root, cfg
	return nil
}

func (a *app) docsDir() string {
	return filepath.Join(a.root, filepath.FromSlash(a.cfg.DocsDir))
}

// relRoot đổi đường dẫn tuyệt đối sang tương đối gốc dự án, dạng slash.
func (a *app) relRoot(p string) string {
	rel, err := filepath.Rel(a.root, p)
	if err != nil {
		return p
	}
	return filepath.ToSlash(rel)
}

func (a *app) printJSON(v any) error {
	enc := json.NewEncoder(a.out)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func (a *app) printf(format string, args ...any) {
	fmt.Fprintf(a.out, format, args...)
}
