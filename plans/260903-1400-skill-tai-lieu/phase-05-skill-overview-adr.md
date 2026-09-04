---
phase: 5
title: "Phase 5: Skill nền sản phẩm: doc-overview, doc-adr"
status: done
priority: P1
effort: "3d"
dependencies: [4]
---

# Phase 5: Skill nền sản phẩm: doc-overview, doc-adr

## Overview

Tầng 2: `doc-overview` soạn Product overview từ Product brief, Architecture overview, Glossary và duy trì Feature catalog qua `dk index features`; `doc-adr` ghi ADR đánh số, bất biến, có chỉ mục.

## Requirements

- Functional: template `product-overview`, `architecture`, `glossary`, `adr`; `dk new overview --from <brief>` (loại `product-overview`), `dk new adr <slug>` ra `ADR-0001-slug.md`; `dk index features|adr`; `dk check` quy tắc `adr-immutable` (ADR `status: đã chốt` đổi thân là lỗi, dựa `git diff` với HEAD) và `glossary-term` (thuật ngữ in đậm lần đầu trong Feature Spec phải có trong Glossary, warning).
- Non-functional: như phase 4; Product overview một trang, cảnh báo khi vượt 150 dòng qua `types.toml` `warn_lines=150`.

## Architecture

- `types.toml`: `product-overview` (dir `overview`, name `product-overview.md`, id `none`, `from.brief`: title, outcome, constraints, non_goals, acceptance), `architecture` (name `architecture.md`, id `none`), `glossary` (name `glossary.md`, id `none`, `check.warn_lines=300` theo số thuật ngữ thì để `check` đếm dòng bảng), `adr` (dir `adr`, id `seq:ADR-{n:04}`, name `{id}-{slug}.md`, statuses `đề xuất | đã chốt | thay thế | bãi bỏ`, trường `supersedes`, `superseded_by`).
- Template `adr`: Bối cảnh, Lựa chọn đã xét (bảng), Quyết định, Hệ quả, Liên kết (brief hoặc CR nguồn). Template `architecture`: sơ đồ Mermaid thành phần trước, văn xuôi sau, mục tech stack dạng bảng, ranh giới hệ thống. Template `glossary`: bảng `Thuật ngữ | Định nghĩa | Không nhầm với | Nguồn`.
- Quy trình `doc-overview`: `new` Product overview chỉ từ Product brief đã duyệt (`level: project`, `kind: product`, `status: đã duyệt`), khác thì dừng; soạn nháp, người chốt. Architecture: soạn từ mã nguồn thật khi có (đọc cấu trúc thư mục, file cấu hình), sơ đồ thành phần, đánh dấu "chưa xác nhận" cho thành phần suy đoán. Glossary: quét Feature Spec và brief lấy thuật ngữ in đậm lần đầu, đề xuất định nghĩa, người chốt. Feature catalog: chỉ chạy `dk index features`, không viết tay. `update`: chỉ khi có CR hoặc brief mới đã duyệt trỏ đến; ghi `source`.
- Quy trình `doc-adr`: `new` khi người hoặc CR nêu quyết định kỹ thuật quan trọng; soạn từ thảo luận trong CR hoặc chat; `status: đề xuất`; người chốt đổi trạng thái. Quyết định mới thay quyết định cũ: `dk new adr` mới với `supersedes: ADR-000X`, skill sửa ADR cũ **chỉ** frontmatter `status: thay thế`, `superseded_by`, không sửa thân; `dk index adr`.
- `check` mới: `adr-immutable` so thân file ADR đã chốt với HEAD qua `gitx`; `glossary-term` regexp `\*\*[^*]+\*\*` lần xuất hiện đầu trong mỗi Feature Spec, tra bảng Glossary.
- `index features`: cột `mã | tên | trạng thái | chủ sở hữu | brief hoặc CR nguồn | cập nhật`; `index adr`: cột `mã | tiêu đề | trạng thái | thay thế | ngày`.

## Related Code Files

- Create: `assets/templates/{product-overview,architecture,glossary,adr}.md`, `assets/skills/doc-overview/{SKILL.md,references/rules.md}`, `assets/skills/doc-adr/{SKILL.md,references/rules.md}`
- Modify: `assets/types.toml`, `internal/check/` (thêm `rules_adr.go`, `rules_glossary.go` và test), `internal/cli/index.go`
- Create test: fixture ADR đã chốt bị sửa thân → lỗi; Feature Spec có thuật ngữ in đậm không trong Glossary → warning

## Implementation Steps

1. `types.toml` và 4 template; test render.
2. `rules_adr.go`, `rules_glossary.go` và test.
3. `index features`, `index adr`; test số dòng và thứ tự (ADR theo số, features theo mã).
4. `rules.md` cho hai skill từ báo cáo Tầng 2, mục 7 (Product overview, Architecture, Glossary, ADR), mục 10.
5. `SKILL.md` `doc-overview`; mô tả: "Product overview, Architecture overview, Glossary và Feature catalog cấp dự án: soạn từ Product brief đã duyệt hoặc từ mã nguồn, sinh catalog bằng chỉ mục. Không viết Feature Spec, không ghi ADR."
6. `SKILL.md` `doc-adr`; mô tả: "Architecture Decision Record: ghi bối cảnh, lựa chọn, quyết định, hệ quả; đánh số, không sửa sau khi chốt, ADR mới thay ADR cũ; chỉ mục ADR. Không dùng cho quyết định nhỏ (Decision log)."
7. `skills_test.go` qua; cài, chạy thử: từ Product brief mẫu ra Product overview dưới 150 dòng; 2 Feature Spec mẫu → catalog 2 dòng; ADR-0001 chốt rồi tạo ADR-0002 thay thế → ADR-0001 chỉ đổi frontmatter, `dk check` qua; cố sửa thân ADR-0001 → `dk check` lỗi.
8. Ghi report chạy thử, chỉnh SKILL.md.

## Success Criteria

- [x] Feature catalog khớp số file trong `docs/features/`, sinh hoàn toàn bằng `dk index`
- [x] ADR mới thay ADR cũ không sửa thân ADR cũ; `dk check` bắt thân ADR đã chốt bị đổi
- [x] Product overview tạo từ Product brief đã duyệt, chép đúng trường, dưới 150 dòng; brief chưa duyệt → skill dừng (phần "skill dừng" là quy tắc trong SKILL.md, chưa kiểm thật trong Claude Code; xem `reports/phase-05-run.md`)
- [x] Glossary có thuật ngữ do skill quét từ Feature Spec (quy tắc trong SKILL.md, chưa kiểm thật); `dk check` cảnh báo thuật ngữ in đậm chưa định nghĩa
- [x] `skills_test.go` qua với 4 skill (tên thật: `internal/skill/content_test.go`)

## Risk Assessment

- **`adr-immutable` báo giả khi đổi định dạng (trailing space, CRLF)**: so sau khi chuẩn hóa xuống dòng và cắt khoảng trắng cuối dòng.
- **Architecture overview suy đoán sai từ mã**: skill bắt buộc đánh dấu "chưa xác nhận" và người xác nhận từng thành phần; không tự chốt.
- **Glossary phình vì quét mọi chữ in đậm**: chỉ lấy lần in đậm đầu trong mục 2, 5, 8 của Feature Spec; skill đề xuất, người chọn giữ.
