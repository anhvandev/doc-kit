package docs

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/anhvandev/doc-kit/internal/frontmatter"
	"github.com/anhvandev/doc-kit/internal/tmpl"
)

// briefsDir là thư mục Release brief tương đối docs/ (khớp dir của release-brief).
const briefsDir = "release/briefs"

// collectBriefs gom mọi Release brief `status: ready` chưa có `released_in`
// trong docsDir, cộng brief đã có `released_in` đúng version (để `--force`
// sinh lại notes cùng phiên bản không làm mất dòng cũ); Link tương đối từ
// notesDir. Trả về riêng các brief cần ghi released_in. Không có brief nào là lỗi.
func collectBriefs(docsDir, notesDir, version string) ([]tmpl.Collected, []Meta, error) {
	metas, err := Scan(docsDir, briefsDir)
	if err != nil {
		return nil, nil, err
	}
	var out []tmpl.Collected
	var picked []Meta
	for _, m := range metas {
		released := frontmatter.GetString(m.FM, "released_in")
		if !m.HasFM || m.Type != "release-brief" || m.Status != "ready" || (released != "" && released != version) {
			continue
		}
		link, _ := filepath.Rel(notesDir, m.Path)
		kind := frontmatter.GetString(m.FM, "kind")
		if kind != "fix" {
			kind = "feature"
		}
		out = append(out, tmpl.Collected{Kind: kind, Feature: frontmatter.GetString(m.FM, "feature"), Title: m.Title, Link: filepath.ToSlash(link)})
		if released == "" {
			picked = append(picked, m)
		}
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("không có Release brief nào status ready chưa có released_in trong %s", filepath.Join(docsDir, briefsDir))
	}
	return out, picked, nil
}

// markReleased ghi released_in vào từng brief và bump updated. Nội dung mới
// dựng xong cho mọi brief rồi mới ghi, nên lỗi ghi giữa chừng chỉ còn là lỗi
// I/O; khi đó thông điệp nêu các file đã sửa để người sửa tay. Trả về đường
// dẫn tuyệt đối các file đã sửa để lệnh gọi ghi changelog.
func markReleased(briefs []Meta, version string, now time.Time) ([]string, error) {
	outs := make([][]byte, len(briefs))
	for i, m := range briefs {
		frontmatter.SetString(m.FM, "released_in", version)
		frontmatter.SetString(m.FM, "updated", now.Format("2006-01-02 15:04"))
		out, err := frontmatter.JoinFile(m.Path, m.FM, m.Body)
		if err != nil {
			return nil, err
		}
		outs[i] = out
	}
	var paths []string
	for i, m := range briefs {
		if err := os.WriteFile(m.Path, outs[i], 0o644); err != nil {
			return paths, fmt.Errorf("ghi released_in vào %s: %w (đã sửa: %s)", m.Path, err, strings.Join(paths, ", "))
		}
		paths = append(paths, m.Path)
	}
	return paths, nil
}

var nonSlugRe = regexp.MustCompile(`[^a-z0-9]+`)

// VersionSlug đổi phiên bản thành slug hợp lệ: "v1.0.0" thành "v1-0-0".
func VersionSlug(version string) string {
	return strings.Trim(nonSlugRe.ReplaceAllString(strings.ToLower(version), "-"), "-")
}
