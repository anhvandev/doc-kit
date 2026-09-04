# Code review — Phase 11: target Codex + Windows

Ngày: 2026-09-04 · Reviewer: code-reviewer · Phạm vi: `internal/target`, `internal/hook`,
`internal/cli`, `internal/gitx`, Makefile/.goreleaser/CI, README + docs.
Không sửa file nào; chỉ báo cáo.

## Kết quả kiểm chạy được

| Lệnh | Kết quả |
|---|---|
| `go test ./...` | PASS (mọi package) |
| `go vet ./...` | PASS |
| `GOOS=windows go vet ./...` | PASS |
| `make lint-skills` | sạch |
| `make build-all` | 6 nền tảng build được (linux/darwin/windows × amd64/arm64) |

Acceptance (a):

- `skill install --target codex` → `uninstall` để lại `.codex/` rỗng (thư mục vẫn còn, xem H-2).
- `status --target claude,codex` ra cả hai target (`skill.Row.Target`, `status.go:48`).
- Golden `hooks.json` giữ nguyên mục người dùng, gỡ về byte-for-byte
  (`internal/target/codex_test.go:36`), idempotent 2 vòng.
- `doctor` có dòng skill/hook cho codex, tự thêm codex khi có `.codex/` (xem H-2).
- Ma trận build 6 nền tảng: Makefile `PLATFORMS`, `.goreleaser.yaml` goos×goarch, CI
  `windows-latest` trong ma trận test.

Không có regression đường Claude: hành vi `--target` mặc định `claude` giữ nguyên;
`TestInstallEndToEnd` và `TestDoctor` cover cả đường cũ.

---

## Findings (xếp theo mức độ)

### C-1 (Critical–correctness) `*** Move to:` không được đọc → thủng cả hai hook

`internal/hook/run.go:52` (`patchFiles`) chỉ nhận `*** Add File: ` và `*** Update File: `.
Định dạng `apply_patch` của Codex đổi tên/di chuyển file bằng
`*** Update File: <cũ>` + `*** Move to: <mới>`. Hậu quả, đã kiểm bằng binary thật:

```
*** Update File: notes.md
*** Move to: docs/features/moved.md
```

- `pre-write`: không deny → agent tạo tài liệu mới trong `docs/` mà không qua `dk new`
  chỉ bằng cách tạo file ngoài `docs/` rồi move vào. Đây đúng là hành vi mà hook sinh ra
  để chặn.
- `post-edit`: changelog không ghi dòng nào cho `features/moved.md` (kiểm: `grep -c
  moved docs/CHANGELOG-DOCS.md` = 0) → pre-commit sẽ báo thiếu changelog sau đó, hoặc
  file lọt qua nếu `Tracks` không thấy.

Sửa: đọc thêm `*** Move to: ` và trả cả đường dẫn đích. Cùng chỗ nên cân nhắc
`*** Delete File:` cho `post-edit` (hiện bỏ qua có chủ đích — hợp lý, nhưng nên nói rõ
trong comment vì hiện comment chỉ giải thích "không tạo tài liệu").

### H-1 (High) Hook Codex không bắt được ghi file qua `shell`/`exec_command`

`internal/target/codex.go:19` cố định matcher `apply_patch`. Codex vẫn ghi file được
bằng tool shell (`bash -c 'cat > docs/x.md'`). Plan bước 6 yêu cầu ghi rõ giới hạn này
trong README nếu hook không bắt hết; README (dòng 44-49) và `docs/kien-truc.md:87` chỉ mô
tả matcher `apply_patch`, không nói "hook không phải rào chắn đủ, pre-commit mới là bảo
đảm". Đề nghị thêm một câu giới hạn, đây là kỳ vọng an toàn của người dùng.

### H-2 (High) `.codex/` rỗng sót lại làm `dk doctor` fail vĩnh viễn

Kịch bản kiểm thật:

```
dk skill install --target codex && dk hook install --target codex
dk skill uninstall --target codex && dk hook uninstall --target codex
dk doctor            # không có --target
```

`uninstall` xóa hết nội dung nhưng để lại thư mục `.codex/` rỗng. `internal/cli/doctor.go:92`
tự thêm codex khi `os.Stat(<root>/.codex)` thành công → doctor in 11 dòng
`skill … (codex) chưa cài` + `hook (codex, dự án) có 0/2` và thoát 3, dù người dùng vừa
chủ động gỡ codex. Người dùng không có cách nào tắt ngoài `--target claude` mỗi lần.

Sửa (một trong hai, hoặc cả hai):
- điều kiện tự thêm nên là `.codex/skills` hoặc `.codex/hooks.json` tồn tại, không phải
  thư mục `.codex/`;
- `uninstall` xóa thư mục cha khi rỗng (nhất quán với `.claude/`, hiện cũng sót).

### H-3 (High) Mã thoát không nhất quán cho `--target` sai/rỗng

Hợp đồng nêu trong task: target lạ = 2, `--target` rỗng = 2. Thực tế đo:

| Lệnh | Mã |
|---|---|
| `dk skill install --target x` | 2 ✔ |
| `dk skill install --target ""` | 2 ✔ |
| `dk doctor --target x` | 3 ✘ |
| `dk doctor --target ""` | 3 ✘ |

`internal/cli/doctor.go:98-101` nuốt `*ExitError` mã 2 của `targetsOf` thành một dòng
bảng rồi trả `codeCheck`. Sai cờ (usage) không nên biến thành "kiểm tra không qua".
`internal/cli/phase10_test.go:161` đang đóng đinh hành vi này, nên nếu giữ thì phải là
quyết định có chủ đích và ghi vào `docs/lenh.md` — hiện docs không nói.

### M-1 (Medium) Plan bước 7 "gitx (CRLF)" chưa làm

`internal/gitx/gitx.go:31` `strings.TrimRight(out, "\n\x00")` không cắt `\r`. Chỉ
`filepath.ToSlash` cho pathspec được thêm. Nếu bất kỳ đầu ra git nào trên Windows về
CRLF thì `IsRepo` so `out == "true"` sẽ false và `Root()`/`HooksDir()` nhận đường dẫn
dính `\r`. Rủi ro thấp (git plumbing thường ra LF) nhưng plan liệt kê mục này; hoặc làm
(`TrimRight(out, "\r\n\x00")`), hoặc ghi rõ trong report là đã đánh giá và bỏ.

### M-2 (Medium) `assets/skills/` bị chạm sau mốc 10:00 — cần xác nhận

Ràng buộc "không đổi byte nào trong `assets/skills/`". 6 file có mtime **2026-09-04
10:05:43**:

```
assets/skills/{doc-design-flow,doc-design-system,doc-cr,doc-feature-spec,doc-plan-report,doc-test}/SKILL.md
```

`make lint-skills` sạch, nhưng lint không chứng minh nội dung không đổi, và repo chưa có
commit nào nên không diff được. Không report nào ghi lại `embed_hash` của
`dk self-check` trước Phase 11 để đối chiếu. Đề nghị: commit đầu tiên trước, hoặc ghi
`dk self-check --json` embed_hash vào report và đối chiếu ở lần sau. Hiện tại tiêu chí
"assets/skills không đổi" **chưa chứng minh được**.

### M-3 (Medium) `hook install` báo số hook không đúng thực tế

`internal/cli/hook.go:39` in `len(hook.Entries())*len(ts)` bất kể có thêm mới hay không.
Chạy lần hai in "đã cài 4 hook" dù không thêm gì. Là hành vi cũ mở rộng cho nhiều target,
nhưng với `--json` (`{"action":"đã cài","hooks":4}`) thì con số này là hợp đồng máy đọc và
đang sai. Đề nghị đếm thật (trả số mục thực thêm từ `installHooksFile`).

### M-4 (Medium) Plan phase-11 còn mô tả sai so với mã đã làm

`phase-11-target-codex.md` mục Architecture vẫn ghi `"command": ["dk","hook","run",...]`
(mảng) và matcher rộng `.*`. Mã dùng chuỗi và `apply_patch`, đúng theo sự thật đã xác minh
với Codex 0.153.2. Plan là bản ghi trạng thái, nên sửa hai câu đó để lần sau không ai
"sửa lại cho khớp plan".

### L-1 (Low) `payload.files()` nhận nhầm patch từ tool khác

`internal/hook/run.go:44`: `strings.Contains(p.ToolInput.Command, "*** Begin Patch")` bắt
mọi tool có `command` chứa chuỗi này. Hiện vô hại vì matcher cài chỉ là
`Write`/`Edit|Write`/`apply_patch`, nhưng nếu ai đó cài tay hook cho `Bash`, một lệnh
`cat <<EOF` chứa patch sẽ bị parse. Đủ để ghi một dòng comment.

### L-2 (Low) `doctor --global` vẫn dò `.codex/` theo gốc dự án

`doctor.go:92` kiểm `<root>/.codex` để tự thêm codex kể cả khi `--global`. Nên dò theo
scope đang xét (`$CODEX_HOME`/`~/.codex`) hoặc bỏ auto-detect ở scope toàn máy.

### L-3 (Low) `doctor` in 11 dòng khi target chưa cài gì

Với target hoàn toàn chưa cài, mỗi skill một dòng `chưa cài` (xem output H-2). Nên gộp
thành một dòng `skill (codex, dự án) | chưa cài 0/11` khi `current == 0`.

---

## Edge case đã probe (kết quả)

| Trường hợp | Kết quả |
|---|---|
| `*** Move to:` | **thủng** (C-1) |
| Đường dẫn có khoảng trắng trong patch | OK, deny đúng |
| Patch CRLF | OK (`TrimRight(ln, "\r")`, run.go:55) |
| Claude gửi cả `file_path` lẫn `command` | `file_path` thắng, đúng (run.go:39) |
| Deny chỉ một lần cho một patch nhiều file | OK (cờ `denied`, run.go:83) |
| Patch có 1 doc ngoài `docs/` + 1 doc mới trong `docs/` | deny đúng một lần (test `TestCodexApplyPatchDeniesNewDoc`) |
| `--target "claude,,codex"`, `--target " codex ,,claude"` | OK, trim + bỏ trùng (skill.go:32) |
| `doctor --target codex` ở dự án không có `.codex` | in đầy đủ dòng thiếu, mã 3 — OK |
| `skill status` ngoài dự án | chỉ scope toàn máy, không lỗi (skill.go:170) |
| `--target ""` | xem H-3 |

## Windows (không chạy được ở đây — Linux)

Kiểm tĩnh, chỉ liệt kê:

- `internal/cli/init_test.go:38` `exec.Command("sh", hook)` — đã guard `runtime.GOOS != "windows"`. OK.
- `internal/cli/init_test.go:32`, `phase10_test.go` bit thực thi — đã guard. OK.
- `internal/cli/phase10_test.go:73` dùng `dk.bat` trên Windows cho `LookPath`. OK.
- `internal/target/{codex,claude}_test.go:10` literal `/tmp/h`, `/p` — so sánh với
  `filepath.Join` cùng nguồn nên khớp cả hai hệ. OK.
- `internal/cli/doctor.go:179` guard bit chạy cho Windows. OK.
- Còn lại: M-1 (CRLF trong `gitx.run`). CI `windows-latest` sẽ là bằng chứng thật đầu tiên.

## Hợp đồng công khai đã đổi (cần ghi vào changelog release)

1. `doctor` item đổi tên: `skill (dự án)` → `skill (claude, dự án)`; thêm dòng
   `skill <tên> (<target>)`. JSON `doctor --json` field `item` đổi giá trị → breaking cho
   ai parse.
2. `skill status --json` thêm field `target` (thêm field, tương thích ngược).
3. `hook install --json` thêm field tùy chọn `note` (tương thích ngược).
4. `--target` giờ nhận danh sách phân tách phẩy (mở rộng, tương thích ngược).

Docs đã phản ánh 2-4 (`docs/lenh.md:10,25,26`, README:44). Mục 1 chưa thấy ghi ở đâu.

## Phong cách

Comment tiếng Việt, `fail(codeX, ...)`, cấu trúc lỗi khớp mã cũ. Refactor
`installHooksFile`/`uninstallHooksFile`/`installedHooksFile` thành package func dùng chung
là đúng hướng DRY, không tạo abstraction thừa. Không thấy scope drift, không thấy
`any` widening hay catch-and-swallow mới.

## Câu hỏi chưa giải quyết

1. `doctor --target x` = 3 là cố ý hay sót? (H-3)
2. 6 file `assets/skills/*/SKILL.md` mtime 10:05 — có thực sự không đổi nội dung? (M-2)
3. Có chấp nhận rằng agent Codex vẫn ghi file qua shell tool mà hook không thấy, và chỉ
   dựa vào pre-commit? (H-1)
