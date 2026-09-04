package check

import (
	"fmt"
	"strings"

	"github.com/anhvandev/doc-kit/internal/frontmatter"
	"github.com/anhvandev/doc-kit/internal/render"
)

// userflowSteps: mã bước trong sơ đồ và bảng của userflow phải là tập con mã
// bước của feature-spec có id bằng trường feature; feature không trỏ đến spec
// nào là lỗi. Spec không có mã bước (CRUD) thì userflow cũng không được có.
func userflowSteps(c *Context) []Finding {
	var out []Finding
	specCodes := map[string]map[string]bool{}
	for _, m := range c.typed() {
		if m.Type != "feature-spec" || m.ID == "" {
			continue
		}
		diagram, table := render.StepCodes(m.Body)
		all := map[string]bool{}
		for _, x := range append(diagram, table...) {
			all[x] = true
		}
		specCodes[m.ID] = all
	}
	for _, m := range c.typed() {
		if m.Type != "userflow" {
			continue
		}
		feature := frontmatter.GetString(m.FM, "feature")
		if feature == "" {
			continue // frontmatter-required đã báo
		}
		spec, ok := specCodes[feature]
		if !ok {
			out = append(out, finding(m, "userflow-steps", Error, fmt.Sprintf("feature %q không trỏ đến Feature Spec nào", feature)))
			continue
		}
		diagram, table := render.StepCodes(m.Body)
		var bad []string
		for _, x := range append(diagram, table...) {
			if !spec[x] {
				bad = append(bad, x)
			}
		}
		if len(bad) > 0 {
			out = append(out, finding(m, "userflow-steps", Error, fmt.Sprintf("mã bước không có trong %s: %s", feature, strings.Join(unique(bad), ", "))))
		}
	}
	return out
}
