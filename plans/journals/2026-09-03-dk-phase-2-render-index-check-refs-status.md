---
title: "dk phase 2: render, index, check, refs, status"
date: 2026-09-03
summary: "5 lệnh đọc của dk xong: HTML tự chứa nhúng mermaid 11.17.2, 7 quy tắc check, đồ thị refs, chỉ mục sinh ra; sửa 4 should-fix sau review"
---

# dk phase 2: render, index, check, refs, status

## Chuyện gì đã xảy ra
- Thêm `dk render|index|check|refs|status` theo `plans/260903-1400-skill-tai-lieu/phase-02-render-check.md`. Package mới: `internal/render` (goldmark, renderer mermaid, mã bước, chỉ mục HTML), `internal/refs`, `internal/check` (7 quy tắc, mỗi quy tắc một file), `internal/docs/scan.go` dùng chung.
- Khung HTML `assets/html/{page.html,style.css,html-style.md}`; `mermaid.min.js` 11.17.2 bản IIFE, tải một lần, ghi `THIRD_PARTY.md` và `MERMAID-LICENSE.txt`.
- Kiểm offline bằng Chrome headless `--dump-dom` từ `file://`: có `<svg id="mermaid-…">`. 200 feature-spec: `render --all --index` 2 giây, `docs/html/` 685 MB.
- Review (`plans/reports/code-review-260903-phase-02-render-check.md`): DONE_WITH_CONCERNS. Đã sửa: `check <file>` không có trong tập quét trả mã 1 thay vì báo sạch; bảng hành vi ưu tiên bảng dưới tiêu đề "hành vi"; `render <file>` chỉ nhận `.md`; kèm giấy phép MIT vào HTML; `line-threshold` chỉ áp `docs/`; `backlink` đếm cả loại chưa có trong `types.toml`; cỡ chữ `.status` trong metadata.

## Điều bất ngờ
- `html/template` bỏ mọi chú thích HTML, nên giấy phép đặt trong `<!-- -->` biến mất; chuyển thành chú thích JS `/*! … */` ở đầu `<script>` mermaid.
- `dk check` đo 7 giây tưởng chậm; profile chỉ 30 ms CPU. Thủ phạm là đĩa đang ghi 685 MB HTML vừa render. Đo lại: 20 ms.
- Cỡ HTML 3,5 MB mỗi trang có sơ đồ đúng như plan chấp nhận; `--shared-assets` vẫn để chờ có yêu cầu.

## Quyết định
Ghi ở `plan.md` mục 1.8: `docs.Scan` chung, liên kết `.md` ngoài `docs/` trỏ file gốc, trường `final` trong `types.toml`, `check` quét `plans/` nhưng không áp ngưỡng dòng, `index` không ghi `updated`.

## Tiếp theo
- Phase 3: `skill install/uninstall`, `hook install/uninstall`, target `claude`.
- Nit còn lại từ review: tiền tố `..` trong `filepath.Rel`, test `internal/config`, lỗi `Rel` bị nuốt trong `index`.
- Chưa commit: repo chưa có commit nào, người dùng quyết định.

> Historical work record — not durable authority. Prefer docs/specs/ADRs for current decisions.
