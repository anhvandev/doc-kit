---
id: F-001
type: feature-spec
title: "Bộ lọc: đơn hàng"
status: draft
owner: vanna
source: CR-260910-loc
acceptance: []
---

# F-001: Bộ lọc đơn hàng

Xem [ADR](../adr/ADR-0001-x.md#muc-2), [plan](../../plans/p.md), [web](https://example.com/a.md) và [neo](#x).

```mermaid
flowchart TD
    B1[B1 Mở] --> B2[B2 Nhập]
    B2 --> B1
    B2 --> B3{B3 Hợp lệ?}
```

## 5. Giao diện

| Mã | Màn hình |
|---|---|
| B1 | Danh sách |

## 4. Hành vi

| Mã | Hành động | Phản hồi |
|---|---|---|
| B1 | mở | thấy |
| B2 | nhập | ok |
| B9 | thừa | |

```go
fmt.Println("<b>")
```

<script>alert(1)</script>
