---
title: Review phase 2 - render và kiểm tra
date: 2026-09-03
phase: 2
verdict: DONE_WITH_CONCERNS
---

# Review phase 2: render, index, check, refs, status

Phạm vi: toàn bộ file mới/sửa của phase 2 (repo chưa có commit, đọc file trực tiếp).
Kiểm chứng chạy được: `gofmt -l .` sạch, `go vet ./...` sạch, `go test ./...` mọi package pass.
Đo thật: 200 feature-spec, `dk render --all --index` hết **1,55 s** (yêu cầu < 5 s) → đạt; đầu ra `docs/html` **684 MB**.

## a) Success Criteria

| # | Tiêu chí | Verdict | Bằng chứng |
|---|---|---|---|
| 1 | render 1 file ra HTML tự chứa, mở `file://` hiện Mermaid, không request mạng | Đạt (chưa mở trình duyệt thật) | `phase2_test.go:55` chặn `<link>` và `src="http"`; bundle là IIFE gán `globalThis["mermaid"]` (cuối `mermaid.min.js`), không `import` |
| 2 | File không Mermaid không chứa script mermaid | Đạt | `render_test.go:48`, `phase2_test.go:52` |
| 3 | Mã bước có liên kết nhảy đúng dòng bảng hành vi | Đạt, có điều kiện | `render_test.go:75-86`; xem F-2 (chọn nhầm bảng) |
| 4 | `check` bắt mã bước lệch, link hỏng, thiếu source, > 800 dòng; mã thoát 3 | Đạt | `check_test.go:41-53`, `phase2_test.go:120-144` |
| 5 | `index all` ra 4 file `generated: true`; `changelog pending` bỏ qua | Đạt | `phase2_test.go:73-97` |
| 6 | `refs` in đúng liên kết đi/đến | Đạt | `graph_test.go`, `phase2_test.go:146-157` |
| 7 | `status` khớp fixture, `--json` hợp lệ | Đạt | `phase2_test.go:162-175` |
| 8 | `style.css` không vi phạm `html-style.md` theo test | Đạt một phần | `render_test.go:159-176` chỉ kiểm 6 điểm; xem N-3 |

## b) Regression phase 1

Không thấy regression. `init`, `new`, `changelog add/pending` giữ nguyên hành vi; mã thoát 0/1/2/3 giữ nguyên (`root.go:24-29`, `unknown command` → 2 xác nhận bằng tay).
`config.Write` nay sinh thêm bảng `[check]` trong `dk.toml`: **chấp nhận được**. `Load` (`config.go:58-63`) điền mặc định khi thiếu, nên `dk.toml` cũ vẫn chạy; `init --force` vẫn giữ file cũ; `cli_test.go:97,109` không assert nội dung đầy đủ nên không test nào vỡ.
`changelog.go` tách `pendingDocs` và dùng `docs.CountLines` - đúng DRY, không đổi contract.

## c) Public contract

- Cờ và tham số khớp Requirements: `render [<file>] --all --index`, `index [features|adr|cr|intake|all]`, `check [<file>] --strict`, `refs <file>`, `status`; `--json` là cờ persistent.
- JSON: `render` → mảng `{src,out,mermaid}`; `check` → mảng `Finding{file,line?,rule,level,msg}`; `refs` → `{file,out,in}`; `status` → `{docs,open_cr,changelog_pending,check}` với `changelog_pending` null khi không git; `index` → mảng đường dẫn. Đều ổn định, mảng rỗng in `[]` chứ không `null`.
- **Nguồn Markdown không bị sửa**: `render`, `check`, `refs`, `status` chỉ đọc; ghi duy nhất vào `docs/html/` (`render.go:114`). `index` ghi `README.md` là bản sinh có chủ đích. Đúng hợp đồng.

## d) Theo pattern phase 1

Đạt: `fail(code, ...)`, `a.printJSON`, `a.relRoot`, `a.requireProject`, `exactArgs`, comment và thông điệp tiếng Việt, tên file snake_case, một file một nhóm lệnh.

## e) Rủi ro đúng sai

### Should-fix

**F-1. `dk check <file>` với đường dẫn sai hoặc ngoài dự án trả 0 "không lỗi"** - `internal/cli/check.go:53-65`.
Lọc theo `f.File == rel`; file không nằm trong tập đã quét thì danh sách rỗng → in "0 lỗi, 0 cảnh báo", mã thoát 0. Kiểm thật: `dk check /etc/hosts` và `dk check docs/features/nope.md` đều trả 0. Hook hoặc skill gõ sai tên file sẽ được báo "sạch" (false green).
Sửa: sau khi tính `rel`, nếu không có `meta.Rel == rel` trong kết quả `docs.Scan` thì `fail(codeError, ...)` như `runRefs` đang làm (`refs.go:33-42`).

**F-2. `findBehaviorTable` bỏ điều kiện tiêu đề "Bảng hành vi" của spec** - `internal/render/steps.go:33-52`.
Spec ghi "bảng đầu tiên sau tiêu đề chứa 'Bảng hành vi' **hoặc** cột đầu khớp `^B\d+$`"; code chỉ còn vế sau. Trong `assets/templates/feature-spec.md` bảng hành vi (dòng 49) đứng trước bảng Giao diện (dòng 61) nên hiện đúng, nhưng chỉ do thứ tự template. Fixture `render_test.go:84-86` đang khẳng định chính hành vi chọn nhầm bảng này. Hệ quả: nếu người viết đảo mục, cả `id="step-*"` lẫn quy tắc `step-codes` đều bám nhầm bảng và báo lỗi giả.
Sửa: ưu tiên bảng nằm dưới heading chứa "Bảng hành vi", chỉ fallback sang "cột đầu khớp `^B\d+$`" khi không có.

**F-3. HTML sinh ra nhúng mermaid.min.js mà không kèm thông báo giấy phép MIT** - `assets/embed.go:6`, `assets/html/page.html:24`.
Bundle minified **không** chứa dòng bản quyền của mermaid (`grep "Knut Sveidqvist"` = 0); `MERMAID-LICENSE.txt` không nằm trong `//go:embed` nên không đi theo binary. Mỗi file HTML là một bản phân phối lại mermaid, MIT yêu cầu kèm notice.
Sửa: thêm `html/MERMAID-LICENSE.txt` vào `//go:embed`, in ra một `<!-- ... -->` ngay trước `<script>` mermaid (chỉ khi có nhúng).

**F-4. `dk render <file>` với file không phải `.md` im lặng trả 0** - `internal/cli/render.go:82`.
Kiểm thật: `dk render docs/html/index.html` → không in gì, mã 0, không ghi file. Bộ lọc `.md` đúng cho `--all` nhưng sai cho một file chỉ định rõ.
Sửa: khi `file != ""` mà đuôi không phải `.md` thì `fail(codeUsage, ...)`.

### Nit

**N-1. `filepath.Rel` prefix `".."` bắt nhầm tên thư mục bắt đầu bằng `..`** - `render.go:85`, `check/rules_links.go:50`, `changelog.go:70` (có sẵn từ phase 1). Thư mục tên `..data` (phổ biến trong volume k8s) bị coi là nằm ngoài. Sửa: so `rel == ".."` hoặc `strings.HasPrefix(rel, ".."+string(filepath.Separator))`.

**N-2. `line-threshold` áp cho `plans/`**. `runChecks` quét `docs/` + `plans/` (`check.go:37`), `lineThreshold` duyệt toàn bộ `Metas`. Kiểm thật: file `plans/big.md` 902 dòng → `error`, `dk check` mã 3. README dòng 44 đã nói rõ check quét `plans/`, nên đây là quyết định có chủ đích; nhưng ngưỡng tài liệu 500/800 áp cho plan file sẽ làm `dk check` đỏ vì plan dài. Cân nhắc miễn `plans/` khỏi `line-threshold` hoặc cho ngưỡng riêng.

**N-3. Test style chỉ phủ 6/12 quy tắc `html-style.md`**. Không kiểm `!important` ngoài `@media print`, không kiểm cỡ chữ tối thiểu, không kiểm emoji trong tiêu đề. Vi phạm nhỏ đang có: `.status` = 0.85em trong `.meta` 0.85rem → ~0.72rem, dưới ngưỡng "không nhỏ hơn 0.8rem cho nội dung đọc" (`style.css:59,65`).

**N-4. `render --all` render cả `README.md` sinh ra và `CHANGELOG-DOCS.md`** thành HTML, nhưng `Index` bỏ qua `generated` (`index.go:22`) nên các trang này không được liên kết từ đâu. Rác nhỏ trong `docs/html/`.

**N-5. `refs`: id ngắn khớp nhầm khi id dài không có trong tập đã quét.** `idRegexp` sắp dài trước (`graph.go:87`) xử lý đúng khi cả hai id tồn tại (test `graph_test.go:38` chứng minh), nhưng nhắc `CR-260910-loc-theo` trong lúc chỉ `CR-260910-loc` tồn tại vẫn tạo cạnh sai vì `\b` khớp trước dấu `-`. Ít gặp, ghi nhận.

**N-6. `backlink` "tài liệu chết" chỉ đếm tham chiếu từ `typed()`.** `rules_backlink.go:16-23` bỏ qua tài liệu loại chưa có trong `types.toml` (ví dụ `docs/test/*` hiện chưa có loại `test`), nên một brief chỉ được file test trỏ về vẫn bị cảnh báo giả.

**N-7. `docs.Resolve` quét tuyến tính trong vòng lặp** (`rules_backlink.go:20,30`, `graph.go:52,57`) → O(n²)/O(n²·ids). 200 file vẫn nhanh (đo 1,55 s cho render, check tức thì), chưa cần sửa; dùng `docs.ByID` nếu tập tài liệu lớn lên.

**N-8. Lỗi `filepath.Rel` bị nuốt** ở `render/index.go:46` và `cli/index.go:102` (`rel, _ :=`) → href rỗng thay vì báo lỗi.

**N-9. `config.Load` coi `warn_lines = 0` là chưa đặt** (`config.go:58`), không tắt được ngưỡng. Ghi vào README hoặc dùng con trỏ.

**N-10. `Windows` (v0.2.0, chỉ ghi nhận)**: `IsRelativeLink("C:/x")` trả false (coi `C:` là scheme) - đúng cho URL, nhưng đường dẫn tuyệt đối Windows và dấu `\` trong link Markdown chưa có đường xử lý riêng.

### Đã kiểm và **không** phải vấn đề

- **XSS / HTML thô**: `html.WithUnsafe()` không bật, `render_test.go:94` chứng minh `<script>alert` bị escape. Link `javascript:` bị goldmark chặn sẵn (`renderer/html/html.go:518,591,620` gọi `IsDangerousURL`).
- **Metadata**: khung dùng `html/template`, mọi trường frontmatter đi qua auto-escape (`page.html:19`); `render/index.go` escape thủ công đủ cho attribute trong nháy kép.
- **`template.JS` cho mermaid**: nội dung là asset nhúng, tin cậy; `render.go:77` còn thay `</script` phòng khi nâng phiên bản.
- **CRLF**: `frontmatter.Split` tách theo newline phát hiện được và ghép lại nguyên dạng, nên `offset = len(Raw)-len(Body)` (`rules_links.go:21`) và số dòng của `link-broken` đúng cho cả LF và CRLF; test dòng 10 trong `check_test.go:75` khớp.
- **Link có `#fragment`**: `render.go:201-216` giữ neo và đổi `.md` → `.html` (test `href="../adr/ADR-0001-x.html#muc-2"`); `check` và `refs` cắt neo trước khi stat/so khớp.
- **Concurrency**: không có goroutine; `sync.Once` cho asset là an toàn.
- **`docs.Scan` bỏ `html/`** chỉ ở cấp một dưới base (`scan.go:51`) - đúng yêu cầu, `docs/features/html/` vẫn được quét.
- **`--index` và `dk index`** đều loại `generated`; `changelog pending` loại qua `tracksChangelog` (`changelog.go:81-93`).
- **`status` gọi git** có bảo vệ `gitx.IsRepo` (`status.go:50`), ngoài repo trả `changelog_pending: null` chứ không lỗi.

## f) Lint / vet / build

`gofmt -l .` không ra file nào. `go vet ./...` sạch. `go test ./...` pass toàn bộ (12 package, 2 package không có test: `assets`, `cmd/dk`, `internal/config`).
Gợi ý nhỏ: `internal/config` chưa có test dù nay có logic mặc định `[check]`.

## Câu hỏi chưa giải quyết

1. `dk check` có nên đỏ vì file trong `plans/` vượt 800 dòng không (N-2)?
2. `docs/html` 684 MB cho 200 spec - giữ nguyên hay kéo `--shared-assets` từ mục Risk lên sớm hơn?
3. Tiêu chí 1 chưa được xác nhận bằng trình duyệt thật ở chế độ offline; ai chạy bước 12 của Implementation Steps?
