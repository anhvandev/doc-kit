# Code review phase 8: skill tiến độ và test (doc-plan-report, doc-test)

- Ngày: 2026-09-04
- Phạm vi: `assets/types.toml` (11 loại mới), 11 template, `assets/skills/doc-plan-report/`,
  `assets/skills/doc-test/`, `internal/frontmatter/kinds.go`, `internal/doctype/doctype.go`,
  `internal/config/config.go`, `internal/docs/{new.go,scan.go,from_spec.go}`,
  `internal/tmpl/tmpl.go`, `internal/cli/new.go`, `internal/check/{check.go,rules_report.go,rules_backlink.go}`,
  fixture và test đi kèm, `README.md`.
- Cách kiểm: đọc mã; `gofmt -l .` rỗng, `go vet ./...` sạch, `go test ./...` toàn bộ ok;
  `go build -o /tmp/dkr ./cmd/dk` rồi chạy thật trong repo git tạm (`/tmp/dkt`): tạo đủ
  11 loại mới, `dk check`, `dk render --all`, `dk changelog add`, `dk status`, `dk index`;
  thêm một test tạm trong `internal/frontmatter` để thử CRLF và thân rỗng của `.feature`
  (đã xóa sau khi chạy).
- Kết quả cổng: sạch. Không có lỗi build, vet, gofmt, test.

## Tiêu chí chấp nhận

| # | Tiêu chí | Kết quả |
|---|---|---|
| 1 | `dk new plan` ra `plans/{yymmdd}-{hhmm}-{slug}/plan.md`; `plan-phase --in` đánh số trong thư mục | đạt (`phase-01`, `phase-02`, `phase-03` liên tiếp) |
| 2 | `test-case --from` một Scenario mỗi AC, tag `@F-xxx @ACn`; AC lệch khung để lại dấu | đạt cho AC dạng gạch đầu dòng và khối gherkin; **không đạt** cho AC dạng danh sách đánh số hoặc bảng (H1) |
| 3 | Dry-run Gherkin tùy chọn qua `[test] bdd_cmd`; không có thì báo "chưa kiểm chạy được" | đạt (`config.Test.BddCmd`, `doc-test/SKILL.md` bước 4 và `rules.md`) |
| 4 | Report thiếu bằng chứng bị cảnh báo `report-evidence` | đạt (chú thích HTML và khối ``` rỗng không tính) |
| 5 | `decision-log --append` nối dòng không phá nội dung cũ | đạt (bump `updated`, giữ nguyên thân) |
| 6 | `spec-has-test` không cảnh báo khi spec đã có `.feature` | đạt (chỉ `.feature`/`test-case-table`/loại lạ được tính; checklist UI và test report không tính) |

Không thấy hồi quy: `docs.New` cho loại cũ (subdir theo nguồn, `beside_source`,
`--from`, `id_prefix` cho loại `seq` trong `docs/`) giữ nguyên; `render --all` bỏ
`.feature`; `changelog add` nhận `.feature` và vẫn từ chối file trong `plans/`;
`hook run` vẫn chỉ chặn `.md`; `link-broken`, `line-threshold`, `backlink`, `index`,
`refs`, `status` không đổi hành vi khi tập quét thêm `.feature`. Hợp đồng công khai
giữ tương thích: `Options` chỉ thêm trường, khóa `types.toml` chỉ thêm, `frontmatter`
chỉ thêm hàm, mã thoát 0/1/2/3 không đổi.

## High

### H1. `from_spec` bỏ im lặng toàn bộ AC khi mục 9 không dùng gạch đầu dòng

`internal/docs/from_spec.go:17` (`acBulletRe`), dùng ở `:102`.

Nguyên nhân: chỉ nhận `^\s*[-*]\s*...AC\d+`. Mục 9 viết bằng danh sách đánh số
(`1. AC1. ...`) hoặc bằng bảng `| AC3 | ... |` không khớp, `scenarios()` trả rỗng,
template rơi vào nhánh `{{else}}` và sinh đúng một `Scenario: AC1` trống. CLI chỉ in
"Đã tạo ...", không có tín hiệu nào. Tái hiện thật:

```
$ dkr new test-case bang-ac --from docs/features/F-002-bang-ac.md --set owner=a
Đã tạo docs/test/F-002-cases.feature   # spec có AC1, AC2 (đánh số) và AC3 (bảng)
# file sinh ra: đúng 1 Scenario rỗng
```

Đây là kiểu hỏng nguy hiểm nhất của tính năng này: file test trông hoàn chỉnh nhưng
mất hết tiêu chí, trái với ràng buộc "báo rõ AC nào không parse được thay vì bỏ qua
im lặng" ghi trong phần Rủi ro của phase.

Sửa: trong `ExtractSpec`, quét mục 9 tìm mọi dòng chứa `AC\d+` (kể cả `\d+\.` và
dòng bảng); dòng nào không sinh được `Scenario` thì thêm `tmpl.Scenario{Code, Raw}`
để rơi vào nhánh `# chưa tách được`. Đồng thời cho `docs.New` trả về số Scenario và
số dòng chưa tách, `internal/cli/new.go:runNew` in ra: "N Scenario, M dòng AC chưa
tách được".

## Medium

### M1. `--in` không bị ràng buộc trong `plans_dir`

`internal/docs/new.go:114` (nhánh `{in}`), `internal/cli/new.go:81` (`a.resolve`).

Nguyên nhân: `--in` được `resolve` thành đường dẫn tuyệt đối rồi dùng thẳng làm gốc,
không so với `PlansDir` và không đòi thư mục có sẵn. Hai hệ quả thật:

```
$ dkr new plan-phase evil --in /tmp/outside-dk   # ngoài repo, dk check không thấy
Đã tạo ../outside-dk/phase-01-evil.md (id phase-01)
$ dkr new report typo --in plans/khong-ton-tai   # gõ sai tên thư mục plan
Đã tạo plans/khong-ton-tai/reports/report-260904-0859-typo.md
```

Trái với cách phase 7 siết `feature`, `step`, `layer` bằng `placeholderValRe` để giá
trị không thoát khỏi thư mục loại.

Sửa: trong `docs.New`, sau khi có `o.In`, yêu cầu `o.In == o.PlansDir || inside(o.PlansDir, o.In)`
và yêu cầu `os.Stat(o.In)` là thư mục có sẵn (chỉ `<in>/reports` được tạo mới); lỗi
thì trả thông điệp nêu `plans_dir` hiện hành.

### M2. `doc-test/SKILL.md` bảo đặt số phiên bản làm slug, `dk` từ chối

`assets/skills/doc-test/SKILL.md`, mục "`update`: Test report":
`dk new test-report <phien-ban>`.

Nguyên nhân: `slugRe` chỉ nhận `a-z0-9-`, số phiên bản có dấu chấm.

```
$ dkr new test-report 1.2.0 --set owner=a
slug "1.2.0" không hợp lệ: chỉ a-z, 0-9 và dấu gạch ngang   # mã thoát 2
```

Sửa: đổi ví dụ trong SKILL.md thành `dk new test-report v1-2-0 --set title="Đợt 1.2.0"`
và ghi rõ "slug không nhận dấu chấm".

### M3. `.feature` mất metadata mà không quy tắc nào báo

`internal/frontmatter/kinds.go:88` (`SplitFeatureComment`), `internal/check/check.go:47` (bộ `Rules`).

Nguyên nhân: khối `# dk:` kết thúc ở dòng đầu tiên không bắt đầu bằng `#`. Một chú
thích Gherkin tự do đặt ngay sau khối (không có dòng trống) bị hút vào YAML; có dấu
hai chấm là YAML hỏng, `SplitFile` trả `ok=false`, cả file thành thân, `HasFM=false`.
Khi đó `frontmatter-required` không chạy (chỉ chạy trên `typed()`), không có quy tắc
nào báo file `.feature` thiếu khối `# dk:`, và triệu chứng duy nhất là `spec-has-test`
lại cảnh báo trên Feature Spec — sai chỗ. Quyết định 1.13 đã đặt luật ngược lại cho
HTML (`.html` trong `docs/` thiếu `<!-- dk: -->` là lỗi `mockup-tokens`).

Sửa: thêm quy tắc (hoặc mở rộng `frontmatterRequired`) báo lỗi cho mọi `.feature`
trong `docs/` mà `HasFM == false`, thông điệp nêu "thiếu hoặc hỏng khối `# dk:`".

### M4. `Scenario Outline` bị hạ cấp thành `Scenario`, mất bảng `Examples`

`internal/docs/from_spec.go:18` (`acScenarioRe`), `assets/templates/test-case.feature`.

Nguyên nhân: regexp nhận `Scenario Outline:` nhưng không ghi lại từ khóa, và
`scenarios()` không đọc `Examples:`. Spec dùng Outline sẽ sinh `Scenario:` với tham
số `<...>` không có bảng — cú pháp qua được parser nhưng chạy sẽ hỏng, đúng thứ mà
tiêu chí 3 (dry-run) muốn tránh.

Sửa: thêm `Outline bool` và `Examples []string` vào `tmpl.Scenario`, template in
`Scenario Outline:` cùng khối `Examples:`; hoặc tối thiểu đưa Outline vào `Raw` với
dấu `# chưa tách được` thay vì hạ cấp im lặng.

### M5. `bullets()` cắt ở dấu hai chấm đầu tiên, mất chữ của Background

`internal/docs/from_spec.go:160`.

Nguyên nhân: `strings.Cut(t, ":")` được dùng để bỏ nhãn "Tác nhân:", nhưng áp cho mọi
dòng. Bullet mục 3 không có nhãn mà có dấu hai chấm, ví dụ
`- Đã bật cờ tính năng: beta`, thành `Given beta`.

Sửa: chỉ cắt khi phần trước dấu hai chấm khớp danh sách nhãn đã biết
(`Tác nhân`, `Điều kiện tiên quyết`, `Quyền`, `Dữ liệu`) và không chứa khoảng trắng
thừa; còn lại giữ nguyên dòng.

### M6. `report-evidence` coi mọi liên kết Markdown là bằng chứng

`internal/check/rules_report.go:10` (`mdLinkRe`).

Nguyên nhân: `\]\([^)\s]+\)` khớp bất kỳ liên kết nào. Một report chỉ có
`- Xem [plan](../plan.md)` đã qua quy tắc, dù không có commit, output hay kết quả
test. Mục 2 của template report lại chính là nơi khuyến khích liên kết file đã đổi,
nên hầu hết report thật sẽ qua mà không có bằng chứng.

Sửa: tách hai điều kiện — liên kết được tính là bằng chứng chỉ khi đích không phải
`.md` (log, ảnh, `.txt`, `.xml`, thư mục output) hoặc nằm ngoài `docs/`; giữ nguyên
`commitRe` và `fenceRe`.

## Low

- **L1** `assets/templates/test-case-table.md`: ô bảng chép thẳng chữ AC, một dấu `|`
  trong AC sẽ vỡ bảng Markdown. Sửa: escape `|` thành `\|` khi render (hàm template
  `escapeCell`).
- **L2** `assets/templates/test-case.feature`: nhánh `Raw` in `Given `, `When `, `Then `
  kèm khoảng trắng cuối dòng và bước rỗng. Sửa: in `Given TODO` (hoặc bỏ hẳn ba dòng,
  chỉ để lại chú thích) và cắt khoảng trắng cuối.
- **L3** `assets/templates/test-case.feature`: `Background` in nhiều dòng `Given` liên
  tiếp; đúng Gherkin nhưng quy ước là `Given` rồi `And`. Sửa: dùng `{{if eq $i 0}}Given{{else}}And{{end}}`.
- **L4** `internal/cli/status.go:33` chỉ `Scan(a.root, a.cfg.DocsDir)`, nên `plan`,
  `plan-phase`, `report` không bao giờ xuất hiện trong `dk status` dù `check` có quét
  `plans/`. Cùng file `:82` dùng `%-14s`, `ui-test-checklist` (17 ký tự) làm lệch cột.
- **L5** `assets/templates/changelog-product.md`: chú thích ghi "Vượt 1000 dòng thì cắt"
  trong khi `line-threshold` báo lỗi ở `max_lines = 800`. Sửa chú thích về 800.
- **L6** `internal/docs/new.go:246` (`appendLine`) + `internal/cli/new.go:100`:
  `--append` khi chưa có file thì tạo file với `owner` rỗng (sinh lỗi
  `frontmatter-required` ngay sau đó) nhưng chỉ in "Đã nối dòng", giấu việc đã tạo file.
  Sửa: in "Đã tạo ... và nối dòng" khi vừa tạo.
- **L7** `internal/docs/new.go:110`: `--in` bị bỏ qua im lặng với loại không dùng `{in}`
  (ví dụ `dk new feature-spec x --in foo`). Sửa: lỗi mã 2 khi `--in` không dùng đến.
- **L8** Thiếu `--in` trả mã 1 (`fail(codeError)`) chứ không phải 2. Nhất quán với cách
  thiếu `--set feature` của phase 7, nên chỉ ghi nhận.
- **L9** `appendLine` đọc rồi ghi cả file, không khóa; hai `--append` song song mất một
  dòng. Cùng loại với M6 đã chấp nhận ở quyết định 1.9.
- **L10** `SplitFeatureComment`: `# dk:` mà dòng kế không phải `#` trả `ok=true` với
  mapping rỗng, nên file có `HasFM=true`, `Type=""`. Vô hại hiện tại nhưng che mất lỗi
  người viết. Sửa: khối rỗng thì trả `ok=false`.
- **L11** Không có test mức CLI cho `--in`, `--append`, `test-case --from`
  (`internal/cli/*_test.go` không nhắc `plan-phase`, `decision-log`, `test-case`).
  Bao phủ hiện chỉ ở `internal/docs`. Nên thêm một ca chạy chuỗi lệnh thật như các phase trước.

## Nit

- `internal/docs/new.go:174`: `ExtractSpec` chạy cho mọi loại có `from.feature-spec`,
  kể cả `userflow` (phase 7) vốn không dùng `Scenarios`/`Background`/`Steps`. Không sai,
  chỉ là việc thừa; thu hẹp bằng cách kiểm `t.Kind == "feature" || strings.HasPrefix(typeName, "test-")`.
- `internal/docs/from_spec.go:43` (`sections`) không phân biệt khối mã, nên một dòng
  `## 10.` bên trong khối ``` ở mục 9 sẽ kết thúc mục sớm. `sectionNumbers` của
  `rules_spec.go` đã xử lý fence; nên dùng chung.
- `gwtRe` nhận `Given|When|Then` ở bất kỳ đâu trong dòng, kể cả khi là chữ nội dung;
  không tránh được với AC không in đậm, chỉ ghi nhận.
- `JoinFeatureComment` với thân rỗng thêm một dòng trống, nên round-trip thân rỗng
  không đối xứng tuyệt đối. CRLF thì round-trip đúng từng byte (đã kiểm).

## Việc nên làm theo thứ tự

1. H1: không để mất AC im lặng (sửa parser + in số Scenario đã sinh).
2. M1: ràng buộc `--in` trong `plans_dir` và đòi thư mục có sẵn.
3. M3: quy tắc `check` cho `.feature` mất khối `# dk:`.
4. M4, M6: `Scenario Outline` và siết `report-evidence`.
5. M2, M5 và nhóm Low.

## Câu hỏi mở

- `dk new report` bắt buộc `--in` (kể cả `--in plans`). Phase mô tả "vào
  `plans/<dir>/reports/` hoặc `plans/reports/`". Có nên để `--in` rỗng mặc định về
  `plans/reports/` cho gọn không?
- `report-evidence` cảnh báo cả report `status: draft` (vừa tạo là đã cảnh báo). Giữ
  nguyên hay chỉ áp cho `submitted`?
- `dk status` có nên quét `plans/` để đếm plan, phase, report như `check` không?
