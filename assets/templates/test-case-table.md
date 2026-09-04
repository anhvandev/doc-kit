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

<!-- gợi ý: mục 1 là frontmatter. Sinh từ mục 9 của Feature Spec {{.Feature}} bằng dk new test-case-table --from; mỗi tiêu chí chấp nhận một dòng, thêm dòng E<n> cho ngoại lệ ở mục 7 của spec. Dòng "chưa tách được" là AC lệch khung Given / When / Then trong spec. Cột Kết quả điền khi chạy: pass | fail | chưa chạy. -->

## 2. Điều kiện chung

<!-- gợi ý: từ mục 3 của spec: tác nhân, dữ liệu có sẵn, quyền -->
{{range .Background}}
- {{.}}{{end}}

## 3. Trường hợp

| Mã | Tên | Given | When | Then | Kết quả |
|---|---|---|---|---|---|
{{range .Scenarios}}{{if .Raw}}| {{.Code}} | chưa tách được: {{cell .Raw}} | | | | chưa chạy |
{{else}}| {{.Code}} | {{cell .Title}} | {{cell .Given}} | {{cell .When}} | {{cell .Then}} | chưa chạy |
{{end}}{{else}}| AC1 | | | | | chưa chạy |
{{end}}
## 4. Trường hợp biên

<!-- gợi ý: từ mục 7 của spec, mã E1, E2... cùng khung bảng trên -->

| Mã | Tên | Given | When | Then | Kết quả |
|---|---|---|---|---|---|
| | | | | | |
