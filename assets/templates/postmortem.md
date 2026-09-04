---
type: postmortem
title: ""
status: draft
owner: ""
created: {{.Created}}
updated: {{.Updated}}
source: ""
created_by: dk
dk_version: {{.DKVersion}}
incident_at: ""
written_within_48h: false
---

# Postmortem: {{.Title}}

<!-- gợi ý: mục 1 là frontmatter. Viết trong 48 giờ sau sự cố (written_within_48h do dk new tính từ incident_at); quá hạn vẫn viết, ghi lý do ở mục 7. Không đổ lỗi cá nhân: nói về hệ thống và quy trình, không nêu tên người gây ra. Mỗi hành động khắc phục có người chịu trách nhiệm và hạn. status: draft | review | closed (mọi hành động đã xong). -->

## 2. Tóm tắt

<!-- gợi ý: ba câu: chuyện gì xảy ra, ảnh hưởng đến ai bao lâu, đã xử lý thế nào -->

## 3. Ảnh hưởng

<!-- gợi ý: số người dùng, thời gian gián đoạn, dữ liệu mất hoặc sai, tiền -->

- Bắt đầu: 
- Kết thúc: 
- Ảnh hưởng: 

## 4. Timeline

<!-- gợi ý: giờ phút, sự kiện quan sát được hoặc hành động; lấy từ log, chat, alert; không suy đoán -->

| Giờ | Sự kiện |
|---|---|
| | |

## 5. Nguyên nhân gốc

<!-- gợi ý: hỏi "vì sao" đến khi chạm hệ thống hoặc quy trình; nguyên nhân kích hoạt và nguyên nhân gốc là hai thứ khác nhau -->

- Nguyên nhân kích hoạt: 
- Nguyên nhân gốc: 
- Vì sao không phát hiện sớm: 

## 6. Hành động khắc phục

<!-- gợi ý: mỗi hành động có người và hạn; loại: ngăn tái diễn | phát hiện sớm hơn | giảm ảnh hưởng -->

| Hành động | Loại | Người | Hạn | Trạng thái |
|---|---|---|---|---|
| | | | | |

## 7. Bài học

<!-- gợi ý: điều làm tốt, điều may mắn, điều cần đổi; lý do nếu viết quá 48 giờ -->

- 
