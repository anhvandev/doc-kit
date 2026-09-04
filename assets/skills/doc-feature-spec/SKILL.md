---
name: doc-feature-spec
description: 'Feature Spec trong docs/features, nguồn sự thật của một tính năng: sơ đồ luồng có mã bước, bảng hành vi, ngoại lệ, tiêu chí chấp nhận; tạo từ Intake brief người đã duyệt, cập nhật theo CR đã chốt. Không dùng cho brief, CR, mockup, test.'
---

# doc-feature-spec

Soạn và cập nhật Feature Spec, một tính năng một file trong `docs/features/`, bằng
`dk new feature-spec` và `dk check`.

**Dừng bắt buộc: không tạo spec khi brief chưa `status: approved`; không sửa spec
`approved` khi chưa có CR `approved` hoặc `in-progress` có dòng "Feature Spec: Có"
trong bảng tác động. Sơ đồ và bảng hành vi cùng một tập mã bước; chèn bước giữa
thì đặt hậu tố (`B2a`), không bao giờ đánh số lại.**

Mọi file trong `docs/` chỉ được tạo bằng `dk new`. Không tự tạo file tay, kể cả khi
`dk` lỗi. Quy tắc nội dung ở `references/rules.md`; đọc trước khi làm.

## Phạm vi

- Làm: Feature Spec 11 mục (mục 1 là frontmatter) từ brief cấp tính năng đã duyệt;
  sửa spec theo CR đã chốt; chọn biến thể `format` theo bối cảnh; render HTML.
- Không làm: brief và phỏng vấn (họ Intake); CR và bảng tác động (họ CR); mockup,
  Design system (họ Design); test case (họ Test). Không sửa test hay brief thay cho
  spec.

## Bước 0. Kiểm `dk`

```
dk --version
```

Lỗi hoặc không tìm thấy lệnh: dừng, báo "Cần cài `dk` trước (xem README của dk)".

## `new`: từ brief đã duyệt đến spec chờ duyệt

### 1. Tìm và kiểm brief

```
dk index intake
```

Đọc `docs/intake/README.md`, mở brief người chỉ. Kiểm ba trường: `level: feature`,
`kind: feature`, `status: approved`. Sai một trường: **dừng**, báo người trường nào
sai và ai cần duyệt; không tạo file. Nguồn là CR `approved` yêu cầu tính năng mới
(bảng tác động không có Feature Spec nào để sửa): dùng `--from <file CR>` ở bước 3.

### 2. Chọn định dạng và loại sơ đồ

Chọn `format` theo bảng ở `rules.md` mục "Định dạng theo bối cảnh". Không rõ thì hỏi
người **một câu**: "Tính năng này thuộc bối cảnh nào: nghiệp vụ nhiều tác nhân và
ngoại lệ, team viết test Gherkin, CRUD đơn giản, hay vòng đời trạng thái phức tạp?"
Không có câu trả lời rõ thì dùng `spec`.

| `format` | Khác `spec` ở |
|---|---|
| `spec` | mặc định: flowchart mục 4, bảng hành vi mục 5, Given / When / Then mục 9 |
| `use-case` | mục 4 sequenceDiagram, mục 5 Use Case Specification có bảng luồng chính theo mã bước |
| `story` | mục 9 Gherkin thuần trong khối ```gherkin |
| `crud` | bỏ mục 4; mục 5 là bảng field và quyền, không có mã bước |
| `state` | như `spec`, thêm stateDiagram "Vòng đời trạng thái" trong mục 4 |

Loại sơ đồ theo bản chất: chuỗi hành động có rẽ nhánh dùng flowchart; nhiều thành
phần tương tác theo thứ tự dùng sequenceDiagram; vòng đời trạng thái dùng
stateDiagram. Luồng tuyến tính dưới bốn bước: giữ `spec` nhưng xóa khối mermaid ở
mục 4 và ghi một dòng "Luồng tuyến tính, xem bảng mục 5"; bảng mục 5 vẫn dùng mã
bước (`step-codes` chỉ so khi có sơ đồ). Tính năng không có giao
diện (API, CLI, batch): thêm `--set has_ui=false`.

### 3. Tạo file

```
dk new feature-spec <slug> --from docs/intake/<yymmdd>-<slug>/brief.md --set owner="<người phụ trách>" --set format=<format>
```

`title`, `purpose`, `acceptance`, `source` được chép từ brief. Mã `F-xxx` do `dk` cấp.
Không sửa `created`, `created_by`, `dk_version`.

### 4. Điền thân

| Mục | Lấy từ | Cách viết |
|---|---|---|
| 2 Mục đích | brief mục 1 | hai đến ba câu, người dùng làm được gì mà trước không |
| 3 Tác nhân | brief, idea.md cùng thư mục | vai trò, đăng nhập, dữ liệu, quyền |
| 4 Sơ đồ và 5 Bảng hành vi | lời người, brief | soạn **cùng lúc**: thêm một nút là thêm một dòng bảng cùng mã |
| 6 Giao diện | mockup trong `docs/design/mockups/` | chỉ liên kết theo mã bước; chưa có thì ghi "chưa có, xem họ Design" |
| 7 Ngoại lệ | suy từ từng bước: rỗng, lỗi, hết hạn, mất mạng, trùng, quyền | mỗi dòng: mã, tại bước, điều kiện, phản hồi quan sát được; ít nhất nửa số dòng bảng mục 5 |
| 8 Quy tắc | brief mục 2 Ràng buộc | mã R1, R2; quy tắc chung liên kết đến nơi định nghĩa |
| 9 Tiêu chí | brief mục 4 Tiêu chí chấp nhận | AC1, AC2, Given / When / Then; mỗi tiêu chí nêu thông báo cụ thể, dữ liệu đổi, sự kiện phát |
| 10 Dữ liệu và API | ADR, API spec có sẵn | chỉ liên kết |
| 11 Ngoài phạm vi | brief mục 3 | chép, thêm thứ người nói không làm |

Mục 4 và 5 ràng buộc bằng mã bước: mỗi nút trong sơ đồ có đúng một dòng trong bảng
mục 5, cùng mã; nhánh thay thế cần bước riêng thì đặt `B2a`, `B2b` theo bước gốc và
ghi vào cả hai nơi. Bảng mục 5 không ghi "hệ thống xử lý" mà không nói kết quả nhìn
thấy. Phần thân có thuật ngữ nghiệp vụ in đậm ở mục 2, 5, 8 thì thuật ngữ đó phải có
trong Glossary; `dk check` cảnh báo nếu thiếu.

Xong thân: đặt `status: review`, chạy khối "Sau mỗi lần sửa spec". **Dừng.** Báo người
đường dẫn, số bước, số ngoại lệ, số tiêu chí; hỏi người duyệt (đổi `status: approved`).

### Sau mỗi lần sửa spec

```
dk changelog add docs/features/<F-xxx>-<slug>.md --summary "<tóm tắt thật>" --source <brief hoặc CR>
dk render docs/features/<F-xxx>-<slug>.md
dk index features
dk check docs/features/<F-xxx>-<slug>.md
```

`--source` là `<yymmdd>-<slug>` của thư mục intake khi tạo từ brief, là mã CR khi sửa
theo CR. `check` báo `step-codes` thì sửa sơ đồ hoặc bảng cho khớp; báo
`spec-section-order` thì khôi phục đúng tiêu đề số; báo lỗi khác thì sửa nội dung và
chạy lại. Không báo người khi `check` còn lỗi.

## `update`: sửa spec theo CR đã chốt

### 1. Kiểm CR

Mở file CR người chỉ. Cần đủ ba điều: `status` là `approved` hoặc `in-progress`;
bảng tác động mục 4 có dòng "Feature Spec" với "Bị ảnh hưởng: Có"; cột "Cần làm gì"
nêu mã bước hoặc mục cụ thể. Thiếu một điều: **dừng**, báo người CR cần gì; không sửa
spec. Spec `approved` mà không có CR như vậy: từ chối sửa, chỉ người sang họ CR.

### 2. Sửa đúng chỗ CR nêu

- Chỉ sửa mã bước và mục ghi ở "Cần làm gì". Bước khác giữ nguyên chữ.
- Thêm bước giữa `B2` và `B3`: đặt `B2a` trong sơ đồ và bảng mục 5; không đổi `B3`
  trở đi. Bỏ bước: xóa ở cả hai nơi, không dồn số.
- Thêm ngoại lệ: dòng mới ở mục 7 với mã `E` kế tiếp, tại bước CR nêu.
- Đổi `source` thành mã CR. Thêm (hoặc nối vào) mục "Lịch sử nguồn" ở cuối file, sau
  mục 11, dạng `### Lịch sử nguồn` với một dòng cho mỗi nguồn cũ: `- <source cũ>:
  <ngày>`. Không đổi tiêu đề số của 11 mục.
- Giữ `status` hiện có; spec `approved` sửa theo CR vẫn `approved`; spec `implemented`
  sửa theo CR đổi về `approved` cho đến khi phát hành lại.

Chạy khối "Sau mỗi lần sửa spec" với `--source <CR-id>`. `dk check` báo
`cr-approval-order` là spec bị sửa khi CR còn `review`: khôi phục spec từ git và dừng.
Rule so `updated` đến phút và chỉ khi CR chưa chốt, nên không thay được bước 1: luôn
kiểm CR trước khi sửa.
Báo người: mã bước đã đổi, đường dẫn spec; nhắc họ CR đánh dấu dòng mục 7 của CR.

### 3. Spec `draft` hoặc `review`

Sửa theo lời người, không cần CR; chạy khối "Sau mỗi lần sửa spec" với `--source`
là thư mục intake của brief.

## `html`

```
dk render docs/features/<F-xxx>-<slug>.md
dk index features
```

Mở `docs/html/features/<F-xxx>-<slug>.html` kiểm: dưới sơ đồ có dòng "Bước:" liệt kê
mã, mỗi mã bấm nhảy đến đúng dòng bảng hành vi. Thiếu mã nào thì sơ đồ và bảng
đang lệch: sửa Markdown, không sửa HTML.
