---
title: "Phase 6 dk: doc-feature-spec"
date: 2026-09-03
summary: "Template feature-spec 5 biến thể format và has_ui, mã bước hậu tố B2a, rule spec-section-order và cr-approval-order, skill doc-feature-spec; review 2 High cùng gốc step-codes, sửa 7 finding"
---

# Phase 6 dk: doc-feature-spec

## Đã làm
- `assets/templates/feature-spec.md` (160 dòng, một file): biến thể `format` `spec | use-case | story | crud | state` qua `{{if}}` ở mức mục; `has_ui: false` bỏ mục 6; frontmatter thêm `has_ui`, `format`.
- `tmpl.Data` thêm `Format`, `HasUI`; `docs.New` đọc `--set format`, `--set has_ui` trước khi render, validate cả hai; `--set k=true|false` ghi YAML bool (`frontmatter.SetBool`).
- `render.StepCodes` nhận mã hậu tố `B2a` (chèn bước giữa, không đánh số lại); neo HTML `#step-B3a`.
- `internal/check`: `spec-section-order` (tiêu đề `## 2.` đến `## 11.` đúng thứ tự, bỏ 6 khi `has_ui: false`, bỏ 4 khi `crud`), `cr-approval-order` (spec `source` CR chưa chốt mà `updated` mới hơn CR là lỗi). Fixture F-005, F-006, F-007, CR-260901-cho.
- Skill `doc-feature-spec` (SKILL.md 156 dòng, rules.md 132 dòng): `new` từ brief `approved` cấp tính năng, `update` chỉ theo CR `approved`/`in-progress` có dòng Feature Spec `Có`, `html`.
- Chạy thử tay ở `/tmp`: đủ 6 tiêu chí; CR chèn B3a chỉ đổi 4 dòng nội dung + `source`; changelog `+10 −4` khớp numstat. Report `plans/260903-1400-skill-tai-lieu/reports/phase-06-run.md`.

## Lỗi gặp
- Regex mã bước `^B\d+$` không nhận `B2a` dù plan yêu cầu hậu tố; nới ở một chỗ dùng chung cho render và check.
- Kịch bản thử `cr-approval-order` lần đầu không bắt vì `dk changelog add` cho CR chạy sau khi sửa spec, bump `updated` CR lên cùng phút. Giới hạn thật: `updated` chỉ đến phút.
- Reviewer đầu treo 36 phút không ghi report (kẹt lúc dựng dự án mẫu); dừng, chạy lại với phạm vi chỉ đọc mã, xong trong 5 phút.

## Review
2 High cùng gốc: `step-codes` báo giả khi spec cố ý không có sơ đồ (`crud` mặc định vì mục 6 có `B1`; luồng tuyến tính xóa mermaid theo đúng SKILL.md). Sửa: chỉ so khi có sơ đồ; mục 6 của `crud` dùng cột "Màn hình". 3 Medium sửa hết (validate `has_ui`, tiêu đề "use case" là bảng hành vi, fence ``` và ~~~ theo độ dài). 4 Low chấp nhận, ghi plan 1.12. Định dạng `updated` giữ nguyên, không thêm offset như phase text dự kiến.

## Tiếp theo
Phase 7 `doc-design-system`, `doc-design-flow` (song song được với 8, 9 vì 6 đã xong). Phần AI thật của phase 3 đến 6 gộp vào phase 10.

> Historical work record — not durable authority. Prefer docs/specs/ADRs for current decisions.
