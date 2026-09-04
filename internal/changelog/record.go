package changelog

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anhvandev/doc-kit/internal/docs"
	"github.com/anhvandev/doc-kit/internal/frontmatter"
	"github.com/anhvandev/doc-kit/internal/gitx"
)

// FileName là tên file changelog trong docs/.
const FileName = "CHANGELOG-DOCS.md"

// MergeWindow là khoảng gộp hai lần ghi cùng file cùng nguồn.
const MergeWindow = 10 * time.Minute

// GeneratedMarker là dòng đầu của file không có frontmatter do dk sinh ra
// (tokens.css); file có dòng này không cần changelog.
const GeneratedMarker = "/* generated: dk tokens css */"

// Tracks báo file (đường dẫn tương đối docsDir, dạng slash) có cần dòng
// changelog không: bỏ html/, file ẩn, chính changelog, file generated và
// file bắt đầu bằng GeneratedMarker. Mockup .html và tokens .json là tài liệu.
func Tracks(docsDir, rel string) bool {
	if rel == FileName || strings.HasPrefix(rel, "html/") || strings.HasPrefix(filepath.Base(rel), ".") {
		return false
	}
	b, err := os.ReadFile(filepath.Join(docsDir, filepath.FromSlash(rel)))
	if err != nil {
		return true
	}
	if strings.HasPrefix(string(b), GeneratedMarker) {
		return false
	}
	if fm, _, ok := frontmatter.SplitFile(rel, b); ok && frontmatter.GetString(fm, "generated") == "true" {
		return false
	}
	return true
}

// Load đọc changelog tại docsDir; thiếu file cho kết quả rỗng.
func Load(docsDir string) (*File, error) {
	b, err := os.ReadFile(filepath.Join(docsDir, FileName))
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return Parse(b)
}

// Record ghi một mục cho file rel (tương đối docsDir) rồi gộp vào changelog
// theo MergeWindow. bump=true cập nhật `updated:` trong frontmatter trước khi
// đếm để số dòng khớp git sau khi lệnh chạy xong; hook agent truyền false để
// không sửa file agent vừa ghi. Trả về mục thực sự được ghi.
func Record(root, docsDir, rel, summary, source string, now time.Time, bump bool) (Entry, error) {
	abs := filepath.Join(docsDir, filepath.FromSlash(rel))
	b, err := os.ReadFile(abs)
	if err != nil {
		return Entry{}, err
	}
	if fm, body, ok := frontmatter.SplitFile(rel, b); ok && bump {
		frontmatter.SetString(fm, "updated", now.Format("2006-01-02 15:04"))
		if b, err = frontmatter.JoinFile(rel, fm, body); err != nil {
			return Entry{}, err
		}
		if err := os.WriteFile(abs, b, 0o644); err != nil {
			return Entry{}, err
		}
	}
	st, err := gitx.NumStat(root, abs)
	if err != nil {
		return Entry{}, err
	}
	e := Entry{Time: now, Path: rel, Summary: summary, Source: source,
		Added: st.Added, Deleted: st.Deleted, NoGit: st.NoGit, New: !st.NoGit && !st.Tracked}
	if e.NoGit || e.New {
		e.Lines = docs.CountLines(b)
	}
	f, err := Load(docsDir)
	if err != nil {
		return Entry{}, err
	}
	e = f.Add(e, MergeWindow)
	if err := os.WriteFile(filepath.Join(docsDir, FileName), f.Format(), 0o644); err != nil {
		return Entry{}, err
	}
	return e, nil
}
