---
title: "Skill doc-feature-list: dựng feature list từ Product brief"
description: "Thêm loại tài liệu feature-list (docs/overview/feature-list.md) và skill thứ 12 doc-feature-list: từ Product brief approved, agent tách sản phẩm thành bảng tính năng dự kiến có mã tạm, nhóm, ưu tiên MoSCoW, dẫn mục brief; dừng ở review"
status: completed
priority: P2
effort: "1d"
tags: [go, skills, docs, template]
created: 2026-09-04
blockedBy: []
blocks: []
---

# Skill doc-feature-list: dựng feature list từ Product brief

- Nguồn: `plans/reports/brainstorm-260904-2328-skill-feature-list.md` (hợp đồng đã chốt, phương án A). Plan không chép lại phần so phương án.
- Plan liên quan: `plans/260903-1400-skill-tai-lieu/` (11 skill, mã xong). Plan này nối thêm, không sửa gì plan cũ đã làm; không chặn nhau.

## Hợp đồng

- **Kết quả**: người nói "Dựng feature list từ Product brief", agent tạo `docs/overview/feature-list.md` bằng `dk new feature-list --from <brief>`, điền bảng tính năng, dừng ở `status: review`. Mỗi dòng sau đó đi tiếp qua `doc-intake` cấp `feature`; cột Spec điền `F-xxx` khi spec ra đời.
- **Ràng buộc**: chỉ nhận brief `level: project`, `kind: product`, `status: approved`; mọi file `docs/` tạo bằng `dk new`; skill trung lập target, SKILL.md dưới 300 dòng, có `references/rules.md`, bước 0 kiểm `dk --version`, câu dừng in đậm đầu file, ba phần `new`, `update`, `html`; mô tả kích hoạt không trùng 3 từ liên tiếp với 11 skill kia; không sửa `doc-feature-spec`, không thêm `from.feature-list` vào `feature-spec`; không đụng lệnh `index`.
- **Ngoài phạm vi**: tách Feature brief thành tính năng con; tạo Feature Spec thẳng từ dòng list; ước lượng effort; tự xếp roadmap; rule `dk check` mới; đồng bộ ngược trạng thái từ Feature catalog vào list.
- **Nghiệm thu**: xem Success Criteria.

## Quyết định thiết kế (dùng chung ba phase)

- Loại `[feature-list]` trong `assets/types.toml`, đặt sau `[glossary]`: `dir = "overview"`, `name = "feature-list.md"`, `id = "none"`, `required = ["type", "title", "status", "owner", "created", "source"]`, `statuses = ["draft", "review", "approved"]`, `warn_lines = 200`, `from.brief = { title = "title" }`. Không đặt `final` (chưa loại nào `source` trỏ về list, tránh cảnh báo "tài liệu chết"). `source` do `dk new --from` tự điền, không cần khai trong `from`.
- Template `assets/templates/feature-list.md`: mục 1 frontmatter; `## 2. Nguồn` (liên kết brief); `## 3. Bảng tính năng` 7 cột `Mã | Tên | Mô tả | Nhóm | Ưu tiên | Nguồn | Spec`; `## 4. Nhóm` (mỗi nhóm một câu, khớp Product overview mục 4); `## 5. Chưa rõ` (điều brief nói nhưng chưa tách được thành tính năng, để hỏi người). Mã tạm `FL-01`, `FL-02`; ưu tiên `must | should | could | won't`; Nguồn ghi mục brief (`brief §1`, `brief §4`) hoặc `idea §3`.
- Skill `doc-feature-list`: `new` kiểm brief bằng `dk index intake` và ba trường, sai thì dừng; tạo file; đọc `idea.md` cùng thư mục và bốn mục brief; mỗi dòng tính năng dẫn đúng mục sinh ra nó; không bịa tính năng brief không nhắc, thứ suy ra ghi ở mục 5. `update`: `draft`/`review` sửa theo lời người; `approved` chỉ sửa theo CR `approved`/`in-progress` có dòng "Feature list: Có", riêng điền cột Spec không cần CR. `html`: `dk render` + mở kiểm.
- Khối "Sau mỗi lần sửa": `dk changelog add docs/overview/feature-list.md --summary "<tóm tắt thật>" --source <yymmdd>-<slug>`, `dk render`, `dk check docs/overview/feature-list.md`. Không có `dk index overview` (thư mục này không có chỉ mục), không gọi.
- Số đếm đổi: template 41 → 42 (`internal/cli/cli_test.go`, `docs/lenh.md`), skill 11 → 12 (`internal/skill/install_test.go`, README, `docs/skill.md`, `docs/kien-truc.md`, `assets/agent-context.md`).

## Phases

| # | Phase | Status | Phụ thuộc |
|---|-------|--------|---|
| 1 | [Loại và template feature-list](./phase-01-loai-va-template.md) | Completed | |
| 2 | [Skill doc-feature-list](./phase-02-skill-doc-feature-list.md) | Completed | 1 |
| 3 | [Docs và chạy vòng thật](./phase-03-docs-va-vong-that.md) | Completed | 1, 2 |

## Success Criteria

- [x] `dk new feature-list <slug> --from docs/intake/<d>/brief.md` tạo `docs/overview/feature-list.md`, `title` và `source` chép từ brief; `dk template list` đếm 42.
- [x] Bảng mục 3 có 7 cột; `dk check docs/overview/feature-list.md` không lỗi; `dk render` ra HTML; changelog có dòng với tóm tắt nội dung.
- [x] Skill dừng, không tạo file khi brief sai một trong ba trường `level`, `kind`, `status`.
- [x] `dk skill install` cài 12 skill; `make lint-skills`, `go test ./...`, `go vet ./...` qua.
- [x] README, `docs/skill.md`, `docs/quy-trinh.md`, `docs/lenh.md`, `docs/kien-truc.md`, `assets/agent-context.md` nói 12 skill và có dòng `doc-feature-list`.
- [x] Chạy thật một lần trên dự án mẫu: brief mẫu → feature list `review` → người duyệt → một dòng đi qua `doc-intake` thành brief tính năng; ghi kết quả vào report của plan.

## Câu hỏi chưa chốt

Không.
