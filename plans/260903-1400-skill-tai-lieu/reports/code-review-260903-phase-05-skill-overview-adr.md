# Code review phase 5: doc-overview, doc-adr

Ngày: 2026-09-03. Phạm vi: 4 loại mới trong `assets/types.toml`, 4 template, 2 skill, `internal/check/rules_adr.go`, `rules_glossary.go`, `rules_threshold.go`, `internal/gitx/gitx.go` (`HeadFile`), `internal/cli/index.go` (`flatRows`), test kèm theo. Khoảng 1460 dòng đọc.

Kiểm chạy: `go vet ./...` sạch, `go test ./...` qua toàn bộ 14 gói. Không sửa mã.

Kết luận: đủ 5 tiêu chí nghiệm thu bằng CLI. Bốn vấn đề mức Medium nên xử lý trước phase 6, không có Blocker.

## High

Không có.

## Medium

### M1. `glossary-term` chỉ đọc **một** file Glossary, trong khi rules.md dạy tách Glossary

`internal/check/rules_glossary.go:24-32` lấy `glossary` đầu tiên theo thứ tự `Rel` rồi `break`. `assets/skills/doc-overview/references/rules.md:40` lại quy định "Vượt 300 thuật ngữ thì tách theo miền nghiệp vụ". Khi có 2 file Glossary, thuật ngữ định nghĩa ở file sắp sau bị báo thiếu.

Kiểm chứng: `docs/overview/glossary.md` định nghĩa "Đơn hàng", thêm `docs/overview/glossary-kho.md` (sắp trước theo chữ cái) thì `dk check` báo `thuật ngữ **đơn hàng** chưa có trong docs/overview/glossary-kho.md`.

Sửa: gộp thuật ngữ của mọi tài liệu loại `glossary`, thông điệp nói "chưa có trong Glossary" hoặc liệt kê file đã quét.

```go
var terms = map[string]bool{}
var found int
for _, m := range c.typed() {
    if m.Type != "glossary" { continue }
    found++
    for t := range glossaryTerms(m.Body) { terms[t] = true }
}
if found == 0 { return nil }
```

### M2. Thứ tự finding của `glossary-term` không ổn định giữa các lần chạy

`internal/check/rules_glossary.go:36` duyệt map trả về từ `firstBold`; `internal/check/check.go:64-72` chỉ phân định theo `File`, `Line`, `Rule` nên nhiều thuật ngữ **cùng một dòng** bị hoán vị ngẫu nhiên. Ảnh hưởng `dk check --json`, diff của pre-commit và mọi test golden sau này.

Kiểm chứng: một dòng mục 2 với 5 thuật ngữ chưa định nghĩa, chạy 6 lần ra 5 thứ tự khác nhau.

Sửa: sắp khóa trước khi tạo finding, hoặc thêm nấc `out[i].Msg < out[j].Msg` vào `sort.SliceStable` ở `check.go:71`.

```go
terms := firstBold(m)
keys := make([]string, 0, len(terms))
for t := range terms { keys = append(keys, t) }
sort.Strings(keys)
for _, term := range keys { ... }
```

### M3. `adr-immutable` chỉ bắt được thay đổi **chưa commit**

`internal/check/rules_adr.go:26-38` so với HEAD, đúng như quyết định 1.11. Nhưng `assets/hooks/pre-commit.sh:5` chỉ chạy `dk changelog pending`, không chạy `dk check`. Vì vậy sửa thân một ADR `accepted`, thêm dòng changelog rồi commit là qua được; sau commit HEAD chính là thân đã sửa nên `dk check` sạch vĩnh viễn. Tiêu chí (b) chỉ đạt cho bản làm việc.

Đây là hệ quả của cơ chế đã chốt, không phải lỗi cài đặt, nên nêu để chủ dự án chọn:

- Thêm `dk check` vào `assets/hooks/pre-commit.sh` (chặt nhất, làm pre-commit chậm hơn).
- Hoặc ghi rõ giới hạn này ở `README.md:69` và `assets/skills/doc-adr/references/rules.md:19`, giữ nguyên mã.

### M4. `HeadFile` sinh 6 tiến trình git cho **mỗi** ADR

`internal/gitx/gitx.go:150-176`: `IsRepo` (1), `Root` gọi `IsRepo` lần nữa (2) rồi `rev-parse --show-toplevel` (3), `hasHead` (4), `cat-file -e` (5), `git show` (6). `rules_adr.go:26` gọi trong vòng lặp nên chi phí tuyến tính theo số ADR: repo 50 ADR mất 0.29s chỉ cho `check`, 200 ADR sẽ hơn 1s cho mỗi lần chạy (kể cả trong hook).

Sửa: giải một lần cho cả `Context` rồi tái dùng. Ví dụ thêm `gitx.Open(root) (*Repo, error)` giữ `root` và `hasHead`, `Repo.File(rel)` chỉ chạy `cat-file` và `show`; hoặc lọc trước bằng một lệnh `git diff --name-only HEAD -- <docs>/adr` rồi chỉ đọc HEAD của file có tên trong danh sách.

## Low

- **L1** `rules_glossary.go:78-92` quét cả chữ in đậm trong khối mã. Kiểm chứng: `// **macro**` trong khối ```go ở mục 5 bị báo. Sửa: bật cờ khi gặp dòng bắt đầu bằng ``` hoặc ~~~ và bỏ qua tới dòng đóng.
- **L2** `rules_glossary.go:36-42` khử trùng theo đúng hoa thường trong khi tra cứu lại hạ chữ, nên `**Kho**` và `**kho**` chưa định nghĩa cho 2 cảnh báo. Biến `seen` là mã chết: khóa của `firstBold` vốn đã duy nhất. Sửa: dùng `strings.ToLower(term)` làm khóa của `first` và bỏ `seen`.
- **L3** `rules_glossary.go:25` nhận cả tài liệu `generated: true` và cả file trong `plans/` có `type: glossary`. Sửa: duyệt `c.typed()` thay cho `c.Metas`.
- **L4** `gitx.go:172-176` dùng `cmd.Output()` nên mất stderr của git; cảnh báo `không đọc được bản HEAD: exit status 128` không nói được lý do. Sửa: dùng lại `run(dir, "show", spec)` như các hàm khác trong gói.
- **L5** Đổi tên file một ADR đã `accepted` làm `HeadFile` trả `ok=false` và rule bỏ qua, thân sửa kèm theo không bị bắt. Sửa nhẹ: ghi giới hạn vào `rules.md`, hoặc dò đổi tên bằng `git diff --find-renames --name-status HEAD`.
- **L6** Không có quy tắc nào kiểm `supersedes` và `superseded_by` trỏ tới ADR có thật và khớp nhau (ADR-0002 `supersedes: ADR-0001` mà ADR-0001 vẫn `accepted` thì im lặng). Ngoài phạm vi phase 5; đề xuất đưa vào phase 6 cùng `backlink`.
- **L7** `glossaryTerms` (`rules_glossary.go:52-67`) chỉ bỏ dòng tiêu đề khi ô đầu đúng chữ "Thuật ngữ", và nhận cột đầu của **mọi** bảng trong file Glossary. Bảng khác trong mục 3 sẽ lặng lẽ thành thuật ngữ hợp lệ (báo sót). Sửa: chỉ quét bảng nằm dưới tiêu đề mục 2.

## Nit

- **N1** `internal/doctype/registry_test.go` không khẳng định `warn_lines` đọc đúng cho `product-overview`, `glossary`, `architecture`; khóa mới chỉ được phủ gián tiếp qua `check_test.go:54`. Thêm một dòng `if reg["glossary"].WarnLines != 300`.
- **N2** `assets/skills/doc-overview/SKILL.md:51,71,95` truyền `<slug-san-pham>` nhưng ba loại này có `name` cố định nên slug chỉ ảnh hưởng `title` mặc định (và bị `--from` ghi đè). Nên nói rõ để người dùng không tưởng slug đổi tên file.
- **N3** `phase-05-skill-overview-adr.md:14` còn ghi `dk new overview --from <brief>`, tên loại thật là `product-overview` (đúng ở mục Architecture và ở 1.11).
- **N4** `rules_adr.go:33` bỏ qua im lặng khi frontmatter bản HEAD không phân tích được; ADR hỏng YAML ở HEAD mất luôn lớp bảo vệ bất biến.

## Đã kiểm và đạt

- `gitx.HeadFile` đúng khi gốc dự án lồng trong gốc git: `docs.Meta.Path` là tuyệt đối nên `filepath.Rel(gitRoot, path)` cho đúng spec `HEAD:<path>`. Kiểm chứng bằng repo `/tmp/gnest` với dự án ở `proj/`: sửa thân ADR `accepted` vẫn báo lỗi.
- Số dòng của `glossary-term` đúng: `fmLines = CountLines(Raw) - CountLines(Body)` khớp cả khi file không có xuống dòng cuối và khi dùng CRLF.
- `adr-immutable` dung thứ CRLF và khoảng trắng cuối dòng, bỏ qua khi không có git, chưa có HEAD, file chưa vào HEAD, và cho sửa thân ở lần chuyển `proposed` sang `accepted` (`rules_adr_test.go` phủ đủ 5 nhánh).
- Không có từ cấm (`ak-`, `Claude Code`, `.claude/`, `.codex/`, `Edit`, `Write`, `MultiEdit`) trong `assets/skills/`; `content_test.go` chạy đủ 4 skill.
- Lệnh trong hai SKILL.md khớp CLI thật: `dk new <type> <slug> --from --set --force`, `dk index [features|adr|cr|intake|all]`, `dk check [<file>]`, `dk changelog add <file> --summary --source`, `dk render [--all --index]`, `dk refs <file>`, `dk --version`.
- `dk check <file>` chạy toàn bộ quy tắc rồi mới lọc theo file, nên `glossary-term` và `adr-immutable` vẫn hoạt động khi chỉ định một file.
- Không thấy hồi quy cách nhóm của `cr` và `intake`: `phase2_test.go:100-119` vẫn khẳng định nhóm `## chưa có brief` trước `## approved` và interview không thành dòng CR.
- Hợp đồng công khai đổi có chủ đích và đã ghi ở `README.md:69,71`: khóa `warn_lines` trong `types.toml`, `dk index features|adr` chuyển từ nhóm theo trạng thái sang bảng phẳng.

## Tiêu chí nghiệm thu

| # | Tiêu chí | Kết quả |
|---|---|---|
| a | Feature catalog khớp số file, chỉ sinh bằng `dk index features` | Đạt (`phase5_test.go:32-38` so với `filepath.Glob`) |
| b | ADR mới thay ADR cũ, `adr-immutable` bắt sửa thân | Đạt cho bản chưa commit; xem M3 |
| c | Product overview từ brief, chép 4 trường, `warn_lines` 150 | Đạt |
| d | `glossary-term` mục 2, 5, 8, bỏ qua khi chưa có Glossary | Đạt; xem M1, L1 |
| e | Test nội dung skill qua với 4 skill | Đạt |

## Việc nên làm, theo thứ tự

1. M1 gộp mọi file Glossary.
2. M2 sắp thuật ngữ trước khi tạo finding.
3. M4 tái dùng một lần giải repo cho cả lượt `check`.
4. M3 chọn phương án: thêm `dk check` vào pre-commit, hay ghi rõ giới hạn.
5. L1, L2, L3, L4 gom một lượt trong `rules_glossary.go` và `gitx.go`.

## Câu hỏi chưa giải quyết

- Dự án được phép có nhiều file Glossary hay không? rules.md nói có, `types.toml` cho `name = "glossary.md"` cố định nên `dk new glossary` lần hai sẽ trùng file. Nếu chỉ cho một file thì M1 chuyển thành: báo lỗi khi thấy file loại `glossary` thứ hai.
- M3 để chủ dự án quyết vì nó đổi thời gian chạy pre-commit.
