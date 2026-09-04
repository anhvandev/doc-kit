---
title: "Phase 1 CLI lõi dk: dựng xong, review sửa 11 lỗi"
date: 2026-09-03
summary: "Module Go, 5 lệnh nền, template, changelog; review tìm 3 lỗi cao về tên file có dấu, .gitkeep, tiêu đề có dấu hai chấm"
---

# Phase 1 CLI lõi dk: dựng xong, review sửa 11 lỗi

# Phase 1 CLI lõi `dk`

Dựng module `github.com/vner/dk`: cobra root, embed template và `types.toml`, 5 lệnh `init`, `template`, `new`, `changelog add|pending`, `--version`. 7 package test, `make vet test build` sạch, binary tĩnh.

## Chỗ mất thời gian

- Test end-to-end so `changelog pending` với mốc commit HEAD bị vướng khi commit và dòng changelog cùng phút (changelog chỉ lưu `HH:MM`). Giữ quy tắc "từ phút của HEAD trở đi" vì quy trình thật là add rồi commit ngay; test ép `GIT_COMMITTER_DATE` sang phút sau; ghi giới hạn vào README.
- Kỳ vọng test viết sai hai lần khi sửa theo review (thứ tự assertion sau khi thêm mục gộp; thay chuỗi hai lần thành `SHOP-SHOP-`).

## Review bắt được gì tôi bỏ sót

- `git status --porcelain` mã hóa tên có dấu dạng octal: `pending` không bao giờ sạch với tài liệu tiếng Việt. Cần `-z` và `core.quotePath=false`.
- `init` tạo `.gitkeep` rồi chính `pending` báo 18 file chưa có changelog.
- Render tiêu đề thẳng vào YAML qua text template vỡ khi có dấu hai chấm; trường chuỗi tự do phải đi qua `yaml.Node`.
- `go mod tidy` xóa dependency khai báo sớm; không "khóa phiên bản" bằng cách đó được.

Báo cáo: `plans/reports/code-review-260903-phase-01-cli-core.md`.

> Historical work record — not durable authority. Prefer docs/specs/ADRs for current decisions.
