# Quy tắc doc-adr

Chép từ báo cáo `plans/reports/consult-260903-1032-tai-lieu-du-an-nguoi-va-ai.md`.
Chỉ quy tắc, không giải thích; giải thích đọc ở mục báo cáo ghi cuối file.

## Nguyên tắc nền

- Mỗi tài liệu tồn tại để trả lời một câu hỏi cụ thể. ADR trả lời "vì sao hệ thống làm theo cách này".
- ADR là tài liệu thẩm quyền bất biến: người chốt; AI soạn nháp từ thảo luận.
- Mỗi ADR có một chủ sở hữu (`owner`). Không có người chịu trách nhiệm thì không viết.
- AI không tự chốt ADR, không sửa thân ADR đã chốt.

## ADR

- Ghi bối cảnh, lựa chọn đã xét, quyết định, hệ quả của mỗi quyết định kỹ thuật quan trọng.
- Một quyết định một file trong `docs/adr/`, đánh số tăng dần `ADR-0001`, tên file mang mã.
- Cấu trúc: (1) frontmatter; (2) Bối cảnh; (3) Lựa chọn đã xét (bảng); (4) Quyết định; (5) Hệ quả; (6) Liên kết.
- Không sửa sau khi chốt. Quyết định mới thay thế thì ghi ADR mới với `supersedes`; ADR cũ chỉ đổi `status` và `superseded_by`.
- `dk check` quy tắc `adr-immutable`: thân ADR có `status` đã chốt ở HEAD mà khác HEAD là lỗi; frontmatter được đổi.
- Quyết định nhỏ chưa đến mức ADR (kể cả quyết định giao diện) ghi Decision log theo ngày, một dòng một quyết định; không tạo ADR.
- Chỉ mục ADR liệt kê trạng thái, sinh bằng `dk index adr`; không viết tay.

## Khi nào là ADR

- Chọn công nghệ, thư viện nền, cơ sở dữ liệu, giao thức, cách triển khai, ranh giới dịch vụ, quy ước ảnh hưởng nhiều tính năng.
- Thay đổi khó đảo ngược hoặc ảnh hưởng nhiều thành phần trong Architecture overview.
- CR có dòng ADR `Có` trong bảng tác động.
- Không phải ADR: đặt tên biến, chọn màu, sắp xếp một màn hình, sửa lỗi cục bộ.

## Luồng một chiều

- Nguồn ADR: brief đã duyệt, CR đã duyệt (dòng ADR trong bảng tác động), hoặc thảo luận trực tiếp với người.
- ADR chốt rồi thì Architecture overview và Feature Spec liên kết đến ADR, không chép nội dung.
- CR chạm ADR: người duyệt bảng tác động trước khi AI soạn ADR.
- ADR bị thay thế vẫn giữ file, làm lịch sử.

## Tách file và frontmatter

- Một quyết định một file, áp dụng từ ADR đầu tiên. Chi tiết từng thành phần của Architecture tách sang ADR khi Architecture vượt 400 dòng.
- Frontmatter thống nhất: mã, tiêu đề, trạng thái, chủ sở hữu, ngày cập nhật, nguồn, `supersedes`, `superseded_by`.
- ADR ngắn: dưới 100 dòng là thường; trần cứng 800 dòng.

## Trạng thái và trường (khóa trong frontmatter)

- `status`: `proposed` (đề xuất), `accepted` (đã chốt), `superseded` (thay thế), `deprecated` (bãi bỏ).
- Thân bất biến từ khi `accepted`; `superseded` và `deprecated` cũng bất biến.
- `supersedes`: mã ADR bị thay; `superseded_by`: mã ADR thay thế. Rỗng khi không có.
- `source`: mã CR hoặc thư mục brief nêu quyết định; rỗng khi từ thảo luận.

## Ghi file

- Mọi file trong `docs/` tạo bằng `dk new`. `dk` không chạy được thì dừng và báo người cài; không tạo file tay.
- Mỗi lần sửa ADR có một dòng changelog với tóm tắt thật (`dk changelog add`); nguồn là mã CR hoặc thư mục brief khi có.
- Chỉ mục ADR không ghi changelog (file sinh).
- Không sửa `created`, `created_by`, `dk_version` trong frontmatter.

## Quy tắc → mục báo cáo

| Quy tắc | Mục báo cáo |
|---|---|
| Nguyên tắc nền | 1, Tầng 2 (dòng ADR, cột Người, AI) |
| ADR, Decision log | 2, Tầng 2 và Tầng 5 (dòng Decision log); 5 (dòng "Vì sao", "Đã làm gì") |
| Khi nào là ADR, CR chạm ADR | 2, Tầng 1, 1b (bảng tác động, dòng ADR) |
| Luồng một chiều | 4 |
| Tách file và frontmatter | 7 (bảng "bắt buộc tách", dòng ADR; bảng "giữ một file", dòng Architecture overview) |
| Trạng thái và trường, ghi file | plan `plans/260903-1400-skill-tai-lieu/plan.md` mục 1.1, 1.3, 1.5; `assets/types.toml` |
