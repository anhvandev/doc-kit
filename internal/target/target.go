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
	// MapHooks đổi entries sang dạng target này ghi (matcher riêng của target);
	// kết quả chỉ để đọc, có thể dùng chung bộ nhớ với entries.
	MapHooks(entries []HookEntry) []HookEntry
	// InstallHooks ghi các hook vào cấu hình: bỏ mọi lệnh dk đang có rồi ghi
	// entries, nên chạy lại sau khi đổi matcher hay lệnh vẫn ra đúng bản hiện tại.
	InstallHooks(global bool, entries []HookEntry) error
	// UninstallHooks gỡ mọi hook do dk cài, giữ nguyên phần còn lại.
	UninstallHooks(global bool) error
	// InstalledHooks trả về các hook dk đang có trong cấu hình (event, matcher, lệnh).
	InstalledHooks(global bool) ([]HookEntry, error)
}
