---
phase: 3
title: "Docs và chạy vòng thật"
status: completed
priority: P2
effort: "2h"
dependencies: [1, 2]
---

# Phase 3: Docs và chạy vòng thật

## Overview

Cập nhật tài liệu viết tay của repo (12 skill, 42 loại), nối `doc-feature-list` vào vòng đời trong `quy-trinh.md`, rồi chạy thật một vòng brief → feature list → duyệt → intake tính năng trên dự án mẫu và ghi report.

## Requirements

- Functional: mọi chỗ nói "11 skill" thành 12 và có dòng cho `doc-feature-list`; `docs/quy-trinh.md` nói khi nào gọi skill này và người phải làm gì; vòng thật chạy chỉ bằng skill gọi `dk`.
- Non-functional: docs viết tay, không dùng `dk` (lý do ở `docs/kien-truc.md`); giữ giọng và độ dài hiện có; commit theo conventional commit, không nhắc AI.

## Architecture

Không có mã. Tài liệu chạm: README, `docs/skill.md`, `docs/quy-trinh.md`, `docs/kien-truc.md`, `docs/lenh.md` (đã đổi 42 ở phase 1). Vòng thật chạy trong thư mục tạm ngoài repo, dùng binary `bin/dk` vừa build.

## Related Code Files

- Modify: `README.md` (dòng 3, dòng 38: 11 → 12), `docs/skill.md` (tiêu đề "11 skill nhúng" → 12; thêm dòng bảng sau `doc-overview`), `docs/quy-trinh.md` (sơ đồ vòng đời, bảng "Khi nào gọi gì", mục "Lần đầu với dự án mới"), `docs/kien-truc.md` (hai chỗ "11 skill")
- Create: `plans/260904-1658-skill-feature-list/reports/report-<yymmdd-hhmm>-vong-that-feature-list.md`

## Implementation Steps

1. `docs/skill.md`: tiêu đề `## 12 skill nhúng`; dòng bảng: ``| `doc-feature-list` | `docs/overview/` | Feature list: bảng tính năng tách từ Product brief, mã tạm `FL-xx`, nhóm, ưu tiên MoSCoW, cột Spec | list `review`; brief chưa `approved` |``.
2. `docs/quy-trinh.md`:
   - Sơ đồ vòng đời: giữ nguyên (vòng đời một tính năng không đổi). Thêm sau sơ đồ một đoạn ngắn "Cấp dự án: sau khi Product brief `approved`, gọi `doc-feature-list` để có bảng tính năng dự kiến; mỗi dòng của bảng là một lượt `doc-intake` cấp tính năng."
   - Bảng "Khi nào gọi gì", sau dòng "Dự án mới, chưa có tài liệu nền": ``| Product brief đã `approved`, chưa biết có những tính năng gì | `doc-feature-list` | "Dựng feature list từ Product brief" | Duyệt bảng: bỏ dòng bịa, quyết mục Chưa rõ, đổi `status: approved`; rồi gọi `doc-intake` cho từng dòng theo ưu tiên |``.
   - Mục "Lần đầu với dự án mới": bước 2 thêm "rồi `doc-feature-list`" sau `doc-overview`.
3. README dòng 3 và 38, `docs/kien-truc.md` hai chỗ: 11 → 12. `grep -rn "11 skill" README.md docs/ assets/agent-context.md` phải rỗng.
4. `make build`, rồi trong thư mục tạm: `dk init`, `dk skill install`, `dk init --agent-context`, `dk doctor`. Tạo Product brief mẫu qua `doc-intake` cấp `project` (hoặc `dk new idea` + `dk new brief` điền tay để rút ngắn), đặt `approved`.
5. Mở harness agent trong thư mục tạm, nói "Dựng feature list từ Product brief". Quan sát: skill kiểm ba trường, gọi `dk new feature-list`, bảng có Nguồn từng dòng, mục 5 có nội dung hay "không", dừng ở `review`, changelog có dòng tóm tắt thật, `dk check` sạch. Thử lần hai với brief `status: review`: skill phải dừng, không tạo file.
6. Người đổi `status: approved`; nói "Tôi có ý tưởng: <tên dòng FL-01>" để `doc-intake` chạy cấp tính năng; sau đó nói "Điền mã spec cho FL-01" khi có spec giả lập (hoặc bỏ nếu chưa có spec). Ghi lại lệnh skill đã chạy và lệch so với SKILL.md.
7. Sửa SKILL.md theo lệch quan sát (chỉ phần gây lệch), chạy lại `go test ./internal/skill/`, `make lint-skills`.
8. Report theo template `report.md` của plan cũ: bằng chứng là mã commit, đường dẫn file sinh, khối kết quả `dk check`, `dk status`; danh sách lệch và cách sửa.
9. Commit tách theo loại: `feat(types): thêm loại feature-list và template`, `feat(skill): thêm doc-feature-list`, `docs: 12 skill, quy trình feature list`.

## Success Criteria

- [x] `grep -rn "11 skill" README.md docs/ assets/` rỗng; `docs/skill.md` và `docs/quy-trinh.md` có `doc-feature-list`.
- [x] Vòng thật: feature list `review` có đủ 7 cột, mọi dòng có Nguồn, không dòng nào bịa ngoài brief và idea; brief `review` thì skill dừng không tạo file.
- [x] `dk doctor` trên dự án mẫu báo đủ 12 skill đúng phiên bản và khối ngữ cảnh agent mới.
- [x] Report có bằng chứng (`report-evidence` không cảnh báo nếu chạy `dk check` lên nó).
- [x] `go test ./...`, `go vet ./...`, `make lint-skills` qua ở commit cuối.

## Risk Assessment

- Harness agent không có sẵn để chạy vòng thật trong phiên: dấu hiệu là bước 5 không thực hiện được. Ứng phó: làm bước 4 và đọc SKILL.md theo tay với brief mẫu (chạy từng lệnh `dk` như skill ghi), ghi rõ trong report "chưa chạy trong harness"; không đánh dấu tiêu chí vòng thật là đạt.
- `docs/quy-trinh.md` vượt ngưỡng dòng hoặc phá giọng: giữ mỗi bổ sung dưới 5 dòng.
