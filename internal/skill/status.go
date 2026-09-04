package skill

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/anhvandev/doc-kit/internal/target"
)

// Trạng thái cài của một skill.
const (
	StateMissing  = "chưa cài"
	StateCurrent  = "đúng phiên bản"
	StateModified = "đã sửa tay"
	StateForeign  = "không do dk"
)

// Row là một dòng của bảng `dk skill status`.
type Row struct {
	Name   string `json:"name"`
	Target string `json:"target"`
	Scope  string `json:"scope"` // dự án, toàn máy
	State  string `json:"state"`
}

// Status so mọi skill nhúng với bản đã cài trong t theo scope.
func Status(t target.Target, global bool, version string) ([]Row, error) {
	metas, err := List()
	if err != nil {
		return nil, err
	}
	dir, err := t.SkillDir(global)
	if err != nil {
		return nil, err
	}
	scope := "dự án"
	if global {
		scope = "toàn máy"
	}
	var rows []Row
	for _, m := range metas {
		src, err := Files(m.Name)
		if err != nil {
			return nil, err
		}
		row := Row{Name: m.Name, Target: t.Name(), Scope: scope}
		cur, err := readDir(filepath.Join(dir, m.Name))
		tr := readTrace(cur)
		switch {
		case errors.Is(err, os.ErrNotExist):
			row.State = StateMissing
		case err != nil:
			return nil, err
		case tr.By != installedBy:
			row.State = StateForeign
		case Hash(cur) != tr.Hash:
			row.State = StateModified
		case Hash(cur) == Hash(src) && tr.Version == version:
			row.State = StateCurrent
		default:
			row.State = fmt.Sprintf("cũ (v%s)", tr.Version)
		}
		rows = append(rows, row)
	}
	return rows, nil
}
