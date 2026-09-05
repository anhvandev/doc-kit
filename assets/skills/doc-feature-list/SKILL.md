---
name: doc-feature-list
description: 'Feature list trong docs/overview: tách brief sản phẩm đã được duyệt thành bảng tính năng dự kiến có mã tạm, nhóm, ưu tiên MoSCoW, dẫn về mục brief; dừng chờ duyệt; điền mã spec khi tính năng có Feature Spec. Không dùng cho brief, spec, roadmap.'
---

# doc-feature-list

Soạn và cập nhật Feature list, một file `docs/overview/feature-list.md` cho một sản
phẩm, bằng `dk new feature-list` và `dk check`.

**Dừng bắt buộc: chỉ soạn từ brief `level: project`, `kind: product`,
`status: approved`; list chỉ đến `status: review`, không tự đặt `approved`; không
tạo idea, brief hay Feature Spec từ dòng list. Mã tạm `FL-xx` không đánh số lại;
dòng bỏ thì ghi `won't`, không xóa.**

Mọi file trong `docs/` chỉ được tạo bằng `dk new`. Không tự tạo file tay, kể cả khi
`dk` lỗi. Quy tắc nội dung ở `references/rules.md`; đọc trước khi làm.

## Phạm vi

- Làm: Feature list từ Product brief đã duyệt; điền cột Spec khi tính năng đã có
  Feature Spec; sửa list theo lời người khi chưa chốt, theo CR khi đã chốt.
- Không làm: brief và phỏng vấn (họ Intake); Feature Spec (họ Feature Spec);
  roadmap (họ Plan); Product overview (họ Overview). Một dòng list muốn thành spec:
  nói người chạy intake cấp tính năng cho dòng đó.

## Bước 0. Kiểm `dk`

```
dk --version
```

Lỗi hoặc không tìm thấy lệnh: dừng, báo "Cần cài `dk` trước (xem README của dk)".

## `new`: từ Product brief đã duyệt đến list chờ duyệt

### 1. Tìm và kiểm brief

```
dk index intake
```

Đọc `docs/intake/README.md`, mở brief người chỉ (hoặc brief `approved` duy nhất có
`level: project`). Kiểm ba trường: `level: project`, `kind: product`,
`status: approved`. Sai một trường: **dừng**, báo người trường nào sai và ai cần
duyệt; không tạo file. Đã có `docs/overview/feature-list.md`: không tạo lại, sang
phần `update`.

### 2. Tạo file

```
dk new feature-list <slug-san-pham> --from docs/intake/<yymmdd>-<slug>/brief.md --set owner="<người phụ trách>"
```

`title`, `source` được chép từ brief. Không sửa `created`, `created_by`, `dk_version`.

### 3. Đọc nguồn và điền thân

Đọc bốn mục của brief và `idea.md` cùng thư mục (Ai gặp, Vấn đề, Kết quả, Giao diện).

| Mục list | Lấy từ | Cách viết |
|---|---|---|
| 1 Nguồn | brief, idea.md | liên kết tương đối đến hai file |
| 2 Bảng tính năng | brief §1 Kết quả mong muốn, §4 Tiêu chí chấp nhận; idea mục Giao diện | mỗi kết quả kiểm chứng được là một hoặc vài tính năng; mỗi màn hình idea nêu là một ứng viên |
| 3 Nhóm | bảng mục 2 | gom thành 3 đến 7 nhóm theo việc người dùng làm; mỗi nhóm một câu |
| 4 Chưa rõ | brief, idea | câu không tách được thành tính năng; tính năng suy ra mà brief không nhắc; không có thì ghi "không" |

Quy tắc điền một dòng bảng:

- Mã: `FL-01`, `FL-02`, theo thứ tự xuất hiện trong brief.
- Mô tả: một câu "người dùng <làm gì> để <được gì>", có kết quả nhìn thấy.
- Nhóm: tên trùng với mục 3.
- Ưu tiên: `must` khi brief §4 có tiêu chí chấp nhận ứng với nó; `should` khi brief
  §1 nhắc mà §4 không đo; `could` khi chỉ idea nhắc; `won't` khi brief §3 Ngoài phạm
  vi nhắc (giữ dòng để không ai hỏi lại).
- Nguồn: đúng mục sinh ra dòng (`brief §1`, `brief §4`, `idea §8`). Dòng không dẫn
  được mục nào thì không vào bảng, chuyển sang mục 4.
- Spec: để trống.

Không bịa tính năng brief và idea không nhắc. Quá 25 dòng thì gom dòng cùng việc.

| Sai | Đúng |
|---|---|
| Thêm "Đăng nhập" vì sản phẩm nào cũng có | Mục 4: "brief không nhắc đăng nhập, cần người quyết" |
| Mô tả "hệ thống xử lý đơn" | "người bán xác nhận đơn để kho bắt đầu đóng gói" |
| Bỏ FL-03 rồi đánh số lại FL-04 thành FL-03 | FL-03 ghi `won't`, giữ số |

Xong thân: đặt `status: review`, chạy khối "Sau mỗi lần sửa". **Dừng.** Báo người:
đường dẫn, số dòng theo từng ưu tiên, nội dung mục 4; cách chốt: đổi
`status: approved`; nhắc mỗi dòng đi tiếp bằng intake cấp tính năng theo thứ tự
ưu tiên.

## Sau mỗi lần sửa

```
dk changelog add docs/overview/feature-list.md --summary "<tóm tắt thật>" --source <nguồn>
dk render docs/overview/feature-list.md
dk check docs/overview/feature-list.md
```

Tóm tắt thật nói nội dung ("Tách 9 tính năng, 4 nhóm, 2 chưa rõ"), không "cập nhật
file". `--source`: thư mục brief (`<yymmdd>-<slug>`) khi tạo hoặc sửa theo brief; mã
CR khi sửa theo CR; mã spec (`F-xxx`) khi điền cột Spec. `check` có lỗi thì sửa nội
dung và chạy lại; cảnh báo `line-threshold` thì gom dòng.

## `update`

- List `draft` hoặc `review`: sửa theo lời người; dòng mới lấy mã kế tiếp; chạy khối
  "Sau mỗi lần sửa" với `--source` là thư mục brief.
- Điền cột Spec: khi người báo, hoặc `dk index features` cho thấy spec mới có `title`
  khớp tên dòng, ghi `F-xxx` vào cột Spec của dòng đó. Không cần CR, không đổi
  `status`; `--source` là mã spec.
- List `approved` đổi nội dung khác cột Spec: chỉ sửa khi người đưa CR `approved`
  hoặc `in-progress` nêu Feature list (mục mô tả hoặc bảng tác động), hoặc Product
  brief mới `approved`. Thiếu: **dừng**, nói người thay đổi trên list đã chốt là
  Change Request, chỉ sang họ CR. Có: sửa đúng dòng CR nêu, đổi `source` thành mã
  CR hoặc thư mục brief mới, chạy khối với `--source` cùng giá trị.

## `html`

```
dk render docs/overview/feature-list.md
```

Mở `docs/html/overview/feature-list.html` kiểm: bảng mục 2 hiện đủ 7 cột, liên kết
ở mục 1 bấm được. Lệch thì sửa Markdown, không sửa HTML.
