# Quy tắc doc-ops

Chép từ báo cáo `plans/reports/consult-260903-1032-tai-lieu-du-an-nguoi-va-ai.md`.
Chỉ quy tắc, không giải thích; giải thích đọc ở mục báo cáo ghi cuối file.

## Nguyên tắc nền

- Tài liệu vận hành là tài liệu thẩm quyền (luôn đúng với production, sửa tại chỗ),
  trừ Postmortem là tài liệu trạng thái (một sự cố một file, không sửa sau `closed`).
- Soạn từ cấu hình và log thật trong repo hoặc do người cung cấp; không bịa lệnh,
  URL, ngưỡng.
- Mỗi tài liệu có một `owner`. Không có người chịu trách nhiệm thì không viết.
- Một sự cố một file (Runbook, Postmortem) từ file đầu tiên; không gộp rồi tách.

## Tầng 8: Vận hành

| Tài liệu | Cách dùng | Loại |
|---|---|---|
| Deployment guide | chạy được từng bước, không giả định kiến thức ngầm, kèm rollback | thẩm quyền |
| Environment | tên biến, ý nghĩa, khác biệt dev/staging/prod; không giá trị secret | thẩm quyền |
| Runbook | một sự cố một file; viết sau mỗi incident | thẩm quyền |
| Monitoring | metric, ngưỡng, người nhận; liên kết dashboard thật | thẩm quyền |
| Postmortem | nguyên nhân gốc, timeline, hành động; không đổ lỗi; trong 48 giờ | trạng thái |
| Backup và DR | RPO, RTO, quy trình khôi phục; diễn tập ít nhất một lần | thẩm quyền |

## Deployment guide

- Một file `docs/ops/deployment.md`. Nguồn: cấu hình CI, Dockerfile, script deploy.
- Mỗi bước một lệnh trong khối ``` hoặc một hành động quan sát được, kèm kết quả
  mong đợi. Người mới làm theo được không cần hỏi.
- Có mục rollback đủ bước và điều kiện quyết định rollback; có bước xác nhận sau
  triển khai.

## Environment

- Một file `docs/ops/environment.md`. Chỉ tên biến, ý nghĩa, môi trường, nơi giữ
  giá trị thật, có phải secret.
- Không giá trị secret. `dk check env-no-secret` (lỗi): dòng `KEY=value` với
  value không rỗng và không phải placeholder `<...>`. Rule là lớp phụ, không thay
  công cụ quét secret của dự án; secret không theo mẫu `KEY=value` không bắt.
- Không mở `.env` thật hay secret store để soạn; người dán giá trị thật thì không
  chép và báo xoay vòng.

## Runbook

- Một sự cố một file trong `docs/ops/runbooks/`; người trực mở đúng file, không
  cuộn tìm. Viết sau mỗi incident hoặc cho sự cố lường trước.
- Thứ tự cố định: dấu hiệu, bước chẩn đoán, bước xử lý, bước xác nhận, leo thang.
  Mỗi bước một lệnh hoặc một hành động quan sát được kèm kết quả nếu đúng.
- Bước có rủi ro ghi rõ cần ai xác nhận trước. Không giải thích dài; giải thích ở
  postmortem.
- Trỏ tên cảnh báo trong `monitoring.md` và postmortem liên quan.

## Monitoring

- Một file `docs/ops/monitoring.md`. Mỗi metric: ý nghĩa, ngưỡng cảnh báo, ngưỡng
  nghiêm trọng, người nhận, runbook.
- Liên kết dashboard thật bằng URL; không mô tả dashboard bằng chữ.
- Người chốt ngưỡng và người nhận; skill không tự đặt.

## Postmortem

- Một sự cố một file `docs/ops/postmortems/<yymmdd>-<slug>.md`; `incident_at`
  là giờ bắt đầu sự cố; `written_within_48h` do `dk new` tính, `false` thì báo
  người ngay và vẫn viết, ghi lý do ở mục Bài học.
- Không đổ lỗi cá nhân: nguyên nhân nói về hệ thống và quy trình, không tên người
  gây lỗi. Timeline từ log, alert, chat, không suy đoán.
- Phân biệt nguyên nhân kích hoạt và nguyên nhân gốc; ghi vì sao không phát hiện sớm.
- Mỗi hành động khắc phục có loại (ngăn tái diễn, phát hiện sớm, giảm ảnh hưởng),
  người và hạn do người chốt. `closed` khi mọi hành động xong; sau đó không sửa.
- Sự cố chưa có runbook thì tạo runbook sau postmortem.

## Backup và DR

- Một file `docs/ops/backup-dr.md`. RPO, RTO do người chốt theo hệ thống.
- Quy trình khôi phục từng bước như runbook; mỗi kịch bản thảm họa trỏ về bước
  khôi phục.
- Phải diễn tập ít nhất một lần; mỗi lần ghi mục Lịch sử diễn tập và cập nhật
  `last_drill: yyyy-mm-dd`. `dk status` báo file quá 6 tháng chưa diễn tập; skill
  nhắc người, không tự đổi ngày.

## Tách file và ngưỡng dòng

- Deployment, Environment, Monitoring, Backup và DR: mỗi thứ một file. Runbook và
  Postmortem: một sự cố một file.
- Vượt 300 đến 500 dòng thì xem lại; trần cứng 800 dòng (`line-threshold`).
- Frontmatter thống nhất: `type`, `title`, `status`, `owner`, `created`,
  `updated`, `source`; postmortem thêm `incident_at`, `written_within_48h`;
  backup-dr thêm `last_drill`. `dk new` điền, skill không sửa `created`,
  `created_by`, `dk_version`.

## Khi nào không áp dụng

- Chưa lên production: chưa cần Tầng 8; Deployment guide có thể viết sớm khi đã có CI.
- Nền tảng quản lý hộ (PaaS không cần vận hành): Deployment và Environment ngắn,
  Runbook chỉ cho sự cố phía ứng dụng.
- Ngành có quy chuẩn riêng (y tế, tài chính, chính phủ): mẫu postmortem và DR do quy
  chuẩn quyết định.

## Trạng thái (khóa tiếng Anh trong `types.toml`)

| Khóa | Nghĩa |
|---|---|
| `draft` | đang soạn |
| `review` | chờ người duyệt hoặc chạy thử |
| `approved` | đã chốt, đúng với production |
| `closed` | postmortem xong mọi hành động, không sửa |

Giải thích: báo cáo mục 1 (một đơn vị một file), mục 2 (Tầng 8), mục 7 (Runbook,
Postmortem), mục 10 (quy chuẩn ngành).
