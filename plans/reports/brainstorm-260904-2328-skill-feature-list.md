---
title: "Brainstorm: skill doc-feature-list dựng feature list từ Product brief"
type: brainstorm
created: 2026-09-04 23:28
status: accepted
---

# Skill `doc-feature-list`: dựng feature list từ Product brief

## Tóm tắt

Thêm skill thứ 12 và loại tài liệu `feature-list`: từ Product brief `approved`, agent tách
sản phẩm thành bảng tính năng dự kiến (mã tạm, tên, mô tả, nhóm, ưu tiên MoSCoW, dẫn mục
brief), dừng ở `review`. Lấp khoảng trống hiện có: Product overview chỉ có "nhóm tính
năng", Roadmap cần `F-xxx` hoặc tên brief, Feature catalog chỉ sinh từ spec đã tồn tại.

## Hợp đồng

- **Kết quả**: người dùng nói "Dựng feature list từ Product brief", agent tạo
  `docs/overview/feature-list.md` bằng `dk new feature-list --from <brief>`, điền bảng,
  dừng ở `status: review`. Người duyệt, rồi mỗi dòng đi tiếp qua `doc-intake` cấp
  `feature` như hiện nay; cột "Spec" điền `F-xxx` khi spec ra đời.
- **Ràng buộc**: chỉ nhận brief `level: project`, `kind: product`, `status: approved`;
  mọi file trong `docs/` tạo bằng `dk new`; skill trung lập target, `SKILL.md` dưới
  300 dòng, có `references/rules.md`, bước 0 kiểm `dk --version`, câu dừng in đậm đầu
  file, ba phần `new`, `update`, `html`; mô tả kích hoạt không trùng 3 từ liên tiếp với
  11 skill kia (`internal/skill/content_test.go`); không sửa `doc-feature-spec`,
  không thêm `from.feature-list` vào `feature-spec`.
- **Ngoài phạm vi**: tách Feature brief thành tính năng con; tạo Feature Spec thẳng từ
  dòng list; ước lượng effort; tự xếp roadmap; rule `dk check` mới riêng cho loại này;
  đồng bộ ngược trạng thái từ Feature catalog vào list.
- **Nghiệm thu**:
  - `dk new feature-list <slug> --from docs/intake/<d>/brief.md` tạo
    `docs/overview/feature-list.md` với `title`, `source` chép từ brief; brief sai một
    trong ba trường thì skill dừng, không tạo file.
  - Bảng có đủ cột: Mã tạm (`FL-01`...), Tên, Mô tả một câu, Nhóm, Ưu tiên (`must`,
    `should`, `could`, `won't`), Nguồn (mục brief), Spec (rỗng cho đến khi có `F-xxx`).
  - `dk check docs/overview/feature-list.md` không lỗi; `dk render` ra HTML;
    `dk changelog add` có dòng với tóm tắt nội dung.
  - `dk skill install` cài 12 skill; `make lint-skills`, `go test ./...`, `go vet ./...`
    qua; `install_test.go` đếm 12.
  - `docs/skill.md`, `docs/quy-trinh.md`, README, `assets/agent-context.md` nói 12 skill
    và có dòng cho `doc-feature-list`.

## Phương án đã so

| | A. Skill mới + loại `feature-list` trong `docs/overview/` | B. Thêm phần vào `doc-overview` | C. Loại mới đặt trong `docs/features/` |
|---|---|---|---|
| Khớp yêu cầu "thêm một skill" | Có | Không (vẫn 11 skill) | Có |
| Đụng mã CLI | chỉ `types.toml`, template, đếm skill trong test | không | phải sửa `flatRows` của `dk index features` để không liệt kê list như một spec |
| Giả định dễ vỡ nhất | người chấp nhận list ở `overview/` thay vì `features/` | `doc-overview` phình, SKILL.md tiến gần 300 dòng | index Feature catalog lẫn dòng lạ |

Chọn **A**: nhỏ nhất, không đụng lệnh `index`, đúng chỗ tài liệu cấp dự án đọc đầu tiên.

## Quyết định chi tiết

- Loại `[feature-list]`: `dir = "overview"`, `name = "feature-list.md"`, `id = "none"`,
  `required = [type, title, status, owner, created, source]`,
  `statuses = [draft, review, approved]`, `from.brief = { title, source }`, không đặt
  `final` để tránh cảnh báo "tài liệu chết" (chưa có loại nào `source` trỏ về list).
- Template: mục 1 frontmatter; 2 Nguồn (liên kết brief); 3 Bảng tính năng; 4 Nhóm
  (mỗi nhóm một câu, khớp Product overview mục 4); 5 Chưa rõ (dòng brief không tách
  được thành tính năng, để hỏi người). Dưới 200 dòng, đặt `warn_lines = 200`.
- Skill `new`: kiểm brief → tạo file → đọc idea.md cùng thư mục và bốn mục brief → mỗi
  dòng tính năng dẫn đúng mục brief sinh ra nó; không bịa tính năng brief không nhắc,
  thứ suy ra thì ghi ở mục 5 chờ người. `update`: chỉ `draft`/`review`, hoặc `approved`
  khi có CR `approved` (điền cột Spec không cần CR). `html`: `dk render`.
- Roadmap template gợi ý "tên brief khi chưa có spec": thêm "hoặc mã tạm FL-xx" (một
  dòng chú thích, không đổi cấu trúc).

## Tệp sẽ chạm

`assets/types.toml`, `assets/templates/feature-list.md`, `assets/skills/doc-feature-list/
{SKILL.md,references/rules.md}`, `assets/templates/roadmap.md` (chú thích),
`assets/agent-context.md`, `internal/skill/install_test.go`, `docs/skill.md`,
`docs/quy-trinh.md`, `docs/lenh.md` (nếu liệt kê loại), `README.md`.

## Câu hỏi chưa chốt

Không.
