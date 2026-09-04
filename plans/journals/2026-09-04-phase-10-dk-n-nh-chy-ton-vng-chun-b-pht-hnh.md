---
title: "Phase 10 dk: ổn định, chạy toàn vòng, chuẩn bị phát hành"
date: 2026-09-04
summary: "Thêm doctor, self-check, init --agent-context, CI và goreleaser; toàn vòng 12 bước trên dự án mẫu lộ 4 chỗ skill mâu thuẫn với check và 5 lỗi CLI nhỏ, đã sửa; tag v0.1.0 chờ remote"
---

# Phase 10 dk: ổn định, chạy toàn vòng, chuẩn bị phát hành

## Đã làm
- Lệnh mới `dk doctor` (bảng mục | trạng thái | cách sửa, mã 3), `dk self-check` (hash sha256 nội dung nhúng), `dk init --agent-context` (45 dòng, in ra stdout). `Target` thêm `InstalledHooks`.
- `make lint-skills`, `make build-all` (4 binary), `.github/workflows/ci.yml` (ubuntu, macos), `release.yml`, `.goreleaser.yaml`. `goreleaser check` chỉ báo thiếu remote.
- Docs viết tay: `docs/kien-truc.md`, `docs/lenh.md`, `docs/skill.md`; README rút gọn, có mục giới hạn đã biết.
- Toàn vòng do agent riêng đóng vai trên `/tmp/dk-fullrun/quan-ly-don`: 195 lệnh `dk`, 16 commit qua pre-commit thật, `dk check` cuối 0 lỗi. Report `plans/260903-1400-skill-tai-lieu/reports/phase-10-full-run.md`.

## Lỗi tìm được và sửa
- Chuỗi phiên bản lệch giữa `make install` (v0.1.0) và goreleaser (0.1.0) làm `skill status` báo "cũ" sai. Sửa: bỏ tiền tố `v` ở Makefile và `root.go`.
- Doctor bỏ qua scope `--global` và mất bảng khi target lỗi. Sửa.
- Gộp changelog trong 10 phút ghi đè tóm tắt cũ: giờ nối bằng `; ` (đổi quyết định 1.7).
- Skill: `doc-design-flow` không thể qua `userflow-steps` trước khi có spec dù cho phép làm flow trước; mockup không có CSS component chung; `doc-cr` liên kết file chưa có; `doc-feature-spec` im lặng với spec `implemented`. Sửa chữ từng chỗ, giữ trigram test.
- CLI nhỏ: `dk check <thư mục>` thông báo sai, `ui-test-checklist` in số Scenario, `userflow --from brief` mất title, `dk refs` liệt kê changelog, mockup dùng `--space-1` làm độ dày viền.

## Quyết định
- Không tag `v0.1.0` và không commit: repo chưa có remote, module path vẫn `github.com/vner/dk`; cần chủ dự án chốt. Ghi ở plan mục 1.16 và câu hỏi mở.
- Không kiểm được trong phase này: kích hoạt skill thật trong Claude Code, CI thật, `go install @v0.1.0`.

## Tiếp theo
- Chủ dự án: tạo remote, đổi module path (go.mod, import, Makefile, .goreleaser.yaml, README), chọn LICENSE, commit đầu, tag `v0.1.0`.
- Phase 11: target Codex và Windows.

> Historical work record — not durable authority. Prefer docs/specs/ADRs for current decisions.
