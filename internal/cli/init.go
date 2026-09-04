package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anhvandev/doc-kit/assets"
	"github.com/anhvandev/doc-kit/internal/agentctx"
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
				return a.writeAgentContext()
			}
			return a.runInit(force)
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "ghi đè dk.toml nếu đã có")
	cmd.Flags().BoolVar(&agentContext, "agent-context", false, "ghi khối ngữ cảnh agent vào CLAUDE.md và AGENTS.md tại thư mục hiện tại, không đụng phần còn lại của file")
	return cmd
}

// agentContextResult là kết quả ghi một file ngữ cảnh agent.
type agentContextResult struct {
	File   string `json:"file"`
	Result string `json:"result"` // created, updated, unchanged
}

// writeAgentContext ghi khối nhúng vào CLAUDE.md và AGENTS.md ở gốc dự án
// (có dk.toml) hoặc thư mục hiện tại khi chưa init, giữa hai dấu mốc; chạy
// lại thay đúng khối cũ, giữ phần còn lại. Kiểm cả hai file trước khi ghi để
// một file hỏng không để file kia đã bị đổi.
func (a *app) writeAgentContext() error {
	content, err := assets.FS.ReadFile("agent-context.md")
	if err != nil {
		return fail(codeError, "%v", err)
	}
	dir := a.cwd
	if root, ok := findProjectRoot(a.cwd); ok {
		dir = root
	}
	for _, name := range agentctx.Files {
		st, err := agentctx.Check(filepath.Join(dir, name), content, Version)
		if err != nil {
			return fail(codeError, "%v", err)
		}
		if st == agentctx.StateBroken {
			return fail(codeError, "%s: khối thiếu mốc đóng hoặc có nhiều hơn một khối; sửa tay rồi chạy lại", filepath.Join(dir, name))
		}
	}
	var results []agentContextResult
	for _, name := range agentctx.Files {
		r, err := agentctx.Write(filepath.Join(dir, name), content, Version)
		if err != nil {
			return fail(codeError, "%v", err)
		}
		results = append(results, agentContextResult{File: name, Result: r})
	}
	if a.json {
		return a.printJSON(results)
	}
	for _, r := range results {
		a.printf("%s: %s\n", filepath.Join(dir, r.File), r.Result)
	}
	return nil
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
