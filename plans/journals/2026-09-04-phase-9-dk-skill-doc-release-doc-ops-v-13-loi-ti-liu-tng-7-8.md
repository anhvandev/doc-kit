---
title: "Phase 9 dk: skill doc-release, doc-ops và 13 loại tài liệu Tầng 7, 8, 9"
date: 2026-09-04
summary: "Thêm release-brief --from spec, release-notes --collect, no-jargon, env-no-secret, index user-guide, status DR; review sửa H1 H2"
---

# Phase 9 dk: skill doc-release, doc-ops và 13 loại tài liệu Tầng 7, 8, 9

## What happened

- `assets/types.toml` thêm 13 loại (41 loại, 11 skill): release-brief, release-notes, user-guide, faq, deployment, environment, runbook, monitoring, postmortem, backup-dr, charter, risk-register, meeting-notes; 13 template mới.
- `docs.New`: nhánh `release-brief --from <Feature Spec>` chép mục 2, 3, cột hành động mục 5, cột mockup mục 6 (liên kết đổi gốc `rebaseLinks`), mục 11; `--collect <phiên bản>` cho release-notes gom brief `ready` chưa `released_in`, ghi `released_in` và một dòng changelog mỗi brief; `--append` mở rộng cho faq; postmortem tính `written_within_48h`.
- `dk check` thêm `no-jargon` (warning, danh sách `[release] jargon` trong `dk.toml`) và `env-no-secret` (lỗi); `dk index user-guide` theo `task`; `dk status` liệt kê `backup-dr` quá 6 tháng chưa diễn tập.
- Skill `doc-release`, `doc-ops`; Tầng 9 chỉ template, không skill.
- Review phát hiện `--collect` không nguyên tử và `--collect --force` làm mất dòng brief cũ trong notes cùng phiên bản; đã sửa (dựng nội dung trước rồi ghi, gom thêm brief có `released_in` bằng phiên bản đang tạo).

## Decision

- Không có loại `user-guide-index`; mục lục là README sinh bởi `dk index user-guide` như `features`, `adr`.
- Mô tả kích hoạt skill đổi câu phủ định về dạng "Không dùng cho ..." để qua test trigram (quy ước 1.10).
- `env-no-secret` bỏ chú thích HTML nhưng quét khối mã; placeholder phải khớp trọn `<...>`.
- Phần người không kỹ thuật đọc Release brief và runbook trên sự cố thật gộp vào chạy thử toàn vòng phase 10.

## Next steps

- Phase 10: chạy thử toàn vòng (kể cả phần AI thật của phase 4 đến 9), Agent context file, goreleaser, `v0.1.0`.
- Chưa commit; repo vẫn chưa có commit đầu.

> Historical work record — not durable authority. Prefer docs/specs/ADRs for current decisions.
