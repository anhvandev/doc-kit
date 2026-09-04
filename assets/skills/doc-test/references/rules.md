# Quy tắc doc-test

Chép từ báo cáo `plans/reports/consult-260903-1032-tai-lieu-du-an-nguoi-va-ai.md`.
Chỉ quy tắc, không giải thích; giải thích đọc ở mục báo cáo ghi cuối file.

## Nguyên tắc nền

- Một nguồn sự thật cho mỗi tính năng: Feature Spec. Test case sinh ra từ đó,
  không ngược lại. Test lệch spec thì sửa spec qua CR, không sửa test cho khớp
  hành vi lỡ làm.
- Testing strategy là tài liệu thẩm quyền: người chốt công cụ và ngưỡng, AI đề
  xuất từ stack thật. Test case và test giao diện là tài liệu sinh. Test report là
  tài liệu trạng thái.
- Mỗi tài liệu có một `owner`. Không có người chịu trách nhiệm thì không viết.

## Tầng 6: Testing

| Tài liệu | Cách dùng | Loại | Người | AI |
|---|---|---|---|---|
| Testing strategy | một file; cập nhật khi đổi công cụ hoặc ngưỡng | thẩm quyền | chốt | đề xuất |
| Test case | sinh từ mục Tiêu chí chấp nhận của Feature Spec; Gherkin thì file test chính là tài liệu | sinh | duyệt trường hợp biên | sinh, chạy, báo kết quả |
| Test giao diện | visual regression hoặc checklist theo mã bước; chuẩn so sánh là mockup, không phải cảm tính | sinh | duyệt sai lệch | chạy, báo sai lệch |
| Test report | theo đợt phát hành; đính kèm output thật | trạng thái | quyết định xử lý fail | viết |

## Testing strategy

- Chính sách, không phải danh sách test. Một file `docs/test/strategy.md`.
- Nội dung: loại test, công cụ, phạm vi bao phủ mong đợi, cách chạy, định dạng test
  case, ngưỡng và cửa kiểm, cách test giao diện.
- Đề xuất từ stack thật (`package.json`, `go.mod`, cấu hình CI), không đề xuất
  công cụ ngoài stack khi chưa hỏi người.
- `format: gherkin | table` trong frontmatter quyết định loại test case sinh ra.
  Công cụ BDD đã chốt ghi vào `dk.toml` `[test] bdd_cmd`.

## Test case

- Sinh từ Feature Spec `approved`: mỗi AC ở mục 9 một Scenario (Gherkin) hoặc một
  dòng (bảng), mã `ACn` giữ nguyên để truy vết; tag `@<mã spec> @ACn`.
- Thêm trường hợp biên từ mục 7 ngoại lệ với mã `En`. Trường hợp biên do người
  duyệt.
- Mỗi trường hợp: bước thực hiện, dữ liệu vào, kết quả mong đợi cụ thể (thông báo,
  dữ liệu đổi, sự kiện phát ra). "Hoạt động đúng" không phải kết quả mong đợi.
- Gherkin: file `.feature` là tài liệu và là test; phải chạy được bởi công cụ BDD
  đã chốt. Có `bdd_cmd` thì chạy dry-run; không có thì báo "chưa kiểm chạy được".
  Từ khóa Gherkin tiếng Anh (Feature, Background, Scenario, Given, When, Then),
  nội dung bước tiếng Việt; không dùng `# language:` vì dòng đầu file là khối
  metadata `# dk:`.
- Bảng: dùng khi không có công cụ BDD; cột Kết quả điền khi chạy.
- Một tính năng một file test case, cùng mã với Feature Spec
  (`<F-xxx>-cases.feature` hoặc `.md`). Không gộp nhiều tính năng.
- `dk check` quy tắc `spec-has-test`: spec rời `draft` phải có test case trong
  `docs/test/` có `source` là mã spec; checklist giao diện và test report không
  thay được.
- Yêu cầu thay đổi hàng ngày: spec giữ mục 1, 2, 9; chi tiết hành vi để trong
  test. Test case khi đó là nơi chi tiết, vẫn truy vết về AC.

## Test giao diện

- Chỉ khi phần mềm có giao diện (`has_ui: true`). Không có giao diện: bỏ toàn bộ.
- Checklist theo mã bước: mỗi mã bước ở mục 6 của spec một dòng, so màn hình thật
  với mockup đã duyệt và Design system. Chuẩn so sánh là mockup.
- Sai lệch ghi mã bước, điểm lệch, ảnh; người duyệt sai lệch, quyết định sửa màn
  hình hay sửa mockup (qua CR).
- Có công cụ visual regression thì checklist ghi liên kết kết quả, không chép ảnh.

## Test report

- Một đợt kiểm thử một file, theo đợt phát hành, trong `docs/test/reports/`.
- Đính output thật của công cụ; bảng đếm lấy từ output.
- Mỗi fail: mã test case, hiện tượng, nguyên nhân nếu biết, quyết định (sửa trước
  phát hành, chấp nhận, hoãn), người quyết định. AI không tự quyết.
- Sau `submitted` không sửa; đợt sau viết report mới.

## Tách file và ngưỡng dòng

- Testing strategy một file. Test case theo tính năng, đi kèm Feature Spec cùng
  mã. Test report một đợt một file.
- Vượt 300 đến 500 dòng thì xem lại; trần cứng 800 dòng (`line-threshold`, chỉ
  áp cho Markdown; `.feature` dài thì tách theo nhóm Scenario nhưng vẫn một tính
  năng một file chính).
- Frontmatter thống nhất: `type`, `title`, `status`, `owner`, `created`,
  `updated`, `source`, `feature`. File `.feature` giữ metadata trong khối chú
  thích `# dk:` đầu file; `dk new` điền, skill không sửa `created`, `created_by`,
  `dk_version`.

## Khi nào không áp dụng

- Không có giao diện (API, CLI, batch): bỏ test giao diện.
- Dự án dưới 5 tính năng hoặc dưới 2 tháng: vẫn giữ test case theo tính năng, có
  thể bỏ test report riêng nếu CI đã lưu output.
- Ngành có quy chuẩn riêng: danh sách test và report do quy chuẩn quyết định.

## Trạng thái (khóa tiếng Anh trong `types.toml`)

| Khóa | Nghĩa |
|---|---|
| `draft` | đang soạn hoặc vừa sinh |
| `review` | chờ người duyệt (công cụ, trường hợp biên) |
| `approved` | đã chốt |
| `submitted` | test report đã nộp, không sửa |

Giải thích: báo cáo mục 2 (Tầng 6), mục 3 (mẫu Feature Spec, mục 9 nuôi test),
mục 7 (Tách file), mục 10 (Khi nào khuyến nghị không còn đúng).
