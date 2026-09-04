---
title: "Report: vòng thật doc-feature-list trên dự án mẫu"
type: report
created: 2026-09-05 00:10
plan: plans/260904-1658-skill-feature-list
---

# Vòng thật doc-feature-list

## Cách chạy

Dự án mẫu `/tmp/dk-vong-that`, binary `bin/dk` build từ working tree (chưa commit).
Agent (phiên này) làm theo `assets/skills/doc-feature-list/SKILL.md` từng bước, gõ lệnh
`dk` như skill ghi; chưa chạy trong harness riêng với skill đã cài.

## Bằng chứng

- `dk init`, `dk skill install` (12 skill), `dk init --agent-context`, `dk doctor`:
  ```
  OK agent context (CLAUDE.md)  | có, đúng phiên bản
  OK agent context (AGENTS.md)  | có, đúng phiên bản
  OK skill (claude, dự án)      | đủ, đúng phiên bản
  ```
  (1 mục chưa đạt là `dk` trên PATH của pre-commit, do dùng `bin/dk` tạm.)
- Brief `docs/intake/260905-quan-ly-don/brief.md`, `level: project`, `kind: product`.
- Kiểm âm: đặt brief `status: review` → bước 1 skill dừng; `docs/overview/feature-list.md`
  không tồn tại.
- Kiểm dương: brief `approved` → `dk new feature-list quan-ly-don --from ... --set owner=An`
  tạo `docs/overview/feature-list.md`, `title` và `source` chép đúng.
- Điền thân theo bước 3: 5 dòng (FL-01..03 `must` từ brief §1, §4; FL-04, FL-05 `won't`
  từ brief §3), 2 nhóm, mục 5 có 2 điều chưa rõ (đăng nhập không được nhắc; nhập Zalo
  gõ tay hay tự động). Không dòng nào thiếu Nguồn.
- Khối "Sau mỗi lần sửa":
  ```
  - 00:05 | overview/feature-list.md | mới, 40 dòng | Tách 5 tính năng (3 must, 2 won't), 2 nhóm, 2 chưa rõ | 260905-quan-ly-don
  docs/html/overview/feature-list.html
  0 lỗi, 0 cảnh báo
  ```
  HTML: bảng mục 3 có 5 hàng dữ liệu, 7 cột; hàng mẫu `FL-01` trống của template đã
  được thay, không còn trong Markdown.
- Intake tính năng cho FL-01: `dk new idea nhap-don-tu-zalo --set level=feature` tạo
  `docs/intake/260905-nhap-don-tu-zalo/idea.md`. `dk status`:
  ```
  brief          1  approved=1
  feature-list   1  review=1
  idea           2  briefed=1  draft=1
  ```
- Pre-commit chặn commit đầu vì 2 idea thiếu changelog (đúng hành vi); thêm changelog rồi
  commit `233c29c` trong dự án mẫu; `dk check` sau commit: 0 lỗi, 0 cảnh báo.

## Lệch so với SKILL.md

Sau review: cổng `update` cho list `approved` từng đòi dòng "Feature list" trong bảng tác
động CR, mà template CR không có dòng đó; đã đổi sang cách nói của `doc-overview` (CR nêu
Feature list, hoặc Product brief mới `approved`). Thêm test `TestFeatureListFrom`. Ghi nhận: skill không nhắc changelog cho idea do intake tạo
sau đó, đó là việc của `doc-intake` (đúng phân công).

## Kết quả kiểm trong repo

`go test ./...`, `go vet ./...`, `make lint-skills`: qua. Test đếm đổi: template 41 → 42
(`cli_test.go`, `registry_test.go`), skill 11 → 12 (`install_test.go`, `init_test.go`
22 → 24 dòng trạng thái, `phase10_test.go` 9/11 → 10/12).

## Câu hỏi chưa chốt

Không.
