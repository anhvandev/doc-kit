package cli

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/anhvandev/doc-kit/internal/changelog"
	"github.com/anhvandev/doc-kit/internal/docs"
)

func newNewCmd(a *app) *cobra.Command {
	var (
		from    string
		in      string
		appnd   string
		collect string
		sets    []string
		force   bool
	)
	cmd := &cobra.Command{
		Use:   "new <type> [<slug>]",
		Short: "Tạo tài liệu mới từ template theo loại",
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) == 2 || (len(args) == 1 && (appnd != "" || collect != "")) {
				return nil
			}
			return fail(codeUsage, "new cần <type> và <slug> (slug bỏ được khi --append hoặc --collect)")
		},
		RunE: func(_ *cobra.Command, args []string) error {
			set, err := parseSets(sets)
			if err != nil {
				return err
			}
			slug := args[0]
			if len(args) == 2 {
				slug = args[1]
			} else if collect != "" {
				slug = docs.VersionSlug(collect)
			}
			return a.runNew(args[0], slug, newFlags{from: from, in: in, appnd: appnd, collect: collect, set: set, force: force})
		},
	}
	cmd.Flags().StringVar(&from, "from", "", "file nguồn để chép trường theo types.toml")
	cmd.Flags().StringVar(&in, "in", "", "thư mục plan chứa file (plan-phase, report)")
	cmd.Flags().StringVar(&appnd, "append", "", "decision-log, faq: nối một dòng vào cuối file")
	cmd.Flags().StringVar(&collect, "collect", "", "release-notes: phiên bản; gom Release brief ready chưa có released_in và ghi released_in vào chúng")
	cmd.Flags().StringArrayVar(&sets, "set", nil, "đặt trường frontmatter, dạng k=v (lặp được)")
	cmd.Flags().BoolVar(&force, "force", false, "ghi đè nếu file đích đã có")
	return cmd
}

type newFlags struct {
	from, in, appnd, collect string
	set                      map[string]string
	force                    bool
}

func parseSets(kvs []string) (map[string]string, error) {
	m := map[string]string{}
	for _, kv := range kvs {
		k, v, ok := strings.Cut(kv, "=")
		if !ok || k == "" {
			return nil, fail(codeUsage, "--set cần dạng k=v, nhận %q", kv)
		}
		m[k] = v
	}
	return m, nil
}

func (a *app) runNew(typeName, slug string, f newFlags) error {
	if err := a.requireProject(); err != nil {
		return err
	}
	if _, err := a.reg.Get(typeName); err != nil {
		return fail(codeUsage, "%v", err)
	}
	if !docs.ValidSlug(slug) {
		return fail(codeUsage, "slug %q không hợp lệ: chỉ a-z, 0-9 và dấu gạch ngang", slug)
	}
	if strings.ContainsAny(f.appnd, "\n\r") {
		return fail(codeUsage, "--append không được chứa xuống dòng")
	}
	now := time.Now()
	res, err := docs.New(a.reg, typeName, slug, docs.Options{
		DocsDir:  a.docsDir(),
		PlansDir: filepath.Join(a.root, filepath.FromSlash(a.cfg.PlansDir)),
		In:       a.resolve(f.in, ""),
		From:     a.resolve(f.from, ""),
		Append:   f.appnd,
		Collect:  f.collect,
		Now:      now,
		Set:      f.set,
		Force:    f.force,
		Owner:    a.cfg.DefaultOwner,
		IDPrefix: a.cfg.IDPrefix,
		Version:  Version,
	})
	if err != nil {
		return fail(codeError, "%v", err)
	}
	// --collect sửa từng brief (released_in), là sửa tài liệu khác nên dk tự ghi
	// changelog cho mỗi brief với nguồn là phiên bản.
	for i, p := range res.Released {
		relDocs, _ := filepath.Rel(a.docsDir(), p)
		if _, err := changelog.Record(a.root, a.docsDir(), filepath.ToSlash(relDocs), "Phát hành trong "+f.collect, f.collect, now, false); err != nil {
			return fail(codeError, "ghi changelog cho %s: %v", a.relRoot(p), err)
		}
		res.Released[i] = a.relRoot(p)
	}
	rel := a.relRoot(res.Path)
	if a.json {
		res.Path = rel
		return a.printJSON(res)
	}
	if res.Appended {
		a.printf("Đã nối dòng vào %s\n", rel)
		return nil
	}
	a.printf("Đã tạo %s", rel)
	if res.ID != "" {
		a.printf(" (id %s)", res.ID)
	}
	if f.appnd != "" {
		a.printf(" và nối dòng")
	}
	if res.Scenarios > 0 || res.Unparsed > 0 {
		a.printf(": %d Scenario, %d dòng AC chưa tách được", res.Scenarios, res.Unparsed)
	}
	if len(res.Released) > 0 {
		a.printf(": gom %d Release brief, đã ghi released_in", len(res.Released))
	}
	a.printf("\n")
	return nil
}
