# Quy tắc khung HTML và CSS của `dk render`

Áp cho `page.html` và `style.css` trong thư mục này. Test tự động kiểm một phần; phần còn lại kiểm bằng mắt khi đổi.

## Bắt buộc

1. **Một file, không tài nguyên ngoài.** Không `<link>`, không `@import`, không `url(http…)`, không font tải về. CSS nhúng trong `<style>`, script Mermaid nhúng trong `<script>` và chỉ khi trang có khối `mermaid`.
2. **Khai báo font rõ ràng.** Mọi vùng chữ có `font-family` khai báo; không dựa vào font mặc định của trình duyệt. Tiêu đề dùng bộ serif, thân dùng bộ sans, mã bước và mã tài liệu dùng bộ mono. Chỉ dùng font có sẵn phổ biến trên Linux, macOS, Windows; luôn kết thúc bằng họ generic.
3. **Tương phản tối thiểu 4.5:1** cho chữ thường, 3:1 cho chữ lớn và đường viền có nghĩa. Màu khai báo bằng biến ở `:root`, không rải mã màu trong quy tắc.
4. **Bố cục bất đối xứng có chủ đích.** Cột nội dung giới hạn khoảng 70 ký tự, không căn giữa toàn trang; metadata đứng ở lề phải như cột ghi chú. Dưới 60rem thì xếp dọc, metadata lên trước thân.
5. **Sơ đồ và bảng đọc được.** Bảng có đường kẻ ngang, dòng chẵn lệch nền nhẹ; ô mã bước (`id="step-…"`) khi là đích neo (`:target`) được tô nền nổi. Sơ đồ Mermaid không bị cắt: khối `pre.mermaid` cuộn ngang khi cần.
6. **Chuyển động tôn trọng người dùng.** Mọi transition, animation nằm trong `@media (prefers-reduced-motion: no-preference)`.
7. **In được.** `@media print`: một cột, bỏ nền tối, bỏ sticky, liên kết ngoài in kèm địa chỉ, bảng không vỡ dòng giữa trang.

## Cấm

- Gradient trong chữ: không `background-clip: text`, không `linear-gradient` trong `color` hoặc `background` của chữ. Gradient chỉ được dùng cho nền trang trí không chứa chữ, và hiện tại không dùng.
- Emoji hoặc ký tự trang trí trong tiêu đề và nhãn.
- Bóng đổ nhiều lớp, viền bo tròn lớn hơn 6px, hiệu ứng "kính mờ".
- `!important`, trừ trong `@media print`.
- Chữ nhỏ hơn 0.8rem cho nội dung đọc.

## Cách kiểm

- Test `internal/render` đọc `style.css` và thất bại nếu thấy `background-clip: text`, `-webkit-background-clip: text`, `@import`, `url(http`, hoặc `linear-gradient` trong khai báo `color:`.
- Kiểm bằng mắt: mở một Feature Spec có sơ đồ, bảng hành vi, chú thích ở chế độ offline; bấm mã bước dưới sơ đồ phải nhảy đúng dòng bảng và dòng được tô.
