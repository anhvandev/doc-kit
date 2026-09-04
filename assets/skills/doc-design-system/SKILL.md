---
name: doc-design-system
description: 'Design system toàn dự án theo Atomic Design (thư mục docs/design): design tokens, quy tắc nền, atoms, molecules, organisms, templates, patterns, giọng văn giao diện; dựng khi Design brief (do họ Intake soạn) đã được chốt; đổi theo CR chốt. Không dùng cho user flow, wireframe, mockup hay màn hình của một tính năng.'
---

# doc-design-system

Dựng và duy trì Design system của dự án trong `docs/design/` bằng `dk new`,
`dk tokens css` và `dk check`. Bỏ toàn bộ skill này khi phần mềm không có giao diện.

**Dừng bắt buộc: không tạo tokens hay component khi Design brief chưa `status:
approved` với `kind: design`; không sửa tokens, foundations, component `approved` khi
chưa có CR `approved` có dòng "Design: mockup, Design system" là Có. Mọi giá trị hiển
thị đặt tên trong `tokens.json`; không gõ hex hay px ở bất kỳ file nào khác.**

Mọi file trong `docs/` chỉ được tạo bằng `dk new`. Không tự tạo file tay, kể cả khi
`dk` lỗi. Quy tắc nội dung ở `references/rules.md`; đọc trước khi làm.

## Phạm vi

- Làm: `tokens.json` và `tokens.css` sinh ra; `foundations.md`; component theo lớp
  (`atoms/`, `molecules/`, `organisms/`, `templates/`); `patterns/`; giọng văn giao diện;
  kiểm tra thật bằng một mockup trang từ template.
- Không làm: user flow, wireframe, mockup, UI spec của một tính năng (họ Design flow);
  Feature Spec; lint code UI của dự án đích; Tailwind config, theme object, đồng bộ Figma.

## Bước 0. Kiểm `dk`

```
dk --version
```

Lỗi hoặc không tìm thấy lệnh: dừng, báo "Cần cài `dk` trước (xem README của dk)".

## `new`: từ Design brief đến bộ tối thiểu

### 1. Kiểm Design brief

```
dk index intake
```

Mở brief người chỉ. Cần `level: project`, `kind: design`, `status: approved`. Sai một
trường: **dừng**, báo người trường nào sai và ai cần duyệt. Chưa có Design brief: báo
người soạn bằng họ Intake với `level=project`, `kind=design` trước, rồi quay lại. Đọc brief mục 2 Ràng buộc:
có tên thư viện UI (MUI, Ant, shadcn) thì đi nhánh "Thư viện UI có sẵn" ở `rules.md`.

### 2. Lớp 0: tokens

```
dk new design-tokens tokens --from docs/intake/<yymmdd>-<slug>/brief.md --set owner="<người phụ trách>"
```

Sửa `docs/design/tokens/tokens.json` theo brief: thang khoảng cách và chữ, bảng màu
primitive, semantic cho nền, chữ, hành động, viền, bốn trạng thái, bo góc, bóng. Chỉ
sửa giá trị và thêm token; giữ khung nhóm và tên semantic có sẵn để mockup template
dùng được. Mỗi thang dưới 10 bước. Theme tối: đặt giá trị thay thế ở
`$extensions.dk.theme.dark` của token semantic, không đổi primitive. Rồi:

```
dk tokens css
```

Lỗi alias vòng hoặc alias lạ: sửa `tokens.json`, chạy lại. Không sửa `tokens.css`.

### 3. Lớp 0: quy tắc nền

```
dk new design-foundations foundations --from docs/intake/<yymmdd>-<slug>/brief.md --set owner="<người phụ trách>"
```

Điền 7 mục: bảng semantic theo vai trò (tên token, không chép giá trị), lưới và điểm
gãy, accessibility (tương phản 4.5:1, vùng chạm `--size-touch`, focus ring), một bộ
icon, giọng văn (thống nhất với Glossary mục 3), quy tắc lớp. Thư viện UI có sẵn: ghi
`ui_library` và chỉ ghi ngoại lệ.

### 4. Lớp 1 đến 4: component

Đúng thứ tự, mỗi component một lệnh, **hỏi người duyệt sau mỗi lớp** trước khi lên lớp
trên. Component mà Design brief mục 3 loại ra thì bỏ, ghi ở `foundations.md` mục 7:

```
dk new design-component text --set layer=atom --set owner="<người phụ trách>" --set source=<yymmdd-slug của Design brief>
dk new design-component icon --set layer=atom ...
dk new design-component button --set layer=atom ...
dk new design-component input --set layer=atom ...
dk new design-component form-field --set layer=molecule ...
dk new design-component form --set layer=organism ...
dk new design-component empty-state --set layer=organism ...
dk new design-component toast --set layer=organism ...
dk new design-component list-page --set layer=template ...
dk new design-component form-page --set layer=template ...
```

Mỗi file đủ 5 mục nội dung: biến thể, kích thước, trạng thái, quy tắc dùng và không
dùng, accessibility; `uses` liệt kê component lớp ngay dưới. Mọi ô token ghi tên
`--...` có trong `tokens.css`; không ghi giá trị. Dùng thư viện UI có sẵn: lớp 1 đến 3
chỉ ghi component thư viện tương ứng và ngoại lệ; vẫn viết đủ template và pattern.

### 5. Patterns

```
dk new design-pattern <slug> --set owner="<người phụ trách>"
```

Tối thiểu: trạng thái tải, trạng thái rỗng, lỗi hệ thống, xác nhận xóa. Pattern chỉ
ghép organism đã có; mô tả hành vi, không chỉ hình.

### 6. Kiểm tra thật

Tạo một mockup mẫu từ template `list-page`, chỉ dùng organism đã có và biến semantic;
`source` là Design brief:

```
dk new mockup mau-danh-sach --from docs/intake/<yymmdd>-<slug>/brief.md --set feature=F-000 --set step=B0 --set owner="<người phụ trách>" --set title="Mẫu trang danh sách"
```

Điền HTML bằng organism trong `docs/design/organisms/`; mọi màu và khoảng cách qua
`var(--...)`. Chạy khối "Sau mỗi lần sửa". `mockup-tokens` báo lỗi thì hệ chưa đủ
token: bổ sung `tokens.json`, `dk tokens css`, sửa mockup, không gõ giá trị. Qua thì
giữ file `docs/design/mockups/F-000-B0.html` làm mẫu tham chiếu cho họ Design flow
(không xóa; `F-000` là mã dành riêng cho mẫu, không phải tính năng). **Dừng**, báo
người số token, số component theo lớp, số pattern; hỏi duyệt (đổi `status: approved`
cho tokens, foundations, từng component).

### Sau mỗi lần sửa

```
dk tokens css
dk changelog add docs/design/<đường dẫn> --summary "<tóm tắt thật>" --source <thư mục intake của brief hoặc mã CR>
dk render docs/design/<file .md>
dk check docs/design/<đường dẫn>
```

`tokens.json` và mockup `.html` không render; chỉ changelog và check. `tokens.css` là
bản sinh, không ghi changelog. `check` còn lỗi thì sửa và chạy lại, không báo người.

## `update`: sửa theo CR đã chốt

Mở CR người chỉ: cần `status` là `approved` hoặc `in-progress`, bảng tác động có dòng
"Design: mockup, Design system" với "Bị ảnh hưởng: Có", cột "Cần làm gì" nêu token,
component hay pattern cụ thể. Thiếu một điều: **dừng**, báo người CR cần gì.

- Đổi token: sửa `tokens.json`, `dk tokens css`; mockup đã có tự nhận giá trị mới, không
  sửa mockup.
- Thêm biến thể hoặc trạng thái: sửa đúng mục của file component; không đổi lớp của
  component đã có, cần lớp khác thì tạo component mới và `deprecated` cái cũ.
- Component `approved` sửa xong vẫn `approved`; đổi `source` sang mã CR.

Chạy khối "Sau mỗi lần sửa" với `--source <CR-id>`. Báo người: token và component đã
đổi; nhắc họ CR đánh dấu dòng mục 7 của CR và họ Design flow xem mockup nào cần cập nhật.

## `html`

```
dk render docs/design/foundations.md
dk render docs/design/atoms/button.md
```

Mockup `.html` mở thẳng từ đĩa, không cần render; kiểm bằng cách mở file và xác nhận
màu, khoảng cách theo `tokens.css` (đổi một token, `dk tokens css`, mở lại thấy đổi).
