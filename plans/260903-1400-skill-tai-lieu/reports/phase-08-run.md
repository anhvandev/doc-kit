# Phase 8: kết quả chạy thử doc-plan-report, doc-test

Ngày: 2026-09-04. Dự án mẫu `/tmp/tmp.*` (git init, `dk init`, `dk skill install`), binary `dk` build từ `go build ./cmd/dk`. Máy không có `cucumber`, `behave`, `godog`, `cucumber-js` nên bước dry-run Gherkin chưa kiểm được (đúng nhánh "chưa kiểm chạy được" của skill).

## Chuỗi lệnh hai SKILL.md quy định, chạy tay

| Bước | Lệnh | Kết quả |
|---|---|---|
| Cài | `dk skill install` | 9 skill "đã cài"; `skill status` 18 dòng |
| Plan | `dk new plan dot-mot --set owner=an --set title="Đợt một"` | `plans/260904-0852-dot-mot/plan.md` (thư mục `{yymmdd}-{hhmm}-{slug}`) |
| 2 phase | `dk new plan-phase khoi-tao --in <plan>`; `... hoan-thien --in <plan>` | `phase-01-khoi-tao.md` (id `phase-01`), `phase-02-hoan-thien.md` (id `phase-02`), đếm trong thư mục |
| Thiếu `--in` | `dk new plan-phase x` | `loại plan-phase cần --in <thư mục plan>`, mã 1 |
| Report thiếu bằng chứng | `dk new report phase-01-run --in <plan> --set source=phase-01` rồi `dk check <report>` | `warning report-evidence: report không có bằng chứng`, mã 0 |
| Report có commit | ghi `commit <hash HEAD>` vào mục 3, `dk check` | 0 lỗi, 0 cảnh báo |
| Report có liên kết file | thay bằng `[output](...)` | `report-evidence` hết; `link-broken` bắt đường dẫn sai (kịch bản gõ sai, không phải lỗi rule) |
| Report ngoài plan | `dk new report tong-hop --in plans` | `plans/reports/report-260904-0852-tong-hop.md` |
| Decision log | `dk new decision-log --append "Nút Lưu \| an \| F-001" --set owner=an`; lần hai `--append "Bỏ cột \| an \| -"` | lần đầu tạo file và nối; lần hai "Đã nối dòng vào", 2 dòng `- 2026-09-04 \| ...`, phần cũ nguyên; `check` sạch |
| `--append` loại khác | `dk new adr x --append y` | `--append chỉ dùng cho decision-log`, mã 1 |
| Spec 4 AC | spec `approved`, mục 9 có AC1..AC3 khung Given/When/Then, AC4 "Người dùng thấy đúng" | `dk check` spec: `warning spec-has-test` |
| Test case Gherkin | `dk new test-case bo-loc --from <spec> --set owner=an` | `docs/test/F-001-cases.feature`: khối `# dk:` 11 trường (`source: F-001`, `feature: F-001`), `@F-001` trên `Feature:`, `Background` 2 dòng Given từ mục 3 (bỏ nhãn "Tác nhân:"), 4 `Scenario: ACn` tag `@F-001 @ACn`; AC4 có `# chưa tách được: Người dùng thấy đúng` và ba bước trống |
| `spec-has-test` | `dk check <spec>` sau khi có `.feature` | hết cảnh báo; chỉ còn `link-broken` do mockup kịch bản chưa tạo |
| Checklist UI | `dk new ui-test-checklist bo-loc --from <spec>` | `F-001-ui.md`: `- [ ] B1: khớp mockup [B1](../design/mockups/F-001-B1.html)`, đủ mã bước có ở mục 6 |
| Bảng | `dk new test-case-table bo-loc --from <spec>` | `F-001-cases.md`: 4 dòng `\| ACn \| ...`, AC4 "chưa tách được" |
| Loại còn lại | `roadmap`, `testing-strategy`, `changelog-product`, `test-report v1` | `docs/plan/roadmap.md`, `docs/test/strategy.md`, `docs/CHANGELOG.md`, `docs/test/reports/260904-v1.md` |
| Changelog | `dk changelog add` cho 8 file kể cả `.feature`; `dk changelog add plans/x.md` | mỗi file một dòng; file trong `plans/` bị từ chối "nằm ngoài docs/" |
| HTML | `dk render --all` | render Markdown, bỏ `.feature` |
| `status` | `dk status` | đếm đủ 9 loại mới theo trạng thái |
| Không có công cụ BDD | `which cucumber behave godog cucumber-js` | không có; skill phải báo "chưa kiểm chạy được" |

## Tiêu chí phase 8

| Tiêu chí | Kết quả |
|---|---|
| `dk new plan` và `plan-phase` đúng cấu trúc thư mục và số thứ tự | đạt |
| `test-case --from <spec>` một Scenario mỗi AC, tag đúng mã tính năng | đạt (4 AC → 4 Scenario, `@F-001 @ACn`) |
| Gherkin qua dry-run công cụ BDD khi có; không có thì báo rõ | công cụ chưa có trên máy; nhánh "có" chưa kiểm (chờ chốt công cụ, câu hỏi mở); nhánh "không" ghi trong SKILL.md bước 4 |
| Report thiếu bằng chứng bị `dk check` cảnh báo | đạt |
| Decision log nối dòng không phá nội dung cũ | đạt (test `TestDecisionLogAppend` và chạy tay) |
| `spec-has-test` không còn cảnh báo khi có `.feature` tương ứng | đạt |

## Lỗi phát hiện và đã sửa khi chạy thử

- `report-evidence` không cảnh báo report trống vì chú thích gợi ý của template chứa "commit abc1234"; sửa rule bỏ chú thích HTML trước khi so, thêm fixture.
- `decision-log --append` lần đầu tạo file thiếu `owner` nên `check` báo lỗi; SKILL.md ghi lần đầu thêm `--set owner`.
- Background chép cả nhãn "Tác nhân:"; sửa `bullets` bỏ nhãn trước dấu hai chấm. `Scenario: AC1 ` thừa khoảng trắng khi AC không có tên; sửa template.
- Thông báo `--append` in "Đã tạo" ở lần nối; đổi thành "Đã nối dòng vào".
- Test đếm: `template list` 17 → 28, skill 7 → 9, `skill status` 14 → 18 dòng.

## Sau code review

`plans/reports/code-review-260904-phase-08-skill-plan-test.md`: 1 High, 6 Medium, 11 Low, 4 nit. Đã sửa và chạy lại (`go test ./...`, `go vet`, `gofmt -l` sạch; thêm `internal/cli/phase8_test.go` chạy chuỗi lệnh thật):

- H1: mục 9 viết bằng danh sách đánh số hoặc bảng không còn bị bỏ im lặng: mọi dòng nhắc `ACn` chưa gặp thành Scenario `# chưa tách được` (mã chỉ tính một lần; chú thích HTML gợi ý của template bị bỏ trước khi quét). `dk new` in `N Scenario, M dòng AC chưa tách được`; `--json` có `scenarios`, `unparsed`.
- M1: `--in` phải là thư mục có sẵn nằm trong `plans_dir` (hoặc chính `plans_dir`); loại không dùng `{in}` mà có `--in` là lỗi (L7).
- M2: SKILL.md `doc-test` ví dụ `test-report v1-2-0 --set title="Đợt 1.2.0"`, ghi rõ slug không nhận dấu chấm.
- M3: `.feature` trong `docs/` không đọc được khối `# dk:` là lỗi `frontmatter-required` (fixture `test/hong.feature`); khối `# dk:` rỗng coi như thiếu (L10).
- M4: `Scenario Outline` giữ ở Raw kèm ghi chú "chép tay cùng bảng Examples" thay vì hạ cấp thành `Scenario`.
- M5: `bullets` chỉ bỏ nhãn đã biết (`Tác nhân`, `Điều kiện tiên quyết`, `Quyền`, `Dữ liệu`); bullet khác có dấu hai chấm giữ nguyên.
- M6: `report-evidence` chỉ tính liên kết tới thứ không phải `.md` (log, ảnh, thư mục output); liên kết tài liệu không phải bằng chứng.
- L1 hàm template `cell` escape `|`; L2 bước trống in `TODO`; L3 Background `Given` rồi `And`; L5 chú thích CHANGELOG về ngưỡng 500/800; L6 in "Đã tạo ... và nối dòng" khi `--append` tạo file mới; L11 test CLI.
- Chấp nhận, không sửa: L4 `dk status` không quét `plans/` (plan, phase, report không đếm; cột `%-14s` lệch với tên loại dài); L8 thiếu `--in` mã 1 (nhất quán với thiếu `--set feature` ở phase 7); L9 `--append` không khóa file (cùng loại M6 quyết định 1.9); nit `gwtRe` nhận Given/When/Then ở chữ nội dung khi AC không in đậm; `JoinFeatureComment` thân rỗng thêm một dòng trống.
- Câu hỏi mở của reviewer, giữ quyết định hiện tại: `report` bắt buộc `--in` (kể cả `--in plans`) để không tạo nhầm chỗ; `report-evidence` cảnh báo cả report `draft` vì bằng chứng phải có từ lúc viết; `dk status` quét `plans/` để khi có nhu cầu thật.
