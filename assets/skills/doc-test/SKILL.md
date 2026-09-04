---
name: doc-test
description: 'Testing strategy đề xuất từ stack thật, test case Gherkin hoặc bảng sinh từ AC của spec đã duyệt, checklist test giao diện so màn hình với mockup theo mã bước, test report đính output thật; sửa test case sau CR duyệt. Không dùng cho Feature Spec, mockup hay plan.'
---

# doc-test

Soạn Testing strategy, sinh Test case từ mục 9 Feature Spec, checklist test giao
diện theo mã bước, Test report; tất cả trong `docs/test/` bằng `dk new` và `dk check`.

**Dừng bắt buộc: không sinh test case từ spec chưa `status: approved`; không sửa
Feature Spec để hợp thức hóa test (lệch thì báo người, sửa spec qua họ CR); người
chốt công cụ test và quyết định xử lý fail, skill không tự quyết.**

Mọi file trong `docs/` chỉ được tạo bằng `dk new`. Không tự tạo file tay, kể cả khi
`dk` lỗi. Quy tắc nội dung ở `references/rules.md`; đọc trước khi làm.

## Phạm vi

- Làm: `docs/test/strategy.md`; `docs/test/<F-xxx>-cases.feature` (Gherkin) hoặc
  `<F-xxx>-cases.md` (bảng); `docs/test/<F-xxx>-ui.md`; `docs/test/reports/`.
- Không làm: Feature Spec và tiêu chí chấp nhận (họ Feature Spec); mockup và
  Design system (họ Design); plan, report tiến độ (họ Plan); mã test tự động ngoài
  file `.feature` (việc của dự án đích).

## Bước 0. Kiểm `dk`

```
dk --version
```

Lỗi hoặc không tìm thấy lệnh: dừng, báo "Cần cài `dk` trước (xem README của dk)".

## `new`: Testing strategy

### 1. Đọc stack thật

Đọc `package.json`, `go.mod`, `pyproject.toml`, `Gemfile`, cấu hình CI
(`.github/workflows/`, `.gitlab-ci.yml`) và thư mục test có sẵn. Ghi lại công cụ
test đang dùng, lệnh chạy, công cụ BDD nếu có (cucumber, behave, godog, pytest-bdd).

### 2. Tạo và điền

```
dk new testing-strategy strategy --set owner="<người chốt>"
```

| Mục | Cách viết |
|---|---|
| 2 Loại test | unit, integration, e2e, giao diện; mục đích, phạm vi, ai viết; chỉ loại stack đang có hoặc người muốn thêm |
| 3 Công cụ | từ stack thật, ghi phiên bản và lệnh chạy; đề xuất công cụ BDD nếu chưa có |
| 4 Định dạng | `gherkin` khi có (hoặc sẽ có) công cụ BDD, `table` khi không; ghi `format` vào frontmatter |
| 5 Ngưỡng | bao phủ, test chặn merge, test trước phát hành |
| 6 Test giao diện | bỏ mục khi `has_ui: false` toàn dự án |

Chạy khối "Sau mỗi lần sửa". Đặt `status: review`. **Dừng.** Hỏi người chốt công
cụ và định dạng; khi `approved`, nhắc người ghi lệnh BDD vào `dk.toml`:

```
[test]
bdd_cmd = "cucumber"
```

## `new`: Test case từ Feature Spec

### 1. Kiểm spec và strategy

Mở spec người chỉ: cần `status: approved` và mục 9 có `AC1`, `AC2`... Thiếu:
**dừng**, báo người. Đọc `format` trong `docs/test/strategy.md`; chưa có strategy
thì dùng `gherkin` và nói rõ với người.

### 2. Sinh khung

```
dk new test-case <slug> --from docs/features/<F-xxx>-<slug>.md --set owner="<người viết>"
```

`format: table` thì thay `test-case` bằng `test-case-table`. `dk` chép mục 9 thành
một Scenario (hoặc một dòng bảng) cho mỗi AC với tag `@F-xxx @ACn`, mục 3 thành
Background, đặt `source` và `feature` là mã spec.

### 3. Hoàn thiện

- Dòng `# chưa tách được: ...` (hoặc ô "chưa tách được") là AC trong spec không
  theo khung Given / When / Then: điền tay từ chữ gốc và báo người AC đó cần sửa
  trong spec qua họ CR. Không sửa spec.
- Thêm Scenario cho từng dòng ngoại lệ mục 7 của spec với tag `@F-xxx @En`; dữ
  liệu vào và kết quả mong đợi phải cụ thể (thông báo, dữ liệu đổi, sự kiện).
- Giữ nguyên tag và mã; không đổi tên Scenario đã sinh.

### 4. Kiểm chạy được

Gherkin: đọc `dk.toml` `[test] bdd_cmd`. Có giá trị thì chạy dry-run của công cụ
đó với file vừa sinh, ví dụ:

```
cucumber --dry-run docs/test/<F-xxx>-cases.feature
```

Lỗi cú pháp thì sửa file `.feature`, không sửa template hay spec. `bdd_cmd` rỗng
hoặc lệnh không có trên máy: **báo rõ "chưa kiểm chạy được"** trong câu trả lời và
ở mục 5 của test report sau này; không được im lặng.

Chạy khối "Sau mỗi lần sửa". Đặt `status: review`. **Dừng.** Báo người số Scenario,
số ngoại lệ đã thêm; hỏi người duyệt trường hợp biên (`status: approved`).

## `new`: Checklist test giao diện

Chỉ khi spec có mục 6 (`has_ui: true`) và mockup ở cột Mockup đã `approved`:

```
dk new ui-test-checklist <slug> --from docs/features/<F-xxx>-<slug>.md --set owner="<người duyệt>"
```

`dk` sinh một dòng `- [ ] Bn: khớp mockup <liên kết>` cho mỗi mã bước ở mục 6.
Khi chạy: so màn hình thật với mockup và Design system theo từng mã bước; lệch
thì ghi mục 3 (mã bước, điểm lệch, ảnh); không tự sửa mockup hay spec. **Dừng**,
người duyệt sai lệch.

## `update`: Test report theo đợt phát hành

```
dk new test-report v1-2-0 --set title="Đợt 1.2.0" --set owner="<người viết>" --set source=<plan hoặc CR>
```

Slug chỉ nhận `a-z`, `0-9`, gạch ngang (không dấu chấm); số phiên bản đặt ở `title`.

- Mục 2: phiên bản, môi trường, liên kết file test case đã chạy.
- Mục 3: chạy bộ test, dán output thật của công cụ vào khối ```; bảng đếm từ
  output, không ước lượng.
- Mục 4: mỗi fail một dòng, để trống cột Quyết định và Người.
- Mục 5: ghi "Gherkin chưa kiểm chạy được" nếu bước 4 ở trên chưa làm được.

Chạy khối "Sau mỗi lần sửa". **Dừng.** Người quyết định xử lý từng fail (sửa
trước phát hành, chấp nhận, hoãn); điền cột Quyết định theo lời người, đặt
`status: submitted`. Sau đó không sửa; đợt sau viết report mới.

## `update`: Spec đổi theo CR

CR `approved` có dòng "Test case" trong bảng tác động: mở file test case cùng mã
spec, sửa đúng Scenario của AC hoặc mã bước CR nêu; AC mới thì thêm Scenario với
tag `@ACn` kế tiếp. Không sinh lại cả file (mất phần ngoại lệ đã thêm).
`--source <CR-id>` ở changelog.

## Sau mỗi lần sửa

```
dk changelog add docs/test/<file> --summary "<tóm tắt thật>" --source <F-xxx hoặc CR>
dk render docs/test/<file>.md
dk check docs/test/<file>
```

File `.feature` không render HTML (không phải Markdown), vẫn `changelog add` và
`check`. `check` báo lỗi thì sửa rồi chạy lại. Sau khi có test case, `dk check`
spec tương ứng không còn cảnh báo `spec-has-test`; còn thì `source` trong khối
`# dk:` của `.feature` sai mã.

## `html`

```
dk render docs/test/strategy.md
dk render --all
```

Mở `docs/html/test/strategy.html` và `<F-xxx>-ui.html` kiểm liên kết mockup mở
được từ đĩa.
