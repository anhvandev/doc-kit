---
name: doc-ops
description: 'Tài liệu vận hành production: Deployment guide, Environment không secret, Runbook một sự cố một file, Monitoring, Postmortem trong 48 giờ, Backup và DR. Không dùng cho tài liệu người dùng cuối.'
---

# doc-ops

Soạn Deployment guide, Environment, Runbook, Monitoring, Postmortem, Backup và DR
trong `docs/ops/` bằng `dk new` và `dk check`; từng bước chạy được.

**Dừng bắt buộc: không ghi giá trị secret vào tài liệu dù người đưa; không bịa
lệnh hoặc URL dashboard chưa kiểm được trong repo hay cấu hình thật; Postmortem
không nêu tên người gây lỗi; người chốt RPO, RTO, ngưỡng cảnh báo và hành động
khắc phục.**

Mọi file trong `docs/` chỉ được tạo bằng `dk new`. Không tự tạo file tay, kể cả khi
`dk` lỗi. Quy tắc nội dung ở `references/rules.md`; đọc trước khi làm.

## Phạm vi

- Làm: `docs/ops/deployment.md`, `environment.md`, `monitoring.md`,
  `backup-dr.md`; `docs/ops/runbooks/<slug>.md`; `docs/ops/postmortems/<yymmdd>-<slug>.md`.
- Không làm: Release brief, User guide (họ Người dùng cuối); Architecture overview,
  ADR (họ Nền sản phẩm); script deploy hay cấu hình CI (mã của dự án đích).
- Dự án chưa lên production: báo người và chỉ làm khi người vẫn muốn.

## Bước 0. Kiểm `dk`

```
dk --version
```

Lỗi hoặc không tìm thấy lệnh: dừng, báo "Cần cài `dk` trước (xem README của dk)".

## `new`: Deployment guide

1. Đọc cấu hình thật: `.github/workflows/`, `.gitlab-ci.yml`, `Dockerfile`,
   `docker-compose*.yml`, script deploy, `Makefile`. Không có thì **hỏi người**
   cách đang triển khai; không bịa.
2. `dk new deployment deployment --set owner="<người vận hành>"`.
3. Mục 3 và 5: mỗi bước một lệnh trong khối ``` chép từ cấu hình thật, kèm kết quả
   mong đợi; rollback đủ bước quay về phiên bản trước. Không giả định kiến thức ngầm.
4. Chạy khối "Sau mỗi lần sửa". `status: review`. **Dừng.** Người chạy thử một lần
   theo tài liệu rồi `approved`.

## `new`: Environment

1. Đọc `.env.example`, cấu hình CI, code đọc biến môi trường; **không mở** file
   `.env` thật hay secret store.
2. `dk new environment environment --set owner="<người vận hành>"`.
3. Mỗi biến một dòng bảng: ý nghĩa, khác biệt dev/staging/prod, có phải secret.
   Mục 3 chỉ `KEY=<mô tả nơi lấy>`. Người dán giá trị thật vào chat thì không chép,
   báo người xoay vòng secret đó.
4. Chạy khối "Sau mỗi lần sửa". `env-no-secret` báo lỗi dòng nào thì thay bằng
   placeholder `<...>`. Nhắc người: rule là lớp phụ, không thay công cụ quét secret.

## `new`: Runbook

Sau mỗi sự cố hoặc khi người nêu một sự cố có thể xảy ra:

```
dk new runbook <slug> --set owner="<người trực>" --set title="<tên sự cố>"
```

Một sự cố một file. Ba khối bước theo thứ tự: chẩn đoán (xác nhận đúng sự cố),
xử lý, xác nhận (metric về ngưỡng, thử thao tác). Mỗi bước một lệnh hoặc một hành
động quan sát được kèm kết quả nếu đúng; lệnh lấy từ postmortem, log, chat sự cố
thật. Mục 1 trỏ tên cảnh báo trong `monitoring.md`. Chạy khối "Sau mỗi lần sửa".
**Dừng.** Người chạy thử từng bước rồi `approved`.

## `new`: Monitoring

1. Hỏi người URL dashboard và kênh cảnh báo thật; đọc cấu hình alert trong repo nếu có.
2. `dk new monitoring monitoring --set owner="<người vận hành>"`.
3. Mỗi metric một dòng: ngưỡng, người nhận, runbook tương ứng (liên kết file trong
   `runbooks/`). Người chốt ngưỡng; không tự đặt.
4. Chạy khối "Sau mỗi lần sửa". **Dừng.** Người duyệt.

## `new`: Postmortem

Ngay khi sự cố kết thúc:

```
dk new postmortem <slug> --set owner="<người điều phối>" --set incident_at="<yyyy-mm-dd hh:mm>"
```

`dk` tính `written_within_48h`. Giá trị `false`: **báo người ngay** là đã quá 48
giờ (hoặc `incident_at` sai định dạng), vẫn viết và ghi lý do ở mục 6.

- Timeline từ log, alert, chat; không suy đoán giờ.
- Nguyên nhân gốc nói về hệ thống và quy trình; không nêu tên người gây lỗi.
- Mỗi hành động khắc phục có người và hạn do người chốt; để trống chờ người, không
  tự điền tên.
- Sự cố chưa có runbook: sau postmortem tạo runbook theo mục ở trên.

Chạy khối "Sau mỗi lần sửa". `status: review`. **Dừng.** Người duyệt; `closed`
khi mọi hành động xong.

## `new` và `update`: Backup và DR

```
dk new backup-dr backup-dr --set owner="<người vận hành>"
```

Người chốt RPO, RTO; skill điền quy trình khôi phục từng bước từ cấu hình backup
thật. Sau mỗi lần diễn tập: ghi dòng mục 5 và cập nhật `last_drill: <yyyy-mm-dd>`.
`dk status` báo "DR chưa diễn tập quá 6 tháng": nhắc người lên lịch diễn tập,
không tự đổi `last_drill`.

## Sau mỗi lần sửa

```
dk changelog add docs/ops/<file> --summary "<tóm tắt thật>" --source <postmortem, CR hoặc ->
dk render docs/ops/<file>
dk check docs/ops/<file>
```

`check` báo lỗi thì sửa rồi chạy lại.

## `html`

```
dk render --all
```

Mở `docs/html/ops/runbooks/<slug>.html` kiểm khối lệnh đọc được và liên kết
`monitoring.html`, postmortem mở được từ đĩa.
