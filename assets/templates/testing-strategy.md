---
type: testing-strategy
title: ""
status: draft
owner: ""
created: {{.Created}}
updated: {{.Updated}}
source: ""
created_by: dk
dk_version: {{.DKVersion}}
format: gherkin
---

# Testing strategy: {{.Title}}

<!-- gợi ý: mục 1 là frontmatter. Chính sách, không phải danh sách test; một file, cập nhật khi đổi công cụ hoặc ngưỡng. AI đề xuất từ stack thật (package.json, go.mod, cấu hình CI); người chốt công cụ. format: gherkin (test case là file .feature chạy được) | table (bảng Markdown). Công cụ BDD đã chốt ghi vào dk.toml [test] bdd_cmd để skill chạy dry-run. status: draft | review | approved. -->

## 2. Loại test và phạm vi

<!-- gợi ý: unit, integration, e2e, giao diện, hiệu năng; mỗi loại: mục đích, phạm vi bao phủ mong đợi, ai viết -->

| Loại | Mục đích | Phạm vi mong đợi | Ai viết |
|---|---|---|---|
| | | | |

## 3. Công cụ

<!-- gợi ý: lấy từ stack thật; ghi phiên bản; công cụ BDD nếu format gherkin -->

| Loại | Công cụ | Lệnh chạy |
|---|---|---|
| | | |

## 4. Định dạng test case

<!-- gợi ý: gherkin hay table, lý do; quy ước tag (@F-xxx @ACn), nơi đặt file (docs/test/) -->

- Định dạng: gherkin
- Tag truy vết: `@<mã Feature Spec> @AC<n>`

## 5. Ngưỡng và cửa kiểm

<!-- gợi ý: ngưỡng bao phủ, test nào chặn merge, test nào chạy trước phát hành -->

- 

## 6. Test giao diện

<!-- gợi ý: checklist theo mã bước so với mockup đã duyệt hoặc visual regression; chuẩn so sánh là mockup. Bỏ mục này khi không có giao diện -->

- 
