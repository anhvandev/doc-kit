---
name: doc-design-flow
description: 'User flow, wireframe, mockup HTML tĩnh, prototype, UI spec cho giao diện từng tính năng trong docs/design, mã bước chung với Feature Spec, theo design system đã có; sửa khi CR chốt. Không dùng cho token, component, pattern hay Feature Spec.'
---

# doc-design-flow

Soạn User flow, Wireframe, Mockup, Prototype, UI spec cho một tính năng trong
`docs/design/` bằng `dk new` và `dk check`. Cần Design system đã có (`docs/design/tokens/
tokens.css`, `foundations.md`, `organisms/`); chưa có thì dừng, chỉ người sang họ Design
system. Bỏ skill này khi phần mềm không có giao diện.

**Dừng bắt buộc: không làm khi brief chưa `status: approved` hoặc CR chưa `approved`
có dòng "Design: mockup, Design system" là Có; không làm mockup khi wireframe chưa
được người duyệt; mã bước trong user flow là tập con mã bước Feature Spec cùng
`feature`; mockup chỉ dùng organism đã có và biến `var(--...)` từ `tokens.css`, không
gõ hex hay px; không tạo token hay component mới, thiếu thì báo họ Design system.**

Mọi file trong `docs/` chỉ được tạo bằng `dk new`. Không tự tạo file tay, kể cả khi
`dk` lỗi. Quy tắc nội dung ở `references/rules.md`; đọc trước khi làm.

## Phạm vi

- Làm: `flows/<F>-flow.md`, `wireframes/<F>-<B>.md`, `mockups/<F>-<B>.html` mỗi trạng
  thái quan trọng, prototype (liên kết giữa mockup), `<F>-ui-spec.md` khi có công cụ sinh.
- Không làm: tokens, foundations, component, pattern (họ Design system); Feature Spec
  (họ Feature Spec); test giao diện (họ Test).

## Bước 0. Kiểm `dk`

```
dk --version
```

Lỗi hoặc không tìm thấy lệnh: dừng, báo "Cần cài `dk` trước (xem README của dk)".

## `new`: từ brief hoặc CR đến mockup

### 1. Kiểm đầu vào

Tính năng mới: brief `level: feature`, `status: approved`, có mục Giao diện. Tính năng
đã có: CR `approved` hoặc `in-progress` với dòng "Design" là Có và "Cần làm gì" nêu màn
hình. Sai: **dừng**, báo người. Tìm Feature Spec cùng tính năng:

```
dk index features
```

Có spec: mã bước lấy từ mục 3 và 4 của spec. Chưa có spec: đề xuất mã bước ở user flow,
báo người rằng Feature Spec phải dùng lại đúng mã này.

### 2. User flow

```
dk new userflow <slug> --from docs/features/<F-xxx>-<slug>.md --set owner="<người phụ trách>"
```

Chưa có spec: `--from <brief hoặc CR>` và thêm `--set feature=<F-xxx dự kiến>`. Mục 1 là
Mermaid, mỗi nút một màn hình hoặc hành động mang mã bước; mục 2 mỗi mã bước một dòng
với trạng thái cần mockup (bình thường, rỗng, lỗi, đang tải). Chạy khối "Sau mỗi lần
sửa". `userflow-steps` báo mã lạ: sửa flow theo spec, không sửa spec. Chưa có spec thì
`dk check` báo `userflow-steps` cho đến khi Feature Spec tồn tại: bỏ qua lỗi này, chạy
lại `check` sau khi spec được tạo. **Dừng**, hỏi người duyệt flow.

### 3. Wireframe

Mỗi màn hình trong flow một file, đen trắng, text hoặc ASCII:

```
dk new wireframe <slug> --from docs/design/flows/<F-xxx>-flow.md --set step=<B> --set owner="<người phụ trách>"
```

Mục 2 ghi organism dự kiến, tên đúng file trong `docs/design/organisms/`. Chạy khối
"Sau mỗi lần sửa". **Dừng**, hỏi người duyệt từng wireframe (đổi `status: approved`).
Chưa duyệt thì không làm mockup.

### 4. Mockup HTML tĩnh (mặc định)

Mỗi trạng thái quan trọng một file, `state` là `normal`, `empty`, `error`, `loading`;
tên file theo mã bước, trạng thái khác `normal` thêm hậu tố vào slug và `title`:

```
dk new mockup <slug> --from docs/design/wireframes/<F-xxx>-<B>.md --set owner="<người phụ trách>" --set state=empty --set title="<màn hình>: rỗng"
```

Cùng mã bước nhiều trạng thái: file thứ hai trở đi tạo với `--set step=<B>` và `--force`
không được dùng; đặt `step` là mã bước kèm hậu tố trạng thái được ghi ở `rules.md`
("Đặt tên mockup"). Điền thân: chỉ `<section data-organism="<tên>">` với organism đã
có, mọi màu và khoảng cách qua `var(--...)`; giọng văn theo `foundations.md` mục 5.
Style cho organism chép từ mockup mẫu `F-000-B0.html` của họ Design system, không viết
style mới; thiếu token hay organism: **dừng**, báo họ Design system, không gõ giá trị tạm. Chạy khối "Sau mỗi lần sửa"; `mockup-tokens` báo lỗi thì
thay bằng biến. Mở file từ đĩa kiểm màu và khoảng cách.

Dự án có designer dùng công cụ ngoài (Design brief ghi rõ): thay bước này bằng

```
dk new mockup <slug> --from docs/design/wireframes/<F-xxx>-<B>.md --set external=<url> --set owner="<người phụ trách>"
```

và nhờ designer đặt ảnh xuất tĩnh cùng tên đuôi `.png` cạnh file. Không sinh HTML, không
qua lint; người duyệt chịu trách nhiệm nhất quán.

### 5. Prototype

Thêm thẻ `<a href="<F-xxx>-<B kế>.html">` vào nút hành động của mockup để bấm chuyển
màn hình theo flow. Không tạo file riêng; prototype là công cụ ra quyết định, không lưu
lâu: sau khi duyệt xong có thể bỏ liên kết.

### 6. UI spec

Chỉ khi dự án có công cụ thiết kế sinh được (ghi ở Design brief):

```
dk new ui-spec <slug> --from docs/design/flows/<F-xxx>-flow.md --set owner="<người phụ trách>" --set tool="<công cụ>"
```

Dán bản sinh; không viết tay. Mockup HTML tĩnh thì bỏ bước này.

### 7. Báo và bàn giao

**Dừng.** Báo người: flow, số wireframe, số mockup theo trạng thái, liên kết mở từ đĩa;
hỏi duyệt mockup (đổi `status: approved`). Nhắc họ Feature Spec điền mục 5 Giao diện
bằng liên kết `../design/mockups/<F-xxx>-<B>.html` theo mã bước.

### Sau mỗi lần sửa

```
dk changelog add docs/design/<đường dẫn> --summary "<tóm tắt thật>" --source <thư mục intake của brief hoặc mã CR>
dk render docs/design/<file .md>
dk check docs/design/<đường dẫn>
```

Mockup `.html` không render, chỉ changelog và check. `check` còn lỗi thì sửa và chạy
lại, không báo người.

## `update`: sửa theo CR đã chốt

Kiểm CR như bước 1. Sửa giao diện thì sửa mockup trước: đúng màn hình và trạng thái CR
nêu; thêm bước thì thêm nút mã `B2a` ở flow (khớp spec sau khi họ Feature Spec sửa),
wireframe và mockup mới cho bước đó. Đổi `source` sang mã CR; `approved` giữ nguyên.
Chạy khối "Sau mỗi lần sửa" với `--source <CR-id>`. Báo người mockup đã đổi; nhắc họ CR
đánh dấu dòng mục 6 và họ Test cập nhật test giao diện.

## `html`

```
dk render docs/design/flows/<F-xxx>-flow.md
```

Mở `docs/html/design/flows/<F-xxx>-flow.html` kiểm sơ đồ có dòng "Bước:" liệt kê mã.
Mockup mở thẳng file `.html` từ đĩa; đổi một token rồi `dk tokens css` thấy mockup đổi
theo là đúng.
