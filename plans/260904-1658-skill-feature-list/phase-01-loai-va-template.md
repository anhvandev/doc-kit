---
phase: 1
title: "Loại và template feature-list"
status: completed
priority: P2
effort: "2h"
dependencies: []
---

# Phase 1: Loại và template feature-list

## Overview

Đăng ký loại `feature-list` trong `types.toml`, viết template, sửa chú thích roadmap; `dk new feature-list --from <brief>` chạy được và test đếm template lên 42.

## Requirements

- Functional: `dk new feature-list <slug> --from <brief>` tạo `docs/overview/feature-list.md`, chép `title`, điền `source`; `dk check` không lỗi trên file mới tạo (sau khi điền `owner`); `dk render` ra HTML.
- Non-functional: template dưới 60 dòng, tiếng Việt có dấu, chú thích gợi ý theo kiểu các template khác; không đụng `internal/cli/index.go`.

## Architecture

Không có mã Go mới. `internal/tmpl` tìm template theo `templates/<loại>.*`; `internal/doctype` đọc `types.toml`; `internal/docs/new.go` chép trường theo `from` và tự điền `source` từ đường dẫn `--from`. Thêm một bảng trong `types.toml` và một file template là đủ. Test đếm 41 template ở `internal/cli/cli_test.go:77` đổi thành 42; `phase10_test.go:82` dùng `< 41` giữ nguyên.

## Related Code Files

- Create: `assets/templates/feature-list.md`
- Modify: `assets/types.toml` (thêm `[feature-list]` sau `[glossary]`), `assets/templates/roadmap.md` (chú thích mục 3), `internal/cli/cli_test.go` (41 → 42), `docs/lenh.md` dòng 13 (41 → 42)

## Implementation Steps

1. `assets/types.toml`, sau `[glossary]`:
   ```toml
   [feature-list]
   dir = "overview"
   name = "feature-list.md"
   id = "none"
   description = "Feature list: bảng tính năng dự kiến tách từ Product brief đã duyệt; mã tạm FL-xx, nhóm, ưu tiên MoSCoW, dẫn mục brief; mỗi dòng đi tiếp qua intake tính năng"
   required = ["type", "title", "status", "owner", "created", "source"]
   statuses = ["draft", "review", "approved"]
   warn_lines = 200
   from.brief = { title = "title" }
   ```
2. `assets/templates/feature-list.md`, frontmatter như `product-overview.md` (bỏ `outcome`, `acceptance`), thân:
   ```markdown
   # Feature list: {{.Title}}

   <!-- gợi ý: mục 1 là frontmatter. Bảng tính năng dự kiến tách từ Product brief đã duyệt (source); mỗi dòng là một tính năng sẽ đi qua intake riêng rồi mới có Feature Spec. AI soạn nháp (status: review), người duyệt (approved). Sau khi chốt chỉ sửa theo CR; riêng cột Spec điền mã F-xxx khi spec ra đời, không cần CR. Mã tạm FL-01, FL-02 không đánh số lại; bỏ dòng thì ghi won't, không xóa. -->

   ## 2. Nguồn

   <!-- gợi ý: liên kết brief và idea.md cùng thư mục -->

   - Brief: 

   ## 3. Bảng tính năng

   <!-- gợi ý: Mô tả một câu, người dùng làm được gì. Nhóm khớp mục 4. Ưu tiên: must | should | could | won't. Nguồn: mục brief hoặc idea sinh ra dòng này (brief §1, idea §4). Spec: mã F-xxx khi đã có, để trống khi chưa. -->

   | Mã | Tên | Mô tả | Nhóm | Ưu tiên | Nguồn | Spec |
   |---|---|---|---|---|---|---|
   | FL-01 | | | | | | |

   ## 4. Nhóm

   <!-- gợi ý: mỗi nhóm một dòng, một câu; 3 đến 7 nhóm; Product overview mục 4 dùng cùng tên nhóm -->

   - 

   ## 5. Chưa rõ

   <!-- gợi ý: điều brief hoặc idea nói mà chưa tách được thành tính năng, hoặc tính năng suy ra nhưng brief không nhắc; để người quyết. Không có thì ghi "không". -->

   - 
   ```
3. `assets/templates/roadmap.md` dòng 27: "mã Feature Spec (F-xxx) hoặc tên brief khi chưa có spec" → "mã Feature Spec (F-xxx), hoặc mã tạm FL-xx của Feature list / tên brief khi chưa có spec".
4. `internal/cli/cli_test.go:77` `41` → `42`; `docs/lenh.md:13` `(41)` → `(42)`.
5. Chạy `go test ./internal/cli/ ./internal/doctype/ ./internal/tmpl/`, rồi `go test ./...`, `go vet ./...`.
6. Kiểm tay trong thư mục tạm: `dk init`, tạo brief `--set level=project --set kind=product`, đặt `status: approved`, `dk new feature-list san-pham --from docs/intake/<d>/brief.md --set owner=x`, `dk check docs/overview/feature-list.md`, `dk render docs/overview/feature-list.md`, `dk status`.

## Success Criteria

- [x] `dk template list` liệt kê `feature-list`, tổng 42; `dk template show feature-list` in template.
- [x] `dk new feature-list --from <brief>` ra `docs/overview/feature-list.md`, `title` và `source` đúng.
- [x] `dk check` trên file mới: không lỗi `frontmatter-required`, `status-valid`, `backlink`.
- [x] `go test ./...`, `go vet ./...` qua.

## Risk Assessment

- `dk status` hoặc `dk check` có quy tắc chỉ áp cho thư mục `overview` theo tên loại cứng (ví dụ `line-threshold` theo `warn_lines`): dấu hiệu là test `internal/check` fail hoặc `dk status` không đếm loại mới. Ứng phó: đọc rule liên quan, thêm loại vào bảng nếu là bảng cứng; không đổi thiết kế.
- `from.brief` không chép được gì ngoài `title` (brief không có trường khác cần): chấp nhận, `source` do `--from` điền.
