# Quy tắc doc-cr

Chép từ báo cáo `plans/reports/consult-260903-1032-tai-lieu-du-an-nguoi-va-ai.md`.
Chỉ quy tắc, không giải thích; giải thích đọc ở mục báo cáo ghi cuối file.

## Nguyên tắc nền

- Mỗi tài liệu tồn tại để trả lời một câu hỏi cụ thể. Không có câu hỏi thì không viết.
- Một nguồn sự thật cho mỗi tính năng (Feature Spec). Các tài liệu khác sinh ra hoặc tham chiếu từ đó.
- CR là tài liệu trạng thái; phần quyết định là thẩm quyền cho đến khi tài liệu đích cập nhật xong. Không sửa phần quyết định sau khi chốt.
- Mỗi tài liệu có một chủ sở hữu (`owner`). Không có người chịu trách nhiệm thì không viết.
- Change Request là tài liệu riêng cho thay đổi trên thứ đã có, có phân tích tác động và tự rẽ sang Feature Spec, Design, ADR hay Test tùy tác động. Ý tưởng mới không đi qua CR.
- AI không tự sửa tài liệu thẩm quyền, mục hành vi của Feature Spec, Design system hoặc mockup đã duyệt nếu chưa có người duyệt.

## Change Request

- CR khác ý tưởng mới ở ba điểm: có trạng thái hiện tại để so sánh, có phân tích tác động lên tài liệu và code đang tồn tại, đích đến không cố định mà tùy tác động.
- Một CR một file `cr/<CR-id>.md`. Tách `interview.md` ra thư mục cùng tên CR chỉ khi hỏi đáp dài.
- Template CR một file gốc, dưới một trang.
- Cấu trúc CR: (1) frontmatter: mã, tiêu đề, người yêu cầu, ngày, trạng thái, mức ưu tiên; (2) Hiện trạng: liên kết đến Feature Spec, mockup, ADR đang mô tả thứ sắp đổi, không chép lại nội dung; (3) Thay đổi đề xuất và lý do; (4) Phân tích tác động; (5) Hỏi đáp; (6) Quyết định; (7) Tài liệu đích đã cập nhật.
- Bảng tác động do AI soạn nháp từ tài liệu hiện có, người duyệt. Sáu dòng: Feature Spec (sửa mã bước nào, thêm ngoại lệ nào); Design: mockup, Design system (sửa màn hình nào, cần organism hay template mới không); ADR (quyết định nào bị thay thế); Data model, API (migration, contract đổi); Test case (test nào sửa, thêm); Release brief, User guide (đoạn nào cập nhật).
- Hỏi đáp chỉ khi thay đổi hoặc tác động còn "chưa rõ". Cùng quy tắc với interview: từng câu, kết bằng quyết định rút ra; câu về giao diện đưa 2 đến 3 phương án tham chiếu.
- Quyết định: chấp nhận, từ chối, hoãn. Ghi lý do và người chốt.
- Mục 7 liên kết đến từng tài liệu đích đã sửa, đánh dấu khi xong.
- Trạng thái CR: mới, đang phân tích, chờ duyệt, đã chốt, đang triển khai, đã đóng, từ chối. Chỉ mục CR sinh từ frontmatter trả lời "đang có thay đổi nào".
- Người điền yêu cầu, trả lời, duyệt tác động và quyết định. AI soạn phân tích tác động, dẫn hỏi đáp, cập nhật tài liệu đích sau khi duyệt.

## Quy tắc tỷ lệ

- Sửa lỗi nhỏ, tác động chỉ một Feature Spec và test: bảng tác động rút còn hai dòng, bỏ hỏi đáp.
- Tác động chạm Design system hoặc ADR: bắt buộc người duyệt bảng tác động trước khi AI sửa bất kỳ tài liệu đích nào.
- CR bị từ chối vẫn giữ file, làm lịch sử vì sao không đổi.

## Luồng một chiều

- CR đi: hiện trạng, thay đổi, bảng tác động (người duyệt), quyết định, rồi mới rẽ sang Feature Spec, Mockup và Design system, ADR, Data và API, Test case, Release brief.
- Sửa hành vi thì sửa Feature Spec trước, rồi cập nhật tài liệu con. Sửa giao diện thì sửa mockup trước, rồi cập nhật spec và code.
- Không sửa test hay brief mà bỏ qua spec.
- Bảng tác động chưa được người duyệt thì AI không sửa bất kỳ tài liệu đích nào.
- Tài liệu đích phải liên kết ngược (`source`) về CR sửa nó. CR chốt rồi không đổ đi đâu là tài liệu chết.

## Tách file và frontmatter

- Một CR một file, áp dụng từ CR đầu tiên. Chỉ mục CR sinh từ frontmatter trạng thái (`dk index cr`).
- AI đọc chỉ mục trước rồi mở đúng file.
- Frontmatter thống nhất: mã, tên, trạng thái, chủ sở hữu, ngày cập nhật. Tên file mang mã để liên kết chéo ổn định.
- File đơn vượt 300 đến 500 dòng thì xem lại; trần cứng 800 dòng.

## Trạng thái và trường (khóa trong frontmatter)

- `status`: `new` (mới), `analyzing` (đang phân tích), `review` (chờ duyệt), `approved` (đã chốt), `in-progress` (đang triển khai), `closed` (đã đóng), `rejected` (từ chối).
- `requester`: người yêu cầu. `priority`: `low`, `medium`, `high`. `owner`: người chịu trách nhiệm CR.
- interview của CR: `source` là mã CR; `status` `open` hoặc `closed`.

## Ghi file

- Mọi file trong `docs/` tạo bằng `dk new`. `dk` không chạy được thì dừng và báo người cài; không tạo file tay.
- Mỗi lần sửa CR có một dòng changelog với tóm tắt thật (`dk changelog add`), nguồn là mã CR.
- Tài liệu đích cập nhật theo CR ghi changelog với `--source <CR-id>`.
- Không sửa `created`, `created_by`, `dk_version` trong frontmatter.

## Quy tắc → mục báo cáo

| Quy tắc | Mục báo cáo |
|---|---|
| Nguyên tắc nền | 1 |
| Change Request, quy tắc tỷ lệ | 2, Tầng 1, 1b |
| CR sửa lỗi nhỏ | 10 |
| Luồng một chiều | 4 (đoạn "Mũi tên chỉ đi một chiều") |
| Tách file và frontmatter | 7 (dòng "Change Request", "Quy tắc đi kèm") |
| Trạng thái và trường, ghi file | plan `plans/260903-1400-skill-tai-lieu/plan.md` mục 1.1, 1.3, 1.5; `assets/types.toml` |
