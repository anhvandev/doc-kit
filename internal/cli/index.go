package cli

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anhvandev/doc-kit/internal/docs"
	"github.com/anhvandev/doc-kit/internal/frontmatter"
)

// indexDirs là các chỉ mục Markdown sinh ra, theo thứ tự chạy khi `all`;
// indexDirOf ánh xạ tên chỉ mục sang thư mục trong docs/.
var (
	indexDirs  = []string{"features", "adr", "cr", "intake", "user-guide"}
	indexDirOf = map[string]string{"features": "features", "adr": "adr", "cr": "cr", "intake": "intake", "user-guide": "release/guide"}
)

func newIndexCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "index [features|adr|cr|intake|user-guide|all]",
		Short: "Sinh README.md chỉ mục (generated: true) trong các thư mục docs/",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return fail(codeUsage, "index nhận tối đa một tham số")
			}
			if len(args) == 1 && args[0] != "all" && !contains(indexDirs, args[0]) {
				return fail(codeUsage, "index chỉ nhận %s hoặc all", strings.Join(indexDirs, ", "))
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			which := "all"
			if len(args) == 1 {
				which = args[0]
			}
			return a.runIndex(which)
		},
	}
}

func (a *app) runIndex(which string) error {
	if err := a.requireProject(); err != nil {
		return err
	}
	dirs := indexDirs
	if which != "all" {
		dirs = []string{which}
	}
	written := []string{}
	for _, name := range dirs {
		d := indexDirOf[name]
		rel := a.cfg.DocsDir + "/" + d
		metas, err := docs.Scan(a.root, rel)
		if err != nil {
			return fail(codeError, "%v", err)
		}
		out := filepath.Join(a.docsDir(), filepath.FromSlash(d), "README.md")
		if err := writeFile(out, a.indexMarkdown(name, filepath.Dir(out), metas)); err != nil {
			return fail(codeError, "%v", err)
		}
		written = append(written, a.relRoot(out))
	}
	if a.json {
		return a.printJSON(written)
	}
	for _, w := range written {
		a.printf("%s\n", w)
	}
	return nil
}

// indexMarkdown dựng chỉ mục. features (Feature catalog) và adr: một bảng
// phẳng sắp theo mã, trạng thái là cột. cr và intake: nhóm theo trạng thái
// theo thứ tự statuses của loại; trong nhóm sắp theo ngày tạo giảm dần rồi
// đường dẫn.
func (a *app) indexMarkdown(name, dir string, metas []docs.Meta) []byte {
	var items []docs.Meta
	for _, m := range metas {
		if !m.HasFM || m.Generated {
			continue
		}
		// Loại có dir khác thư mục này (ví dụ interview trong cr/<CR>/) không
		// phải đơn vị của chỉ mục; loại lạ vẫn liệt kê.
		if t, ok := a.reg[m.Type]; ok && t.Dir != indexDirOf[name] {
			continue
		}
		items = append(items, m)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "---\ngenerated: true\ntitle: Chỉ mục %s\n---\n\n# Chỉ mục %s\n\n", name, name)
	fmt.Fprintf(&sb, "<!-- sinh bởi `dk index %s`; không sửa tay -->\n", name)
	if len(items) == 0 {
		sb.WriteString("\nChưa có tài liệu.\n")
		return []byte(sb.String())
	}
	if name == "features" || name == "adr" {
		flatRows(&sb, name, dir, items)
		return []byte(sb.String())
	}
	if name == "user-guide" {
		guideRows(&sb, dir, items)
		return []byte(sb.String())
	}
	if name == "intake" {
		items = intakeRows(items)
	}
	order := a.statusOrder(items)
	sort.SliceStable(items, func(i, j int) bool {
		if order[items[i].Status] != order[items[j].Status] {
			return order[items[i].Status] < order[items[j].Status]
		}
		if items[i].Created != items[j].Created {
			return items[i].Created > items[j].Created
		}
		return items[i].Rel < items[j].Rel
	})
	cur := "\x00"
	for _, m := range items {
		if m.Status != cur {
			cur = m.Status
			fmt.Fprintf(&sb, "\n## %s\n\n| Mã | Tiêu đề | Loại | Chủ | Cập nhật | Nguồn |\n|---|---|---|---|---|---|\n", cur)
		}
		link, _ := filepath.Rel(dir, m.Path)
		id := m.ID
		if id == "" {
			id = filepath.ToSlash(link)
		}
		fmt.Fprintf(&sb, "| [%s](%s) | %s | %s | %s | %s | %s |\n", id, filepath.ToSlash(link), cell(m.Title), m.Type, cell(m.Owner), m.Updated, cell(m.Source))
	}
	return []byte(sb.String())
}

// flatRows ghi bảng phẳng sắp theo mã (ADR theo số, feature theo mã); tài
// liệu không có mã xếp cuối theo đường dẫn. features: Feature catalog với
// nguồn là brief hoặc CR; adr: kèm cột thay thế (supersedes) và ngày tạo.
func flatRows(sb *strings.Builder, name, dir string, items []docs.Meta) {
	sort.SliceStable(items, func(i, j int) bool {
		if (items[i].ID == "") != (items[j].ID == "") {
			return items[i].ID != ""
		}
		if items[i].ID != items[j].ID {
			return items[i].ID < items[j].ID
		}
		return items[i].Rel < items[j].Rel
	})
	if name == "adr" {
		sb.WriteString("\n| Mã | Tiêu đề | Trạng thái | Thay thế | Ngày |\n|---|---|---|---|---|\n")
	} else {
		sb.WriteString("\n| Mã | Tên | Trạng thái | Chủ sở hữu | Brief hoặc CR nguồn | Cập nhật |\n|---|---|---|---|---|---|\n")
	}
	for _, m := range items {
		link, _ := filepath.Rel(dir, m.Path)
		id := m.ID
		if id == "" {
			id = filepath.ToSlash(link)
		}
		if name == "adr" {
			fmt.Fprintf(sb, "| [%s](%s) | %s | %s | %s | %s |\n", id, filepath.ToSlash(link), cell(m.Title), m.Status, cell(frontmatter.GetString(m.FM, "supersedes")), m.Created)
		} else {
			fmt.Fprintf(sb, "| [%s](%s) | %s | %s | %s | %s | %s |\n", id, filepath.ToSlash(link), cell(m.Title), m.Status, cell(m.Owner), cell(m.Source), m.Updated)
		}
	}
}

// guideRows ghi mục lục User guide theo nhiệm vụ (trường task), trong nhóm
// sắp theo tiêu đề; trang chưa có task xếp cuối dưới "Chưa phân nhóm".
func guideRows(sb *strings.Builder, dir string, items []docs.Meta) {
	task := func(m docs.Meta) string { return frontmatter.GetString(m.FM, "task") }
	sort.SliceStable(items, func(i, j int) bool {
		ti, tj := task(items[i]), task(items[j])
		if (ti == "") != (tj == "") {
			return ti != ""
		}
		if ti != tj {
			return ti < tj
		}
		return items[i].Title < items[j].Title
	})
	cur := "\x00"
	for _, m := range items {
		t := task(m)
		if t == "" {
			t = "Chưa phân nhóm"
		}
		if t != cur {
			cur = t
			fmt.Fprintf(sb, "\n## %s\n\n| Trang | Trạng thái | Cập nhật |\n|---|---|---|\n", cell(t))
		}
		link, _ := filepath.Rel(dir, m.Path)
		fmt.Fprintf(sb, "| [%s](%s) | %s | %s |\n", cell(m.Title), filepath.ToSlash(link), m.Status, m.Updated)
	}
}

// noBrief là trạng thái chỉ mục của thư mục intake chưa có brief.md.
const noBrief = "chưa có brief"

// intakeRows gộp mỗi thư mục intake thành một dòng: đại diện là brief.md
// (trạng thái brief), không có thì interview.md rồi idea.md với trạng thái
// "chưa có brief".
func intakeRows(items []docs.Meta) []docs.Meta {
	rank := map[string]int{"brief": 0, "interview": 1, "idea": 2}
	best := map[string]docs.Meta{}
	var dirs []string
	for _, m := range items {
		d := filepath.Dir(m.Rel)
		cur, seen := best[d]
		if !seen {
			dirs = append(dirs, d)
		}
		r, ok := rank[m.Type]
		if !ok {
			r = len(rank)
		}
		if !seen || r < rankOf(rank, cur.Type) {
			best[d] = m
		}
	}
	rows := make([]docs.Meta, 0, len(dirs))
	for _, d := range dirs {
		m := best[d]
		if m.Type != "brief" {
			m.Status = noBrief
		}
		rows = append(rows, m)
	}
	return rows
}

func rankOf(rank map[string]int, typ string) int {
	if r, ok := rank[typ]; ok {
		return r
	}
	return len(rank)
}

// statusOrder xếp trạng thái theo statuses của loại trong types.toml; trạng
// thái lạ hoặc loại lạ xếp sau theo chữ cái.
func (a *app) statusOrder(items []docs.Meta) map[string]int {
	order := map[string]int{}
	n := 0
	for _, m := range items {
		if m.Status == noBrief {
			order[noBrief] = n
			n++
			break
		}
	}
	for _, m := range items {
		if t, ok := a.reg[m.Type]; ok {
			for _, s := range t.Statuses {
				if _, seen := order[s]; !seen {
					order[s] = n
					n++
				}
			}
		}
	}
	var rest []string
	for _, m := range items {
		if _, ok := order[m.Status]; !ok && !contains(rest, m.Status) {
			rest = append(rest, m.Status)
		}
	}
	sort.Strings(rest)
	for _, s := range rest {
		order[s] = n
		n++
	}
	return order
}

func cell(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "|", "\\|"), "\n", " ")
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
