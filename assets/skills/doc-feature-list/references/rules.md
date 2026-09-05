# Quy tắc doc-feature-list

Chép từ báo cáo `plans/reports/brainstorm-260904-2328-skill-feature-list.md` và
`plans/reports/consult-260903-1032-tai-lieu-du-an-nguoi-va-ai.md`. Chỉ quy tắc,
không giải thích.

## Nguyên tắc nền

- Feature list là tài liệu định hướng cấp dự án: trả lời "sản phẩm dự kiến có những
  tính năng gì, cái nào trước". Không phải Feature Spec; không mô tả hành vi.
- Nguồn sự thật của list là Product brief đã duyệt (`source`). Mỗi dòng phải dẫn
  được về một mục của brief hoặc idea; không dẫn được thì vào mục 4 Chưa rõ.
- Người chốt (`approved`). AI soạn nháp, dừng ở `review`. Sau khi chốt chỉ sửa theo
  CR; riêng cột Spec điền được không cần CR.

## Feature list

- Một file cho một sản phẩm: `docs/overview/feature-list.md`, dưới 200 dòng
  (`dk check` cảnh báo). Đã có file thì cập nhật, không tạo file thứ hai.
- Bảng mục 2 đủ 7 cột: Mã, Tên, Mô tả, Nhóm, Ưu tiên, Nguồn, Spec.
- Mã tạm `FL-01`, `FL-02` cấp một lần, không đánh số lại, không dùng lại. Dòng bỏ
  ghi ưu tiên `won't` và giữ nguyên chỗ.
- Ưu tiên theo MoSCoW: `must` (có tiêu chí chấp nhận trong brief §4), `should` (brief
  §1 nhắc), `could` (chỉ idea nhắc), `won't` (brief §3 loại). Không thêm mức khác.
- Mô tả một câu, có hành động của người dùng và kết quả nhìn thấy.
- Không quá 25 dòng; hơn thì gom dòng cùng việc thành một dòng.
- Không ước lượng effort, không ghi mốc thời gian, không xếp thứ tự trong nhóm ngoài
  thứ tự xuất hiện ở brief.

## Quan hệ với tài liệu khác

- Product overview mục 3 dùng cùng tên nhóm với mục 3 của list; list là nơi liệt kê
  đầy đủ, overview chỉ một dòng một nhóm.
- Roadmap tham chiếu `FL-xx` khi tính năng chưa có `F-xxx`; có spec rồi thì dùng
  mã spec.
- Feature catalog (`docs/features/README.md`) là bản sinh từ Feature Spec, không thay
  list và không được chép vào list.
- Mỗi dòng list đi tiếp bằng intake cấp tính năng (idea, phỏng vấn nếu cần, brief),
  rồi mới có Feature Spec. List không rút ngắn chuỗi đó.

## Không làm

- Không tạo idea, brief, Feature Spec, roadmap, Product overview từ list.
- Không tự đặt `approved`, không điền `approved_by`.
- Không bịa tính năng brief và idea không nhắc, kể cả tính năng phổ biến.
- Không sửa `created`, `created_by`, `dk_version` trong frontmatter.
