# Quy tắc doc-intake

Chép từ báo cáo `plans/reports/consult-260903-1032-tai-lieu-du-an-nguoi-va-ai.md`.
Chỉ quy tắc, không giải thích; giải thích đọc ở mục báo cáo ghi cuối file.

## Nguyên tắc nền

- Mỗi tài liệu tồn tại để trả lời một câu hỏi cụ thể. Không có câu hỏi thì không viết.
- idea.md và interview.md là tài liệu trạng thái: có ngày, không sửa lại sau khi chốt. brief.md là tài liệu thẩm quyền: luôn phải đúng, sửa tại chỗ.
- Tài liệu phải máy đọc được: cấu trúc cố định, Markdown trong repo, frontmatter.
- Mỗi tài liệu có một chủ sở hữu (`owner`). Không có người chịu trách nhiệm thì không viết.
- Intake brief là đầu vào (trước Feature Spec). Release brief là đầu ra (sau khi hoàn thành). Không nhầm hai loại.
- Ý tưởng mới (chưa có gì tồn tại) đi qua chuỗi ba bước idea, interview, brief. Thay đổi trên thứ đã có là Change Request, không đi qua chuỗi này.
- Template viết một lần; mỗi bản điền nằm trong thư mục riêng.
- AI không tự chốt tài liệu thẩm quyền. Brief do người duyệt.

## Chuỗi ba bước

- Ba file `idea.md`, `interview.md`, `brief.md` trong `intake/<yymmdd>-<slug>/`. Một ý tưởng một thư mục.
- idea.md: điền trong 10 phút; trường chưa biết ghi đúng chữ "chưa rõ". Người điền, AI không điền hộ.
- Template ý tưởng gồm: vấn đề, ai gặp, cách xoay xở hiện tại, kết quả mong muốn, điều không làm, ràng buộc, thước đo thành công, Giao diện (màn hình liên quan, template trang dùng, có màn hình hoặc organism mới không), câu hỏi mở.
- interview.md làm rõ từng trường "chưa rõ" hoặc mâu thuẫn trong idea.md. Câu hỏi sinh từ chỗ trống của template, không từ danh sách cố định.
- Hỏi từng câu một. Mỗi cặp hỏi đáp kết thúc bằng dòng "Quyết định rút ra" và trạng thái "đã chốt" hoặc "còn mở".
- Câu hỏi về giao diện phải hỏi bằng hình: AI đưa 2 đến 3 phương án tham chiếu để người chọn.
- interview.md có ngày, không sửa sau khi kết thúc.
- brief.md có bốn mục cố định: kết quả mong muốn, ràng buộc, ngoài phạm vi, tiêu chí chấp nhận; cộng "Quyết định đã chốt" liên kết về interview.md.
- brief.md thay thế idea.md làm nguồn sự thật. Là đầu vào cho Feature Spec và User flow.
- AI soạn nháp brief, không tự chốt. Người duyệt và chốt.

## Hai cấp

- Cấp dự án, một lần lúc khởi đầu: ý tưởng sản phẩm, brief viết thành Product overview; ý tưởng giao diện toàn sản phẩm, brief là Design brief, đầu vào cho Design system lớp 0 và templates.
- Design brief chạy sau hoặc song song ý tưởng sản phẩm, không chạy trước.
- Cấp tính năng, lặp lại mỗi tính năng mới: brief đổ vào User flow, Wireframe, Mockup rồi Feature Spec.

## Quy tắc tỷ lệ

- idea.md không còn trường "chưa rõ" và không mâu thuẫn: bỏ phỏng vấn, AI soạn thẳng brief.md, người duyệt.
- Có trường "chưa rõ" hoặc mâu thuẫn: phỏng vấn, chỉ hỏi đúng những trường đó.
- Tính năng dùng template trang có sẵn: mục Giao diện trong idea.md là đủ, không cần Design brief riêng.
- Dùng thư viện UI có sẵn và chấp nhận mặc định: bỏ chuỗi ý tưởng giao diện cấp dự án; Design brief rút xuống một trang ghi thư viện, ngoại lệ, giọng văn.
- Dự án đã chạy lâu, tài liệu tiếp nhận không còn ai đọc: Tầng 1 thành phụ lục quy trình.

## Luồng một chiều

- Mũi tên chỉ đi một chiều: idea.md, rồi interview.md (khi có trường trống), rồi brief.md (người duyệt).
- brief.md chưa được người duyệt thì không đi tiếp sang Design hay Feature Spec.
- Không sửa test hay brief mà bỏ qua spec.
- Tài liệu đích (Feature Spec, mockup, ADR) phải liên kết ngược về brief sinh ra nó. Brief chốt rồi không đổ đi đâu là tài liệu chết.

## Tách file và frontmatter

- Một ý tưởng một thư mục, ba file, áp dụng từ ý tưởng đầu tiên. Template gốc một file duy nhất.
- Mỗi thư mục có chỉ mục sinh tự động (`dk index intake`). AI đọc chỉ mục trước rồi mở đúng file.
- Frontmatter thống nhất: mã, tên, trạng thái, chủ sở hữu, ngày cập nhật.
- Không tách theo mục con của một đơn vị.
- File đơn vượt 300 đến 500 dòng thì xem lại; trần cứng 800 dòng.

## Trạng thái (khóa trong frontmatter)

- idea: `draft` (mới ghi), `interviewing` (đang phỏng vấn), `briefed` (đã có brief), `dropped` (bỏ).
- interview: `open`, `closed`.
- brief: `draft` (đang soạn), `review` (chờ duyệt), `approved` (đã duyệt, người ghi `approved_by`, `approved_at`), `rejected`.
- brief: `level` là `project` hoặc `feature`; `kind` là `product`, `design` hoặc `feature`.

## Ghi file

- Mọi file trong `docs/` tạo bằng `dk new`. `dk` không chạy được thì dừng và báo người cài; không tạo file tay.
- Mỗi file tạo hoặc sửa có một dòng changelog với tóm tắt thật (`dk changelog add`), nguồn là thư mục intake.
- Không sửa `created`, `created_by`, `dk_version` trong frontmatter.

## Quy tắc → mục báo cáo

| Quy tắc | Mục báo cáo |
|---|---|
| Nguyên tắc nền | 1 |
| Chuỗi ba bước, hai cấp, quy tắc tỷ lệ | 2, Tầng 1, 1a |
| Dùng thư viện UI có sẵn; dự án chạy lâu | 10 |
| Luồng một chiều | 4 (đoạn "Mũi tên chỉ đi một chiều") |
| Tách file và frontmatter | 7 (dòng "Ý tưởng mới", "Quy tắc đi kèm") |
| Trạng thái, ghi file | plan `plans/260903-1400-skill-tai-lieu/plan.md` mục 1.1, 1.3, 1.5; `assets/types.toml` |
