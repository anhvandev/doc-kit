# Kiến trúc `dk`

Binary Go duy nhất nhúng template tài liệu, bảng loại, khung HTML và 11 skill. Không gọi
mạng lúc chạy, không dependency trên máy đích ngoài `git` (tùy chọn, để đếm dòng changelog).

## Ba lớp

```
binary dk (Go)
  ├── assets/templates/      template tài liệu, chỉ đọc, nhúng bằng embed
  ├── assets/types.toml      bảng loại: thư mục, tên file, mã, trường bắt buộc, trạng thái, from
  ├── assets/skills/         SKILL.md + references/rules.md của 11 skill
  ├── assets/html/           khung trang, CSS, mermaid.min.js
  ├── assets/hooks/          script pre-commit
  └── assets/agent-context.md  khối dán vào file ngữ cảnh agent (dk init --agent-context)
        ▲ gọi
.claude/skills/doc-*/        bản do dk cài ra (Codex: .codex/skills/)
```

| Lớp | Việc | Đổi thế nào |
|---|---|---|
| Template | cấu trúc và câu hỏi gợi ý của từng loại | sửa repo này, phát hành phiên bản mới; không tùy biến theo dự án |
| CLI | việc xác định, lặp lại được: tạo file, changelog, render, chỉ mục, kiểm tra, cài skill và hook | thêm lệnh khi việc không cần suy luận |
| Skill | việc cần đọc hiểu: phỏng vấn, soạn nội dung, phân tích tác động | skill chỉ gọi `dk`, không tự ghi file vào `docs/` |

Quy tắc chọn lớp: việc mà một script làm đúng 100 lần như nhau thì vào CLI; việc cần
đọc nội dung mới quyết được thì vào skill. `dk` không chạy được thì skill dừng và báo
người cài, không tạo file tay.

## Bố cục mã

```
cmd/dk/main.go        gọi internal/cli.Execute(), in lỗi ra stderr, thoát theo mã
internal/cli/         một file một nhóm lệnh (init, template, new, render, changelog, index,
                      check, refs, status, skill, hook, tokens, doctor, selfcheck)
internal/config/      dk.toml
internal/doctype/     registry loại tài liệu từ types.toml
internal/frontmatter/ tách, đọc, ghi metadata: YAML (.md), <!-- dk: --> (.html), $dk (.json), # dk: (.feature)
internal/tmpl/        render text/template từ assets/templates
internal/docs/        quét docs/, tạo file theo loại, chép trường --from, trích mục từ Feature Spec
internal/changelog/   đọc, thêm, gộp dòng CHANGELOG-DOCS.md
internal/gitx/        exec git: numstat, file đổi, hooks dir, nội dung tại HEAD
internal/render/      Markdown sang HTML tự chứa, chỉ mục, mã bước
internal/check/       16 quy tắc kiểm tra (xem docs/skill.md)
internal/refs/        đồ thị liên kết giữa file
internal/skill/       liệt kê, cài, gỡ, trạng thái skill; hash nội dung
internal/target/      interface Target; claude.go
internal/hook/        mục hook cần cài; `dk hook run` đọc stdin
internal/tokens/      tokens.json sang tokens.css
```

Thư viện: cobra, goldmark, yaml.v3, go-toml/v2. Test bằng `testing` chuẩn, `t.TempDir()`,
chạy `git` thật.

## Frontmatter chung

```yaml
id: F-012            # theo types.toml: seq, date hoặc none
type: feature-spec
title: Bộ lọc đơn hàng
status: draft        # tập trạng thái theo loại, khóa tiếng Anh
owner: ""
created: 2026-09-03
updated: 2026-09-03 14:05
source: CR-260910-loc-theo-trang-thai   # brief hoặc CR sinh ra hoặc sửa tài liệu này
created_by: dk       # dk check báo file thiếu
dk_version: 0.1.0
```

`dk new` điền toàn bộ; skill không sửa `created`, `created_by`, `dk_version`. Trường riêng
của loại khai trong `types.toml`.

## Changelog tài liệu

`docs/CHANGELOG-DOCS.md`, mới nhất trên, nhóm theo ngày, mỗi dòng
`HH:MM | file | +a −b | tóm tắt | nguồn`. Ba lớp bảo đảm: pre-commit chạy
`dk changelog pending` và chặn commit (cứng khi máy có `dk`); skill gọi `dk changelog add`
với tóm tắt thật; hook agent `PostToolUse` ghi dòng "chưa tóm tắt" để người hoặc skill thay
sau. Không ghi cho `docs/html/`, file `generated: true`, `plans/`.

## Target

| | `claude` | `codex` |
|---|---|---|
| Skill, scope dự án | `.claude/skills/` | `.codex/skills/` |
| Skill, `--global` | `~/.claude/skills/` | `$CODEX_HOME/skills/`, mặc định `~/.codex/skills/` |
| Hook | `.claude/settings.json` khóa `hooks`, matcher `Write`, `Edit\|Write` | `.codex/hooks.json` cùng cấu trúc khối `hooks`, matcher `apply_patch`; cần người dùng trust (`/hooks`) |
| Payload `hook run` | `tool_input.file_path` | `tool_input.command` chứa patch; đọc `*** Add File:` và `*** Update File:` |
| File ngữ cảnh agent | `CLAUDE.md` | `AGENTS.md` |

Hai file cấu hình hook dùng chung mã đọc, ghi và gỡ theo tiền tố `dk hook run`
(`internal/target/claude.go`); `codex.go` chỉ đổi đường dẫn và matcher. `--target` nhận
danh sách phân tách phẩy để cài một lệnh cho cả hai.

Nội dung skill không biết target. Dấu vết cài ghi trong frontmatter SKILL.md khóa
`metadata`: `dk_installed_by`, `dk_version`, `dk_hash` (sha256 mọi file, bỏ khối
metadata). Cài từ chối ghi đè skill đã sửa tay trừ `--force`; gỡ chỉ xóa skill có dấu vết.

## Vì sao repo này không dùng `dk` cho docs của chính nó

`docs/` của repo này là tài liệu thường, không có `dk.toml`. Dùng `dk` để viết tài liệu
cho `dk` tạo phụ thuộc vòng: sửa template hoặc lệnh là phải sinh lại tài liệu bằng chính
bản đang sửa, và lỗi ở CLI làm hỏng luôn tài liệu mô tả cách sửa lỗi đó. Ba file
`docs/kien-truc.md`, `docs/lenh.md`, `docs/skill.md` viết tay, cập nhật cùng commit đổi
hành vi.
