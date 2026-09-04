## Tài liệu dự án: `dk`

Mọi tài liệu trong `docs/` do CLI `dk` tạo từ template nhúng; agent không tự tạo file
tay ở đó. Ba lớp: **template** (khung, sửa bằng phát hành `dk` mới), **CLI** (việc xác
định: tạo file, changelog, render, chỉ mục, kiểm tra), **skill** (việc cần suy luận:
phỏng vấn, soạn nội dung, phân tích tác động). Thiếu `dk` thì dừng và báo người cài.

### Skill (cài bằng `dk skill install`)

| Skill | Dùng khi |
|---|---|
| `doc-intake` | ý tưởng mới chưa có gì: idea, phỏng vấn, brief chờ duyệt (`docs/intake/`) |
| `doc-cr` | thay đổi trên thứ đã có: Change Request, bảng tác động chờ duyệt (`docs/cr/`) |
| `doc-overview` | Product overview, Architecture overview, Glossary (`docs/overview/`) |
| `doc-adr` | quyết định kỹ thuật lớn, đánh số, bất biến sau khi chốt (`docs/adr/`) |
| `doc-feature-spec` | Feature Spec 11 mục từ brief đã duyệt hoặc CR đã chốt (`docs/features/`) |
| `doc-design-system` | tokens, foundations, atoms đến templates, patterns (`docs/design/`) |
| `doc-design-flow` | user flow, wireframe, mockup HTML, prototype cho một tính năng |
| `doc-test` | testing strategy, test case Gherkin hoặc bảng, checklist UI, test report (`docs/test/`) |
| `doc-plan-report` | roadmap, plan và phase, report có bằng chứng, decision log, CHANGELOG sản phẩm |
| `doc-release` | release brief, release notes, user guide, FAQ cho người dùng cuối (`docs/release/`) |
| `doc-ops` | deployment, environment, runbook, monitoring, postmortem, backup và DR (`docs/ops/`) |

### Lệnh `dk` chính

```
dk new <loại> <slug> [--from <file>] [--set k=v]   # tạo file từ template, điền frontmatter
dk changelog add <file> --summary "<tóm tắt thật>" --source <CR-id|brief>
dk render <file> | --all --index                    # HTML tự chứa vào docs/html/
dk index [features|adr|cr|intake|user-guide|all]    # README.md chỉ mục sinh ra
dk check [<file>]                                   # mã thoát 3 khi có lỗi
dk refs <file>                                      # liên kết đi và đến
dk status | dk doctor                               # tổng quan tài liệu | kiểm cài đặt
```

### Quy tắc

- Sau mỗi lần sửa một tài liệu: `changelog add` với tóm tắt nói nội dung, `render`,
  `index` thư mục đó, `check` file đó. Pre-commit chặn commit khi thiếu dòng changelog.
- Trạng thái dùng khóa tiếng Anh của `types.toml` (`draft`, `review`, `approved`...);
  chỉ người đổi sang trạng thái đã chốt. Skill dừng ở "chờ duyệt".
- Không sửa `created`, `created_by`, `dk_version` trong frontmatter; không sửa thân
  ADR đã `accepted`; không sửa file `generated: true` hay `docs/html/`.
- Ngưỡng dòng: cảnh báo trên 500, lỗi trên 800 (`dk.toml` `[check]`); loại có ngưỡng
  riêng ghi trong `types.toml`.
