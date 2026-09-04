---
title: "Phase 11: target Codex và build Windows cho dk"
date: 2026-09-04
summary: "Thêm target codex (hooks.json, apply_patch), --target danh sách, doctor codex, build 6 nền tảng; review sửa bypass Move to và doctor .codex rỗng"
---

# Phase 11: target Codex và build Windows cho dk

## What happened
- Codex CLI 0.153.2 cài qua npm; tài liệu hooks và `apply_patch.rs` cho thấy `hooks.json` cùng cấu trúc khối `hooks` của Claude, `command` là chuỗi (plan dự kiến mảng), sửa file phát `tool_name: apply_patch` với patch trong `tool_input.command`.
- `internal/target/codex.go` chỉ đổi đường dẫn (`.codex/`, `$CODEX_HOME`) và matcher `apply_patch`; mã đọc/ghi/gỡ hook tách thành `installHooksFile`/`uninstallHooksFile`/`installedHooksFile` trong `claude.go`.
- `hook run` đọc `*** Add File:`, `*** Update File:`, `*** Move to:` từ patch, deny một lần cho mỗi patch.
- `--target claude,codex` cho skill/hook/doctor; `hook install --target codex` in nhắc trust `/hooks`; `doctor` tự kiểm codex khi có `.codex/skills/` hoặc `.codex/hooks.json`.
- Windows: Makefile, goreleaser, CI `windows-latest`; `doctor` bỏ kiểm bit chạy; `gitx` pathspec slash và cắt `\r`; test đặt `USERPROFILE`. Chỉ kiểm bằng `GOOS=windows go vet` và biên dịch test.
- Review bắt bypass `*** Move to:` (đổi tên vào `docs/` không bị chặn) và vòng lặp doctor báo thiếu khi `.codex/` rỗng sau gỡ; đã sửa kèm test.

## Decision
- Giữ `doctor --target lạ` mã 3 kèm bảng (quyết định phase 10) và số hook in cố định (nit 1.9).
- Không thêm trường `shell` cho hook Claude trên Windows: `dk hook run` không phụ thuộc shell.

## Next steps
- `codex login` rồi chạy vòng Tiếp nhận thật, ghi payload thật vào `internal/hook/testdata/`.
- Tạo remote, đổi module path, tag `v0.1.0` rồi `v0.2.0` sau khi CI Windows xanh.

> Historical work record — not durable authority. Prefer docs/specs/ADRs for current decisions.
