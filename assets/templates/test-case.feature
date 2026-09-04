# dk:
# type: test-case
# title: ""
# status: draft
# owner: ""
# created: {{.Created}}
# updated: {{.Updated}}
# source: ""
# created_by: dk
# dk_version: {{.DKVersion}}
# feature: ""

# Sinh từ mục 9 của Feature Spec {{.Feature}} bằng dk new test-case --from. Mỗi tiêu chí chấp nhận một Scenario, tag @<mã tính năng> @AC<n> để truy vết; thêm Scenario cho ngoại lệ ở mục 7 của spec với tag @E<n>. Scenario có dòng "# chưa tách được" là AC lệch khung Given / When / Then trong spec: sửa spec hoặc điền tay. File này là tài liệu và là test: không sửa Feature Spec để hợp thức hóa test.

@{{.Feature}}
Feature: {{.Title}}
{{if .Background}}
  Background:
{{range $i, $b := .Background}}    {{if eq $i 0}}Given{{else}}And{{end}} {{$b}}
{{end}}{{end}}{{range .Scenarios}}
  @{{$.Feature}} @{{.Code}}
  Scenario: {{.Code}}{{if .Title}} {{.Title}}{{end}}
{{if .Raw}}    # chưa tách được: {{.Raw}}
    Given TODO
    When TODO
    Then TODO
{{else}}    Given {{.Given}}
    When {{.When}}
    Then {{.Then}}
{{end}}{{else}}
  @{{.Feature}} @AC1
  Scenario: AC1
    Given TODO
    When TODO
    Then TODO
{{end}}
