---
name: doc-overview
description: 'Product overview, Architecture overview, Glossary và Feature catalog cấp dự án: soạn từ Product brief đã duyệt hoặc từ mã nguồn, sinh catalog bằng chỉ mục. Không viết Feature Spec, không ghi ADR.'
---

# doc-overview

Soạn ba tài liệu cấp dự án trong `docs/overview/` (Product overview, Architecture
overview, Glossary) và sinh Feature catalog `docs/features/README.md`.

**Dừng bắt buộc: Product overview chỉ soạn từ brief có `level: project`,
`kind: product`, `status: approved`; brief khác thì dừng, báo người. Ba tài liệu chỉ
đến `status: review`; người chốt `approved`. Thành phần kiến trúc suy đoán từ mã
giữ chữ `chưa xác nhận` cho đến khi người xác nhận.**

Mọi file trong `docs/` chỉ được tạo bằng `dk new`. Không tự tạo file tay, kể cả khi
`dk` lỗi. Quy tắc nội dung ở `references/rules.md`; đọc trước khi làm.

## Phạm vi

- Làm: Product overview từ Product brief đã duyệt; Architecture overview từ mã nguồn
  hoặc lời người; Glossary từ thuật ngữ in đậm trong Feature Spec và brief; Feature
  catalog bằng `dk index features`.
- Không làm: Feature Spec (họ Feature Spec); quyết định kỹ thuật (họ ADR); tài liệu
  thiết kế giao diện (họ Design). Product brief chưa có: chuỗi idea, interview, brief.

## Bước 0. Kiểm `dk`

```
dk --version
```

Lỗi hoặc không tìm thấy lệnh: dừng, báo "Cần cài `dk` trước (xem README của dk)".

## `new` Product overview

### 1. Tìm và kiểm brief

```
dk index intake
```

Đọc `docs/intake/README.md`, mở brief người chỉ (hoặc brief `approved` duy nhất có
`level: project`). Kiểm ba trường: `level: project`, `kind: product`,
`status: approved`. Sai một trường: **dừng**, báo người trường nào sai và ai cần
duyệt; không tạo file.

### 2. Tạo và chép

```
dk new product-overview <slug-san-pham> --from docs/intake/<yymmdd>-<slug>/brief.md --set owner="<người phụ trách>"
```

`title`, `outcome`, `acceptance`, `source` được chép từ brief. Điền thân:

| Mục overview | Lấy từ |
|---|---|
| 2 Vấn đề và người dùng | idea.md cùng thư mục (mục Ai gặp, Vấn đề) và lời người |
| 3 Kết quả mong muốn | brief mục 1, chép nguyên |
| 4 Nhóm tính năng chính | brief và Feature catalog: gom tính năng thành 3 đến 7 nhóm, mỗi nhóm một dòng |
| 5 Thước đo thành công | brief mục 4, mỗi dòng đo được |
| 6 Ràng buộc | brief mục 2, chép nguyên |
| 7 Ngoài phạm vi | brief mục 3, chép nguyên |

Không chép hành vi từng tính năng. Giữ dưới 150 dòng; vượt thì cắt mục 4 về một
dòng một nhóm. Đặt `status: review`. Chạy khối "Sau mỗi lần sửa".

## `new` Architecture overview

```
dk new architecture <slug-san-pham> --set owner="<người phụ trách>"
```

Có mã nguồn trong dự án: đọc cấu trúc thư mục hai cấp đầu, manifest phụ thuộc
(`go.mod`, `package.json`, `pyproject.toml`, tương tự), file cấu hình dịch vụ và
triển khai (`docker-compose*`, `Dockerfile`, thư mục `deploy/`, `infra/`), file cấu
hình môi trường mẫu. Từ đó:

1. Mục 2: sơ đồ Mermaid, mỗi nút một thành phần chạy độc lập (dịch vụ, cơ sở dữ
   liệu, hàng đợi, client, dịch vụ ngoài); mũi tên là luồng dữ liệu chính.
2. Mục 3: mỗi nút một dòng; cột "Nơi trong mã" là đường dẫn thật; cột "Xác nhận"
   ghi `đã xác nhận` chỉ khi có bằng chứng trực tiếp (file cấu hình khai báo thành
   phần đó), còn lại `chưa xác nhận`.
3. Mục 4, 5, 6 từ cùng bằng chứng; phiên bản tech stack lấy từ manifest.

Không có mã nguồn: hỏi người **một câu một lượt** theo mục 2 đến 6; không đoán.

Đưa người danh sách dòng `chưa xác nhận`, chờ người xác nhận từng dòng; đổi chữ
theo câu trả lời. Không tự đổi thành `đã xác nhận`. Đặt `status: review`. Chạy
khối "Sau mỗi lần sửa".

## `new` Glossary

```
dk new glossary <slug-san-pham> --set owner="<người phụ trách>"
```

Quét thuật ngữ:

1. Mở từng file `docs/features/F-*.md`, lấy chữ in đậm **lần đầu** trong mục 2, 5,
   8; mở brief `approved` trong `docs/intake/`, lấy chữ in đậm.
2. Mỗi thuật ngữ một dòng bảng mục 2: định nghĩa đề xuất từ ngữ cảnh câu chứa nó;
   cột "Nguồn" là mã Feature Spec hoặc thư mục brief; cột "Không nhầm với" chỉ khi
   có thuật ngữ khác trong danh sách dễ lẫn.
3. Sắp theo bảng chữ cái. Đưa người danh sách để người giữ, sửa hoặc bỏ; không tự
   chốt định nghĩa.

Mục 3 Giọng văn: hỏi người ba điểm (cách gọi người dùng, thì, độ dài câu) nếu chưa có
trong Design brief. Đặt `status: review`. Chạy khối "Sau mỗi lần sửa". `dk check`
sau đó cảnh báo `glossary-term` cho thuật ngữ trong Feature Spec còn thiếu; thêm
dòng hoặc báo người bỏ in đậm.

## Feature catalog

```
dk index features
```

Chỉ chạy lệnh; không viết hay sửa tay `docs/features/README.md`. Trạng thái sai thì
sửa frontmatter Feature Spec tương ứng (việc của họ Feature Spec) rồi chạy lại.

## Sau mỗi lần sửa

Với từng file trong `docs/overview/` vừa tạo hoặc sửa:

```
dk changelog add docs/overview/<file>.md --summary "<tóm tắt thật>" --source <brief-hoặc-CR>
dk render docs/overview/<file>.md
dk check docs/overview/<file>.md
```

Tóm tắt thật nói nội dung ("5 nhóm tính năng, 4 thước đo"), không "cập nhật file".
`--source`: thư mục brief (`<yymmdd>-<slug>`) hoặc mã CR; soạn từ mã thì bỏ cờ.
`check` có lỗi thì sửa nội dung và chạy lại; cảnh báo `line-threshold` thì cắt bớt.

Cuối cùng báo người: đường dẫn file, các mục cần duyệt, dòng `chưa xác nhận` hoặc
thuật ngữ chờ chốt, và cách chốt: đổi `status: approved`.

## `update`

- Tài liệu `draft` hoặc `review`: sửa theo lời người, chạy khối "Sau mỗi lần sửa".
- Tài liệu `approved`: chỉ sửa khi người đưa CR `approved` hoặc brief mới `approved`
  trỏ đến; đặt `source` là mã CR hoặc thư mục brief đó, `--source` cùng giá trị.
  Không có: dừng, nói người thay đổi trên tài liệu đã chốt là Change Request.
- Glossary: thêm thuật ngữ mới bất kỳ lúc nào từ cảnh báo `glossary-term`; định
  nghĩa vẫn do người chốt.
- Feature catalog: chỉ chạy lại `dk index features`.

## `html`

```
dk render docs/overview/product-overview.md
dk render docs/overview/architecture.md
dk render docs/overview/glossary.md
dk index features
```

`dk render --all --index` khi người muốn xem toàn bộ.
