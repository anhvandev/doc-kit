package hook

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anhvandev/doc-kit/internal/changelog"
	"github.com/anhvandev/doc-kit/internal/config"
)

// DenyReason là lý do từ chối khi agent tạo tài liệu mới không qua `dk new`.
const DenyReason = "Tạo tài liệu mới bằng `dk new <type> <slug>` rồi sửa file đó."

// DirectSource là nguồn ghi vào changelog khi hook ghi thay skill.
const DirectSource = "trực-tiếp"

// payload là phần dk cần trong JSON stdin của harness. Claude Code gửi
// tool_input.file_path; Codex gửi apply_patch với toàn bộ patch trong
// tool_input.command (kiểm với Codex 0.153.2).
type payload struct {
	ToolName  string `json:"tool_name"`
	ToolInput struct {
		FilePath string `json:"file_path"`
		Path     string `json:"path"`
		Command  string `json:"command"`
	} `json:"tool_input"`
	CWD string `json:"cwd"`
}

// files trả về các đường dẫn tool sắp ghi hoặc vừa ghi.
func (p payload) files() []string {
	switch {
	case p.ToolInput.FilePath != "":
		return []string{p.ToolInput.FilePath}
	case p.ToolInput.Path != "":
		return []string{p.ToolInput.Path}
	case p.ToolName == "apply_patch" || strings.Contains(p.ToolInput.Command, "*** Begin Patch"):
		return patchFiles(p.ToolInput.Command)
	}
	return nil
}

// patchFiles đọc đường dẫn sau `*** Add File:`, `*** Update File:` và
// `*** Move to:` (đích khi đổi tên) trong patch dạng apply_patch;
// `*** Delete File:` bỏ qua vì không tạo tài liệu.
func patchFiles(patch string) []string {
	var out []string
	for _, ln := range strings.Split(patch, "\n") {
		ln = strings.TrimRight(ln, "\r")
		for _, prefix := range []string{"*** Add File: ", "*** Update File: ", "*** Move to: "} {
			if strings.HasPrefix(ln, prefix) {
				if f := strings.TrimSpace(strings.TrimPrefix(ln, prefix)); f != "" {
					out = append(out, f)
				}
			}
		}
	}
	return out
}

// Run xử lý một sự kiện hook. Luôn trả 0: lỗi nội bộ in cảnh báo ra stderr
// (fail-open); quyết định từ chối in JSON ra stdout theo định dạng Claude Code.
func Run(event string, stdin io.Reader, stdout, stderr io.Writer) int {
	if err := run(event, stdin, stdout); err != nil {
		fmt.Fprintf(stderr, "dk hook %s: %v\n", event, err)
	}
	return 0
}

func run(event string, stdin io.Reader, stdout io.Writer) error {
	var p payload
	if err := json.NewDecoder(stdin).Decode(&p); err != nil {
		return fmt.Errorf("đọc stdin: %w", err)
	}
	denied := false
	for _, f := range p.files() {
		deny, err := handle(event, f, p.CWD)
		if err != nil {
			return err
		}
		if deny && !denied {
			denied = true
			if err := writeDeny(stdout); err != nil {
				return err
			}
		}
	}
	return nil
}

// handle xử lý một file; trả deny=true khi pre-write gặp tài liệu mới.
func handle(event, file, cwd string) (bool, error) {
	if !filepath.IsAbs(file) {
		file = filepath.Join(cwd, file)
	}
	file = filepath.Clean(file)
	root, ok := findRoot(filepath.Dir(file))
	if !ok {
		return false, nil
	}
	cfg, err := config.Load(filepath.Join(root, config.FileName))
	if err != nil {
		return false, err
	}
	docsDir := filepath.Join(root, filepath.FromSlash(cfg.DocsDir))
	rel, err := filepath.Rel(docsDir, file)
	if err != nil || strings.HasPrefix(rel, "..") || !strings.EqualFold(filepath.Ext(rel), ".md") {
		return false, nil
	}
	rel = filepath.ToSlash(rel)
	if !changelog.Tracks(docsDir, rel) {
		return false, nil
	}
	switch event {
	case PreWrite:
		_, err := os.Stat(file)
		return errors.Is(err, os.ErrNotExist), nil
	case PostEdit:
		_, err := changelog.Record(root, docsDir, rel, changelog.NoSummary, DirectSource, time.Now(), false)
		return false, err
	}
	return false, fmt.Errorf("sự kiện lạ %q", event)
}

func writeDeny(w io.Writer) error {
	out := map[string]any{"hookSpecificOutput": map[string]string{
		"hookEventName":            "PreToolUse",
		"permissionDecision":       "deny",
		"permissionDecisionReason": DenyReason,
	}}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(out)
}

// findRoot đi lên từ dir đến thư mục có dk.toml.
func findRoot(dir string) (string, bool) {
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
