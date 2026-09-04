# Báo cáo tổng hợp: Hệ thống tài liệu cho dự án phần mềm kết hợp người và AI

- Ngày: 2026-09-03
- Cập nhật: 2026-09-03 11:14, bổ sung chuỗi tiếp nhận yêu cầu (Intake brief) và cấu trúc Feature Spec có sơ đồ luồng
- Cập nhật: 2026-09-03 11:24, bổ sung quy tắc tách file; một đơn vị một file áp dụng bất kể số lượng tính năng
- Cập nhật: 2026-09-03 11:40, sắp xếp lại các tầng theo vòng đời (tiếp nhận ý tưởng lên Tầng 1, chạy hai cấp), thêm Tầng 3 Design
- Cập nhật: 2026-09-03 13:28, bổ sung cấu trúc Design system theo Atomic Design, bộ tối thiểu và quy tắc nhất quán
- Cập nhật: 2026-09-03 13:55, Tầng 1 viết lại thành mô hình Change Request: một template gốc, mỗi CR một thư mục, bốn loại CR, quy tắc tỷ lệ
- Cập nhật: 2026-09-03 14:10, tách Change Request thành loại tài liệu riêng có phân tích tác động; Tầng 1 gồm hai nhánh ý tưởng mới và CR
- Loại: tư vấn kiến trúc tài liệu
- Phạm vi: danh sách loại tài liệu, mục đích, cách dùng, vai trò người và AI, quy tắc tách file, lộ trình áp dụng

## 1. Nguyên tắc nền

- Mỗi tài liệu tồn tại để trả lời một câu hỏi cụ thể. Không có câu hỏi thì không viết.
- Một nguồn sự thật cho mỗi tính năng (Feature Spec). Các tài liệu khác sinh ra hoặc tham chiếu từ đó.
- Phân biệt hai loại:
  - **Thẩm quyền**: mô tả trạng thái hiện tại, luôn phải đúng, được sửa tại chỗ.
  - **Trạng thái**: ghi lại quá trình, có ngày tháng, không sửa lại sau khi chốt.
- Tài liệu phải máy đọc được: cấu trúc cố định, Markdown trong repo, frontmatter khi cần sinh tự động.
- Mỗi tài liệu có một chủ sở hữu. Không có người chịu trách nhiệm thì không viết.
- Hai loại brief, không nhầm lẫn: **Intake brief** là đầu vào (tổng kết của một Change Request, trước Feature Spec), **Release brief** là đầu ra (giới thiệu tính năng cho người dùng, sau khi hoàn thành).
- Hai nhánh vào hệ thống: **ý tưởng mới** (chưa có gì tồn tại) đi qua chuỗi ba bước ý tưởng, phỏng vấn, tổng kết; **Change Request (CR)** là tài liệu riêng cho thay đổi trên thứ đã có, có phân tích tác động và tự rẽ sang Feature Spec, Design, ADR hay Test tùy tác động. Template viết một lần; mỗi bản điền nằm trong thư mục riêng.
- AI không tự sửa tài liệu thẩm quyền, mục hành vi của Feature Spec, Design system hoặc mockup đã duyệt nếu chưa có người duyệt.
- **Một đơn vị một file, bất kể số lượng tính năng.** Feature Spec, ADR, Plan, Report, Runbook luôn tách từ file đầu tiên; không gộp lúc dự án nhỏ rồi tách sau, vì chi phí tách sau luôn cao hơn.
- Các tầng xếp theo **vòng đời**: ý tưởng sinh ra trước, vận hành đến sau. Người mới cần đọc nhanh thì đi từ Tầng 0, nơi README và Agent context file trỏ thẳng đến Product overview.

## 2. Danh sách tài liệu

### Tầng 0. Điểm vào

| Tài liệu | Dùng để làm gì | Cách dùng | Loại | Người | AI |
|---|---|---|---|---|---|
| Agent context file (CLAUDE.md / AGENTS.md) | Điểm vào duy nhất của AI: quy ước code, lệnh chạy, cấu trúc thư mục, nơi tìm từng loại tài liệu, ngưỡng dòng và quy tắc tách file | Đặt ở gốc repo. Ngắn, chỉ chứa quy tắc và con trỏ, không chép nội dung tài liệu khác | Thẩm quyền | Viết, chốt | Đề xuất cập nhật khi quy ước đổi |
| README | Điểm vào của người: dự án là gì, cài đặt, chạy, test trong 5 phút | Giữ dưới 1 màn hình cuộn. Liên kết đến docs chi tiết | Thẩm quyền | Viết, chốt | Soạn nháp, kiểm tra lệnh còn chạy được |

### Tầng 1. Tiếp nhận yêu cầu

Hai nhánh, hai loại tài liệu khác nhau:

| Nhánh | Khi nào | Tài liệu | Đầu ra |
|---|---|---|---|
| Ý tưởng mới | Chưa có gì tồn tại: sản phẩm mới, tính năng mới, hệ giao diện lần đầu | Chuỗi ba bước idea, interview, brief | Product overview, Feature Spec mới, Design brief |
| Change Request | Thay đổi trên thứ đã có: sửa hành vi, đổi giao diện, đổi kỹ thuật, sửa lỗi | Một file CR có phân tích tác động | Cập nhật Feature Spec, mockup, ADR, Test case bị ảnh hưởng |

Phân biệt hai khái niệm dùng chung cho cả hai nhánh:

- **Template**: mẫu cố định, viết một lần, có các trường trống kèm câu hỏi gợi ý cho từng trường. Trường mới phải thay trường cũ, không phình quá một trang.
- **Bản điền**: mỗi ý tưởng hoặc mỗi CR tạo một bản mới từ template, nằm trong thư mục riêng, frontmatter có trạng thái.

#### 1a. Ý tưởng mới (chuỗi ba bước)

```
intake/
  260903-bo-loc-don-hang/
    idea.md        ← điền từ template ý tưởng
    interview.md   ← hỏi đáp, mỗi câu kết bằng quyết định
    brief.md       ← tổng kết (Intake brief), người duyệt
```

| Tài liệu | Dùng để làm gì | Cách dùng | Loại | Người | AI |
|---|---|---|---|---|---|
| Template ý tưởng | Mẫu chung: vấn đề, ai gặp, cách xoay xở hiện tại, kết quả mong muốn, điều không làm, ràng buộc, thước đo thành công, mục Giao diện (màn hình liên quan, template trang dùng, có màn hình hoặc organism mới không), câu hỏi mở | Một file gốc. Mỗi trường có câu hỏi gợi ý | Thẩm quyền | Chốt | Đề xuất sửa khi nhiều bản điền bỏ trống cùng một trường |
| idea.md | Ghi nhanh ý tưởng theo template | Điền trong 10 phút, trường chưa biết ghi "chưa rõ" | Trạng thái | Điền | Không |
| interview.md | Làm rõ từng trường "chưa rõ" hoặc mâu thuẫn trong idea.md | Câu hỏi sinh từ chỗ trống của template, không từ danh sách cố định. Hỏi từng câu một, mỗi cặp hỏi đáp kết thúc bằng dòng **Quyết định rút ra** và trạng thái đã chốt hoặc còn mở. Câu hỏi về giao diện phải hỏi bằng hình: AI đưa 2 đến 3 phương án tham chiếu để người chọn. Có ngày, không sửa sau khi kết thúc | Trạng thái | Trả lời, chọn | Dẫn phỏng vấn, ghi quyết định |
| brief.md (Intake brief) | Nguồn sự thật của ý tưởng: kết quả mong muốn, ràng buộc, ngoài phạm vi, tiêu chí chấp nhận, quyết định đã chốt (liên kết về interview.md) | Bốn mục cố định. Thay thế idea.md làm nguồn sự thật. Đầu vào cho Feature Spec và User flow | Thẩm quyền | Duyệt và chốt | Soạn nháp, không tự chốt |

Chuỗi này chạy ở hai cấp với cùng template:

- **Cấp dự án**, một lần lúc khởi đầu. Ý tưởng sản phẩm: brief viết thành **Product overview** (Tầng 2). Ý tưởng giao diện toàn sản phẩm: brief là **Design brief**, đầu vào cho Design system lớp 0 và templates (Tầng 3). Design brief chạy sau hoặc song song ý tưởng sản phẩm, không chạy trước.
- **Cấp tính năng**, lặp lại mỗi tính năng mới. Brief đổ vào User flow, Wireframe, Mockup rồi Feature Spec.

Quy tắc tỷ lệ:

- idea.md không còn trường "chưa rõ" và không mâu thuẫn: bỏ phỏng vấn, AI soạn thẳng brief.md, người duyệt.
- Có trường "chưa rõ" hoặc mâu thuẫn: phỏng vấn, chỉ hỏi đúng những trường đó.
- Tính năng dùng template trang có sẵn: mục Giao diện trong idea.md là đủ, không cần Design brief riêng.

#### 1b. Change Request

CR là tài liệu riêng vì khác ý tưởng mới ở ba điểm: có **trạng thái hiện tại** để so sánh, có **phân tích tác động** lên tài liệu và code đang tồn tại, và **đích đến không cố định** mà tùy tác động. Một CR có thể chỉ chạm Feature Spec, chỉ chạm mockup, hoặc cả hai cộng thêm ADR.

```
cr/
  260910-doi-cach-loc-theo-trang-thai.md     ← một CR một file
  260910-doi-cach-loc-theo-trang-thai/       ← chỉ khi phỏng vấn dài
    interview.md
```

| Tài liệu | Dùng để làm gì | Cách dùng | Loại | Người | AI |
|---|---|---|---|---|---|
| Template CR | Mẫu cố định cho mọi CR | Một file gốc, dưới một trang, mỗi trường có câu hỏi gợi ý | Thẩm quyền | Chốt | Đề xuất sửa |
| CR (bản điền) | Ghi yêu cầu thay đổi, phân tích tác động, hỏi đáp, quyết định, trạng thái | Một file. Tách interview.md ra thư mục riêng khi hỏi đáp dài. Không sửa phần quyết định sau khi chốt | Trạng thái, phần quyết định là thẩm quyền cho đến khi tài liệu đích cập nhật xong | Điền yêu cầu, trả lời, duyệt tác động và quyết định | Soạn phân tích tác động từ tài liệu hiện có, dẫn hỏi đáp, cập nhật tài liệu đích sau khi duyệt |

Cấu trúc một CR:

1. Frontmatter: mã CR, tiêu đề, người yêu cầu, ngày, trạng thái, mức ưu tiên.
2. Hiện trạng: liên kết đến Feature Spec, mockup, ADR đang mô tả thứ sắp đổi. Không chép lại nội dung.
3. Thay đổi đề xuất và lý do.
4. **Phân tích tác động**, bảng do AI soạn nháp từ tài liệu hiện có, người duyệt:

| Tài liệu hoặc thành phần | Bị ảnh hưởng | Cần làm gì |
|---|---|---|
| Feature Spec | Có / Không | Sửa mã bước nào, thêm ngoại lệ nào |
| Design: mockup, Design system | Có / Không | Sửa màn hình nào, cần organism hay template mới không |
| ADR | Có / Không | Quyết định nào bị thay thế |
| Data model, API | Có / Không | Migration, contract đổi |
| Test case | Có / Không | Test nào sửa, thêm |
| Release brief, User guide | Có / Không | Đoạn nào cập nhật |

5. Hỏi đáp: chỉ khi thay đổi hoặc tác động còn "chưa rõ". Cùng quy tắc với interview.md: từng câu, kết bằng quyết định.
6. Quyết định: chấp nhận, từ chối, hoãn. Ghi lý do và người chốt.
7. Danh sách tài liệu đích đã cập nhật, liên kết đến từng cái, đánh dấu khi xong.

Trạng thái CR: mới, đang phân tích, chờ duyệt, đã chốt, đang triển khai, đã đóng, từ chối. Chỉ mục CR sinh từ frontmatter trả lời "đang có thay đổi nào".

Quy tắc tỷ lệ:

- Sửa lỗi nhỏ, tác động chỉ một Feature Spec và test: mục 4 rút gọn còn hai dòng, bỏ mục 5.
- Tác động chạm Design system hoặc ADR: bắt buộc người duyệt mục 4 trước khi AI sửa bất kỳ tài liệu đích nào.
- CR bị từ chối vẫn giữ file, làm lịch sử vì sao không đổi.

Liên kết ngược bắt buộc cho cả hai nhánh: Feature Spec, mockup, ADR có trường liên kết về Intake brief hoặc CR sinh ra hoặc sửa nó. Brief hoặc CR chốt rồi không đổ đi đâu là tài liệu chết.

### Tầng 2. Sản phẩm và hệ thống

| Tài liệu | Dùng để làm gì | Cách dùng | Loại | Người | AI |
|---|---|---|---|---|---|
| Product overview | Sản phẩm giải quyết vấn đề gì, cho ai, nhóm tính năng chính, thước đo thành công | Một trang. Là Intake brief cấp dự án (Tầng 1a). Đọc đầu tiên khi tham gia dự án | Thẩm quyền | Viết, chốt | Đề xuất nháp |
| Feature catalog | Bảng điều khiển trả lời "phần mềm có gì": mã, tên, trạng thái, liên kết Feature Spec | Sinh tự động từ frontmatter của Feature Spec. Cập nhật trạng thái từ Report | Thẩm quyền (sinh) | Duyệt | Sinh và cập nhật |
| Architecture overview | Sơ đồ thành phần, luồng dữ liệu, ranh giới hệ thống, tech stack | Sơ đồ trước, văn xuôi sau. Cập nhật khi thêm hoặc bỏ thành phần | Thẩm quyền | Viết, chốt | Đề xuất nháp, phát hiện lệch với code |
| ADR (Architecture Decision Record) | Ghi bối cảnh, lựa chọn, hệ quả của mỗi quyết định kỹ thuật quan trọng | Một file cho một quyết định, đánh số, không sửa sau khi chốt; quyết định mới thay thế thì ghi ADR mới | Thẩm quyền (bất biến) | Chốt | Soạn nháp từ thảo luận |
| Domain glossary | Thuật ngữ nghiệp vụ và định nghĩa thống nhất | Bắt buộc khi có AI. AI dùng sai thuật ngữ sinh code sai | Thẩm quyền | Viết, chốt | Phát hiện thuật ngữ chưa định nghĩa |
| Data model / Schema | Mô hình dữ liệu, quy ước đặt tên, chiến lược migration | Sinh từ migration hoặc ORM. Không viết tay | Thẩm quyền (sinh) | Duyệt | Sinh |
| API specification | Contract giữa client và server hoặc giữa các service | OpenAPI, GraphQL schema. Sinh từ code | Thẩm quyền (sinh) | Duyệt | Sinh |
| Security & threat model | Tài sản cần bảo vệ, bề mặt tấn công, biện pháp kiểm soát | Viết khi hệ thống có dữ liệu người dùng hoặc thanh toán | Thẩm quyền | Viết, chốt | Đề xuất bề mặt tấn công |

### Tầng 3. Design (giao diện và trải nghiệm)

Chỉ áp dụng khi phần mềm có giao diện người dùng. Design system là cấp dự án, viết một lần; các tài liệu còn lại là cấp tính năng hoặc cấp màn hình.

| Tài liệu | Dùng để làm gì | Cách dùng | Loại | Người | AI |
|---|---|---|---|---|---|
| Design system | Mọi màn hình dùng chung quy ước gì: tokens, quy tắc nền, component theo Atomic Design (atoms, molecules, organisms, templates), patterns, giọng văn giao diện. Cấu trúc chi tiết ở mục con bên dưới | Tokens dạng máy đọc được (JSON sinh CSS variables) trong repo. Component sống trong code, tài liệu sinh từ Storybook hoặc tương đương. Tài liệu AI cần nhất, không có thì AI sinh mỗi màn hình một kiểu | Thẩm quyền | Chốt tokens, pattern, giọng văn | Đề xuất, kiểm tra code UI có lệch tokens không |
| User flow / Wireflow | Chuỗi màn hình và hành động chuyển giữa chúng cho một tính năng | Mermaid trong Markdown, dùng chung mã bước với sơ đồ luồng của Feature Spec | Thẩm quyền | Duyệt | Soạn từ Intake brief |
| Wireframe | Màn hình có những gì, bố cục ra sao | Đen trắng, không màu. Chốt "có gì trên màn hình" trước khi vẽ đẹp. Một màn hình một frame, đặt tên theo mã tính năng và mã bước | Trạng thái | Vẽ hoặc duyệt | Sinh nháp dạng text hoặc HTML thô |
| Mockup | Màn hình trông như thế nào | Theo Design system. Mỗi trạng thái quan trọng một mockup: bình thường, rỗng, lỗi, đang tải. Lưu liên kết công cụ thiết kế kèm ảnh xuất tĩnh trong repo để AI đọc được | Thẩm quyền | Vẽ, duyệt | Sinh HTML theo Design system cho màn hình đơn giản; không thay designer cho màn hình quan trọng |
| Prototype | Bấm vào thì chuyện gì xảy ra | Bản bấm được từ mockup, dùng để test với người dùng hoặc duyệt với bên liên quan trước khi code. Công cụ ra quyết định, không lưu lâu dài | Trạng thái | Duyệt, chạy test với người dùng | Dựng từ mockup |
| UI spec / Redline | Kích thước, khoảng cách, hành vi responsive, animation | Sinh từ công cụ thiết kế, không viết tay | Thẩm quyền (sinh) | Duyệt | Sinh |
| Design review record | Quyết định thay đổi giao diện và lý do | Gộp vào Decision log (Tầng 5), không cần file riêng | Trạng thái | Ghi | Ghi khi được ủy quyền |

Thứ tự trong tầng: User flow, Wireframe (người duyệt), Mockup, Prototype. Đầu vào là mục Giao diện trong Intake brief với tính năng mới, hoặc bảng tác động của CR với thay đổi giao diện trên tính năng đã có. Feature Spec liên kết đến mockup theo mã bước, không chép ảnh vào spec. Test giao diện so với mockup đã duyệt.

#### Cấu trúc Design system (Atomic Design)

Nhất quán đến từ việc mọi giá trị hiển thị đều được đặt tên và tham chiếu, không có giá trị gõ tay. Xây từ dưới lên, mỗi lớp chỉ tham chiếu lớp ngay dưới nó.

| Lớp | Là gì | Ví dụ | Quy tắc |
|---|---|---|---|
| 0. Primitive tokens | Giá trị thô theo thang ít bước: màu, chữ, khoảng cách, kích thước, bo góc, bóng, viền, chuyển động, điểm gãy, z-index | `blue.500`, `space.4`, `radius.md` | Chỉ semantic tokens được tham chiếu. Mỗi thang dưới 10 bước |
| 0. Semantic tokens | Giá trị thô gắn ý nghĩa dùng: màu theo vai trò và trạng thái, chữ theo vai trò, khoảng cách theo vai trò, độ nổi theo vai trò. Theme sáng tối và đa thương hiệu đổi ở đây | `color.action.primary`, `text.body`, `space.inset.md`, `elevation.card` | Lớp duy nhất component được tham chiếu |
| 0. Quy tắc nền | Lưới và bố cục theo điểm gãy, accessibility, một bộ icon duy nhất, giọng văn giao diện | Tương phản 4.5:1, vùng chạm 44px, focus ring bắt buộc | Áp cho mọi lớp trên |
| 1. Atoms | Phần tử không chia nhỏ hơn được | Text, Icon, Button, Input, Checkbox, Radio, Switch, Badge, Avatar, Spinner, Divider | Không chứa atom khác. Đủ biến thể, kích thước, trạng thái, quy tắc dùng, ARIA |
| 2. Molecules | Vài atom ghép thành một đơn vị có chức năng | Form field (label + input + lỗi + gợi ý), Select, Dropdown menu, Tooltip, Toast, Tag, Tabs, Pagination | Chỉ chứa atom |
| 3. Organisms | Khối giao diện hoàn chỉnh, có thể mang logic | Card, Modal, Drawer, Table, Navigation, Empty state, Form | Chứa molecule và atom |
| 4. Templates | Khung trang, chưa có dữ liệu thật | Trang danh sách, chi tiết, form, dashboard, xác thực | Chứa organism, định nghĩa vùng đặt |
| 5. Pages | Template đổ dữ liệu thật của một tính năng | Trang danh sách đơn hàng | Không định nghĩa style mới |

Patterns không phải lớp riêng, là công thức ghép organism cho một tình huống: form dài, xác nhận xóa, tìm kiếm và lọc bảng, trạng thái tải, trạng thái rỗng, thông báo lỗi hệ thống, luồng nhiều bước. Viết thành tài liệu đi kèm lớp 3, giữ nhất quán về hành vi chứ không chỉ hình ảnh.

Mỗi component định nghĩa đủ: biến thể (primary, secondary, ghost, destructive), kích thước (sm, md, lg theo thang chung), trạng thái (mặc định, hover, active, focus, disabled, loading, lỗi), quy tắc dùng và không dùng, accessibility (vai trò ARIA, điều hướng bàn phím).

Bộ tối thiểu để bắt đầu, theo thứ tự:

1. Lớp 0: thang khoảng cách và thang chữ, bảng màu và semantic cho nền, chữ, hành động, viền, bốn trạng thái, bo góc, bóng, icon, accessibility cơ bản.
2. Lớp 1: Text, Icon, Button, Input.
3. Lớp 2: Form field.
4. Lớp 3: Form, Empty state, Toast.
5. Lớp 4: một template trang danh sách và một template trang form.

Ba quy tắc giữ nhất quán:

- Mỗi lớp chỉ tham chiếu lớp ngay dưới. Organism không gọi thẳng token, page không định nghĩa style, component không tham chiếu primitive.
- Tên token trong Figma (Figma Variables) và trong code giống hệt nhau. Đây là điểm hay đứt nhất ở ranh giới bàn giao.
- Lint cấm giá trị gõ tay trong code UI và cấm import vượt lớp. Tài liệu không giữ được nhất quán, lint mới giữ được.

Dùng thư viện UI có sẵn (MUI, Ant, shadcn): lớp 0 primitive, lớp 1 đến 3 đã có. Chỉ viết lớp semantic ánh xạ sang thư viện, quy tắc dùng, patterns, templates và giọng văn. Không xây lại.

Kiểm tra thật: cho AI sinh một page từ template có sẵn, chỉ dùng organism và semantic tokens, lint bật. Qua được là hệ đã nhất quán.

### Tầng 4. Feature Spec (trung tâm)

| Tài liệu | Dùng để làm gì | Cách dùng | Loại | Người | AI |
|---|---|---|---|---|---|
| Feature Spec | Nguồn sự thật cho một tính năng: sơ đồ luồng, hành vi, action, kết quả, giao diện, ngoại lệ, tiêu chí chấp nhận | Một file cho một tính năng, sơ đồ Mermaid ở trên, text spec ở dưới, cấu trúc cố định (xem mục 3). Sửa hành vi thì sửa đây trước | Thẩm quyền | Viết mục đích, hành vi, quy tắc; duyệt toàn bộ | Soạn ngoại lệ, tiêu chí chấp nhận, kiểm tra thiếu sót |

Các định dạng thay thế tùy bối cảnh:

| Bối cảnh | Định dạng |
|---|---|
| Nghiệp vụ phức tạp, nhiều tác nhân, nhiều ngoại lệ | Use Case Specification + Activity / Sequence diagram |
| Team Agile, muốn tài liệu chạy được thành test | User Story + Acceptance Criteria dạng Gherkin |
| Team nhỏ, tài liệu sống trong repo | Feature Spec Markdown (khuyến nghị mặc định) |
| Outsource, cần bàn giao và nghiệm thu | Functional Spec đầy đủ + Test Case Specification |
| Tính năng có trạng thái phức tạp | Bổ sung State diagram |
| CRUD đơn giản | Bảng field và quyền, không viết use case |

### Tầng 5. Kế hoạch và tiến độ

| Tài liệu | Dùng để làm gì | Cách dùng | Loại | Người | AI |
|---|---|---|---|---|---|
| Roadmap | Thứ tự ưu tiên tính năng theo mốc thời gian | Cập nhật theo quý hoặc theo đợt phát hành | Thẩm quyền | Chốt | Không |
| Plan | Trả lời "cần làm gì": mục tiêu, phạm vi, non-goals, giai đoạn, tiêu chí hoàn thành | Một thư mục cho một đợt việc. AI làm việc theo plan này. Chốt phạm vi trước khi bắt đầu | Trạng thái | Duyệt phạm vi và non-goals | Soạn giai đoạn và bước |
| Report | Trả lời "đã làm gì": việc đã xong, cách kiểm chứng, còn gì chưa xong, câu hỏi mở | Viết ngay khi kết thúc một giai đoạn, kèm bằng chứng (test output, commit). Không sửa sau khi nộp | Trạng thái | Đọc, xác nhận | Viết |
| Decision log | Các quyết định nhỏ chưa đến mức ADR, gồm cả quyết định về giao diện | Ghi theo ngày, một dòng cho một quyết định | Trạng thái | Ghi | Ghi khi được ủy quyền |
| CHANGELOG | Thay đổi theo phiên bản, ngôn ngữ hướng người dùng | Tổng hợp từ Report và commit. Theo Keep a Changelog | Trạng thái | Duyệt | Tổng hợp |

### Tầng 6. Testing

| Tài liệu | Dùng để làm gì | Cách dùng | Loại | Người | AI |
|---|---|---|---|---|---|
| Testing strategy | Loại test, công cụ, phạm vi bao phủ mong đợi, cách chạy | Một file. Cập nhật khi đổi công cụ hoặc ngưỡng | Thẩm quyền | Chốt | Đề xuất |
| Test case | Bước thực hiện, dữ liệu vào, kết quả mong đợi cho từng trường hợp | Sinh từ mục Tiêu chí chấp nhận của Feature Spec. Nếu dùng Gherkin thì file test chính là tài liệu | Sinh | Duyệt trường hợp biên | Sinh, chạy, báo kết quả |
| Test giao diện | Màn hình thật có khớp mockup đã duyệt và Design system không | Visual regression hoặc checklist theo mã bước. Chuẩn so sánh là mockup, không phải cảm tính | Sinh | Duyệt sai lệch | Chạy, báo sai lệch |
| Test report | Kết quả một đợt kiểm thử: cái gì chạy, fail, quyết định xử lý | Theo đợt phát hành. Đính kèm output thật | Trạng thái | Quyết định xử lý fail | Viết |

### Tầng 7. Người dùng cuối

| Tài liệu | Dùng để làm gì | Cách dùng | Loại | Người | AI |
|---|---|---|---|---|---|
| Release brief | Giới thiệu một tính năng cho người dùng: giúp gì, ai dùng, bước cơ bản, giới hạn đã biết | Sinh từ mục Mục đích, Hành vi, Giao diện của Feature Spec, viết lại bằng ngôn ngữ người dùng, không thuật ngữ kỹ thuật, có ảnh màn hình | Thẩm quyền | Duyệt giọng văn, thêm ảnh | Soạn nháp |
| Release notes | Tập hợp Release brief và sửa lỗi trong một đợt phát hành | Gửi cho người dùng hoặc đăng công khai | Trạng thái | Duyệt | Tổng hợp |
| User guide | Hướng dẫn sử dụng tích lũy | Ghép dần từ các Release brief. Tổ chức theo nhiệm vụ người dùng, không theo menu | Thẩm quyền | Duyệt | Soạn, cập nhật |
| FAQ / Troubleshooting | Câu hỏi và lỗi thường gặp | Viết khi câu hỏi lặp lại lần thứ hai | Thẩm quyền | Duyệt | Soạn từ ticket hỗ trợ |

### Tầng 8. Vận hành (khi đã lên production)

| Tài liệu | Dùng để làm gì | Cách dùng | Loại |
|---|---|---|---|
| Deployment guide | Quy trình triển khai, rollback, CI/CD | Chạy được từng bước, không giả định kiến thức ngầm | Thẩm quyền |
| Environment & configuration | Biến môi trường, nơi giữ secrets, khác biệt dev/staging/prod | Không chứa giá trị secret, chỉ tên và ý nghĩa | Thẩm quyền |
| Runbook | Xử lý từng sự cố đã biết theo từng bước | Một sự cố một mục. Viết sau mỗi incident | Thẩm quyền |
| Monitoring & alerting | Metric, ngưỡng, người nhận cảnh báo | Liên kết đến dashboard thật | Thẩm quyền |
| Incident postmortem | Nguyên nhân gốc, timeline, hành động khắc phục | Không đổ lỗi cá nhân. Viết trong 48 giờ sau sự cố | Trạng thái |
| Backup & disaster recovery | RPO, RTO, quy trình khôi phục | Phải diễn tập ít nhất một lần | Thẩm quyền |

### Tầng 9. Quản trị và pháp lý (tùy quy mô, ngành)

| Tài liệu | Dùng để làm gì | Khi nào cần |
|---|---|---|
| Project charter / SOW | Phạm vi, ngân sách, cột mốc, bên liên quan | Có khách hàng hoặc nhiều bên |
| Risk register | Rủi ro, xác suất, tác động, giảm thiểu | Dự án trên 3 tháng hoặc nhiều bên |
| Meeting notes | Quyết định nào, ai chốt, khi nào | Có họp định kỳ |
| Privacy policy, Terms of Service | Nghĩa vụ pháp lý với người dùng | Có người dùng bên ngoài |
| Tài liệu tuân thủ (ISO 27001, SOC 2, Nghị định 13/2023) | Chứng minh tuân thủ | Ngành yêu cầu |
| Giấy phép mã nguồn mở của dependency | Tránh vi phạm license | Phát hành thương mại |

## 3. Mẫu Feature Spec

Cấu trúc cố định, cả người và AI đọc theo cùng một mẫu:

1. Frontmatter: mã, tên tính năng, trạng thái, chủ sở hữu, liên kết Intake brief (để sinh Feature catalog).
2. Mục đích và giá trị cho người dùng.
3. Tác nhân và điều kiện tiên quyết.
4. **Sơ đồ luồng chính**: Mermaid nhúng trong Markdown, mỗi nút mang mã bước (B1, B2...).
5. **Bảng hành vi theo mã bước**: mỗi dòng là một cặp **hành động của tác nhân** và **phản hồi quan sát được của hệ thống**, dùng đúng mã bước trong sơ đồ.
6. **Giao diện**: mỗi mã bước liên kết đến mockup của màn hình tương ứng và trạng thái hiển thị. Chỉ liên kết, không chép ảnh. Bỏ mục này nếu tính năng không có giao diện.
7. Luồng thay thế và ngoại lệ: text hoặc sơ đồ nhỏ riêng, không nhồi vào sơ đồ chính.
8. Quy tắc nghiệp vụ, tham chiếu bằng mã, không lặp lại giữa các spec.
9. Tiêu chí chấp nhận dạng Given / When / Then.
10. Dữ liệu và API liên quan: chỉ liên kết, không chép lại.
11. Ngoài phạm vi.

Quy tắc sơ đồ:
- Dạng text (Mermaid), không dạng ảnh: AI đọc và sửa được, diff được, nằm cùng file nên không lệch với text.
- Sơ đồ là mục lục, text là diễn giải từng nút. Ràng buộc bằng mã bước để hai phần không tách nhau. User flow ở Tầng 3 dùng chung mã bước này.
- Một sơ đồ cho luồng chính. Luồng tuyến tính dưới bốn bước hoặc CRUD đơn giản thì bỏ sơ đồ, bảng hành vi là đủ.

| Bản chất tính năng | Loại sơ đồ |
|---|---|
| Chuỗi hành động có rẽ nhánh | Flowchart |
| Nhiều thành phần tương tác theo thứ tự | Sequence diagram |
| Xoay quanh vòng đời trạng thái của một đối tượng | State diagram |

Nguyên tắc viết:
- Không viết "hệ thống xử lý" mà không nói kết quả nhìn thấy được.
- Ngoại lệ và trường hợp biên chiếm ít nhất một nửa giá trị tài liệu.
- Kết quả mong đợi phải kiểm chứng được: thông báo cụ thể, dữ liệu thay đổi thế nào, sự kiện nào phát ra.
- Dùng bảng và sơ đồ thay cho văn xuôi dài.

Mục 5, 7, 9 nuôi Test case. Mục 6 nuôi Test giao diện. Mục 2, 4, 5, 6 nuôi Release brief.

## 4. Luồng dữ liệu giữa các tài liệu

```
[Ý tưởng mới]  idea.md ──► interview.md (khi có trường trống) ──► brief.md (người duyệt)
                                                                      │
              ┌───────────────────────────┬───────────────────────────┤
              ▼                           ▼                           ▼
     Cấp dự án: sản phẩm        Cấp dự án: giao diện         Cấp tính năng
              │                           │                           │
              ▼                           ▼                           │
      Product overview        Design system + Templates               │
              │                           │                           │
              ▼                           │                           ▼
      Feature catalog                     └──► User flow ──► Wireframe (duyệt) ──► Mockup ──► Prototype
              │                                                       │
              └──────────────► Feature Spec ◄─────────────────────────┘ (mục Giao diện liên kết mockup)
                                    │
                    ┌───────────────┼───────────────┬───────────────┐
                    ▼               ▼               ▼               ▼
                Plan / Report    Test case     Test giao diện   Release brief
                    │               │               │               │
                    └───────────────┴───────────────┴───────────────┘
                                            ▼
                                  CHANGELOG / Release notes

[Change Request]  cr.md: hiện trạng ──► thay đổi ──► bảng tác động (người duyệt) ──► quyết định
                                                            │
                     ┌──────────────┬──────────────┬────────┴──────┬──────────────┐
                     ▼              ▼              ▼               ▼              ▼
               Feature Spec    Mockup /        ADR           Data / API      Test case
               (sửa mã bước)   Design system                                 Release brief
                     └──────────────┴──────────────┴───────────────┴──────────────┘
                                            ▼
                                   Plan / Report ──► CHANGELOG
```

Mũi tên chỉ đi một chiều. Sửa hành vi thì sửa Feature Spec trước, rồi cập nhật tài liệu con. Sửa giao diện thì sửa mockup trước, rồi cập nhật spec và code. Không sửa test hay brief mà bỏ qua spec. brief.md chưa được người duyệt thì không đi tiếp sang Design hay Feature Spec. Bảng tác động của CR chưa được người duyệt thì AI không sửa bất kỳ tài liệu đích nào. Tài liệu đích phải liên kết ngược về Intake brief hoặc CR sinh ra hoặc sửa nó.

## 5. Ánh xạ nhu cầu sang tài liệu

| Nhu cầu | Tài liệu chủ | Tài liệu hỗ trợ |
|---|---|---|
| Tiếp nhận ý tưởng mới | Intake brief (brief.md) | Template ý tưởng, idea.md, interview.md |
| Thay đổi thứ đã có | Change Request (bảng tác động) | Chỉ mục CR theo trạng thái, tài liệu đích bị ảnh hưởng |
| Hiểu tính năng của phần mềm | Feature Spec (sơ đồ + bảng hành vi) | Product overview, Feature catalog, Glossary |
| Hình dung giao diện trước khi code | Mockup | Design system, User flow, Wireframe, Prototype |
| Cần làm gì | Plan | Roadmap, Feature catalog (trạng thái) |
| Đã làm gì | Report | Decision log, CHANGELOG, ADR |
| Hỗ trợ testing | Test case (từ mục 9 Feature Spec) | Testing strategy, Test giao diện, Test report |
| Brief cho người dùng | Release brief | Release notes, User guide |

## 6. Công cụ

- **Markdown trong repo**: mặc định cho mọi tầng. Cùng version control với code, AI đọc trực tiếp, review qua pull request.
- **Sinh tự động**: API spec từ OpenAPI, schema từ migration, Feature catalog từ frontmatter, test từ Gherkin, component library từ code, UI spec từ công cụ thiết kế. Không viết tay những gì sinh được.
- **Design tokens**: một file JSON theo chuẩn W3C Design Tokens làm nguồn sự thật, sinh CSS variables, Tailwind config hoặc theme object bằng Style Dictionary. Là việc đầu tiên của Tầng 3 trước mọi mockup. Figma Variables ánh xạ đúng tên token.
- **Lint giao diện**: cấm giá trị màu, khoảng cách gõ tay trong code UI và cấm import vượt lớp Atomic Design. Đây là chốt chặn thật của nhất quán.
- **Wireframe, mockup, prototype**: Figma là chuẩn ngành. Lưu liên kết trong repo kèm ảnh xuất tĩnh để AI đọc được và không phụ thuộc tài khoản Figma. Không có designer thì mockup dạng HTML tĩnh trong repo, AI sinh theo Design system, người duyệt; đủ cho công cụ nội bộ.
- **Site tĩnh** (Docusaurus, MkDocs, Mintlify): thêm khi có người đọc không kỹ thuật, hiển thị Tầng 2, 3, 4, 7.
- **Notion, Confluence**: chỉ khi phần lớn người đọc không dùng git, và chỉ cho Tầng 7 và Product overview. Còn lại giữ trong repo để AI truy cập trực tiếp.
- **CI**: kiểm tra liên kết hỏng, frontmatter thiếu, spec không có test tương ứng, code UI lệch design tokens.

## 7. Tách file

Quy tắc chung: **một đơn vị độc lập một file, áp dụng từ file đầu tiên, bất kể dự án có 1 hay 100 tính năng.** Lý do: AI phải tải cả file để dùng một phần, nhiều người sửa cùng file gây xung đột merge, phần đã chốt nằm chung phần đang sửa dễ bị động vào. Tách ngay từ đầu rẻ hơn tách sau.

Ngưỡng bổ sung cho file đơn: vượt 300 đến 500 dòng thì xem lại; trần cứng 800 dòng.

### Tài liệu bắt buộc tách

| Tài liệu | Đơn vị tách | Cách tổ chức |
|---|---|---|
| Feature Spec | Một tính năng một file | Thư mục `features/`, tên file theo mã tính năng. Feature catalog là chỉ mục sinh từ frontmatter |
| ADR | Một quyết định một file | Đánh số tăng dần, không sửa sau khi chốt. Chỉ mục liệt kê trạng thái |
| Ý tưởng mới | Một ý tưởng một thư mục | Ba file idea.md, interview.md, brief.md trong `intake/{ngày}-{slug}/`. Template gốc một file duy nhất |
| Change Request | Một CR một file | `cr/{ngày}-{slug}.md`, tách interview.md ra thư mục cùng tên khi dài. Template CR một file gốc. Chỉ mục CR sinh từ frontmatter trạng thái |
| Design system | Tokens một file JSON; mỗi component một thư mục theo lớp `atoms/`, `molecules/`, `organisms/`, `templates/`; mỗi pattern một file | Chỉ mục ở đầu. Tokens là nguồn sự thật, CSS variables và theme sinh ra từ đó. Tài liệu component sinh từ Storybook, không viết tay |
| Wireframe, Mockup | Một màn hình một file hoặc một frame | Đặt tên theo mã tính năng và mã bước |
| User flow | Một tính năng một file | Cùng mã tính năng với Feature Spec, hoặc nằm trong Feature Spec khi ngắn |
| Plan | Một đợt việc một thư mục | Chỉ mục ngắn, mỗi giai đoạn một file |
| Report | Một giai đoạn một file | Theo ngày giờ và slug, không gộp |
| Runbook | Một sự cố một file | Người trực mở đúng file, không cuộn tìm |
| Test case | Theo tính năng | Đi kèm Feature Spec tương ứng, hoặc file Gherkin cùng tên |
| Release brief | Một tính năng một file | Release notes gộp từ các file này khi phát hành |
| User guide | Theo nhiệm vụ người dùng | Mỗi nhiệm vụ một trang, mục lục riêng |
| API spec | Theo nhóm tài nguyên | OpenAPI tách file và tham chiếu bằng `$ref`. Sinh từ code nên ít khi tách tay |
| Incident postmortem | Một sự cố một file | Theo ngày |

### Tài liệu giữ một file

| Tài liệu | Lý do | Giới hạn |
|---|---|---|
| README | Điểm vào, đọc liền mạch | Dưới 150 dòng, dài hơn thì đẩy sang docs và để liên kết |
| Agent context file | AI tải mỗi phiên, tách ra AI không tự tìm | Dưới 200 dòng, chỉ quy tắc và con trỏ |
| Product overview | Một trang mới có giá trị | Dài hơn hai trang là đang lẫn Feature Spec vào |
| Architecture overview | Cần nhìn toàn cảnh một lần | Chi tiết từng thành phần tách sang ADR hoặc file riêng khi vượt 400 dòng |
| Glossary | Tra cứu theo bảng chữ cái | Vượt 300 thuật ngữ thì tách theo miền nghiệp vụ |
| Testing strategy | Chính sách, không phải danh sách | Một file |
| CHANGELOG | Thứ tự thời gian liên tục | Vượt 1000 dòng thì cắt theo năm hoặc phiên bản lớn |
| Roadmap | Toàn cảnh | Một file |
| Decision log | Nhật ký ngắn | Cắt theo quý khi dài |

### Quy tắc đi kèm

- Mỗi thư mục có file chỉ mục, sinh tự động hoặc ngắn gọn: tên, trạng thái, một dòng mô tả, liên kết. AI đọc chỉ mục trước rồi mở đúng file.
- Frontmatter thống nhất trên mọi file con: mã, tên, trạng thái, chủ sở hữu, ngày cập nhật. Định nghĩa trước khi viết file đầu tiên.
- Tên file mang mã định danh để liên kết chéo ổn định khi đổi tên hiển thị.
- Không tách theo mục con của một đơn vị: sơ đồ, bảng hành vi, tiêu chí chấp nhận của một Feature Spec phải ở cùng file vì ràng buộc bằng mã bước.
- Không tách vì file "trông dài": Product overview 150 dòng tốt hơn ba file 50 dòng phải nhảy qua lại.

## 8. Lộ trình áp dụng

| Giai đoạn | Tài liệu | Tiêu chí hoàn thành |
|---|---|---|
| 1. Nền tảng | Agent context file, README, template ý tưởng có mục Giao diện, template brief bốn mục cố định, template CR có bảng tác động, Product overview (chạy chuỗi ý tưởng cấp dự án), Glossary, mẫu Feature Spec có sơ đồ, frontmatter chung và cấu trúc thư mục tách file | Chạy thử một ý tưởng mới và một CR trên tính năng vừa tạo; AI sinh đúng câu hỏi từ trường trống và soạn đúng bảng tác động từ tài liệu hiện có; người duyệt đọc brief hoặc CR trong 15 phút; AI sinh được test và brief từ spec |
| 2. Design | Chuỗi ý tưởng giao diện cấp dự án (phỏng vấn có hình) ra Design brief, rồi bộ tối thiểu Design system theo thứ tự lớp 0 đến lớp 4 (tokens, quy tắc nền, 4 atoms, Form field, 3 organisms, 2 templates), lint giao diện, mẫu User flow và Mockup | AI sinh được một page từ template chỉ bằng organism và semantic tokens, lint bật không lỗi; mockup liên kết được theo mã bước từ Feature Spec |
| 3. Vận hành | Feature catalog, cặp Plan và Report, Testing strategy, Test giao diện | AI viết Report sau mỗi giai đoạn, catalog phản ánh đúng trạng thái, test giao diện so được với mockup |
| 4. Người dùng | Mẫu Release brief, CHANGELOG, Release notes | Brief sinh từ Feature Spec không cần viết lại từ đầu |
| 5. Tự động hóa | Sinh catalog từ frontmatter, sinh test từ Gherkin, component library từ code, CI kiểm tra liên kết và lệch tokens | Không còn tài liệu viết tay có thể sinh được |
| 6. Production | Tầng 8 vận hành | Có runbook cho mọi sự cố đã gặp, DR đã diễn tập |

Không có giao diện người dùng thì bỏ giai đoạn 2 và các mục test giao diện.

## 9. Ưu tiên theo quy mô

| Quy mô | Bắt buộc | Nên có | Bỏ qua được |
|---|---|---|---|
| Cá nhân / MVP | README, Agent context file, Feature Spec mỗi tính năng một file (nội dung rút gọn), Design tokens, ADR ngắn, CHANGELOG | Architecture 1 trang, mockup cho màn hình phức tạp | Wireframe, Prototype, Tầng 7, 8, 9 |
| Team 3 đến 10 người | Tầng 0 đến 6 đầy đủ | Release brief, Threat model, Deployment guide, Prototype cho tính năng lớn | Tầng 9 hình thức |
| Doanh nghiệp / có khách hàng | Tầng 0 đến 8 | Tầng 9 theo yêu cầu ngành | Không có |

## 10. Khi nào khuyến nghị không còn đúng

- Dự án dưới 5 tính năng hoặc dưới 2 tháng: bỏ Feature catalog, Release brief và Prototype. Vẫn giữ mỗi tính năng một file.
- idea.md điền đủ, không mâu thuẫn: bỏ phỏng vấn, AI soạn thẳng brief. CR sửa lỗi nhỏ: bảng tác động hai dòng, không hỏi đáp.
- Dùng thư viện UI có sẵn và chấp nhận mặc định: bỏ chuỗi ý tưởng giao diện cấp dự án, Design brief rút xuống một trang ghi thư viện, ngoại lệ, giọng văn.
- Luồng tuyến tính dưới bốn bước hoặc CRUD: bỏ sơ đồ, bảng hành vi là đủ.
- Không có giao diện người dùng (API, CLI, batch job): bỏ toàn bộ Tầng 3 và test giao diện, giữ giọng văn cho thông báo lỗi trong Glossary.
- Công cụ nội bộ dùng thư viện UI có sẵn: Design system rút gọn thành "dùng thư viện X, các ngoại lệ sau", bỏ wireframe, mockup chỉ cho màn hình phức tạp.
- Dùng Notion hoặc Confluence: đơn vị là trang, ngưỡng dòng không áp dụng, thay bằng một trang một đơn vị.
- Công cụ nội bộ không có người dùng ngoài: bỏ Tầng 7, CHANGELOG là đủ.
- Yêu cầu thay đổi hàng ngày: giữ Feature Spec ở mức mục 1, 2, 9; chi tiết hành vi để trong test.
- Ngành có quy chuẩn riêng (y tế, tài chính, chính phủ): danh sách do quy chuẩn ngành quyết định.
- Team hoàn toàn không dùng AI: bỏ Agent context file, Report thay bằng ghi chú trong công cụ quản lý task.
- Dự án đã chạy lâu, tài liệu tiếp nhận ban đầu không còn ai đọc: xếp lại theo thẩm quyền, Tầng 1 chuyển xuống thành phụ lục quy trình.

## 11. Bước tiếp theo

1. Chốt danh sách tài liệu Giai đoạn 1 và người chịu trách nhiệm từng tài liệu.
2. Viết template ý tưởng có mục Giao diện, template brief bốn mục cố định, template CR có bảng tác động, mẫu Feature Spec có sơ đồ; mỗi template dưới một trang. Chạy chuỗi ý tưởng cấp dự án để ra Product overview, rồi chạy thử một ý tưởng tính năng thật và một CR trên chính tính năng đó, đo xem AI có sinh đúng câu hỏi từ trường trống và đúng bảng tác động không, người duyệt đọc trong 15 phút, AI sinh được test và brief từ spec.
3. Nếu có giao diện: chốt thang khoảng cách và thang chữ trước, viết tokens JSON lớp primitive và semantic, sinh CSS variables, bật lint cấm giá trị gõ tay, rồi cho AI sinh một màn hình chỉ bằng semantic tokens.
4. Nếu đạt, nhân rộng và đưa quy ước vào Agent context file, kèm ngưỡng dòng và quy tắc tách file.
5. Viết script sinh chỉ mục cho `features/` và `adr/` ngay từ đầu, không viết tay chỉ mục.

## Câu hỏi chưa giải quyết

- Dự án áp dụng thuộc quy mô nào trong mục 9? Quyết định này chốt danh sách tối thiểu.
- Phần mềm có giao diện người dùng không, và có designer không? Quyết định Tầng 3 nặng hay nhẹ, và mockup vẽ bằng Figma hay AI sinh HTML.
- Người đọc chính của Tầng 7 có dùng git không? Quyết định có cần site tĩnh hay Notion ngay từ đầu.
- Định dạng Feature Spec chọn Markdown thuần hay Gherkin để test chạy trực tiếp?
