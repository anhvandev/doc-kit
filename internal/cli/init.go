package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anhvandev/doc-kit/assets"
	"github.com/anhvandev/doc-kit/internal/changelog"
	"github.com/anhvandev/doc-kit/internal/config"
	"github.com/anhvandev/doc-kit/internal/gitx"
)

// ChangelogFile là tên file changelog trong docs/.
const ChangelogFile = changelog.FileName

// docsTree là các thư mục do `dk init` tạo trong docs/.
var docsTree = []string{
	"intake", "cr", "features", "adr", "overview",
	"design/tokens", "design/atoms", "design/molecules", "design/organisms", "design/templates",
	"design/patterns", "design/flows", "design/wireframes", "design/mockups",
	"test", "release", "ops", "governance", "html",
}

func newInitCmd(a *app) *cobra.Command {
	var force, agentContext bool
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Khởi tạo dk.toml, cây docs/ và changelog tài liệu",
		Args:  exactArgs(0),
		RunE: func(_ *cobra.Command, _ []string) error {
			if agentContext {
				return a.printAgentContext()
			}
			return a.runInit(force)
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "ghi đè dk.toml nếu đã có")
	cmd.Flags().BoolVar(&agentContext, "agent-context", false, "chỉ in khối Markdown để dán vào file ngữ cảnh agent, không ghi gì")
	return cmd
}

// printAgentContext in template nhúng ra stdout; người tự dán vào CLAUDE.md
// hoặc AGENTS.md, dk không sửa file đó.
func (a *app) printAgentContext() error {
	b, err := assets.FS.ReadFile("agent-context.md")
	if err != nil {
		return fail(codeError, "%v", err)
	}
	if a.json {
		return a.printJSON(map[string]string{"content": string(b)})
	}
	_, err = a.out.Write(b)
	return err
}

func (a *app) runInit(force bool) error {
	cfgPath := filepath.Join(a.cwd, config.FileName)
	cfg := config.Default(filepath.Base(a.cwd))
	if _, err := os.Stat(cfgPath); err == nil {
		if !force {
			return fail(codeError, "%s đã có; dùng --force để ghi đè", cfgPath)
		}
		// --force giữ cấu hình đã có, chỉ bổ sung cây thư mục còn thiếu.
		if cfg, err = config.Load(cfgPath); err != nil {
			return fail(codeError, "%v", err)
		}
	} else if root, ok := findProjectRoot(filepath.Dir(a.cwd)); ok {
		return fail(codeError, "%s đã nằm trong dự án %s; không tạo dự án lồng nhau", a.cwd, root)
	}
	if err := config.Write(cfgPath, cfg); err != nil {
		return fail(codeError, "ghi %s: %v", cfgPath, err)
	}
	docs := filepath.Join(a.cwd, cfg.DocsDir)
	var created []string
	for _, d := range docsTree {
		p := filepath.Join(docs, filepath.FromSlash(d))
		if err := os.MkdirAll(p, 0o755); err != nil {
			return fail(codeError, "tạo %s: %v", p, err)
		}
		keep := filepath.Join(p, ".gitkeep")
		if _, err := os.Stat(keep); errors.Is(err, os.ErrNotExist) {
			if err := os.WriteFile(keep, nil, 0o644); err != nil {
				return fail(codeError, "%v", err)
			}
		}
		created = append(created, filepath.ToSlash(filepath.Join(cfg.DocsDir, d)))
	}
	cl := filepath.Join(docs, ChangelogFile)
	if _, err := os.Stat(cl); errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(cl, []byte(changelog.Title+"\n"), 0o644); err != nil {
			return fail(codeError, "%v", err)
		}
	}
	plans := filepath.Join(a.cwd, cfg.PlansDir)
	if _, err := os.Stat(plans); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(plans, 0o755); err != nil {
			return fail(codeError, "%v", err)
		}
		_ = os.WriteFile(filepath.Join(plans, ".gitkeep"), nil, 0o644)
	}

	preCommit, err := a.installPreCommit()
	if err != nil {
		return err
	}

	next := "Tiếp theo: `dk skill install` và `dk hook install`"
	if a.json {
		return a.printJSON(map[string]any{
			"config":     config.FileName,
			"docs_dir":   cfg.DocsDir,
			"changelog":  filepath.ToSlash(filepath.Join(cfg.DocsDir, ChangelogFile)),
			"dirs":       created,
			"pre_commit": preCommit,
			"next":       next,
		})
	}
	a.printf("Đã tạo %s, %s/ (%d thư mục) và %s/%s\n", config.FileName, cfg.DocsDir, len(created), cfg.DocsDir, ChangelogFile)
	switch preCommit.Status {
	case "installed":
		a.printf("Đã cài pre-commit: %s\n", preCommit.Path)
	case "exists":
		a.printf("%s đã có; thêm đoạn sau vào cuối file đó:\n\n%s\n", preCommit.Path, preCommit.Snippet)
	}
	a.printf("%s\n", next)
	return nil
}

// preCommitResult mô tả việc cài pre-commit trong `dk init`.
type preCommitResult struct {
	Status  string `json:"status"` // installed, exists, no-git
	Path    string `json:"path,omitempty"`
	Snippet string `json:"snippet,omitempty"`
}

// installPreCommit chép script nhúng vào .git/hooks/pre-commit khi có repo git
// và chưa có hook; đã có thì trả đoạn cần thêm để người nối tay.
func (a *app) installPreCommit() (preCommitResult, error) {
	if !gitx.IsRepo(a.cwd) {
		return preCommitResult{Status: "no-git"}, nil
	}
	hooks, err := gitx.HooksDir(a.cwd)
	if err != nil {
		return preCommitResult{}, fail(codeError, "%v", err)
	}
	gitRoot, err := gitx.Root(a.cwd)
	if err != nil {
		return preCommitResult{}, fail(codeError, "%v", err)
	}
	// Git chạy hook tại gốc repo; --cwd trỏ về dự án khi dk.toml nằm ở thư mục con.
	rel, err := filepath.Rel(gitRoot, a.cwd)
	if err != nil {
		rel = "."
	}
	raw, err := assets.FS.ReadFile("hooks/pre-commit.sh")
	if err != nil {
		return preCommitResult{}, fail(codeError, "%v", err)
	}
	script := strings.ReplaceAll(string(raw), "__DK_CWD__", filepath.ToSlash(rel))
	path := filepath.Join(hooks, "pre-commit")
	if _, err := os.Stat(path); err == nil {
		lines := strings.Split(strings.TrimSpace(script), "\n")
		return preCommitResult{Status: "exists", Path: path, Snippet: strings.Join(lines[1:], "\n")}, nil
	}
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		return preCommitResult{}, fail(codeError, "%v", err)
	}
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		return preCommitResult{}, fail(codeError, "%v", err)
	}
	return preCommitResult{Status: "installed", Path: path}, nil
}
