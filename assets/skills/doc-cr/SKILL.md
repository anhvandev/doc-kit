---
name: doc-cr
description: 'Change Request trong docs/cr cho thay đổi trên tính năng, giao diện, kỹ thuật đã có: hiện trạng, đề xuất, bảng tác động chờ duyệt, quyết định, tài liệu đích. Không dùng cho ý tưởng mới.'
---

# doc-cr

Soạn một Change Request (CR) một file trong `docs/cr/`, bảng tác động từ `dk refs`.

**Dừng bắt buộc: sau khi điền bảng tác động, đặt `status: review` và chờ người duyệt
bảng. Không sửa bất kỳ tài liệu đích nào (Feature Spec, mockup, Design system, ADR,
test, guide) trước khi người duyệt bảng và ghi mục Quyết định.**

Mọi file trong `docs/` chỉ được tạo bằng `dk new`. Không tự tạo file tay, kể cả khi
`dk` lỗi. Quy tắc nội dung ở `references/rules.md`; đọc trước khi làm.

## Phạm vi

- Làm: thay đổi trên thứ đã có tài liệu mô tả: sửa hành vi, đổi giao diện, đổi kỹ
  thuật, sửa lỗi. Kể cả thay đổi trên brief đã duyệt.
- Không làm: ý tưởng chưa có gì tồn tại. Việc đó là chuỗi idea, interview, brief.

## Bước 0. Kiểm `dk`

```
dk --version
```

Lỗi hoặc không tìm thấy lệnh: dừng, báo "Cần cài `dk` trước (xem README của dk)".

## `new`: từ yêu cầu đến CR chờ duyệt

### 1. Tạo CR và ghi hiện trạng

```
dk new cr <slug> --set priority=medium --set requester="<người yêu cầu>" --set owner="<người phụ trách>"
```

`priority`: `low`, `medium`, `high`. `owner` là người đang làm việc cùng. Sửa `title`
thành câu đầy đủ có dấu.

Chạy `dk index features`; chỉ có một spec hoặc tiêu đề khớp yêu cầu thì lấy luôn.
Nhiều thì hỏi người **một câu**: "Thứ sắp đổi đang được mô tả ở đâu?" (Feature Spec,
mockup, ADR). Người không biết thì:

```
dk status
dk index features
```

rồi đọc `docs/features/README.md`, gợi ý 2 đến 3 Feature Spec có tiêu đề gần nhất để
người chọn. Mục 1 Hiện trạng: mỗi dòng một liên kết tương đối đến tài liệu đang mô tả
thứ sắp đổi, kèm mã bước hoặc màn hình cụ thể. Không chép nội dung tài liệu vào CR.

### 2. Thay đổi đề xuất và lý do

Mục 2 ghi từ lời người: đổi gì, vì sao, ai yêu cầu, chuyện gì xảy ra nếu không đổi.
Thiếu ý nào thì hỏi đúng ý đó, một câu một lần. Đặt `status: analyzing`.

### 3. Bảng tác động

Với **mỗi** tài liệu ở mục 1:

```
dk refs docs/features/<file>.md
```

Đọc mọi file trong hai danh sách "Đi" và "Đến" (mở từng file, tìm phần nói về thứ
sắp đổi). Điền 6 dòng bảng mục 3, cột "Bị ảnh hưởng" ghi `Có` hoặc `Không`, cột
"Cần làm gì" ghi việc cụ thể kèm liên kết file đã có; file sẽ tạo ghi đường dẫn trong
backtick (liên kết tới file chưa có là lỗi `link-broken`):

| Dòng | Nguồn xét | Ví dụ "Cần làm gì" |
|---|---|---|
| Feature Spec | file hiện trạng và spec khác trong "Đến" | Sửa B3, thêm ngoại lệ hết hàng ở B4 |
| Design: mockup, Design system | `docs/design/`, mockup liên kết từ mục Giao diện của spec | Sửa màn hình F-012-B3; không cần organism mới |
| ADR | `docs/adr/` có trong "Đi" hoặc "Đến" | ADR-0004 bị thay thế một phần |
| Data model, API | mục Data và API của spec, `docs/api/` | Thêm cột `status_filter`, contract `/orders` đổi |
| Test case | `docs/test/` có `source` trỏ về spec | Sửa test B3, thêm test B4 |
| Release brief, User guide | `docs/release/`, `docs/guide/` liên kết đến spec | Đoạn "Lọc đơn" trong guide |

Không tìm thấy file cho một dòng thì ghi `Không` và "không có tài liệu". Tác động
không xác định được từ tài liệu thì ghi `chưa rõ` ở cột "Bị ảnh hưởng".

Đặt `status: review`, chạy khối lệnh "Sau mỗi lần sửa CR" bên dưới. **Dừng.** Báo
người: đường dẫn CR, bảng tác động, các dòng `chưa rõ`; hỏi người duyệt bảng (đồng
ý, sửa, hay từ chối CR).

### Sau mỗi lần sửa CR

Chạy ngay sau mọi lần ghi vào file CR (kể cả khi đang chờ duyệt), để file luôn có
dòng changelog trước khi người commit:

```
dk changelog add docs/cr/<CR-id>.md --summary "<tóm tắt thật>" --source <CR-id>
dk render docs/cr/<CR-id>.md
dk index cr
dk check docs/cr/<CR-id>.md
```

Tóm tắt thật nói nội dung ("Bảng tác động: 3 Có, ADR chưa rõ"), không "cập nhật file".
`check` có lỗi thì sửa nội dung và chạy lại.

### 4. Hỏi đáp (chỉ khi còn `chưa rõ`)

Chỉ cho dòng bảng còn `chưa rõ` hoặc thay đổi ở mục 2 còn mơ hồ. Từng câu một, mở
đầu bằng tên dòng đang làm rõ, chờ trả lời, ghi khối "Câu N" vào mục 4 gồm câu hỏi,
trả lời, **Quyết định rút ra**. Cập nhật dòng bảng tương ứng. Câu về giao diện
kèm 2 đến 3 phương án tham chiếu để người chọn.

Mục 4 dài quá 40 dòng:

```
dk new interview <slug> --from docs/cr/<CR-id>.md
```

File vào `docs/cr/<CR-id>/interview.md`; mục 4 chỉ giữ một liên kết đến đó.
Không còn `chưa rõ` và không có hỏi đáp thì xóa mục 4.

### 5. Sau khi người duyệt

Người trả lời bảng tác động. Ghi mục 5 Quyết định: kết luận `chấp nhận`, `từ chối`
hoặc `hoãn`; lý do; người chốt; ngày. Đặt `status`:

- chấp nhận: `approved`, rồi `in-progress` khi bắt đầu cập nhật tài liệu đích;
- từ chối: `rejected`, giữ nguyên file, không xóa;
- hoãn: giữ `review`, ghi lý do hoãn ở mục 5.

Với CR `approved`: mục 6 liệt kê từng tài liệu đích (mỗi dòng `Có` ở bảng là ít
nhất một dòng `- [ ] [file](đường dẫn)` kèm việc). Việc cập nhật từng tài liệu đích
thuộc skill của họ tài liệu đó (họ Feature Spec, họ Design, họ ADR, họ Test, họ
Release và User guide); giao việc bằng cách nêu tên họ, đường dẫn file và dòng
"Cần làm gì", kèm `--source <CR-id>` để tài liệu đích liên kết ngược. Khi một tài
liệu đích đã cập nhật, đánh dấu `- [x]` dòng đó. Mọi dòng xong: `status: closed`.
Mỗi lần sửa chạy lại khối "Sau mỗi lần sửa CR".

### 6. Quy tắc tỷ lệ

- Sửa lỗi nhỏ, tác động chỉ một Feature Spec và test: bảng mục 3 giữ hai dòng
  Feature Spec và Test case, xóa bốn dòng còn lại, bỏ mục 4. Vẫn dừng ở `review`.
- Bảng có `Có` ở dòng Design system hoặc ADR: nhắc rõ trong báo cáo cho người
  "cần duyệt bảng trước mọi thay đổi"; không giao việc cho tài liệu đích nào,
  kể cả Feature Spec, trước khi người duyệt.
- Người từ chối: giữ file với `status: rejected`, ghi lý do ở mục 5.

## `update`: sửa CR đang mở

CR `new`, `analyzing`, `review`: sửa theo lời người, chạy lại bước 3 nếu mục 1 hoặc
3 đổi, rồi dừng ở `review` như cũ. CR `approved` hoặc `in-progress`: chỉ sửa mục 6
(đánh dấu xong) và không đổi mục 5. CR `closed` hoặc `rejected`: không sửa; thay
đổi mới là CR mới.

## `html`

```
dk render docs/cr/<CR-id>.md
dk index cr
```
