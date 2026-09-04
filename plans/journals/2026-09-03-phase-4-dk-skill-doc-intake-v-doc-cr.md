---
title: "Phase 4 dk: skill doc-intake và doc-cr"
date: 2026-09-03
summary: "Hai skill Tầng 1 viết xong, template và types.toml hoàn thiện, chỉ mục intake/cr sửa, test nội dung skill; chưa chạy thật trong Claude Code"
---

# Phase 4 dk: skill doc-intake và doc-cr

## Chuyện gì đã xảy ra

- Viết `assets/skills/doc-intake` và `assets/skills/doc-cr` (SKILL.md 144 và 156 dòng, rules.md chép quy tắc từ báo cáo tư vấn, bảng quy tắc → mục báo cáo). Xóa `doc-smoke`, test cơ chế cài chuyển sang `doc-cr`.
- `types.toml`: `idea` bắt buộc `level`, `brief` bắt buộc `level`, `kind` (thêm `approved_by`, `approved_at`), `cr` bắt buộc `requester`. Khóa mới `beside_source` để `dk new interview --from <CR>` đặt vào `cr/<CR-id>/interview.md`; quyết định 1.7 giữ nguyên cho loại khác.
- `dk index intake` một dòng một thư mục theo trạng thái brief; `dk index <dir>` chỉ liệt kê loại có `dir` trùng thư mục. Lỗi phát hiện khi chạy thử: interview trong thư mục CR bị liệt kê như một CR.
- `internal/skill/content_test.go`: frontmatter hợp lệ, dưới 300 dòng, từ cấm target, mô tả không trùng 3 từ liên tiếp (bỏ câu "Không dùng cho..." vì nó cố ý nhắc phạm vi skill kia).
- Frontmatter SKILL.md không parse vì mô tả có dấu hai chấm chưa quote; test mới bắt được.
- Code review: M1 (CR chờ duyệt chưa có changelog, pre-commit sẽ chặn) → tách khối "Sau mỗi lần sửa CR"; M2 (required mới phá tương thích) → ghi nhận có chủ đích ở plan 1.10; L1 thêm `--set owner`; L2 bỏ `level` mặc định của interview.

## Quyết định

- Test skill đặt ở `internal/skill/` chứ không trong `assets/skills/` vì embed nhúng nguyên thư mục.
- Interview CR dùng `--from` thay `--set for=cr` để có `source` cho `dk refs`.

## Bước tiếp theo

- Chạy thật trong Claude Code với người đối thoại (gộp với bước 9 phase 3 vào phase 10), ghi vào `plans/260903-1400-skill-tai-lieu/reports/phase-04-run.md`.
- Phase 5: `doc-overview`, `doc-adr`.

> Historical work record — not durable authority. Prefer docs/specs/ADRs for current decisions.
