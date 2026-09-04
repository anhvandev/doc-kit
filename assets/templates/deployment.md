---
type: deployment
title: ""
status: draft
owner: ""
created: {{.Created}}
updated: {{.Updated}}
source: ""
created_by: dk
dk_version: {{.DKVersion}}
---

# Deployment guide: {{.Title}}

<!-- gợi ý: mục 1 là frontmatter. Soạn từ cấu hình CI thật (.github/workflows, .gitlab-ci.yml, Dockerfile, script deploy), không từ trí nhớ. Từng bước chạy được: mỗi bước một lệnh trong khối ``` hoặc một hành động quan sát được, ghi kết quả mong đợi. Không giả định kiến thức ngầm; người mới vào dự án phải làm theo được. -->

## 2. Tổng quan

<!-- gợi ý: môi trường (dev, staging, prod), cách mã lên từng môi trường (CI tự động, lệnh tay), ai được phép -->

| Môi trường | Cách triển khai | Ai được phép | Thời gian ước tính |
|---|---|---|---|
| | | | |

## 3. Điều kiện trước khi triển khai

<!-- gợi ý: quyền, công cụ cài trên máy, biến môi trường (xem environment.md), test phải qua -->

- 

## 4. Các bước triển khai

<!-- gợi ý: mỗi bước một lệnh và kết quả mong đợi -->

1. 

```
```

## 5. Xác nhận sau triển khai

<!-- gợi ý: URL kiểm tra, dashboard, log cần nhìn, trong bao lâu -->

- 

## 6. Rollback

<!-- gợi ý: từng bước quay về phiên bản trước; điều kiện quyết định rollback; ai quyết -->

1. 

## 7. Liên kết

- Environment: [environment.md](environment.md)
- Runbook: 
