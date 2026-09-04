# Code review phase 1: CLI lõi `dk`

- Ngày: 2026-09-03
- Reviewer: subagent `code-reviewer`, chạy binary thật trong repo tạm; không sửa file.
- Kết quả cổng: `go vet`, `gofmt`, `go test ./...` sạch; binary tĩnh; 14 trường hợp mã thoát đúng; frontmatter giữ thứ tự khóa, CRLF, thân; số dòng changelog khớp `git diff --numstat HEAD` ở hai lần sửa; gộp 10 phút đúng.
- Coverage: changelog 91.7, cli 81.8, docs 84.0, doctype 78.4, frontmatter 80.0, gitx 77.8, tmpl 71.4 (%).

## Phát hiện và xử lý

| # | Mức | Vấn đề | Xử lý |
|---|---|---|---|
| 1 | cao | `git status --porcelain` trả tên file có dấu dạng octal, `pending` không bao giờ sạch với tên tiếng Việt | Sửa: `-c core.quotePath=false status -z`; test `gitx` với `bộ lọc đơn.md` |
| 2 | cao | 18 file `.gitkeep` do `init` tạo bị tính pending; `add .gitkeep` được chấp nhận | Sửa: bỏ file ẩn khỏi changelog; test `pending` sạch ngay sau `init` |
| 3 | cao | Tiêu đề có dấu hai chấm làm hỏng frontmatter khi render qua text template | Sửa: `title`, `owner`, `source` đặt qua `yaml.Node` sau render; test `Sửa lỗi: mất đơn` |
| 4 | vừa | `--from` nguồn ngoài thư mục loại vẫn đặt file cạnh nguồn (`docs/cr/brief.md`) | Sửa: chỉ đặt cạnh nguồn khi nguồn nằm trong `dir` của loại; test |
| 5 | vừa | `--summary`, `--source` chứa xuống dòng hoặc ` \| ` làm hỏng file changelog | Sửa: từ chối, mã 2; test |
| 6 | vừa | stdout và JSON của `add` in mục trước khi gộp | Sửa: `File.Add` trả mục đã ghi; test |
| 7 | vừa | `init --force` ghi đè `dk.toml` bằng mặc định | Sửa: nạp và giữ cấu hình cũ; test |
| 8 | thấp | `init` trong dự án có sẵn tạo dự án lồng | Sửa: từ chối khi thư mục cha có `dk.toml`; test |
| 9 | thấp | slug sai trả mã 1 thay vì 2 | Sửa: kiểm tra ở CLI, mã 2; test |
| 10 | thấp | `owner` rỗng ghi thành `owner:` (null) | Không tái hiện với mã hiện tại: `SetString` tag `!!str` ghi `owner: ""`; reviewer chạy bản trước khi sửa mục 3 |
| 11 | thấp | goldmark khai báo sớm bị `go mod tidy` xóa | Sửa: bỏ khỏi `go.mod`, phase 2 thêm khi dùng |
| 12 | thấp | Gộp cả khi mục mới có giờ cũ hơn dòng đầu | Sửa: yêu cầu delta không âm; test |
| 13 | thấp | Mỗi `add` gọi git khoảng 6 lần | Bỏ qua: không ảnh hưởng hợp đồng, xét lại nếu chậm thấy được |

Sau sửa: `make vet test build` qua, binary tĩnh.
