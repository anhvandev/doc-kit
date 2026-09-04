# Phase 5: kết quả chạy thử doc-overview, doc-adr

Ngày: 2026-09-03. Dự án mẫu `/tmp/dk-p5.*` (git init, `dk init`, `dk skill install doc-overview doc-adr`), binary `bin/dk` từ `make build`.

## Chuỗi lệnh mà SKILL.md quy định, chạy tay

| Bước | Lệnh | Kết quả |
|---|---|---|
| Cài | `dk skill install doc-overview doc-adr` | 2 skill "đã cài"; `skill status` 8 dòng (4 skill x 2 scope) |
| Brief | `dk new brief san-pham --set level=project --set kind=product --set status=approved` | brief `approved` cấp dự án |
| Product overview | `dk new product-overview quan-ly-don --from <brief> --set owner=an` | `docs/overview/product-overview.md`, 57 dòng; `source`, `title`, `outcome`, `acceptance` chép từ brief; `dk check` 0 lỗi |
| Catalog | 2 `dk new feature-spec --from <brief>`; `dk index features` | `docs/features/README.md` bảng phẳng 2 dòng F-001, F-002, đúng cột `Mã | Tên | Trạng thái | Chủ sở hữu | Brief hoặc CR nguồn | Cập nhật` |
| Glossary | `dk new glossary`; thêm dòng "Đơn hàng"; spec F-001 mục 2 có **đơn hàng**, **kho** | `dk check`: 1 cảnh báo `glossary-term` cho **kho** đúng dòng 22; **đơn hàng** không báo |
| ADR chốt | `dk new adr postgres`; đổi `status: accepted`; commit | ADR-0001 trong HEAD |
| Thay thế | `dk new adr sqlite --set supersedes=ADR-0001`; ADR-0001 chỉ đổi `status: superseded`, `superseded_by: ADR-0002`; `dk index adr` | chỉ mục 2 dòng theo số, cột Thay thế của ADR-0002 là ADR-0001; `dk check` 0 lỗi |
| Sửa thân ADR-0001 | `echo >> ADR-0001`; `dk check` | lỗi `adr-immutable`, mã thoát 3 |
| Changelog, render | `dk changelog add` x5, `dk render` x5, `dk index all` | HTML trong `docs/html/overview/`, `docs/html/adr/`; `dk check` toàn bộ chỉ còn lỗi `./interview.md` của brief mẫu (template brief phase 4, skill xóa mục đó khi không phỏng vấn) |

## Lỗi phát hiện và đã sửa

- Template `product-overview` và `architecture` liên kết sẵn đến `architecture.md`, `glossary.md`, `../adr/README.md`, `../features/README.md`; file chưa có thì `dk check` báo `link-broken` ngay sau `dk new`. Sửa: bỏ liên kết cứng, thay bằng gợi ý trong chú thích.
- Chạy thử lần đầu commit ADR bị pre-commit chặn (thiếu changelog) nên `adr-immutable` không có HEAD để so; kịch bản thật phải `dk changelog add` trước khi commit, đúng như SKILL.md yêu cầu.
- Chữ in đậm ngoài mục 2, 5, 8 không được quét (cố ý, tránh **Given/When/Then** ở mục 9); SKILL.md và rules.md ghi rõ.

## Sau code review

Sửa M1, M2, M4 và hai Low của `reports/code-review-260903-phase-05-skill-overview-adr.md` (gộp mọi Glossary, thứ tự ổn định, một tiến trình git mỗi ADR, bỏ khối mã, gộp hoa thường); test bổ sung ở `internal/check/testdata`. M3 để chủ dự án quyết (xem plan mục 1.11).

## Chưa chạy được ở phiên này

Phần AI thật trong Claude Code: skill dừng khi brief chưa `approved`, hỏi từng câu một khi soạn ADR, quét mã nguồn thật để ra Architecture overview với dòng `chưa xác nhận`. Gộp vào đợt chạy thử toàn vòng phase 10 cùng bước 9 phase 3 và bước 7 phase 4. Điều cần quan sát:

- Với brief `level: feature` hoặc `status: review`, skill có dừng không tạo file không.
- Architecture overview từ mã: cột "Xác nhận" có giữ `chưa xác nhận` cho đến khi người trả lời không.
- Khi thay thế ADR, skill có chỉ sửa frontmatter ADR cũ không; nếu sửa thân, `dk check` bắt được nhưng skill cần tự khôi phục từ git.
