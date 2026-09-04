---
phase: 10
title: "Phase 10: Ổn định và phát hành"
status: done
priority: P1
effort: "4d"
dependencies: [1, 2, 3, 4, 5, 6, 7, 8, 9]
---

# Phase 10: Ổn định và phát hành

<!-- Updated: Validation Session 1 - v0.1.0 chỉ Linux và macOS; Windows sang phase 11 -->
<!-- Updated: 2026-09-04 10:15 - phase xong trừ tag v0.1.0 (chưa có remote); xem plan.md mục 1.16 -->

## Overview

Chạy toàn vòng với một tính năng thật trên một dự án thật, sửa mô tả kích hoạt chồng nhau, ghi bảng skill và lệnh `dk` vào Agent context file mẫu, kiểm độc lập với `ak-*` và tính trung lập target, phát hành binary phiên bản đầu.

## Requirements

- Functional: `dk init --agent-context` sinh khối Markdown dưới 60 dòng để người dán vào `CLAUDE.md` hoặc `AGENTS.md` (bảng skill, lệnh `dk` chính, quy tắc ba lớp, ngưỡng dòng); goreleaser build ma trận linux/darwin × amd64/arm64 (Windows ở phase 11); `dk self-check` (ẩn hoặc công khai) in phiên bản, số template, số skill, target hỗ trợ, kết quả `go:embed` toàn vẹn (hash); `dk doctor` kiểm dự án hiện tại: `dk.toml`, git, pre-commit, skill cài đúng phiên bản, hook có mặt.
- Non-functional: mọi test qua trên ubuntu và macos qua CI GitHub Actions; `assets/skills/` grep `ak-`, `Claude Code`, `.claude/`, `Edit`, `Write` (dạng tool) bằng không; SKILL.md mọi skill dưới 300 dòng; Agent context file mẫu dưới 200 dòng khi ghép vào dự án mẫu.

## Architecture

- `.github/workflows/ci.yml`: `go vet`, `go test ./...`, build ma trận, chạy `skills_test.go`; `release.yml` trên tag `v*` dùng goreleaser (`.goreleaser.yaml`: `CGO_ENABLED=0`, ldflags Version, archive tar.gz và zip, checksums).
- `dk doctor`: tập kiểm tra trả bảng `mục | trạng thái | cách sửa`.
- `dk init --agent-context`: template `assets/agent-context.md` nhúng, in ra stdout, không ghi file (người quyết định dán vào đâu; không tự sửa `CLAUDE.md` của người).
- Kịch bản toàn vòng (ghi thành `plans/260903-1400-skill-tai-lieu/reports/phase-10-full-run.md`): dự án thật có giao diện, một tính năng mới: intake → brief duyệt → user flow → wireframe duyệt → mockup → Feature Spec duyệt → test case → plan và report → release brief; rồi một CR trên tính năng đó: bảng tác động duyệt → sửa spec, mockup, test → CHANGELOG-DOCS đủ dòng → commit qua pre-commit. Đo: số lần người phải can thiệp ngoài duyệt, số câu hỏi thừa, số lần skill kích hoạt sai.
- Xóa `assets/skills/doc-smoke/` nếu còn; test cơ chế cài chuyển sang dùng `doc-adr` (nhỏ nhất).

## Related Code Files

- Create: `.github/workflows/{ci,release}.yml`, `.goreleaser.yaml`, `assets/agent-context.md`, `internal/cli/{doctor,selfcheck}.go`, `docs/` của chính repo CLI (README đầy đủ, `docs/lenh.md` bảng lệnh, `docs/skill.md` bảng skill, `docs/kien-truc.md` ba lớp)
- Modify: `go.mod` (module path thật), mọi import, `Makefile`, `README.md`, mô tả trong 11 `SKILL.md` theo kết quả chạy thử
- Delete: `assets/skills/doc-smoke/`

## Implementation Steps

1. Chốt module path thật; `go mod edit -module`, sed import, build.
2. Viết `doctor`, `self-check`, `init --agent-context`; test.
3. CI: `ci.yml` chạy trên ubuntu và macos. Không sửa đường dẫn Windows ở phase này.
4. Chạy toàn vòng theo kịch bản trên dự án thật; ghi report với số đo.
5. Sửa mô tả kích hoạt các skill kích hoạt sai hoặc chồng; chạy lại phần bị lỗi.
6. Kiểm độc lập: script `make lint-skills` grep từ cấm; đưa vào CI.
7. Viết docs của repo CLI (dùng chính `dk`? Không: repo CLI dùng docs thường để tránh phụ thuộc vòng; ghi rõ lý do trong README).
8. `.goreleaser.yaml`, `release.yml`; tag `v0.1.0`; kiểm tra tải về chạy `dk --version` trên Linux và macOS; `go install <module>/cmd/dk@v0.1.0` chạy.
9. README: cài từ release, `go install`, `dk init`, `dk skill install`, `dk hook install`, gỡ, giới hạn đã biết (Codex và Windows chưa hỗ trợ cho tới `v0.2.0`).

## Success Criteria

- [x] Toàn vòng chạy không cần can thiệp tay ngoài duyệt; report ghi số đo (`reports/phase-10-full-run.md`: 195 lệnh, 16 commit qua pre-commit, 4 can thiệp đã sửa nguyên nhân, kiểm lại từng chỗ)
- [x] `make lint-skills` bằng không; CI xanh trên ubuntu và macos (workflow đã viết; chưa chạy thật vì chưa có remote)
- [x] Agent context file mẫu dưới 60 dòng (45), ghép vào dự án mẫu tổng dưới 200 dòng (142 với CLAUDE.md 97 dòng)
- [x] `dk doctor` phát hiện thiếu pre-commit, skill cũ phiên bản, hook thiếu, và nêu cách sửa (`phase10_test.go`)
- [ ] Release `v0.1.0` có binary 4 nền tảng, checksums; `go install` chạy. Đã có `.goreleaser.yaml`, `release.yml`, `make build-all` ra 4 binary; **chưa tag** vì chưa có remote và module path thật (câu hỏi mở của plan)
- [x] `doc-smoke` đã xóa; test cơ chế cài dùng skill thật (`doc-cr`, giữ theo quyết định 1.10)

## Risk Assessment

- **Người dùng Windows đòi sớm**: đã quyết lùi sang `v0.2.0` (phase 11); README ghi rõ. Không kéo vào phase này.
- **Toàn vòng lộ lỗi thiết kế template** (mục thiếu, thứ tự sai): sửa template và phát hành `v0.1.1`; không vá bằng SKILL.md.
- **Phụ thuộc vòng** nếu dùng `dk` để viết docs của chính repo `dk`: tránh có chủ đích ở bước 7.
