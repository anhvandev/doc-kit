# Code review phase 4: skill tiếp nhận doc-intake, doc-cr

- Ngày: 2026-09-03
- Phạm vi: `assets/skills/doc-intake/`, `assets/skills/doc-cr/`, `assets/types.toml`,
  `assets/templates/{idea,interview,brief,cr}.md`, `internal/doctype/doctype.go`,
  `internal/docs/new.go`, `internal/cli/index.go`, test kèm theo, README.
- Repo chưa có commit nên review theo nội dung file, không dùng `git diff`.
- Đã chạy: `go vet ./...` (sạch), `go test ./...` (mọi package `ok`), và build binary
  chạy thử trong `/tmp/dkchk` để kiểm hành vi `new`, `check`, `index`, `refs`.

Không có finding mức High.

## Medium

### M1. `doc-cr` bỏ sót changelog tại điểm dừng bắt buộc

`assets/skills/doc-cr/SKILL.md:130-137`. Khối `dk changelog add`, `dk render`,
`dk index cr`, `dk check` nằm trong mục "### 5. Sau khi người duyệt". Quy trình dừng
bắt buộc ở cuối mục 3 với `status: review`, nên CR đang chờ duyệt không có dòng
changelog nào. Vi phạm hợp đồng "mọi thay đổi tài liệu đều ghi changelog" (plan mục
Hợp đồng) và Success Criteria 6 của phase; hệ quả thực tế là pre-commit chặn commit
của người dùng. `doc-intake` không mắc lỗi này vì bước 5 đứng riêng.

Đề xuất: tách khối bốn lệnh thành mục riêng "Sau mỗi lần sửa CR", đặt trước mục 5,
và dẫn chiếu lại từ mục 5.

### M2. Trường `required` mới phá tương thích với tài liệu phiên bản trước

`assets/types.toml:12` (idea `level`), `:38` (brief `level`, `kind`), `:50`
(cr `requester`). Kiểm chứng thật trên file thiếu trường:

```
docs/intake/260903-cu/brief.md: error frontmatter-required: thiếu trường level, kind
docs/intake/260903-y-tuong/idea.md: error frontmatter-required: thiếu trường level
```

Mức `error`, không phải cảnh báo, nên `dk check` trả mã 3. Repo chưa phát hành nên
tác động thực tế bằng không, nhưng đây là thay đổi hợp đồng công khai và chưa được
ghi ở phần quyết định của `plan.md` hay README.

Đề xuất: thêm một dòng quyết định phase 4 vào `plan.md`, ghi rõ tài liệu tạo trước
phase 4 phải bổ sung `level`, `kind`, `requester` bằng tay.

## Low

### L1. `dk new cr` mới tạo không qua `dk check`

`assets/templates/cr.md:11-12` để `requester: ""` và `owner: ""`, cả hai nằm trong
`required`. Lệnh mẫu `assets/skills/doc-cr/SKILL.md:36` chỉ set `priority` và
`requester`. Kết quả chạy thật:

```
docs/cr/CR-260903-thu-nghiem.md: error frontmatter-required: thiếu trường owner, requester
```

Success Criteria 4 ("`dk check` qua trên CR mẫu") chỉ đúng sau khi người điền `owner`.

Đề xuất: thêm `--set owner="<người phụ trách>"` vào lệnh mẫu trong SKILL.md.

### L2. Interview sinh từ CR mang trường `level` vô nghĩa

`assets/templates/interview.md:11` mặc định `level: feature`; `from.cr` chỉ chép
`title`, nên file trong `docs/cr/<CR>/` giữ `level: feature` không có ý nghĩa. Không
gây lỗi vì `level` không thuộc `required` của interview.

Đề xuất: doc-cr xóa dòng `level` khi tạo interview cho CR, hoặc bỏ mặc định trong
template và để `from.idea` chép vào.

### L3. `beside_source` rộng hơn phạm vi đã chốt

`internal/docs/new.go:77` chỉ yêu cầu nguồn nằm trong `docs/`, nên
`dk new interview x --from docs/features/F-001-y.md` tạo
`docs/features/F-001-y/interview.md`. Quyết định 1.7 và phase chỉ nói tới nguồn CR.
Chỉ mục vẫn lọc đúng theo `dir` nên không hồi quy, nhưng hành vi chưa được ghi.

### L4. Vị trí test lệch phase spec

Phase yêu cầu `assets/skills/skills_test.go`; thực tế là
`internal/skill/content_test.go`. Vị trí mới hợp lý hơn vì `assets` chỉ có
`embed.go`. Nên sửa dòng "Create test" trong phase file cho khớp.

## Nit

- `internal/skill/content_test.go:44` đếm ký tự xuống dòng, nên file đúng 300 dòng
  không có newline cuối vẫn qua ngưỡng.
- `assets/skills/doc-intake/SKILL.md:120` dùng `--source <yymmdd>-<slug>` cho
  changelog, trong khi `source:` frontmatter là `<thư mục>/<file>`. Lệch dạng, không
  gây lỗi.
- README không nhắc `level`, `kind`, `requester` là trường bắt buộc mới.

## Đối chiếu theo yêu cầu review

| Mục | Kết quả |
|---|---|
| (a) Success Criteria | 4 tiêu chí kiểm được không cần Claude Code đều đạt, tiêu chí 4 phụ thuộc L1; hai tiêu chí cần người đối thoại được ghi trung thực là chưa chạy trong `reports/phase-04-run.md` |
| (b) SKILL.md | Khớp mục Architecture trừ M1. Mọi lệnh `dk` có thật, đúng cú pháp: `--version`, `new`, `changelog add --summary --source`, `render [--all --index]`, `index features\|cr\|intake`, `refs`, `status`, `check`. Câu dừng in đậm ở đầu cả hai file. Không có `Claude`, `Codex`, `target`, `.claude/`, `ak-`, tên tool. Dài 144 và 148 dòng |
| (c) rules.md | Không có quy tắc bịa hay chép sai so với báo cáo mục 1, 1a, 1b, 4, 7, 10. Bảng "quy tắc → mục báo cáo" có ở cả hai file. Một chỗ diễn giải lại: intake rules ghi "Intake brief là đầu vào (trước Feature Spec)", bỏ cụm "tổng kết của một Change Request"; nên giữ nguyên vì cụm đó tự mâu thuẫn trong báo cáo |
| (d) Hồi quy | `go test ./...` và `go vet ./...` qua. Chỉ mục `features` và `adr` không đổi; loại `adr` chưa có trong `types.toml` nên bộ lọc `dir` không chạm. Quyết định 1.7 giữ nguyên, có test `internal/docs/new_test.go:141-144` |
| (e) Public contract | Có phá tương thích, xem M2 |
| (f) Mã Go | Theo pattern hiện có: comment và lỗi tiếng Việt, file snake_case, cobra, không thêm thư viện |

## Câu hỏi chưa giải quyết

- M2 xử lý theo hướng nào: ghi nhận là thay đổi có chủ đích, hay hạ mức trường mới
  xuống cảnh báo cho file `created_by` khác phiên bản hiện tại?
- L3 có nên siết `beside_source` chỉ cho nguồn loại `cr` không?

Status: DONE_WITH_CONCERNS
Summary: Phase 4 đạt yêu cầu về nội dung skill, template, chỉ mục và test; `go test`
và `go vet` qua. Còn hai vấn đề cần xử lý trước khi sang phase 5.
Concerns/Blockers: M1 làm CR chờ duyệt không có dòng changelog nên pre-commit sẽ
chặn commit; M2 là thay đổi hợp đồng `required` chưa ghi vào quyết định của plan.
