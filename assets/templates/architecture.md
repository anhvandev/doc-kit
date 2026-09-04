---
type: architecture
title: ""
status: draft
owner: ""
created: {{.Created}}
updated: {{.Updated}}
source: ""
created_by: dk
dk_version: {{.DKVersion}}
---

# Architecture overview: {{.Title}}

<!-- gợi ý: mục 1 là frontmatter. Sơ đồ trước, văn xuôi sau. Soạn từ mã nguồn thật khi có (cấu trúc thư mục, file cấu hình); thành phần suy đoán đánh dấu "chưa xác nhận" ở bảng mục 3, người xác nhận từng dòng. Cập nhật khi thêm hoặc bỏ thành phần. Chi tiết từng thành phần tách sang ADR khi vượt 400 dòng. -->

## 2. Sơ đồ thành phần

<!-- gợi ý: mỗi nút một thành phần chạy độc lập (dịch vụ, cơ sở dữ liệu, hàng đợi, client); mũi tên là luồng dữ liệu chính. Không vẽ chi tiết bên trong một thành phần. -->

```mermaid
flowchart LR
    client[Client] --> api[API]
    api --> db[(Cơ sở dữ liệu)]
```

## 3. Thành phần

<!-- gợi ý: mỗi nút trong sơ đồ một dòng. Cột "Xác nhận": đã xác nhận | chưa xác nhận (suy đoán từ mã, chờ người xác nhận). -->

| Thành phần | Vai trò | Nơi trong mã | Xác nhận |
|---|---|---|---|
| | | | chưa xác nhận |

## 4. Luồng dữ liệu

<!-- gợi ý: hai đến ba luồng chính, mỗi luồng một đoạn ngắn theo thứ tự thành phần đi qua -->

## 5. Ranh giới hệ thống

<!-- gợi ý: hệ thống ngoài (dịch vụ thứ ba, hệ thống cũ) và giao tiếp qua gì; thứ nằm ngoài trách nhiệm của dự án -->

- 

## 6. Tech stack

| Lớp | Công nghệ | Phiên bản | Lý do hoặc ADR |
|---|---|---|---|
| | | | |

## 7. Quyết định liên quan

<!-- gợi ý: liên kết đến từng ADR còn hiệu lực trong docs/adr/; không chép nội dung, không liên kết file chưa có -->

- 
