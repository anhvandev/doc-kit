# Code review — phase 3 "Cài đặt cho Claude"

Ngày: 2026-09-03. Phạm vi: mọi file phase 3 (chưa commit), đối chiếu
`plans/260903-1400-skill-tai-lieu/phase-03-install-claude.md` và `plan.md` §1.6.
Kiểm chứng bằng đọc mã + chạy `go vet ./...`, `go test ./...` và thử binary
`go build -o /tmp/dk ./cmd/dk` trong repo tạm.

## Kết quả cổng chất lượng

- `go vet ./...`: sạch.
- `go test ./...`: mọi package qua (`internal/cli` 0.28s), không test nào cần Claude Code.
- Success Criteria tự động: đều đạt (init không đụng `.claude/`; metadata `dk_*`;
  cài lại "không đổi"; sửa tay bị từ chối, `--force` ghi đè; uninstall bỏ qua thư mục
  lạ; settings golden về đúng bản gốc; pre-commit fail-open khi rút `dk` khỏi PATH).
  Bước 9 (kiểm thật trong Claude Code) chưa làm — xem H2, M5, M6 trước khi làm.
- Hợp đồng công khai: mã thoát 0/1/2/3 giữ nguyên; `Target` khớp `plan.md` §1.6
  (Claude có thêm `SettingsPath` ngoài interface, chấp nhận được); `init --json`
  giữ đủ khóa cũ (`config`, `docs_dir`, `changelog`, `dirs`) và thêm `pre_commit`,
  `next`; `changelog add` không đổi hành vi sau khi tách sang `changelog.Record`
  (test phase 1–2 vẫn qua).
- `settings.json`: giữ khóa lạ, thứ tự khóa cấp một, thụt 2 khoảng, không escape
  HTML (`echo 'a & b' > /dev/null` trong golden giữ nguyên). Xóa file chỉ khi object rỗng.
- `HooksDir` dùng `git rev-parse --git-path hooks`, đúng cả với worktree, thư mục con
  và `core.hooksPath` (đã thử: husky → trả `.husky`).

## Lỗi thật

### H1 — `dk hook uninstall` xóa luôn hook của người dùng nằm chung một mục

`internal/target/claude.go:322-327`: vòng lọc bỏ **cả nhóm** (`matcher` + mọi
`hooks[]`) khi nhóm chứa một lệnh `dk hook run`, thay vì chỉ bỏ lệnh của dk.

Tái hiện:

```json
{"hooks":{"PreToolUse":[{"matcher":"Write","hooks":[
  {"type":"command","command":"dk hook run pre-write"},
  {"type":"command","command":"my-own-guard.sh"}]}]},"model":"opus"}
```

`dk hook uninstall` → `{"model":"opus"}`; `my-own-guard.sh` biến mất không cảnh báo.
Kịch bản đời thật: người dùng thêm lệnh của mình vào đúng mục `matcher: "Write"` do
dk tạo (cách tự nhiên nhất để hai hook cùng chạy trên Write), rồi gỡ dk → mất hook.
Yêu cầu phase ghi rõ "xóa **đúng các mục có `command` bắt đầu bằng `dk hook run`**".

Sửa: lọc ở mức phần tử `hooks[]`, chỉ bỏ nhóm khi mảng `hooks` rỗng sau khi lọc;
giữ nguyên các khóa khác của nhóm (raw) khi ghi lại.

### H2 — Hook post-edit làm cổng `changelog pending` (và pre-commit) luôn xanh

`internal/hook/run.go:78` ghi mục "chưa tóm tắt / trực-tiếp" cho mọi lần Edit|Write
trong `docs/`. `internal/changelog/changelog.go:150-161` (`Since`) coi file là "đã ghi"
khi có **bất kỳ** mục nào từ phút HEAD trở đi, không xét tóm tắt.

Hệ quả: sau khi agent sửa tài liệu, `dk changelog pending` luôn rỗng → pre-commit
không bao giờ chặn nữa, dù không ai viết tóm tắt thật. Cổng chất lượng do phase 1–2
dựng trở thành trang trí ngay khi phase 3 được cài.

Đây là hệ quả của thiết kế đã chốt trong phase file (hook ghi "chưa tóm tắt"), nên
**không tự sửa**. Hai lựa chọn để chủ dự án quyết:

1. `Since` bỏ qua mục có `Summary == NoSummary` (và/hoặc `Source == "trực-tiếp"`)
   → pre-commit vẫn ép người/skill viết tóm tắt thật.
2. Giữ nguyên và ghi rõ trong README rằng pre-commit chỉ bảo đảm "có dòng", không
   bảo đảm "có tóm tắt".

### M3 — Dòng của hook không bao giờ được thay bằng tóm tắt thật (changelog nhân đôi)

`internal/changelog/merge.go:29` chỉ gộp khi `last.Source == e.Source`. Hook dùng
source `trực-tiếp`, còn `dk changelog add --summary ...` thường có source rỗng hoặc
`brief-x` → hai dòng cho cùng một lần sửa.

Tái hiện (cùng phút, cùng file):

```
- 16:40 | features/a.md | không git, 5 dòng | thêm mục tiêu | -
- 16:40 | features/a.md | không git, 5 dòng | chưa tóm tắt | trực-tiếp
```

Trái với ý ghi trong README ("để skill hoặc người thay bằng tóm tắt thật sau").
Sửa gợi ý: khi mục mới trùng file trong cửa sổ gộp và mục cũ là `trực-tiếp` +
`chưa tóm tắt`, cho phép thay bất kể `Source`.

### M4 — `dk skill uninstall` nhận tên chưa kiểm, `RemoveAll` ra ngoài thư mục skill

`internal/skill/install.go:113` `dest := filepath.Join(dir, name)` với `name` lấy
thẳng từ tham số dòng lệnh, không chuẩn hóa; `:124` gọi `os.RemoveAll(dest)`.

Tái hiện:

```
$ dk skill uninstall ../../other/doc-smoke
../../other/doc-smoke	đã gỡ	/tmp/t1/other/doc-smoke
```

Thư mục ngoài `.claude/skills` bị xóa. Giới hạn duy nhất là đích phải có
`metadata.dk_installed_by: dk`, nên ví dụ `../../../.claude/skills/<tên>` của dự án
khác hoặc `~/.claude/skills/<tên>` vẫn xóa được. Lệnh này hay do agent sinh tham số,
nên tên không tin cậy là biên thật.

Sửa: từ chối tên chứa `/`, `\`, `..` hoặc đường dẫn tuyệt đối (áp cho cả
`Install`, `Uninstall`, `Status`); `Install` hiện được che nhờ `fs.Stat` trong
`Files()` chứ không nhờ kiểm tra chủ ý.

### M5 — Cập nhật skill không nguyên tử; đứt giữa chừng thì `--force` cũng không cứu được

`internal/skill/install.go:68-76`: `os.RemoveAll(dest)` rồi `write()` từng file theo
thứ tự map ngẫu nhiên. Nếu đứt (Ctrl-C, hết đĩa, lỗi quyền) sau khi ghi
`references/rules.md` mà chưa ghi `SKILL.md`, thư mục đích mất dấu vết →
`install.go:58-59` báo "skill không do dk cài" và `--force` **không** vượt qua được
(nhánh `tr.By != installedBy` đứng trước nhánh force).

Tái hiện tương đương với thư mục rỗng:

```
$ mkdir -p .claude/skills/doc-smoke && dk skill install --force
/tmp/t1/.claude/skills/doc-smoke: skill không do dk cài; dùng tên khác hoặc xóa tay  (exit 1)
```

Người dùng phải `rm -rf` bằng tay. Sửa: ghi vào thư mục tạm cạnh đích rồi
`os.Rename`, hoặc cho `--force` vượt khi thư mục đích rỗng.

### M6 — Không có khóa cho ghi `CHANGELOG-DOCS.md` và `updated:` khi hook chạy song song

`internal/changelog/record.go:48-79` là read-modify-write toàn file, không khóa.
Claude Code chạy nhiều tool call song song trong một lượt; hai PostToolUse cho hai
file `docs/` khác nhau sẽ cùng đọc `CHANGELOG-DOCS.md`, cùng ghi đè → mất một dòng.
Trước phase 3 chỉ có người gõ lệnh nên thực tế không đụng; hook làm rủi ro thành thật.

Sửa: khóa file (`flock`/lock file trong `docs/`) quanh đoạn load → format → write,
hoặc ghi qua temp + rename với retry.

### M7 — Hook post-edit sửa lại chính file Claude vừa ghi (rủi ro "file modified since read")

`internal/changelog/record.go:54-61` ghi `updated:` vào tài liệu **sau** khi
Write/Edit hoàn tất. Claude Code so nội dung/mtime với bản đã đọc và từ chối Edit
tiếp theo trên file bị đổi ngoài phiên. Chưa kiểm chứng được ngoài Claude Code —
đưa vào bước 9 làm mục kiểm bắt buộc: sửa cùng một file `docs/` hai lần liên tiếp
bằng Edit. Nếu vỡ, cân nhắc bump `updated:` chỉ khi nguồn không phải `trực-tiếp`.

## Nit (không chặn)

- `internal/target/claude.go:224` `"hooks": null` trong settings có sẵn → thoát 1
  ("khóa hooks: không phải JSON object"), không có đường tự phục hồi. Coi `null` như rỗng.
- `internal/target/claude.go:218` `os.WriteFile` thẳng lên `settings.json`: đứt điện
  giữa chừng thì cụt file. Nên temp + rename.
- `internal/cli/hook.go:33,49` in "đã cài/gỡ 2 hook" là hằng số, không phản ánh số
  mục thật sự thay đổi (cài lần hai idempotent vẫn nói "đã cài 2 hook").
- Sau `skill uninstall` còn thư mục rỗng `.claude/skills/`; tiêu chí bước 8 nói
  `.claude/` chỉ còn `settings.json` hoặc bị xóa. Cân nhắc `os.Remove` thư mục cha khi rỗng.
- README mục "Gỡ" ghi cứng `rm .git/hooks/pre-commit`, sai với repo đặt
  `core.hooksPath` (mã thì xử lý đúng). Nên nói "xóa file mà `git rev-parse --git-path hooks` trỏ tới".
- `internal/cli/hook.go:13` `var stdin io.Reader` mức package, test gán trực tiếp →
  trạng thái toàn cục; truyền qua `app` sẽ sạch hơn.
- `internal/cli/skill.go:130` `skillStatus` gán `a.root` mà không nạp `cfg`; hiện vô hại.
- `internal/cli/skill.go:108` `--json` in `null` khi `res` rỗng (lỗi trước khi cài
  skill nào); `[]` sẽ dễ tiêu thụ hơn.
- `internal/target/claude.go:198-205` file `settings.json` do người dùng tạo sẵn với
  nội dung `{}` cũng bị xóa khi gỡ hook (không phân biệt "do dk tạo").

## Câu hỏi còn treo

1. H2: chọn phương án 1 (siết `pending`) hay 2 (ghi rõ trong README)?
2. M3: có chấp nhận gộp bỏ qua `Source` khi mục cũ là `trực-tiếp/chưa tóm tắt` không?
3. M7 chỉ kiểm được ở bước 9 — ai chạy và khi nào?
