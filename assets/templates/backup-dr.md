---
type: backup-dr
title: ""
status: draft
owner: ""
created: {{.Created}}
updated: {{.Updated}}
source: ""
created_by: dk
dk_version: {{.DKVersion}}
last_drill: ""
---

# Backup và disaster recovery: {{.Title}}

<!-- gợi ý: mục 1 là frontmatter. RPO (mất tối đa bao nhiêu dữ liệu), RTO (khôi phục trong bao lâu) do người chốt. Quy trình khôi phục phải diễn tập ít nhất một lần; ghi ngày diễn tập cuối vào last_drill (yyyy-mm-dd), dk status báo khi quá 6 tháng. -->

## 2. Mục tiêu

| Hệ thống | RPO | RTO | Người chốt |
|---|---|---|---|
| | | | |

## 3. Backup

<!-- gợi ý: cái gì được sao lưu, tần suất, giữ bao lâu, ở đâu, ai kiểm tra backup còn đọc được -->

| Dữ liệu | Tần suất | Giữ | Nơi lưu | Kiểm tra |
|---|---|---|---|---|
| | | | | |

## 4. Quy trình khôi phục

<!-- gợi ý: từng bước chạy được như runbook; mỗi bước một lệnh và kết quả mong đợi -->

1. 

```
```

## 5. Kịch bản thảm họa

<!-- gợi ý: mất một máy, mất một vùng, xóa nhầm dữ liệu, bị mã hóa tống tiền; mỗi kịch bản trỏ về bước khôi phục nào -->

| Kịch bản | Cách khôi phục | Đã diễn tập |
|---|---|---|
| | | |

## 6. Lịch sử diễn tập

<!-- gợi ý: ngày, kịch bản, thời gian thực tế so RTO, điều học được; ngày mới nhất ghi vào last_drill -->

| Ngày | Kịch bản | Thời gian | Kết quả |
|---|---|---|---|
| | | | |
