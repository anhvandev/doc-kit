---
phase: 2
title: "Skill doc-feature-list"
status: completed
priority: P2
effort: "3h"
dependencies: [1]
---

# Phase 2: Skill doc-feature-list

## Overview

Viết skill thứ 12 `doc-feature-list` (SKILL.md và `references/rules.md`) theo khuôn 11 skill hiện có; `dk skill install` cài 12 skill, test nội dung và đếm skill qua.

## Requirements

- Functional: ba phần `new`, `update`, `html`; bước 0 kiểm `dk --version`; câu dừng in đậm đầu file; `new` dừng khi brief sai `level: project`, `kind: product`, `status: approved`; mọi file tạo bằng `dk new`.
- Non-functional: SKILL.md dưới 300 dòng (nhắm 120 đến 160), rules.md dưới 300; không chứa `ak-`, `Claude Code`, `.claude/`, `.codex/`, `Edit`, `Write`, `MultiEdit` (regex `internal/skill/content_test.go:18`); `description` không trùng 3 từ liên tiếp với 11 mô tả kia (câu "Không dùng ..." cuối được bỏ khi so); tiếng Việt có dấu.

## Architecture

`internal/skill` liệt kê skill bằng cách duyệt `assets/skills/`; thêm thư mục là đủ, không sửa mã. Test `install_test.go:59` đếm 11 và lấy `res[10]`: đổi 11 → 12, `res[11]`. Bảng skill trong `assets/agent-context.md` thêm một dòng (hash khối đổi theo, `dk doctor` sẽ báo khối cũ trên dự án đã cài, đúng hành vi).

## Related Code Files

- Create: `assets/skills/doc-feature-list/SKILL.md`, `assets/skills/doc-feature-list/references/rules.md`
- Modify: `internal/skill/install_test.go` (11 → 12, `res[10]` → `res[11]`), `assets/agent-context.md` (dòng bảng sau `doc-overview`)

## Implementation Steps

1. Frontmatter SKILL.md:
   ```yaml
   name: doc-feature-list
   description: 'Feature list trong docs/overview: tách Product brief người đã duyệt thành bảng tính năng dự kiến có mã tạm, nhóm, ưu tiên MoSCoW, dẫn về mục brief; dừng chờ duyệt; điền mã spec khi tính năng có Feature Spec. Không dùng cho brief, spec, roadmap.'
   ```
   Chạy `go test ./internal/skill/` ngay để bắt trigram trùng; trùng thì đổi chữ, không đổi ý.
2. Đầu thân, sau tiêu đề: câu dừng in đậm: "**Dừng bắt buộc: chỉ soạn từ brief `level: project`, `kind: product`, `status: approved`; list chỉ đến `status: review`; không tự đặt `approved`; không tạo idea, brief hay Feature Spec từ dòng list. Mã tạm `FL-xx` không đánh số lại, dòng bỏ ghi `won't`.**" Rồi đoạn "Mọi file trong `docs/` chỉ được tạo bằng `dk new`..." chép nguyên từ `doc-overview`.
3. `## Phạm vi`: Làm: feature list từ Product brief đã duyệt; điền cột Spec; sửa list theo CR. Không làm: brief, phỏng vấn (họ Intake); Feature Spec (họ Feature Spec); roadmap (họ Plan); Product overview (họ Overview). Mỗi dòng list muốn thành spec: nói người gọi intake cấp tính năng cho dòng đó.
4. `## Bước 0. Kiểm dk` chép nguyên khuôn.
5. `## new`:
   1. Tìm và kiểm brief: `dk index intake`, đọc `docs/intake/README.md`, mở brief người chỉ (hoặc brief `approved` duy nhất `level: project`); kiểm ba trường, sai thì **dừng**, báo trường sai và ai duyệt, không tạo file. Đã có `docs/overview/feature-list.md`: không tạo lại, chuyển sang `update`.
   2. Tạo: `dk new feature-list <slug-san-pham> --from docs/intake/<yymmdd>-<slug>/brief.md --set owner="<người phụ trách>"`.
   3. Đọc nguồn: bốn mục brief và `idea.md` cùng thư mục (mục Ai gặp, Vấn đề, Kết quả, Giao diện). Bảng "Lấy từ": mục 2 Nguồn ← liên kết brief, idea; mục 3 Bảng ← brief §1 Kết quả mong muốn và §4 Tiêu chí chấp nhận (mỗi kết quả kiểm chứng được thường là một hoặc vài tính năng), idea mục Giao diện (mỗi màn hình là ứng viên); mục 4 Nhóm ← gom bảng thành 3 đến 7 nhóm theo việc người dùng làm; mục 5 Chưa rõ ← câu trong brief không tách được, tính năng suy ra mà brief không nhắc.
   4. Quy tắc điền dòng: Mô tả một câu "người dùng <làm gì> để <được gì>"; Nguồn ghi đúng mục sinh ra (`brief §1`, `idea §8`), dòng không dẫn được mục nào thì không vào bảng mà vào mục 5; Ưu tiên: `must` khi brief §4 có tiêu chí chấp nhận ứng với nó, `should` khi brief §1 nhắc nhưng §4 không đo, `could` khi chỉ idea nhắc, `won't` khi brief §3 Ngoài phạm vi nhắc (giữ dòng để không ai hỏi lại); Spec để trống. Không bịa tính năng; không quá 25 dòng, hơn thì gom.
   5. Đặt `status: review`, chạy khối "Sau mỗi lần sửa". **Dừng.** Báo người: đường dẫn, số dòng theo ưu tiên, mục 5; cách chốt (`status: approved`); nhắc mỗi dòng đi tiếp bằng intake cấp tính năng.
   Bảng đối chiếu Sai / Đúng ba dòng: "bịa tính năng đăng nhập vì sản phẩm nào cũng có" / "vào mục 5: brief không nhắc đăng nhập, cần người quyết"; "mô tả: hệ thống xử lý đơn" / "người bán xác nhận đơn để kho bắt đầu đóng gói"; "đánh số lại sau khi bỏ FL-03" / "FL-03 ghi won't, giữ số".
6. `## Sau mỗi lần sửa`:
   ```
   dk changelog add docs/overview/feature-list.md --summary "<tóm tắt thật>" --source <yymmdd>-<slug>
   dk render docs/overview/feature-list.md
   dk check docs/overview/feature-list.md
   ```
   `--source` là thư mục intake của brief khi tạo hoặc sửa theo brief, mã CR khi sửa theo CR. Tóm tắt thật: "Tách 9 tính năng, 4 nhóm, 2 chưa rõ".
7. `## update`: (a) list `draft`/`review`: sửa theo lời người, thêm dòng lấy mã kế tiếp; (b) điền cột Spec: khi người báo hoặc `dk index features` có spec mới với `title` khớp tên dòng, ghi `F-xxx` vào cột Spec, không cần CR, `--source` là mã spec; (c) list `approved` đổi nội dung khác cột Spec: cần CR `approved` hoặc `in-progress` có dòng "Feature list: Có" trong bảng tác động, thiếu thì **dừng** chỉ người sang họ CR; đổi `source` thành mã CR. Mỗi trường hợp chạy khối "Sau mỗi lần sửa".
8. `## html`: `dk render docs/overview/feature-list.md`; mở `docs/html/overview/feature-list.html` kiểm bảng 7 cột hiện đủ, liên kết brief bấm được.
9. `references/rules.md` (dưới 80 dòng): đầu file chép hai dòng nguồn như rules.md khác; mục "Nguyên tắc nền" (nguồn sự thật là brief, list là tài liệu định hướng không phải spec, người chốt); mục "Feature list" (một file một sản phẩm; 7 cột; mã tạm bất biến; ưu tiên MoSCoW và cách suy; không dòng nào thiếu Nguồn; dưới 200 dòng); mục "Quan hệ" (Product overview mục 4 dùng cùng tên nhóm; Roadmap tham chiếu `FL-xx` khi chưa có `F-xxx`; Feature catalog là bản sinh từ spec, không thay list); mục "Không làm" (không tạo intake, spec, roadmap; không xếp mốc; không ước effort).
10. `assets/agent-context.md`, sau dòng `doc-overview`: ``| `doc-feature-list` | planned feature table split from an approved Product brief: temporary FL-xx codes, groups, MoSCoW priority, link back to brief sections (`docs/overview/feature-list.md`) |``. Kiểm file vẫn dưới 120 dòng.
11. `internal/skill/install_test.go:59`: `11` → `12`, `res[10]` → `res[11]` (danh sách sắp theo tên; `doc-feature-list` đứng trước `doc-feature-spec`, nên `res[11]` vẫn là skill cuối theo bảng chữ cái; kiểm bằng cách chạy test).
12. `make lint-skills`, `go test ./...`, `go vet ./...`; `wc -l` hai file skill.

## Success Criteria

- [x] `go test ./internal/skill/` qua: name trùng thư mục, dưới 300 dòng, không từ cấm, không trigram trùng, đếm 12.
- [x] `make lint-skills` qua.
- [x] Trong thư mục tạm: `dk skill install` in 12 dòng "đã cài", có `doc-feature-list`; `dk skill uninstall` gỡ sạch.
- [x] Đọc SKILL.md theo tay với brief mẫu: mỗi bước có lệnh `dk` cụ thể hoặc quy tắc điền, không bước nào cần đoán.

## Risk Assessment

- Mô tả kích hoạt trùng trigram với `doc-overview` hoặc `doc-feature-spec` (cùng nói "Product brief", "Feature Spec"): dấu hiệu là test `content_test` fail nêu cụm trùng. Ứng phó: đổi cách nói ở mô tả, giữ ý; không nới test.
- Agent thực tế tự bịa tính năng phổ biến (đăng nhập, phân quyền): mục 5 và bảng Sai / Đúng là lớp chặn; kiểm ở phase 3 vòng thật, nếu vẫn bịa thì siết câu dừng.
