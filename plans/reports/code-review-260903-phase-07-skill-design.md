# Code review phase 7: Skill Design (doc-design-system, doc-design-flow)

- Ngày: 2026-09-03
- Phạm vi: toàn bộ mã phase 7 theo danh sách trong yêu cầu (8 loại họ Design, `dk tokens css`, rule `mockup-tokens` và `userflow-steps`, hai skill), cộng các file bị sửa (`docs.Scan`, `docs.New`, `tmpl.Raw`, `changelog.Record`, `check`, `cli/render`, `cli/root`).
- Cơ sở: `plan.md` mục Hợp đồng, 1.1 đến 1.12; `phase-07-skill-design.md`; `reports/phase-07-run.md`.
- Kiểm chạy thật: `go build ./...`, `go vet ./...` sạch, `gofmt -l .` trống, `go test ./...` xanh (17 package). Không dependency mới (`go.mod` vẫn chỉ cobra, goldmark, yaml.v3, go-toml/v2).
- Kiểm bằng dự án mẫu `/tmp/proj` (git init, `dk init`, binary build từ `cmd/dk`): chạy lại toàn chuỗi lệnh của hai SKILL.md và các ca biên; mọi finding dưới đây có bước tái tạo.

## Kết luận ngắn

Phần lõi (`internal/tokens`, `frontmatter` kinds, `docs.New` cho `json|html`, `changelog.Record`, lọc Markdown ở `render`) làm đúng, test là test thật (chứng minh hành vi, không phải test rỗng), vòng JSON và HTML round-trip ổn định và giữ JSON hợp lệ. Ba vấn đề thật cần xử lý trước khi coi phase 7 xong: (1) `mockup-tokens` có ba đường vòng và một false positive không thể sửa được, nên tiêu chí 2 chỉ đạt một phần; (2) quy tắc "Đặt tên mockup" trong `doc-design-flow/references/rules.md` xung đột trực tiếp với rule `userflow-steps`; (3) một lệnh mẫu trong `doc-design-system/SKILL.md` chạy là lỗi `dk check`, và một bước của cùng skill không thể thực hiện được bằng CLI hiện có.

## High

### H1. `mockup-tokens` bị vô hiệu hoàn toàn khi mockup không có khối `<!-- dk:` hợp lệ

- `internal/check/check.go:93` (`typed()` yêu cầu `m.HasFM`), `internal/frontmatter/kinds.go:46` (`SplitHTMLComment` yêu cầu dòng đầu đúng `<!-- dk:`).
- Vấn đề: file `docs/design/mockups/bypass.html` chỉ có `<style>body { color: #333; padding: 12px; }</style>` và không có khối metadata → `dk check` báo **0 lỗi, 0 cảnh báo**. Không rule nào báo "file trong `docs/` thiếu metadata", nên mockup viết tay (hoặc metadata gõ sai, thiếu dấu `-->`) thoát sạch lint. Tiêu chí phase "dk check bắt giá trị hex hoặc px gõ tay" chỉ đúng với mockup do `dk new` sinh. Tái tạo: tạo file trên rồi `dk check`.
- Cách sửa: thêm một rule (hoặc mở rộng `frontmatter-required`) báo `error` cho mọi file trong `docs/` có đuôi được quét mà `HasFM == false` (trừ `CHANGELOG-DOCS.md`, file `generated`); tối thiểu áp cho `design/mockups/*.html` và `design/tokens/*.json`.

### H2. `styleAttrRe` chỉ khớp dấu nháy kép, nháy đơn thoát lint

- `internal/check/rules_mockup.go:13`: `\bstyle\s*=\s*"[^"]*"`.
- Vấn đề: `<div style='color:#333;padding:9px'>` không bị bắt (đã tái tạo: `dk check` không báo dòng đó, trong khi cùng file bắt `#abc`, `0px`, `768px` ở `<style>`). HTML cho phép nháy đơn, nên đây là đường vòng một dòng.
- Cách sửa: `\bstyle\s*=\s*("[^"]*"|'[^']*')`.

### H3. `px` trong `@media` bị báo lỗi mà không có cách sửa

- `internal/check/rules_mockup.go:15`.
- Vấn đề: `@media (min-width: 768px)` trong `<style>` bị báo `error mockup-tokens: giá trị gõ tay "768px"` (đã tái tạo). CSS custom properties **không dùng được** trong điều kiện media query, nên mockup có điểm gãy (đúng `foundations.md` mục 3 yêu cầu) không bao giờ qua được `dk check`; agent chỉ còn cách bỏ responsive hoặc bỏ lint.
- Cách sửa: bỏ các đoạn `@media (...)` và `@container (...)` khỏi vùng quét trước khi tìm `px`, và ghi giới hạn này vào `rules.md` (điểm gãy gõ số là hợp lệ).

### H4. "Đặt tên mockup" (hậu tố trạng thái `B3e`) xung đột với rule `userflow-steps`

- `assets/skills/doc-design-flow/references/rules.md:54-57` ("đặt `step` là mã bước kèm hậu tố ... `B3e` cho rỗng ... **và ghi ở mục 3 của user flow**") so với `internal/check/rules_userflow.go:41-50` và `render.StepCodes` (cột đầu bảng mã bước).
- Vấn đề: mục 3 của template `userflow` là bảng có cột đầu là "Mã bước", đúng bảng mà `render.StepCodes` đọc. Ghi `| B3e | Kết quả rỗng | ... |` vào bảng đó cho ra `error userflow-steps: mã bước không có trong F-001: B3e` (đã tái tạo). Ngoài ra `B3e` dùng chung không gian tên với hậu tố "bước chèn giữa" của quyết định 1.12 (`^B\d+[a-z]?$`), nên khi spec thật thêm bước `B3e` thì tên file mockup đụng nghĩa.
- Cách sửa: ghi rõ trong `rules.md` rằng mã có hậu tố trạng thái **chỉ** được xuất hiện ở cột "Mockup" (dạng liên kết), không được xuất hiện ở cột đầu; hoặc bỏ hậu tố kiểu bước và đổi tên file mockup theo `state` (`<F>-<B>-empty.html`, cần thêm `{state}` vào `name` của loại `mockup` trong `types.toml`). Phương án hai giải luôn L6.

### H5. Lệnh mẫu và một bước của `doc-design-system` không chạy được

- `assets/skills/doc-design-system/SKILL.md:113` và `SKILL.md:119-120`.
- Vấn đề 1: `dk new mockup kiem-tra --set feature=F-000 --set step=B0 --set owner=... --set title=...` không đặt `source`, mà `mockup` có `source` trong `required` → ngay khối "Sau mỗi lần sửa" (`dk check`) trả `error frontmatter-required: thiếu trường source`, mã 3 (đã tái tạo). Skill nói "check còn lỗi thì sửa và chạy lại, không báo người" nên agent sẽ loay hoay. `reports/phase-07-run.md` không chạy lệnh này (bản chạy thử dùng `--from <wireframe>`, nên có `source`).
- Vấn đề 2: bước "xóa file kiểm tra ... và ghi changelog 'xóa mockup kiểm tra'" không thực hiện được: `changelog.Record` đọc file trước khi ghi (`internal/changelog/record.go:59`), nên `dk changelog add` trên file đã xóa trả `open ...: no such file or directory`, mã 1 (đã tái tạo). `dk` cũng không có lệnh xóa tài liệu.
- Cách sửa: thêm `--set source=<thư mục intake của brief>` (hoặc `--from <brief>`) vào lệnh mẫu; đổi bước 6 thành "giữ mockup kiểm tra như một tài liệu bình thường" hoặc "xóa file và xóa dòng changelog tương ứng bằng tay", đừng yêu cầu `changelog add` cho file đã xóa.

## Medium

### M1. Giá trị `--set` đi thẳng vào đường dẫn file, không kiểm tra

- `internal/docs/new.go:289-307` (`fill`), dùng ở `new.go:98` và `new.go:121-125`.
- Vấn đề: chỉ `step` (`stepRe`) và `layer` (`Layers`) được kiểm; `feature` thì không. `dk new mockup e2 --set feature=../wireframes/esc --set step=B1` ghi ra `docs/design/wireframes/esc-B1.html` (đã tái tạo); `--set feature=../../../../tmp/evil` cố ghi `/tmp/tmp/evil-B1.html` (chỉ trượt vì `MkdirAll` tạo thư mục loại, không tạo thư mục của đường dẫn cuối). Khoảng trắng cũng qua: `--set feature="F 0 0"` tạo `docs/design/flows/F 0 0-flow.md`.
- Cách sửa: kiểm mọi giá trị dùng cho `dir`/`name` theo một biểu thức chặt (ví dụ `^[A-Za-z0-9][A-Za-z0-9._-]*$`), từ chối `/`, `\`, `..` và khoảng trắng với mã 2.

### M2. `line-threshold` áp cho `.json` và `.html`, `max_lines` 800 không tắt được, `design-tokens` không có `warn_lines` riêng

- `internal/check/rules_threshold.go:62-79` (quét `c.Metas`, không cần `HasFM`), `assets/types.toml:100-108`.
- Vấn đề: `tokens.json` thật của một design system dễ vượt 500 đến 800 dòng (khung mẫu đã 123 dòng cho 57 token); vượt 800 là `error` cứng và theo quyết định 1.8 không tắt được. Đã tái tạo bằng một `.json` 901 dòng trong `docs/design/`: `error line-threshold: 902 dòng, vượt max_lines 800`. Ngoài ra mọi `.json` khác nằm trong `docs/` (ví dụ `tokens.schema.json`) cũng chịu ngưỡng dù không phải tài liệu.
- Cách sửa: bỏ ngưỡng cho loại có `kind != md` (hoặc thêm `warn_lines`/`max_lines` theo loại cho `design-tokens`), và chỉ áp `line-threshold` cho file có metadata dk.

### M3. Theme không giải alias theo theme

- `internal/tokens/css.go:29` (chuỗi mồi là `path@name`) và `internal/tokens/tokens.go:112` (giải `target.Value`, tức giá trị `:root`).
- Vấn đề: nếu giá trị theme của một token semantic trỏ đến một token semantic khác **cũng có** theme override, khối `[data-theme]` lấy giá trị sáng của token đó → CSS theme sai âm thầm. Ngoài ra alias theme tự trỏ về chính token (`theme.dark = "{c.a}"`) không bị coi là vòng, trả về giá trị sáng (đã tái tạo). Khung template hiện chỉ trỏ primitive nên chưa lộ, nhưng skill hướng agent tự viết theme.
- Cách sửa: khi giải giá trị theme, ưu tiên `Theme[name]` của token đích rồi mới đến `Value`; đưa `path@name` vào chuỗi kiểm vòng cùng dạng với `path`.

### M4. `internal/tokens` bỏ token âm thầm và không kiểm tên biến CSS

- `internal/tokens/tokens.go:49-81`, `internal/tokens/css.go:27`.
- Vấn đề (đã tái tạo bằng probe test): (a) nhóm rỗng `{}` hoặc token gõ sai khóa (`$vaule`) cho ra 0 biến, **không** lỗi, không cảnh báo — token biến mất im lặng; (b) khóa có khoảng trắng hoặc dấu nháy sinh ra CSS không hợp lệ `--a b-c"d: #111;`; (c) hai đường dẫn khác nhau (`a.b-c` và `a-b.c`) sinh cùng tên `--a-b-c`, ghi hai lần, cái sau thắng, không báo.
- Cách sửa: lỗi khi một nhóm không chứa token nào; kiểm khóa theo `^[A-Za-z0-9][A-Za-z0-9-]*$`; phát hiện trùng tên biến CSS và báo lỗi.

### M5. `$dk` trùng khóa cho ra JSON có hai `$dk`

- `internal/frontmatter/kinds.go:81-136`, `kinds.go:161-191` (`findJSONKey` dừng ở lần khớp đầu).
- Vấn đề (đã tái tạo): với `{"$dk": {"a":1}, "$dk": {"b":2}, "c": 1}`, `SplitJSONKey` chỉ bỏ khối đầu, `JoinJSONKey` chèn lại → file có hai khóa `$dk` (JSON parser lấy khóa cuối, tức metadata dk bị mất) và phần thân bị dồn một dòng. Với `$dk` không phải object (`"$dk": 5`) hoặc `$dk` lồng trong nhóm con thì `ok=false`, file rơi vào H1 (không rule nào báo).
- Cách sửa: từ chối (lỗi rõ) khi thấy `$dk` xuất hiện nhiều hơn một lần hoặc không phải object, thay vì im lặng.

### M6. `hook run` chưa biết `.html`/`.json`, lớp chặn "chỉ tạo bằng dk new" hở với họ Design

- `internal/hook/run.go:64` (`!strings.EqualFold(filepath.Ext(rel), ".md")` → im lặng).
- Vấn đề: cả hai skill khẳng định "Mọi file trong `docs/` chỉ được tạo bằng `dk new`", và phase 3 dựa vào hook `PreToolUse` để chặn tạo tay. Với `.html`/`.json` hook không chặn tạo mới và cũng không ghi dòng changelog `post-edit`. Kết hợp H1, agent hoàn toàn có thể tự ghi một mockup không metadata mà không lệnh nào báo.
- Cách sửa: mở rộng bộ đuôi của `hook run` theo `doctype` (`.md`, `.html`, `.json`), giữ nguyên các điều kiện im lặng khác; hoặc ghi rõ giới hạn này vào rules.md.

### M7. `mockup-tokens` chỉ bảo đảm "không hex, không px", chữ trong rules.md hứa nhiều hơn

- `internal/check/rules_mockup.go:14-15`; `assets/skills/doc-design-system/references/rules.md:11-14` ("không có giá trị gõ tay"), `doc-design-flow/references/rules.md:41-43`.
- Vấn đề (đã tái tạo): `rgb(37,99,235)`, `red`, `1.5rem`, `50%`, `200ms` đều không bị bắt. Agent bị chặn `#2563eb` sẽ rất dễ chuyển sang `rgb(...)` và vẫn qua lint.
- Cách sửa: thêm `rgb(`, `rgba(`, `hsl(`, `hsla(` và các đơn vị `rem|em|%` vào bộ khớp (hoặc chốt lại phạm vi), và sửa chữ trong hai `rules.md` thành đúng những gì rule bắt.

### M8. Thiếu test cho hai điểm tích hợp mới nhiều rủi ro nhất

- `internal/cli/tokens.go` (không có test nào trong `internal/cli` gọi `tokens css`), `internal/changelog/record.go:20-42` và `:57-71` (không test `Tracks` với `GeneratedMarker`, không test `Record` với `.html`/`.json`).
- Vấn đề: đường `--in/--out`, `requireProject`, `resolve`, và việc bump `updated` trong `$dk`/chú thích HTML rồi giữ JSON hợp lệ chỉ được chứng minh bằng chạy tay (tôi kiểm lại: đúng, JSON sau bump vẫn `json.tool` qua, `changelog pending` bỏ `tokens.css`). Không có test thì lần refactor sau không ai bắt được.
- Cách sửa: một test `cli` cho `tokens css` (mặc định, `--in/--out`, ngoài dự án) và một test `changelog` cho `Record` trên `.json` và `.html` cộng `Tracks` với dòng `GeneratedMarker`.

### M9. Metadata sai làm tài liệu vô hình với `dk check`

- `internal/check/check.go:93`.
- Vấn đề: đây là gốc chung của H1 và M5, và cũng đúng với `.md` (đã có từ phase 1), nhưng với họ Design thì hậu quả nặng hơn vì metadata nằm trong chú thích HTML hoặc khóa JSON, rất dễ gõ sai và không có dấu hiệu nào. `tokens.json` hỏng cú pháp cũng không lệnh nào báo ngoài `dk tokens css`.
- Cách sửa: cùng H1; ngoài ra `dk check` nên báo lỗi khi `.json` trong `docs/` không parse được.

## Low

- **L1. `tmpl.Raw` dùng glob.** `internal/tmpl/tmpl.go:37-43`: `templates/<loại>.*`, `len(matches) != 1` đều trả "không có template cho loại %q". Một file lạc (`mockup.html.bak`, hoặc thêm `mockup.md`) làm `dk new` chết với thông điệp sai, dù `doctype.validate` (`internal/doctype/doctype.go:93`) vẫn qua vì nó `Stat` đúng `t.Ext()`. Sửa: truyền `Type.Ext()` vào `Raw` thay vì glob; không có rủi ro trùng tên loại nào hiện tại (không loại nào là tiền tố của loại khác trước dấu chấm).
- **L2. Finding trùng.** `internal/check/rules_mockup.go:28`: `styleAttrRe` cũng khớp bên trong `<style>` (ví dụ selector `a[style="color:#333"]`), nên cùng một giá trị bị báo hai lần cùng dòng. Sửa: bỏ các span lồng nhau trước khi quét.
- **L3. False positive nhỏ.** `#abc`, `#deadbeef` (id selector) và `0px` bị báo lỗi (đã tái tạo). Sửa: bỏ qua `#...` khi ký tự trước là đầu selector (`{`, `}`, `;`, đầu dòng, dấu cách sau `,`), và bỏ qua `0px`.
- **L4. Trích dẫn plan sai.** `doc-design-system/references/rules.md:128` và `doc-design-flow/references/rules.md:118` trỏ "mục 1.1, 1.3, 1.5, 1.10, 1.13" — plan chỉ có tới 1.12. (Việc rules.md trỏ đường dẫn nội bộ repo `dk` là quy ước có từ phase 4, dùng cho cả 7 skill; không phải hồi quy phase 7, nhưng nên tính vào phase 10 vì bản cài nằm trong dự án khác.)
- **L5. Lệnh mẫu không bọc nháy.** `doc-design-flow/SKILL.md:96`: `--set external=<url>`; URL có `?` hoặc `&` làm shell lỗi (`no matches found` trong zsh, đã tái tạo). Sửa: `--set 'external=<url>'`.
- **L6. `state` không vào tên file.** `doc-design-flow/SKILL.md:82` tạo mockup đầu tiên với `--set state=empty` nhưng tên vẫn `F-001-B1.html`, trong khi `rules.md:54` dành tên đó cho `normal` (đã tái tạo). Xem H4 để sửa gộp.
- **L7. `userflow` thiếu `from.brief`.** `assets/types.toml:142-143` chỉ có `from.feature-spec` và `from.cr`, còn `doc-design-flow/SKILL.md:58` nói "`--from <brief hoặc CR>`". Lệnh vẫn chạy (nhờ `--set feature=`) nhưng `title` không được chép từ brief. Sửa: thêm `from.brief = { title = "title" }` hoặc bỏ chữ "brief" trong SKILL.md.
- **L8. `resolve` trùng logic sẵn có.** `internal/cli/tokens.go:64-72` thêm helper `(*app).resolve` chỉ dùng ở một chỗ, trong khi `internal/cli/check.go:64` và `internal/cli/render.go:62` vẫn nội tuyến cùng logic. Sửa: đưa helper về `root.go` và dùng lại ở ba chỗ. Ngoài ra `--out` nhận đường dẫn bất kỳ ngoài `docs/`, không kiểm.
- **L9. `go.mod` chưa tidy.** `goldmark` vẫn nằm ở khối `// indirect` dù `internal/render` dùng trực tiếp; `go mod tidy` chuyển nó lên khối trực tiếp (đã kiểm trên bản sao). Không thêm dependency mới. Nợ từ phase 2.
- **L10. Chỉ mục HTML không có mockup.** `internal/cli/render.go:80` lọc Markdown trước khi dựng `index.html`, nên `docs/html/index.html` không có đường vào mockup và tokens; đúng chủ ý (mockup mở trực tiếp) nhưng không ghi ở đâu cho người đọc.
- **L11. Mô tả skill lệch chữ phase.** `phase-07` bước 6, 7 quy định mô tả cụ thể; bản ship (`SKILL.md:3` của cả hai skill) đổi chữ để qua test không trùng 3 từ liên tiếp (đã ghi trong `phase-07-run.md`). Chấp nhận, nêu để plan cập nhật quyết định.
- **L12. `dk template list` in `design/{layer}s`** với placeholder chưa giải (`assets/types.toml:120`), người đọc dễ tưởng là tên thư mục thật.

## Xác nhận không hồi quy (mục b của yêu cầu)

Đã chạy thật trên dự án mẫu, không thấy hồi quy ngoài các mục đã nêu:

| Consumer | Kết quả với `.html` / `.json` |
|---|---|
| `dk render --all --index` | chỉ render `.md`, `docs/html/design/` không sinh (đúng, `render.go:80`) |
| `dk render <file .json>` | mã 2 "render chỉ nhận file .md" (đúng) |
| `dk refs` | `tokens.json` có "Đi (1) → brief.md" qua `source`; `refs.LinkRe` không khớp `href=` nên HTML không sinh liên kết rác |
| `dk index` | 4 thư mục cố định không gồm `design/`, không ảnh hưởng |
| `dk status` | đếm `design-tokens`, `mockup` theo trạng thái |
| `dk check` | rule mới chạy đúng; `link-broken` áp cho thân mockup nhưng chỉ khớp liên kết Markdown; `line-threshold` là vấn đề M2 |
| `dk changelog add` `.json` | bump `updated` trong `$dk`, file vẫn `json.tool` qua, dòng `mới, 123 dòng` |
| `dk changelog add` `.html` | bump `updated` trong chú thích, thứ tự khóa giữ nguyên |
| `dk changelog pending` | bỏ `tokens.css` nhờ `GeneratedMarker` (`record.go:22`); **không** thấy file bị xóa (liên quan H5 vấn đề 2) |
| `hook run` | vẫn chỉ `.md` → M6 |

`frontmatter.SplitJSONKey` với JSON lạ: `$dk` không phải object → `ok=false` (im lặng, xem M5/H1); `$dk` lồng trong nhóm con → `ok=false` (đúng, chỉ nhận cấp cao nhất); chuỗi chứa "$dk" ở giá trị khác → tách đúng (dùng `json.Decoder`, không tìm chuỗi thô); trùng khóa → M5. `JoinJSONKey` với body `{}` → `{\n  "$dk": {...}\n}` hợp lệ; body còn khóa khác → chèn `,` đúng; bỏ `$dk` ở cuối object thì cắt luôn dấu phẩy trước đó (`kinds.go:102-108`), JSON còn hợp lệ (có test).

`internal/tokens` (mục f): alias lồng trong chuỗi đúng (`0 1px 2px {c.s}` → `0 1px 2px #111827`); `$type` kế thừa và ghi đè ở nhóm con đúng; `$extensions.dk.theme` sinh `[data-theme]` theo tên đã sắp; số cho `dimension` thành `px`, `{value,unit}` thành `4rem`; alias vòng, alias lạ, giá trị array/null báo lỗi rõ và không ghi file. Còn lại: M3, M4.

## Tiêu chí phase 7

| Tiêu chí | Đánh giá |
|---|---|
| `dk tokens css` đúng, alias vòng báo rõ | Đạt |
| Mockup dùng `tokens.css`; `check` bắt hex, px gõ tay | **Một phần**: bắt được ca mẫu, nhưng H1, H2 là đường vòng và H3 là chặn sai |
| User flow đúng mã bước spec, lệch bị bắt | Đạt về mã, nhưng H4 làm quy tắc skill và rule đá nhau |
| Bộ tối thiểu đúng thứ tự lớp, mỗi component đủ 5 mục | Template đủ mục; phần AI thật để phase 10 (chấp nhận theo tiền lệ 1.11, 1.12) |
| AI sinh mockup từ template, `check` qua | Khung template qua; H5 làm lệnh mẫu của bước này lỗi |
| `--set external=<url>` giữ liên kết và ảnh | Đạt (đã tái tạo) |

## Hành động đề xuất theo thứ tự

1. H1 và H2: đóng hai đường vòng của lint (rule báo file thiếu metadata; nhận nháy đơn).
2. H3: bỏ `@media`/`@container` khỏi vùng quét `px`, ghi giới hạn vào `rules.md`.
3. H4: chốt cách đặt tên mockup theo trạng thái, sửa `rules.md` (và `types.toml` nếu chọn `{state}`).
4. H5: sửa lệnh mẫu bước 6 và bỏ yêu cầu `changelog add` cho file đã xóa.
5. M1: kiểm giá trị `--set` trước khi ghép đường dẫn.
6. M2, M3, M4, M5: ngưỡng dòng cho file không Markdown; theme giải theo theme; báo lỗi thay vì bỏ token/trùng tên/trùng `$dk`.
7. M6, M7, M8: mở rộng `hook run`, sửa chữ hứa trong `rules.md`, thêm hai test còn thiếu.
8. Low: gom vào một lượt (L1, L4, L5, L7, L8, L9).

Không sửa mã trong lượt review này; không chỉnh trạng thái phase (để lead hoặc planner cập nhật `plan.md` mục 1.13 và bảng phase).

## Câu hỏi chưa giải quyết

- Mockup nhiều trạng thái: chốt `{state}` trong tên file (đổi `types.toml`) hay giữ hậu tố kiểu bước? Ảnh hưởng H4, L6 và cách họ Feature Spec liên kết mục 6.
- `mockup-tokens` nên mở rộng tới `rgb()`/`rem`/`%` (siết hơn, rủi ro false positive cao hơn) hay chốt phạm vi "chỉ hex và px" và sửa chữ trong `rules.md`?
- `tokens.json` vượt `max_lines` 800: miễn ngưỡng cho file không Markdown, hay cho `design-tokens` một `max_lines` riêng (đụng quyết định 1.8 "không tắt được")?
