// Package config đọc và ghi dk.toml.
package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// FileName là tên file cấu hình tại gốc dự án.
const FileName = "dk.toml"

// Config là nội dung dk.toml.
type Config struct {
	ProjectName  string  `toml:"project_name"`
	IDPrefix     string  `toml:"id_prefix"`
	DefaultOwner string  `toml:"default_owner"`
	DocsDir      string  `toml:"docs_dir"`
	PlansDir     string  `toml:"plans_dir"`
	Language     string  `toml:"language"`
	Check        Check   `toml:"check"`
	Test         Test    `toml:"test"`
	Release      Release `toml:"release"`
}

// Release là cấu hình cho họ Người dùng cuối.
type Release struct {
	// Jargon: từ kỹ thuật không được xuất hiện trong Release brief và User guide
	// (dk check no-jargon, warning); rỗng lấy DefaultJargon. Tên sản phẩm có chứa
	// từ này thì bỏ từ đó khỏi danh sách.
	Jargon []string `toml:"jargon"`
}

// DefaultJargon là danh sách thuật ngữ mặc định của no-jargon.
var DefaultJargon = []string{"API", "endpoint", "database", "migration", "backend", "frontend", "JSON"}

// Test là cấu hình cho họ Test.
type Test struct {
	// BddCmd: lệnh công cụ BDD đã chốt trong Testing strategy (ví dụ "cucumber",
	// "behave"); skill doc-test chạy `<lệnh> --dry-run` cho file .feature sinh ra.
	// Rỗng: chưa chốt, skill báo "chưa kiểm chạy được".
	BddCmd string `toml:"bdd_cmd"`
}

// Check là ngưỡng của `dk check`.
type Check struct {
	WarnLines int `toml:"warn_lines"` // trên ngưỡng này là warning
	MaxLines  int `toml:"max_lines"`  // trên ngưỡng này là lỗi
}

// Default trả về cấu hình mặc định cho dự án có tên đã cho.
func Default(projectName string) Config {
	return Config{
		ProjectName: projectName,
		DocsDir:     "docs",
		PlansDir:    "plans",
		Language:    "vi",
		Check:       Check{WarnLines: 500, MaxLines: 800},
		Release:     Release{Jargon: DefaultJargon},
	}
}

// Load đọc dk.toml; trường thiếu lấy mặc định.
func Load(path string) (Config, error) {
	cfg := Default("")
	b, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("đọc %s: %w", path, err)
	}
	if err := toml.Unmarshal(b, &cfg); err != nil {
		return cfg, fmt.Errorf("phân tích %s: %w", path, err)
	}
	if cfg.DocsDir == "" {
		cfg.DocsDir = "docs"
	}
	if cfg.PlansDir == "" {
		cfg.PlansDir = "plans"
	}
	if cfg.Check.WarnLines <= 0 {
		cfg.Check.WarnLines = 500
	}
	if cfg.Check.MaxLines <= 0 {
		cfg.Check.MaxLines = 800
	}
	if len(cfg.Release.Jargon) == 0 {
		cfg.Release.Jargon = DefaultJargon
	}
	return cfg, nil
}

// Write ghi cấu hình ra file.
func Write(path string, cfg Config) error {
	b, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
