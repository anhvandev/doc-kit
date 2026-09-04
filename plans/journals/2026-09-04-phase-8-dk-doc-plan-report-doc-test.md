---
title: "Phase 8 dk: doc-plan-report, doc-test"
date: 2026-09-04
summary: "11 loại tài liệu mới, kind feature cho Gherkin, sinh test case từ mục 9 spec, rule report-evidence; review 1 High 6 Medium đã sửa"
---

# Phase 8 dk: doc-plan-report, doc-test

## What happened

- Thêm 11 loại vào `assets/types.toml` (28 loại): roadmap, plan, plan-phase, report, decision-log, changelog-product, testing-strategy, test-case (`.feature`), test-case-table, ui-test-checklist, test-report. `dir` nhận tiền tố gốc `{plans_dir}` và `{in}`; `--in`, `--append`, `{hhmm}` mới trong `docs.New` và `dk new`.
- `kind = "feature"`: metadata YAML trong khối `# dk:` đầu file Gherkin (`internal/frontmatter/kinds.go`), `docs.Scan` quét thêm `.feature`.
- `internal/docs/from_spec.go`: mục 9 → Scenario, mục 3 → Background, mục 6 → checklist mã bước. AC lệch khung giữ `Raw` và in `# chưa tách được`.
- `check`: rule `report-evidence` (warning), `spec-has-test` nâng cấp đọc `.feature`; `.feature` hỏng khối `# dk:` là lỗi.
- Hai skill `doc-plan-report`, `doc-test` (9 skill). README thêm mục "Plan, report và test".

## Bugs found while running

- `report-evidence` bị chú thích gợi ý "commit abc1234" trong template đánh lừa → bỏ chú thích HTML trước khi so.
- Parser mục 9 cũng khớp `AC1, AC2` trong chú thích gợi ý → bỏ chú thích và dedupe mã AC.
- Review H1: AC dạng danh sách đánh số hoặc bảng bị bỏ im lặng → mọi dòng nhắc `ACn` chưa gặp thành Scenario Raw; CLI in số Scenario và số dòng chưa tách.
- Review M1: `--in` không ràng buộc trong `plans_dir` → siết thư mục có sẵn trong `plans_dir`.

## Decision

- `test-case` tách hai loại (Gherkin, bảng) thay vì một loại hai đuôi để giữ cơ chế `templates/<loại>.<kind>`.
- `report` luôn cần `--in` (`--in plans` cho `plans/reports/`); `id_prefix` không áp cho loại ngoài `docs/`.
- Không dùng `# language:` trong `.feature` vì dòng đầu là `# dk:`; từ khóa Gherkin tiếng Anh.
- Chấp nhận: `dk status` không quét `plans/`; `--append` không khóa file.

## Next steps

- Phase 9 (doc-release, doc-ops). Phase 10: kiểm nhánh dry-run BDD khi có công cụ, phần AI thật của các skill.
- Chưa commit (repo chưa có commit nào); `.serena/` đã thêm vào `.gitignore`.

> Historical work record — not durable authority. Prefer docs/specs/ADRs for current decisions.
