---
title: "dk phase 3: cài skill và hook cho Claude Code"
date: 2026-09-03
summary: "skill/hook install-uninstall theo target claude, hook run fail-open, pre-commit trong init; 5 lỗi review đã sửa, 1 chờ quyết định"
---

# dk phase 3: cài skill và hook cho Claude Code

## What happened
- Thêm `internal/target` (interface Target, registry, `claude.go` merge `settings.json` giữ thứ tự khóa bằng object thô), `internal/skill` (list, hash bỏ khối metadata, install/uninstall/status với dấu vết `metadata.dk_*`), `internal/hook` (`Entries`, `Run` fail-open đọc JSON stdin), lệnh `dk skill *`, `dk hook *`, `dk hook run` ẩn, pre-commit trong `dk init`.
- Tách `changelog.Record`/`Tracks`/`Load` khỏi `internal/cli` để hook gọi nội bộ, không `os/exec`.
- Kiểm bằng binary thật trong repo tạm: deny JSON đúng định dạng Claude Code, post-edit ghi dòng "chưa tóm tắt | trực-tiếp", pre-commit chặn khi pending, cho qua khi thiếu `dk`.

## Review (plans/reports/code-review-260903-phase-03-install-claude.md)
- Sửa: `hook uninstall` xóa cả nhóm matcher (mất hook người dùng cùng mục) → lọc từng lệnh; tên skill traversal → chỉ nhận một thành phần đường dẫn; cập nhật không nguyên tử → ghi temp rồi rename; hook bump `updated:` gây "file modified since read" → `Record(..., bump=false)` cho hook; mục giữ chỗ không được thay bởi tóm tắt thật khác nguồn → gộp bất kể nguồn khi mục cũ là "chưa tóm tắt".
- Lần đầu `json.Marshal` escape HTML làm lệnh `>` của người dùng thành `>` và map làm đảo thứ tự khóa; đổi sang Encoder `SetEscapeHTML(false)` và struct có thứ tự.

## Decision
- Ghi ở plan.md mục 1.9. Chưa quyết: hook post-edit làm `changelog pending` luôn xanh với file agent sửa (H2); không khóa file changelog khi hook chạy song song (M6, chấp nhận).

## Next steps
- Bước 9: kiểm thật trong Claude Code (Write mới bị từ chối, Edit có dòng changelog, Edit liên tiếp không báo "modified since read").
- Chủ dự án quyết H2 rồi sang phase 4.

> Historical work record — not durable authority. Prefer docs/specs/ADRs for current decisions.
