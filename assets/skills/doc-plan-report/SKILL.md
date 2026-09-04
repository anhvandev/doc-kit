---
name: doc-plan-report
description: 'Roadmap, Plan theo đợt việc kèm file phase, Report kết thúc giai đoạn có bằng chứng, Decision log nối dòng, CHANGELOG sản phẩm dạng Keep a Changelog; plan và report ở plans/. Không dùng cho Feature Spec, test hay ADR.'
---

# doc-plan-report

Soạn Roadmap, Plan (một đợt việc một thư mục, mỗi giai đoạn một file), Report kèm
bằng chứng, Decision log và CHANGELOG sản phẩm bằng `dk new` và `dk check`.

**Dừng bắt buộc: không làm việc theo plan khi người chưa duyệt phạm vi và
non-goals (`status: approved`); không sửa report đã `submitted`; Roadmap chỉ
người chốt, skill chỉ đề xuất; không ghi Decision log khi chưa được ủy quyền.**

Mọi file trong `docs/` và `plans/` chỉ được tạo bằng `dk new`. Không tự tạo file
tay, kể cả khi `dk` lỗi. Quy tắc nội dung ở `references/rules.md`; đọc trước khi
làm. Plan và report nằm ở `plans/` ngoài `docs/`, **không** ghi `dk changelog add`
cho chúng; Roadmap, Decision log, CHANGELOG nằm trong `docs/` thì có.

## Phạm vi

- Làm: `docs/plan/roadmap.md`; `plans/<yymmdd>-<hhmm>-<slug>/plan.md` và
  `phase-NN-<slug>.md`; `plans/<đợt>/reports/report-*.md` hoặc `plans/reports/`;
  `docs/plan/decision-log.md`; `docs/CHANGELOG.md`.
- Không làm: Feature Spec, brief, CR (họ Intake, CR, Feature Spec); ADR (quyết
  định lớn, họ ADR); test case và test report (họ Test); release notes cho người
  dùng cuối (họ Release).

## Bước 0. Kiểm `dk` và công cụ plan có sẵn

```
dk --version
```

Lỗi hoặc không tìm thấy lệnh: dừng, báo "Cần cài `dk` trước (xem README của dk)".

Đọc `dk.toml` khóa `plans_dir` (mặc định `plans`). Mở thư mục đó: nếu đã có plan
theo quy ước khác (tên thư mục, frontmatter khác template `dk`), repo đang dùng
công cụ plan riêng. Khi đó skill **không tạo plan hay phase**, chỉ viết Report
(`--in <thư mục plan có sẵn>`) và Decision log; nêu rõ với người.

## `new`: Plan

### 1. Tạo thư mục plan

```
dk new plan <slug> --set owner="<người phụ trách>" --set title="<tên đợt việc>"
```

Nguồn là brief hoặc CR đã chốt: thêm `--set source=<id>`. `dk` tạo
`plans/<yymmdd>-<hhmm>-<slug>/plan.md`.

### 2. Điền plan.md

| Mục | Cách viết |
|---|---|
| 2 Mục tiêu | một đến ba câu, xong đợt này thì có gì |
| 3 Phạm vi | việc cụ thể, liên kết Feature Spec, CR, ADR |
| 4 Non-goals | thứ cố ý không làm |
| 5 Giai đoạn | một dòng một phase, liên kết file phase, trạng thái chép từ phase |
| 6 Tiêu chí hoàn thành | kiểm chứng được: test, tài liệu, số đo |
| 7 Rủi ro | rủi ro, dấu hiệu, phản ứng |
| 8 Report | liên kết report từng giai đoạn, thêm dần |

### 3. Tạo file phase

```
dk new plan-phase <slug> --in plans/<yymmdd>-<hhmm>-<slug> --set owner="<người phụ trách>" --set title="<tên giai đoạn có dấu>"
```

Một giai đoạn một file; số thứ tự `phase-01`, `phase-02` do `dk` đếm trong thư mục.
Điền mục 2 đến 7 của phase: chỉ thứ cần để thực hiện an toàn, không chép plan.md.
Cập nhật bảng mục 5 của plan.md.

### 4. Kiểm và dừng

```
dk check plans/<thư mục plan>/plan.md
```

Đặt `status: draft`. **Dừng.** Báo người đường dẫn, số phase, mục 3 và 4 để duyệt.
Người duyệt thì đổi `status: approved`; chỉ khi đó mới làm việc theo plan.

## `update`: Report kết thúc giai đoạn

### 1. Tạo report ngay khi giai đoạn xong

```
dk new report <phase-slug>-run --in plans/<thư mục plan> --set owner="<người viết>" --set source=phase-NN --set title="<tên giai đoạn có dấu>"
```

Không thuộc plan nào: `--in plans` để vào `plans/reports/`.

### 2. Điền và nộp

- Mục 2 Đã làm: từng việc, liên kết file đã đổi.
- Mục 3 Bằng chứng: mã commit (`commit <hash>`), đường dẫn output test, hoặc khối
  kết quả thật trong ```. Không tóm tắt thay output.
- Mục 4 Chưa xong: việc còn trong phạm vi, lý do, ai làm tiếp.
- Mục 5 Câu hỏi mở: quyết định cần người; không tự quyết trong report.

```
dk check plans/<thư mục plan>/reports/<file report>
```

`check` báo `report-evidence`: bổ sung mục 3 rồi chạy lại. Xong đặt
`status: submitted`; đổi `status` của file phase thành `done` (hoặc `skipped`),
cập nhật cột Trạng thái mục 5 và liên kết mục 8 của plan.md. Sau `submitted`
không sửa report; sai thì viết report mới.

## `update`: Decision log

Chỉ khi người ủy quyền ghi (nói rõ trong phiên hoặc trong plan). Một quyết định
một dòng, kể cả quyết định giao diện:

```
dk new decision-log --append "<quyết định> | <người chốt> | <liên kết CR, spec, mockup>"
dk changelog add docs/plan/decision-log.md --summary "<quyết định>" --source <CR hoặc plan>
```

Chưa có `docs/plan/decision-log.md`: lần đầu thêm `--set owner="<người giữ log>"`,
`dk` tạo file rồi nối dòng; các lần sau chỉ nối dòng có ngày (`--set` bị bỏ qua).
Không sửa dòng cũ. Quyết định đủ lớn (đổi kiến trúc, thư viện, hợp đồng dữ liệu)
thì dừng và chỉ người sang họ ADR.

## `update`: CHANGELOG sản phẩm

Khi phát hành một phiên bản:

```
dk new changelog-product changelog --set owner="<người duyệt>"
```

(chỉ lần đầu; đã có thì sửa `docs/CHANGELOG.md`). Đọc mọi report của đợt và
`git log` từ phiên bản trước; gom vào mục `## [x.y.z] - yyyy-mm-dd` với nhóm
Added, Changed, Deprecated, Removed, Fixed, Security; mỗi dòng nói người dùng
làm được gì, không nói file hay hàm. Giữ `## [Unreleased]` ở trên cùng.

```
dk changelog add docs/CHANGELOG.md --summary "Phiên bản x.y.z" --source <plan hoặc CR>
dk render docs/CHANGELOG.md
dk check docs/CHANGELOG.md
```

**Dừng.** Người duyệt chữ trước khi phát hành.

## `new`: Roadmap

Roadmap là tài liệu thẩm quyền, người chốt. Skill chỉ đề xuất:

```
dk index features
dk new roadmap roadmap --set owner="<người chốt>"
```

Điền mục 3 từ `docs/features/README.md` (Feature catalog): mã, tên, trạng thái;
xếp mốc theo thứ tự người nói hoặc theo phụ thuộc giữa spec. Chưa có người chốt
mốc thì để mục 4 Chưa xếp. Chạy khối "Sau mỗi lần sửa file trong docs/" rồi
**dừng**, hỏi người chốt thứ tự và mốc (`status: approved`).

## Sau mỗi lần sửa file trong `docs/`

```
dk changelog add docs/<file> --summary "<tóm tắt thật>" --source <plan, CR hoặc ->
dk render docs/<file>
dk check docs/<file>
```

Không chạy `changelog add` cho file trong `plans/`. `check` còn lỗi thì sửa rồi
chạy lại, không báo người khi còn lỗi.

## `html`

```
dk render docs/plan/roadmap.md
dk render docs/plan/decision-log.md
dk render docs/CHANGELOG.md
```

Mở `docs/html/plan/*.html`, `docs/html/CHANGELOG.html` kiểm liên kết mã tính năng
trong Roadmap nhảy đúng spec. Plan và report không render (ngoài `docs/`).
