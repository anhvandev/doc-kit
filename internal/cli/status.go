package cli

import (
	"sort"
	"time"

	"github.com/spf13/cobra"

	"github.com/anhvandev/doc-kit/internal/check"
	"github.com/anhvandev/doc-kit/internal/docs"
	"github.com/anhvandev/doc-kit/internal/frontmatter"
	"github.com/anhvandev/doc-kit/internal/gitx"
)

func newStatusCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Bảng tổng: tài liệu theo loại và trạng thái, CR mở, changelog pending, finding của check",
		Args:  exactArgs(0),
		RunE:  func(_ *cobra.Command, _ []string) error { return a.runStatus() },
	}
}

type statusReport struct {
	Docs             map[string]map[string]int `json:"docs"`
	OpenCR           int                       `json:"open_cr"`
	ChangelogPending *int                      `json:"changelog_pending"` // null khi không có git
	Check            map[string]int            `json:"check"`
	// DROverdue: file backup-dr có last_drill rỗng hoặc quá 6 tháng; last_drill
	// sai định dạng cũng tính, kèm ghi chú riêng ở bản in.
	DROverdue []string `json:"dr_overdue"`
	drBadDate map[string]bool
}

// drDrillMonths là số tháng tối đa giữa hai lần diễn tập khôi phục.
const drDrillMonths = 6

func (a *app) runStatus() error {
	if err := a.requireProject(); err != nil {
		return err
	}
	metas, err := docs.Scan(a.root, a.cfg.DocsDir)
	if err != nil {
		return fail(codeError, "%v", err)
	}
	r := statusReport{Docs: map[string]map[string]int{}, Check: map[string]int{}, DROverdue: []string{}}
	cutoff := time.Now().AddDate(0, -drDrillMonths, 0)
	for _, m := range metas {
		if !m.HasFM || m.Generated || m.Type == "" {
			continue
		}
		if m.Type == "backup-dr" {
			raw := frontmatter.GetString(m.FM, "last_drill")
			if d, err := time.Parse("2006-01-02", raw); err != nil || d.Before(cutoff) {
				r.DROverdue = append(r.DROverdue, m.Rel)
				if raw != "" && err != nil {
					if r.drBadDate == nil {
						r.drBadDate = map[string]bool{}
					}
					r.drBadDate[m.Rel] = true
				}
			}
		}
		if r.Docs[m.Type] == nil {
			r.Docs[m.Type] = map[string]int{}
		}
		r.Docs[m.Type][m.Status]++
		if m.Type == "cr" && m.Status != "closed" && m.Status != "rejected" {
			r.OpenCR++
		}
	}
	if gitx.IsRepo(a.root) {
		pending, err := a.pendingDocs()
		if err != nil {
			return err
		}
		n := len(pending)
		r.ChangelogPending = &n
	}
	findings, _, err := a.runChecks(false)
	if err != nil {
		return err
	}
	r.Check[check.Error], r.Check[check.Warning] = check.Count(findings)
	if a.json {
		return a.printJSON(r)
	}
	types := make([]string, 0, len(r.Docs))
	for t := range r.Docs {
		types = append(types, t)
	}
	sort.Strings(types)
	a.printf("Tài liệu theo loại và trạng thái:\n")
	for _, t := range types {
		statuses := make([]string, 0, len(r.Docs[t]))
		for s := range r.Docs[t] {
			statuses = append(statuses, s)
		}
		sort.Strings(statuses)
		total := 0
		for _, s := range statuses {
			total += r.Docs[t][s]
		}
		a.printf("  %-14s %d", t, total)
		for _, s := range statuses {
			a.printf("  %s=%d", s, r.Docs[t][s])
		}
		a.printf("\n")
	}
	if len(types) == 0 {
		a.printf("  (chưa có)\n")
	}
	a.printf("CR đang mở: %d\n", r.OpenCR)
	if r.ChangelogPending != nil {
		a.printf("Changelog pending: %d\n", *r.ChangelogPending)
	} else {
		a.printf("Changelog pending: không có git\n")
	}
	a.printf("Check: %d lỗi, %d cảnh báo\n", r.Check[check.Error], r.Check[check.Warning])
	for _, f := range r.DROverdue {
		if r.drBadDate[f] {
			a.printf("DR last_drill sai định dạng (cần yyyy-mm-dd): %s\n", f)
			continue
		}
		a.printf("DR chưa diễn tập quá %d tháng: %s\n", drDrillMonths, f)
	}
	return nil
}
