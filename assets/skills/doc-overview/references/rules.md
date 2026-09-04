# Quy tắc doc-overview

Chép từ báo cáo `plans/reports/consult-260903-1032-tai-lieu-du-an-nguoi-va-ai.md`.
Chỉ quy tắc, không giải thích; giải thích đọc ở mục báo cáo ghi cuối file.

## Nguyên tắc nền

- Mỗi tài liệu tồn tại để trả lời một câu hỏi cụ thể. Không có câu hỏi thì không viết.
- Tài liệu thẩm quyền (Product overview, Architecture overview, Glossary) do người viết và chốt; AI đề xuất nháp, phát hiện lệch với mã, phát hiện thuật ngữ chưa định nghĩa.
- Tài liệu sinh (Feature catalog) không viết tay; sinh từ frontmatter của Feature Spec.
- Mỗi tài liệu có một chủ sở hữu (`owner`). Không có người chịu trách nhiệm thì không viết.
- AI không tự sửa tài liệu thẩm quyền đã chốt nếu chưa có CR hoặc brief mới đã duyệt.

## Product overview

- Trả lời: sản phẩm giải quyết vấn đề gì, cho ai, nhóm tính năng chính, thước đo thành công. Đọc đầu tiên khi tham gia dự án.
- Là Product brief cấp dự án viết thành trang: chỉ soạn từ brief `level: project`, `kind: product`, `status: approved`.
- Một trang, dưới 150 dòng (`dk check` cảnh báo). Dài hơn hai trang là đang lẫn Feature Spec vào; không tách thành nhiều file.
- Bốn mục của brief chép sang: Kết quả mong muốn, Ràng buộc, Ngoài phạm vi, Tiêu chí chấp nhận (thành Thước đo thành công). Frontmatter `outcome`, `acceptance`, `source` do `dk new --from` chép.
- Nhóm tính năng: mỗi nhóm một dòng; chi tiết ở Feature catalog, không chép vào đây.

## Feature catalog

- Bảng điều khiển trả lời "phần mềm có gì": mã, tên, trạng thái, chủ sở hữu, liên kết Feature Spec, nguồn.
- Sinh tự động bằng `dk index features`, ghi `docs/features/README.md` với `generated: true`. Không viết tay, không sửa tay.
- Cập nhật trạng thái bằng cách sửa frontmatter Feature Spec rồi chạy lại chỉ mục.
- Dự án dưới 5 tính năng hoặc dưới 2 tháng: bỏ được; vẫn mỗi tính năng một file.

## Architecture overview

- Trả lời: sơ đồ thành phần, luồng dữ liệu, ranh giới hệ thống, tech stack.
- Sơ đồ trước, văn xuôi sau. Mỗi nút sơ đồ một thành phần chạy độc lập; không vẽ bên trong thành phần.
- Soạn từ mã nguồn thật khi có: cấu trúc thư mục, file cấu hình, manifest phụ thuộc, file triển khai. Thành phần suy đoán đánh dấu `chưa xác nhận`; người xác nhận từng dòng; AI không tự chốt.
- Cập nhật khi thêm hoặc bỏ thành phần. Chi tiết từng thành phần tách sang ADR hoặc file riêng khi vượt 400 dòng (`dk check` cảnh báo).
- Quyết định kỹ thuật trong Architecture chỉ liên kết đến ADR, không chép nội dung.

## Glossary

- Thuật ngữ nghiệp vụ và định nghĩa thống nhất. Bắt buộc khi có AI: dùng sai thuật ngữ sinh mã sai.
- Một bảng, một file, sắp theo bảng chữ cái: `Thuật ngữ | Định nghĩa | Không nhầm với | Nguồn`. Vượt 300 thuật ngữ thì tách theo miền nghiệp vụ.
- Nguồn thuật ngữ: chữ in đậm lần đầu trong mục 2, 5, 8 của Feature Spec và trong brief. AI quét và đề xuất định nghĩa; người chọn giữ và chốt.
- `dk check` quy tắc `glossary-term`: thuật ngữ in đậm trong mục 2, 5, 8 của Feature Spec chưa có ở cột đầu bảng là cảnh báo.
- Không có giao diện người dùng vẫn giữ giọng văn cho thông báo lỗi trong Glossary.

## Luồng một chiều

- Brief cấp dự án (người duyệt) đổ vào Product overview; Product overview đổ vào Feature catalog; Feature Spec sinh từ brief hoặc CR, không từ Product overview.
- Sửa Product overview hoặc Architecture sau khi chốt chỉ khi có CR hoặc brief mới đã duyệt trỏ đến; ghi `source`.
- Tài liệu cấp dự án không chứa hành vi từng tính năng; việc đó là Feature Spec.

## Tách file và frontmatter

- Product overview, Architecture overview, Glossary: mỗi thứ một file trong `docs/overview/`, không tách theo mục.
- Feature catalog là `docs/features/README.md` sinh ra; chỉ mục ADR là `docs/adr/README.md` sinh ra.
- Frontmatter thống nhất: loại, tiêu đề, trạng thái, chủ sở hữu, ngày cập nhật, nguồn.
- File đơn vượt 300 đến 500 dòng thì xem lại; trần cứng 800 dòng.

## Trạng thái và trường (khóa trong frontmatter)

- `status` (cả ba loại): `draft` (nháp), `review` (chờ duyệt), `approved` (đã chốt).
- Product overview: `source` là brief nguồn (bắt buộc); `outcome`, `acceptance` chép từ brief.
- Architecture, Glossary: `source` rỗng khi soạn từ mã hoặc từ Feature Spec; là mã CR khi sửa theo CR.
- Ngưỡng dòng riêng: Product overview 150, Glossary 300, Architecture 400 (`warn_lines` trong `types.toml`).

## Ghi file

- Mọi file trong `docs/` tạo bằng `dk new`. `dk` không chạy được thì dừng và báo người cài; không tạo file tay.
- Mỗi lần sửa có một dòng changelog với tóm tắt thật (`dk changelog add`); nguồn là brief hoặc CR khi có, `-` khi soạn từ mã.
- Feature catalog và chỉ mục ADR không ghi changelog (file sinh).
- Không sửa `created`, `created_by`, `dk_version` trong frontmatter.

## Quy tắc → mục báo cáo

| Quy tắc | Mục báo cáo |
|---|---|
| Nguyên tắc nền | 1, Tầng 2 (cột Người, AI) |
| Product overview, Feature catalog, Architecture, Glossary | 2, Tầng 2; 5 |
| Bỏ Feature catalog khi dự án nhỏ; giọng văn khi không có giao diện | 10 |
| Luồng một chiều | 4 |
| Tách file và frontmatter | 7 (bảng "giữ một file", "Quy tắc đi kèm") |
| Trạng thái và trường, ghi file | plan `plans/260903-1400-skill-tai-lieu/plan.md` mục 1.1, 1.3, 1.5; `assets/types.toml` |
