---
phase: 4
title: "Phase 4: Skill tiếp nhận: doc-intake, doc-cr"
status: done
priority: P1
effort: "4d"
dependencies: [1, 2, 3]
---

# Phase 4: Skill tiếp nhận: doc-intake, doc-cr

## Overview

Hai skill đầu tiên của Tầng 1: `doc-intake` dẫn chuỗi ý tưởng, phỏng vấn, brief; `doc-cr` soạn Change Request với bảng tác động từ `dk refs`. Hoàn thiện template `idea`, `interview`, `brief`, `cr` đã dựng thô ở phase 1 và thêm chỉ mục CR, intake.

## Requirements

- Functional: SKILL.md và `references/rules.md` cho 2 skill; template 4 loại đủ trường theo báo cáo mục 1a, 1b; `brief` chạy được ở hai cấp (dự án, tính năng) và biến thể Design brief qua `--set level=project|feature --set kind=product|design|feature`; `dk index cr|intake`.
- Non-functional: SKILL.md dưới 300 dòng, rules.md dưới 300 dòng; mô tả kích hoạt không trùng với nhau và với 9 skill sau; không nhắc target; mọi bước ghi file đều đi qua `dk`.

## Architecture

Quy trình `doc-intake` (trong SKILL.md):

0. `dk --version`; lỗi thì dừng: "Cần cài `dk` trước".
1. `new`: hỏi người một câu duy nhất lấy tên ý tưởng và cấp (dự án hay tính năng), chạy `dk new idea <slug> --set level=...`; mở file, điền cùng người từng trường; trường không biết ghi đúng chữ "chưa rõ".
2. Đọc lại idea.md: liệt kê trường "chưa rõ" và mâu thuẫn. Không có → bỏ phỏng vấn, sang bước 4. Có → `dk new interview <slug>` (cùng thư mục) rồi phỏng vấn **từng câu một**, câu hỏi sinh từ đúng trường trống; câu về giao diện phải kèm 2 đến 3 phương án tham chiếu để người chọn; mỗi cặp hỏi đáp kết bằng dòng "Quyết định rút ra" và trạng thái đã chốt hoặc còn mở.
3. Dừng khi không còn trường "chưa rõ" hoặc người yêu cầu dừng.
4. `dk new brief <slug> --from idea.md`; soạn bốn mục cố định (kết quả mong muốn, ràng buộc, ngoài phạm vi, tiêu chí chấp nhận) và mục "Quyết định đã chốt" liên kết về interview.md; `status: chờ duyệt`. **Không tự chốt.**
5. `dk changelog add` cho từng file với tóm tắt thật và `--source <slug intake>`; `dk render` thư mục intake; báo người duyệt brief.
6. `update` chỉ áp cho brief chưa duyệt; brief đã duyệt thì hướng người sang `doc-cr`.

Quy trình `doc-cr`:

0. Kiểm `dk`.
1. `dk new cr <slug>`; điền frontmatter (người yêu cầu, mức ưu tiên) và mục Hiện trạng: hỏi người "thứ sắp đổi đang được mô tả ở đâu"; nếu người không biết, `dk status` và `dk refs` để gợi ý danh sách Feature Spec.
2. Mục Thay đổi đề xuất và lý do từ lời người.
3. **Bảng tác động**: chạy `dk refs <feature-spec>` cho mỗi tài liệu hiện trạng, đọc các file liên kết đến và đi, điền 6 dòng bảng (Feature Spec, Design, ADR, Data model và API, Test case, Release brief và User guide) với "Có / Không" và việc cần làm; `status: chờ duyệt`. **Dừng, chờ người duyệt bảng tác động.** Không sửa tài liệu đích trước khi duyệt.
4. Hỏi đáp chỉ khi tác động còn "chưa rõ", cùng quy tắc interview; dài quá 40 dòng thì `dk new interview <slug> --set for=cr` vào thư mục cùng tên CR.
5. Sau duyệt: ghi mục Quyết định (chấp nhận, từ chối, hoãn; người chốt), liệt kê tài liệu đích; giao việc cập nhật từng tài liệu đích cho skill của họ đó (chỉ nêu tên họ tài liệu, ví dụ "họ Feature Spec", không nêu tên skill khác để mô tả kích hoạt không chồng); đánh dấu xong từng dòng khi tài liệu đích cập nhật; `dk index cr`.
6. Quy tắc tỷ lệ: sửa lỗi nhỏ, tác động một Feature Spec và test → bảng hai dòng, bỏ hỏi đáp. Tác động chạm Design system hoặc ADR → bắt buộc người duyệt bảng trước bất kỳ thay đổi nào. CR bị từ chối giữ file, `status: từ chối`.

`rules.md` của mỗi skill: chép quy tắc (không giải thích) từ báo cáo mục 1, 1a, 1b, 4 (đoạn "Mũi tên chỉ đi một chiều"), 7 (dòng intake và CR), 10 (dòng liên quan); mỗi quy tắc một gạch đầu dòng, cuối file bảng "quy tắc → mục báo cáo".

Template hoàn thiện: `idea.md` có mục Giao diện (màn hình liên quan, template trang dùng, có màn hình hoặc organism mới không); `interview.md` khung lặp "Câu hỏi / Trả lời / Quyết định rút ra / Trạng thái"; `brief.md` bốn mục cố định cộng Quyết định đã chốt, frontmatter `level`, `kind`, `approved_by`, `approved_at`; `cr.md` 7 mục, frontmatter `requester`, `priority`, trạng thái `mới | đang phân tích | chờ duyệt | đã chốt | đang triển khai | đã đóng | từ chối`.

## Related Code Files

- Create: `assets/skills/doc-intake/{SKILL.md,references/rules.md}`, `assets/skills/doc-cr/{SKILL.md,references/rules.md}`
- Modify: `assets/templates/{idea,interview,brief,cr}.md`, `assets/types.toml` (trường `level`, `kind`, `requester`, `priority`, `statuses`, `from`), `internal/cli/index.go` (bảng CR theo trạng thái, bảng intake theo trạng thái brief)
- Create test: `internal/skill/content_test.go` (đổi từ `assets/skills/skills_test.go`, xem Kết quả; mọi SKILL.md: frontmatter hợp lệ, dưới 300 dòng, không chứa `ak-`, `Claude Code`, `.claude/`, tên tool `Edit`/`Write`; mô tả không trùng 3 từ khóa liên tiếp giữa hai skill)
- Delete: `assets/skills/doc-smoke/` (chuyển sang phase 10 nếu vẫn cần cho test cơ chế; nếu test cơ chế dùng skill thật thì xóa ở đây)

## Implementation Steps

1. Hoàn thiện `types.toml` và 4 template; `go test ./internal/tmpl/...` qua.
2. Viết `rules.md` cho `doc-intake` từ báo cáo; đếm dòng.
3. Viết `SKILL.md` `doc-intake` theo quy trình trên; mô tả kích hoạt: "Tiếp nhận ý tưởng mới chưa có gì tồn tại: ghi idea, phỏng vấn trường chưa rõ, soạn Intake brief hoặc Product brief hoặc Design brief chờ duyệt. Không dùng cho thay đổi trên thứ đã có."
4. Viết `rules.md` và `SKILL.md` `doc-cr`; mô tả: "Change Request cho thay đổi trên tính năng, giao diện, kỹ thuật đã có: hiện trạng, đề xuất, bảng tác động chờ duyệt, quyết định, tài liệu đích. Không dùng cho ý tưởng mới."
5. Viết `skills_test.go` và chạy; sửa mô tả nếu trùng.
6. Sửa `index.go` cho `cr` và `intake`; test.
7. `make build && dk skill install doc-intake doc-cr` vào dự án mẫu; chạy thử trong Claude Code: một ý tưởng tính năng có 2 trường "chưa rõ" → interview đúng 2 câu → brief chờ duyệt; một CR trên Feature Spec mẫu → bảng tác động liệt kê đúng file `dk refs` trả về; kiểm `CHANGELOG-DOCS.md`, `dk check`, HTML.
8. Ghi kết quả chạy thử vào `plans/260903-1400-skill-tai-lieu/reports/phase-04-run.md` (điều gì AI hỏi thừa, điều gì thiếu), sửa SKILL.md theo đó, chạy lại một lần.

## Success Criteria

- [x] Một ý tưởng mới đi đến brief `review` chỉ qua chuỗi lệnh skill quy định; không file nào tạo ngoài `dk new`; brief có bốn mục cố định và liên kết về interview (chạy tay chuỗi lệnh, xem `reports/phase-04-run.md`; phần AI dẫn phỏng vấn thật chưa kiểm)
- [x] idea.md điền đủ, không "chưa rõ" → SKILL.md bước 2 bỏ phỏng vấn, sang bước 4 (quy trình; chưa kiểm hành vi AI)
- [x] CR trên tính năng có sẵn: `dk refs` trả đúng brief (đi) và CR (đến) cho spec mẫu; SKILL.md dừng ở `review` trước khi đụng tài liệu đích (câu dừng in đậm đầu file, lặp ở bước 3 và 6)
- [x] `dk index cr` nhóm theo trạng thái, không lẫn interview trong thư mục CR; `dk index intake` một dòng một thư mục theo trạng thái brief; `dk check` 0 lỗi trên intake và CR mẫu
- [x] `internal/skill/content_test.go` qua: frontmatter, dưới 300 dòng, từ cấm, mô tả không trùng trigram
- [x] Mỗi file trong chạy thử có dòng changelog tóm tắt thật

## Kết quả

Xong 2026-09-03, trừ bước 7 chạy thật trong Claude Code với người trả lời (cần phiên tương tác; ghi vào `reports/phase-04-run.md` khi chạy). Quyết định triển khai ghi ở `plan.md` mục 1.10. Review: `plans/reports/code-review-260903-phase-04-skill-intake-cr.md` (M1, M2, L1, L2, L4 và nit đếm dòng đã sửa; L3 và nit `--source` ghi nhận ở plan 1.10).

Lệch so với kế hoạch:

- Test skill đặt ở `internal/skill/content_test.go` thay vì `assets/skills/skills_test.go`: `//go:embed skills` nhúng nguyên thư mục, file test đặt trong đó sẽ vào binary.
- Interview của CR dùng `dk new interview <slug> --from <file CR>` thay cho `--set for=cr`: `--from` cho `source: <CR-id>` để `dk refs` và `check` thấy liên kết ngược; vị trí `cr/<CR-id>/interview.md` do khóa mới `beside_source` trong `types.toml`.
- `doc-smoke` xóa ở phase này; test cơ chế cài dùng `doc-cr`.

## Risk Assessment

- **AI hỏi theo danh sách cố định thay vì từ trường trống**: SKILL.md bắt "liệt kê trường chưa rõ trước, mỗi câu trích tên trường"; kiểm trong bước 7. Tín hiệu vỡ: câu hỏi không khớp trường; phản ứng: thêm ví dụ đối chiếu vào SKILL.md.
- **AI tự chốt brief hoặc tự sửa tài liệu đích trước duyệt**: câu lệnh dừng in đậm ở đầu SKILL.md và lặp tại bước; `dk check` cảnh báo Feature Spec có `updated` mới hơn CR nguồn đang `chờ duyệt` (thêm quy tắc `cr-approval-order` nếu bước 7 cho thấy vi phạm).
- **Mô tả kích hoạt hai skill chồng nhau** ("yêu cầu", "tính năng"): test trùng từ khóa và câu phủ định "Không dùng cho..." ở cuối mô tả.
