package check

import (
	"strings"

	"github.com/anhvandev/doc-kit/internal/render"
)

// stepCodes: feature-spec: tập mã bước trong sơ đồ bằng tập trong bảng hành vi.
// Spec cố ý không có sơ đồ (luồng tuyến tính, CRUD) chỉ có bảng thì bỏ qua.
// Userflow so tập con với spec ở quy tắc userflow-steps.
func stepCodes(c *Context) []Finding {
	var out []Finding
	for _, m := range c.typed() {
		if m.Type != "feature-spec" {
			continue
		}
		diagram, table := render.StepCodes(m.Body)
		if len(diagram) == 0 {
			continue
		}
		if missing := diff(diagram, table); len(missing) > 0 {
			out = append(out, finding(m, "step-codes", Error, "mã có trong sơ đồ nhưng thiếu ở bảng hành vi: "+strings.Join(missing, ", ")))
		}
		if extra := diff(table, diagram); len(extra) > 0 {
			out = append(out, finding(m, "step-codes", Error, "mã có trong bảng hành vi nhưng không có trong sơ đồ: "+strings.Join(extra, ", ")))
		}
	}
	return out
}

// diff trả về phần tử của a không có trong b, giữ thứ tự.
func diff(a, b []string) []string {
	set := map[string]bool{}
	for _, x := range b {
		set[x] = true
	}
	var out []string
	for _, x := range a {
		if !set[x] {
			out = append(out, x)
		}
	}
	return unique(out)
}

func unique(a []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, x := range a {
		if !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	return out
}
