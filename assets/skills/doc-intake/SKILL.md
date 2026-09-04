---
name: doc-intake
description: 'Tiếp nhận ý tưởng mới chưa có gì tồn tại trong docs/intake: ghi idea, phỏng vấn từng trường "chưa rõ", soạn Intake brief, Product brief hoặc Design brief ở trạng thái chờ duyệt. Không dùng cho thay đổi trên thứ đã có.'
---

# doc-intake

Dẫn chuỗi ba bước `idea.md`, `interview.md`, `brief.md` trong `docs/intake/<yymmdd>-<slug>/`.

**Dừng bắt buộc: brief chỉ đến `status: review` (chờ duyệt). Không tự đặt `approved`,
không điền `approved_by`, không tạo Feature Spec hay Design từ brief chưa duyệt.**

Mọi file trong `docs/` chỉ được tạo bằng `dk new`. Không tự tạo file tay, kể cả khi
`dk` lỗi. Quy tắc nội dung ở `references/rules.md`; đọc trước khi làm.

## Phạm vi

- Làm: ý tưởng mới, chưa có Feature Spec, mockup hay ADR nào mô tả nó. Cấp dự án
  (sản phẩm mới, hệ giao diện lần đầu) hoặc cấp tính năng.
- Không làm: thay đổi trên thứ đã có (sửa hành vi, đổi giao diện, sửa lỗi). Việc đó
  là Change Request; nói người dùng `doc-cr`.

## Bước 0. Kiểm `dk`

```
dk --version
```

Lỗi hoặc không tìm thấy lệnh: dừng, báo "Cần cài `dk` trước (xem README của dk)".
Không làm gì tiếp.

## `new`: ý tưởng mới đến brief chờ duyệt

### 1. Ghi idea

Hỏi người **một câu duy nhất**: tên ý tưởng (để đặt slug) và cấp: dự án hay tính năng.
Không hỏi gì khác ở bước này.

```
dk new idea <slug> --set level=project      # hoặc level=feature
```

Mở file vừa tạo, đi qua **từng mục cùng người**, mỗi mục một lượt: đọc câu gợi ý
trong chú thích, ghi câu trả lời của người. Người không biết thì để nguyên đúng chữ
`chưa rõ`; không đoán, không điền hộ. Mục 8 Giao diện: nếu ý tưởng không có giao
diện ghi `không`. Sửa `title` thành tên đầy đủ có dấu.

### 2. Rà idea, quyết định có phỏng vấn không

Đọc lại `idea.md`. Liệt kê ra cho người:

- mọi mục (và dòng con trong mục 8) còn `chưa rõ`;
- mâu thuẫn giữa các mục (ví dụ mục 4 hứa điều mục 5 loại bỏ).

Không có gì trong danh sách: bỏ phỏng vấn, sang bước 4. Có: sang bước 3.

### 3. Phỏng vấn từng câu một

```
dk new interview <slug> --from docs/intake/<yymmdd>-<slug>/idea.md
```

Đổi `status` của `idea.md` sang `interviewing`. Chép danh sách ở bước 2 vào mục
"Trường cần làm rõ" của `interview.md`.

Rồi với **từng** dòng trong danh sách đó, theo thứ tự:

1. Đặt **một** câu hỏi, mở đầu bằng tên mục đang làm rõ, ví dụ "Mục 2 Ai gặp: ...".
   Câu hỏi sinh từ đúng chỗ trống hoặc mâu thuẫn, không từ danh sách câu hỏi có sẵn.
2. Câu về giao diện (mục 8): kèm 2 đến 3 phương án tham chiếu cụ thể (màn hình có
   sẵn trong dự án, template trang, hoặc mẫu phổ biến) để người chọn, không hỏi mở.
3. Chờ người trả lời. Ghi vào `interview.md` khối "Câu N" gồm câu hỏi, trả lời,
   **Quyết định rút ra** (một câu, dạng khẳng định) và **Trạng thái** `đã chốt`
   hoặc `còn mở`.
4. Quyết định `đã chốt` thì ghi ngược vào mục tương ứng của `idea.md` thay cho
   `chưa rõ`.

Dừng phỏng vấn khi hết danh sách hoặc người yêu cầu dừng. Điền mục "Tổng kết"
(mỗi quyết định một dòng, dẫn số câu), đổi `status` interview sang `closed`.
Sau `closed` không sửa `interview.md` nữa. Trường còn `chưa rõ` do người dừng sớm
thì giữ nguyên chữ đó trong idea và nêu trong brief.

Đối chiếu để tự kiểm:

| Sai | Đúng |
|---|---|
| "Bạn có thể mô tả người dùng mục tiêu, tần suất và bối cảnh sử dụng không?" khi mục 2 đã điền | Không hỏi; mục 2 không còn `chưa rõ` |
| Hỏi 5 câu liền rồi chờ | Một câu, chờ, ghi, câu kế |
| "Giao diện trông thế nào?" | "Mục 8 Màn hình liên quan: đặt bộ lọc ở (a) thanh bên trái như trang Sản phẩm, (b) hàng chip trên bảng, hay (c) hộp thoại như trang Báo cáo?" |

### 4. Soạn brief

```
dk new brief <slug> --from docs/intake/<yymmdd>-<slug>/idea.md --set kind=feature --set owner="<người phụ trách>"
```

`kind`: `product` khi ý tưởng sản phẩm cấp dự án; `design` khi ý tưởng giao diện
toàn sản phẩm (Design brief); `feature` cho tính năng. `level` chép từ idea.
Design brief chạy sau hoặc song song ý tưởng sản phẩm, không trước.

Soạn bốn mục cố định từ idea và các quyết định đã chốt:

1. Kết quả mong muốn (từ mục 4 idea, kiểm chứng được);
2. Ràng buộc (mục 6);
3. Ngoài phạm vi (mục 5);
4. Tiêu chí chấp nhận (từ mục 7 và các quyết định; mỗi dòng kiểm tra được).

Mục "Quyết định đã chốt": mỗi quyết định một dòng, dẫn `[Câu N](./interview.md)`.
Không phỏng vấn thì xóa mục này. Điền `outcome` (một dòng) và `acceptance` (danh
sách) trong frontmatter cho khớp mục 1 và 4. `owner` là người đang làm việc cùng
(bắt buộc, `dk check` báo lỗi khi rỗng). Đặt `status: review`. Đổi `status` của `idea.md` sang `briefed`.

### 5. Changelog, render, bàn giao

Với **từng** file đã tạo hoặc sửa (idea, interview, brief):

```
dk changelog add docs/intake/<yymmdd>-<slug>/<file>.md --summary "<tóm tắt thật>" --source <yymmdd>-<slug>
dk render docs/intake/<yymmdd>-<slug>/<file>.md
```

Tóm tắt thật: nói nội dung ("Ghi 7 mục, 2 mục chưa rõ: ai gặp, thước đo"), không
"cập nhật file". Rồi `dk check docs/intake/<yymmdd>-<slug>/brief.md`; có lỗi thì
sửa nội dung và chạy lại. Cuối cùng báo người: đường dẫn brief, bốn mục để duyệt,
trường còn `chưa rõ` (nếu có), và cách chốt: đổi `status: approved`, ghi
`approved_by`, `approved_at`.

## `update`: sửa brief

Chỉ áp cho brief `status` là `draft` hoặc `review`. Sửa theo lời người, ghi thêm
quyết định mới vào interview (nếu interview còn `open`) hoặc vào mục "Quyết định
đã chốt" kèm ngày; rồi làm lại bước 5 cho file đã sửa.

Brief đã `approved`: không sửa. Nói người: thay đổi trên brief đã duyệt là Change
Request, dùng `doc-cr`.

## `html`

```
dk render docs/intake/<yymmdd>-<slug>/brief.md
dk index intake
```

`dk render --all --index` khi người muốn xem toàn bộ.
