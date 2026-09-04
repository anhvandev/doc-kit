// Package gitx gọi git qua os/exec cho những việc CLI cần: numstat, file đổi,
// mốc commit HEAD. Máy không có git hoặc thư mục ngoài repo thì trả NoGit.
package gitx

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Stat là số dòng thêm, xóa của một file so với HEAD.
type Stat struct {
	Added, Deleted int
	Tracked        bool // file có trong HEAD
	NoGit          bool // không có git hoặc ngoài repo
}

func run(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), strings.TrimSpace(errb.String()))
	}
	return strings.TrimRight(out.String(), "\r\n\x00"), nil
}

// IsRepo báo dir nằm trong một work tree git và có git trên máy.
func IsRepo(dir string) bool {
	if _, err := exec.LookPath("git"); err != nil {
		return false
	}
	out, err := run(dir, "rev-parse", "--is-inside-work-tree")
	return err == nil && out == "true"
}

// Root trả về gốc work tree.
func Root(dir string) (string, error) {
	if !IsRepo(dir) {
		return "", errors.New("không phải repo git")
	}
	return run(dir, "rev-parse", "--show-toplevel")
}

func hasHead(dir string) bool {
	_, err := run(dir, "rev-parse", "--verify", "-q", "HEAD")
	return err == nil
}

// NumStat đếm dòng thêm, xóa của path (tuyệt đối hoặc tương đối dir) so với
// HEAD, gồm cả phần đã stage. File chưa có trong HEAD thì Tracked=false.
func NumStat(dir, path string) (Stat, error) {
	if !IsRepo(dir) {
		return Stat{NoGit: true}, nil
	}
	rel := path
	if filepath.IsAbs(path) {
		root, err := Root(dir)
		if err != nil {
			return Stat{}, err
		}
		if rel, err = filepath.Rel(root, path); err != nil {
			return Stat{}, err
		}
		dir = root
	}
	if !hasHead(dir) {
		return Stat{}, nil
	}
	if _, err := run(dir, "cat-file", "-e", "HEAD:"+filepath.ToSlash(rel)); err != nil {
		return Stat{}, nil
	}
	out, err := run(dir, "diff", "--numstat", "HEAD", "--", filepath.ToSlash(rel))
	if err != nil {
		return Stat{}, err
	}
	st := Stat{Tracked: true}
	if out == "" {
		return st, nil
	}
	f := strings.Fields(out)
	if len(f) < 2 {
		return st, fmt.Errorf("numstat lạ: %q", out)
	}
	st.Added, _ = strconv.Atoi(f[0])
	st.Deleted, _ = strconv.Atoi(f[1])
	return st, nil
}

// ChangedDocs liệt kê file trong docsDir đã đổi, stage hoặc chưa theo dõi,
// đường dẫn tương đối gốc repo.
func ChangedDocs(root, docsDir string) ([]string, error) {
	// -z và quotePath=false để đường dẫn có dấu hoặc khoảng trắng về nguyên dạng.
	out, err := run(root, "-c", "core.quotePath=false", "status", "--porcelain", "-z", "--untracked-files=all", "--", docsDir)
	if err != nil {
		return nil, err
	}
	var files []string
	fields := strings.Split(out, "\x00")
	for i := 0; i < len(fields); i++ {
		ln := fields[i]
		if len(ln) < 4 {
			continue
		}
		if ln[0] == 'R' || ln[0] == 'C' {
			i++ // trường kế tiếp là đường dẫn cũ
		}
		if ln[0] == 'D' || ln[1] == 'D' {
			continue
		}
		files = append(files, ln[3:])
	}
	return files, nil
}

// HeadTime trả về thời điểm commit HEAD; zero nếu chưa có commit.
func HeadTime(root string) (time.Time, error) {
	if !hasHead(root) {
		return time.Time{}, nil
	}
	out, err := run(root, "log", "-1", "--format=%cI")
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339, out)
}

// HooksDir trả về thư mục hooks của repo chứa dir (đúng cả với worktree).
func HooksDir(dir string) (string, error) {
	out, err := run(dir, "rev-parse", "--git-path", "hooks")
	if err != nil {
		return "", err
	}
	if !filepath.IsAbs(out) {
		out = filepath.Join(dir, out)
	}
	return out, nil
}

// HeadFile trả về nội dung của path (tuyệt đối hoặc tương đối dir) tại HEAD;
// ok=false khi không có git, chưa có commit hoặc file chưa trong HEAD (mọi
// thất bại của `git show` đều coi là chưa có, để mỗi lần gọi chỉ tốn một
// tiến trình git ngoài việc tìm gốc repo).
func HeadFile(dir, path string) (content []byte, ok bool, err error) {
	if _, err := exec.LookPath("git"); err != nil {
		return nil, false, nil
	}
	rel := path
	if filepath.IsAbs(path) {
		root, err := Root(dir)
		if err != nil {
			return nil, false, nil
		}
		if rel, err = filepath.Rel(root, path); err != nil {
			return nil, false, err
		}
		dir = root
	}
	out, err := run(dir, "show", "HEAD:"+filepath.ToSlash(rel))
	if err != nil {
		return nil, false, nil
	}
	return []byte(out + "\n"), true, nil
}
