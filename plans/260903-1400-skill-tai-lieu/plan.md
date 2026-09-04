---
title: "dk: CLI Go nhúng template và bộ skill tạo, duy trì tài liệu dự án"
description: "Binary Go duy nhất nhúng template tài liệu và 11 skill; cài/gỡ skill và hook theo target (Claude trước, Codex sau); mọi thao tác máy làm được trên docs/ đi qua CLI; changelog tài liệu bắt buộc"
status: in-progress
priority: P1
effort: "6w"
tags: [go, cli, skills, docs, claude-code, codex]
created: 2026-09-03
blockedBy: []
blocks: []
---

# dk: CLI Go nhúng template và bộ skill tạo, duy trì tài liệu dự án

- Ngày: 2026-09-03
- Cập nhật: 14:20 thêm CLI làm cửa duy nhất tạo tài liệu; 14:25 CLI Go, template và skill nhúng trong binary, bỏ `docs/_templates/`; 14:45 tách thành 11 phase chi tiết; 15:05 validate phiên 1: `init` không tự cài skill, Windows lùi v0.2.0, mockup HTML mặc định; 16:00 phase 1 xong, ghi quyết định triển khai mục 1.7; 16:25 phase 2 xong, quyết định mục 1.8; 16:50 phase 3 xong (trừ bước 9 kiểm thật trong Claude Code), quyết định mục 1.9; 17:05 phase 4 xong (trừ bước 7 chạy thật trong Claude Code), quyết định mục 1.10; 17:35 phase 5 xong (phần AI thật gộp vào phase 10), quyết định mục 1.11; 18:00 phase 6 xong (phần AI thật gộp vào phase 10), quyết định mục 1.12; 22:55 phase 7 xong (phần AI thật gộp vào phase 10), quyết định mục 1.13; 2026-09-04 09:10 phase 8 xong (phần AI thật và dry-run BDD gộp vào phase 10), quyết định mục 1.14; 09:40 phase 9 xong (người không kỹ thuật đọc brief và runbook trên sự cố thật gộp vào phase 10), quyết định mục 1.15; 10:15 phase 10 xong trừ tag `v0.1.0` (chờ remote và module path), quyết định mục 1.16; 10:50 phase 11 xong trừ kiểm thật trong Codex đã đăng nhập, Windows máy thật và tag `v0.2.0`, quyết định mục 1.17; 2026-09-04 16:40 `--agent-context` ghi file thay vì in, quyết định mục 1.18
- Trạng thái: mã xong cả 11 phase; chờ remote để tag `v0.1.0`, `v0.2.0`, chạy CI Windows; chờ `codex login` để kiểm hook Codex thật
- Nguồn: `plans/reports/consult-260903-1032-tai-lieu-du-an-nguoi-va-ai.md` (hệ thống tài liệu 10 tầng). Mọi quy tắc nội dung tài liệu lấy từ báo cáo này; plan không chép lại.

## Hợp đồng

- **Kết quả**: một binary `dk` viết bằng Go, nhúng bộ template tài liệu và bộ skill; cài và gỡ skill, hook theo target (`claude`, `codex`) và scope (dự án, toàn máy); thực hiện mọi thao tác máy làm được trên tài liệu (tạo từ template, render HTML, changelog, chỉ mục, kiểm tra); skill là lớp mỏng gọi CLI thay vì tự ghi file; mỗi tài liệu có bản HTML tự chứa; mọi thay đổi tài liệu đều ghi changelog với file, ngày giờ, tóm tắt, số dòng thêm bớt.
- **Ràng buộc**: Go, một binary tĩnh, không dependency runtime trên máy đích, không gọi mạng lúc chạy; CLI và skill là mã mới, độc lập hoàn toàn với các skill `ak-*` và script của chúng; CLI chạy được không cần Claude Code; template chỉ tồn tại trong binary, không chép ra repo dự án; SKILL.md theo định dạng chung (frontmatter `name`, `description`), thân skill trung lập target (không nhắc Claude Code, tên tool `Edit`/`Write`, đường dẫn `.claude/`); Claude Code hoàn chỉnh trước, Codex là target thứ hai dùng cùng nội dung skill; Markdown là nguồn sự thật, HTML là bản sinh; tài liệu tiếng Việt có dấu; tên file Go theo snake_case, package lowercase.
- **Ngoài phạm vi**: target khác ngoài Claude và Codex; site tĩnh; đồng bộ Notion; sinh mockup bằng AI; lint code UI (việc của dự án đích); giao diện web cho CLI; tùy biến template theo từng dự án; sinh Schema và API spec từ code.
- **Nghiệm thu**: chạy trọn một vòng ý tưởng mới đến Feature Spec đến CR chỉ bằng skill gọi CLI; `dk skill install` và `dk hook install` rồi `uninstall` cả hai trả `.claude/` về đúng trạng thái cũ, không đụng skill người dùng tự viết; rút `dk` khỏi PATH thì skill dừng có thông báo, không sinh file nào vào `docs/`; mọi file trong `docs/` có dòng changelog trước khi commit, số dòng khớp `git diff --numstat`; HTML mở từ đĩa không cần mạng; `go test ./...` và `go vet ./...` qua ở mọi phase.

Quyết định validate phiên 1: `dk init` chỉ tạo `docs/`, `dk.toml`, changelog, pre-commit và in nhắc hai lệnh tiếp theo; cài skill và hook là việc riêng của `dk skill install` và `dk hook install` (mặc định target `claude`, scope dự án, cờ `--global`).

## 1. Quyết định kiến trúc dùng chung

### 1.1 Ba lớp

```
binary dk (Go)
  ├── assets/templates/      ← template tài liệu, chỉ đọc, nhúng bằng embed
  ├── assets/skills/         ← SKILL.md + references/rules.md cho 11 skill
  ├── assets/html/           ← khung HTML, CSS, mermaid.min.js, quy tắc html-style
  └── lệnh                   ← init, template, new, render, changelog, index, check, refs, status, skill, hook
        ▲ gọi
.claude/skills/doc-*/        ← bản do dk cài ra; Codex: .codex/skills/doc-*/
```

- **Template**: cấu trúc và câu hỏi gợi ý. Sửa template là sửa hợp đồng tài liệu: sửa repo CLI, người duyệt, phát hành phiên bản mới.
- **CLI**: việc xác định, lặp lại được, không suy luận. Việc cần đọc hiểu nội dung (phân tích tác động, soạn câu hỏi, viết tóm tắt) không đưa vào CLI.
- **Skill**: việc cần suy luận. Không tự ghi file vào `docs/`; gọi `dk new` lấy khung, điền nội dung, gọi `dk changelog add` khi xong. `dk` không chạy được thì dừng và báo người cài, không tạo file tay.

### 1.2 Bố cục repo Go

```
go.mod                         module path chốt theo remote thật (placeholder github.com/vner/dk)
cmd/dk/main.go                 gọi internal/cli.Execute()
internal/cli/                  một file một nhóm lệnh: root.go init.go template.go new.go render.go
                               changelog.go index.go check.go refs.go status.go skill.go hook.go
internal/config/               đọc dk.toml
internal/doctype/              registry loại tài liệu, đọc từ assets/types.toml
internal/frontmatter/          tách, parse, ghi lại YAML frontmatter, giữ nguyên thân
internal/tmpl/                 render text/template từ assets/templates
internal/docs/                 tạo file theo loại (new), sinh mã, chép trường --from
internal/changelog/            đọc, thêm, gộp dòng CHANGELOG-DOCS.md
internal/gitx/                 exec git: numstat, changed files, untracked
internal/render/               goldmark sang HTML, khung, chỉ mục
internal/check/                liên kết, mã bước, liên kết ngược, ngưỡng dòng
internal/refs/                 đồ thị liên kết giữa file
internal/skill/                embed skill, install, uninstall, status, dấu vết
internal/target/               interface Target; claude.go, codex.go
internal/hook/                 ghi/gỡ cấu hình hook theo target; hook run đọc stdin
assets/                        embed.go khai báo //go:embed
```

Thư viện: `spf13/cobra` (lệnh và cờ), `yuin/goldmark` (Markdown), `gopkg.in/yaml.v3` (frontmatter), `pelletier/go-toml/v2` (dk.toml, types.toml). Không thêm thư viện khác nếu stdlib làm được. Test bằng `testing` chuẩn, `t.TempDir()`, chạy `git` thật trong thư mục tạm.

### 1.3 Frontmatter chung cho mọi tài liệu

```yaml
id: F-012                 # mã theo loại, xem types.toml
type: feature-spec
title: Bộ lọc đơn hàng
status: draft             # tập trạng thái theo loại
owner: ""                 # người chịu trách nhiệm, bắt buộc trước khi status rời draft
created: 2026-09-03
updated: 2026-09-03 14:05
source: CR-260910-loc-theo-trang-thai   # brief hoặc CR sinh ra hoặc sửa tài liệu này; rỗng với intake và cr
created_by: dk            # dấu vết; dk check báo file thiếu
dk_version: 0.1.0
```

Loại nào cần thêm trường thì khai trong `types.toml`. `dk new` điền toàn bộ; skill không sửa `created`, `created_by`, `dk_version`.

### 1.4 `types.toml` nhúng

Mỗi loại một bảng: `dir` (thư mục đích tương đối `docs/`, riêng `plan` và `report` trỏ `plans/`), `name` (mẫu tên file, ví dụ `{id}-{slug}.md`), `id` (`seq:F-{n:03}` đếm theo file có sẵn; `date:CR-{yymmdd}-{slug}`; `seq:ADR-{n:04}`; `none`), `subdir` (intake tạo thư mục ba file), `required` (trường frontmatter bắt buộc), `statuses`, `from` (ánh xạ trường chép khi `--from`, ví dụ brief sang feature-spec: `title`, `outcome`, `acceptance`). Danh sách loại theo phase 4 đến 9.

### 1.5 Changelog tài liệu

File `docs/CHANGELOG-DOCS.md`, mới nhất ở trên, nhóm theo ngày:

```
## 2026-09-03

- 14:10 | features/F-012-bo-loc-don-hang.md | +18 −4 | Thêm ngoại lệ B4 khi hết hàng | CR-260910-loc-theo-trang-thai
- 14:05 | design/mockups/F-012-B3.html | mới, 210 dòng | Mockup trạng thái rỗng | CR-260910-loc-theo-trang-thai
```

Ba lớp bảo đảm, cứng đến tiện: (1) pre-commit chạy `dk changelog pending`, chặn commit; (2) skill gọi `dk changelog add` với tóm tắt thật; (3) hook agent do `dk hook install` cài, fail-open khi `dk` lỗi. Không ghi changelog cho `docs/html/`, chỉ mục sinh ra, `plans/`.

### 1.6 Target

| | `claude` | `codex` |
|---|---|---|
| Skill, scope dự án | `.claude/skills/` | `.codex/skills/` |
| Skill, `--global` | `~/.claude/skills/` | `~/.codex/skills/` |
| Hook | `.claude/settings.json` hoặc `~/.claude/settings.json`, khóa `hooks` | `.codex/hooks.json` hoặc `~/.codex/hooks.json`; cần người dùng trust |
| Agent context file | `CLAUDE.md` | `AGENTS.md` |

```go
type Target interface {
    Name() string
    SkillDir(global bool) (string, error)
    InstallHooks(global bool, entries []HookEntry) error
    UninstallHooks(global bool) error
}
```

Nội dung skill không biết target. Dấu vết cài trong frontmatter SKILL.md theo trường `metadata` của Agent Skills spec (map chuỗi sang chuỗi, khóa nên có tiền tố riêng): `metadata.dk_installed_by: "dk"`, `metadata.dk_version`, `metadata.dk_hash` (sha256 mọi file trong skill, đã bỏ khối metadata).

### 1.7 Quyết định triển khai chốt ở phase 1

Các phase sau dùng đúng hợp đồng này, không đọc lại chữ mô tả trong `phase-01`:

- `types.toml`: `subdir` là một chuỗi mẫu (`{yymmdd}-{slug}`), không phải mảng; `from.<loại nguồn> = { <trường đích> = "<trường nguồn>" }`. Ví dụ `feature-spec.from.brief = { title = "title", purpose = "outcome", acceptance = "acceptance" }`. Mỗi loại có thêm `description` cho `template list`.
- `dk new <loại có subdir> --from <file>` đặt file mới vào đúng thư mục của file nguồn (idea, interview, brief cùng thư mục dù khác ngày). `source:` là `id` nguồn; nguồn không có id (intake) thì là `<thư mục>/<file>`.
- Loại `seq` luôn tăng số, nên trùng slug không bao giờ đụng file; `--force` chỉ có nghĩa với loại `none` và `date`.
- `--set k=v` ghi đè trường frontmatter sau khi render; `title` còn dùng để render thân.
- `dk changelog add` bump `updated:` trước khi đếm `git diff --numstat HEAD -- <file>` để số dòng khớp git sau khi lệnh chạy xong. Repo chưa có commit hoặc file chưa trong HEAD: `mới, N dòng`; không có git: `không git, N dòng`. Nguồn rỗng ghi `-`.
- `dk changelog pending` coi file đã ghi khi có dòng từ phút của commit HEAD trở đi (changelog chỉ lưu `HH:MM`); ghi rõ trong README. `pending` cần git, ngoài repo trả mã 1.
- `init` tạo thêm `plans/` (kèm `.gitkeep`) khi chưa có; không đụng khi đã có.
- Mã thoát: `unknown command` và sai số tham số đều là 2.
- Template: frontmatter là mục 1; thân Feature Spec đánh số 2 đến 11, CR 2 đến 7, khớp báo cáo nguồn. Delims giữ `{{ }}` (Mermaid mẫu không chứa `{{`). Trường `title`, `owner`, `source` trong frontmatter template là `""`; `dk new` đặt giá trị qua `yaml.Node` để tiêu đề có dấu hai chấm được quote đúng; thân vẫn dùng `{{.Title}}`.
- Sau review phase 1: `git status` chạy với `-z` và `core.quotePath=false` để tên file có dấu về nguyên dạng; file ẩn (`.gitkeep`) không ghi và không tính pending; `--summary`, `--source` không được chứa xuống dòng hoặc ` | ` (mã 2); `add` in và trả JSON đúng mục đã gộp; `init --force` giữ `dk.toml` cũ, chỉ bổ sung cây; `init` từ chối tạo dự án lồng trong dự án có sẵn; slug sai là mã 2; không gộp với mục có giờ mới hơn mục đang thêm; goldmark bỏ khỏi `go.mod` (tidy luôn xóa dependency chưa dùng), phase 2 thêm khi dùng. Báo cáo: `plans/reports/code-review-260903-phase-01-cli-core.md`.

### 1.8 Quyết định triển khai chốt ở phase 2

- `docs.Scan(root, dirs...)` là hàm quét chung: `Meta` giữ frontmatter, thân, số dòng, `Generated`; bỏ `docs/html/` và file ẩn; đường dẫn `Rel` tương đối gốc dự án. `docs.Resolve(metas, ref)` giải `source:` theo `id` hoặc `<thư mục>/<file>`.
- HTML: `docs/html/<đường dẫn>.html` giữ cây; liên kết `.md` trong `docs/` đổi sang `.html` tương đối, giữ `#neo`; liên kết `.md` ra ngoài `docs/` (ví dụ `plans/`) trỏ thẳng file gốc theo đường dẫn tương đối từ trang HTML. HTML thô trong Markdown bị escape (không `WithUnsafe`). Metadata là toàn bộ frontmatter theo thứ tự khóa, bỏ trường rỗng.
- Mermaid 11.17.2 bản IIFE (`globalThis.mermaid`), nhúng `<script>` chỉ khi trang có khối; giấy phép MIT đi kèm dưới dạng chú thích JS ở đầu script vì `html/template` bỏ chú thích HTML. `assets/html/THIRD_PARTY.md` và `MERMAID-LICENSE.txt` ghi nguồn.
- Bảng hành vi: bảng đầu tiên dưới tiêu đề chứa "hành vi"; không có thì bảng đầu tiên có cột đầu khớp `^B\d+$`. Chỉ bảng này nhận `id="step-Bn"`; dòng "Bước:" dưới mỗi sơ đồ liệt kê mã trong sơ đồ đó. `render.StepCodes` dùng chung cho `check`.
- `check` quét `docs/` và `plans/`; `line-threshold` chỉ áp cho `docs/`. Loại bắt buộc `source` lấy từ `required` trong `types.toml`; trạng thái "đã chốt" lấy từ trường mới `final` (brief: approved; cr: approved, in-progress, closed). Tài liệu loại chưa có trong `types.toml` (ví dụ `docs/test/`) vẫn được đếm khi xét "tài liệu chết" nhưng không chịu quy tắc theo loại. `check <file>` với file không nằm trong tập quét trả mã 1, không báo sạch. `dk.toml` có `[check] warn_lines = 500, max_lines = 800`; giá trị 0 hoặc thiếu lấy mặc định, không tắt được.
- `index`: 4 thư mục cố định `features`, `adr`, `cr`, `intake`; nhóm theo trạng thái theo thứ tự `statuses` của loại, trong nhóm mới nhất trước; không ghi `updated` để file không đổi mỗi lần chạy. `refs` và `index` HTML bỏ qua file `generated: true`.
- `render <file>` chỉ nhận `.md` (khác là mã 2); `render --all` render cả `CHANGELOG-DOCS.md` và README sinh ra (không liên kết từ index, chấp nhận).
- `status`: `changelog_pending` là `null` khi không có git; CR "đang mở" là mọi trạng thái ngoài `closed`, `rejected`.
- Nit chưa sửa từ review (`plans/reports/code-review-260903-phase-02-render-check.md`): so sánh tiền tố `..` trong `filepath.Rel` (có từ phase 1); `refs` có thể khớp id ngắn khi id dài hơn không tồn tại; `filepath.Rel` lỗi bị nuốt trong `index`; chưa có test `internal/config`; Windows (`C:/`, dấu `\`) để phase 11.

### 1.9 Quyết định triển khai chốt ở phase 3

- `target.Get(name, root)`: target giữ gốc dự án; scope dự án cần `dk.toml` (`requireProject`), `--global` thì không. `codex` trả lỗi "chưa hỗ trợ" (mã 2) tới phase 11.
- `settings.json`: đọc bằng object giữ thứ tự khóa, giá trị để thô; chỉ khối `hooks` được phân tích lại. Ghi bằng `json.Indent` thụt 2, không escape HTML, kết thúc bằng xuống dòng. `hooks: null` coi như rỗng. `hook uninstall` chỉ bỏ lệnh có tiền tố `dk hook run` trong từng mục matcher, giữ lệnh khác cùng mục; mục rỗng thì bỏ, sự kiện rỗng thì bỏ, khóa `hooks` rỗng thì bỏ; object rỗng thì xóa file (coi là do dk tạo).
- Skill: hash sha256 theo tên file đã sắp, SKILL.md băm sau khi bỏ `metadata` và ghi lại frontmatter dạng chuẩn (nguồn và bản cài đi qua cùng chuẩn hóa). Cài ghi vào thư mục tạm cạnh đích rồi `rename`, nên đứt giữa chừng không để lại thư mục thiếu dấu vết. Tên skill phải là một thành phần đường dẫn (không `/`, `..`). `uninstall` xóa `skills/` khi rỗng. `skill status` in scope dự án (khi đang trong dự án) và toàn máy. Trạng thái "đúng phiên bản" cần hash khớp nguồn nhúng và `dk_version` khớp binary; hash khớp `dk_hash` nhưng khác nguồn hoặc khác version là "cũ (vX)".
- `changelog.Record(root, docsDir, rel, summary, source, now, bump)` tách từ `dk changelog add` để hook gọi nội bộ; `changelog.Tracks`, `changelog.Load`, `changelog.FileName`, `changelog.MergeWindow` là hàm dùng chung. Hook `post-edit` truyền `bump=false`: không sửa `updated:` của file agent vừa ghi để tránh Claude Code báo "file modified since read" ở lần Edit kế; số dòng vì thế có thể lệch 1 so với lần `dk changelog add` sau đó, chấp nhận vì mục sẽ được thay.
- Gộp changelog: mục giữ chỗ (tóm tắt "chưa tóm tắt") được thay bởi mục kế tiếp cùng file trong 10 phút bất kể nguồn; hai mục có tóm tắt thật vẫn chỉ gộp khi cùng nguồn (giữ quyết định 1.7).
- `hook run`: đường dẫn tương đối ghép với `cwd` trong payload; tìm `dk.toml` đi lên từ thư mục chứa file; không tìm thấy hoặc file ngoài `docs/`, không `.md`, trong `html/`, file ẩn, `generated: true` hay chính changelog thì im lặng. Deny in JSON `hookSpecificOutput.permissionDecision: "deny"` không escape HTML, thoát 0. Mọi lỗi in `dk hook <sự kiện>: ...` ra stderr và thoát 0.
- Pre-commit: `assets/hooks/pre-commit.sh` với `__DK_CWD__` thay bằng đường dẫn dự án tương đối gốc git (git chạy hook tại gốc repo). Thư mục hooks lấy từ `git rev-parse --git-path hooks` (đúng với worktree và `core.hooksPath`). `init --force` khi đã có pre-commit in đoạn cần thêm, không ghi đè. Không có git thì bỏ qua, `--json` ghi `pre_commit.status: no-git`.
- Chưa giải quyết, chờ chủ dự án (từ `plans/reports/code-review-260903-phase-03-install-claude.md`): (H2) hook post-edit ghi dòng "chưa tóm tắt" làm `changelog pending` và pre-commit coi file đã ghi, lớp cứng mất tác dụng với file agent sửa; phương án là `pending` không tính mục "chưa tóm tắt", hoặc chấp nhận. (M6) `CHANGELOG-DOCS.md` không khóa khi hook chạy song song; chấp nhận vì pre-commit bắt lại. Nit còn lại: `hook install` luôn in "đã cài 2 hook" dù không đổi; `settings.json` ghi thẳng không qua temp+rename; `stdin` của `hook run` là biến mức package cho test.

### 1.10 Quyết định triển khai chốt ở phase 4

- Skill nhúng: `doc-intake`, `doc-cr`; `doc-smoke` đã xóa, test cơ chế cài dùng `doc-cr`. Test nội dung skill ở `internal/skill/content_test.go` (không đặt trong `assets/skills/` vì `//go:embed skills` nhúng nguyên thư mục): frontmatter với `name` trùng thư mục, mọi `.md` dưới 300 dòng, cấm `ak-`, `Claude Code`, `.claude/`, `.codex/`, `Edit`, `Write`, `MultiEdit`; hai mô tả không trùng 3 từ liên tiếp sau khi bỏ câu phủ định từ "Không dùng" trở đi (câu đó cố ý nhắc phạm vi skill kia). Mô tả có dấu hai chấm phải quote trong YAML.
- Mỗi skill có bước 0 `dk --version`, câu dừng bắt buộc in đậm đầu file và lặp tại bước; phần `new`, `update`, `html`. Trạng thái dùng khóa tiếng Anh của `types.toml` (`review` là chờ duyệt, `approved` là đã chốt); rules.md ghi bảng khóa và nghĩa.
- `types.toml`: `idea` bắt buộc `level`; `brief` bắt buộc `level`, `kind` (frontmatter thêm `approved_by`, `approved_at`); `cr` bắt buộc `requester`. Đây là thay đổi hợp đồng có chủ đích: tài liệu tạo trước phase 4 thiếu trường sẽ bị `dk check` báo `error`, phải bổ sung tay; chấp nhận vì chưa phát hành. `interview.from.idea` chép `level`, `interview.from.cr` chỉ chép `title`; template interview không đặt mặc định `level`. Khóa mới `beside_source = true` (chỉ `interview`): `--from` file ngoài thư mục loại nhưng trong `docs/` thì đặt vào thư mục cùng tên file nguồn (`cr/<CR-id>/interview.md`, và tương tự với nguồn bất kỳ như `features/F-001-x/interview.md`; không siết riêng cho CR), `source` là id nguồn; loại không có khóa giữ quyết định 1.7 (tạo subdir mới). Interview của CR tạo bằng `--from <file CR>`, không dùng `--set for=cr`.
- Lệnh mẫu trong SKILL.md đặt `owner` (và `requester` cho CR) ngay ở `dk new --set` vì các trường này bắt buộc; `doc-cr` chạy khối changelog, render, index, check sau **mọi** lần sửa CR, kể cả khi dừng chờ duyệt, để pre-commit không chặn. Changelog của intake dùng `--source <yymmdd>-<slug>` (tên thư mục), lệch dạng với `source:` frontmatter `<thư mục>/<file>`; chấp nhận vì `--source` là chữ tự do.
- `dk index <thư mục>` chỉ liệt kê loại có `dir` trùng thư mục (interview trong `cr/<CR>/` không thành dòng CR); loại lạ vẫn liệt kê. Chỉ mục `intake` một dòng một thư mục: đại diện là `brief.md` với trạng thái brief, không có thì `interview.md` rồi `idea.md` với nhóm `chưa có brief` xếp đầu.
- Skill giao việc cập nhật tài liệu đích cho "họ" tài liệu (họ Feature Spec, họ Design, họ ADR, họ Test, họ Release và User guide) kèm `--source <CR-id>`, không nêu tên skill khác; riêng `doc-intake` được chỉ sang `doc-cr` khi brief đã `approved`.
- Chạy thử tay chuỗi lệnh của hai skill trong dự án mẫu (`reports/phase-04-run.md` trong thư mục plan): `check` sạch, `refs` đúng, HTML sinh, changelog đủ. Phần AI dẫn phỏng vấn thật trong Claude Code chưa kiểm; gộp với bước 9 phase 3 vào đợt chạy thử toàn vòng phase 10.

### 1.11 Quyết định triển khai chốt ở phase 5

- Loại mới: `product-overview`, `architecture`, `glossary` (dir `overview`, id `none`, tên file cố định; `status` `draft | review | approved`, không `final` vì không tài liệu nào `source` trỏ về chúng), `adr` (dir `adr`, `seq:ADR-{n:04}`, `{id}-{slug}.md`, trường `supersedes`, `superseded_by`). Trạng thái ADR dùng khóa tiếng Anh theo 1.10: `proposed`, `accepted`, `superseded`, `deprecated`.
- `product-overview.from.brief` chỉ chép trường có thật trong frontmatter brief: `title`, `outcome`, `acceptance` (`source` do `--from` đặt); Ràng buộc và Ngoài phạm vi skill chép từ mục 2, 3 của thân brief. Không thêm `constraints`, `non_goals` vào frontmatter brief để không đổi hợp đồng phase 4. Việc "chỉ soạn từ brief `level: project`, `kind: product`, `status: approved`" là quy tắc skill, CLI không chặn.
- `types.toml` thêm khóa `warn_lines` theo loại (`product-overview` 150, `glossary` 300, `architecture` 400), `line-threshold` dùng thay `[check] warn_lines`; `max_lines` vẫn toàn cục. Glossary đếm dòng file, không đếm dòng bảng (bảng chiếm gần hết file, chênh vài dòng, chấp nhận).
- `adr-immutable`: so thân với bản HEAD qua `gitx.HeadFile` (`git show HEAD:<path>`) chỉ khi `status` **ở HEAD** thuộc `accepted | superseded | deprecated`; nhờ vậy lần chuyển `proposed` sang `accepted` vẫn được sửa thân. Chuẩn hóa CRLF về LF, cắt khoảng trắng cuối dòng và dòng trống cuối file. Không có git, chưa có HEAD hoặc file chưa trong HEAD thì bỏ qua; lỗi đọc HEAD là warning. Mã thoát vẫn 3 khi lỗi.
- `glossary-term`: chỉ quét chữ in đậm trong mục `## 2.`, `## 5.`, `## 8.` của Feature Spec (tránh **Given/When/Then** ở mục 9), mỗi thuật ngữ báo một lần tại dòng xuất hiện đầu, so không phân biệt hoa thường với cột đầu mọi dòng bảng của file loại `glossary`; chưa có Glossary thì bỏ qua rule.
- `dk index features` và `dk index adr` đổi sang bảng phẳng sắp theo mã (thay cho nhóm theo trạng thái ở 1.8): features `Mã | Tên | Trạng thái | Chủ sở hữu | Brief hoặc CR nguồn | Cập nhật`; adr `Mã | Tiêu đề | Trạng thái | Thay thế | Ngày`. `cr` và `intake` giữ cách nhóm cũ. Không có chỉ mục cho `overview/` (ba file cố định).
- Template không đặt liên kết sẵn đến file có thể chưa tồn tại (`architecture.md`, `glossary.md`, `../adr/README.md`, `../features/README.md`) vì `link-broken` là lỗi; chỉ gợi ý trong chú thích. Phát hiện khi chạy thử.
- Test cài đặt và `template list` không đếm cứng số skill hoặc loại theo thứ tự chữ cái nữa: tra theo tên (`rowOf`, `resOf`); số loại trong `template list` là 9, skill status 8 dòng.
- Sau review (`reports/code-review-260903-phase-05-skill-overview-adr.md`): `glossary-term` gộp mọi file loại `glossary` (tách Glossary theo miền vẫn đúng; `dk new glossary` lần hai vẫn trùng tên `glossary.md`, file thứ hai phải tạo bằng `--force` với tên khác chưa hỗ trợ, chờ nhu cầu thật), thứ tự finding theo thứ tự xuất hiện, bỏ khối mã, gộp hoa thường. `gitx.HeadFile` chỉ chạy `git show HEAD:<path>` (cộng `rev-parse` tìm gốc), mọi thất bại coi là chưa có trong HEAD. Chưa giải quyết, chờ chủ dự án (M3): `adr-immutable` chỉ so với HEAD nên thân ADR đã chốt bị sửa **và commit** thì thoát; phương án là thêm `dk check` vào pre-commit (hiện chỉ chạy `changelog pending`, quyết định 1.5) hoặc chấp nhận và ghi giới hạn. Nit chưa sửa: đổi tên file ADR đã chốt làm rule bỏ qua; `HeadFile` không phân biệt lỗi git thật với file chưa có.
- Chạy thử tay ở `reports/phase-05-run.md`: đủ 5 tiêu chí bằng CLI; phần AI thật (skill dừng khi brief chưa duyệt, Architecture từ mã, hỏi từng câu khi soạn ADR) gộp vào đợt chạy thử toàn vòng phase 10.

### 1.12 Quyết định triển khai chốt ở phase 6

- Template `feature-spec` một file 157 dòng với 5 biến thể qua `{{if}}` ở mức mục (chưa cần tách `feature-spec-<format>.md`): `spec` mặc định; `use-case` mục 4 là `sequenceDiagram`, mục 5 là Use Case Specification có bảng luồng chính theo mã bước (bảng này là "bảng hành vi" của `render` và `check` qua cột đầu `B\d+`); `story` mục 9 là khối ```gherkin; `crud` bỏ mục 4, mục 5 là bảng field và quyền, không mã bước nên `step-codes` bỏ qua; `state` thêm `stateDiagram-v2` trong mục 4 dưới tiêu đề cấp 3 "Vòng đời trạng thái" để giữ 11 tiêu đề số. Số mục cố định 2 đến 11, chữ tiêu đề theo biến thể.
- Frontmatter feature-spec thêm `has_ui: true`, `format: spec`; hai trường không bắt buộc (spec tạo trước phase 6 vẫn qua `check`, coi như `has_ui: true`, `format: spec`). `tmpl.Data` thêm `Format`, `HasUI`; `docs.New` đọc `--set format`, `--set has_ui` trước khi render; `format` ngoài `docs.SpecFormats` là lỗi. `--set k=true|false` ghi YAML bool (`frontmatter.SetBool`) cho mọi khóa, không quote.
- `from.brief` giữ `title`, `purpose`, `acceptance` theo 1.11; Ràng buộc và Ngoài phạm vi chép từ thân brief mục 2, 3 là việc của skill. Trạng thái giữ khóa tiếng Anh theo 1.10 (`draft`, `review`, `approved`, `implemented`, `deprecated`); statuses tiếng Việt trong chữ phase 6 không áp dụng.
- Mã bước nhận hậu tố chữ thường (`^B\d+[a-z]?$` trong `render.StepCodes`, dùng chung cho `check` và neo HTML `#step-B3a`) để chèn bước giữa không đánh số lại. Chỉ một chữ hậu tố; `B3aa` không được nhận.
- `spec-section-order`: chỉ xét số của tiêu đề `## N.` ngoài khối mã; tập mong đợi 2..11 bỏ 6 khi `has_ui: false`, bỏ 4 khi `format: crud`; thông điệp nêu thiếu, lạ, lặp hoặc sai thứ tự. Mục "Lịch sử nguồn" cuối file là tiêu đề cấp 3 nên không tính.
- `cr-approval-order`: spec có `source` trỏ CR (qua `docs.Resolve`) mà CR ngoài `final` của loại `cr` và `updated` spec sau `updated` CR là lỗi; `updated` không đọc được là warning. Định dạng `updated` giữ `YYYY-MM-DD HH:MM` (không thêm offset như phase 6 dự kiến, vì đổi định dạng chạm changelog và mọi tài liệu đã tạo); parse thử thêm dạng có offset và RFC3339, dạng không offset đọc theo giờ máy chạy `check`. Giới hạn chấp nhận: hai lần sửa cùng phút không phân biệt; máy khác múi giờ so lệch tối đa vài giờ.
- Skill `doc-feature-spec`: `update` chỉ nhận CR `approved` hoặc `in-progress` có dòng Feature Spec `Có`; đổi `source` sang mã CR và ghi nguồn cũ vào `### Lịch sử nguồn` sau mục 11; `--source` của changelog là thư mục intake khi tạo từ brief, mã CR khi sửa theo CR. Không có rule đo "ngoại lệ chiếm nửa bảng hành vi", kiểm bằng mắt.
- Test đếm skill đổi 4 sang 5 (`install_test`, `init_test` 10 dòng status). Chạy thử tay ở `reports/phase-06-run.md`: đủ 6 tiêu chí bằng CLI; phần AI thật gộp vào phase 10.
- Sau review (`reports/code-review-260903-phase-06-skill-feature-spec.md`): `step-codes` chỉ so khi spec có sơ đồ, nên "bỏ sơ đồ, bảng hành vi là đủ" (báo cáo mục 3, 10) là trạng thái hợp lệ; mã bước của spec không sơ đồ vẫn được ghi để userflow so tập con. Mục 6 của `crud` dùng cột "Màn hình" thay mã bước. `--set has_ui` chỉ nhận `true`/`false`. `findBehaviorTable` nhận tiêu đề chứa "hành vi" hoặc "use case". `sectionNumbers` xử lý fence ``` và ~~~ theo độ dài. Chưa sửa, chấp nhận: `--set k=true|false` ép YAML bool cho mọi khóa (không khóa nào trong `types.toml` cần chuỗi "true"); `updated` chưa nằm trong `required` của feature-spec nên spec viết tay thiếu `updated` chỉ nhận warning `cr-approval-order`; `format` chỉ validate cho feature-spec.

### 1.13 Quyết định triển khai chốt ở phase 7

- `types.toml` thêm khóa `kind = md | html | json` (mặc định `md`); template là `templates/<loại>.<kind>`, `tmpl.Raw` tìm theo glob `templates/<loại>.*` nên tên loại không được là tiền tố có dấu chấm của loại khác. Metadata chỉ hai định dạng ngoài YAML: HTML là YAML trong chú thích `<!-- dk:` ... `-->` đầu file; JSON là object dưới khóa `$dk` cấp cao nhất (đặt làm khóa đầu khi ghi, đọc được ở vị trí bất kỳ). `frontmatter.SplitFile`, `JoinFile` chọn theo đuôi; `docs.Read`, `docs.New`, `changelog.Record`, `changelog.Tracks` dùng chúng. `docs.Scan` quét thêm `.html`, `.json`; `Meta.IsMarkdown` để `render --all` chỉ render Markdown. `hook run` vẫn chỉ chặn `.md`.
- Tên file có placeholder `{feature}`, `{step}`, `{layer}` (và `dir = "design/{layer}s"`): giá trị chép từ nguồn `--from` theo bảng `from` rồi `--set` ghi đè; thiếu là lỗi nêu rõ cờ cần thêm. `step` phải khớp `^B\d+[a-z]?$`; `layer` thuộc `atom | molecule | organism | template`. `--from` chép: `userflow.from.feature-spec` lấy `feature = id`; `wireframe`, `mockup`, `ui-spec` `from.userflow` lấy `feature`; `mockup.from.wireframe` lấy cả `step`. `source` của các file này là `id` hoặc `<thư mục>/<file>` của nguồn `--from` (ví dụ `flows/F-001-flow.md`), skill đổi về brief hoặc CR khi biết.
- `dk tokens css [--in] [--out]`: mặc định `docs/design/tokens/tokens.json` sang `tokens.css`; parser JSON giữ thứ tự khóa (`internal/tokens/ordered.go`), `$type` kế thừa từ nhóm, alias `{a.b.c}` giải cả khi lồng trong chuỗi (bóng), alias vòng in chuỗi `a -> b -> a`, alias lạ là lỗi. Chuẩn hóa `dimension` (số thành `px`, object `{value, unit}`), `color` object lấy `hex`; loại khác chép nguyên chuỗi hoặc số, object loại khác là lỗi. Theme: `$extensions.dk.theme.<tên>` trên từng token thành khối `[data-theme="<tên>"]`. Dòng đầu `tokens.css` là `/* generated: dk tokens css */` (`changelog.GeneratedMarker`), `Tracks` bỏ file có dòng này; `.css` không thuộc tập quét.
- Mockup `external`: vẫn tạo `<F>-<B>.html` qua template nhánh `{{if .External}}` chỉ gồm liên kết và `<img src="<F>-<B>.png">`, không `tokens.css`, không `<style>`; rule `mockup-tokens` bỏ qua file có `external`. Trạng thái màn hình ghi ở trường `state` (`normal`, `empty`, `error`, `loading`); cùng mã bước nhiều trạng thái đặt `step` kèm hậu tố chữ thường chưa dùng trong spec (quy ước ở rules.md của `doc-design-flow`), là tên file mockup, không phải bước mới.
- `check`: `mockup-tokens` quét `<style>…</style>` và `style="…"` bằng regexp `#[0-9a-fA-F]{3,8}\b` và `\b\d+(\.\d+)?px\b`, báo từng giá trị kèm dòng; id selector toàn chữ hex (`#abc`) bị báo giả, chấp nhận. `userflow-steps` thay nhánh userflow cũ trong `step-codes` (so theo `source`): so tập con với spec có id bằng `feature`, `feature` không trỏ spec nào là lỗi; spec không mã bước thì userflow cũng không được có. `link-broken`, `refs` không xét `<a href>` trong mockup.
- Skill: `doc-design-system` bước 2 đến 6 theo thứ tự lớp 0 đến 4 rồi kiểm tra thật bằng một mockup `feature=F-000 step=B0` từ template, xóa sau khi qua; `doc-design-flow` bước 1 đến 7 (flow, wireframe người duyệt, mockup mỗi trạng thái, prototype bằng thẻ `a`, UI spec chỉ khi có công cụ sinh). `test` đếm skill 7, `skill status` 14 dòng, `template list` 17 loại. README thêm mục "Design tokens và mockup", 13 quy tắc `check`.
- Sau review (`plans/reports/code-review-260903-phase-07-skill-design.md`): file `.html` trong `docs/` thiếu khối `<!-- dk: -->` là lỗi `mockup-tokens` và vẫn bị lint (không lách bằng cách bỏ metadata); `style='...'` nháy đơn cũng lint; px trong prelude `@media`, `@container` không báo; giá trị `feature`, `step`, `layer` chỉ nhận `[A-Za-z0-9][A-Za-z0-9_.-]*` (không `/`, `..`); `line-threshold` chỉ áp cho Markdown. Mockup mẫu của `doc-design-system` giữ lại là `F-000-B0.html` (`F-000` dành riêng cho mẫu). Chấp nhận: `hook run` chỉ chặn `.md`; lint chỉ hex và px; alias đến token có theme không kế thừa theme; `$dk` trùng khóa trong JSON viết tay không xử lý.
- Chạy thử tay ở `reports/phase-07-run.md`: đủ 6 tiêu chí bằng CLI; phần AI thật (điền component 5 mục, sinh mockup chỉ bằng organism) gộp vào đợt chạy thử toàn vòng phase 10.

### 1.14 Quyết định triển khai chốt ở phase 8

- `types.toml` thêm 11 loại (28 loại, 9 skill, `skill status` 18 dòng). `dir` có hai tiền tố gốc mới: `{plans_dir}` (gốc là `plans_dir` của `dk.toml`, loại `plan`) và `{in}` (gốc là thư mục `--in`, loại `plan-phase` và `report`; thiếu `--in` là lỗi mã 1). `report` không thuộc plan nào dùng `--in plans` để vào `plans/reports/`. `subdir` và `name` nhận thêm `{hhmm}`. Loại ngoài `docs/` không áp `id_prefix` (`phase-01` đếm cục bộ trong thư mục plan). Không thêm `plan_name_pattern`: repo có công cụ plan riêng thì skill không tạo plan, chỉ viết report và decision log (quy tắc skill, CLI không chặn).
- `kind = "feature"` cho Gherkin: metadata YAML trong khối chú thích đầu file mở bằng `# dk:`, mỗi dòng `# <yaml>`, kết thúc ở dòng đầu không bắt đầu bằng `#` (hoặc hết file); `JoinFeatureComment` chèn dòng trống khi thân bắt đầu bằng `#`. `docs.Scan` quét thêm `.feature`; không render (`IsMarkdown` false), `changelog add` và `check` bình thường; `hook run` vẫn chỉ chặn `.md`. Vì dòng đầu file là `# dk:` nên không dùng `# language:`; từ khóa Gherkin tiếng Anh, nội dung bước tiếng Việt.
- Test case là hai loại thay vì một loại hai đuôi: `test-case` (`.feature`, template `test-case.feature`) và `test-case-table` (`.md`); Testing strategy `format: gherkin | table` quyết định skill gọi loại nào. Cơ chế `templates/<loại>.<kind>` và `doctype.validate` giữ nguyên.
- `docs.ExtractSpec` (`from_spec.go`) chỉ chạy khi `--from` là `feature-spec` và loại đích có `from.feature-spec`: mục 9 thành `Scenario` (bullet `- ACn. **Given** a **When** b **Then** c`, in đậm tùy chọn, dấu phẩy cuối bước bị bỏ; khối ```gherkin với `Scenario: ACn ...` và `And`/`But` nối vào bước trước), mục 3 thành `Background` (bỏ nhãn trước dấu hai chấm), mục 6 thành checklist (dòng bảng có cột đầu là mã bước, cột hai là mockup). AC có mã nhưng thiếu một trong ba phần giữ ở `Raw`: template ghi `# chưa tách được: <chữ gốc>` và ba bước trống, không bỏ qua im lặng. Không có `--from` thì một `Scenario: AC1` trống và cần `--set feature`. Ngoại lệ mục 7 do skill thêm tay (tag `@En`), CLI không sinh.
- `--append` chỉ cho `decision-log`: file chưa có thì tạo từ template rồi nối; đã có thì nối dòng `- <yyyy-mm-dd> | <chữ>` vào cuối, bump `updated`, không render lại template (`--set` bị bỏ qua khi nối). Slug bỏ được khi có `--append` (`dk new decision-log --append "..."`); `--append` không được chứa xuống dòng (mã 2). `dk.toml` thêm `[test] bdd_cmd` (chỉ skill đọc; rỗng là "chưa kiểm chạy được").
- `check`: `report-evidence` (warning) khi thân report, sau khi bỏ chú thích HTML, không có `commit` kèm hex 7 đến 40 ký tự trong 20 ký tự sau, không có liên kết Markdown `](...)`, không có khối ``` có nội dung; khối ``` trống của template không tính. `spec-has-test` tính test case là file trong `docs/test/` có `source` và loại `test-case*` hoặc loại chưa có trong `types.toml`; `ui-test-checklist`, `test-report`, `testing-strategy` không thay được test case. Rule mới đứng cuối danh sách (14 rule). README thêm mục "Plan, report và test".
- Sau review (`plans/reports/code-review-260904-phase-08-skill-plan-test.md`): mọi dòng mục 9 nhắc `ACn` (danh sách đánh số, bảng, `Scenario Outline`) chưa tách được đều thành Scenario `# chưa tách được`, không bỏ im lặng; `dk new` in số Scenario và số dòng chưa tách (`--json`: `scenarios`, `unparsed`). `--in` phải là thư mục có sẵn trong `plans_dir`; loại không dùng `{in}` mà có `--in` là lỗi. `.feature` trong `docs/` không đọc được khối `# dk:` (kể cả khối rỗng) là lỗi `frontmatter-required`. `report-evidence` không tính liên kết `.md`. Chấp nhận: `dk status` không quét `plans/`; thiếu `--in` mã 1; `--append` không khóa file; `report` luôn cần `--in`; `report-evidence` cảnh báo cả `draft`.
- Chạy thử tay ở `reports/phase-08-run.md`: đủ 6 tiêu chí bằng CLI; máy không có công cụ BDD nên nhánh dry-run "có công cụ" chưa kiểm, gộp vào phase 10 cùng phần AI thật (đề xuất strategy từ stack, thêm Scenario ngoại lệ, chạy checklist UI). Câu hỏi mở "công cụ BDD" vẫn mở: định dạng đã chốt là Gherkin chuẩn không phụ thuộc công cụ.

### 1.15 Quyết định triển khai chốt ở phase 9

- `types.toml` thêm 13 loại (41 loại, 11 skill, `skill status` 22 dòng). Không có loại `user-guide-index`: mục lục User guide là `docs/release/guide/README.md` sinh bởi `dk index user-guide` (`generated: true`, nhóm theo trường `task`, trang thiếu `task` xếp cuối "Chưa phân nhóm"), cùng cách với `features`, `adr`. `index` ánh xạ tên chỉ mục sang thư mục (`indexDirOf`), `index all` ghi 5 file.
- `release-brief`: `dir release/briefs`, `name {feature}.md`, `from.feature-spec = { title, feature = id }`; `source` là mã spec. `docs.ExtractRelease` (trong `from_spec.go`, dùng lại `sections`, `bullets`, `uiSteps`) chép mục 2 (mọi dòng không rỗng, bỏ chú thích HTML), mục 3 (bullet, bỏ nhãn), cột hành động của bảng hành vi mục 5 (dòng có cột đầu là mã bước), cột mockup mục 6 với liên kết tương đối đổi gốc từ `features/` sang `release/briefs/` (`rebaseLinks`; liên kết tuyệt đối hoặc có scheme giữ nguyên), mục 11 (bullet). Chữ chép sang là nháp để skill viết lại; `tmpl.Data` thêm `Purpose`, `Actors`, `Actions`, `Limits`, dùng lại `Steps`; template có hàm `inc` để đánh số. Trạng thái brief `draft | review | ready`; `kind: feature | fix` bắt buộc; `released_in` chỉ do `--collect` ghi, trạng thái không đổi (brief `ready` có `released_in` là đã phát hành).
- `release-notes`: `name {version}.md`; `dk new release-notes --collect <phiên bản>` không cần slug (`docs.VersionSlug`: `v1.0.0` thành `v1-0-0`), `version` đi qua kiểm tra giá trị placeholder như `feature`, `step`, `layer`; `title` mặc định là phiên bản. `collectBriefs` quét `docs/release/briefs` (gốc là `docs/`), lấy `type: release-brief`, `status: ready`, `released_in` rỗng; rỗng là lỗi mã 1; `--collect` cho loại khác là lỗi mã 1. Sau khi ghi notes, `markReleased` đặt `released_in` và bump `updated` của từng brief; CLI gọi `changelog.Record(bump=false)` cho mỗi brief với tóm tắt `Phát hành trong <phiên bản>` và nguồn `<phiên bản>`; `Result` JSON thêm `released` (đường dẫn tương đối gốc). `--collect` không lọc theo `feature` hay thư mục con.
- `--append` mở rộng cho `faq` (cùng dạng `- <ngày> | <chữ>` với decision-log, không thêm định dạng riêng). `postmortem` bắt buộc `incident_at` (`yyyy-mm-dd hh:mm` hoặc `yyyy-mm-dd`, giờ máy); `dk new` ghi `written_within_48h` bool: `true` khi `0 <= now - incident_at <= 48h`, sai định dạng hoặc thiếu là `false`; skill báo người khi `false`. `backup-dr` có `last_drill`; `dk status` thêm `dr_overdue` (mảng, luôn có trong JSON) với file `last_drill` rỗng hoặc trước `now - 6 tháng`, in dòng "DR chưa diễn tập quá 6 tháng" ở cuối.
- `dk.toml` thêm `[release] jargon` (`config.Release`, mặc định `config.DefaultJargon` = API, endpoint, database, migration, backend, frontend, JSON; rỗng lấy mặc định, `init` ghi danh sách vào file để người thấy). `check.Context` thêm `Jargon`. `no-jargon` (warning) chỉ áp `release-brief`, `user-guide`: so cả từ không phân biệt hoa thường, bỏ khối mã và chú thích HTML (kể cả nhiều dòng), mỗi từ báo một lần tại dòng đầu; skill không tự sửa `dk.toml`. `env-no-secret` (lỗi) chỉ áp `environment`: dòng `^\s*(export\s+)?KEY\s*=\s*value` với value bỏ nháy không rỗng và không phải `<...>`; bắt cả trong khối mã; bảng không bắt. Hai rule đứng cuối danh sách (16 rule). README ghi rõ rule là lớp phụ, không thay công cụ quét secret.
- Mô tả kích hoạt hai skill đổi câu phủ định cuối thành dạng "Không dùng cho ..." (quy ước 1.10) vì test trigram chỉ bỏ câu bắt đầu "Không dùng"; nội dung phạm vi giữ nguyên chữ phase 9. Tầng 9 (`charter`, `risk-register`, `meeting-notes` trong `docs/governance/`, meeting theo `{yymmdd}-{slug}`) chỉ có template; README và `types.toml` ghi "không có skill vì tài liệu quản trị do người viết". `init` đã tạo `release/`, `ops/`, `governance/` từ phase 1, không đổi.
- Sau review (`plans/reports/code-review-260904-phase-09-skill-release-ops.md`): `markReleased` dựng nội dung mọi brief xong mới ghi, lỗi I/O giữa chừng nêu danh sách file đã sửa (H1); `--collect` gom thêm brief đã có `released_in` bằng đúng phiên bản đang tạo nên `--force` sinh lại notes cùng phiên bản không mất dòng cũ, chỉ ghi `released_in` cho brief mới (H2); `env-no-secret` bỏ chú thích HTML nhưng vẫn quét khối mã (mẫu cấu hình nằm đó), placeholder phải khớp trọn `^<[^<>]*>$`; `dk status` tách thông điệp `last_drill` sai định dạng; `ExtractRelease` bỏ ô mockup không có liên kết (chữ mẫu "chưa có, xem họ Design") và `rebaseLinks` giữ nguyên mọi liên kết có scheme (`mailto:`). Chấp nhận, ghi nhận: `no-jargon` không tắt được bằng danh sách rỗng (như `warn_lines`, README ghi); sai loại cho `--collect`/`--append` là mã 1 như phase 8; `index all` tạo `release/guide/README.md` dù chưa có trang; `--collect ""` báo lỗi slug rỗng không nhắc cờ; `filepath.Rel` nuốt lỗi khi in `released`; `glossary-term` chưa dùng `countFMLines`.
- Chạy thử tay ở `reports/phase-09-run.md`: đủ 6 tiêu chí bằng CLI; phần người không kỹ thuật đọc Release brief và Runbook chạy trên sự cố thật gộp vào đợt chạy thử toàn vòng phase 10.

### 1.16 Quyết định triển khai chốt ở phase 10

- Lệnh mới: `dk doctor [--target] [--global]` in bảng `mục | trạng thái | cách sửa` (dk.toml, docs/, `dk` trên PATH, git, pre-commit có và gọi `changelog pending` và có quyền chạy, từng skill không đúng phiên bản kèm lệnh cài, hook đủ 2), mã 3 khi có mục chưa đạt; target lạ hay lỗi đọc skill vẫn in bảng đã gom. `dk self-check`: phiên bản, số template, số skill, target, sha256 mọi file nhúng (so hai binary), mã 1 khi template hoặc skill nhúng hỏng. `dk init --agent-context` in `assets/agent-context.md` (45 dòng, bảng 11 skill viết tay, test bắt thiếu tên skill) ra stdout, không ghi file. `Target` thêm `InstalledHooks(global)`.
- Phiên bản: chuỗi không có tiền tố `v` ở mọi đường (Makefile `sed s/^v//`, goreleaser `{{ .Version }}`, `root.go` lấy `debug.ReadBuildInfo().Main.Version` bỏ `v` khi build bằng `go install pkg@vX`); lệch chuỗi làm `skill status` báo "cũ" sai.
- Release: `.goreleaser.yaml` v2 (linux, darwin × amd64, arm64; `tar.gz` và `zip`; `checksums.txt`), `release.yml` trên tag `v*` chạy test và lint trước; `ci.yml` vet, test, `make lint-skills` trên ubuntu và macos, `make build-all` cross-build 4 binary. `goreleaser check` chỉ báo thiếu remote. **Chưa tag `v0.1.0`**: cần remote thật và chốt module path (đang `github.com/vner/dk`, đổi ở `go.mod`, import, Makefile, `.goreleaser.yaml`, README).
- `make lint-skills` (grep từ cấm với `-w` cho tên tool để chạy trên BSD grep, đếm dòng bằng awk) chạy cùng `content_test.go` trong CI. Docs của repo CLI viết tay (`docs/kien-truc.md`, `docs/lenh.md`, `docs/skill.md`), README rút gọn; lý do không dùng `dk` cho docs của `dk` ghi ở `kien-truc.md`.
- Toàn vòng (`reports/phase-10-full-run.md`, agent đóng vai trên dự án mẫu `/tmp/dk-fullrun/quan-ly-don`, 12 bước, 195 lệnh `dk`, 16 commit qua pre-commit, 25 câu hỏi trong đó 1 thừa, 4 can thiệp): không lệnh `dk new` nào thất bại, không tạo file tay; mô tả skill không có cặp nào cùng khớp mà không loại trừ nhau. Sửa theo kết quả: (skill) `doc-design-flow` bỏ qua `userflow-steps` cho đến khi có spec và chép style organism từ `F-000-B0.html`; `doc-design-system` nói cách có Design brief (họ Intake) và bỏ component bị brief loại; `doc-cr` lấy spec từ `dk index features` khi chỉ có một, "Cần làm gì" chỉ liên kết file đã có; `doc-feature-spec` spec `implemented` sửa theo CR về `approved`; `doc-plan-report` thêm `--set owner`, `--set title` có dấu cho phase và report; mô tả `doc-design-system`, `doc-design-flow`, `doc-test` thêm ý cập nhật theo CR bằng chữ khác nhau (test trigram). (CLI) gộp changelog trong 10 phút nối hai tóm tắt thật khác nhau bằng `; ` thay vì ghi đè (đổi quyết định 1.7, test cập nhật); `dk check <thư mục>` là mã 2 với thông báo đúng; `ui-test-checklist` không in số Scenario; `userflow.from.brief` chép `title`; `dk refs` bỏ `CHANGELOG-DOCS.md`. (Template) `design-tokens` thêm `border-width.thin`, mockup dùng `var(--border-width-thin)` thay `--space-1`; `design-component` mẫu lệnh có `--set source`. Chưa sửa, chấp nhận: `test-case-table` cột Tên trống và giữ dấu chấm cuối Then; `plan-phase` không kế thừa owner từ plan.
- Chưa kiểm được ở phase này, chuyển sang khi có remote hoặc Claude Code thật: kích hoạt skill thật trong Claude Code (hook agent chỉ kiểm bằng unit test và `hook run` qua stdin); CI xanh thật; `go install @v0.1.0`; người không kỹ thuật đọc Release brief; dry-run BDD có công cụ.
- Sau review (`plans/reports/code-review-260904-phase-10-stabilize-release.md`): H1 chuỗi phiên bản, M1 doctor scope `--global`, M2 doctor không mất bảng khi target lỗi, M3 `permissions: contents: read` trong `ci.yml` đã sửa. Chấp nhận: dự án không git luôn mã 3 ở doctor (git là điều kiện để pre-commit chạy); `self-check` ghi cứng target; action chưa ghim SHA; repo chưa có LICENSE (cần chủ dự án chọn).

### 1.17 Quyết định triển khai chốt ở phase 11

- Định dạng hook Codex (Codex CLI 0.153.2, tài liệu hooks và `apply_patch.rs`): `hooks.json` cùng cấu trúc khối `hooks` với `settings.json` của Claude, **`command` là chuỗi** (phase 11 dự kiến mảng, bỏ); sửa file phát hook với `tool_name: "apply_patch"`, patch nằm trong `tool_input.command`; deny cùng JSON `hookSpecificOutput.permissionDecision`; hook cấp dự án cần trust qua `/hooks` trong Codex. Đóng mục UNVERIFIED số 8 của validate phiên 1 bằng tài liệu và mã nguồn, chưa bằng payload ghi từ phiên thật (máy chưa `codex login`).
- `internal/target/claude.go` tách `installHooksFile`, `uninstallHooksFile`, `installedHooksFile` theo đường dẫn; `codex.go` chỉ đổi đường dẫn (`.codex/hooks.json`, `$CODEX_HOME` rồi `~/.codex`) và matcher `apply_patch` cho cả hai sự kiện (`hook run` tự lọc theo đường dẫn). Không ghi `config.toml`. `target.Names = [claude, codex]`; `self-check` in từ đó.
- `hook run`: `payload.files()` lấy `tool_input.file_path`, dự phòng `tool_input.path`, với `apply_patch` (hoặc `command` chứa `*** Begin Patch`) đọc `*** Add File:` và `*** Update File:` (bỏ `Delete`); mỗi file qua `handle`, deny in đúng một lần dù patch tạo nhiều tài liệu. Không có file nào thì im lặng kể cả sự kiện lạ (như trước với payload thiếu `file_path`).
- `--target` nhận danh sách phẩy (`targetFlags.names()`: cắt khoảng trắng, bỏ rỗng, bỏ trùng; tên lạ hoặc rỗng là mã 2 trước khi đụng target nào) cho `skill install|uninstall|status`, `hook install|uninstall`, `doctor`. Kết quả gộp theo thứ tự target; lỗi giữa chừng vẫn in phần đã làm. `hook install` in `đã cài <2×số target> hook` và với codex thêm dòng `target.CodexTrustNote` (JSON: khóa `note`).
- `doctor`: không nêu `--target` mà dự án có `.codex/` thì kiểm thêm codex. Nhãn dòng đổi (hợp đồng có chủ đích): `skill (claude, dự án)`, `hook (codex, toàn máy)`, `skill doc-x (claude)`; cách sửa chỉ thêm `--target` khi khác target mặc định; hook codex thiếu thì cách sửa kèm nhắc trust. Windows không kiểm bit chạy của pre-commit.
- Windows: Makefile và goreleaser 6 nền tảng (`.exe`, archive `zip` cho Windows qua `format_overrides`), CI `windows-latest` chạy vet và test (bỏ `lint-skills` vì thiếu make/awk GNU), `gitx.NumStat` truyền pathspec dạng slash, test đặt `USERPROFILE` cùng `HOME`, test `sh` và bit chạy bỏ qua trên Windows, `TestDoctor` tạo `dk.bat` trên Windows. Hook Claude Code trên Windows chạy qua Git Bash hoặc PowerShell, lệnh `dk hook run` không phụ thuộc shell nên không thêm trường `shell`. Chỉ kiểm bằng `GOOS=windows go vet` và biên dịch test, chưa chạy trên máy Windows.
- `assets/agent-context.md` không đổi: khối đã trung lập target, dán vào `CLAUDE.md` hay `AGENTS.md` đều được.
- Sau review (`plans/reports/code-review-260904-phase-11-target-codex.md`): `patchFiles` đọc thêm `*** Move to:` (C1, đổi tên vào `docs/` bị chặn); `doctor` chỉ tự kiểm codex khi có `.codex/skills/` hoặc `.codex/hooks.json` (H2, `.codex/` rỗng sau gỡ không kích hoạt); `gitx.run` cắt `\r`; README ghi hook không bắt ghi file qua shell ở cả hai target. Giữ: `doctor --target lạ` mã 3 kèm bảng (M2 phase 10); số hook in cố định (nit 1.9). Chấp nhận: `doctor --global` vẫn dò `.codex/` của dự án; target chưa cài in 11 dòng thiếu.
- Kết quả kiểm ở `reports/phase-11-codex.md`. Chưa kiểm: vòng Tiếp nhận trong Codex đã đăng nhập (deny khi `apply_patch` tạo tài liệu mới, changelog sau sửa); CI Windows và máy Windows thật; tag `v0.2.0` (chặn cùng `v0.1.0`).

### 1.18 Quyết định sau phase 11: khối ngữ cảnh agent ghi thẳng vào file

- `dk init --agent-context` không in ra stdout nữa; ghi `assets/agent-context.md` vào cả `CLAUDE.md` và `AGENTS.md` tại thư mục hiện tại (không cần `dk.toml`). Khối nằm giữa `<!-- dk:agent-context start version=<dk> hash=<sha256 16 hex> -->` và `<!-- dk:agent-context end -->` (`internal/agentctx`): file chưa có thì tạo; có khối thì thay tại chỗ; chưa có khối thì nối cuối, phần ngoài khối giữ nguyên; có mốc mở thiếu mốc đóng là lỗi mã 1. In `created | updated | unchanged` từng file, `--json` là mảng.
- `dk doctor` thêm hai dòng `agent context (CLAUDE.md)`, `agent context (AGENTS.md)`: `missing-file`, `missing-block`, `outdated` (khác phiên bản, hash, hoặc thân khối bị sửa tay), khối thiếu mốc đóng; cách sửa `dk init --agent-context`. Test doctor "đủ cài đặt" thêm bước này.
- Sau review: ghi ở gốc dự án khi có `dk.toml` (từ thư mục con vẫn đúng chỗ doctor kiểm), chưa init thì thư mục hiện tại; kiểm cả hai file trước khi ghi để file hỏng không để file kia đã đổi; nhận CRLF (regexp `\r?\n`, so khối sau khi đổi CRLF về LF); nhiều hơn một khối là `broken`; hash tính trên thân đã bỏ xuống dòng cuối đúng như ghi vào file; trạng thái doctor in tiếng Việt (`thiếu file`, `có file, chưa có khối`, `cũ hoặc bị sửa tay`). Chấp nhận có chủ đích: dự án chưa chạy `--agent-context` thì `doctor` mã 3, như thiếu skill hay hook.
- Mục lục HTML: `render` gom tiêu đề cấp 2 và 3 ở mức đầu tài liệu (`render.toc`, id do goldmark gán, bỏ ký tự ngoài ASCII như trước), đưa vào `<nav class="toc">` dưới bảng metadata trong `aside` dính cột phải; dưới 3 tiêu đề thì không hiện; `aside` có khi có metadata hoặc mục lục; trang chỉ mục không có. Golden `plain` đổi vì CSS nhúng thêm `.toc`.
- Nội dung khối viết lại bằng tiếng Anh: rút gọn 6 nguyên tắc của `docs/rule.md` (nghĩ trước khi code, tối giản, sửa đúng chỗ, mục tiêu kiểm được, hạn chế subagent, trả lời đúng câu hỏi) rồi phần `dk` như cũ đã dịch. Ngưỡng test đổi từ dưới 60 lên dưới 120 dòng; vẫn bắt đủ 11 tên skill.

## 2. Phase

| # | Phase | Sản phẩm | Phụ thuộc | Trạng thái |
|---|---|---|---|---|
| 1 | [CLI lõi](./phase-01-cli-core.md) | module Go, embed, `init` (chưa cài skill), `template`, `new`, `changelog`, `dk.toml`, template 1a, 1b, feature-spec | | Done |
| 2 | [Render và kiểm tra](./phase-02-render-check.md) | `render`, `index`, `check`, `refs`, `status` | 1 | Done |
| 3 | [Cài đặt cho Claude](./phase-03-install-claude.md) | `skill *`, `hook *`, target `claude`, `hook run`, pre-commit trong `init` | 1 | Done (bước 9 chưa kiểm thật) |
| 4 | [Skill tiếp nhận](./phase-04-skill-intake-cr.md) | `doc-intake`, `doc-cr` | 1, 2, 3 | Done (bước 7 chưa chạy thật) |
| 5 | [Skill nền sản phẩm](./phase-05-skill-overview-adr.md) | `doc-overview`, `doc-adr` | 4 | Done (phần AI thật chưa kiểm) |
| 6 | [Skill Feature Spec](./phase-06-skill-feature-spec.md) | `doc-feature-spec` | 4, 5 | Done (phần AI thật chưa kiểm) |
| 7 | [Skill Design](./phase-07-skill-design.md) | `doc-design-system`, `doc-design-flow` | 1, 4 | Done (phần AI thật chưa kiểm) |
| 8 | [Skill tiến độ và test](./phase-08-skill-plan-test.md) | `doc-plan-report`, `doc-test` | 6 | Done (phần AI thật, dry-run BDD chưa kiểm) |
| 9 | [Skill phát hành và vận hành](./phase-09-skill-release-ops.md) | `doc-release`, `doc-ops`, template Tầng 9 | 6 | Done (người đọc thật, runbook trên sự cố thật chưa kiểm) |
| 10 | [Ổn định và phát hành](./phase-10-stabilize-release.md) | chạy thử toàn vòng, Agent context file, goreleaser, `v0.1.0` Linux và macOS | 1 đến 9 | Done (tag `v0.1.0` chờ remote) |
| 11 | [Target Codex và Windows](./phase-11-target-codex.md) | `codex.go`, hook `hooks.json`, kiểm thật, build Windows, `v0.2.0` | 10 | Done (kiểm thật trong Codex, Windows máy thật, tag `v0.2.0` chờ) |

Thứ tự bắt buộc: 1, 2, 3 trước mọi skill; 4 trước 5 và 6; 6 trước 8 và 9. Phase 7 song song với 5 và 6 nhưng xong trước khi `doc-feature-spec` điền mục Giao diện. Không có giao diện người dùng thì bỏ phase 7 và phần test giao diện của phase 8. Phase 11 chỉ bắt đầu khi 10 xong.

## 3. Quy trình chung cho mỗi skill (phase 4 đến 9)

1. Thêm loại vào `assets/types.toml`, viết template vào `assets/templates/`, test đơn vị template qua (đủ trường bắt buộc, render không lỗi).
2. Viết `references/rules.md` từ báo cáo: chỉ quy tắc, liên kết về mục của báo cáo cho phần giải thích, dưới 300 dòng.
3. Viết `SKILL.md` dưới 300 dòng: phạm vi (họ này, không họ khác), bước 0 kiểm tra `dk --version`, quy trình `new`, `update`, `html`, mỗi bước ghi rõ lệnh `dk` và khi nào dừng hỏi người. Không nhắc target.
4. `go build`, `dk skill install <tên>`, chạy thử với một tài liệu thật, `dk check` qua, HTML đọc được.
5. Mô tả kích hoạt không trùng từ khóa với skill khác trong bộ; grep `ak-` và `Claude Code` trong `assets/skills/` bằng không.

## 4. Rủi ro

- **CLI phình thành framework**: quy tắc "chỉ việc xác định, không suy luận". Lint code UI là của dự án đích.
- **Template không tùy biến theo dự án**: chấp nhận có chủ đích; nhu cầu đổi đi qua phát hành. Nếu sau này cần, thêm ghi đè theo dự án không đổi kiến trúc.
- **Skill sửa tay bị ghi đè**: hash trong frontmatter, `install` từ chối khi lệch nếu thiếu `--force`.
- **`dk` không trên PATH lúc hook chạy**: hook fail-open kèm cảnh báo; pre-commit bắt.
- **Hook Codex không bắt ghi file hoặc chưa trust**: không coi hook là bảo đảm; Codex 0.153.2 bắt `apply_patch` theo tài liệu, chưa kiểm trong phiên đã đăng nhập (mục 1.17).
- **HTML nặng vì nhúng mermaid.min.js** (khoảng 2 đến 3 MB mỗi file): chấp nhận theo quyết định "tự chứa"; chỉ nhúng khi file có khối Mermaid.
- **Hook chặn nhầm**: chỉ chặn tạo file mới trong `docs/**/*.md`, không chặn sửa; `dk new --force` cho trường hợp đặc biệt.
- **Changelog rác**: gộp trong 10 phút cùng file cùng nguồn.

## 5. Rollback

Thứ tự: `dk skill uninstall` và `dk hook uninstall` cho từng target và scope, gỡ pre-commit, rồi mới gỡ binary. `docs/` đọc được không cần CLI; `docs/html/` là bản sinh, xóa được.

## Câu hỏi chưa giải quyết

- Module path Go thật (remote GitHub nào)? Repo là chính thư mục này, `git init` ở phase 1; vẫn `github.com/vner/dk` sau phase 10 vì chưa có remote. **Chặn tag `v0.1.0`**: cần chủ dự án tạo remote, đổi module path (5 chỗ, mục 1.16), commit đầu, rồi `git tag v0.1.0 && git push --tags`.
- LICENSE cho repo và archive release: chưa chọn.
- Tên lệnh `dk`: chưa trùng trên máy này; xác nhận trên máy khác trong team.
- Công cụ BDD cho Gherkin? Phase 8 chốt định dạng Gherkin chuẩn (từ khóa tiếng Anh, tag `@F-xxx @ACn`), công cụ ghi ở `dk.toml` `[test] bdd_cmd` khi chốt; nhánh dry-run kiểm ở phase 10.
- Changelog một file toàn dự án hay theo thư mục con? Đề xuất một file, cắt theo tháng khi vượt 1000 dòng.

## Validation Log

### Session 1 — 2026-09-03
**Trigger:** người dùng chọn `/ak:plan validate` sau khi plan 11 phase được viết.
**Questions asked:** 8

### Verification Results
- **Tier:** Full (11 phase). Repo chưa có mã, mọi đường dẫn trong plan là `Create`, nên kiểm chứng tập trung vào tuyên bố về công cụ ngoài và môi trường.
- **Claims checked:** 11
- **Verified:** 9 | **Failed:** 1 | **Unverified:** 1

Đã kiểm:
1. VERIFIED Claude Code hook PreToolUse từ chối bằng JSON `hookSpecificOutput.permissionDecision: "deny"` với `permissionDecisionReason`; exit 2 cũng chặn (tài liệu hooks Claude Code).
2. VERIFIED stdin hook có `tool_name`, `tool_input.file_path`, `cwd`, `hook_event_name`.
3. VERIFIED schema `settings.json`: `hooks.<Event>[].matcher`, `hooks[].type: "command"`, `command`.
4. FAILED `MultiEdit` không còn trong danh sách tool của tài liệu hooks; matcher phase 3 đã sửa thành `Edit|Write`.
5. VERIFIED Agent Skills spec: `metadata` là map chuỗi sang chuỗi, khóa nên có tiền tố riêng; `name` phải trùng tên thư mục. Khóa đổi thành `dk_installed_by`, `dk_version`, `dk_hash`.
6. VERIFIED Codex đọc skill từ `~/.codex/skills` và `.codex/skills`, cùng định dạng SKILL.md.
7. VERIFIED Codex hooks tồn tại, sự kiện `PreToolUse`, `PostToolUse`, ổn định từ 0.145, hook dự án cần trust.
8. UNVERIFIED định dạng `hooks.json` của Codex (trường `command` dạng mảng) và payload đường dẫn file với `apply_patch`; phase 11 bước 1 ghi payload thật trước khi viết mã.
9. VERIFIED Go 1.26.3 có trên máy phát triển.
10. VERIFIED không có lệnh `dk` trên PATH máy này.
11. VERIFIED `ak plan use` cần git repo; thư mục chưa `git init`, phase 1 bước 1 xử lý.

#### Questions & Answers

1. **[Assumptions]** Plan giả định `dk init` tự gọi `skill install` và `hook install` cho Claude, scope dự án. Chốt thế nào?
   - Options: Cài luôn trong init (Recommended) | init chỉ tạo docs/ | init hỏi tương tác
   - **Answer:** init chỉ tạo docs/
   - **Rationale:** tách rõ khởi tạo tài liệu và cài công cụ agent; `init` không đụng `.claude/`; ba lệnh rõ ràng trong README.
2. **[Risks]** Hook agent khi `dk` không chạy được: xử lý thế nào?
   - Options: Fail-open kèm cảnh báo (Recommended) | Fail-closed | Fail-open im lặng
   - **Answer:** Fail-open kèm cảnh báo
   - **Rationale:** pre-commit là lớp bảo đảm cứng; hook chỉ tiện.
3. **[Tradeoffs]** Mermaid trong HTML sinh ra: cách nhúng nào?
   - Options: Nhúng inline mỗi file (Recommended) | File chung trong docs/html/_assets/ | Cả hai, cờ --shared-assets
   - **Answer:** Nhúng inline mỗi file
   - **Rationale:** giữ hợp đồng "HTML tự chứa"; chấp nhận 2 đến 3 MB mỗi file có sơ đồ.
4. **[Architecture]** Plan và Report của `doc-plan-report` đặt ở `plans/` ngoài `docs/`, không vào CHANGELOG-DOCS. Giữ không?
   - Options: Giữ ở plans/ (Recommended) | Đưa vào docs/plans/ | Bỏ doc-plan-report
   - **Answer:** Giữ ở plans/
   - **Rationale:** khớp quy ước `plans/` đang dùng; `plans_dir` cấu hình được.
5. **[Architecture]** Repo mã nguồn của CLI `dk` là repo `agent-skill` này hay repo riêng?
   - Options: Repo này, git init tại đây (Recommended) | Repo riêng | Chưa quyết
   - **Answer:** Repo này, git init tại đây
   - **Rationale:** plan và code cùng chỗ; `ak plan use` chạy được sau `git init`; module path chốt khi có remote.
6. **[Scope]** Hỗ trợ Windows ở `v0.1.0`?
   - Options: Linux và macOS trước (Recommended) | Cả 3 hệ từ v0.1.0 | Chỉ Linux
   - **Answer:** Linux và macOS trước
   - **Rationale:** tránh kéo dài phase 10; Windows gộp vào `v0.2.0` ở phase 11.
7. **[Verification]** Tài liệu hooks Claude Code không còn `MultiEdit`. Sửa matcher?
   - Options: Sửa thành Edit|Write (Recommended) | Giữ Edit|Write|MultiEdit
   - **Answer:** Sửa thành Edit|Write
   - **Rationale:** khớp tài liệu hiện hành.
8. **[Assumptions]** Mockup mặc định cho `doc-design-flow`?
   - Options: HTML tĩnh theo tokens (Recommended) | Figma liên kết + ảnh xuất | Quyết theo từng dự án
   - **Answer:** HTML tĩnh theo tokens
   - **Rationale:** lint được bằng `dk check`; Figma vẫn dùng qua `--set external`.

#### Confirmed Decisions
- `dk init` chỉ khởi tạo tài liệu và pre-commit; skill và hook cài bằng lệnh riêng.
- Hook fail-open kèm cảnh báo; pre-commit là lớp cứng.
- Mermaid nhúng inline.
- Plan và Report ở `plans/`, không vào changelog tài liệu.
- Repo CLI là repo này; `git init` ở phase 1.
- `v0.1.0` Linux và macOS; Windows và Codex ở `v0.2.0`.
- Matcher PostToolUse `Edit|Write`.
- Mockup mặc định HTML tĩnh theo tokens.
- Khóa metadata skill: `dk_installed_by`, `dk_version`, `dk_hash`.

#### Action Items
- [x] Sửa hợp đồng, bảng phase, câu hỏi mở trong `plan.md`
- [x] Phase 1: thêm `git init`, `init` in nhắc hai lệnh, không cài skill
- [x] Phase 3: bỏ cài skill và hook khỏi `init`, matcher `Edit|Write`, khóa `dk_*`, success criteria tách `init` và `install`
- [x] Phase 7: mockup mặc định HTML, Figma là ngoại lệ
- [x] Phase 10: bỏ Windows, 4 nền tảng, CI ubuntu và macos
- [x] Phase 11: thêm Windows, bỏ `dk init --target`, release `v0.2.0` 6 nền tảng
- [x] Phase 11 bước 1: đóng mục UNVERIFIED số 8 bằng tài liệu hooks và mã nguồn Codex 0.153.2 (mục 1.17); payload từ phiên thật chờ `codex login`

#### Impact on Phases
- Phase 1: bước 1 và 8 đổi.
- Phase 3: Overview, Requirements, Architecture, bước 8 và 10, Success Criteria đổi.
- Phase 7: quy trình `doc-design-flow` và Risk đổi.
- Phase 10: Requirements, bước 3, 8, 9, Success Criteria, Risk đổi.
- Phase 11: tiêu đề, Overview, Requirements, bước 4, 5, 7, 8, Success Criteria đổi.

### Whole-Plan Consistency Sweep
- Files reread: plan.md, phase-01 đến phase-11
- Decision deltas checked: 6 (init không cài; matcher; khóa metadata; Windows lùi; mockup mặc định; repo và git init)
- Reconciled stale references: 15 (grep `MultiEdit`, `init gọi`, `dk init --target`, `--no-skill`, `--no-hook`, `installed_by`, `6 nền tảng`, `3 hệ điều hành` đều về không ngoài dòng đích)
- Unresolved contradictions: 0
