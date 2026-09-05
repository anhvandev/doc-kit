# dk

CLI Go nhúng template tài liệu và 12 skill tạo, duy trì tài liệu dự án. Một binary tĩnh,
không gọi mạng, không dependency lúc chạy. Tài liệu tiếng Việt, Markdown là nguồn sự
thật, HTML là bản sinh, mọi thay đổi trong `docs/` có dòng changelog trước khi commit.

- Kiến trúc ba lớp template, CLI, skill: [`docs/kien-truc.md`](docs/kien-truc.md)
- Bảng lệnh và hành vi từng lệnh: [`docs/lenh.md`](docs/lenh.md)
- Bảng skill và 16 quy tắc `dk check`: [`docs/skill.md`](docs/skill.md)
- Quy trình dùng skill, khi nào gọi gì, ai duyệt: [`docs/quy-trinh.md`](docs/quy-trinh.md)

## Cài binary

Từ release (Linux, macOS, Windows; amd64, arm64), thay phiên bản và nền tảng:

```sh
curl -sSL -o dk.tar.gz https://github.com/anhvandev/doc-kit/releases/download/v0.1.0/dk_0.1.0_linux_amd64.tar.gz
tar -xzf dk.tar.gz dk && install -m 0755 dk ~/.local/bin/dk   # hoặc thư mục khác trong PATH
dk --version
```

Windows: tải `dk_<phiên bản>_windows_amd64.zip`, giải nén `dk.exe` vào thư mục trong `PATH`;
pre-commit chạy qua `sh` của Git for Windows. Kiểm `checksums.txt` cùng release bằng
`sha256sum -c` trước khi cài. Hoặc bằng Go 1.26 trở lên:

```sh
go install github.com/anhvandev/doc-kit/cmd/dk@v0.1.0
```

Từ mã nguồn: `make build` ra `bin/dk`, `make install` vào `$GOBIN`. `dk self-check` in
phiên bản, số template, số skill và hash nội dung nhúng để so hai bản cài.

## Cài vào dự án

```sh
cd du-an
dk init                 # dk.toml, docs/, docs/CHANGELOG-DOCS.md, plans/, pre-commit; không đụng .claude/
dk skill install        # 12 skill vào .claude/skills/, dấu vết dk_* trong frontmatter
dk hook install         # 2 hook vào .claude/settings.json, giữ nguyên khóa khác; chạy lại thay mục dk cũ
dk init --agent-context # ghi khối quy tắc (tiếng Anh) vào CLAUDE.md và AGENTS.md giữa hai dấu mốc; chạy lại chỉ thay khối
dk doctor               # kiểm dk.toml, git, pre-commit, dk trên PATH, khối ngữ cảnh agent, skill đúng phiên bản, hook
```

Mặc định target `claude`, scope dự án; `--global` cài vào `~/.claude/`. Codex: thêm
`--target codex` (hoặc `--target claude,codex` cài cả hai trong một lệnh): skill vào
`.codex/skills/`, hook vào `.codex/hooks.json` với matcher `apply_patch`; scope toàn máy
theo `$CODEX_HOME`, mặc định `~/.codex/`. Codex chỉ chạy hook cấp dự án sau khi người
dùng trust: mở Codex trong dự án, gõ `/hooks` và duyệt. `dk doctor` tự kiểm thêm Codex
khi dự án có `.codex/`. Cùng bộ skill, không đổi nội dung theo target. `git` cần có để
đếm dòng changelog và chạy pre-commit; không có git vẫn chạy, changelog ghi "không git".
Pre-commit và hook gọi `dk` qua PATH và cho qua có chủ đích khi thiếu `dk`; `dk doctor`
báo mục này.

Rồi mở harness agent trong dự án: skill `doc-*` gọi `dk new` lấy khung, điền nội dung,
`dk changelog add` khi xong, dừng ở trạng thái chờ duyệt để người chốt.

## Gỡ

```sh
dk hook uninstall       # chỉ xóa mục có command bắt đầu bằng `dk hook run`
dk skill uninstall      # chỉ xóa skill có dấu vết dk; skill tự viết giữ nguyên
rm "$(git rev-parse --git-path hooks)/pre-commit"
```

`docs/` đọc được không cần CLI; `docs/html/` là bản sinh, xóa được. Gỡ xong mới gỡ binary.

## Phát triển

```sh
make test          # go test ./...
make vet
make lint-skills   # assets/skills/ không nhắc ak-, Claude Code, .claude/, .codex/, Edit, Write; mọi .md dưới 300 dòng
make build-all     # 6 binary vào dist/ (linux, darwin, windows × amd64, arm64)
```

CI (`.github/workflows/ci.yml`) chạy vet và test trên ubuntu, macos, windows (lint-skills
chỉ trên ubuntu và macos), và build 6 nền tảng. Push tag `v*` chạy goreleaser
(`.goreleaser.yaml`): archive `tar.gz` (Windows: `zip`), `checksums.txt`. Docs của repo này viết tay, không dùng `dk`, để tránh phụ thuộc
vòng; lý do ở `docs/kien-truc.md`.

## Giới hạn đã biết (`v0.2.0`)

- Hook Codex: cấu trúc `hooks.json` và payload `apply_patch` (`tool_input.command` chứa
  patch, `*** Add File:` / `*** Update File:`) theo tài liệu hooks và mã nguồn Codex
  0.153.2; chưa kiểm trong phiên Codex đã đăng nhập. Hook không phải lớp bảo đảm: pre-commit
  và skill mới là lớp cứng; hook chỉ chặn khi Codex bắt `apply_patch` và hook đã được trust.
- Windows: build và test qua CI; chưa chạy tay trên máy Windows thật. Hook Claude Code trên
  Windows chạy qua Git Bash (mặc định) hoặc PowerShell, lệnh `dk hook run` không phụ thuộc
  shell.
- `dk changelog pending` so theo giờ phút của commit HEAD: sửa, commit, sửa tiếp trong cùng
  một phút không bị báo cho đến phút kế.
- `adr-immutable` chỉ so với HEAD; thân ADR đã chốt bị sửa và đã commit thì không bắt được.
- Hook agent chỉ chặn tạo `.md` mới trong `docs/`; `.html`, `.json`, `.feature` không chặn.
  Cả hai target đều không bắt ghi file qua lệnh shell (`Bash`, `exec`); Codex chỉ bắt
  `apply_patch`. Pre-commit là lớp cứng.
- Template không tùy biến theo dự án; nhu cầu đổi đi qua phát hành phiên bản mới.
