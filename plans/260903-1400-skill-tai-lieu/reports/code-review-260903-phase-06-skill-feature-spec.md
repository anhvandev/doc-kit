# Review phase 6: skill doc-feature-spec

Ngày: 2026-09-03. Phạm vi: các file đổi ở phase 6 (template feature-spec, `internal/tmpl`, `internal/docs/new.go`, `internal/frontmatter.SetBool`, `internal/render/steps.go`, hai rule mới trong `internal/check`, skill `doc-feature-spec`, README, types.toml). Chỉ đọc mã và chạy `go test ./... && go vet ./...` một lần: **tất cả pass, vet sạch**.

Kết luận: hai lỗi mức High làm `dk check` báo lỗi giả ngay trên file `dk new` vừa sinh ra hoặc trên đường đi mà SKILL.md dạy. Cả hai lọt lưới vì test và bản chạy thử chỉ thử `crud` kèm `has_ui=false`.

## High

**H1. `dk new feature-spec --set format=crud` (không kèm `has_ui=false`) sinh file bị `step-codes` báo lỗi ngay.**
`assets/templates/feature-spec.md:103-111` — khối `{{if .HasUI}}` mục 6 độc lập với `format`, nên bản `crud` vẫn có bảng "Giao diện" với ô cột đầu `B1`. `internal/render/steps.go:36-56` — bản `crud` không có tiêu đề chứa "hành vi" nên `findBehaviorTable` rơi xuống nhánh `byCode` và chọn **bảng mục 6**; `StepCodes` trả `table=[B1]`, `diagram=[]`. `internal/check/rules_steps.go:21` chỉ bỏ qua khi cả hai rỗng, nên `internal/check/rules_steps.go:34` báo `mã có trong bảng hành vi nhưng không có trong sơ đồ: B1`.
Nguyên nhân: mục 6 mẫu có sẵn mã bước, còn `crud` (theo plan 1.12 là "không mã bước nên step-codes bỏ qua") không có bảng hành vi để giành quyền.
Sửa (một trong hai): bỏ dòng mẫu `| B1 | ... |` ở mục 6 của template (để bảng rỗng), hoặc cho `stepCodes` bỏ qua spec có `format: crud`. Thêm case test `{"crud", true, ...}` chạy qua `check`, không chỉ qua `tmpl.Render`.

**H2. Đường đi "luồng tuyến tính dưới bốn bước: bỏ sơ đồ" mà SKILL.md và rules.md dạy sẽ làm `dk check` fail.**
`assets/skills/doc-feature-spec/SKILL.md:65-66` ("giữ `spec` nhưng xóa khối mermaid ở mục 4") và `references/rules.md:44-45`, `90` ("bỏ sơ đồ, bảng hành vi là đủ"). Sau khi xóa mermaid, `diagram=[]` còn bảng mục 5 vẫn `B1..B5` → `internal/check/rules_steps.go:34` báo lỗi, trái với chỉ dẫn "Không báo người khi `check` còn lỗi" ở SKILL.md:113.
Nguyên nhân: `stepCodes` coi "bảng có mã mà không có sơ đồ" luôn là lệch, không có nhánh "cố ý không vẽ sơ đồ".
Sửa: khi `len(diagram) == 0` thì bỏ qua (chỉ so khi có sơ đồ), hoặc SKILL.md yêu cầu xóa luôn cột mã bước khi bỏ sơ đồ. Cùng gốc với H1 nên nên sửa một lần.

## Medium

**M1. `--set has_ui=<giá trị khác "false">` im lặng đảo ý người dùng.**
`internal/docs/new.go:107` (`o.Set["has_ui"] != "false"`) và `:148-154`. `--set has_ui=False`/`0`/`no` cho `HasUI=true`, ghi frontmatter thành chuỗi `"False"`, mục 6 vẫn render, `check` cũng coi là true (`rules_spec.go:21` so đúng chuỗi `"false"`). Không có lỗi nào báo cho người dùng, khác hẳn cách `format` được validate ở `:115-117`.
Sửa: validate `has_ui` chỉ nhận `true`/`false`, sai thì lỗi như `format`.

**M2. `findBehaviorTable` chọn bảng theo vị trí cho biến thể `use-case`, dễ chọn nhầm.**
`internal/render/steps.go:36-56`. Với `use-case`, tiêu đề mục 5 là "Use Case Specification" (không chứa "hành vi"), nên bảng hành vi được chọn bằng "bảng đầu tiên có cột đầu khớp `^B\d+[a-z]?$`". Hiện đúng vì bảng mục 5 đứng trước bảng mục 6, nhưng hỏng ngay khi tác giả để bảng mục 5 trống mã (chưa điền) hoặc thêm bảng nào có ô `B1` ở mục 3: khi đó bảng "Giao diện" thành bảng hành vi, `check` báo lệch và HTML gắn `id="step-*"` sai chỗ (`internal/render/mermaid.go:64-81`).
Sửa: cho `findBehaviorTable` nhận thêm tiêu đề "use case" (hoặc "luồng chính"), hoặc ưu tiên bảng nằm dưới tiêu đề `## 5.`.

**M3. `sectionNumbers` đếm fence quá đơn giản.**
`internal/check/rules_spec.go:46-49`. Chỉ toggle trên tiền tố ```` ``` ````: fence `~~~`, fence 4 backtick lồng nhau, và khối mã thụt 4 khoảng trắng đều không được xử lý; một fence không đóng làm toàn bộ phần còn lại của file bị bỏ, cho ra thông điệp "thiếu mục 7, 8, 9, 10, 11" gây hiểu nhầm.
Sửa: dùng luôn AST của goldmark (đã có sẵn trong `render`) để lấy heading, hoặc ít nhất khớp đúng độ dài fence khi đóng.

## Low

- **L1. `--set k=true|false` đổi kiểu cho *mọi* khóa** — `internal/docs/new.go:148-152`. Rà `assets/types.toml`: không khóa nào (title, owner, source, status, level, kind, priority, requester...) cần chuỗi `"true"`/`"false"`, nên hiện không có hồi quy thật; rủi ro còn lại là `--set title=false` sinh `title: false` rồi `frontmatter-required` vẫn cho qua (GetString trả `"false"`). Nếu muốn chặt, chỉ ép bool cho khóa khai báo kiểu bool trong types.toml.
- **L2. `updated` thiếu hoặc sai dạng chỉ là warning và bỏ luôn việc kiểm tra thứ tự** — `internal/check/rules_cr_order.go:154-162`; `updated` không nằm trong `required` của `feature-spec` (`assets/types.toml`), nên spec viết tay thiếu `updated` vừa ồn vừa không được kiểm. Cân nhắc thêm `updated` vào `required`.
- **L3. Giới hạn của `cr-approval-order` (đã ghi ở plan 1.12) vẫn nên nêu trong SKILL.md**: chạm vào CR sau khi sửa spec, hoặc sửa cùng phút, làm rule im lặng. SKILL.md:136 nói như thể rule luôn bắt được.
- **L4. `--source` của changelog không thống nhất**: SKILL.md:110 nói `<yymmdd>-<slug>`, còn `source` trong frontmatter do `readSource` sinh là `<yymmdd>-<slug>/brief.md` (`internal/docs/new.go:190`) và rules.md:108 ghi theo dạng thứ hai. Chọn một dạng.
- **L5. `format` chỉ được validate khi `typeName == "feature-spec"`** (`internal/docs/new.go:115`); `--set format=xyz` cho loại khác ghi thẳng vào frontmatter. Không hại hiện tại.
- **L6. README.md:23 thiếu dấu phẩy** giữa cụm `doc-overview` và `doc-adr`.

## Không phải lỗi (đã kiểm)

- **`\b(B\d+[a-z]?)\b` không rò mã giả từ mermaid**: `stateDiagram-v2` mẫu chỉ có `Nhap`, `ChoDuyet`, `DaDuyet`, `[*]`; `sequenceDiagram` mẫu chỉ có `B1..B4` trong nội dung thông điệp; `flowchart` mẫu trùng id và nhãn nhưng `stepCodesIn` khử trùng. `B3aa` không khớp (đúng ý plan 1.12: chỉ một chữ hậu tố).
- **`findBehaviorTable` cho `spec`, `story`, `state`**: tiêu đề mục 5 chứa "hành vi" nên nhánh `byHeading` thắng bảng mục 6 đứng sau; `### Vòng đời trạng thái` của `state` chỉ reset cờ tạm, không ảnh hưởng.
- **`spec-section-order` với `has_ui` vắng mặt**: `GetString` trả `""` ≠ `"false"` → vẫn đòi mục 6, đúng mặc định `has_ui: true` của plan 1.12; spec tạo trước phase 6 vẫn qua.
- **`## 12.`, mục lặp, sai thứ tự**: `compareSections` phân biệt đủ bốn trường hợp; fixture F-006 chứng minh (thiếu 9, lạ 12, lặp 8).
- **`### Lịch sử nguồn`** là tiêu đề cấp 3 nên không lọt vào `sectionRe` — đúng thiết kế `update`.
- **Múi giờ**: `parseUpdated` thử dạng có offset và RFC3339 trước khi rơi về `ParseInLocation(..., time.Local)`; khớp quyết định 1.12 (không đổi định dạng `updated`).
- **CR `rejected` không nằm trong `final`** (`assets/types.toml`) nên spec trỏ CR bị từ chối vẫn bị bắt — đúng.
- **Hợp đồng công khai**: `tmpl.Data` chỉ thêm trường (một caller duy nhất là `docs.New`, `missingkey=error` không bị ảnh hưởng); `frontmatter.SetBool` và `docs.SpecFormats` là bổ sung; nới `^B\d+$` thành `^B\d+[a-z]?$` chỉ mở rộng tập nhận, không phá tài liệu cũ. Không thấy thay đổi phá vỡ nào khác.
- **Đếm skill 4→5** khớp `assets/skills/` (5 thư mục), README và `dk skill status` 10 dòng.
- **Lệnh và cờ SKILL.md dùng đều có thật**: `dk --version` (`root.go:85`), `dk new --from/--set` (`new.go:30-32`), `dk changelog add --summary --source` (`changelog.go:28-29`), `dk render <file>`, `dk index features|intake`, `dk check <file>`.
- **README mục quy tắc**: 11 rule liệt kê khớp đúng bảng đăng ký ở `internal/check/check.go:44-57`.

## Câu hỏi còn treo

1. `crud` có `has_ui: true` là tổ hợp hợp lệ (CRUD có màn hình) hay skill nên luôn ép `has_ui=false` cho `crud`? Câu trả lời quyết định cách sửa H1.
2. Bỏ sơ đồ mà vẫn giữ mã bước ở bảng (H2) có phải trạng thái muốn hỗ trợ lâu dài không? Nếu có thì `step-codes` nên bỏ qua khi không có mermaid.
