---
phase: 1
title: "Phase 1: CLI lõi"
status: done
priority: P1
effort: "5d"
dependencies: []
---

# Phase 1: CLI lõi

## Overview

Dựng module Go, khung lệnh cobra, nhúng template bằng `embed`, và 4 lệnh nền: `init` (chưa cài skill), `template`, `new`, `changelog`. Kết thúc phase này `dk new feature-spec bo-loc-don-hang` tạo đúng file, và `dk changelog add` ghi dòng có số dòng khớp git.

## Requirements

- Functional: `dk --version`; `dk init [--force]`; `dk template list|show <type>`; `dk new <type> <slug> [--from <file>] [--set k=v] [--force]`; `dk changelog add <file> --summary <s> --source <src>`; `dk changelog pending`; mọi lệnh có `--json`.
- Non-functional: binary tĩnh (`CGO_ENABLED=0`); không gọi mạng; mã thoát 0 thành công, 1 lỗi I/O hoặc logic, 2 sai cờ, 3 kiểm tra không qua (dùng từ phase 2); không ghi đè file có sẵn nếu thiếu `--force`; thông điệp tiếng Việt có dấu; `go vet` sạch.

## Architecture

- `cmd/dk/main.go` chỉ gọi `cli.Execute()` và `os.Exit` theo mã lỗi kiểu `*cli.ExitError{Code, Msg}`.
- `internal/cli/root.go`: cobra root, cờ toàn cục `--json`, `--cwd` (mặc định thư mục hiện tại), tìm gốc dự án bằng cách đi lên đến thư mục có `dk.toml`; không có thì báo lỗi "chưa `dk init`" trừ lệnh `init`, `template`, `--version`.
- `assets/embed.go`: `//go:embed templates/* types.toml` xuất `var FS embed.FS`.
- `internal/doctype`: struct `Type{Name, Dir, NamePattern, IDScheme, Subdir []string, Required []string, Statuses []string, From map[string]string}`; `Load(fs) (Registry, error)`; validate lúc load: dir không rỗng, IDScheme hợp lệ, template tương ứng tồn tại.
- `internal/frontmatter`: `Split(b []byte) (fm map[string]any, body []byte, ok bool)`; `Join(fm, body)`; giữ thứ tự khóa bằng `yaml.Node`, không dùng map thuần khi ghi lại để diff sạch.
- `internal/tmpl`: `Render(typeName string, data Data) ([]byte, error)` với `text/template`, delims mặc định `{{ }}`, `missingkey=error` để template thiếu trường thì fail ngay trong test.
- `internal/docs`: `New(reg, typeName, slug, opts) (paths []string, error)`: sinh `id` theo scheme (`seq`: quét `Dir` lấy số lớn nhất + 1; `date`: yymmdd theo giờ máy), tính đường dẫn theo `NamePattern`, tạo `Subdir` khi có, chép trường `From` từ frontmatter file nguồn, ghi `created_by: dk`, `dk_version`.
- `internal/gitx`: `Root(dir)`, `NumStat(path) (add, del int, tracked bool, err)`, `ChangedDocs(docsDir) []string` (`git status --porcelain -- docs/`), `IsRepo`. Chạy `git` qua `os/exec`; không có git thì trả cờ `noGit` để changelog ghi "không git" và đếm dòng toàn file.
- `internal/changelog`: model `Entry{Time, Path, Added, Deleted, New bool, Lines int, Summary, Source}`; `Parse(b)`, `Add(e, mergeWindow 10*time.Minute)`, `Format()`. Quy tắc gộp: dòng đầu tiên cùng ngày, cùng `Path`, cùng `Source`, cách nhau dưới 10 phút thì thay dòng đó: giờ mới, số dòng tính lại, tóm tắt mới thay "chưa tóm tắt" hoặc thay tóm tắt cũ. Sau khi ghi, bump `updated:` trong frontmatter file tài liệu (không ghi changelog cho việc bump này).
- `internal/config`: `dk.toml` gồm `project_name`, `id_prefix` (mặc định rỗng), `default_owner`, `docs_dir = "docs"`, `plans_dir = "plans"`, `language = "vi"`.
- Phiên bản: `var Version = "dev"` trong `internal/cli`, gán qua `-ldflags "-X .../internal/cli.Version=..."`.

Cây `docs/` do `init` tạo: `intake/ cr/ features/ adr/ overview/ design/{tokens,atoms,molecules,organisms,templates,patterns,flows,wireframes,mockups}/ test/ release/ ops/ governance/ html/` và `CHANGELOG-DOCS.md` rỗng có tiêu đề. Thư mục rỗng giữ bằng `.gitkeep`. `plans/` không do `init` tạo nếu đã có. Thư mục con thêm ở các phase sau (`plan/`, `release/briefs/`, `ops/runbooks/`...) do `dk new` tự tạo khi thiếu, `init` không cần biết trước.

## Related Code Files

- Create: `go.mod`, `cmd/dk/main.go`, `assets/embed.go`, `assets/types.toml`, `assets/templates/{idea,interview,brief,cr,feature-spec}.md`
- Create: `internal/cli/{root,init,template,new,changelog}.go`, `internal/config/config.go`, `internal/doctype/{doctype,registry}.go`, `internal/frontmatter/frontmatter.go`, `internal/tmpl/tmpl.go`, `internal/docs/new.go`, `internal/gitx/gitx.go`, `internal/changelog/{changelog,merge}.go`
- Create test: `internal/doctype/registry_test.go`, `internal/frontmatter/frontmatter_test.go`, `internal/tmpl/tmpl_test.go`, `internal/docs/new_test.go`, `internal/changelog/changelog_test.go`, `internal/gitx/gitx_test.go`, `internal/cli/cli_test.go` (chạy binary trong `t.TempDir()` có `git init`)
- Create: `Makefile` (`build`, `test`, `vet`, `install`), `.gitignore`, `README.md` ngắn (cài, 5 lệnh đầu)

## Implementation Steps

1. `git init` tại gốc repo này (quyết định validate: repo CLI là repo `agent-skill`); `go mod init` với module path placeholder; `go get` cobra, goldmark (dùng ở phase 2 nhưng khai báo sớm để khóa phiên bản), yaml.v3, go-toml/v2. Ghi `Makefile` với `CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X ...Version=$(git describe --tags --always)" -o bin/dk ./cmd/dk`.
2. Viết `internal/cli/root.go`: root command, `--json`, `--version`, hàm `findProjectRoot`, kiểu `ExitError`. Test: chạy `dk --version` in `dev`.
3. Viết `internal/frontmatter` và test: tách đúng khi có và không có frontmatter, giữ thứ tự khóa khi ghi lại, không đổi thân, CRLF giữ nguyên.
4. Viết `assets/types.toml` cho 5 loại đầu: `idea`, `interview`, `brief` (dir `intake`, subdir `{yymmdd}-{slug}/`, id `none`, name `idea.md|interview.md|brief.md`), `cr` (dir `cr`, id `date:CR-{yymmdd}-{slug}`, name `{id}.md`), `feature-spec` (dir `features`, id `seq:F-{n:03}`, name `{id}-{slug}.md`, `from.brief = {title: title, outcome: purpose, acceptance: acceptance}`). Viết `internal/doctype` và test load, validate, lỗi khi template thiếu.
5. Viết 5 template theo báo cáo mục 1a, 1b, 3: frontmatter dùng `{{.ID}} {{.Title}} {{.Status}} {{.Owner}} {{.Created}} {{.Source}}`; thân là các mục cố định, mỗi mục một chú thích HTML `<!-- gợi ý: ... -->` chứa câu hỏi. Feature Spec đủ 11 mục, mục 4 là khối ```mermaid có nút mẫu `B1[...]`, mục 5 bảng mã bước. CR đủ 7 mục với bảng tác động 6 dòng. Test: render mọi template với dữ liệu mẫu không lỗi, `missingkey=error`.
6. Viết `internal/tmpl` (`Render`) và `internal/docs.New`: sinh id `seq` bằng regexp trên tên file trong `Dir`, `date` theo `time.Now()` (cho phép inject clock để test); từ chối khi file đích tồn tại và không `--force`; với `--from`: đọc frontmatter nguồn, ánh xạ theo `From`, ghi `source: <id nguồn>`. Test: tạo feature-spec đầu tiên ra `F-001`, tạo tiếp ra `F-002`; tạo intake ra thư mục 3 file; `--from brief` chép đúng trường; file tồn tại thì lỗi mã 1.
7. Viết `internal/cli/new.go` và `template.go`: `template list` in bảng `type | dir | id | mô tả`; `template show <type>` in nội dung template thô. `--json` in mảng đối tượng hoặc `{path, id}`.
8. Viết `internal/config` và `internal/cli/init.go`: tạo `dk.toml` (từ chối nếu có, trừ `--force`), tạo cây thư mục, `CHANGELOG-DOCS.md` với dòng đầu `# Changelog tài liệu`; cuối cùng in nhắc: "Tiếp theo: `dk skill install` và `dk hook install`". Phase 3 chỉ thêm cài pre-commit vào lệnh này, không cài skill hay hook.
9. Viết `internal/gitx` và test trong repo tạm: file tracked sửa 3 dòng thêm 1 xóa ra `(3,1,true)`; file untracked ra `tracked=false`; ngoài repo ra `noGit`.
10. Viết `internal/changelog`: `Parse` chịu được file rỗng và file có nhiều ngày; `Add` chèn nhóm ngày mới lên đầu; gộp theo quy tắc; `Format` giữ đúng định dạng cột. Test: thêm 2 dòng cách 5 phút cùng file cùng nguồn ra 1 dòng; cách 11 phút ra 2 dòng; khác nguồn ra 2 dòng; tóm tắt "chưa tóm tắt" bị thay bởi tóm tắt thật.
11. Viết `internal/cli/changelog.go`: `add` nhận đường dẫn tuyệt đối hoặc tương đối, chuẩn hóa về tương đối `docs/`; từ chối file ngoài `docs/`, file trong `docs/html/`, file chỉ mục (đánh dấu `generated: true` trong frontmatter); bump `updated:`. `pending`: `gitx.ChangedDocs` trừ những file đã có dòng trong changelog kể từ commit HEAD (so bằng đường dẫn và mốc thời gian commit HEAD), in danh sách, mã thoát 1 nếu còn.
12. Viết `internal/cli/cli_test.go`: kịch bản end-to-end trong thư mục tạm: `git init`, `dk init`, `dk new feature-spec bo-loc-don-hang`, sửa 2 dòng, `dk changelog add`, đọc `CHANGELOG-DOCS.md` khớp `git diff --numstat`, `dk changelog pending` trả rỗng; sửa tiếp không add, `pending` liệt kê file và mã thoát 1.
13. `make vet test build`; chạy tay `bin/dk template list`, `bin/dk template show cr`.

## Success Criteria

- [x] `make build` ra binary tĩnh (`file bin/dk` báo statically linked); `dk --version` in phiên bản
- [x] `dk init` tạo đủ cây `docs/`, `dk.toml`, `CHANGELOG-DOCS.md`; chạy lại báo lỗi nếu thiếu `--force`
- [x] `dk template list` liệt kê 5 loại; `dk template show feature-spec` in đủ 11 mục
- [x] `dk new feature-spec bo-loc-don-hang` tạo `docs/features/F-001-bo-loc-don-hang.md` với frontmatter đủ trường mục 1.3 plan; lần hai ra `F-002`
- [x] `dk new brief --from docs/intake/<dir>/idea.md` chép đúng trường theo `types.toml`
- [x] `dk changelog add` ghi dòng số dòng khớp `git diff --numstat`; gộp trong 10 phút hoạt động; `updated:` được bump
- [x] `dk changelog pending` mã thoát 1 khi còn file chưa có changelog, 0 khi sạch
- [x] `go test ./...` và `go vet ./...` qua; test không phụ thuộc mạng

## Risk Assessment

- **Đếm dòng lệch với git** khi file có thay đổi đã stage một phần: dùng `git diff --numstat HEAD -- <file>` (working tree so HEAD) để bao cả staged và unstaged. Tín hiệu vỡ: test end-to-end lệch số; phản ứng: chốt một cách đếm duy nhất và ghi vào README.
- **Thứ tự khóa frontmatter bị xáo** khi ghi lại bằng map: dùng `yaml.Node`. Tín hiệu vỡ: diff tài liệu có dòng đổi chỗ; phản ứng: sửa trong `frontmatter`, thêm test golden.
- **Template Feature Spec chứa `{{` trong Mermaid hoặc ví dụ**: đặt delims `text/template` thành `[[ ]]` nếu xảy ra; quyết định trong bước 5, không đổi sau.
- **Module path placeholder** làm `go install` từ xa không chạy: chỉ ảnh hưởng phase 10; đổi bằng `go mod edit -module` và sed import một lần.
