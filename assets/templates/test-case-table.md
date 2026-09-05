---
type: test-case-table
title: ""
status: draft
owner: ""
created: {{.Created}}
updated: {{.Updated}}
source: ""
created_by: dk
dk_version: {{.DKVersion}}
feature: ""
---

# Test case: {{.Title}}

<!-- gợi ý: Sinh từ mục 8 của Feature Spec {{.Feature}} bằng dk new test-case-table --from; mỗi tiêu chí chấp nhận một dòng, thêm dòng E<n> cho ngoại lệ ở mục 6 của spec. Dòng "chưa tách được" là AC lệch khung Given / When / Then trong spec. Cột Kết quả điền khi chạy: pass | fail | chưa chạy. -->

## 1. Điều kiện chung

<!-- gợi ý: từ mục 2 của spec: tác nhân, dữ liệu có sẵn, quyền -->
{{range .Background}}
- {{.}}{{end}}

## 2. Trường hợp

| Mã | Tên | Given | When | Then | Kết quả |
|---|---|---|---|---|---|
{{range .Scenarios}}{{if .Raw}}| {{.Code}} | chưa tách được: {{cell .Raw}} | | | | chưa chạy |
{{else}}| {{.Code}} | {{cell .Title}} | {{cell .Given}} | {{cell .When}} | {{cell .Then}} | chưa chạy |
{{end}}{{else}}| AC1 | | | | | chưa chạy |
{{end}}
## 3. Trường hợp biên

<!-- gợi ý: từ mục 6 của spec, mã E1, E2... cùng khung bảng trên -->

| Mã | Tên | Given | When | Then | Kết quả |
|---|---|---|---|---|---|
| | | | | | |
