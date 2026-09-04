// Package skill liệt kê skill nhúng trong binary và cài, gỡ, kiểm tra
// chúng vào thư mục skill của một target. Dấu vết cài nằm trong frontmatter
// SKILL.md dưới khóa metadata: dk_installed_by, dk_version, dk_hash.
package skill

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/anhvandev/doc-kit/assets"
	"github.com/anhvandev/doc-kit/internal/frontmatter"
)

const (
	skillFile   = "SKILL.md"
	embedDir    = "skills"
	metaKey     = "metadata"
	keyBy       = "dk_installed_by"
	keyVersion  = "dk_version"
	keyHash     = "dk_hash"
	installedBy = "dk"
)

// Meta là thông tin một skill nhúng.
type Meta struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// List liệt kê skill nhúng theo tên thư mục; description đọc từ frontmatter.
func List() ([]Meta, error) {
	entries, err := fs.ReadDir(assets.FS, embedDir)
	if err != nil {
		return nil, err
	}
	var out []Meta
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		b, err := fs.ReadFile(assets.FS, path.Join(embedDir, e.Name(), skillFile))
		if err != nil {
			return nil, fmt.Errorf("skill %s thiếu %s", e.Name(), skillFile)
		}
		fm, _, _ := frontmatter.Split(b)
		out = append(out, Meta{Name: e.Name(), Description: frontmatter.GetString(fm, "description")})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// Files đọc mọi file của skill nhúng, khóa là đường dẫn tương đối dạng slash.
func Files(name string) (map[string][]byte, error) {
	root := path.Join(embedDir, name)
	if _, err := fs.Stat(assets.FS, path.Join(root, skillFile)); err != nil {
		return nil, fmt.Errorf("không có skill %q; xem `dk skill list`", name)
	}
	files := map[string][]byte{}
	err := fs.WalkDir(assets.FS, root, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := fs.ReadFile(assets.FS, p)
		if err != nil {
			return err
		}
		files[strings.TrimPrefix(p, root+"/")] = b
		return nil
	})
	return files, err
}

// readDir đọc mọi file thường trong thư mục skill đã cài, bỏ file ẩn.
func readDir(dir string) (map[string][]byte, error) {
	files := map[string][]byte{}
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(d.Name(), ".") && p != dir {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(dir, p)
		files[filepath.ToSlash(rel)] = b
		return nil
	})
	return files, err
}

// trace là dấu vết cài đọc từ frontmatter SKILL.md.
type trace struct {
	By, Version, Hash string
}

func readTrace(files map[string][]byte) trace {
	fm, _, ok := frontmatter.Split(files[skillFile])
	if !ok {
		return trace{}
	}
	m, ok := frontmatter.Get(fm, metaKey)
	if !ok || m.Kind != yaml.MappingNode {
		return trace{}
	}
	return trace{
		By:      frontmatter.GetString(m, keyBy),
		Version: frontmatter.GetString(m, keyVersion),
		Hash:    frontmatter.GetString(m, keyHash),
	}
}

// withTrace trả về SKILL.md đã chèn khối metadata dk_*.
func withTrace(skillMD []byte, version, hash string) ([]byte, error) {
	fm, body, ok := frontmatter.Split(skillMD)
	if !ok {
		return nil, fmt.Errorf("%s thiếu frontmatter", skillFile)
	}
	m, ok := frontmatter.Get(fm, metaKey)
	if !ok || m.Kind != yaml.MappingNode {
		m = &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		frontmatter.Set(fm, metaKey, m)
	}
	frontmatter.SetString(m, keyBy, installedBy)
	frontmatter.SetString(m, keyVersion, version)
	frontmatter.SetString(m, keyHash, hash)
	return frontmatter.Join(fm, body)
}
