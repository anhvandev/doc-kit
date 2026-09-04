---
phase: 9
title: "Phase 9: Skill phát hành và vận hành: doc-release, doc-ops"
status: done
priority: P2
effort: "3d"
dependencies: [6]
---

# Phase 9: Skill phát hành và vận hành: doc-release, doc-ops

## Overview

Tầng 7 và 8, cộng template Tầng 9 không có skill. `doc-release` viết Release brief từ Feature Spec bằng ngôn ngữ người dùng, gộp Release notes, ghép User guide theo nhiệm vụ, FAQ. `doc-ops` viết Deployment guide, Environment không secret, Runbook một sự cố một file, Monitoring, Postmortem trong 48 giờ, Backup và DR.

## Requirements

- Functional: template `release-brief`, `release-notes`, `user-guide` (một nhiệm vụ một trang, `user-guide-index`), `faq`, `deployment`, `environment`, `runbook`, `monitoring`, `postmortem`, `backup-dr`, và Tầng 9 `charter`, `risk-register`, `meeting-notes` (chỉ `dk new`, không skill); `dk new release-brief --from <feature-spec>` chép mục 2, 4, 5, 6; `dk new release-notes --collect <version>` gom mọi release-brief `status: sẵn sàng` chưa có `released_in`; `dk check` quy tắc `no-jargon` (Release brief và User guide không chứa từ trong danh sách thuật ngữ kỹ thuật cấu hình được `dk.toml` `[release] jargon = [...]`, mặc định danh sách ngắn: API, endpoint, database, migration, backend, frontend, JSON; warning) và `env-no-secret` (Environment không có dòng `KEY=value` với value không phải placeholder `<...>`; lỗi).
- Non-functional: Release brief đọc được bởi người không kỹ thuật (kiểm bằng người trong chạy thử); Runbook chạy được từng bước, mỗi bước một lệnh hoặc một hành động quan sát được.

## Architecture

- `types.toml`: `release-brief` (dir `release/briefs`, name `{feature}.md`, `from.feature-spec`: title, purpose→"Giúp gì", actors→"Ai dùng", behavior→"Bước cơ bản" (từ bảng hành vi lấy cột hành động), ui→"Ảnh màn hình" (liên kết mockup), non_goals→"Giới hạn đã biết"; trường `released_in`), `release-notes` (dir `release`, name `{version}.md`), `user-guide` (dir `release/guide`, name `{slug}.md`, trường `task`), `user-guide-index` (generated), `faq` (dir `release`, name `faq.md`, `--append`), `deployment` (dir `ops`, name `deployment.md`), `environment` (name `environment.md`), `runbook` (dir `ops/runbooks`, name `{slug}.md`), `monitoring` (name `monitoring.md`), `postmortem` (dir `ops/postmortems`, name `{yymmdd}-{slug}.md`, trường `incident_at`, `written_within_48h` do `dk new` tính), `backup-dr` (name `backup-dr.md`, trường `last_drill`), `charter`, `risk-register`, `meeting-notes` (dir `governance`).
- `--collect`: `docs.New` mở rộng đọc mọi `release-brief` theo điều kiện, nối vào template `release-notes` theo nhóm "Mới", "Sửa lỗi" (từ `kind` trong brief), ghi `released_in` vào từng brief (đây là sửa tài liệu khác, `dk` tự ghi changelog cho từng brief với nguồn `<version>`).
- Quy trình `doc-release`: Release brief chỉ từ Feature Spec `đã phát hành` hoặc `đang triển khai` sắp xong; viết lại mục 2, 4, 5, 6 bằng ngôn ngữ người dùng, không thuật ngữ, có ảnh (liên kết mockup hoặc ảnh chụp thật do người thêm), giới hạn đã biết từ mục 11 và ngoại lệ chưa xử lý; người duyệt giọng văn. Release notes: `--collect`, người duyệt. User guide: ghép dần từ Release brief, tổ chức theo nhiệm vụ người dùng, không theo menu; `dk index user-guide`. FAQ: chỉ thêm khi câu hỏi lặp lần thứ hai, `--append`.
- Quy trình `doc-ops`: Deployment: từng bước chạy được, không giả định kiến thức ngầm, kèm rollback; đọc CI config thật để soạn. Environment: tên biến, ý nghĩa, khác biệt dev/staging/prod, **không giá trị**, `dk check` chặn. Runbook: một sự cố một file, viết sau mỗi incident, bước chẩn đoán rồi bước xử lý rồi bước xác nhận. Monitoring: metric, ngưỡng, người nhận, liên kết dashboard thật. Postmortem: không đổ lỗi cá nhân, timeline, nguyên nhân gốc, hành động khắc phục có người và hạn; skill cảnh báo khi quá 48 giờ. Backup và DR: RPO, RTO, quy trình, ngày diễn tập cuối; `dk status` liệt kê DR chưa diễn tập quá 6 tháng.
- Tầng 9: chỉ template và `dk new`; ghi trong README của skill bundle rằng không có skill vì tài liệu quản trị do người viết.

## Related Code Files

- Create: 15 template trong `assets/templates/`, `assets/skills/doc-release/{SKILL.md,references/rules.md}`, `assets/skills/doc-ops/{SKILL.md,references/rules.md}`, `internal/docs/collect.go`
- Modify: `assets/types.toml`, `internal/docs/new.go`, `internal/check/` (`rules_release.go`, `rules_env.go`), `internal/cli/index.go` (`user-guide`), `internal/cli/status.go` (DR quá hạn), `internal/config` (`[release] jargon`)
- Create test: `--collect` gom 2 brief đúng nhóm và ghi `released_in`; Environment có `DB_PASSWORD=abc` → lỗi, `DB_PASSWORD=<secret>` → qua; brief chứa "endpoint" → warning

## Implementation Steps

1. `types.toml` 15 loại; `collect.go`; test.
2. Template theo báo cáo Tầng 7, 8, 9; Runbook khung ba bước; Postmortem khung không đổ lỗi.
3. `rules_release.go`, `rules_env.go`; test.
4. `index user-guide`; `status` DR quá hạn; test.
5. `rules.md` hai skill từ báo cáo Tầng 7, 8, mục 7 (Release brief, User guide, Runbook, Postmortem), mục 10 (công cụ nội bộ không người dùng ngoài).
6. `SKILL.md` `doc-release`; mô tả: "Tài liệu cho người dùng cuối: Release brief từ Feature Spec bằng ngôn ngữ không kỹ thuật, Release notes gộp theo phiên bản, User guide theo nhiệm vụ, FAQ. Không viết tài liệu vận hành."
7. `SKILL.md` `doc-ops`; mô tả: "Tài liệu vận hành production: Deployment guide, Environment không secret, Runbook một sự cố một file, Monitoring, Postmortem trong 48 giờ, Backup và DR. Không viết cho người dùng cuối."
8. `skills_test.go`; cài; chạy thử: từ spec mẫu ra Release brief, đưa một người không kỹ thuật đọc và ghi nhận xét vào report; `--collect v1.0.0` gom 2 brief; Environment có secret bị chặn; Runbook mẫu chạy từng bước trên dự án thử.
9. Report chạy thử, chỉnh SKILL.md.

## Success Criteria

- [x] Release brief đọc được bởi người không kỹ thuật (ghi nhận trong report), không có từ trong danh sách jargon (phần người đọc thật gộp phase 10, xem `reports/phase-09-run.md`)
- [x] `--collect` gom đúng brief và ghi `released_in`, changelog có dòng cho từng brief
- [x] `env-no-secret` chặn giá trị thật, cho qua placeholder
- [x] Runbook chạy được từng bước trên dự án thử; Postmortem quá 48 giờ được skill cảnh báo (chạy trên sự cố thật gộp phase 10)
- [x] Tầng 9 tạo được bằng `dk new charter|risk-register|meeting-notes`, không có skill
- [x] `skills_test.go` qua với 11 skill

## Risk Assessment

- **Danh sách jargon gây báo giả** (từ "API" trong tên sản phẩm): danh sách cấu hình được, warning không phải lỗi, skill không tự xóa từ mà hỏi người.
- **`env-no-secret` bỏ sót** secret không theo mẫu `KEY=value`: quy tắc chỉ là lớp phụ; README ghi rõ không thay gitleaks hoặc công cụ quét secret của dự án.
- **`--collect` sửa nhiều brief một lúc** làm changelog nhiều dòng: chấp nhận, mỗi brief một dòng với nguồn là phiên bản, đúng tinh thần một file một dòng.
