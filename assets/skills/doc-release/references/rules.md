# Quy tắc doc-release

Chép từ báo cáo `plans/reports/consult-260903-1032-tai-lieu-du-an-nguoi-va-ai.md`.
Chỉ quy tắc, không giải thích; giải thích đọc ở mục báo cáo ghi cuối file.

## Nguyên tắc nền

- Hai loại brief, không nhầm: Intake brief là đầu vào (trước Feature Spec);
  Release brief là đầu ra (giới thiệu tính năng cho người dùng, sau khi xong).
- Release brief và User guide là tài liệu thẩm quyền: luôn đúng với sản phẩm đang
  chạy, sửa tại chỗ. Release notes là tài liệu trạng thái: một phiên bản một file,
  không sửa sau khi `published`.
- Nguồn sự thật là Feature Spec; Release brief chép từ đó rồi viết lại, không thêm
  hành vi spec không có.
- Mỗi tài liệu có một `owner`. Không có người chịu trách nhiệm thì không viết.

## Tầng 7: Người dùng cuối

| Tài liệu | Cách dùng | Loại | Người | AI |
|---|---|---|---|---|
| Release brief | sinh từ mục Mục đích, Hành vi, Giao diện của Feature Spec; viết lại bằng ngôn ngữ người dùng, có ảnh | thẩm quyền | duyệt giọng văn, thêm ảnh | soạn nháp |
| Release notes | gộp Release brief và sửa lỗi một đợt phát hành | trạng thái | duyệt | tổng hợp |
| User guide | ghép dần từ Release brief; theo nhiệm vụ, không theo menu | thẩm quyền | duyệt | soạn, cập nhật |
| FAQ | viết khi câu hỏi lặp lần thứ hai | thẩm quyền | duyệt | soạn từ ticket |

## Release brief

- Chỉ từ Feature Spec `implemented`, hoặc `approved` sắp phát hành do người chỉ.
- Một tính năng một file, tên theo mã spec (`docs/release/briefs/<F-xxx>.md`);
  `feature` và `source` là mã spec; `kind: feature | fix`.
- Nội dung: giúp gì, ai dùng, bước cơ bản, ảnh màn hình, giới hạn đã biết. Không
  thuật ngữ kỹ thuật, không mã bước, không tên role hay bảng. `dk check
  no-jargon` cảnh báo từ trong `dk.toml` `[release] jargon` (mặc định: API,
  endpoint, database, migration, backend, frontend, JSON). Từ nằm trong tên sản
  phẩm: người quyết định bỏ khỏi danh sách, skill không tự sửa.
- Giới hạn đã biết lấy từ mục 11 và ngoại lệ mục 7 chưa xử lý; nói thẳng, không hứa.
- Ảnh: liên kết mockup đã duyệt hoặc ảnh chụp thật do người thêm.
- `released_in` do `dk new release-notes --collect` ghi; không sửa tay.

## Release notes

- `dk new release-notes --collect <phiên bản>`: gom mọi brief `status: ready` chưa
  có `released_in`, nhóm Mới (`kind: feature`) và Sửa lỗi (`kind: fix`); ghi
  `released_in` và một dòng changelog nguồn `<phiên bản>` cho từng brief.
- Mỗi dòng một câu cho người dùng, liên kết về brief. Mục "Cần lưu ý" ghi điều
  người dùng phải làm khác đi.
- Người duyệt rồi `published`; phiên bản sau viết file mới.

## User guide

- Một nhiệm vụ người dùng một trang trong `docs/release/guide/`; `task` là nhóm
  nhiệm vụ để `dk index user-guide` xếp mục lục. Không tổ chức theo menu.
- Ghép dần từ Release brief sau mỗi lần phát hành; trang đã có thì sửa đúng mục,
  không tạo trang trùng.
- Liên kết FAQ thay vì lặp câu trả lời. Chịu `no-jargon` như Release brief.

## FAQ

- Một file `docs/release/faq.md`, nối bằng `--append "<câu hỏi> | <trả lời> |
  <liên kết>"`; chỉ thêm khi câu hỏi lặp lần thứ hai và người xác nhận.
- Không sửa dòng cũ; trả lời đổi thì nối dòng mới, người xóa dòng cũ.

## Tách file và ngưỡng dòng

- Release brief một tính năng một file; Release notes một phiên bản một file;
  User guide mỗi nhiệm vụ một trang, mục lục riêng; FAQ một file.
- Vượt 300 đến 500 dòng thì xem lại; trần cứng 800 dòng (`line-threshold`).
- Frontmatter thống nhất: `type`, `title`, `status`, `owner`, `created`,
  `updated`, `source`; brief thêm `feature`, `kind`, `released_in`; guide thêm
  `task`; notes thêm `version`. `dk new` điền, skill không sửa `created`,
  `created_by`, `dk_version`.

## Khi nào không áp dụng

- Dự án dưới 5 tính năng hoặc dưới 2 tháng: bỏ Release brief, giữ `CHANGELOG.md`.
- Công cụ nội bộ không có người dùng ngoài: bỏ Tầng 7, `CHANGELOG.md` là đủ.
- Dùng Notion hoặc Confluence cho tài liệu người dùng: đơn vị là trang, ngưỡng dòng
  không áp dụng; vẫn giữ một nhiệm vụ một trang.

## Trạng thái (khóa tiếng Anh trong `types.toml`)

| Khóa | Nghĩa |
|---|---|
| `draft` | đang soạn |
| `review` | chờ người duyệt giọng văn |
| `ready` | brief sẵn sàng gom vào Release notes |
| `approved` | user guide đã chốt |
| `published` | release notes đã phát hành, không sửa |
| `active` | FAQ đang dùng |

Giải thích: báo cáo mục 1 (hai loại brief), mục 2 (Tầng 7), mục 7 (Release brief,
User guide), mục 10 (dự án nhỏ, công cụ nội bộ).
