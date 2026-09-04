package target

import (
	"errors"
	"os"
	"path/filepath"
)

// CodexMatcher là tên tool sửa file của Codex; hook bắt mọi apply_patch và
// `dk hook run` tự lọc theo đường dẫn trong patch.
const CodexMatcher = "apply_patch"

// CodexTrustNote nhắc người dùng trust hook: Codex chỉ chạy hook cấp dự án
// sau khi người dùng duyệt trong phiên Codex.
const CodexTrustNote = "Codex chỉ chạy hook sau khi được trust: mở Codex trong dự án, gõ /hooks và duyệt hook dk."

// Codex là target Codex CLI: skill trong .codex/skills, hook trong
// .codex/hooks.json (cùng cấu trúc khối hooks với Claude, command là chuỗi).
// Scope toàn máy theo $CODEX_HOME, mặc định ~/.codex. Không ghi config.toml.
type Codex struct {
	Root string
}

func (c *Codex) Name() string { return "codex" }

func (c *Codex) base(global bool) (string, error) {
	if global {
		if h := os.Getenv("CODEX_HOME"); h != "" {
			return h, nil
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".codex"), nil
	}
	if c.Root == "" {
		return "", errors.New("chưa có gốc dự án cho scope dự án")
	}
	return filepath.Join(c.Root, ".codex"), nil
}

func (c *Codex) SkillDir(global bool) (string, error) {
	b, err := c.base(global)
	if err != nil {
		return "", err
	}
	return filepath.Join(b, "skills"), nil
}

// HooksPath trả về đường dẫn hooks.json theo scope.
func (c *Codex) HooksPath(global bool) (string, error) {
	b, err := c.base(global)
	if err != nil {
		return "", err
	}
	return filepath.Join(b, "hooks.json"), nil
}

// InstallHooks ghi entries với matcher apply_patch: Codex gọi PreToolUse và
// PostToolUse cho apply_patch (tài liệu hooks Codex, kiểm với 0.153.2).
func (c *Codex) InstallHooks(global bool, entries []HookEntry) error {
	path, err := c.HooksPath(global)
	if err != nil {
		return err
	}
	mapped := make([]HookEntry, len(entries))
	for i, e := range entries {
		e.Matcher = CodexMatcher
		mapped[i] = e
	}
	return installHooksFile(path, mapped)
}

func (c *Codex) UninstallHooks(global bool) error {
	path, err := c.HooksPath(global)
	if err != nil {
		return err
	}
	return uninstallHooksFile(path)
}

func (c *Codex) InstalledHooks(global bool) ([]string, error) {
	path, err := c.HooksPath(global)
	if err != nil {
		return nil, err
	}
	return installedHooksFile(path)
}
