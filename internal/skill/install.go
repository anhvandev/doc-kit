package skill

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/anhvandev/doc-kit/internal/target"
)

// Result là kết quả cài hoặc gỡ một skill.
type Result struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Action string `json:"action"` // đã cài, đã cập nhật, không đổi, đã gỡ, bỏ qua
	Note   string `json:"note,omitempty"`
}

// resolveNames trả về danh sách tên; rỗng nghĩa là mọi skill nhúng. Tên phải
// là một thành phần đường dẫn để không đụng ngoài thư mục skill.
func resolveNames(names []string) ([]string, error) {
	if len(names) > 0 {
		for _, n := range names {
			if n == "" || n == "." || n == ".." || n != filepath.Base(n) {
				return nil, fmt.Errorf("tên skill %q không hợp lệ", n)
			}
		}
		return names, nil
	}
	metas, err := List()
	if err != nil {
		return nil, err
	}
	for _, m := range metas {
		names = append(names, m.Name)
	}
	return names, nil
}

// Install cài các skill vào t. Thư mục đích không có dấu vết dk thì lỗi và
// không đụng; đã sửa tay thì lỗi trừ khi force; đúng phiên bản thì bỏ qua.
func Install(t target.Target, names []string, global, force bool, version string) ([]Result, error) {
	names, err := resolveNames(names)
	if err != nil {
		return nil, err
	}
	dir, err := t.SkillDir(global)
	if err != nil {
		return nil, err
	}
	var out []Result
	for _, name := range names {
		src, err := Files(name)
		if err != nil {
			return out, err
		}
		hash := Hash(src)
		dest := filepath.Join(dir, name)
		r := Result{Name: name, Path: dest, Action: "đã cài"}
		if cur, err := readDir(dest); err == nil {
			tr := readTrace(cur)
			switch {
			case tr.By != installedBy:
				return out, fmt.Errorf("%s: skill không do dk cài; dùng tên khác hoặc xóa tay", dest)
			case Hash(cur) != tr.Hash && !force:
				return out, fmt.Errorf("%s: skill đã sửa tay; dùng --force để ghi đè", dest)
			case Hash(cur) == hash && tr.Version == version:
				r.Action = "không đổi"
				out = append(out, r)
				continue
			}
			r.Action = "đã cập nhật"
		} else if !errors.Is(err, os.ErrNotExist) {
			return out, err
		}
		if err := write(dest, src, version, hash); err != nil {
			return out, err
		}
		out = append(out, r)
	}
	return out, nil
}

// write ghi skill vào thư mục tạm cạnh đích rồi đổi tên, để đứt giữa chừng
// không để lại thư mục thiếu dấu vết.
func write(dest string, files map[string][]byte, version, hash string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	tmp, err := os.MkdirTemp(filepath.Dir(dest), "."+filepath.Base(dest)+"-dk-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	for rel, b := range files {
		p := filepath.Join(tmp, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return err
		}
		if rel == skillFile {
			var err error
			if b, err = withTrace(b, version, hash); err != nil {
				return err
			}
		}
		if err := os.WriteFile(p, b, 0o644); err != nil {
			return err
		}
	}
	if err := os.RemoveAll(dest); err != nil {
		return err
	}
	return os.Rename(tmp, dest)
}

// Uninstall gỡ các skill có dấu vết dk; thư mục khác thì bỏ qua và báo.
func Uninstall(t target.Target, names []string, global bool) ([]Result, error) {
	names, err := resolveNames(names)
	if err != nil {
		return nil, err
	}
	dir, err := t.SkillDir(global)
	if err != nil {
		return nil, err
	}
	var out []Result
	for _, name := range names {
		dest := filepath.Join(dir, name)
		r := Result{Name: name, Path: dest, Action: "bỏ qua"}
		cur, err := readDir(dest)
		switch {
		case errors.Is(err, os.ErrNotExist):
			r.Note = "chưa cài"
		case err != nil:
			return out, err
		case readTrace(cur).By != installedBy:
			r.Note = "không do dk cài, giữ nguyên"
		default:
			if err := os.RemoveAll(dest); err != nil {
				return out, err
			}
			r.Action = "đã gỡ"
		}
		out = append(out, r)
	}
	_ = os.Remove(dir) // chỉ thành công khi thư mục skill đã rỗng
	return out, nil
}
