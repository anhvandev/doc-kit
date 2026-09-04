---
phase: 11
title: "Phase 11: Target Codex và Windows"
status: done
priority: P2
effort: "3d"
dependencies: [10]
---

# Phase 11: Target Codex và Windows

<!-- Updated: Validation Session 1 - init không tự cài; thêm build Windows vào v0.2.0 -->
<!-- Updated: 2026-09-04 10:50 - xong trừ kiểm thật trong Codex đã đăng nhập, Windows máy thật, tag v0.2.0; kết quả ở reports/phase-11-codex.md -->

## Overview

Thêm target `codex`: cùng nội dung skill, thư mục `.codex/skills/`, hook trong `hooks.json`, `AGENTS.md`. Kiểm thật với Codex CLI cài trên máy. Không đổi SKILL.md. Cùng bản `v0.2.0`: bật build và CI Windows.

## Requirements

- Functional: `dk skill install --target codex [--global]`, `uninstall`, `status`; `dk hook install --target codex [--global]` ghi `.codex/hooks.json` hoặc `~/.codex/hooks.json` với cấu trúc sự kiện `PreToolUse`, `PostToolUse`; `dk hook run` đọc payload Codex (kiểm tra trường đường dẫn thật khi làm); `--target claude,codex` cài cả hai trong một lệnh `skill install` hoặc `hook install`; `dk doctor` kiểm Codex khi `.codex/` có mặt; Windows: đường dẫn `%USERPROFILE%`, pre-commit trên Git for Windows, hook Claude qua `shell: powershell` hoặc `cmd`, CI windows-latest, goreleaser thêm windows/amd64 và arm64; `init --agent-context` in cùng khối cho `AGENTS.md`.
- Non-functional: không sửa bất kỳ file nào trong `assets/skills/`; test không cần Codex cài; kiểm thật ghi report.

## Architecture

- `internal/target/codex.go`: `SkillDir`: `<root>/.codex/skills` hoặc `$CODEX_HOME/skills` mặc định `~/.codex/skills`; `InstallHooks`: ghi `hooks.json` dạng `{"hooks": {"PreToolUse": [{"matcher": "apply_patch", "hooks": [{"type": "command", "command": "dk hook run pre-write"}]}]}}`; idempotent và gỡ theo tiền tố `dk hook run` như Claude. Không ghi vào `config.toml` để không đụng cấu hình khác của người dùng. *(Sửa khi làm: `command` là chuỗi, không phải mảng như dự kiến; matcher `apply_patch` thay `.*`; xem plan mục 1.17.)*
- `internal/hook/run.go`: nhận cả hai dạng payload; đọc `tool_input.file_path`, dự phòng `tool_input.path`, và với `apply_patch` đọc danh sách file trong patch (`*** Update File:` / `*** Add File:`); trả deny theo định dạng Codex nếu khác Claude (xác minh lúc làm, tài liệu hooks Codex phiên bản hiện hành).
- Matcher: Codex có thể chỉ bắt `Bash`/`exec_command` và `apply_patch` tùy phiên bản; `hook install` ghi matcher rộng (`.*`) và để `hook run` tự lọc theo đường dẫn.
- `registry.go`: đăng ký `codex`; `--target` nhận danh sách phân tách phẩy.
- Trust: sau `hook install --target codex`, in hướng dẫn người dùng trust hook trong Codex (lệnh hoặc lần chạy đầu Codex hỏi).

## Related Code Files

- Create: `internal/target/codex.go`, `internal/target/codex_test.go` (golden `hooks.json`), `internal/hook/run_codex_test.go` (payload mẫu ghi lại từ Codex thật)
- Modify: `internal/target/registry.go`, `internal/hook/run.go`, `internal/cli/{init,skill,hook,doctor}.go`, `assets/agent-context.md` (ghi chú AGENTS.md), `README.md`, `docs/` của repo CLI

## Implementation Steps

1. Cài Codex CLI trên máy phát triển; đọc tài liệu hooks và skills phiên bản đang cài; ghi lại một payload PreToolUse và PostToolUse thật vào `testdata/` (bằng hook tạm `cat > /tmp/payload.json`).
2. Viết `codex.go`; test golden với `hooks.json` có mục người dùng sẵn.
3. Sửa `run.go` cho payload Codex và `apply_patch`; test.
4. `--target` danh sách cho `skill install` và `hook install`; `doctor`; test.
5. Kiểm thật: `dk skill install --target codex && dk hook install --target codex` trong dự án mẫu; mở Codex, `$doc-intake` chạy vòng Tiếp nhận; thử tạo file mới vào `docs/` bằng apply_patch → quan sát hook có deny không; sửa file có sẵn → changelog có dòng; ghi kết quả vào `plans/260903-1400-skill-tai-lieu/reports/phase-11-codex.md`, gồm phiên bản Codex và phạm vi hook thật sự bắt được.
6. Nếu hook Codex không bắt ghi file: ghi rõ trong README "Codex: bảo đảm bằng pre-commit và skill, hook chỉ hỗ trợ khi phiên bản Codex bắt apply_patch"; không cố vá.
7. Windows: sửa `target.SkillDir`, `hook` (đường dẫn và shell), `gitx` (CRLF), pre-commit; CI thêm windows-latest; goreleaser thêm windows; kiểm `dk --version` và `dk init` trên máy Windows thật hoặc VM.
8. README, docs, release `v0.2.0`.

## Success Criteria

- [x] `dk skill install --target codex` rồi `uninstall` sạch; `status` đúng cho cả hai target (test `TestInstallEndToEnd`)
- [ ] Cùng bộ skill (không đổi byte nào trong `assets/skills/`) chạy vòng Tiếp nhận trong Codex, tạo file đúng qua `dk new` — `assets/skills/` không đổi (kiểm mtime, lint sạch); vòng chạy thật cần `codex login`, chưa kiểm
- [x] `hooks.json` giữ nguyên mục người dùng; gỡ về nguyên trạng (golden `TestCodexHooksGolden`)
- [x] Report ghi rõ hook Codex bắt được gì với phiên bản nào; README phản ánh đúng (`reports/phase-11-codex.md`, Codex 0.153.2, `apply_patch`)
- [x] Pre-commit chặn commit thiếu changelog trong dự án Codex như với Claude (pre-commit không phụ thuộc target; `dk init` giống nhau)
- [ ] CI xanh trên windows-latest; binary windows/amd64 chạy `dk init` và `dk skill install` — `GOOS=windows` vet và biên dịch test qua, 2 binary Windows build được; CI và máy Windows thật chờ remote
- [ ] Release `v0.2.0` có 6 nền tảng — goreleaser, Makefile, CI đã 6 nền tảng; tag chặn cùng lý do `v0.1.0`

## Risk Assessment

- **Định dạng payload và deny của Codex khác tài liệu hoặc đổi giữa phiên bản**: ghi payload thật vào testdata ở bước 1 trước khi viết mã; test bám payload thật; README ghi phiên bản đã kiểm.
- **Hook cấp dự án cần trust**, người dùng bỏ qua: `doctor` không kiểm được trust từ ngoài; in hướng dẫn rõ sau `hook install`; không coi hook là bảo đảm.
- **`$CODEX_HOME` khác mặc định**: đọc biến môi trường trước, fallback `~/.codex`.
