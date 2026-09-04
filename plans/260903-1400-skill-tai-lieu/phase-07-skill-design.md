---
phase: 7
title: "Phase 7: Skill Design: doc-design-system, doc-design-flow"
status: done
priority: P2
effort: "5d"
dependencies: [1, 4]
---

# Phase 7: Skill Design: doc-design-system, doc-design-flow

<!-- Updated: Validation Session 1 - mockup mặc định là HTML tĩnh theo tokens; Figma qua --set external -->
<!-- Updated: 2026-09-03 22:55 - xong; quyết định triển khai ở plan.md mục 1.13; chạy thử reports/phase-07-run.md -->

## Overview

Tầng 3, chỉ khi phần mềm có giao diện. `doc-design-system` dựng Design system theo Atomic Design từ Design brief: tokens JSON, quy tắc nền, component theo lớp, patterns, templates. `doc-design-flow` làm User flow, Wireframe, Mockup HTML tĩnh, Prototype, UI spec cấp tính năng, dùng chung mã bước với Feature Spec. Bỏ toàn bộ phase nếu dự án đích không có giao diện.

## Requirements

- Functional: template `design-foundations` (quy tắc nền), `design-component` (một component, `--set layer=atom|molecule|organism|template`), `design-pattern`, `design-tokens` (sinh `tokens.json` khung W3C Design Tokens và `tokens.css` từ nó bằng `dk tokens css`), `userflow`, `wireframe`, `mockup` (file HTML tĩnh chỉ dùng CSS variables từ `tokens.css`), `ui-spec`; `dk check` quy tắc `mockup-tokens` (HTML mockup không có giá trị màu hex, px khoảng cách gõ tay ngoài `tokens.css`) và `userflow-steps` (mã bước là tập con của Feature Spec cùng `feature`).
- Non-functional: mockup mở từ đĩa không mạng; tokens.css sinh, không viết tay; lint code UI của dự án đích **không** thuộc phạm vi, chỉ lint mockup do `dk` quản.

## Architecture

- `types.toml`: `design-tokens` (dir `design/tokens`, name `tokens.json`, id `none`, không phải Markdown: `kind = "json"`), `design-foundations` (dir `design`, name `foundations.md`), `design-component` (dir `design/{layer}s`, name `{slug}.md`), `design-pattern` (dir `design/patterns`), `userflow` (dir `design/flows`, name `{feature}-flow.md`, trường `feature` bắt buộc), `wireframe` (dir `design/wireframes`, name `{feature}-{step}.md`, dạng text hoặc ASCII), `mockup` (dir `design/mockups`, name `{feature}-{step}.html`, `kind = "html"`, trường frontmatter đặt trong `<!-- dk: ... -->` đầu file), `ui-spec` (dir `design`, name `{feature}-ui-spec.md`).
- `internal/docs/new.go`: hỗ trợ `kind = json|html` với template riêng và cách ghi metadata riêng (JSON: khóa `$dk`; HTML: chú thích đầu file); `frontmatter` thêm `SplitHTMLComment`, `SplitJSONKey`.
- `dk tokens css [--in docs/design/tokens/tokens.json] [--out docs/design/tokens/tokens.css]`: đọc W3C Design Tokens (nhóm lồng, `$value`, `$type`, alias `{color.blue.500}`), sinh `:root { --color-action-primary: ... }` theo tên đường dẫn nối gạch ngang; semantic alias giải về primitive; theme sáng tối qua nhóm `$extensions.dk.theme` nếu có. Chỉ hai loại: color và dimension, đủ cho mockup; loại khác chép nguyên.
- Quy trình `doc-design-system`: bắt đầu từ Design brief `đã duyệt` (`kind: design`); thứ tự bộ tối thiểu của báo cáo: lớp 0 (thang khoảng cách, thang chữ, màu semantic bốn trạng thái, bo góc, bóng, icon, accessibility) → `tokens.json` → `dk tokens css` → `foundations.md` → 4 atoms (Text, Icon, Button, Input) → Form field → Form, Empty state, Toast → 2 templates (danh sách, form); mỗi component đủ biến thể, kích thước, trạng thái, quy tắc dùng và không dùng, ARIA; quy tắc mỗi lớp chỉ tham chiếu lớp ngay dưới; dùng thư viện UI có sẵn → chỉ viết semantic ánh xạ, quy tắc dùng, patterns, templates, giọng văn. Kiểm tra thật: sinh một mockup trang từ template chỉ bằng organism và semantic tokens, `dk check` qua.
- Quy trình `doc-design-flow`: từ mục Giao diện của brief hoặc bảng tác động CR; thứ tự User flow (Mermaid, mã bước lấy từ Feature Spec nếu có, hoặc đề xuất và sau đó Feature Spec dùng lại) → Wireframe (đen trắng, text) người duyệt → Mockup **mặc định là HTML tĩnh** mỗi trạng thái quan trọng (bình thường, rỗng, lỗi, đang tải) theo `foundations.md` và `tokens.css`, chỉ dùng organism đã có → Prototype (liên kết giữa mockup bằng thẻ a, không lưu lâu) → UI spec nếu có công cụ thiết kế sinh, không viết tay. Feature Spec liên kết mockup theo mã bước qua họ Feature Spec. Ngoại lệ có designer dùng Figma: lưu liên kết Figma và ảnh xuất tĩnh vào `design/mockups/` qua `dk new mockup --set external=<url>`; skill không sinh HTML. Design brief ghi rõ khi dự án chọn ngoại lệ này.
- `check`: `mockup-tokens` regexp `#[0-9a-fA-F]{3,8}` và `\d+px` trong `<style>` hoặc `style=""` của mockup, trừ trong `tokens.css`; `userflow-steps` so với Feature Spec `feature` cùng mã.

## Related Code Files

- Create: `assets/templates/{design-tokens.json,design-foundations,design-component,design-pattern,userflow,wireframe,mockup.html,ui-spec}.md` (đuôi theo `kind`), `assets/skills/doc-design-system/{SKILL.md,references/rules.md}`, `assets/skills/doc-design-flow/{SKILL.md,references/rules.md}`, `internal/tokens/{tokens,css}.go`, `internal/cli/tokens.go`
- Modify: `assets/types.toml`, `internal/docs/new.go`, `internal/frontmatter/frontmatter.go`, `internal/check/` (`rules_mockup.go`, `rules_userflow.go`), `internal/changelog` (đếm dòng cho `.html`, `.json` trong `docs/design/` vì đây là tài liệu, không phải bản sinh; `tokens.css` là bản sinh, bỏ qua)
- Create test: `internal/tokens/*_test.go` (alias, lồng, theme), fixture mockup có hex gõ tay → lỗi

## Implementation Steps

1. `types.toml` cho 8 loại; mở rộng `docs.New` và `frontmatter` cho `json`, `html`; test.
2. `internal/tokens`: parse W3C Design Tokens, giải alias, sinh CSS; test với file mẫu 3 cấp và 1 alias vòng → lỗi rõ.
3. Template `design-tokens.json` khung tối thiểu (space, font, color primitive và semantic, radius, shadow); `mockup.html` khung có `<link rel="stylesheet" href="../tokens/tokens.css">` (đường dẫn tương đối, mở từ đĩa được) và các vùng đặt organism; các template Markdown còn lại theo báo cáo Tầng 3.
4. `rules_mockup.go`, `rules_userflow.go`; test.
5. `rules.md` cho hai skill từ báo cáo Tầng 3, "Cấu trúc Design system", mục 6 (Design tokens, Wireframe mockup prototype), mục 7 (Design system, Wireframe Mockup, User flow), mục 10 (thư viện UI có sẵn, công cụ nội bộ, không giao diện).
6. `SKILL.md` `doc-design-system`; mô tả: "Design system cấp dự án theo Atomic Design: design tokens, quy tắc nền, atoms, molecules, organisms, templates, patterns, giọng văn giao diện; từ Design brief đã duyệt. Không làm màn hình của một tính năng."
7. `SKILL.md` `doc-design-flow`; mô tả: "Giao diện của một tính năng: user flow cùng mã bước với Feature Spec, wireframe, mockup HTML theo design system, prototype, UI spec. Không định nghĩa token hay component mới."
8. `skills_test.go`; cài; chạy thử: Design brief mẫu → tokens.json → tokens.css → foundations → 4 atoms → mockup trang danh sách; `dk check` qua; cố gõ `#333` vào mockup → lỗi; user flow có mã bước ngoài Feature Spec → lỗi.
9. Report chạy thử, chỉnh SKILL.md.

## Success Criteria

- [x] `dk tokens css` sinh CSS variables đúng từ tokens.json có alias; alias vòng báo lỗi rõ
- [x] Mockup HTML mở từ đĩa dùng đúng `tokens.css`; `dk check` bắt giá trị hex hoặc px gõ tay
- [x] User flow dùng đúng mã bước Feature Spec; lệch bị bắt
- [x] Bộ tối thiểu Design system tạo đúng thứ tự lớp 0 đến 4, mỗi component đủ 5 mục (biến thể, kích thước, trạng thái, quy tắc, ARIA) (template và thứ tự trong SKILL.md; AI điền thật gộp phase 10)
- [x] AI sinh được một mockup trang từ template chỉ bằng organism và semantic tokens, `dk check` qua (khung template qua `check`; AI điền thật gộp phase 10)
- [x] Dự án có designer Figma: `dk new mockup --set external=<url>` lưu liên kết và ảnh, không sinh HTML

## Risk Assessment

- **Phạm vi phình sang Style Dictionary**: `dk tokens css` chỉ hai loại color và dimension, một đầu ra CSS; Tailwind, theme object, Figma sync là ngoài phạm vi, ghi vào rules.md để skill không hứa.
- **Mockup do AI sinh lệch Design system** dù qua lint: lint chỉ bắt giá trị gõ tay; quy tắc "chỉ dùng organism đã có" người kiểm khi duyệt wireframe và mockup; ghi rõ mockup quan trọng cần designer.
- **`kind = html|json` làm `frontmatter` và `changelog` phức tạp**: giới hạn hai định dạng metadata (chú thích HTML đầu file, khóa `$dk` trong JSON); không thêm định dạng khác.
- **Dự án có designer Figma** lệch mặc định HTML: `external` đủ dùng; mockup Figma không qua lint `mockup-tokens`, người duyệt chịu trách nhiệm nhất quán.
