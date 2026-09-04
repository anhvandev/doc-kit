---
title: "Lên plan chi tiết dk: CLI Go nhúng template và skill tài liệu"
date: 2026-09-03
summary: "Chốt Go, template và skill nhúng trong binary, bỏ docs/_templates, Claude trước Codex sau; tách 11 phase"
---

# Lên plan chi tiết dk: CLI Go nhúng template và skill tài liệu

## Bối cảnh

Từ báo cáo tư vấn hệ thống tài liệu 10 tầng và roadmap nháp, người dùng chốt: CLI viết bằng Go, template và skill nhúng trong binary, không có `docs/_templates/`, `dk skill install` mặc định scope dự án với cờ `--global`, Claude Code trước, Codex sau.

## Quyết định

- Template chỉ sống trong binary; đổi template là phát hành phiên bản mới. Chấp nhận mất tùy biến theo dự án.
- Skill chỉ có `SKILL.md` và `references/rules.md`; không mang bản sao template; `dk` lỗi thì skill dừng.
- Lớp bảo đảm changelog: pre-commit là lớp cứng, hook agent fail-open. Không dựa vào hook Codex.
- Giả định `dk init` gọi luôn `skill install` và `hook install` cho Claude; chờ người dùng xác nhận.

## Kết quả

`plans/260903-1400-skill-tai-lieu/`: `plan.md` chỉ mục và kiến trúc dùng chung, 11 file phase. `ak plan validate` qua. `ak plan use` không chạy vì thư mục chưa là git repo.

> Historical work record — not durable authority. Prefer docs/specs/ADRs for current decisions.
