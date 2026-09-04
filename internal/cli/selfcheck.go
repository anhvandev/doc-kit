package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"sort"

	"github.com/spf13/cobra"

	"github.com/anhvandev/doc-kit/assets"
	"github.com/anhvandev/doc-kit/internal/skill"
	"github.com/anhvandev/doc-kit/internal/target"
	"github.com/anhvandev/doc-kit/internal/tmpl"
)

// selfCheckResult là kết quả `dk self-check`.
type selfCheckResult struct {
	Version   string   `json:"version"`
	Templates int      `json:"templates"`
	Skills    int      `json:"skills"`
	Targets   []string `json:"targets"`
	EmbedHash string   `json:"embed_hash"`
	Errors    []string `json:"errors"`
}

func newSelfCheckCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "self-check",
		Short: "Kiểm binary: phiên bản, số template, số skill, target, hash nội dung nhúng",
		Args:  exactArgs(0),
		RunE:  func(_ *cobra.Command, _ []string) error { return a.runSelfCheck() },
	}
}

// runSelfCheck render thử mọi template, đọc mọi skill và băm toàn bộ nội dung
// nhúng để so giữa hai binary; mã thoát 1 khi có lỗi.
func (a *app) runSelfCheck() error {
	res := selfCheckResult{Version: Version, Targets: target.Names, Errors: []string{}}
	for _, n := range a.reg.Names() {
		if _, err := tmpl.Raw(n); err != nil {
			res.Errors = append(res.Errors, err.Error())
		}
	}
	res.Templates = len(a.reg)
	metas, err := skill.List()
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
	}
	for _, m := range metas {
		if files, err := skill.Files(m.Name); err != nil {
			res.Errors = append(res.Errors, err.Error())
		} else if _, ok := files["references/rules.md"]; !ok {
			res.Errors = append(res.Errors, m.Name+": thiếu references/rules.md")
		}
	}
	res.Skills = len(metas)
	h, err := embedHash()
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
	}
	res.EmbedHash = h
	if a.json {
		if err := a.printJSON(res); err != nil {
			return err
		}
	} else {
		a.printf("dk %s\ntemplate: %d\nskill: %d\ntarget: %v\nembed sha256: %s\n", res.Version, res.Templates, res.Skills, res.Targets, res.EmbedHash)
		for _, e := range res.Errors {
			a.printf("lỗi: %s\n", e)
		}
	}
	if len(res.Errors) > 0 {
		return fail(codeError, "%d lỗi nội dung nhúng", len(res.Errors))
	}
	return nil
}

// embedHash băm sha256 mọi file nhúng theo đường dẫn đã sắp.
func embedHash() (string, error) {
	var paths []string
	err := fs.WalkDir(assets.FS, ".", func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			paths = append(paths, p)
		}
		return err
	})
	if err != nil {
		return "", err
	}
	sort.Strings(paths)
	h := sha256.New()
	for _, p := range paths {
		b, err := fs.ReadFile(assets.FS, p)
		if err != nil {
			return "", err
		}
		h.Write([]byte(p))
		h.Write([]byte{0})
		h.Write(b)
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
