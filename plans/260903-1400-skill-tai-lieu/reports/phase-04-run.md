# Phase 4: kết quả chạy thử doc-intake, doc-cr

Ngày: 2026-09-03. Dự án mẫu `/tmp/dk-p4` (git init, `dk init`, `dk skill install doc-intake doc-cr`).

## Phần chạy được không cần Claude Code

Chuỗi lệnh mà SKILL.md quy định, chạy tay bằng `bin/dk` (build `make build`):

| Bước | Lệnh | Kết quả |
|---|---|---|
| Cài | `dk skill install doc-intake doc-cr`; `dk skill status` | 2 skill "đúng phiên bản" scope dự án; `uninstall` trả `.claude/` về rỗng |
| Idea | `dk new idea bo-loc-don-hang --set level=feature` | `docs/intake/260903-bo-loc-don-hang/idea.md`, `level: feature`, 8 mục có `chưa rõ` |
| Interview | `dk new interview ... --from .../idea.md` | cùng thư mục, `source: 260903-bo-loc-don-hang/idea.md`, chép `level` |
| Brief | `dk new brief ... --from .../idea.md --set kind=feature` | cùng thư mục; frontmatter có `level`, `kind`, `approved_by`, `approved_at`, `outcome`, `acceptance` |
| Changelog, render | `dk changelog add` x3 với `--source 260903-bo-loc-don-hang`; `dk render` x3 | 3 dòng changelog tóm tắt thật; 3 HTML trong `docs/html/intake/...` |
| Check brief `review` | `dk check docs/intake/.../brief.md` | 0 lỗi 0 cảnh báo (liên kết `./interview.md` tồn tại) |
| CR | `dk new cr doi-cach-loc --set priority=high --set requester=an` | `docs/cr/CR-260903-doi-cach-loc.md`, `requester`, `priority` trong frontmatter |
| Interview CR | `dk new interview doi-cach-loc --from docs/cr/CR-260903-doi-cach-loc.md` | `docs/cr/CR-260903-doi-cach-loc/interview.md`, `source: CR-260903-doi-cach-loc` (nhờ `beside_source` trong `types.toml`) |
| Refs | `dk new feature-spec --from brief`; CR liên kết tới spec; `dk refs docs/features/F-001-...md` | Đi: brief; Đến: CR. Đúng danh sách để điền bảng tác động |
| Chỉ mục | `dk index cr`; `dk index intake` | CR không lẫn `interview.md` (lỗi phát hiện và sửa trong lúc chạy thử); intake một dòng một thư mục theo trạng thái brief |
| Check toàn bộ | `dk check` | 0 lỗi 0 cảnh báo |

## Lỗi phát hiện và đã sửa

- `dk index cr` liệt kê `cr/<CR>/interview.md` như một CR. Sửa: chỉ mục chỉ nhận loại có `dir` trùng thư mục; thêm test.
- Mô tả kích hoạt có dấu hai chấm không quote làm frontmatter SKILL.md không parse. Sửa: quote đơn; `content_test.go` bắt lỗi này (name phải trùng thư mục).
- Hai câu phủ định "Không dùng cho ..." cố ý lặp từ khóa của skill kia nên luôn trùng 3 từ; test trigram bỏ phần từ "Không dùng" trở đi.

## Chưa chạy được ở phiên này

Bước 7 của phase yêu cầu chạy thật trong Claude Code với một người trả lời phỏng vấn:
ý tưởng có 2 trường "chưa rõ" → đúng 2 câu hỏi → brief `review`; CR trên spec mẫu →
bảng tác động khớp `dk refs`. Phiên này tự động, không có người đối thoại nên phần
đó chưa kiểm. Điều cần quan sát khi chạy thật (từ Risk Assessment của phase):

- AI có hỏi đúng tên mục ("Mục 2 Ai gặp: ...") thay vì danh sách cố định không; bảng
  "Sai / Đúng" trong SKILL.md doc-intake là chỗ thêm ví dụ nếu vỡ.
- AI có tự đặt `approved` hoặc sửa Feature Spec trước khi người duyệt bảng không; nếu
  có, thêm quy tắc `cr-approval-order` vào `dk check` (spec `updated` mới hơn CR nguồn
  đang `review`).
- Câu hỏi giao diện có kèm 2 đến 3 phương án tham chiếu không.

Ghi lại vào file này sau lần chạy thật, rồi sửa SKILL.md và chạy lại một lần.
