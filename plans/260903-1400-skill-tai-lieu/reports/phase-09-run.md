# Phase 9: kết quả chạy thử doc-release, doc-ops

Ngày: 2026-09-04. Dự án mẫu `/tmp/tmp.6txSNrkHi9` (git init, `dk init`, `dk skill install`), binary `dk` build từ `go build ./cmd/dk`. Phần "đưa một người không kỹ thuật đọc Release brief" chưa làm được trong phiên tự động; gộp vào đợt chạy thử toàn vòng phase 10 cùng phần AI thật (viết lại giọng người dùng, soạn runbook từ postmortem thật).

## Chuỗi lệnh hai SKILL.md quy định, chạy tay

| Bước | Lệnh | Kết quả |
|---|---|---|
| Cài | `dk skill install` | 11 skill "đã cài"; `skill status` 11 dòng "đúng phiên bản" scope dự án (22 dòng tổng); `template list` 41 loại |
| Release brief | `dk new release-brief bo-loc --from <spec implemented>` | `docs/release/briefs/F-001.md`: `feature: F-001`, `source: F-001`, mục 2 chép câu mục đích, mục 3 hai bullet tác nhân, mục 4 đánh số `1. Mở danh sách đơn`, `2. Chọn một trạng thái` từ cột hành động, mục 5 liên kết mockup đổi gốc `../../design/mockups/F-001-B1.html`, mục 6 hai dòng từ mục 11 |
| `no-jargon` | `dk check` brief có chữ "API filter" | `warning no-jargon: từ kỹ thuật "API"` kèm dòng; bỏ chữ, cảnh báo hết. `link-broken` báo mockup chưa tạo là kịch bản, không phải lỗi rule |
| `--collect` | brief 1 `ready` (feature), brief 2 `ready` (fix); `dk new release-notes --collect v1.0.0` | `docs/release/v1.0.0.md`: "gom 2 Release brief, đã ghi released_in"; mục Mới có F-001, mục Sửa lỗi có F-002; cả hai brief `released_in: v1.0.0`; changelog 2 dòng `Phát hành trong v1.0.0 \| v1.0.0` |
| `--collect` rỗng | `dk new release-notes --collect v1.1.0` | `không có Release brief nào status ready chưa có released_in`, mã 1 |
| User guide | `dk new user-guide loc-don --from <brief> --set task="Bán hàng"`; `dk index user-guide` | `docs/release/guide/README.md` nhóm `## Bán hàng` bảng Trang, Trạng thái, Cập nhật |
| FAQ | `dk new faq --append "Quên mật khẩu? \| ..."` hai lần | lần đầu tạo và nối; lần hai "Đã nối dòng"; 2 dòng `- 2026-09-04 \| ...` |
| `env-no-secret` | `environment.md` có `DB_PASSWORD=abc123` và `DB_HOST=<host từ vault>` | `error env-no-secret: DB_PASSWORD có giá trị thật`, mã 3; đổi thành `<lấy từ vault>` thì 0 lỗi |
| Postmortem | `--set incident_at="2026-09-03 14:00"` (hôm qua) và `"2026-08-01 10:00"` | `written_within_48h: true` và `false` |
| Runbook | `dk new runbook db-mat-ket-noi` | `docs/ops/runbooks/db-mat-ket-noi.md` khung dấu hiệu, chẩn đoán, xử lý, xác nhận, leo thang; HTML render được |
| DR quá hạn | `dk new backup-dr` (`last_drill` rỗng) rồi `dk status`; đặt `last_drill: 2026-08-15` | có dòng `DR chưa diễn tập quá 6 tháng: docs/ops/backup-dr.md`; sau khi đặt ngày thì hết |
| Tầng 9 | `dk new charter charter`, `risk-register risk-register`, `meeting-notes kickoff` | `docs/governance/charter.md`, `risk-register.md`, `meetings/260904-kickoff.md`; không có skill |
| Loại còn lại | `deployment`, `monitoring` | tạo được; `dk render --all` 19 file HTML kể `release/briefs/*.html`, `ops/runbooks/*.html` |
| Skill sạch | grep `ak-`, `Claude Code` trong hai skill đã cài | 0 |

## Tiêu chí phase 9

| Tiêu chí | Kết quả |
|---|---|
| Release brief đọc được bởi người không kỹ thuật, không có từ jargon | `no-jargon` bắt đúng từ và hết khi viết lại; phần người đọc thật chưa làm, gộp phase 10 |
| `--collect` gom đúng brief, ghi `released_in`, changelog mỗi brief một dòng | đạt (test `TestCollectReleaseNotes`, `TestReleaseAndOps` và chạy tay) |
| `env-no-secret` chặn giá trị thật, cho qua placeholder | đạt |
| Runbook chạy được từng bước; Postmortem quá 48 giờ được cảnh báo | khung runbook ba khối bước có; chạy trên sự cố thật gộp phase 10. `written_within_48h` tính đúng, SKILL.md bắt skill báo người khi `false` |
| Tầng 9 tạo bằng `dk new`, không có skill | đạt |
| `skills_test.go` qua với 11 skill | đạt (`internal/skill/content_test.go`, `install_test.go` 11 skill; `go test ./...` xanh) |

Bằng chứng: `go test ./...` và `go vet ./...` qua toàn bộ sau khi thêm `internal/docs/release_test.go`, `internal/check/rules_release_test.go`, `internal/cli/phase9_test.go`.

## Sau review

`plans/reports/code-review-260904-phase-09-skill-release-ops.md`: 0 Critical, 2 High, 4 Medium; đã sửa H1, H2, M1 đến M4 và phần Low về liên kết, placeholder, ô mockup mẫu (test bổ sung `TestCollectReleaseNotes` nhánh `--force`, `TestExtractRelease` liên kết tuyệt đối và `mailto:`, `TestEnvNoSecret` chú thích và `<a>x<b>`, `TestReleaseAndOps` `last_drill` sai định dạng). `go test ./...`, `go vet ./...` xanh sau sửa. Mục chấp nhận ghi ở `plan.md` 1.15.
