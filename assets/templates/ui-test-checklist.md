---
type: ui-test-checklist
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

# Test giao diện: {{.Title}}

<!-- gợi ý: mục 1 là frontmatter. Sinh từ mục 6 của Feature Spec {{.Feature}} bằng dk new ui-test-checklist --from: mỗi mã bước một dòng so màn hình thật với mockup đã duyệt. Chuẩn so sánh là mockup và Design system, không phải cảm tính. Người duyệt sai lệch; AI chạy và báo sai lệch. -->

## 2. Checklist theo mã bước

<!-- gợi ý: đánh dấu khi màn hình thật khớp mockup; lệch thì ghi ở mục 3 -->
{{range .Steps}}
- [ ] {{.Code}}: khớp mockup {{.Mockup}}{{else}}
- [ ] B1: khớp mockup (spec chưa có mục 6 hoặc chưa có mockup){{end}}

## 3. Sai lệch

<!-- gợi ý: mã bước, điểm lệch (khoảng cách, màu, chữ, trạng thái), ảnh chụp màn hình, quyết định của người duyệt -->

| Mã bước | Sai lệch | Ảnh | Quyết định |
|---|---|---|---|
| | | | |
