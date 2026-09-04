# Quy tắc doc-design-system

Chép từ báo cáo `plans/reports/consult-260903-1032-tai-lieu-du-an-nguoi-va-ai.md`.
Chỉ quy tắc, không giải thích; giải thích đọc ở mục báo cáo ghi cuối file.

## Nguyên tắc nền

- Design system là cấp dự án, viết một lần, là tài liệu thẩm quyền. Người chốt tokens,
  pattern, giọng văn; AI đề xuất và kiểm tra lệch tokens.
- AI không tự sửa Design system đã duyệt nếu chưa có người duyệt (qua CR).
- Nhất quán đến từ việc mọi giá trị hiển thị được đặt tên và tham chiếu, không có giá
  trị gõ tay. Tài liệu không giữ được nhất quán, lint mới giữ được: `dk check` quy tắc
  `mockup-tokens` là lint cho mockup do `dk` quản; lint code UI của dự án đích là việc
  của dự án đích, skill không làm.
- Chỉ áp dụng khi phần mềm có giao diện. Không có giao diện: bỏ toàn bộ, giữ giọng văn
  thông báo lỗi trong Glossary.

## Cấu trúc theo lớp (Atomic Design)

| Lớp | Là gì | Quy tắc |
|---|---|---|
| 0 Primitive tokens | thang thô: màu, chữ, khoảng cách, kích thước, bo góc, bóng, viền, chuyển động, điểm gãy, z-index | mỗi thang dưới 10 bước; chỉ semantic được tham chiếu primitive |
| 0 Semantic tokens | giá trị gắn ý nghĩa: màu theo vai trò và trạng thái, chữ theo vai trò, khoảng cách theo vai trò, độ nổi | lớp duy nhất component được tham chiếu; theme sáng tối và đa thương hiệu đổi ở đây |
| 0 Quy tắc nền | lưới, điểm gãy, accessibility, một bộ icon, giọng văn | áp cho mọi lớp trên |
| 1 Atoms | Text, Icon, Button, Input, Checkbox, Radio, Switch, Badge, Avatar, Spinner, Divider | không chứa atom khác |
| 2 Molecules | Form field, Select, Dropdown menu, Tooltip, Toast, Tag, Tabs, Pagination | chỉ chứa atom |
| 3 Organisms | Card, Modal, Drawer, Table, Navigation, Empty state, Form | chứa molecule và atom |
| 4 Templates | trang danh sách, chi tiết, form, dashboard, xác thực | chứa organism, định nghĩa vùng đặt |
| 5 Pages | template đổ dữ liệu thật của một tính năng | không định nghĩa style mới; là mockup của họ Design flow |

Patterns không phải lớp: công thức ghép organism cho một tình huống (form dài, xác nhận
xóa, tìm kiếm và lọc bảng, trạng thái tải, rỗng, lỗi hệ thống, luồng nhiều bước); đi kèm
lớp 3, giữ nhất quán về hành vi chứ không chỉ hình ảnh.

Ba quy tắc nhất quán:

- Mỗi lớp chỉ tham chiếu lớp ngay dưới: organism không gọi thẳng token, page không định
  nghĩa style, component không tham chiếu primitive.
- Tên token trong công cụ thiết kế (Figma Variables) và trong code giống hệt nhau.
- Cấm giá trị gõ tay và cấm import vượt lớp.

## Mỗi component đủ 5 mục

1. Biến thể: primary, secondary, ghost, destructive.
2. Kích thước: sm, md, lg theo thang chung.
3. Trạng thái: mặc định, hover, active, focus, disabled, loading, lỗi.
4. Quy tắc dùng và không dùng.
5. Accessibility: vai trò ARIA, điều hướng bàn phím.

## Bộ tối thiểu, đúng thứ tự

1. Lớp 0: thang khoảng cách và chữ, bảng màu và semantic cho nền, chữ, hành động, viền,
   bốn trạng thái, bo góc, bóng, icon, accessibility cơ bản (`tokens.json`, `dk tokens
   css`, `foundations.md`).
2. Lớp 1: Text, Icon, Button, Input.
3. Lớp 2: Form field.
4. Lớp 3: Form, Empty state, Toast.
5. Lớp 4: một template trang danh sách và một template trang form.

Kiểm tra thật: AI sinh một page từ template có sẵn, chỉ dùng organism và semantic
tokens, `dk check` qua. Qua được là hệ đã nhất quán.

## Design tokens và `dk tokens css`

- Một file `tokens.json` khung W3C Design Tokens là nguồn sự thật: nhóm lồng, `$type`
  kế thừa từ nhóm, `$value`, alias `{a.b.c}` (kể cả lồng trong chuỗi như bóng). Metadata
  của dk ở khóa `$dk`.
- `dk tokens css` sinh `tokens.css`: `:root { --a-b-c: ... }` theo tên đường dẫn nối
  gạch ngang, alias giải về primitive, theme qua `$extensions.dk.theme.<tên>` thành khối
  `[data-theme="<tên>"]`. Chuẩn hóa hai loại color và dimension; loại khác chép nguyên
  chuỗi hoặc số. Alias vòng hoặc alias lạ là lỗi rõ. Giá trị theme khai trên chính token
  semantic cần đổi (token alias đến một token có theme không tự đổi theo).
- Ngoài phạm vi, không hứa: Tailwind config, theme object, Style Dictionary, đồng bộ
  Figma, sinh tài liệu component từ Storybook.
- `tokens.css` là bản sinh: không sửa tay, không ghi changelog; mockup tham chiếu qua
  `../tokens/tokens.css`.

## Thư viện UI có sẵn và các ngoại lệ

- Dùng thư viện UI có sẵn (MUI, Ant, shadcn): lớp 0 primitive và lớp 1 đến 3 đã có;
  chỉ viết semantic ánh xạ sang thư viện, quy tắc dùng, patterns, templates và giọng
  văn. Không xây lại. Design brief rút xuống một trang ghi thư viện, ngoại lệ, giọng văn.
- Công cụ nội bộ dùng thư viện có sẵn: Design system rút gọn thành "dùng thư viện X,
  các ngoại lệ sau".
- Không có giao diện người dùng: bỏ toàn bộ họ Design.

## Tách file

- Tokens một file JSON. Mỗi component một file theo lớp `atoms/`, `molecules/`,
  `organisms/`, `templates/`; mỗi pattern một file trong `patterns/`. Component sống
  trong code; tài liệu ở đây là hợp đồng, không chép code.
- File vượt 500 dòng `dk check` cảnh báo; trần 800.

## Luồng một chiều

- Design brief `approved`, `kind: design`, `level: project` mới đi tiếp sang tokens và
  foundations; brief chưa duyệt thì không.
- Sửa giao diện thì sửa mockup trước rồi spec; sửa token hoặc component đã duyệt đi qua
  CR có dòng Design là Có, người duyệt bảng tác động trước.
- Tokens và foundations `source` trỏ về Design brief hoặc CR sửa gần nhất.

## Trạng thái và trường (khóa trong frontmatter)

- `status`: `draft`, `review` (chờ duyệt), `approved` (đã chốt); component và pattern
  thêm `deprecated` (bãi bỏ, thay bằng component khác).
- `design-component`: `layer` là `atom`, `molecule`, `organism`, `template`; `uses` liệt
  kê component lớp ngay dưới.
- `design-foundations`: `ui_library` tên thư viện UI dùng sẵn, rỗng khi tự xây.
- `owner`: bắt buộc cho mọi file.

## Ghi file

- Mọi file trong `docs/` tạo bằng `dk new`. `dk` không chạy được thì dừng và báo người
  cài; không tạo file tay.
- Mỗi lần sửa có một dòng changelog với tóm tắt thật (`dk changelog add`); `--source` là
  thư mục intake của Design brief hoặc mã CR.
- Không sửa `created`, `created_by`, `dk_version`.

## Quy tắc → mục báo cáo

| Quy tắc | Mục báo cáo |
|---|---|
| Nguyên tắc nền | 1, Tầng 3 (dòng Design system) |
| Cấu trúc theo lớp, 5 mục, bộ tối thiểu, ba quy tắc | Tầng 3, "Cấu trúc Design system (Atomic Design)" |
| Design tokens và lint | 6 (Design tokens, Lint giao diện) |
| Thư viện UI có sẵn và ngoại lệ | 10 (thư viện UI, công cụ nội bộ, không giao diện), "Cấu trúc Design system" đoạn cuối |
| Tách file | 7 (dòng Design system) |
| Luồng một chiều | 4, Tầng 1 (dòng Tác động chạm Design system) |
| Trạng thái, ghi file | plan `plans/260903-1400-skill-tai-lieu/plan.md` mục 1.1, 1.3, 1.5, 1.10, 1.13; `assets/types.toml` |
