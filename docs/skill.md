# Skill và quy tắc kiểm tra

## 11 skill nhúng

Cài bằng `dk skill install` (mặc định target `claude`, scope dự án; `--global` vào
`~/.claude/`). Mỗi skill: `SKILL.md` dưới 300 dòng, `references/rules.md`, bước 0 kiểm
`dk --version`, câu dừng bắt buộc in đậm đầu file, ba phần `new`, `update`, `html`.

| Skill | Thư mục | Sản phẩm | Dừng chờ người |
|---|---|---|---|
| `doc-intake` | `docs/intake/<yymmdd>-<slug>/` | idea, interview, brief (Intake, Product, Design) | brief `review` |
| `doc-cr` | `docs/cr/` | CR với bảng tác động từ `dk refs`, interview khi cần | bảng tác động `review` |
| `doc-overview` | `docs/overview/` | Product overview từ Product brief, Architecture từ mã, Glossary; Feature catalog bằng `dk index features` | brief chưa `approved` |
| `doc-adr` | `docs/adr/` | ADR đánh số, thân bất biến sau `accepted`, thay bằng ADR mới | ADR `proposed` |
| `doc-feature-spec` | `docs/features/` | Feature Spec 11 mục, 5 biến thể `format` | spec `review`; chỉ sửa theo CR `approved` |
| `doc-design-system` | `docs/design/` | tokens, foundations, atoms đến templates, patterns | Design brief chưa `approved` |
| `doc-design-flow` | `docs/design/flows`, `wireframes`, `mockups` | user flow, wireframe, mockup mỗi trạng thái, prototype, UI spec | wireframe trước mockup |
| `doc-test` | `docs/test/` | strategy, test case `.feature` hoặc bảng, checklist UI, test report | strategy chờ chốt; `bdd_cmd` rỗng thì "chưa kiểm chạy được" |
| `doc-plan-report` | `docs/plan/`, `plans/` | roadmap, plan và phase, report có bằng chứng, decision log, CHANGELOG sản phẩm | repo có công cụ plan riêng thì chỉ viết report |
| `doc-release` | `docs/release/` | release brief, release notes `--collect`, user guide, FAQ | brief `ready` do người đặt |
| `doc-ops` | `docs/ops/` | deployment, environment, runbook, monitoring, postmortem, backup và DR | postmortem quá 48 giờ |

Tầng quản trị (`charter`, `risk-register`, `meeting-notes` trong `docs/governance/`) chỉ có
template, không có skill: tài liệu quản trị do người viết.

## Quy ước chung của skill

- Trạng thái dùng khóa tiếng Anh của `types.toml`; `review` là chờ duyệt, `approved` là đã
  chốt. Skill không tự đổi sang trạng thái đã chốt.
- Sau mỗi lần sửa một tài liệu: `dk changelog add` với tóm tắt nói nội dung, `dk render`,
  `dk index <thư mục>`, `dk check <file>`.
- Skill giao việc cập nhật tài liệu đích cho "họ" tài liệu kèm `--source <CR-id>`, không
  nêu tên skill khác (trừ `doc-intake` chỉ sang `doc-cr` khi brief đã `approved`).
- Nội dung trung lập target: không nhắc tên harness, tên tool, thư mục `.claude/`. Kiểm
  bằng `make lint-skills` và `internal/skill/content_test.go` (cả hai chạy trong CI); test
  còn bắt hai mô tả kích hoạt trùng 3 từ liên tiếp (bỏ câu "Không dùng ..." cuối).

## Khối ngữ cảnh agent

`dk init --agent-context` ghi `assets/agent-context.md` (tiếng Anh, dưới 120 dòng) vào
`CLAUDE.md` và `AGENTS.md`: quy tắc làm việc chung của agent (nghĩ trước khi code, tối
giản, sửa đúng chỗ, mục tiêu kiểm được, trả lời đúng câu hỏi) rồi phần `dk` (ba lớp,
bảng skill, lệnh chính, quy tắc tài liệu, ngưỡng dòng). Khối nằm giữa hai dấu mốc HTML
comment mang phiên bản `dk` và hash nội dung (`internal/agentctx`); chạy lại thay đúng
khối, phần còn lại của file giữ nguyên; `dk doctor` báo khối thiếu, cũ hay bị sửa tay.

## 16 quy tắc của `dk check`

Quét `docs/` và `plans/`; mã thoát 3 khi có lỗi, `--strict` coi warning là lỗi.

| Quy tắc | Mức | Áp cho | Nội dung |
|---|---|---|---|
| `frontmatter-required` | lỗi | mọi loại | thiếu trường `required` của loại; `created_by` khác `dk`; `.feature` thiếu khối `# dk:` |
| `status-valid` | lỗi | mọi loại | `status` ngoài `statuses` của loại |
| `link-broken` | lỗi | Markdown | liên kết tương đối trỏ file không có |
| `step-codes` | lỗi | Feature Spec có sơ đồ | mã bước trong sơ đồ và bảng hành vi lệch nhau |
| `spec-section-order` | lỗi | Feature Spec | tiêu đề `## N.` thiếu, lạ, lặp, sai thứ tự (bỏ 6 khi `has_ui: false`, bỏ 4 khi `crud`) |
| `cr-approval-order` | lỗi | Feature Spec có `source` là CR | spec sửa sau khi CR chưa `approved` |
| `backlink` | lỗi / cảnh báo | có `source` / đã chốt | `source` không trỏ đến đâu; tài liệu chốt không ai trỏ về |
| `spec-has-test` | cảnh báo | Feature Spec `approved` trở lên | chưa có test case trong `docs/test/` có `source` là spec |
| `line-threshold` | cảnh báo / lỗi | Markdown trong `docs/` | vượt `warn_lines` (theo loại hoặc `dk.toml`) / `max_lines` |
| `adr-immutable` | lỗi | ADR chốt tại HEAD | thân khác bản HEAD |
| `glossary-term` | cảnh báo | Feature Spec | chữ in đậm mục 2, 5, 8 chưa có trong Glossary |
| `mockup-tokens` | lỗi | mockup `.html` | hex hoặc px trong `<style>` hay `style=`; thiếu khối `<!-- dk: -->` |
| `userflow-steps` | lỗi | userflow | mã bước không phải tập con của spec cùng `feature` |
| `report-evidence` | cảnh báo | report | không có mã commit, liên kết file output hay khối kết quả |
| `no-jargon` | cảnh báo | release brief, user guide | từ trong `[release] jargon` của `dk.toml` |
| `env-no-secret` | lỗi | environment | dòng `KEY=value` có giá trị thật, không phải `<...>` |

Rule là lớp phụ: `env-no-secret` không thay công cụ quét secret; `mockup-tokens` chỉ bắt
hex và px; `no-jargon` không tắt được bằng danh sách rỗng (rỗng lấy mặc định).
