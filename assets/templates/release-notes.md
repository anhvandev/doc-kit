---
type: release-notes
title: ""
status: draft
owner: ""
created: {{.Created}}
updated: {{.Updated}}
source: ""
created_by: dk
dk_version: {{.DKVersion}}
version: "{{.Version}}"
---

# Release notes {{.Version}}

<!-- gợi ý: Sinh bằng dk new release-notes --collect <phiên bản>: gom mọi Release brief status ready chưa có released_in, nhóm theo kind. Người duyệt rồi đặt published; sau đó không sửa, phiên bản sau viết file mới. Mỗi dòng một câu cho người dùng, liên kết về Release brief để đọc thêm. -->

## 1. Mới
{{range .Collected}}{{if eq .Kind "feature"}}
- [{{.Title}}]({{.Link}})
{{- end}}{{end}}

## 2. Sửa lỗi
{{range .Collected}}{{if eq .Kind "fix"}}
- [{{.Title}}]({{.Link}})
{{- end}}{{end}}

## 3. Cần lưu ý

<!-- gợi ý: thay đổi làm người dùng phải làm khác đi (đổi menu, cần đăng nhập lại, dữ liệu cũ hiển thị khác); không có thì ghi "không" -->

- 
