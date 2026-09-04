// Package hook định nghĩa hook agent của dk: mục cần cài vào target và
// phần chạy khi harness gọi `dk hook run <sự kiện>`.
package hook

import "github.com/anhvandev/doc-kit/internal/target"

// Sự kiện hook mà `dk hook run` nhận.
const (
	PreWrite = "pre-write"
	PostEdit = "post-edit"
)

// Entries là các hook dk cài vào target.
func Entries() []target.HookEntry {
	return []target.HookEntry{
		{Event: "PreToolUse", Matcher: "Write", Command: target.HookCommandPrefix + " " + PreWrite},
		{Event: "PostToolUse", Matcher: "Edit|Write", Command: target.HookCommandPrefix + " " + PostEdit},
	}
}

// Install ghi hook vào cấu hình của t.
func Install(t target.Target, global bool) error {
	return t.InstallHooks(global, Entries())
}

// Uninstall gỡ hook dk khỏi cấu hình của t.
func Uninstall(t target.Target, global bool) error {
	return t.UninstallHooks(global)
}
