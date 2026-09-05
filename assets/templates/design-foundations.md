---
type: design-foundations
title: ""
status: draft
owner: ""
created: {{.Created}}
updated: {{.Updated}}
source: ""
created_by: dk
dk_version: {{.DKVersion}}
ui_library: ""
---

# Quy tắc nền: {{.Title}}

<!-- gợi ý: Lớp 0 của Design system, áp cho mọi lớp trên. Soạn sau tokens.json và trước atoms. source: Design brief đã duyệt. ui_library: tên thư viện UI dùng sẵn (MUI, Ant, shadcn) nếu có; khi đó mục 1 đến 4 chỉ ghi ngoại lệ so với thư viện. Mọi giá trị hiển thị tham chiếu token ở tokens/tokens.css, không gõ tay. -->

## 1. Thang và semantic tokens

<!-- gợi ý: liệt kê tên token theo vai trò (không chép giá trị): màu nền, chữ, hành động, viền, bốn trạng thái; khoảng cách inset và stack; chữ body và heading; bo góc; bóng. Component chỉ được tham chiếu các tên ở đây. -->

| Vai trò | Token | Dùng khi |
|---|---|---|
| Nền trang | `--color-bg-page` | |
| Hành động chính | `--color-action-primary` | |

## 2. Lưới và điểm gãy

<!-- gợi ý: số cột, gutter theo space, điểm gãy (sm, md, lg) và bố cục đổi thế nào ở mỗi điểm -->

| Điểm gãy | Rộng | Bố cục |
|---|---|---|
| | | |

## 3. Accessibility

<!-- gợi ý: tương phản tối thiểu 4.5:1, vùng chạm --size-touch, focus ring bắt buộc (--color-border-focus), điều hướng bàn phím, thứ tự tab, nhãn ARIA cho icon-only -->

- 

## 4. Icon

<!-- gợi ý: một bộ icon duy nhất, kích thước theo thang, cách đặt tên, icon-only phải có nhãn -->

## 5. Giọng văn giao diện

<!-- gợi ý: cách gọi người dùng, thì, độ dài; mẫu câu cho thành công, lỗi, rỗng, xác nhận xóa; thống nhất với Glossary mục 2 -->

| Tình huống | Mẫu câu |
|---|---|
| Lỗi hệ thống | |
| Trạng thái rỗng | |

## 6. Quy tắc lớp

<!-- gợi ý: mỗi lớp chỉ tham chiếu lớp ngay dưới: organism không gọi thẳng token, page không định nghĩa style, component không tham chiếu primitive. Thư viện UI có sẵn: lớp 0 primitive và lớp 1 đến 3 lấy từ thư viện, chỉ viết semantic ánh xạ -->

- 
