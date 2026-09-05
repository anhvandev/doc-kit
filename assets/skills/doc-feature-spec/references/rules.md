# Quy tắc doc-feature-spec

Chép từ báo cáo `plans/reports/consult-260903-1032-tai-lieu-du-an-nguoi-va-ai.md`.
Chỉ quy tắc, không giải thích; giải thích đọc ở mục báo cáo ghi cuối file.

## Nguyên tắc nền

- Một nguồn sự thật cho mỗi tính năng: Feature Spec. Các tài liệu khác (mockup, test,
  release brief) sinh ra hoặc tham chiếu từ đó.
- Feature Spec là tài liệu thẩm quyền. Người viết mục đích, hành vi, quy tắc và duyệt
  toàn bộ; AI soạn ngoại lệ, tiêu chí chấp nhận, kiểm tra thiếu sót.
- AI không tự sửa mục hành vi của Feature Spec đã duyệt nếu chưa có người duyệt (qua CR).
- Mỗi tài liệu có một `owner`. Không có người chịu trách nhiệm thì không viết.

## Cấu trúc 10 mục

Frontmatter (không đánh số): mã, tên, trạng thái, chủ sở hữu, liên kết brief hoặc CR
(`source`). Thân:

1. Mục đích và giá trị cho người dùng.
2. Tác nhân và điều kiện tiên quyết.
3. Sơ đồ luồng chính: Mermaid nhúng, mỗi nút mang mã bước `B1`, `B2`...
4. Bảng hành vi theo mã bước: mỗi dòng một cặp hành động của tác nhân và phản hồi
   quan sát được của hệ thống, đúng mã bước trong sơ đồ.
5. Giao diện: mỗi mã bước liên kết đến mockup và trạng thái hiển thị. Chỉ liên kết,
   không chép ảnh. Bỏ khi không có giao diện (`has_ui: false`).
6. Luồng thay thế và ngoại lệ: text hoặc sơ đồ nhỏ riêng, không nhồi vào sơ đồ chính.
7. Quy tắc nghiệp vụ: tham chiếu bằng mã, không lặp lại giữa các spec.
8. Tiêu chí chấp nhận dạng Given / When / Then.
9. Dữ liệu và API liên quan: chỉ liên kết.
10. Ngoài phạm vi.

Sơ đồ, bảng hành vi, tiêu chí chấp nhận của một spec ở cùng file vì ràng buộc bằng
mã bước; không tách theo mục con. `dk check` quy tắc `spec-section-order`: tiêu đề
cấp 2 `## 1.` đến `## 10.` đúng thứ tự, mỗi số một lần; `has_ui: false` bỏ 5;
`format: crud` bỏ 3.

## Sơ đồ và mã bước

- Dạng text Mermaid, không dạng ảnh.
- Sơ đồ là mục lục, bảng là diễn giải từng nút. Mã bước ràng buộc hai phần; User flow
  ở họ Design dùng chung mã bước này. `dk check` quy tắc `step-codes`: tập mã trong sơ
  đồ bằng tập mã ở cột đầu bảng hành vi.
- Mã bước: `B` + số, hậu tố chữ thường cho bước chèn giữa (`B2a`, `B2b`). Không đánh
  số lại khi chèn hoặc bỏ bước, vì mockup và test đặt tên theo mã.
- Một sơ đồ cho luồng chính. Luồng tuyến tính dưới bốn bước hoặc CRUD đơn giản: bỏ sơ
  đồ, bảng hành vi là đủ.

| Bản chất tính năng | Loại sơ đồ |
|---|---|
| Chuỗi hành động có rẽ nhánh | flowchart |
| Nhiều thành phần tương tác theo thứ tự | sequenceDiagram |
| Xoay quanh vòng đời trạng thái của một đối tượng | stateDiagram |

## Định dạng theo bối cảnh

| Bối cảnh | Định dạng | `format` |
|---|---|---|
| Team nhỏ, tài liệu sống trong repo (mặc định) | Feature Spec Markdown | `spec` |
| Nghiệp vụ phức tạp, nhiều tác nhân, nhiều ngoại lệ | Use Case Specification + Sequence diagram | `use-case` |
| Team Agile, muốn tài liệu chạy được thành test | User Story + Acceptance Criteria dạng Gherkin | `story` |
| CRUD đơn giản | Bảng field và quyền, không viết use case | `crud` |
| Tính năng có trạng thái phức tạp | Bổ sung State diagram | `state` |
| Outsource, cần bàn giao và nghiệm thu | Functional Spec đầy đủ + Test Case Specification | `spec` với mục 6, 8 đầy đủ; test ở họ Test |

## Nguyên tắc viết

- Không viết "hệ thống xử lý" mà không nói kết quả nhìn thấy được.
- Ngoại lệ và trường hợp biên chiếm ít nhất một nửa giá trị tài liệu: số dòng mục 6
  không dưới nửa số dòng bảng mục 4.
- Kết quả mong đợi phải kiểm chứng được: thông báo cụ thể, dữ liệu thay đổi thế nào,
  sự kiện nào phát ra.
- Dùng bảng và sơ đồ thay cho văn xuôi dài.
- Mục 4, 6, 8 nuôi Test case. Mục 5 nuôi Test giao diện. Mục 1, 2, 4, 5, 10 nuôi Release brief.
- Thuật ngữ nghiệp vụ in đậm ở mục 1, 4, 7 phải có trong Glossary (`dk check` quy tắc
  `glossary-term`).

## Luồng một chiều

- brief `approved` cấp tính năng mới đi tiếp sang Feature Spec. brief chưa duyệt thì không.
- Sửa hành vi thì sửa Feature Spec trước, rồi cập nhật tài liệu con (mockup, test,
  release brief). Sửa giao diện thì sửa mockup trước, rồi cập nhật spec.
- Không sửa test hay brief mà bỏ qua spec.
- Thay đổi trên spec đã duyệt đi qua CR: bảng tác động được người duyệt, quyết định
  ghi, rồi mới sửa spec đúng mã bước nêu ở "Cần làm gì". `dk check` quy tắc
  `cr-approval-order`: spec có `source` là CR chưa chốt mà `updated` mới hơn CR là lỗi.
- Spec phải liên kết ngược (`source`) về brief hoặc CR sinh ra hoặc sửa nó gần nhất;
  nguồn cũ ghi ở "Lịch sử nguồn" cuối file.

## Khi nào nới quy tắc

- Luồng tuyến tính dưới bốn bước hoặc CRUD: bỏ sơ đồ, bảng hành vi là đủ.
- Không có giao diện người dùng (API, CLI, batch job): bỏ mục 5, giữ giọng văn thông
  báo lỗi trong Glossary.
- Yêu cầu thay đổi hàng ngày: giữ spec ở mức frontmatter, mục 1, 8; chi tiết hành vi để trong test.
- Dự án dưới 5 tính năng: vẫn một tính năng một file.

## Tách file và frontmatter

- Một tính năng một file trong `features/`, tên file mang mã `F-xxx`. Feature catalog
  (`docs/features/README.md`) sinh bằng `dk index features`, không viết tay.
- Frontmatter thống nhất: mã, tên, trạng thái, chủ sở hữu, ngày cập nhật.
- File vượt 500 dòng `dk check` cảnh báo; trần cứng 800 dòng. Dài quá là đang gộp
  nhiều tính năng: tách spec.

## Trạng thái và trường (khóa trong frontmatter)

- `status`: `draft`, `review` (chờ duyệt), `approved` (đã duyệt), `implemented` (đã
  phát hành), `deprecated` (bãi bỏ).
- `source`: `<yymmdd>-<slug>/brief.md` khi tạo từ brief; mã CR sau khi sửa theo CR.
- `purpose`, `acceptance`: chép từ `outcome`, `acceptance` của brief.
- `has_ui`: `true` mặc định; `false` bỏ mục 5. `format`: `spec`, `use-case`, `story`,
  `crud`, `state`; đặt lúc `dk new`, không đổi sau.
- `owner`: người chịu trách nhiệm spec, bắt buộc.

## Ghi file

- Mọi file trong `docs/` tạo bằng `dk new`. `dk` không chạy được thì dừng và báo người
  cài; không tạo file tay.
- Mỗi lần sửa spec có một dòng changelog với tóm tắt thật (`dk changelog add`);
  `--source` là thư mục intake của brief hoặc mã CR.
- Không sửa `created`, `created_by`, `dk_version`.

## Quy tắc → mục báo cáo

| Quy tắc | Mục báo cáo |
|---|---|
| Nguyên tắc nền | 1, Tầng 4 (cột Người, AI) |
| Cấu trúc 11 mục, sơ đồ và mã bước, nguyên tắc viết | 3 |
| Định dạng theo bối cảnh | Tầng 4 (bảng "định dạng thay thế") |
| Luồng một chiều | 4 (đoạn "Mũi tên chỉ đi một chiều"), Tầng 4 |
| Khi nào nới quy tắc | 10 (dòng luồng tuyến tính, không giao diện, thay đổi hàng ngày, dưới 5 tính năng) |
| Tách file và frontmatter | 7 (dòng "Feature Spec", "Quy tắc đi kèm", ngưỡng dòng) |
| Trạng thái và trường, ghi file | plan `plans/260903-1400-skill-tai-lieu/plan.md` mục 1.1, 1.3, 1.5, 1.10; `assets/types.toml` |
