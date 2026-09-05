---
type: release-brief
title: ""
status: draft
owner: ""
created: {{.Created}}
updated: {{.Updated}}
source: ""
created_by: dk
dk_version: {{.DKVersion}}
feature: {{.Feature}}
kind: feature
released_in: ""
---

# {{.Title}}

<!-- gợi ý: Một tính năng một file, tên theo mã Feature Spec. dk new --from chép mục 1, 2, 4, 5, 10 của spec; viết lại toàn bộ bằng ngôn ngữ người dùng: không thuật ngữ kỹ thuật (dk check no-jargon), không mã bước, không tên bảng hay hàm. kind: feature (tính năng mới) | fix (sửa lỗi). status: draft | review (người duyệt giọng văn) | ready (sẵn sàng gom vào Release notes). released_in do dk new release-notes --collect ghi, không sửa tay. -->

## 1. Giúp gì

<!-- gợi ý: một đến hai câu, người dùng làm được gì mà trước không làm được; chép từ mục 1 của spec rồi bỏ hết chữ kỹ thuật -->
{{range .Purpose}}
{{.}}{{else}}
{{end}}

## 2. Ai dùng

<!-- gợi ý: vai trò người dùng bằng tên gọi họ tự dùng, không phải tên role trong hệ thống -->
{{range .Actors}}
- {{.}}{{else}}
- {{end}}

## 3. Bước cơ bản

<!-- gợi ý: từ cột hành động của bảng hành vi; mỗi bước một câu bắt đầu bằng động từ, nói điều người dùng nhìn thấy sau mỗi bước -->
{{range $i, $a := .Actions}}
{{inc $i}}. {{$a}}{{else}}
1. {{end}}

## 4. Ảnh màn hình

<!-- gợi ý: liên kết mockup đã duyệt hoặc ảnh chụp thật do người thêm; mỗi ảnh một dòng, chú thích bằng việc người dùng đang làm, không ghi mã bước -->
{{range .Steps}}
- {{.Mockup}}{{else}}
- {{end}}

## 5. Giới hạn đã biết

<!-- gợi ý: từ mục 10 (ngoài phạm vi) và ngoại lệ chưa xử lý ở mục 6 của spec; nói thẳng điều chưa làm được, không hứa -->
{{range .Limits}}
- {{.}}{{else}}
- {{end}}
