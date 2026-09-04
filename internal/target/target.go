// Package target trừu tượng nơi cài skill và hook của từng harness agent
// (Claude Code, Codex). Nội dung skill không biết target; target chỉ quyết
// định thư mục skill và cách ghi cấu hình hook.
package target

// HookEntry là một hook cần cài: sự kiện, matcher tool và lệnh shell.
type HookEntry struct {
	Event   string // PreToolUse, PostToolUse
	Matcher string // tên tool, ví dụ "Write" hoặc "Edit|Write"
	Command string // lệnh shell, ví dụ "dk hook run pre-write"
}

// Target là một harness agent có thể nhận skill và hook.
type Target interface {
	Name() string
	// SkillDir trả về thư mục chứa skill theo scope dự án hoặc toàn máy.
	SkillDir(global bool) (string, error)
	// InstallHooks thêm các hook vào cấu hình, bỏ qua mục đã có cùng lệnh.
	InstallHooks(global bool, entries []HookEntry) error
	// UninstallHooks gỡ mọi hook do dk cài, giữ nguyên phần còn lại.
	UninstallHooks(global bool) error
	// InstalledHooks trả về các lệnh hook dk đang có trong cấu hình.
	InstalledHooks(global bool) ([]string, error)
}
