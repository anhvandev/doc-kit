---
name: doc-adr
description: 'Architecture Decision Record: ghi bối cảnh, lựa chọn, quyết định, hệ quả; đánh số, không sửa sau khi chốt, ADR mới thay ADR cũ; chỉ mục ADR. Không dùng cho quyết định nhỏ (Decision log).'
---

# doc-adr

Soạn ADR một file trong `docs/adr/`, đánh số tự động, sinh chỉ mục bằng `dk index adr`.

**Dừng bắt buộc: ADR chỉ đến `status: proposed`; người chốt `accepted`. Không sửa
thân ADR có `status` `accepted`, `superseded` hoặc `deprecated`; đổi ý là ADR mới
với `supersedes`, ADR cũ chỉ đổi frontmatter.**

Mọi file trong `docs/` chỉ được tạo bằng `dk new`. Không tự tạo file tay, kể cả khi
`dk` lỗi. Quy tắc nội dung ở `references/rules.md`; đọc trước khi làm.

## Phạm vi

- Làm: quyết định kỹ thuật quan trọng (công nghệ, cơ sở dữ liệu, giao thức, ranh
  giới dịch vụ, cách triển khai, quy ước ảnh hưởng nhiều tính năng) do người nêu,
  do CR đã duyệt có dòng ADR `Có`, hoặc rút từ thảo luận.
- Không làm: quyết định nhỏ (tên biến, màu, sắp một màn hình, sửa lỗi cục bộ): ghi
  Decision log, nói người. Sơ đồ toàn hệ thống là Architecture overview (họ Overview).

## Bước 0. Kiểm `dk`

```
dk --version
```

Lỗi hoặc không tìm thấy lệnh: dừng, báo "Cần cài `dk` trước (xem README của dk)".

## `new`: từ thảo luận đến ADR đề xuất

### 1. Xác định quyết định và nguồn

Hỏi người **một câu**: "Quyết định cần ghi là gì, và nó đến từ CR, brief hay thảo
luận này?" Có CR hoặc brief: mở file đó, lấy bối cảnh và các phương án đã bàn từ
mục Hỏi đáp hoặc interview. Không thuộc phạm vi ADR (xem Phạm vi): dừng, nói người
ghi Decision log.

```
dk index adr
```

Đọc `docs/adr/README.md`. Có ADR `accepted` cùng chủ đề: đây là thay thế, sang
mục "Thay thế ADR cũ" sau khi soạn.

### 2. Tạo và soạn

```
dk new adr <slug> --set owner="<người phụ trách>" --set source=<CR-id hoặc yymmdd-slug brief>
```

Bỏ `--set source` khi từ thảo luận. Sửa `title` thành câu khẳng định có dấu, ví
dụ "Dùng PostgreSQL làm cơ sở dữ liệu chính". Điền thân từ thảo luận:

| Mục | Nội dung | Thiếu thì hỏi |
|---|---|---|
| 1 Bối cảnh | vấn đề, ràng buộc, vì sao quyết bây giờ | "Điều gì buộc phải quyết lúc này?" |
| 2 Lựa chọn đã xét | 2 đến 4 phương án, ưu, nhược, kết luận | "Ngoài phương án chọn còn xét gì?" |
| 3 Quyết định | một câu khẳng định; người chốt, ngày để trống | không hỏi; người điền khi chốt |
| 4 Hệ quả | điều tốt lên, điều chấp nhận, việc phải làm theo | "Chấp nhận đánh đổi gì?" |
| 5 Liên kết | CR hoặc brief nguồn, Feature Spec và Architecture bị ảnh hưởng, ADR bị thay | tự tìm bằng `dk refs` |

Hỏi **một câu một lượt**, chờ trả lời, ghi rồi hỏi tiếp. Không bịa phương án chưa
bàn. Giữ `status: proposed`. Chạy khối "Sau mỗi lần sửa ADR".

### 3. Bàn giao

Báo người: đường dẫn ADR, mục 2 và 4 để duyệt, cách chốt: đổi `status: accepted`,
điền người chốt và ngày ở mục 3. Sau khi chốt, thân file bất biến; `dk check` báo
lỗi `adr-immutable` nếu đổi.

## Thay thế ADR cũ

Khi quyết định mới thay ADR `accepted` có sẵn (`ADR-000X`):

```
dk new adr <slug> --set owner="<người phụ trách>" --set supersedes=ADR-000X
```

Soạn ADR mới như bước 2; mục 1 nêu vì sao ADR cũ không còn đúng; mục 5 liên kết
đến ADR cũ. Rồi mở ADR cũ, sửa **chỉ frontmatter**: `status: superseded`,
`superseded_by: <mã ADR mới>`. Không đụng thân ADR cũ dù có lỗi chính tả.

Chạy khối "Sau mỗi lần sửa ADR" cho cả hai file. `dk check` phải sạch: có lỗi
`adr-immutable` nghĩa là thân ADR cũ đã đổi; khôi phục thân từ git rồi chạy lại.

ADR không còn áp dụng mà không có gì thay: người đổi `status: deprecated`, ghi lý do
vào changelog; thân vẫn không đổi.

## Sau mỗi lần sửa ADR

```
dk changelog add docs/adr/<ADR-id>-<slug>.md --summary "<tóm tắt thật>" --source <CR-id hoặc brief>
dk render docs/adr/<ADR-id>-<slug>.md
dk index adr
dk check docs/adr/<ADR-id>-<slug>.md
```

Tóm tắt thật nói nội dung ("3 phương án, chọn PostgreSQL vì JSONB"), không "cập
nhật file". Từ thảo luận thì bỏ `--source`. `check` có lỗi thì sửa và chạy lại.

## `update`

- ADR `proposed`: sửa theo lời người, chạy lại khối "Sau mỗi lần sửa ADR".
- ADR `accepted`, `superseded`, `deprecated`: không sửa thân. Đổi ý là "Thay thế
  ADR cũ". Chỉ frontmatter `status`, `superseded_by` được đổi.
- Chỉ mục: chỉ chạy `dk index adr`, không sửa tay `docs/adr/README.md`.

## `html`

```
dk render docs/adr/<ADR-id>-<slug>.md
dk index adr
```

`dk render --all --index` khi người muốn xem toàn bộ.
