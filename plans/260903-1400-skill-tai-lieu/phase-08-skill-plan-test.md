---
phase: 8
title: "Phase 8: Skill tiến độ và test: doc-plan-report, doc-test"
status: done
priority: P2
effort: "4d"
dependencies: [6]
---

# Phase 8: Skill tiến độ và test: doc-plan-report, doc-test

## Overview

Tầng 5 và 6. `doc-plan-report` viết Roadmap, Plan (một thư mục một đợt việc), Report kèm bằng chứng, Decision log, CHANGELOG sản phẩm. `doc-test` viết Testing strategy, sinh Test case từ mục 9 Feature Spec (Gherkin hoặc bảng), checklist test giao diện theo mã bước, Test report.

## Requirements

- Functional: template `roadmap`, `plan`, `plan-phase`, `report`, `decision-log`, `changelog-product`, `testing-strategy`, `test-case`, `ui-test-checklist`, `test-report`; `dk new plan <slug>` tạo thư mục `plans/{yymmdd-hhmm}-{slug}/plan.md`, `dk new plan-phase <slug> --in <plan-dir>` tạo `phase-NN-slug.md`; `dk new test-case --from <feature-spec>` chép mục 9 (AC1..) thành khung Gherkin `Scenario` từng AC; `dk new report` vào `plans/<dir>/reports/` hoặc `plans/reports/`; `dk check` quy tắc `report-evidence` (Report phải có ít nhất một liên kết commit, file output hoặc khối kết quả test; warning) và `spec-has-test` từ phase 2 nâng cấp đọc `docs/test/`.
- Non-functional: `plan` và `report` nằm ở `plans/` ngoài `docs/`, **không** vào `CHANGELOG-DOCS.md` (đã quy định), nhưng `dk check` vẫn kiểm liên kết; nếu repo đích có công cụ plan riêng thì `dk.toml` `plans_dir` trỏ đúng thư mục và skill tôn trọng quy ước tên có sẵn.

## Architecture

- `types.toml`: `roadmap` (dir `plan`, name `roadmap.md`), `plan` (dir `{plans_dir}`, subdir `{yymmdd-hhmm}-{slug}/`, name `plan.md`), `plan-phase` (dir theo `--in`, name `phase-{n:02}-{slug}.md`, id `seq` đếm trong thư mục), `report` (dir `{plans_dir}/{active}/reports` hoặc `{plans_dir}/reports`, name `{type}-{yymmdd-hhmm}-{slug}.md`), `decision-log` (dir `plan`, name `decision-log.md`, nối thêm dòng qua `dk new decision-log --append "..."`), `changelog-product` (dir `docs`, name `CHANGELOG.md`, Keep a Changelog), `testing-strategy` (dir `test`, name `strategy.md`), `test-case` (dir `test`, name `{feature}-cases.feature` khi `format: gherkin`, `{feature}-cases.md` khi `format: table`; `kind` theo đuôi), `ui-test-checklist` (dir `test`, name `{feature}-ui.md`), `test-report` (dir `test/reports`, name `{yymmdd}-{slug}.md`).
- `docs.New` mở rộng: `--in <dir>` cho `plan-phase`; `--append` cho `decision-log` (một dòng `- {ngày} | {quyết định} | {người} | {liên kết}`); `from.feature-spec` cho `test-case`: parse mục 9 lấy từng `AC{n}` với Given/When/Then, sinh `Scenario: AC{n} ...` với `@F-012 @AC1` tag để truy vết; `from.feature-spec` cho `ui-test-checklist`: mỗi mã bước một dòng `- [ ] B1: khớp mockup <liên kết mục 6>`.
- Quy trình `doc-plan-report`: Plan: mục tiêu, phạm vi, non-goals, giai đoạn, tiêu chí hoàn thành; người duyệt phạm vi và non-goals trước khi AI làm việc theo plan; mỗi giai đoạn một file phase. Report: viết ngay khi kết thúc giai đoạn, kèm bằng chứng (đường dẫn output test, mã commit), không sửa sau khi nộp; câu hỏi mở cuối. Decision log: một dòng một quyết định, ghi khi được ủy quyền, gồm quyết định giao diện. CHANGELOG sản phẩm: tổng hợp từ Report và commit, ngôn ngữ người dùng, người duyệt. Roadmap: chỉ người chốt; skill chỉ đề xuất sắp xếp từ Feature catalog.
- Quy trình `doc-test`: Testing strategy một file, đề xuất từ stack thật (đọc `package.json`, `go.mod`, cấu hình CI), người chốt công cụ. Test case: từ Feature Spec `đã duyệt`, mỗi AC một scenario, thêm trường hợp biên từ mục 7 ngoại lệ; định dạng theo Testing strategy (`format` gherkin hay table); nếu Gherkin thì file `.feature` chính là tài liệu và phải chạy được bởi công cụ BDD đã chốt (skill chạy `dry-run` của công cụ nếu có, ví dụ `cucumber --dry-run`; không có thì báo chưa kiểm). Test giao diện: checklist theo mã bước so với mockup đã duyệt. Test report: theo đợt phát hành, đính output thật, người quyết định xử lý fail.
- `check`: `report-evidence`; `spec-has-test` đọc `source` trong `.feature` qua tag `@F-012` hoặc frontmatter `.md`.

## Related Code Files

- Create: `assets/templates/{roadmap,plan,plan-phase,report,decision-log,changelog-product,testing-strategy,test-case.feature,test-case-table,ui-test-checklist,test-report}.md`, `assets/skills/doc-plan-report/{SKILL.md,references/rules.md}`, `assets/skills/doc-test/{SKILL.md,references/rules.md}`, `internal/docs/from_spec.go` (parse mục 9 và mã bước)
- Modify: `assets/types.toml`, `internal/docs/new.go` (`--in`, `--append`, `kind=feature`), `internal/check/` (`rules_report.go`, sửa `spec-has-test`), `internal/config` (`plans_dir`)
- Create test: parse mục 9 với 3 AC → 3 Scenario; `--append` decision-log không phá dòng cũ; Report thiếu bằng chứng → warning

## Implementation Steps

1. `types.toml` 11 loại; `docs.New` với `--in`, `--append`, `kind=feature`; test.
2. `from_spec.go`: parse mục 9 và bảng hành vi; test với spec mẫu phase 6.
3. Template; `test-case.feature` khung `Feature:` từ mục 2, `Background:` từ mục 3, `Scenario:` từng AC; `ui-test-checklist` từ mã bước.
4. `rules_report.go`, sửa `spec-has-test`; test.
5. `rules.md` từ báo cáo Tầng 5, 6, mục 7 (Plan, Report, Test case, CHANGELOG, Decision log, Roadmap), mục 10.
6. `SKILL.md` `doc-plan-report`; mô tả: "Roadmap, Plan theo đợt việc và file phase, Report kèm bằng chứng, Decision log, CHANGELOG sản phẩm theo Keep a Changelog. Không viết Feature Spec, không viết test."
7. `SKILL.md` `doc-test`; mô tả: "Testing strategy, Test case sinh từ tiêu chí chấp nhận của Feature Spec (Gherkin hoặc bảng), checklist test giao diện theo mã bước so với mockup, Test report có output thật. Không sửa Feature Spec."
8. `skills_test.go`; cài; chạy thử: plan 2 phase → report phase 1 có liên kết commit → `dk check` qua; report không bằng chứng → warning; từ spec mẫu 4 AC ra `.feature` 4 scenario, `cucumber --dry-run` hoặc công cụ đã chốt qua nếu có trên máy; checklist UI đủ mã bước.
9. Report chạy thử, chỉnh SKILL.md.

## Success Criteria

- [x] `dk new plan` và `plan-phase` tạo đúng cấu trúc thư mục và số thứ tự
- [x] `dk new test-case --from <spec>` ra một Scenario cho mỗi AC, tag truy vết đúng mã tính năng
- [x] Gherkin sinh ra qua `dry-run` của công cụ BDD đã chốt (khi có); không có công cụ thì skill báo rõ "chưa kiểm chạy được" (nhánh "có công cụ" chưa kiểm vì máy không có; xem `reports/phase-08-run.md`)
- [x] Report có liên kết commit hoặc output test; thiếu thì `dk check` cảnh báo
- [x] Decision log nối dòng không phá nội dung cũ
- [x] `spec-has-test` không còn cảnh báo cho spec đã có `.feature` tương ứng

Kết quả chạy thử: `reports/phase-08-run.md`. Quyết định triển khai: `plan.md` mục 1.14.

## Risk Assessment

- **Công cụ BDD chưa chốt** (câu hỏi mở): template Gherkin chuẩn không phụ thuộc công cụ; bước dry-run là tùy chọn theo `dk.toml` `[test] bdd_cmd = "..."`. Tín hiệu vỡ: `.feature` không parse được bởi công cụ chọn sau; phản ứng: sửa template, không sửa CLI.
- **Xung đột với công cụ plan có sẵn trong repo đích** (thư mục `plans/` đã có quy ước khác): `dk.toml` `plans_dir` và `plan_name_pattern` cho phép khớp quy ước; skill đọc file plan có sẵn để bắt chước tên. Nếu repo dùng công cụ plan riêng, skill chỉ viết Report và Decision log, không tạo plan.
- **Parse mục 9 vỡ khi người viết AC lệch khung**: `from_spec` báo rõ AC nào không parse được và để trống Scenario đó thay vì bỏ qua im lặng.
