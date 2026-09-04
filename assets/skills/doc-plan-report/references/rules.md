# Quy tắc doc-plan-report

Chép từ báo cáo `plans/reports/consult-260903-1032-tai-lieu-du-an-nguoi-va-ai.md`.
Chỉ quy tắc, không giải thích; giải thích đọc ở mục báo cáo ghi cuối file.

## Nguyên tắc nền

- Mỗi tài liệu tồn tại để trả lời một câu hỏi: Plan trả lời "cần làm gì", Report
  trả lời "đã làm gì", Roadmap trả lời "làm gì trước", Decision log trả lời "đã
  quyết gì nhỏ", CHANGELOG trả lời "phiên bản này đổi gì cho người dùng".
- Tài liệu thẩm quyền (Roadmap) người chốt, AI chỉ đề xuất. Tài liệu trạng thái
  (Plan, Report, Decision log, CHANGELOG) AI soạn, người duyệt phần được ghi.
- Mỗi tài liệu có một `owner`. Không có người chịu trách nhiệm thì không viết.
- AI không tự sửa tài liệu thẩm quyền hoặc report đã nộp.

## Tầng 5: Kế hoạch và tiến độ

| Tài liệu | Cách dùng | Loại | Người | AI |
|---|---|---|---|---|
| Roadmap | cập nhật theo quý hoặc đợt phát hành; một file | thẩm quyền | chốt | không |
| Plan | một thư mục một đợt việc; chốt phạm vi trước khi bắt đầu; AI làm theo plan | trạng thái | duyệt phạm vi và non-goals | soạn giai đoạn và bước |
| Report | viết ngay khi kết thúc giai đoạn, kèm bằng chứng; không sửa sau khi nộp | trạng thái | đọc, xác nhận | viết |
| Decision log | ghi theo ngày, một dòng một quyết định, gồm quyết định giao diện | trạng thái | ghi | ghi khi được ủy quyền |
| CHANGELOG | tổng hợp từ Report và commit; Keep a Changelog; ngôn ngữ người dùng | trạng thái | duyệt | tổng hợp |

## Plan

- Cấu trúc: mục tiêu, phạm vi, non-goals, giai đoạn, tiêu chí hoàn thành, rủi ro,
  liên kết report. Chỉ mục ngắn; mỗi giai đoạn một file `phase-NN-<slug>.md`.
- Thư mục `plans/<yymmdd>-<hhmm>-<slug>/`; `dk.toml` `plans_dir` trỏ nơi khác nếu
  repo có quy ước riêng. Repo đã có công cụ plan riêng: không tạo plan, chỉ viết
  Report và Decision log theo thư mục có sẵn.
- Người duyệt mục Phạm vi và Non-goals trước khi AI làm việc theo plan
  (`status: approved`). Đổi phạm vi giữa chừng: sửa plan, người duyệt lại.
- File phase chỉ chứa thứ cần để thực hiện an toàn: bối cảnh, yêu cầu, file, bước,
  kiểm chứng, rủi ro và rollback. Không chép lại plan.md.
- Trạng thái plan: `draft`, `approved`, `in-progress`, `done`, `cancelled`. Phase:
  `todo`, `in-progress`, `done`, `skipped`.
- Plan và report là bản ghi trạng thái; giai đoạn xong không biến chúng thành tài
  liệu thẩm quyền. Quyết định bền vững chuyển sang ADR hoặc Decision log.

## Report

- Một giai đoạn một file, theo ngày giờ và slug, không gộp. Viết ngay khi kết thúc
  giai đoạn.
- Bốn phần: đã làm, bằng chứng, chưa xong, câu hỏi mở. Bằng chứng là mã commit,
  đường dẫn output test hoặc khối kết quả thật; `dk check` quy tắc
  `report-evidence` cảnh báo khi thiếu.
- Không sửa sau khi `submitted`; sai thì viết report mới trỏ về report cũ.
- Câu hỏi mở đặt cuối; AI không tự trả lời trong report.
- Plan và report không ghi `CHANGELOG-DOCS.md` nhưng `dk check` vẫn kiểm liên kết
  và frontmatter.

## Decision log

- Một dòng một quyết định: `- <ngày> | <quyết định> | <người chốt> | <liên kết>`.
  Mới nhất ở cuối. Không sửa dòng cũ; đổi ý thì thêm dòng mới nhắc dòng cũ.
- Gồm quyết định giao diện (màu, vị trí nút, chữ) khi được ủy quyền.
- Quyết định lớn (kiến trúc, thư viện, hợp đồng dữ liệu, chi phí) là ADR, không ghi
  ở đây.
- Cắt theo quý khi dài.

## CHANGELOG sản phẩm

- Theo Keep a Changelog: mới nhất ở trên, `## [Unreleased]` trên cùng, mỗi phiên
  bản một mục `## [x.y.z] - yyyy-mm-dd`, nhóm Added, Changed, Deprecated, Removed,
  Fixed, Security.
- Ngôn ngữ hướng người dùng: người dùng làm được gì, không nói file hay hàm.
- Tổng hợp từ Report và commit; người duyệt trước khi phát hành.
- Vượt 500 dòng thì xem lại, trần cứng 800 dòng (`line-threshold`): cắt theo năm
  hoặc phiên bản lớn.
- Công cụ nội bộ không có người dùng ngoài: CHANGELOG là đủ, không cần Release
  notes.

## Roadmap

- Toàn cảnh, một file. Thứ tự ưu tiên tính năng theo mốc thời gian.
- Người chốt; skill chỉ đề xuất sắp xếp từ Feature catalog
  (`docs/features/README.md`, sinh bằng `dk index features`).
- Dự án dưới 5 tính năng hoặc dưới 2 tháng: bỏ Feature catalog; Roadmap vẫn có
  thể là một bảng ngắn.

## Tách file và ngưỡng dòng

- Plan: một đợt việc một thư mục, mỗi giai đoạn một file. Report: một giai đoạn
  một file. Decision log, Roadmap, CHANGELOG: một file.
- Vượt 300 đến 500 dòng thì xem lại; trần cứng 800 dòng (`dk check`
  `line-threshold` cho file trong `docs/`).
- Frontmatter thống nhất: `type`, `title`, `status`, `owner`, `created`,
  `updated`, `source`. `dk new` điền; skill không sửa `created`, `created_by`,
  `dk_version`.

## Khi nào không áp dụng

- Team hoàn toàn không dùng AI: Report thay bằng ghi chú trong công cụ quản lý
  task.
- Dùng Notion hoặc Confluence: đơn vị là trang, ngưỡng dòng không áp dụng.
- Repo đích có công cụ plan riêng: tôn trọng quy ước tên có sẵn; chỉ viết Report
  và Decision log.

## Trạng thái (khóa tiếng Anh trong `types.toml`)

| Khóa | Nghĩa |
|---|---|
| `draft` | đang soạn |
| `review` | chờ duyệt (roadmap) |
| `approved` | người đã duyệt phạm vi (plan) hoặc chốt (roadmap) |
| `in-progress`, `todo`, `done`, `skipped`, `cancelled` | tiến độ plan và phase |
| `submitted` | report đã nộp, không sửa |
| `active`, `archived` | decision log đang ghi hoặc đã cắt |

Giải thích: báo cáo mục 2 (Tầng 5), mục 7 (Tách file), mục 10 (Khi nào khuyến
nghị không còn đúng).
