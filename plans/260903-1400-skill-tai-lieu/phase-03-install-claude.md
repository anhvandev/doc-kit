---
phase: 3
title: "Phase 3: Cài đặt cho Claude"
status: done
priority: P1
effort: "4d"
dependencies: [1]
---

# Phase 3: Cài đặt cho Claude

## Overview

Nhúng skill vào binary, viết interface `Target` với hiện thực `claude`, thêm `dk skill *`, `dk hook *`, lệnh ẩn `dk hook run` cho harness gọi, pre-commit git trong `dk init`. Sau phase này một dự án chạy ba lệnh `dk init`, `dk skill install`, `dk hook install` là Claude Code có đủ skill và hook.

<!-- Updated: Validation Session 1 - init không tự cài skill và hook; matcher Edit|Write; khóa metadata dk_* -->
<!-- Updated: 2026-09-03 16:50 - phase xong; quyết định triển khai ghi ở plan.md mục 1.9; bước 9 (kiểm thật trong Claude Code) chưa làm -->

## Requirements

- Functional: `dk skill list`; `dk skill install [<tên>...] [--target claude] [--global] [--force]`; `dk skill uninstall [<tên>...] [--target] [--global]`; `dk skill status`; `dk hook install|uninstall [--target] [--global]`; `dk hook run pre-write|post-edit` (ẩn, đọc JSON stdin); `dk init` **không** cài skill hay hook, chỉ cài pre-commit khi có `.git` và chưa có pre-commit (ngược lại in đoạn cần thêm) và in nhắc hai lệnh tiếp theo.
- Non-functional: `uninstall` không bao giờ xóa thư mục thiếu dấu vết `metadata.dk_installed_by: "dk"`; `install` không ghi đè skill có hash lệch nếu thiếu `--force`; ghi `settings.json` phải giữ nguyên mọi khóa khác và định dạng thụt 2 khoảng; hook fail-open: mọi lỗi nội bộ của `hook run` thoát 0 và in cảnh báo ra stderr, trừ quyết định từ chối có chủ đích.

## Architecture

- `assets/skills/<tên>/SKILL.md` và `references/rules.md`; phase này chỉ cần một skill giả `doc-smoke` để test cơ chế, xóa ở phase 10 khi 11 skill thật đã có. `assets/embed.go` thêm `//go:embed skills`.
- `internal/skill`:
  - `List() []Meta` từ embed: đọc frontmatter `name`, `description`.
  - `Hash(files map[string][]byte) string`: sha256 nối theo tên file đã sắp, nội dung SKILL.md đã bỏ khối `metadata`.
  - `Install(t target.Target, names, global, force)`: với mỗi skill: đích `SkillDir(global)/<tên>`; nếu tồn tại: đọc frontmatter, không có `metadata.dk_installed_by` thì lỗi "skill không do dk cài, dùng tên khác hoặc xóa tay"; có nhưng hash hiện tại của đích khác `metadata.dk_hash` thì lỗi "đã sửa tay" trừ `--force`; ghi mọi file, chèn `metadata: {dk_installed_by: "dk", dk_version: "<Version>", dk_hash: "<h>"}` vào frontmatter SKILL.md (giá trị đều là chuỗi theo spec).
  - `Uninstall`: chỉ xóa thư mục có dấu vết; báo bỏ qua phần còn lại.
  - `Status`: bảng `tên | target | scope | trạng thái` với trạng thái `chưa cài | đúng phiên bản | cũ (vX) | đã sửa tay | không do dk`.
- `internal/target`:
  - `claude.go`: `SkillDir(global)`: `<root>/.claude/skills` hoặc `$HOME/.claude/skills`; `InstallHooks`: đọc `settings.json` (tạo nếu chưa có) vào `map[string]any` giữ thứ tự bằng `json.Decoder` + `orderedmap` tự viết nhỏ hoặc chấp nhận sắp khóa lại chỉ trong khối `hooks`; thêm vào `hooks.PreToolUse` mục `{matcher: "Write", hooks: [{type: "command", command: "dk hook run pre-write"}]}` và `hooks.PostToolUse` mục `{matcher: "Edit|Write", hooks: [{type: "command", command: "dk hook run post-edit"}]}`; idempotent: không thêm nếu đã có lệnh cùng chuỗi. `UninstallHooks`: xóa đúng các mục có `command` bắt đầu bằng `dk hook run`; xóa khóa `hooks` nếu rỗng.
  - `Registry`: `Get(name) (Target, error)`; phase này chỉ đăng ký `claude`; `codex` trả lỗi "chưa hỗ trợ" cho tới phase 11.
- `internal/hook/run.go`:
  - Đọc stdin JSON: `tool_name`, `tool_input.file_path` (Claude Code). Chuẩn hóa đường dẫn về gốc dự án (tìm `dk.toml` đi lên từ `cwd` trong payload hoặc từ file).
  - `pre-write`: nếu `file_path` khớp `<docs_dir>/**/*.md`, không nằm trong `docs/html/`, và file **chưa tồn tại** thì từ chối: in JSON `{"hookSpecificOutput":{"hookEventName":"PreToolUse","permissionDecision":"deny","permissionDecisionReason":"Tạo tài liệu mới bằng `dk new <type> <slug>` rồi sửa file đó."}}` và thoát 0 (định dạng quyết định của Claude Code); mọi trường hợp khác thoát 0 im lặng.
  - `post-edit`: nếu `file_path` là `.md` trong `docs/`, không phải `docs/html/`, không `generated: true`: gọi nội bộ `changelog.Add` với summary "chưa tóm tắt", source "trực-tiếp"; lỗi gì cũng in cảnh báo stderr và thoát 0.
  - Không dùng `os/exec` gọi lại `dk`; gọi hàm nội bộ.
- Pre-commit: `assets/hooks/pre-commit.sh`: `command -v dk >/dev/null || { echo "dk chưa cài, bỏ qua kiểm tra changelog"; exit 0; }; dk changelog pending`. `init` chép vào `.git/hooks/pre-commit` với quyền thực thi nếu chưa có file; nếu có, in nội dung để người nối tay. Lưu ý: fail-open khi thiếu `dk` là chủ đích, ghi trong README.
- `dk init` thứ tự: tạo cây và `dk.toml` (phase 1) → pre-commit → in nhắc "Tiếp theo: `dk skill install`, `dk hook install`". Không có cờ `--target` ở `init`.

## Related Code Files

- Create: `assets/skills/doc-smoke/{SKILL.md,references/rules.md}` (tạm), `assets/hooks/pre-commit.sh`
- Create: `internal/skill/{skill,hash,install,status}.go`, `internal/target/{target,registry,claude}.go`, `internal/hook/{install,run}.go`, `internal/cli/{skill,hook}.go`
- Create test: `internal/skill/install_test.go`, `internal/target/claude_test.go` (golden `settings.json` trước và sau), `internal/hook/run_test.go` (payload JSON mẫu), `internal/cli/init_test.go` (end-to-end trong thư mục tạm có `HOME` giả)
- Modify: `internal/cli/init.go`, `assets/embed.go`, `README.md`

## Implementation Steps

1. Viết `doc-smoke` tối thiểu: frontmatter `name: doc-smoke`, `description: Skill thử cơ chế cài của dk, không dùng thật`, thân 5 dòng; `references/rules.md` 3 dòng.
2. Viết `internal/skill/hash.go` và test: đổi một byte trong `rules.md` đổi hash; thêm hoặc bỏ khối `metadata` trong SKILL.md không đổi hash.
3. Viết `internal/target/target.go` (interface, `HookEntry`), `registry.go`, `claude.go` phần `SkillDir`. Test với `HOME` tạm.
4. Viết `internal/skill/install.go`, `status.go`; test: cài mới → có dấu vết; cài lại → không đổi; sửa `rules.md` ở đích rồi cài lại → lỗi; `--force` → ghi đè; thư mục cùng tên không có dấu vết → lỗi và không đụng; `uninstall` bỏ qua thư mục không dấu vết.
5. Viết `claude.go` phần hook: merge vào `settings.json`. Test golden: file có sẵn `permissions` và một hook người dùng → sau cài giữ nguyên chúng, thêm 2 mục `dk`; gỡ → về đúng file ban đầu byte-for-byte (trừ khoảng trắng cuối nếu có).
6. Viết `internal/hook/run.go`; test bằng payload JSON: Write file mới trong `docs/features/` → JSON deny; Write file đã tồn tại → im lặng thoát 0; Edit file `docs/x.md` → changelog thêm dòng "chưa tóm tắt"; Edit `docs/html/a.html` → không thêm; stdin rác → thoát 0 kèm cảnh báo.
7. Viết `internal/cli/skill.go`, `hook.go` (đăng ký `hook run` với `Hidden: true`).
8. Sửa `init.go` thêm pre-commit và dòng nhắc; test end-to-end: `dk init && dk skill install && dk hook install` trong repo tạm rồi kiểm `.claude/skills/doc-smoke/SKILL.md` có metadata `dk_*`, `.claude/settings.json` có 2 hook, `.git/hooks/pre-commit` chạy được; `dk skill uninstall && dk hook uninstall` → `.claude/` chỉ còn `settings.json` rỗng hoặc bị xóa nếu do dk tạo.
9. Kiểm thật trong Claude Code: mở dự án mẫu, nhờ agent tạo file `docs/features/thu.md` bằng Write → bị từ chối với lý do; nhờ sửa file có sẵn bằng Edit → `CHANGELOG-DOCS.md` có dòng "chưa tóm tắt"; commit khi còn pending → bị chặn.
10. README: mục "Cài vào dự án" (cài binary, `dk init`, `dk skill install`, `dk hook install`, mở Claude Code) và mục "Gỡ" (thứ tự ngược).

## Success Criteria

- [x] `dk init` trong dự án mới: đủ cây `docs/`, pre-commit, dòng nhắc; không đụng `.claude/`
- [x] `dk skill install` và `dk hook install`: skill trong `.claude/skills/` có metadata `dk_*`, hook trong `.claude/settings.json`
- [x] `dk skill install` lần hai không đổi gì; sau khi sửa tay skill thì bị từ chối, `--force` ghi đè
- [x] `dk skill uninstall` chỉ xóa skill có dấu vết; thư mục người dùng tự tạo còn nguyên
- [x] `dk hook uninstall` trả `settings.json` về nội dung trước khi cài
- [ ] Trong Claude Code: Write file mới vào `docs/` bị từ chối kèm hướng dẫn; Edit file `docs/` tự có dòng changelog; commit bị chặn khi còn pending (bước 9, chưa kiểm thật; đã kiểm bằng payload JSON qua binary thật)
- [x] Rút `dk` khỏi PATH: hook không làm hỏng phiên (Claude báo lỗi lệnh nhưng vẫn cho thao tác), pre-commit in cảnh báo và cho qua (phần pre-commit đã kiểm; phần Claude thuộc bước 9)
- [x] `go test ./...` qua, không test nào cần Claude Code chạy

## Risk Assessment

- **Định dạng quyết định hook của Claude Code đổi**: kiểm tra tài liệu hooks hiện hành lúc làm bước 6; dự phòng dùng mã thoát 2 kèm stderr, cách cũ vẫn được hỗ trợ. Tín hiệu vỡ: deny không có hiệu lực trong bước 9; phản ứng: đổi sang exit 2.
- **`dk hook run` gọi khi `cwd` của harness không phải gốc dự án** (worktree, thư mục con): tìm `dk.toml` đi lên từ đường dẫn file, không từ `cwd`.
- **Ghi `settings.json` làm mất chú thích hoặc thứ tự** của người dùng: file này là JSON chuẩn không chú thích; thứ tự khóa cấp một giữ bằng decoder theo thứ tự; test golden bảo vệ.
- **`--global` chạm `~/.claude/settings.json` thật** trong test: mọi test đặt `HOME` tạm; kiểm tra bằng `t.Setenv`.
- **Người dùng đã có pre-commit riêng** (husky, lefthook): `init` không ghi đè, chỉ in đoạn cần thêm; ghi hướng dẫn cho husky và lefthook trong README.
