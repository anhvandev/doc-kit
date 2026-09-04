# Phase 6: kết quả chạy thử doc-feature-spec

Ngày: 2026-09-03. Dự án mẫu `/tmp/dk-p6.*` (git init, `dk init`, `dk skill install doc-feature-spec`), binary `bin/dk` từ `make build`. Phần điền thân spec làm tay theo đúng bảng "Điền thân" của SKILL.md.

## Chuỗi lệnh mà SKILL.md quy định, chạy tay

| Bước | Lệnh | Kết quả |
|---|---|---|
| Cài | `dk skill install doc-feature-spec` | "đã cài"; `skill status` 10 dòng (5 skill x 2 scope) |
| Brief | `dk new brief loc-don --set level=feature --set kind=feature --set status=approved --set owner=an --set outcome=...` | brief `approved` cấp tính năng |
| Spec | `dk new feature-spec loc-don --from <brief> --set owner=an` | `F-001-loc-don.md`, 11 mục đúng thứ tự, `has_ui: true`, `format: spec`; `title`, `purpose`, `acceptance`, `source` chép từ brief; không sửa tay cấu trúc |
| Điền thân | sơ đồ 5 nút B1..B5, bảng hành vi 5 dòng cùng mã, mục 7 có 3 ngoại lệ (3/5 = 60% số dòng bảng hành vi), 2 AC, `status: review` | `dk changelog add --source 260903-loc-don`: `mới, 100 dòng`; `dk check`: 0 lỗi, 1 cảnh báo `spec-has-test` (chưa có test, đúng) |
| HTML | `dk render`; `dk index features` | dòng "Bước:" có 5 liên kết `#step-B1..B5`, 5 ô `id="step-Bn"` trong bảng hành vi; catalog 1 dòng |
| CR chốt | `dk new cr loc-theo-chi-nhanh --set status=approved ...`; bảng tác động dòng Feature Spec `Có`: "B3 đếm theo chi nhánh; chèn B3a" | CR `approved` |
| `update` | sửa spec: chèn `B3a` vào sơ đồ và bảng (không đánh số lại), đổi câu B3, `source` → mã CR, thêm `### Lịch sử nguồn`; `dk changelog add --source <CR>` | `git diff -U0`: chỉ đổi `updated`, `source`, 2 dòng sơ đồ, 2 dòng bảng (B3a, B3) và khối Lịch sử nguồn; changelog `+10 −4` khớp `git diff --numstat`; `dk check` 0 lỗi; HTML "Bước:" thành `B1 B2 B3a B3 B4 B5` |
| Sửa khi CR chưa chốt | CR về `status: review`, `updated: 17:40` (cũ hơn spec 17:56); `dk check` | lỗi `cr-approval-order: sửa lúc 17:56 nhưng CR-... còn review (cập nhật 17:40)`, mã thoát 3; CR về `approved` thì 0 lỗi |
| Biến thể | `--set format=crud --set has_ui=false`, `story`, `use-case`, `state` | tập mục: crud `2 3 5 7 8 9 10 11`; ba biến thể còn lại `2..11`; `dk check` không lỗi mục nào; `--set format=xyz` bị từ chối mã 1 |

## Lỗi phát hiện và đã sửa

- Regex mã bước `^B\d+$` không nhận hậu tố `B2a` (plan mục 7, quy tắc không đánh số lại). Sửa `internal/render/steps.go` thành `^B\d+[a-z]?$`, dùng chung cho `render` và `check`.
- Template test cũ `Render("feature-spec", Data{})` cho `HasUI: false` nên còn 9 mục; test đổi sang `Data{Format: "spec", HasUI: true}`.
- Chạy thử lần đầu `cr-approval-order` không bắt vì `dk changelog add` cho CR chạy **sau** khi sửa spec, bump `updated` của CR lên cùng phút với spec. Đây là thứ tự sai của kịch bản thử, không phải lỗi rule; nhưng cho thấy giới hạn: hai lần sửa cùng phút (`updated` chỉ đến phút) không phân biệt được. Ghi ở plan 1.12.

## Sau code review

`reports/code-review-260903-phase-06-skill-feature-spec.md`: 2 High, 3 Medium, 6 Low. Đã sửa:

- H1, H2 (cùng gốc): `step-codes` chỉ so khi spec có sơ đồ; spec cố ý không sơ đồ (luồng tuyến tính, CRUD) giữ mã bước ở bảng không còn bị báo giả. Mục 6 của biến thể `crud` đổi cột đầu thành "Màn hình", không dùng mã bước. Chạy lại: `--set format=crud` (has_ui mặc định) và spec xóa khối mermaid đều `0 lỗi`. Fixture `F-007-tuyen-tinh.md`.
- M1: `--set has_ui` chỉ nhận `true`, `false`; giá trị khác là lỗi (`has_ui "no" không hợp lệ`).
- M2: `findBehaviorTable` nhận cả tiêu đề chứa "use case" (biến thể `use-case`), không còn phụ thuộc thứ tự bảng.
- M3: `sectionNumbers` xử lý fence ``` và ~~~ theo ký tự và độ dài; fixture F-006 có `## 9.` trong khối `~~~~` không được tính.
- L3 (giới hạn `cr-approval-order` ghi vào SKILL.md), L6 (dấu phẩy README).

Không sửa, ghi ở plan 1.12: L1 (`--set k=true|false` ép bool mọi khóa), L2 (`updated` chưa bắt buộc với feature-spec), L4 (`--source` changelog là tên thư mục intake, khác `source` frontmatter; theo quyết định 1.10), L5 (`format` chỉ validate cho feature-spec).

## Chưa chạy được ở phiên này

Phần AI thật trong agent: skill dừng khi brief `review` hoặc `level: project`; hỏi một câu chọn `format`; soạn mục 4 và 5 cùng lúc; từ chối sửa spec `approved` khi không có CR; khôi phục spec từ git khi `check` báo `cr-approval-order`. Gộp vào đợt chạy thử toàn vòng phase 10. Điều cần quan sát:

- Khi chèn bước, AI có giữ hậu tố `B3a` hay tự đánh số lại `B4..B6` (mockup, test lệch). `step-codes` chỉ bắt lệch sơ đồ và bảng, không bắt đánh số lại; nếu chạy thử cho thấy cần, thêm rule `step-removed` (warning khi mã có ở HEAD biến mất) như plan dự phòng.
- AI có viết ngoại lệ đủ nửa bảng hành vi không; `check` không đo tỷ lệ này.
