package check

import (
	"fmt"
	"time"

	"github.com/anhvandev/doc-kit/internal/docs"
)

// updatedLayouts là các dạng của trường updated, thử theo thứ tự. Dạng không
// có offset đọc theo giờ máy đang chạy check; dạng có offset so đúng thời điểm.
var updatedLayouts = []string{"2006-01-02 15:04 -0700", "2006-01-02 15:04", time.RFC3339, "2006-01-02"}

// crApprovalOrder: Feature Spec có source là CR chưa chốt (status ngoài final
// của loại cr) mà updated mới hơn updated của CR là lỗi: spec đã bị sửa trước
// khi CR được duyệt. updated không đọc được là warning.
func crApprovalOrder(c *Context) []Finding {
	var out []Finding
	crType, ok := c.Reg["cr"]
	if !ok {
		return nil
	}
	for _, m := range c.typed() {
		if m.Type != "feature-spec" || m.Source == "" {
			continue
		}
		cr, found := docs.Resolve(c.Metas, m.Source)
		if !found || cr.Type != "cr" || contains(crType.Final, cr.Status) {
			continue
		}
		specAt, err := parseUpdated(m.Updated)
		if err != nil {
			out = append(out, finding(m, "cr-approval-order", Warning, "updated không đọc được: "+err.Error()))
			continue
		}
		crAt, err := parseUpdated(cr.Updated)
		if err != nil {
			out = append(out, finding(m, "cr-approval-order", Warning, fmt.Sprintf("updated của %s không đọc được: %v", cr.ID, err)))
			continue
		}
		if specAt.After(crAt) {
			out = append(out, finding(m, "cr-approval-order", Error, fmt.Sprintf("sửa lúc %s nhưng %s còn %s (cập nhật %s); chờ CR được duyệt rồi mới sửa spec", m.Updated, cr.ID, cr.Status, cr.Updated)))
		}
	}
	return out
}

func parseUpdated(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, fmt.Errorf("trống")
	}
	for _, l := range updatedLayouts {
		if t, err := time.ParseInLocation(l, s, time.Local); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("%q không theo dạng YYYY-MM-DD HH:MM", s)
}
