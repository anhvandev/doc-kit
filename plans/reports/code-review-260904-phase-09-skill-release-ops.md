---
title: "Code review phase 9: skill phát hành và vận hành (doc-release, doc-ops)"
date: 2026-09-04
reviewer: code-reviewer
scope: "phase 9 chưa commit (repo chưa có commit nào, mọi file untracked)"
---

# Code review phase 9: doc-release, doc-ops

## 1. Phạm vi và kiểm chứng

- Đọc: `assets/types.toml` (13 loại mới), 13 template, `assets/skills/doc-release/`, `assets/skills/doc-ops/`, `internal/docs/{collect.go,new.go,from_spec.go}`, `internal/check/{check.go,rules_release.go,rules_env.go}`, `internal/cli/{new.go,index.go,status.go}`, `internal/config/config.go`, `internal/tmpl/tmpl.go`, test mới, `README.md`.
- Chạy: `go vet ./...` sạch, `gofmt -l .` không ra file, `go test ./...` qua toàn bộ 15 package.
- Chạy thử tay trong dự án tạm `/tmp/p9` (git thật): dựng Feature Spec, `release-brief --from`, `--collect`, `check`, `index user-guide`, `status`, `faq --append`, postmortem, 3 loại Tầng 9.

## 2. Nghiệm thu

| Tiêu chí | Kết quả |
|---|---|
| (a) release-brief --from spec: mục 2, 3, cột hành động mục 5, cột mockup mục 6 (đổi gốc), mục 11 | Đạt. Liên kết `../design/...` thành `../../design/...`; liên kết tuyệt đối `/abs/x.html` giữ nguyên |
| (b) --collect gom brief ready chưa released_in, nhóm feature/fix, ghi released_in + 1 dòng changelog/brief, rỗng là mã 1 | Đạt. Changelog `+1 −1`, nguồn đúng phiên bản |
| (c) env-no-secret | Đạt. `DB_PASSWORD=abc123`, `export TOKEN=real` báo lỗi mã 3; `<lấy từ vault>` qua; số dòng đúng tuyệt đối |
| (d) no-jargon | Đạt. Mỗi từ một lần một file, bỏ chú thích HTML và khối mã, gộp hoa thường, cấu hình được |
| (e) postmortem written_within_48h | Đạt (`true`/`false` theo `incident_at`) |
| (f) status DR quá hạn | Đạt (rỗng, quá 6 tháng đều liệt kê; trong hạn thì không) |
| (g) index user-guide theo task | Đạt, nhóm chưa có task xếp cuối |
| (h) faq --append | Đạt, lần đầu tạo từ template rồi nối |
| (i) ràng buộc nội dung skill | Đạt, 11 skill, mọi file dưới 300 dòng, test nội dung qua |

Không có phát hiện mức Critical.

## 3. High

**H1. `--collect` không nguyên tử, lỗi giữa chừng để lại trạng thái nửa vời**
`internal/docs/new.go:288-299`, `internal/docs/collect.go:47-62`, `internal/cli/new.go:106-112`.
Thứ tự là: ghi file Release notes, rồi lặp ghi `released_in` từng brief, rồi lặp ghi changelog từng brief. Lỗi ở brief thứ hai (đĩa đầy, quyền, git hỏng) để lại notes đã ghi, brief 1 đã đánh dấu, brief 2 chưa, changelog lệch; không có hoàn tác và thông báo lỗi không nói file nào đã sửa.
Sửa: đánh dấu toàn bộ brief vào bộ nhớ và ghi sau cùng, hoặc khi lỗi thì in danh sách file đã sửa để người sửa tay.

**H2. `--collect --force` âm thầm loại brief đã phát hành khỏi Release notes cũ**
`internal/docs/new.go:173-180` và `201-206`, `internal/docs/collect.go:28`.
`--force` bỏ qua kiểm tra file đã tồn tại nhưng `collectBriefs` vẫn lọc `released_in` rỗng, nên chạy lại `--collect v1.0.0 --force` sinh lại `v1.0.0.md` chỉ với brief mới ready. Đã kiểm: notes ban đầu có F-001 và F-002, sau `--force` chỉ còn F-003, trong khi F-001 và F-002 vẫn ghi `released_in: v1.0.0`. Thông báo "gom 1 Release brief" không cảnh báo mất 2 dòng.
Sửa: khi `--collect` kèm `--force`, gom thêm brief có `released_in` bằng đúng phiên bản đang tạo (hoặc từ chối `--force` với `--collect`).

## 4. Medium

**M1. README lệch: `dk index` đã sinh 5 thư mục**
`README.md:99` vẫn ghi `docs/features`, `docs/adr`, `docs/cr`, `docs/intake`; `internal/cli/index.go:18-19` đã thêm `release/guide`.
Sửa: bổ sung `docs/release/guide` và một câu về cách nhóm theo `task` vào câu đó.

**M2. `env-no-secret` không bỏ chú thích HTML và khối mã, khác `no-jargon`**
`internal/check/rules_env.go:22-30` quét thẳng `m.Body`, trong khi `internal/check/rules_release.go:58-88` có `proseLines`. Ví dụ minh họa đặt trong `<!-- ... -->` sẽ bị báo lỗi cứng (mã 3), người viết không có cách nào ghi ví dụ giá trị giả không dạng `<...>`.
Sửa: quyết định rõ và ghi vào comment; nếu giữ nguyên thì nói trong template environment rằng ví dụ phải dùng placeholder kể cả trong chú thích.

**M3. Comment sai hợp đồng của trường `Released` trong JSON**
`internal/docs/new.go:60-61` ghi "(tuyệt đối)", nhưng `internal/cli/new.go:111` đổi sang đường dẫn tương đối gốc dự án trước khi in JSON.
Sửa: đổi comment thành "tuyệt đối trong `docs.New`; `dk new --json` in tương đối gốc dự án", giống `Path`.

**M4. `last_drill` sai định dạng bị gộp chung với quá hạn**
`internal/cli/status.go:51`: `time.Parse` lỗi cũng thành `dr_overdue`, người đọc không biết là chưa diễn tập hay gõ sai ngày.
Sửa: tách thông điệp cho trường hợp `err != nil` và `last_drill != ""`.

## 5. Low

- `internal/cli/status.go:33`: comment nói `drDrillMax`, hằng tên `drDrillMonths`. Sửa comment.
- `internal/check/rules_release.go:17-19`: nhánh `len(c.Jargon) == 0` không bao giờ chạy vì `internal/config/config.go:86-87` luôn điền `DefaultJargon`, tức không tắt được rule bằng `dk.toml`. Ghi rõ trong README hoặc bỏ nhánh chết.
- `internal/check/rules_release.go:54` trùng logic `internal/check/rules_glossary.go:84`. Cho glossary dùng `countFMLines`.
- `internal/docs/from_spec.go:260-277`: `mdLinkRe` không khớp liên kết có tiêu đề `](x "t")` (không đổi gốc) và đổi gốc nhầm `](mailto:...)` vì chỉ loại trừ `://`. Thêm `mailto:`/`tel:` vào điều kiện bỏ qua.
- `internal/docs/from_spec.go:199-217`: `uiSteps` không khử trùng mã bước, và giá trị mẫu "chưa có, xem họ Design" của template Feature Spec được chép thẳng vào mục 5 Release brief. Lọc ô mockup không phải liên kết.
- `internal/check/rules_env.go:28`: `KEY=<a>giá-trị-thật<b>` lọt vì chỉ so tiền tố `<` và hậu tố `>`. Đổi sang khớp trọn `^<[^<>]*>$`.
- `internal/cli/new.go:40-42,80-82`: `--collect "..."` cho slug rỗng, lỗi in ra là `slug "" không hợp lệ`, không nhắc đến `--collect`. Thêm thông điệp riêng.
- `internal/docs/new.go:104-105,120-122`: dùng sai loại cho `--collect`/`--append` trả mã 1 chứ không phải 2. Nhất quán với quyết định phase 8, chỉ ghi nhận.
- `internal/cli/index.go:36-46`: `dk index all` tạo `docs/release/guide/README.md` cả khi dự án chưa có trang guide nào. Chấp nhận được, nhất quán với 4 thư mục cũ, chỉ ghi nhận.
- `internal/cli/new.go:107`: `filepath.Rel` nuốt lỗi (nit cũ từ phase 1, nay thêm một chỗ).
- Thiếu test: `rebaseLinks` với liên kết tuyệt đối, neo `](#x)` và có tiêu đề; `--collect` khi file notes đã tồn tại (cả có và không `--force`, tức H2). `internal/docs/release_test.go:48-60` mới phủ nhánh tương đối.

## 6. Hợp đồng công khai và tương thích

- `Result` thêm `released` và `statusReport` thêm `dr_overdue`: chỉ thêm trường, `omitempty` cho `released`; `dr_overdue` luôn có mảng (khởi tạo rỗng ở `internal/cli/status.go:44`) nên không ra `null`. Không phá tương thích.
- `internal/docs/new.go:115-119` mở rộng validate placeholder sang `version`: chỉ chạm loại dùng `{version}` và `--set version=`, không loại nào trước phase 9 dùng khóa này. Không phải regression.
- 13 loại mới đều `id = "none"` và không loại cũ nào đổi `required`/`statuses`. Tài liệu đã tạo không bị `dk check` báo thêm.
- Mã thoát 1/2/3, thông điệp tiếng Việt, `gofmt` theo đúng nếp cũ.

## 7. Câu hỏi chưa giải quyết

1. H2: `--collect --force` nên gom lại brief đã có `released_in` bằng phiên bản đang tạo, hay nên cấm hẳn `--force` khi có `--collect`? Việc này đổi hành vi cờ nên cần chủ dự án chốt.
2. M2: `env-no-secret` có cố ý quét cả chú thích HTML và khối mã không? Nếu có thì template environment nên nói rõ; nếu không thì dùng `proseLines`.
3. Rule `no-jargon` hiện không tắt được bằng `dk.toml` (danh sách rỗng quay về mặc định). Đây là quyết định giống `[check] warn_lines`, xác nhận giữ nguyên?
4. Tiêu chí "Release brief đọc được bởi người không kỹ thuật" và "Runbook chạy được từng bước trên dự án thử" cần người thật; chưa có report chạy thử nào cho phase 9 trong `plans/260903-1400-skill-tai-lieu/reports/`. Gộp vào phase 10 như các phase trước?
