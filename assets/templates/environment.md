---
type: environment
title: ""
status: draft
owner: ""
created: {{.Created}}
updated: {{.Updated}}
source: ""
created_by: dk
dk_version: {{.DKVersion}}
---

# Environment và cấu hình: {{.Title}}

<!-- gợi ý: mục 1 là frontmatter. Chỉ tên biến, ý nghĩa, khác biệt giữa môi trường và nơi giữ giá trị thật. KHÔNG ghi giá trị secret: dòng KEY=value với value không phải placeholder <...> là lỗi dk check env-no-secret. Quy tắc này là lớp phụ, không thay công cụ quét secret của dự án. -->

## 2. Nơi giữ secret

<!-- gợi ý: vault, secret manager của CI, file .env cục bộ không commit; ai có quyền đọc -->

- 

## 3. Biến môi trường

<!-- gợi ý: một biến một dòng; cột giá trị chỉ ghi placeholder hoặc ví dụ không nhạy cảm -->

| Biến | Ý nghĩa | dev | staging | prod | Secret |
|---|---|---|---|---|---|
| | | | | | có / không |

## 4. Mẫu file cấu hình

<!-- gợi ý: giữ dạng KEY=<mô tả>, ví dụ DB_PASSWORD=<lấy từ vault> -->

```
```

## 5. Khác biệt giữa môi trường

<!-- gợi ý: điều gì cố ý khác nhau (log level, tính năng bật tắt, tài nguyên) và vì sao -->

- 
