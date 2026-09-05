# Quy tắc doc-design-flow

Chép từ báo cáo `plans/reports/consult-260903-1032-tai-lieu-du-an-nguoi-va-ai.md`.
Chỉ quy tắc, không giải thích; giải thích đọc ở mục báo cáo ghi cuối file.

## Nguyên tắc nền

- Design system là cấp dự án; User flow, Wireframe, Mockup, Prototype, UI spec là cấp
  tính năng hoặc cấp màn hình, theo Design system, không định nghĩa style mới (lớp 5
  Pages: template đổ dữ liệu thật).
- Mockup là tài liệu thẩm quyền: người vẽ hoặc duyệt; AI sinh HTML theo Design system
  cho màn hình đơn giản, không thay designer cho màn hình quan trọng. AI không tự sửa
  mockup đã duyệt nếu chưa có người duyệt (qua CR).
- Chỉ áp dụng khi phần mềm có giao diện.

## Thứ tự trong tầng

1. User flow: chuỗi màn hình và hành động cho một tính năng; Mermaid trong Markdown,
   dùng chung mã bước với sơ đồ luồng của Feature Spec. `dk check` quy tắc
   `userflow-steps`: mã bước là tập con mã bước của Feature Spec có id bằng `feature`;
   `feature` không trỏ đến spec nào là lỗi.
2. Wireframe: màn hình có những gì, bố cục ra sao. Đen trắng, không màu; chốt "có gì trên
   màn hình" trước khi vẽ đẹp. Một màn hình một file, tên theo mã tính năng và mã bước.
   Người duyệt rồi mới làm mockup.
3. Mockup: màn hình trông như thế nào, theo Design system. Mỗi trạng thái quan trọng một
   mockup: bình thường, rỗng, lỗi, đang tải.
4. Prototype: bản bấm được từ mockup để test với người dùng hoặc duyệt trước khi code.
   Công cụ ra quyết định, không lưu lâu dài.
5. UI spec / Redline: kích thước, khoảng cách, responsive, animation. Sinh từ công cụ
   thiết kế, không viết tay; không có công cụ sinh thì không có file này.

Đầu vào: mục Giao diện trong brief với tính năng mới, hoặc bảng tác động của CR với
thay đổi giao diện trên tính năng đã có. Feature Spec liên kết đến mockup theo mã bước,
không chép ảnh vào spec. Test giao diện so với mockup đã duyệt.

## Mockup HTML tĩnh (mặc định của dự án)

- File `.html` tĩnh trong repo, mở từ đĩa không mạng; metadata trong chú thích
  `<!-- dk: ... -->` đầu file.
- Chỉ dùng `<link rel="stylesheet" href="../tokens/tokens.css">` và biến `var(--...)`.
  `dk check` quy tắc `mockup-tokens`: giá trị hex (`#333`) hoặc px (`12px`) gõ tay trong
  `<style>` hoặc `style=""` là lỗi; file `.html` trong `docs/` thiếu khối `<!-- dk: -->`
  cũng là lỗi. Lint chỉ bắt hex và px (không bắt `rgb()`, `rem`, `%`: cũng không được gõ,
  người kiểm khi duyệt); px trong prelude `@media` được phép vì biến CSS không dùng
  được ở đó. "Chỉ dùng organism đã có" do người kiểm khi duyệt.
- Không định nghĩa style cho component trong mockup; thiếu token hoặc organism là việc
  của họ Design system.
- Ngoại lệ có designer dùng Figma (Design brief ghi rõ): `dk new mockup --set
  external=<url>` lưu liên kết và ảnh xuất tĩnh cùng tên đuôi `.png`; không sinh HTML,
  không qua `mockup-tokens`; tên token trong Figma Variables trùng tên trong code.

## Đặt tên mockup

- Tên file `<F-xxx>-<B>.html` theo mã tính năng và mã bước; `state` ghi trạng thái:
  `normal`, `empty`, `error`, `loading`.
- Cùng mã bước nhiều trạng thái: trạng thái `normal` giữ tên `<F-xxx>-<B>.html`; trạng
  thái khác đặt `step` là mã bước kèm hậu tố chữ thường chưa dùng trong spec cho bước
  chèn (ví dụ `B3e` cho rỗng, `B3x` cho lỗi, `B3l` cho đang tải). Hậu tố này chỉ là tên
  file mockup, không phải bước mới: **không** ghi `B3e` vào cột "Mã bước" của user flow
  hay bảng hành vi (`userflow-steps` sẽ báo lỗi); liên kết file đó ở cột "Mockup" của
  dòng `B3` trong mục 2 user flow và mục 5 Feature Spec.

## Mã bước

- `B` + số, hậu tố chữ thường cho bước chèn giữa (`B2a`); không đánh số lại vì mockup và
  test đặt tên theo mã.
- Spec có sẵn: lấy mã từ mục 3 và 4 của spec. Spec chưa có: user flow đề xuất, Feature
  Spec dùng lại đúng mã.

## Khi nào nới quy tắc

- Tính năng dùng template trang có sẵn: mục Giao diện trong idea.md là đủ, không cần
  Design brief riêng; vẫn làm flow và mockup theo template.
- Công cụ nội bộ dùng thư viện UI có sẵn: bỏ wireframe, mockup chỉ cho màn hình phức tạp.
- Cá nhân, MVP: mockup cho màn hình phức tạp; wireframe và prototype để sau.
- Không có giao diện người dùng: bỏ toàn bộ.

## Luồng một chiều

- brief chưa duyệt thì không đi tiếp sang Design; bảng tác động CR chưa duyệt thì AI
  không sửa mockup.
- Sửa giao diện thì sửa mockup trước, rồi cập nhật spec và code. Sửa hành vi thì sửa
  Feature Spec trước, rồi cập nhật mockup.
- Mọi file liên kết ngược (`source`) về brief hoặc CR; `feature` trỏ Feature Spec.

## Tách file

- Một tính năng một user flow; một màn hình một wireframe; một trạng thái một mockup;
  một tính năng một UI spec. Thư mục `flows/`, `wireframes/`, `mockups/`.
- File vượt 500 dòng `dk check` cảnh báo; trần 800 (mockup dài là đang nhồi nhiều màn
  hình).

## Trạng thái và trường (khóa trong frontmatter)

- `status`: `draft`, `review` (chờ duyệt), `approved` (đã duyệt).
- `feature`: mã Feature Spec (`F-xxx`), bắt buộc ở userflow, wireframe, mockup, ui-spec;
  `step`: mã bước ở wireframe và mockup; `state`: trạng thái màn hình của mockup;
  `external`: liên kết công cụ ngoài, rỗng với mockup HTML; `tool`: công cụ sinh UI spec.
- `source`: thư mục intake của brief hoặc mã CR; `--from` spec hoặc flow đặt `source` là
  file đó, đổi về brief hoặc CR khi biết.
- `owner`: bắt buộc.

## Ghi file

- Mọi file trong `docs/` tạo bằng `dk new`. `dk` không chạy được thì dừng và báo người
  cài; không tạo file tay.
- Mỗi lần sửa có một dòng changelog với tóm tắt thật (`dk changelog add`), kể cả mockup
  `.html`.
- Không sửa `created`, `created_by`, `dk_version`; không sửa `tokens.css`.

## Quy tắc → mục báo cáo

| Quy tắc | Mục báo cáo |
|---|---|
| Nguyên tắc nền | 1, Tầng 3 (dòng Mockup, "Cấu trúc Design system" lớp 5) |
| Thứ tự trong tầng | Tầng 3 (bảng và đoạn "Thứ tự trong tầng") |
| Mockup HTML tĩnh | 6 (Wireframe, mockup, prototype; Lint giao diện), plan validate phiên 1 câu 8 |
| Mã bước | 3 (đoạn sơ đồ và mã bước) |
| Khi nào nới quy tắc | Tầng 1 (dòng template trang có sẵn), 9 (Cá nhân / MVP), 10 (công cụ nội bộ, không giao diện) |
| Luồng một chiều | 4 |
| Tách file | 7 |
| Trạng thái, ghi file | plan `plans/260903-1400-skill-tai-lieu/plan.md` mục 1.1, 1.3, 1.5, 1.10, 1.13; `assets/types.toml` |
