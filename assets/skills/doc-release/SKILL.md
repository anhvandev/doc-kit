---
name: doc-release
description: 'Tài liệu cho người dùng cuối: Release brief từ Feature Spec bằng ngôn ngữ không kỹ thuật, Release notes gộp theo phiên bản, User guide theo nhiệm vụ, FAQ. Không dùng cho tài liệu vận hành.'
---

# doc-release

Soạn Release brief, Release notes, User guide, FAQ trong `docs/release/` bằng
`dk new` và `dk check`; viết cho người không kỹ thuật.

**Dừng bắt buộc: chỉ soạn Release brief từ Feature Spec `status: implemented`
hoặc `approved` sắp phát hành do người chỉ; người duyệt giọng văn và thêm ảnh
thật; không sửa Feature Spec; không tự bỏ từ khỏi danh sách jargon; không tự đặt
`released_in`.**

Mọi file trong `docs/` chỉ được tạo bằng `dk new`. Không tự tạo file tay, kể cả khi
`dk` lỗi. Quy tắc nội dung ở `references/rules.md`; đọc trước khi làm.

## Phạm vi

- Làm: `docs/release/briefs/<F-xxx>.md`; `docs/release/<phiên bản>.md`;
  `docs/release/guide/<slug>.md` và mục lục `dk index user-guide`;
  `docs/release/faq.md`.
- Không làm: Feature Spec (họ Feature Spec); Deployment, Runbook, Postmortem (họ
  Vận hành); `docs/CHANGELOG.md` (họ Plan); tài liệu quản trị (người viết).
- Công cụ nội bộ không có người dùng ngoài: báo người rằng họ này có thể bỏ,
  `CHANGELOG.md` là đủ; chỉ làm khi người vẫn muốn.

## Bước 0. Kiểm `dk`

```
dk --version
```

Lỗi hoặc không tìm thấy lệnh: dừng, báo "Cần cài `dk` trước (xem README của dk)".

## `new`: Release brief từ Feature Spec

### 1. Kiểm spec

Mở spec người chỉ: `status: implemented`, hoặc `approved` và người xác nhận sắp
phát hành. Khác: **dừng**, báo người. Đọc mục 2, 3, 5, 6, 7, 11.

### 2. Sinh khung

```
dk new release-brief <slug> --from docs/features/<F-xxx>-<slug>.md --set owner="<người duyệt>"
```

`dk` đặt tên file theo mã spec, chép mục 2 thành "Giúp gì", mục 3 thành "Ai dùng",
cột hành động của bảng hành vi thành "Bước cơ bản", cột mockup mục 6 thành "Ảnh
màn hình", mục 11 thành "Giới hạn đã biết". Sửa lỗi thì thêm `--set kind=fix`.

### 3. Viết lại bằng ngôn ngữ người dùng

Toàn bộ chữ chép sang phải viết lại: bỏ mã bước, tên role, tên bảng, tên hàm, từ
kỹ thuật; mỗi bước một câu bắt đầu bằng động từ và nói điều người dùng nhìn thấy.
Thêm vào mục 6 các ngoại lệ ở mục 7 của spec chưa xử lý. Không hứa điều spec
không có. Ảnh: giữ liên kết mockup đã duyệt; ảnh chụp thật do người thêm sau.

### 4. Kiểm

Chạy khối "Sau mỗi lần sửa". `no-jargon` cảnh báo từ nào thì viết lại câu đó;
từ nằm trong tên sản phẩm thì **hỏi người** có bỏ từ khỏi `dk.toml` `[release]
jargon` không, không tự sửa `dk.toml`. Đặt `status: review`. **Dừng.** Người duyệt
giọng văn rồi tự đặt `ready`.

## `new`: Release notes theo phiên bản

Khi người cho biết phiên bản sắp phát hành:

```
dk new release-notes --collect <phiên bản> --set owner="<người duyệt>"
```

`dk` gom mọi Release brief `status: ready` chưa có `released_in`, nhóm "Mới" và
"Sửa lỗi" theo `kind`, ghi `released_in` vào từng brief và một dòng changelog cho
mỗi brief. Không có brief nào thì lệnh báo lỗi: **dừng**, hỏi người brief nào còn
thiếu hoặc chưa `ready`. Điền mục 4 "Cần lưu ý" từ các brief. Chạy khối "Sau mỗi
lần sửa" cho file notes. Đặt `status: review`. **Dừng.** Người duyệt và đặt
`published`; sau đó không sửa.

## `update`: User guide theo nhiệm vụ

Sau mỗi Release notes `published`, với từng brief vừa phát hành:

1. Hỏi người nhiệm vụ nào của người dùng thay đổi (tên nhóm `task`, ví dụ "Bán
   hàng"). Không tổ chức theo menu.
2. Trang đã có: sửa đúng mục 3 và 4, `--source <F-xxx>` ở changelog. Chưa có:

```
dk new user-guide <slug> --from docs/release/briefs/<F-xxx>.md --set owner="<người duyệt>" --set task="<nhóm nhiệm vụ>"
```

3. Viết theo việc người dùng muốn làm; một nhiệm vụ một trang; liên kết câu hỏi
   liên quan trong `faq.md` thay vì lặp lại.
4. `dk index user-guide` sinh mục lục `docs/release/guide/README.md`.

Chạy khối "Sau mỗi lần sửa". **Dừng.** Người duyệt (`approved`).

## `update`: FAQ

Chỉ thêm khi một câu hỏi xuất hiện lần thứ hai (ticket, chat) và người xác nhận:

```
dk new faq --append "<câu hỏi> | <trả lời ngắn> | <liên kết trang guide>" --set owner="<người duyệt>"
```

Lần đầu tạo `docs/release/faq.md` rồi nối; lần sau chỉ nối. Không sửa dòng cũ;
câu trả lời đổi thì nối dòng mới và báo người xóa dòng cũ.

## Sau mỗi lần sửa

```
dk changelog add docs/release/<file> --summary "<tóm tắt thật>" --source <F-xxx hoặc phiên bản>
dk render docs/release/<file>
dk check docs/release/<file>
```

`check` báo lỗi thì sửa rồi chạy lại; `no-jargon` là cảnh báo, xử lý theo bước 4.

## `html`

```
dk render --all
```

Mở `docs/html/release/briefs/<F-xxx>.html` kiểm liên kết mockup mở được từ đĩa;
mở `docs/html/release/guide/README.html` kiểm mục lục theo nhiệm vụ.
