package cli

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/anhvandev/doc-kit/internal/tmpl"
)

func newTemplateCmd(a *app) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Xem các loại tài liệu và template nhúng",
	}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "Liệt kê loại tài liệu",
			Args:  exactArgs(0),
			RunE:  func(_ *cobra.Command, _ []string) error { return a.templateList() },
		},
		&cobra.Command{
			Use:   "show <type>",
			Short: "In nội dung template thô của một loại",
			Args:  exactArgs(1),
			RunE:  func(_ *cobra.Command, args []string) error { return a.templateShow(args[0]) },
		},
	)
	return cmd
}

type typeInfo struct {
	Type        string   `json:"type"`
	Dir         string   `json:"dir"`
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Statuses    []string `json:"statuses"`
}

func (a *app) templateList() error {
	var infos []typeInfo
	for _, n := range a.reg.Names() {
		t := a.reg[n]
		infos = append(infos, typeInfo{Type: n, Dir: t.Dir, ID: t.IDScheme, Name: t.NamePattern, Description: t.Description, Statuses: t.Statuses})
	}
	if a.json {
		return a.printJSON(infos)
	}
	w := tabwriter.NewWriter(a.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "type\t| dir\t| id\t| mô tả")
	for _, i := range infos {
		fmt.Fprintf(w, "%s\t| %s\t| %s\t| %s\n", i.Type, i.Dir, i.ID, i.Description)
	}
	return w.Flush()
}

func (a *app) templateShow(name string) error {
	if _, err := a.reg.Get(name); err != nil {
		return fail(codeError, "%v", err)
	}
	raw, err := tmpl.Raw(name)
	if err != nil {
		return fail(codeError, "%v", err)
	}
	if a.json {
		return a.printJSON(map[string]string{"type": name, "content": string(raw)})
	}
	_, err = a.out.Write(raw)
	return err
}
