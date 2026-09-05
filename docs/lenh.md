# Lệnh `dk`

Mọi lệnh nhận `--json` và `--cwd <dir>`. Mã thoát: 0 thành công, 1 lỗi I/O hoặc logic,
2 sai cờ hoặc tham số, 3 kiểm tra không qua.

| Lệnh | Việc |
|---|---|
| `dk --version` | phiên bản; `go install` lấy từ build info |
| `dk self-check` | phiên bản, số template, số skill, target, sha256 nội dung nhúng; mã 1 khi template hoặc skill nhúng hỏng |
| `dk doctor [--target claude,codex] [--global]` | bảng `mục \| trạng thái \| cách sửa`: `dk.toml`, `docs/`, `dk` trên PATH, git, pre-commit, khối ngữ cảnh agent trong `CLAUDE.md` và `AGENTS.md` (thiếu file; có file chưa có khối; cũ hoặc bị sửa tay; khối thiếu mốc đóng hoặc nhiều hơn một khối), từng skill sai phiên bản, hook theo từng target (đủ, thiếu, hoặc `lệch bản hiện tại` khi matcher hay lệnh khác bản nhúng; nhãn dòng `skill (claude, dự án)`, `hook (codex, toàn máy)`, `skill doc-x (claude)`; đổi từ `v0.2.0`); không nêu `--target` mà có `.codex/skills/` hoặc `.codex/hooks.json` thì kiểm cả Codex; target lạ vẫn in bảng, mã 3 |
| `dk init [--force]` | `dk.toml`, cây `docs/`, `docs/CHANGELOG-DOCS.md`, `plans/`, pre-commit khi có git; `--force` giữ `dk.toml` cũ, chỉ bổ sung; từ chối tạo dự án lồng nhau |
| `dk init --agent-context` | ghi `assets/agent-context.md` (tiếng Anh, dưới 120 dòng) vào `CLAUDE.md` và `AGENTS.md` ở gốc dự án (có `dk.toml`) hoặc thư mục hiện tại giữa `<!-- dk:agent-context start version=… hash=… -->` và `<!-- dk:agent-context end -->`; file chưa có thì tạo, có khối thì thay tại chỗ, chưa có khối thì nối cuối; phần ngoài khối giữ nguyên; in `created`, `updated`, `unchanged` từng file; không cần `dk.toml` |
| `dk template list \| show <loại>` | bảng loại (42) và template thô |
| `dk new <loại> <slug> [--from <file>] [--set k=v] [--force]` | tạo file từ template, điền frontmatter, chép trường theo `from` của `types.toml` |
| `dk new plan-phase \| report <slug> --in <thư mục plan>` | file trong `plans/` (`plans_dir` của `dk.toml`); `--in plans` cho `plans/reports/` |
| `dk new decision-log \| faq --append "<chữ>"` | nối một dòng `- <ngày> \| <chữ>` vào cuối file, tạo từ template khi chưa có |
| `dk new release-notes --collect <phiên bản>` | gom Release brief `ready` chưa có `released_in`, ghi `released_in` vào từng brief |
| `dk changelog add <file> --summary <s> [--source <src>]` | bump `updated:`, đếm `git diff --numstat HEAD`, ghi dòng; gộp với dòng cùng file cùng nguồn trong 10 phút |
| `dk changelog pending` | mã 1 khi còn file trong `docs/` đổi mà chưa có dòng từ phút của commit HEAD |
| `dk render <file> \| --all [--index]` | HTML tự chứa vào `docs/html/`; cột phải có bảng metadata và mục lục (tiêu đề cấp 2, 3 ở mức đầu tài liệu, bỏ tiêu đề trong danh sách hay trích dẫn; chỉ hiện khi có từ 3 tiêu đề; trang chỉ mục không có); `--index` thêm `docs/html/index.html` |
| `dk index [features\|adr\|cr\|intake\|user-guide\|all]` | `README.md` chỉ mục (`generated: true`) trong từng thư mục |
| `dk check [<file>] [--strict]` | 16 quy tắc (xem `docs/skill.md`); `--strict` coi warning là lỗi |
| `dk refs <file>` | liên kết đi và đến |
| `dk status` | tài liệu theo loại × trạng thái, CR mở, changelog pending, finding, DR quá 6 tháng |
| `dk skill list \| install [<tên>...] \| uninstall [<tên>...] \| status` | cờ `--target claude` (`codex`, hoặc `claude,codex` cho cả hai), `--global`, `--force` |
| `dk hook install \| uninstall` | 2 hook mỗi target vào cấu hình target (`--target` như trên); `install` bỏ mọi lệnh `dk hook run` đang có rồi ghi bản hiện tại, nên chạy lại sau khi nâng `dk` là đủ, không nối thêm mục trùng; Codex in nhắc trust (`/hooks`); `dk hook run` ẩn, do harness gọi, đọc `tool_input.file_path` (Claude) hoặc patch `apply_patch` trong `tool_input.command` (Codex) |
| `dk tokens css [--in] [--out]` | `tokens.json` sang `tokens.css` |

## `new`: quy ước theo loại

- Mã `seq` luôn tăng, trùng slug không đụng file; `--force` chỉ có nghĩa với loại `none`
  và `date`.
- `--from` trong `intake/` đặt file mới cùng thư mục nguồn; `interview --from <CR>` đặt
  `cr/<CR>/interview.md`. `source:` là `id` nguồn, không có id thì `<thư mục>/<file>`.
- `--set k=true|false` ghi YAML bool. `--set format=spec|use-case|story|crud|state` và
  `--set has_ui=false` cho `feature-spec`.
- Tên file có placeholder `{feature}`, `{step}`, `{layer}`, `{version}`: chép từ `--from`
  hoặc `--set`; thiếu là lỗi nêu cờ cần thêm. `step` khớp `^B\d+[a-z]?$`.
- `test-case --from <Feature Spec>`: mỗi `ACn` mục 8 một `Scenario`, `Background` từ mục 2;
  AC không theo Given / When / Then giữ dưới `# chưa tách được`.
- `release-brief --from <Feature Spec>`: chép nháp mục 1, 2, cột hành động mục 4, mockup
  mục 5, mục 10 để skill viết lại bằng ngôn ngữ người dùng.
- `postmortem --set incident_at="yyyy-mm-dd hh:mm"`: `dk new` tính `written_within_48h`.

## Changelog

`add` bump `updated:` trước khi đếm nên số dòng khớp git sau khi lệnh chạy. File chưa trong
HEAD ghi `mới, N dòng`; không git ghi `không git, N dòng`. `pending` chỉ so giờ phút, nên sửa
rồi commit rồi sửa tiếp trong cùng phút không bị báo; lần chạy sau bắt lại. Hook
`post-edit` ghi dòng "chưa tóm tắt", nguồn `trực-tiếp`, không bump `updated:`.

## Render và chỉ mục

`docs/html/<đường dẫn>.html` giữ cây; liên kết `.md` trong `docs/` đổi sang `.html`, liên
kết ra `plans/` trỏ file gốc. HTML thô trong Markdown bị escape. `mermaid.min.js` (11.17.2,
MIT) chỉ nhúng khi trang có khối ` ```mermaid `, trang đó nặng khoảng 3,5 MB, mở từ
`file://`. Mã bước trong sơ đồ thành dòng "Bước:" liên kết đến bảng hành vi. `index` không
ghi `updated` để file không đổi mỗi lần chạy; `features` và `adr` là bảng phẳng theo mã,
`cr` và `intake` nhóm theo trạng thái, `user-guide` nhóm theo `task`.

## Design tokens

`docs/design/tokens/tokens.json` theo W3C Design Tokens, metadata ở khóa `$dk`. `dk tokens
css` đọc nhóm lồng, `$type` kế thừa, alias `{a.b.c}` kể cả trong chuỗi, theme qua
`$extensions.dk.theme.<tên>` thành `[data-theme="<tên>"]`; alias vòng hoặc lạ là lỗi.
`tokens.css` là bản sinh (dòng đầu `/* generated: dk tokens css */`), không ghi changelog.
Mockup chỉ dùng `var(--...)`; Figma qua `--set external=<url>`.

## Hook agent

`PreToolUse` cho `Write`: từ chối tạo file `.md` mới trong `docs/` (JSON
`permissionDecision: "deny"`). `PostToolUse` cho `Edit|Write`: ghi dòng changelog giữ chỗ.
Cả hai gọi `dk hook run` qua PATH, fail-open: thiếu `dk` hay lỗi nội bộ chỉ in cảnh báo ra
stderr, thoát 0. `hook uninstall` chỉ bỏ lệnh có tiền tố `dk hook run`, giữ lệnh khác
cùng mục; `settings.json` rỗng thì xóa file.

## Pre-commit

`dk init` chép script vào thư mục hooks (`git rev-parse --git-path hooks`, đúng với worktree)
khi chưa có; đã có thì in đoạn cần thêm (husky: `.husky/pre-commit`; lefthook: một
`commands` chạy `dk changelog pending`). Script bỏ qua có chủ đích khi máy không có `dk`.
Gỡ: `rm "$(git rev-parse --git-path hooks)/pre-commit"`.
