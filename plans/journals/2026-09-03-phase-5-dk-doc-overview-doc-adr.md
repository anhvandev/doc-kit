---
title: "Phase 5 dk: doc-overview, doc-adr"
date: 2026-09-03
summary: "4 loại tài liệu mới, quy tắc adr-immutable và glossary-term, chỉ mục phẳng features/adr, 2 skill; review 4 Medium, sửa 3"
---

# Phase 5 dk: doc-overview, doc-adr

## Đã làm
- `assets/types.toml`: `product-overview`, `architecture`, `glossary` (dir `overview`, tên cố định), `adr` (`seq:ADR-{n:04}`, `supersedes`, `superseded_by`); khóa mới `warn_lines` theo loại (150/400/300).
- `internal/check`: `adr-immutable` (so thân với HEAD qua `gitx.HeadFile`, chỉ khi status ở HEAD đã chốt; chuẩn hóa CRLF, khoảng trắng), `glossary-term` (in đậm ở mục 2, 5, 8 của Feature Spec, gộp mọi file Glossary, bỏ khối mã).
- `dk index features|adr`: bảng phẳng theo mã thay cho nhóm theo trạng thái (thay đổi có chủ đích so với quyết định 1.8).
- Skill `doc-overview`, `doc-adr` (SKILL.md + rules.md); `content_test` qua 4 skill.
- Chạy thử tay ở `/tmp`: đủ 5 tiêu chí. Report `plans/260903-1400-skill-tai-lieu/reports/phase-05-run.md`.

## Lỗi gặp
- Template có liên kết sẵn đến file chưa tồn tại → `link-broken` ngay sau `dk new`. Bỏ liên kết, giữ gợi ý trong chú thích.
- Lần chạy thử đầu, pre-commit chặn commit ADR vì thiếu changelog nên `adr-immutable` không có HEAD để so; kịch bản phải `dk changelog add` trước.
- Test cài skill đếm cứng 2 skill và `rows[0]` là `doc-cr`; đổi sang tra theo tên.

## Review
4 Medium: gộp mọi Glossary (M1), thứ tự finding ổn định (M2), `HeadFile` còn 1 tiến trình git mỗi ADR (M4) đã sửa. M3: `adr-immutable` chỉ so với HEAD, thân đã commit thì thoát; pre-commit chỉ chạy `changelog pending`. Chờ chủ dự án chọn thêm `dk check` vào pre-commit hay chấp nhận. Ghi ở plan mục 1.11.

## Tiếp theo
Phase 6 `doc-feature-spec`. Phần AI thật của phase 3, 4, 5 gộp vào phase 10.

> Historical work record — not durable authority. Prefer docs/specs/ADRs for current decisions.
