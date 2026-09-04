---
phase: 6
title: "Phase 6: Skill Feature Spec"
status: done
priority: P1
effort: "4d"
dependencies: [4, 5]
---

# Phase 6: Skill Feature Spec

## Overview

Skill trung tâm `doc-feature-spec`: tạo Feature Spec 11 mục từ Intake brief đã duyệt bằng `dk new feature-spec --from`, soạn sơ đồ Mermaid có mã bước và bảng hành vi khớp, cập nhật theo bảng tác động của CR đã duyệt, và chọn định dạng thay thế theo bối cảnh.

## Requirements

- Functional: template `feature-spec` hoàn thiện (đã dựng ở phase 1) cộng biến thể qua `--set format=spec|use-case|story|crud|state`; quy trình `new`, `update`, `html`; `dk check` quy tắc `cr-approval-order` (Feature Spec có `source` là CR chưa `đã chốt` mà `updated` mới hơn CR là lỗi) và `spec-section-order` (11 tiêu đề cấp 2 đúng thứ tự, cho phép bỏ mục 6 khi `has_ui: false` và mục 4 khi `format: crud`).
- Non-functional: như phase 4; Feature Spec vượt 500 dòng cảnh báo (mặc định `check`).

## Architecture

- `types.toml` `feature-spec`: thêm trường `has_ui` (bool), `format`, `brief` hoặc `cr` qua `source`; statuses `draft | chờ duyệt | đã duyệt | đang triển khai | đã phát hành | bãi bỏ`; `from.brief`: title, outcome→mục 2, actors nếu có, acceptance→mục 9, non_goals→mục 11, constraints→mục 8.
- Template: 11 mục theo báo cáo mục 3; mục 4 khối ```mermaid `flowchart TD` mẫu với `B1[...]`, `B2{...}`; mục 5 bảng `Mã | Hành động của tác nhân | Phản hồi quan sát được`; mục 6 bảng `Mã bước | Mockup | Trạng thái hiển thị` chỉ liên kết; mục 7 luồng thay thế ghi `B2a`, `B2b` kiểu hậu tố; mục 9 Given / When / Then đánh số `AC1..`; chú thích gợi ý cho từng mục lấy từ "Nguyên tắc viết" của báo cáo. Biến thể `format`: `use-case` thay mục 4, 5 bằng Use Case Specification và sequence; `story` thay mục 9 bằng Gherkin thuần; `crud` bỏ mục 4, mục 5 thành bảng field và quyền; `state` thêm state diagram sau mục 4.
- Quy trình `new`: chỉ từ brief `đã duyệt` cấp tính năng (hoặc CR đã chốt yêu cầu tính năng mới); chọn `format` theo bảng bối cảnh của báo cáo, hỏi người một câu nếu không rõ; chọn loại sơ đồ theo bản chất (flowchart, sequence, state); luồng tuyến tính dưới 4 bước hoặc CRUD → bỏ sơ đồ; soạn mục 4 và 5 **cùng lúc**, mỗi nút một dòng bảng, cùng mã; mục 6 chỉ liên kết mockup theo mã bước, chưa có mockup thì ghi "chưa có, xem họ Design"; mục 7 ngoại lệ chiếm ít nhất nửa bảng hành vi; mục 9 kiểm chứng được (thông báo cụ thể, dữ liệu đổi, sự kiện phát); `status: chờ duyệt`; `dk check` phải qua trước khi báo người.
- Quy trình `update`: chỉ khi có CR `đã chốt` với dòng "Feature Spec: Có" trong bảng tác động; đọc "Cần làm gì" của dòng đó; sửa đúng mã bước nêu; thêm mã mới theo hậu tố nếu chèn giữa (`B2a`) để không đánh số lại; cập nhật `source` thành mã CR, giữ `source` cũ trong mục "Lịch sử nguồn" ở cuối; `dk check`; `dk changelog add --source <CR>`.
- Quy trình `html`: `dk render <file>` rồi kiểm liên kết mã bước.
- Không sửa Feature Spec `đã duyệt` mà không có CR; không sửa test hay brief thay cho spec.

## Related Code Files

- Modify: `assets/templates/feature-spec.md` (biến thể qua `{{if eq .Format "crud"}}`), `assets/types.toml`, `internal/check/` (thêm `rules_spec.go`, `rules_cr_order.go` và test), `internal/docs/new.go` nếu `--set` cần kiểu bool
- Create: `assets/skills/doc-feature-spec/{SKILL.md,references/rules.md}`
- Create test: fixture spec 11 mục đúng và sai thứ tự; spec `updated` mới hơn CR `chờ duyệt`

## Implementation Steps

1. Hoàn thiện template với 5 biến thể `format`; test render mỗi biến thể có đúng tập tiêu đề.
2. `rules_spec.go` (`spec-section-order`), `rules_cr_order.go`; test.
3. `rules.md` từ báo cáo mục 3, Tầng 4, mục 4 và 7 (dòng Feature Spec), mục 10 (dòng luồng tuyến tính, CRUD, yêu cầu đổi hàng ngày).
4. `SKILL.md`; mô tả: "Feature Spec, nguồn sự thật của một tính năng: sơ đồ luồng có mã bước, bảng hành vi, ngoại lệ, tiêu chí chấp nhận; tạo từ Intake brief đã duyệt, cập nhật theo CR đã chốt. Không dùng cho brief, CR, mockup, test."
5. `skills_test.go` qua; cài; chạy thử: từ brief mẫu ra spec `chờ duyệt`, `dk check` qua, HTML có liên kết mã bước; CR mẫu đã chốt sửa B3 → spec chỉ đổi B3 và `source`; thử sửa spec khi CR còn `chờ duyệt` → `dk check` lỗi `cr-approval-order`.
6. Chạy thử biến thể `crud` và `story`; sửa template nếu tiêu đề lệch.
7. Report chạy thử, chỉnh SKILL.md.

## Success Criteria

- [x] Spec tạo từ brief đã duyệt không cần sửa tay cấu trúc; 11 mục đúng thứ tự; `dk check` qua
- [x] Mã bước sơ đồ và bảng hành vi khớp 100% trong mọi spec chạy thử
- [x] `update` theo CR chỉ đổi mã bước nêu trong bảng tác động; `source` cập nhật; changelog ghi `--source <CR>`
- [x] Sửa spec khi CR chưa chốt bị `dk check` bắt
- [x] 5 biến thể `format` render đúng tập tiêu đề theo test
- [x] Ngoại lệ chiếm ít nhất nửa số dòng bảng hành vi trong spec chạy thử (kiểm bằng mắt, ghi vào report: 3/5 dòng)

## Risk Assessment

- **Template phình vì 5 biến thể trong một file**: giữ điều kiện `{{if}}` chỉ ở mức mục, không ở mức câu; nếu vượt 200 dòng thì tách `feature-spec-<format>.md` và `types.toml` ánh xạ `format` sang template, không đổi lệnh.
- **AI đánh số lại mã bước khi chèn** làm mockup và test lệch: quy tắc hậu tố `B2a` in đậm trong SKILL.md; `dk check` báo mã bước biến mất so với HEAD (thêm quy tắc `step-removed` warning nếu chạy thử cho thấy cần).
- **`cr-approval-order` sai khi múi giờ khác nhau** giữa máy: `updated` ghi giờ máy kèm offset; so sánh sau khi parse có offset.
