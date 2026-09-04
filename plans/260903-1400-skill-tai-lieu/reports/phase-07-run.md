# Phase 7: kết quả chạy thử doc-design-system, doc-design-flow

Ngày: 2026-09-03. Dự án mẫu `/tmp/tmp.*` (git init, `dk init`, `dk skill install`, `dk hook install`), binary `dk` build từ `go build ./cmd/dk`. Nội dung component và mockup giữ nguyên khung template (chạy thử CLI, không điền tay).

## Chuỗi lệnh hai SKILL.md quy định, chạy tay

| Bước | Lệnh | Kết quả |
|---|---|---|
| Cài | `dk skill install` | 7 skill "đã cài"; `skill status` 14 dòng |
| Design brief | `dk new brief giao-dien --set level=project --set kind=design --set status=approved --set owner=an` | brief `approved` |
| Tokens | `dk new design-tokens tokens --from <brief> --set owner=an` | `docs/design/tokens/tokens.json`, `$dk` là khóa đầu, `source: 260903-giao-dien/brief.md` |
| CSS | `dk tokens css` | `tokens.css` 57 biến, dòng đầu `/* generated: dk tokens css */`, khối `[data-theme="dark"]` 2 biến (`--color-bg-page`, `--color-text-body`) |
| Foundations | `dk new design-foundations foundations --from <brief> --set owner=an` | `docs/design/foundations.md` 7 mục |
| 4 atoms, 1 molecule, 1 template | `dk new design-component text\|icon\|button\|input --set layer=atom`, `form-field --set layer=molecule`, `list-page --set layer=template` | đúng thư mục `atoms/`, `molecules/`, `templates/`; `--set layer=page` bị từ chối mã 1 với danh sách lớp hợp lệ |
| User flow | `dk new userflow bo-loc --from docs/features/F-001-bo-loc.md --set owner=an` | `flows/F-001-flow.md`, `feature: F-001` chép từ `id` của spec, `source: F-001` |
| Wireframe | `dk new wireframe bo-loc --from <flow> --set step=B1 --set owner=an` | `wireframes/F-001-B1.md`, `source: flows/F-001-flow.md` |
| Mockup HTML | `dk new mockup bo-loc --from <wireframe> --set owner=an --set title="Danh sách"` | `mockups/F-001-B1.html`: chú thích `<!-- dk: -->` 14 trường, `<link href="../tokens/tokens.css">` (file tồn tại, mở từ đĩa), 12 chỗ `var(--...)`, không giá trị gõ tay |
| Mockup Figma | `dk new mockup bo-loc-figma --from <wireframe> --set step=B2 --set external=https://figma.com/file/x` | `F-001-B2.html` chỉ có `<a href>` và `<img src="F-001-B2.png">`, không `tokens.css`, không `<style>` |
| Changelog | `dk changelog add docs/design/<file>` cho 12 file kể cả `.json`, `.html` | mỗi file một dòng `mới, N dòng`; `changelog pending` trống dù `tokens.css` chưa commit (bản sinh bỏ qua) |
| `check` sạch | `dk check` | 0 lỗi, 0 cảnh báo, mã 0 |
| Gõ `#333` và `12px` vào mockup | sửa `<style>` rồi `dk check docs/design/mockups/F-001-B1.html` | 2 lỗi `mockup-tokens` đúng dòng 30, mã 3 |
| User flow có mã ngoài spec | đổi `B3` thành `B9` rồi `dk check` | lỗi `userflow-steps: mã bước không có trong F-001: B9`, mã 3 |
| Alias vòng | `color.action.primary` ↔ `primary-hover` rồi `dk tokens css` | `alias vòng: color.action.primary -> color.action.primary-hover -> color.action.primary`, mã 1, không ghi `tokens.css` |
| HTML | `dk render docs/design/flows/F-001-flow.md`; `dk render --all` | trang flow có dòng "Bước:" liên kết `B1 B2 B3`; `--all` chỉ render 12 file Markdown, bỏ `.html`, `.json` |
| `status` | `dk status` | đếm đủ 8 loại mới theo trạng thái |
| Thiếu `--set` | `dk new mockup x` | `cần --set feature, step=... (tên file dùng {feature}-{step}.html)`, mã 1 |

## Tiêu chí phase 7

| Tiêu chí | Kết quả |
|---|---|
| `dk tokens css` đúng từ tokens.json có alias; alias vòng báo rõ | đạt (57 biến, chuỗi alias in ra) |
| Mockup mở từ đĩa dùng `tokens.css`; `check` bắt hex, px gõ tay | đạt |
| User flow đúng mã bước spec; lệch bị bắt | đạt |
| Bộ tối thiểu tạo đúng thứ tự lớp 0 đến 4, mỗi component đủ 5 mục | template `design-component` có mục 3 biến thể, 4 kích thước, 5 trạng thái, 6 quy tắc, 7 accessibility; thứ tự lớp do SKILL.md bước 2 đến 5 quy định; phần AI điền thật gộp phase 10 |
| AI sinh mockup từ template chỉ bằng organism và semantic tokens, `check` qua | mockup khung template qua `check`; phần AI điền thật gộp phase 10 |
| `--set external=<url>` lưu liên kết và ảnh, không sinh HTML mockup | đạt |

## Lỗi phát hiện và đã sửa khi chạy thử

- Test `tokens` đếm sai số biến (10 thay vì 11); sửa test.
- Mô tả hai skill trùng 3 từ liên tiếp với `doc-overview`, `doc-feature-spec` (`cấp dự án`, `brief đã duyệt`, `của một tính`); đổi chữ mô tả, câu phủ định đổi sang "Không dùng cho" để test bỏ qua đúng chủ ý.
- Hai lỗi `check` trong lần chạy đầu là do kịch bản thử (spec tạo không `--from`, brief có liên kết mẫu `./interview.md`), không phải lỗi phase 7.

## Sau code review

`plans/reports/code-review-260903-phase-07-skill-design.md`: 5 High, 9 Medium. Đã sửa và chạy lại trong dự án mẫu:

- H1: file `.html` trong `docs/` thiếu khối `<!-- dk: -->` giờ là lỗi `mockup-tokens` và vẫn bị lint (`lach.html` với `style='color:#444'` báo 2 lỗi: thiếu khối, `#444`).
- H2: `style='...'` nháy đơn bị lint (cùng file trên).
- H3: px trong prelude `@media (min-width: 768px)` không còn báo (fixture `khong-metadata.html`).
- H4: rules.md `doc-design-flow` mục "Đặt tên mockup": hậu tố trạng thái (`B3e`) chỉ là tên file mockup, không ghi vào cột mã bước của user flow; liên kết ở cột Mockup của dòng `B3`.
- H5: `doc-design-system` bước 6 dùng `--from <Design brief>` để có `source`, giữ `F-000-B0.html` làm mẫu thay vì xóa rồi ghi changelog (`dk check` file mẫu: 0 lỗi).
- M: `--set feature=../x` bị từ chối (chỉ chữ, số, `.`, `-`, `_`); `line-threshold` chỉ áp cho Markdown (tokens.json, mockup không tính dòng); thêm `internal/cli/phase7_test.go` (`dk tokens css`, `changelog add` cho `.html`/`.json`, `pending` bỏ `tokens.css`, `render --all` bỏ file không Markdown); rules.md ghi rõ lint chỉ bắt hex và px, theme phải khai trên chính token.
- Chấp nhận, không sửa: `hook run` vẫn chỉ chặn `.md` (phase 7 không đổi hook; mockup tạo tay bị `mockup-tokens` bắt); `$dk` trùng khóa trong JSON tay; alias đến token có theme không kế thừa theme; `rgb()`, `rem`, `%` không lint; `go.mod` chưa `tidy` (goldmark ở khối indirect, chờ phase 10).
