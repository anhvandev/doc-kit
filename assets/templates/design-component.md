---
type: design-component
title: ""
status: draft
owner: ""
created: {{.Created}}
updated: {{.Updated}}
source: ""
created_by: dk
dk_version: {{.DKVersion}}
layer: {{.Layer}}
uses: []
---

# {{.Title}} ({{.Layer}})

<!-- gợi ý: mục 1 là frontmatter. layer: atom (không chứa atom khác) | molecule (chỉ chứa atom) | organism (chứa molecule và atom) | template (chứa organism, định nghĩa vùng đặt). uses: tên component lớp ngay dưới được dùng; không được vượt lớp. Thư viện UI có sẵn: ghi component thư viện tương ứng ở mục 2 và chỉ ghi ngoại lệ. -->

## 2. Mục đích và khi nào dùng

<!-- gợi ý: hai đến ba câu; dùng khi, không dùng khi (chỉ sang component khác) -->

## 3. Biến thể

<!-- gợi ý: primary, secondary, ghost, destructive...; mỗi biến thể ghi token semantic dùng cho nền, chữ, viền -->

| Biến thể | Nền | Chữ | Viền | Dùng khi |
|---|---|---|---|---|
| | | | | |

## 4. Kích thước

<!-- gợi ý: sm, md, lg theo thang chung: --size-control-*, --space-inset-*, --font-size-* -->

| Kích thước | Cao | Đệm | Chữ |
|---|---|---|---|
| md | `--size-control-md` | `--space-inset-md` | `--font-size-md` |

## 5. Trạng thái

<!-- gợi ý: mặc định, hover, active, focus, disabled, loading, lỗi; mỗi trạng thái khác gì và dùng token nào -->

| Trạng thái | Khác biệt nhìn thấy | Token |
|---|---|---|
| focus | focus ring | `--color-border-focus` |

## 6. Quy tắc dùng và không dùng

- Dùng: 
- Không dùng: 

## 7. Accessibility

<!-- gợi ý: vai trò ARIA, thuộc tính bắt buộc, điều hướng bàn phím, nhãn khi icon-only -->

| Mục | Quy định |
|---|---|
| Vai trò | |
| Bàn phím | |

## 8. Cấu tạo

<!-- gợi ý: molecule trở lên: danh sách thành phần lớp dưới và cách sắp; template: vùng đặt (slot) và organism cho phép -->
