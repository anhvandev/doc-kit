# Phase 11: Target Codex và Windows, kết quả kiểm

Ngày: 2026-09-04. Codex CLI cài trên máy phát triển: `codex-cli 0.153.2` (npm `@openai/codex`).

## Định dạng hook Codex (đóng mục UNVERIFIED số 8 của validate phiên 1)

Nguồn: tài liệu hooks Codex (`learn.chatgpt.com/docs/hooks`, bản hiện hành) và mã nguồn
`codex-rs/core/src/tools/handlers/apply_patch.rs` (nhánh main).

- `hooks.json` tại `~/.codex/hooks.json` hoặc `<repo>/.codex/hooks.json`; cấu trúc
  `{"hooks": {"PreToolUse": [{"matcher": "...", "hooks": [{"type": "command", "command": "..."}]}]}}`.
  **`command` là chuỗi**, không phải mảng như phase 11 dự kiến; cùng cấu trúc khối `hooks`
  của Claude Code. Có thêm `commandWindows` (tùy chọn, không dùng).
- Sự kiện: `PreToolUse`, `PostToolUse` bắt shell, `unified_exec`, **`apply_patch`**, MCP.
  `apply_patch` phát hook từ bản sửa issue #16732 (0.118 chưa có; 0.153.2 có).
- Payload: `tool_name: "apply_patch"`, `tool_input.command` là toàn bộ patch (`*** Begin Patch`,
  `*** Add File: <path>`, `*** Update File: <path>`, `*** Delete File:`); `cwd` là gốc phiên.
  Matcher `Edit`/`Write` là bí danh của `apply_patch` theo tài liệu; `dk` dùng thẳng `apply_patch`.
- Deny: cùng JSON `hookSpecificOutput.permissionDecision: "deny"` (hoặc exit 2 + stderr).
- Trust: hook cấp dự án chỉ chạy sau khi người dùng duyệt qua `/hooks` trong Codex; hook
  bật mặc định (`[features] hooks`, bí danh cũ `codex_hooks`).

Payload trong `internal/hook/testdata/codex-pre-tool-use-apply-patch.json` dựng theo tài
liệu và mã nguồn, **không phải ghi từ phiên Codex thật** (máy chưa đăng nhập, xem dưới).

## Đã làm

- `internal/target/codex.go`: `SkillDir` (`.codex/skills`, `$CODEX_HOME/skills`, mặc định
  `~/.codex/skills`), `HooksPath`, `InstallHooks` (matcher `apply_patch`), `UninstallHooks`,
  `InstalledHooks`; dùng chung mã đọc/ghi/gỡ với Claude (tách thành `installHooksFile`,
  `uninstallHooksFile`, `installedHooksFile` trong `claude.go`). Không ghi `config.toml`.
- `internal/hook/run.go`: đọc `tool_input.file_path`, dự phòng `tool_input.path`, và với
  `apply_patch` đọc `*** Add File:` / `*** Update File:` trong `tool_input.command`; mỗi
  file xử lý riêng, deny in một lần.
- `--target` nhận danh sách phẩy cho `skill install|uninstall|status`, `hook install|uninstall`,
  `doctor`; `hook install --target codex` in nhắc trust; `doctor` tự thêm codex khi có `.codex/`;
  nhãn dòng doctor đổi thành `skill (claude, dự án)`, `hook (codex, dự án)`, `skill doc-x (claude)`.
- Windows: Makefile 6 nền tảng (`.exe`), goreleaser `windows` (zip), CI `windows-latest`
  (bỏ `lint-skills` vì thiếu make/awk GNU), `doctor` không kiểm bit chạy trên Windows,
  `gitx.NumStat` pathspec dạng slash, test đặt `USERPROFILE` cùng `HOME`, test `sh` và bit
  chạy bỏ qua trên Windows.
- README, `docs/lenh.md`, `docs/kien-truc.md`.

## Kiểm

| Việc | Kết quả |
|---|---|
| `go vet ./...`, `go test ./...` | qua |
| `GOOS=windows go vet ./...`, `GOOS=windows go test -c` mọi package | biên dịch qua |
| `make build-all` | 6 binary, gồm `dk-windows-amd64.exe`, `dk-windows-arm64.exe` |
| `make lint-skills` | sạch; không file nào trong `assets/skills/` đổi |
| Dự án mẫu `/tmp/dk-codex-probe`: `skill install --target claude,codex`, `hook install --target codex`, `doctor` | 22 skill cài, `hooks.json` đúng golden, doctor OK cho codex |
| `codex exec` trong dự án mẫu | Codex đọc cấu hình dự án không báo lỗi `hooks.json`, dừng ở 401 vì chưa đăng nhập |

## Chưa kiểm được

- Vòng Tiếp nhận trong Codex thật, hook deny khi `apply_patch` tạo file mới trong `docs/`,
  changelog sau `apply_patch` sửa file: cần `codex login`. Bước 5 của phase chuyển sang khi
  có tài khoản. Payload mẫu trong testdata cần thay bằng payload ghi từ phiên thật khi có.
- Windows máy thật: `dk init`, `dk skill install`, pre-commit qua Git for Windows; CI
  `windows-latest` chưa chạy vì repo chưa có remote.
- Release `v0.2.0`: chặn cùng lý do `v0.1.0` (remote, module path).

## Sau review (`plans/reports/code-review-260904-phase-11-target-codex.md`)

- Sửa: `patchFiles` đọc thêm `*** Move to:` (đổi tên vào `docs/` là tạo tài liệu mới, có
  test); `doctor` chỉ tự kiểm codex khi có `.codex/skills/` hoặc `.codex/hooks.json`
  (thư mục `.codex/` rỗng sau khi gỡ không kích hoạt, có test); `gitx.run` cắt thêm `\r`;
  README ghi hook cả hai target không bắt ghi file qua shell; `docs/lenh.md` ghi nhãn dòng
  doctor đổi từ `v0.2.0`; phase file sửa chữ `command` mảng thành chuỗi.
- Giữ theo quyết định cũ: `doctor --target lạ` vẫn in bảng và mã 3 (1.16, M2 phase 10);
  `hook install` in số hook cố định dù không đổi (nit chấp nhận ở 1.9).
- M-2 (mtime `assets/skills/` 10:05) không phải thay đổi của phase này: phiên bắt đầu 10:27,
  `find -newermt 10:27` trả 0; các file đó đổi ở phase 10 (plan ghi 10:15 xong).
- Chấp nhận, ghi nhận: `doctor --global` vẫn dò `.codex/` trong dự án để quyết định kiểm
  codex; target chưa cài gì in 11 dòng skill thiếu; `*** Begin Patch` trong `command` của
  tool khác cũng được đọc (vô hại với matcher hiện tại).
