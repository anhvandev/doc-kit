---
phase: 2
title: "Phase 2: Render và kiểm tra"
status: done
priority: P1
effort: "5d"
dependencies: [1]
---

# Phase 2: Render và kiểm tra

## Overview

Thêm 5 lệnh đọc: `render` (Markdown sang HTML tự chứa), `index` (chỉ mục từ frontmatter), `check` (kiểm tra cấu trúc và liên kết), `refs` (đồ thị liên kết), `status` (bảng tổng). Sau phase này tài liệu đọc được từ đĩa không cần mạng và lệch mã bước bị bắt tự động.

## Requirements

- Functional: `dk render <path>|--all [--index]`; `dk index [features|adr|cr|intake|all]`; `dk check [<path>] [--strict]`; `dk refs <file>`; `dk status`; tất cả có `--json`; `check` mã thoát 3 khi có lỗi.
- Non-functional: HTML một file, không tải tài nguyên ngoài; mở bằng `file://` hiện sơ đồ Mermaid; tuân `assets/html/html-style.md`; `render --all` cho 200 file dưới 5 giây trên máy phát triển; không đụng file nguồn Markdown.

## Architecture

- `internal/render`:
  - goldmark với `extension.GFM`, `extension.Table`, `html.WithUnsafe()` tắt (escape HTML thô), renderer tùy biến cho fenced code `mermaid`: xuất `<pre class="mermaid">...</pre>`.
  - Khung `assets/html/page.html` (`html/template`): `<head>` có `<style>` nhúng từ `assets/html/style.css`, khối metadata sinh từ frontmatter (bảng 2 cột), thân, cuối trang chèn `<script>` mermaid.min.js **chỉ khi** trang có khối mermaid, gọi `mermaid.initialize({startOnLoad:true, theme:'neutral'})`.
  - Mã bước: sau khi render, regexp `\b(B\d+)\b` trong `<pre class="mermaid">` để lập danh sách mã; trong bảng hành vi (bảng đầu tiên sau tiêu đề chứa "Bảng hành vi" hoặc cột đầu khớp `^B\d+$`) gắn `id="step-B1"` cho ô cột đầu; thêm `<a href="#step-B1">` vào chú thích bên dưới sơ đồ liệt kê mã bước (Mermaid text không chèn được thẻ a, nên liên kết đặt ở dòng "Bước: B1 B2 ..." ngay dưới sơ đồ).
  - Đầu ra: `docs/html/<đường dẫn tương đối>.html`, giữ cây; liên kết `.md` nội bộ đổi thành `.html` tương đối.
  - `--index`: trang `docs/html/index.html` từ frontmatter mọi file: nhóm theo `type`, cột id, title, status, owner, updated, liên kết.
- `assets/html/html-style.md`: quy tắc viết CSS và khung, độc lập, không mượn: cấm font hệ thống mặc định không khai báo, cấm gradient chữ, cấm emoji trong tiêu đề, bố cục bất đối xứng có chủ đích, serif cho tiêu đề, mono cho mã bước, tương phản 4.5:1, `prefers-reduced-motion`, in được. `style.css` viết theo file này; test kiểm tra CSS không chứa `linear-gradient` trong khai báo `color`/`background-clip: text`.
- `internal/refs`: quét mọi `.md` trong `docs/` và `plans/`, lấy liên kết Markdown `[..](path)` tương đối và mã tài liệu trong frontmatter `source:` và trong thân (regexp id theo `types.toml`: `F-\d{3}`, `ADR-\d{4}`, `CR-\d{6}-[a-z0-9-]+`); xây `Graph{Out map[file][]file, In map[file][]file}`; giải mã thành đường dẫn bằng quét frontmatter `id`.
- `internal/check`: bộ quy tắc, mỗi quy tắc một hàm `func(ctx) []Finding` với `Finding{File, Line, Rule, Msg, Level}`:
  1. `frontmatter-required`: thiếu trường theo `types.toml`; `created_by` không phải `dk` là warning (`--strict` nâng thành lỗi).
  2. `link-broken`: liên kết tương đối tới file không tồn tại; liên kết `.md` trỏ vào `docs/html/` là lỗi.
  3. `step-codes`: trong feature-spec, tập mã bước trong Mermaid phải bằng tập mã bước cột đầu bảng hành vi; báo thiếu và thừa từng mã. Với user flow (`type: userflow`) cùng `source`/feature: mã bước phải là tập con của feature-spec cùng mã.
  4. `backlink`: feature-spec, mockup, adr, userflow phải có `source` trỏ đến brief hoặc cr tồn tại; brief và cr trạng thái đã chốt mà không có tài liệu nào `source` trỏ về là warning "tài liệu chết".
  5. `spec-has-test`: feature-spec `status` không còn `draft` phải có file trong `docs/test/` có `source` trỏ về nó; warning.
  6. `line-threshold`: trên 500 dòng warning, trên 800 lỗi; ngưỡng đọc từ `dk.toml` `[check] warn_lines=500 max_lines=800`.
  7. `status-valid`: `status` thuộc `statuses` của loại.
- `internal/cli/index.go`: sinh `docs/features/README.md`, `docs/adr/README.md`, `docs/cr/README.md`, `docs/intake/README.md` với frontmatter `generated: true`, bảng theo trạng thái; CR sắp theo trạng thái rồi ngày. File có `generated: true` không vào changelog.
- `internal/cli/status.go`: bảng số tài liệu theo `type` × `status`, CR đang mở, số file `changelog pending`, số finding của `check` theo mức.

## Related Code Files

- Create: `assets/html/{page.html,style.css,html-style.md,mermaid.min.js}` (mermaid lấy đúng một phiên bản, ghi phiên bản và giấy phép MIT vào `assets/html/THIRD_PARTY.md`), cập nhật `assets/embed.go`
- Create: `internal/render/{render,mermaid,steps,index}.go`, `internal/refs/graph.go`, `internal/check/{check,rules_frontmatter,rules_links,rules_steps,rules_backlink,rules_threshold}.go`, `internal/cli/{render,index,check,refs,status}.go`
- Create test: `internal/render/render_test.go` (golden HTML), `internal/refs/graph_test.go`, `internal/check/*_test.go` với fixture trong `internal/check/testdata/`
- Modify: `internal/cli/changelog.go` (bỏ qua `generated: true`), `internal/config/config.go` (`[check]`)

## Implementation Steps

1. Tải `mermaid.min.js` một phiên bản cố định vào `assets/html/`, ghi phiên bản, giấy phép. Đây là lần duy nhất chạm mạng, ở máy phát triển, không phải lúc chạy.
2. Viết `html-style.md` trước, rồi `style.css` và `page.html` theo nó. Kiểm tra bằng mắt với một Feature Spec mẫu có Mermaid, bảng, chú thích.
3. Viết `internal/render/render.go`: đọc file, tách frontmatter, goldmark render, chèn vào khung, đổi liên kết `.md` sang `.html`, ghi ra `docs/html/`. Test golden: input fixture, output so với file `.golden`; cập nhật golden bằng cờ `-update`.
4. Viết `mermaid.go`: renderer cho fenced `mermaid`; cờ `hasMermaid` để quyết định nhúng script. Test: file không có mermaid thì HTML không chứa chuỗi `mermaid.initialize`.
5. Viết `steps.go`: trích mã bước trong sơ đồ, gắn `id` vào bảng hành vi, chèn dòng "Bước:" có liên kết. Test: 3 mã trong sơ đồ ra 3 liên kết, bấm tới đúng `id`.
6. Viết `index.go` cho `--index` và `internal/cli/index.go` cho chỉ mục Markdown; dùng chung hàm quét frontmatter `docs.Scan(root) []DocMeta`. Test: 2 feature-spec, 1 adr ra bảng đúng số dòng và thứ tự.
7. Viết `internal/refs/graph.go` và `dk refs`. Test: A liên kết B, B `source: A` ra `Out[A]={B}`, `In[A]={B}`.
8. Viết `internal/check` theo 7 quy tắc, mỗi quy tắc một file và một test với fixture có lỗi cố ý; `dk check` gom, in theo file, mã thoát 3 khi có lỗi mức error. `--json` in mảng Finding.
9. Viết `dk status`. Test: đếm khớp fixture.
10. Sửa `changelog add` bỏ qua `generated: true` và mọi thứ trong `docs/html/`.
11. Cập nhật `README.md`: 5 lệnh mới, ảnh chụp HTML không cần, chỉ mô tả.
12. `make vet test`; chạy `dk render --all --index` trên `docs/` mẫu 20 file, mở `docs/html/index.html` bằng trình duyệt ở chế độ offline, xác nhận sơ đồ hiện.

## Success Criteria

- [x] `dk render docs/features/F-001-*.md` ra HTML một file; mở `file://` offline hiện Mermaid; không có request mạng trong DevTools (kiểm bằng Chrome headless `--dump-dom` từ `file://`: có `<svg id="mermaid-…">`; `phase2_test.go` chặn `<link>` và `src="http`)
- [x] File không có Mermaid không chứa script mermaid
- [x] Mã bước trong sơ đồ có liên kết nhảy đúng dòng bảng hành vi
- [x] `dk check` bắt được fixture: mã bước lệch, liên kết hỏng, thiếu `source`, vượt 800 dòng; mã thoát 3
- [x] `dk index all` sinh 4 file chỉ mục có `generated: true`; `dk changelog pending` không liệt kê chúng
- [x] `dk refs` in đúng liên kết đến và đi
- [x] `dk status` khớp số liệu fixture; `--json` hợp lệ
- [x] `style.css` không vi phạm `html-style.md` theo test tự động

## Kết quả

- Hoàn thành 2026-09-03 16:25. `make vet test` qua 10 package. `render --all --index` 200 feature-spec: 2 giây, `docs/html/` 685 MB (3,5 MB mỗi trang có sơ đồ, theo quyết định nhúng inline).
- Review: `plans/reports/code-review-260903-phase-02-render-check.md` (DONE_WITH_CONCERNS; 4 should-fix và 3 nit đã sửa, nit còn lại ghi ở plan.md mục 1.8).
- Chưa làm: mở bằng trình duyệt có giao diện để nhìn `:target` tô dòng (CSS chuẩn, chỉ kiểm DOM); test cho `internal/config`.

## Risk Assessment

- **Mermaid không chạy từ `file://`** ở một số trình duyệt vì module ESM: dùng bản UMD `mermaid.min.js`, không dùng `import`. Tín hiệu vỡ: sơ đồ trắng khi mở offline; phản ứng: đổi sang bản UMD hoặc render sẵn SVG lúc `render` bằng thư viện Go nếu có (chưa có thư viện Go ổn định, nên ưu tiên UMD).
- **Kích thước HTML** 2 đến 3 MB mỗi trang có sơ đồ: chấp nhận theo hợp đồng. Nếu người dùng phàn nàn, thêm `--shared-assets` ghi một `mermaid.min.js` chung trong `docs/html/_assets/`; không làm trước khi có yêu cầu.
- **Quy tắc `backlink` báo giả** với tài liệu cấp dự án (Product overview từ Product brief): cho phép `source` trỏ tới brief loại project; ghi rõ trong `types.toml` loại nào bắt buộc `source`.
- **Regexp mã bước bắt nhầm** chữ "B1" trong văn xuôi: chỉ quét trong khối mermaid và cột đầu bảng hành vi, không quét toàn thân.
