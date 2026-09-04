---
title: "Phase 7 dk: skill Design, dk tokens css, metadata JSON và HTML"
date: 2026-09-03
summary: "Xong phase 7: 8 loại tài liệu họ Design, dk tokens css, rule mockup-tokens và userflow-steps, hai skill; sửa 5 High sau review"
---

# Phase 7 dk: skill Design, dk tokens css, metadata JSON và HTML

## What happened

- Thêm 8 loại `types.toml` họ Design với khóa mới `kind = md | html | json`. Metadata ngoài YAML chỉ hai dạng: JSON khóa `$dk` (đọc ở vị trí bất kỳ, ghi làm khóa đầu, giữ thứ tự khóa còn lại), HTML chú thích `<!-- dk: ... -->` đầu file. `frontmatter.SplitFile/JoinFile` chọn theo đuôi; `docs.Scan` quét thêm `.html`, `.json`; `render --all` lọc bằng `Meta.IsMarkdown`.
- `docs.New` mở rộng placeholder `{feature}`, `{step}`, `{layer}` cho `dir` và `name`, giá trị chép từ `--from` theo bảng `from` rồi `--set` ghi đè.
- `internal/tokens`: parser JSON giữ thứ tự khóa, `$type` kế thừa, alias `{a.b.c}` cả khi lồng trong chuỗi, alias vòng in chuỗi `a -> b -> a`, theme qua `$extensions.dk.theme`. `dk tokens css` ghi `tokens.css` với dòng đầu `/* generated: dk tokens css */` để changelog bỏ qua.
- Rule `mockup-tokens` (hex, px trong `<style>` và `style=`) và `userflow-steps` (tập con mã bước của spec theo `feature`; thay nhánh userflow cũ trong `step-codes`).
- Hai skill `doc-design-system`, `doc-design-flow`; test đếm skill 7, `template list` 17.

## Review và sửa

Code review 5 High: mockup thiếu khối dk lách lint, `style='...'` nháy đơn, px trong `@media` bị báo sai, rules.md ghi hậu tố `B3e` vào cột mã bước làm `userflow-steps` báo lỗi, lệnh mẫu bước 6 của doc-design-system thiếu `source` và ghi changelog cho file đã xóa. Đã sửa hết; thêm chặn `--set feature=../x`, `line-threshold` chỉ áp Markdown, test CLI `phase7_test.go`.

## Decision

Chi tiết ở plan mục 1.13. Chấp nhận: `hook run` chỉ chặn `.md`; lint chỉ hex và px; alias đến token có theme không kế thừa theme.

## Next steps

Phase 8 (`doc-plan-report`, `doc-test`); phần AI điền thật gộp vào chạy thử toàn vòng phase 10. Repo vẫn chưa có commit nào.

> Historical work record — not durable authority. Prefer docs/specs/ADRs for current decisions.
