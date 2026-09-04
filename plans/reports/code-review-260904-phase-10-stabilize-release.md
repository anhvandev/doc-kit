# Code review phase 10: ổn định và phát hành

- Ngày: 2026-09-04
- Phạm vi: `internal/cli/{doctor,selfcheck,init,root}.go`, `internal/target/{target,claude}.go`,
  `internal/cli/phase10_test.go`, `assets/agent-context.md`, `assets/embed.go`, `Makefile`,
  `.github/workflows/{ci,release}.yml`, `.goreleaser.yaml`, `README.md`, `docs/*.md`
- Repo chưa có commit nên không dùng `git diff`; đọc toàn bộ file.

## Kiểm tra đã chạy

| Lệnh | Kết quả |
|---|---|
| `gofmt -l .` | không có file |
| `go vet ./...` | sạch |
| `go test ./...` | mọi package `ok` |
| `make lint-skills` | `lint-skills: sạch` |

## Success Criteria phase 10

| Tiêu chí | Trạng thái |
|---|---|
| `dk doctor` phát hiện thiếu pre-commit, skill cũ phiên bản, hook thiếu và nêu cách sửa | đạt, `phase10_test.go` phủ cả bốn trường hợp |
| Agent context dưới 60 dòng | đạt, `assets/agent-context.md` 45 dòng, test chặn ngưỡng |
| `make lint-skills` bằng không, CI ubuntu và macos | đạt (xem M3, L3 về CI) |
| goreleaser 4 nền tảng, checksums | cấu hình đủ; chưa chạy thật goreleaser trong review này |
| `doc-smoke` đã xóa, test cài dùng skill thật | đạt, `assets/skills/` còn 11 skill, test dùng `doc-adr`, `doc-cr` |
| Toàn vòng chạy, report có số đo | **chưa đạt trong repo**: không có `plans/260903-1400-skill-tai-lieu/reports/phase-10-full-run.md` |

## Không hồi quy

- `dk init` không cờ: nhánh `--agent-context` nằm trước `runInit`, chỉ chạy khi cờ bật; chạy thử
  `dk init` vẫn tạo `dk.toml`, cây `docs/`, changelog, `plans/`, pre-commit như cũ.
- `Target` thêm `InstalledHooks`: grep toàn repo chỉ có `*target.Claude` hiện thực interface,
  không có fake trong test. Không vỡ implementer nào.
- Hợp đồng công khai: chỉ thêm hai lệnh (`doctor`, `self-check`) và một cờ (`init --agent-context`);
  mã thoát 0/1/2/3 khớp `docs/lenh.md`; JSON của lệnh cũ không đổi.
- Docs so với mã: 16 quy tắc (`internal/check/check.go` `Rules` có 16 mục) ✓; 41 loại
  (`assets/types.toml` 41 bảng) ✓; 11 skill ✓; đường dẫn tải trong README
  (`dk_0.1.0_linux_amd64.tar.gz`) khớp `name_template` của goreleaser ✓.

## Finding

### High

**H1. Chuỗi phiên bản của `make` và của release lệch nhau, làm mọi skill bị coi là cũ**
`Makefile:1` lấy `git describe --tags` nên `Version = "v0.1.0"` (có `v`), còn
`.goreleaser.yaml:13` dùng `{{ .Version }}` và `internal/cli/root.go:30` cắt `v` khỏi build info
nên binary phát hành có `Version = "0.1.0"`. Kiểm chứng:

```
$ /tmp/dkA --version   # build bằng ldflags v0.1.0, đã dk skill install
v0.1.0
$ /tmp/dkB doctor      # build bằng ldflags 0.1.0
!! skill doc-adr  | cũ (vv0.1.0)  | `dk skill install doc-adr` để lên 0.1.0
```

Hệ quả: người `make install` rồi nâng lên binary release (hoặc ngược lại) thấy cả 11 skill "cũ",
`doctor` trả mã 3; frontmatter ghi `dk_version: v0.1.0` khác mẫu ở plan mục 1.3; chuỗi hiển thị
`vv0.1.0` do `internal/skill/status.go:63` tự thêm `v` vào giá trị đã có `v`.
Sửa: chuẩn hóa một chỗ, gọn nhất là bỏ điều kiện trong `root.go` `init()` và luôn
`Version = strings.TrimPrefix(Version, "v")`; hoặc `Makefile:1` thêm `| sed 's/^v//'`.

### Medium

**M1. `dk doctor` chỉ nhìn scope dự án, báo sai khi người cài `--global`**
`internal/cli/doctor.go:91` và `:119` truyền cứng `false`. Chạy thử với
`dk skill install --global` và `dk hook install --global` trong dự án sạch: doctor báo cả 11 skill
"chưa cài" và "hook có 0/2", cách sửa lại khuyên cài scope dự án, tức đẩy người dùng sang cấu hình
khác với cái họ đã chọn. Sửa: `f.bind(cmd, true)` rồi truyền `f.global` xuống, hoặc kiểm cả hai
scope và đạt khi một trong hai đủ.

**M2. Lỗi cứng ở giữa `doctor` làm mất cả bảng đã gom**
`internal/cli/doctor.go:87-93`: `target.Get` lỗi trả mã 2, `skill.Status` lỗi trả mã 1, cả hai bỏ
luôn các dòng đã thu (`dk.toml`, `docs/`, git, pre-commit). `dk doctor --target codex` vì thế chỉ in
`target "codex" chưa hỗ trợ` chứ không phải bảng chẩn đoán. Sửa: thêm một dòng
`add("skill", err.Error(), false, ...)` rồi `return a.printDoctor(rows)` thay vì `fail`.

**M3. `ci.yml` không giới hạn quyền token**
`.github/workflows/ci.yml` không có khối `permissions`, nên dùng mặc định của repo, có thể là
`contents: write` với repo cũ. `release.yml` đã đặt `contents: write` đúng phạm vi. Sửa: thêm
`permissions:\n  contents: read` ở cấp workflow trong `ci.yml`. Cân nhắc ghim action theo SHA
(`actions/checkout@<sha>`) cho cả hai workflow.

### Low

**L1. Dự án không có git luôn trả mã 3**
`internal/cli/doctor.go:76-79` đánh dấu cả `git` và `pre-commit` là chưa đạt. README lại nói "không
có git vẫn chạy". Kết quả: dự án cố tình không dùng git không bao giờ có `doctor` xanh. Sửa: hoặc
để hai mục này ở dạng ghi chú đạt-có-điều-kiện, hoặc sửa README cho khớp rằng git là bắt buộc.

**L2. `self-check` ghi cứng danh sách target**
`internal/cli/selfcheck.go:52` `Targets: []string{"claude"}`. Khi phase 11 thêm `codex`, chỗ này
lặng lẽ sai. Sửa: đưa danh sách vào `internal/target` (ví dụ `target.Names()`) và dùng chung với
`registry.go`.

**L3. `lint-skills` dựa vào `\b` của GNU grep**
`Makefile:20` dùng `grep -rnE '...\bEdit\b...'`. Trên `macos-latest` (BSD grep) `\b` không được bảo
đảm, có thể không khớp gì và job vẫn xanh dù có từ cấm. Job ubuntu vẫn bắt được nên rủi ro chỉ là
mất lớp kiểm trùng. Sửa: đổi sang `[^[:alnum:]]Edit[^[:alnum:]]` hoặc chỉ chạy lint-skills trên
ubuntu.

**L4. `archives.formats` cần goreleaser ≥ 2.6**
`.goreleaser.yaml:17` dùng khóa số nhiều `formats`. `version: "~> v2"` hiện lấy bản mới nhất nên
chạy được, nhưng nếu ghim bản v2 cũ hơn thì hỏng lúc phát hành. Ghi rõ ràng buộc hoặc ghim
`version: "~> v2.6"`.

**L5. Repo chưa có LICENSE, archive chỉ kèm README**
`.goreleaser.yaml:19` `files: [README.md]`. Repo không có file LICENSE nào. Binary nhúng
`mermaid.min.js` (MIT); giấy phép mermaid được ghi vào HTML sinh ra (`internal/render/render.go:78`)
nên nghĩa vụ MIT coi như đã đáp ứng ở bản HTML, nhưng bản phát hành binary vẫn thiếu giấy phép của
chính `dk`. Sửa: thêm LICENSE và đưa vào `files`.

**L6. Đếm hook chỉ so `<`**
`internal/cli/doctor.go:121`: `len(cmds) < want`. Cấu hình có hook `dk hook run` trùng lặp hoặc sự
kiện lạ vẫn báo "đủ 3 hook". Đủ dùng, chỉ nêu để biết.

### Nit

- `internal/cli/init.go:32-35`: `--agent-context` kèm `--force` thì `--force` bị bỏ im lặng. Có thể
  `cmd.MarkFlagsMutuallyExclusive("force", "agent-context")`.
- `internal/cli/doctor.go:86`: gọi `target.Get` trực tiếp thay vì `a.resolveTarget(f)` như
  `skill.go`, `hook.go`. Hợp lý vì gốc dự án đã nạp, nhưng lệch pattern; một dòng chú thích là đủ.
- `internal/cli/doctor.go:158-163`: bỏ qua giá trị trả về của `w.Write`; lỗi chỉ lộ ở `Flush`.
- `internal/cli/doctor.go:69`: `exec.LookPath("dk")` không xác nhận binary tìm được đúng là `dk`
  hay cùng phiên bản với binary đang chạy.

## Ghi nhận

- `embedHash` (`internal/cli/selfcheck.go:82-104`) băm cả đường dẫn kèm byte 0 phân tách, nên không
  nhập nhằng giữa hai cây file khác nhau. Đúng cách.
- `doctor` xử lý được `dk.toml` hỏng (dòng lỗi + mã 3) và `settings.json` hỏng (`InstalledHooks` trả
  lỗi thành một dòng của bảng), không panic.
- `printAgentContext` không ghi file, có test khẳng định không sinh `dk.toml`.

## Câu hỏi còn treo

1. `doctor` nên bao gồm scope toàn máy mặc định, hay thêm cờ `--global` để người tự chọn (M1)?
2. Repo có coi git là bắt buộc không (L1)? README và `doctor` đang nói khác nhau.
3. Report toàn vòng `phase-10-full-run.md` đang do người khác chạy, hay còn thiếu?
