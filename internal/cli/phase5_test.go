package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Feature catalog và chỉ mục ADR là bảng phẳng theo mã, mỗi tài liệu một dòng.
func TestIndexFeaturesADR(t *testing.T) {
	dir := initProject(t)
	briefs, _ := filepath.Glob(filepath.Join(dir, "docs", "intake", "*", "brief.md"))
	if _, code := run(t, dir, "new", "feature-spec", "xuat-file", "--from", briefs[0], "--set", "owner=v", "--set", "status=review"); code != 0 {
		t.Fatal("new feature-spec 2")
	}
	for _, slug := range []string{"postgres", "hang-doi"} {
		if _, code := run(t, dir, "new", "adr", slug, "--set", "owner=v"); code != 0 {
			t.Fatal("new adr", slug)
		}
	}
	if _, code := run(t, dir, "new", "adr", "redis", "--set", "owner=v", "--set", "supersedes=ADR-0002"); code != 0 {
		t.Fatal("new adr 3")
	}
	if _, code := run(t, dir, "index", "features"); code != 0 {
		t.Fatal("index features")
	}
	if _, code := run(t, dir, "index", "adr"); code != 0 {
		t.Fatal("index adr")
	}
	feat, _ := os.ReadFile(filepath.Join(dir, "docs", "features", "README.md"))
	rows := tableRows(string(feat))
	specs, _ := filepath.Glob(filepath.Join(dir, "docs", "features", "F-*.md"))
	if len(rows) != len(specs) || len(rows) != 2 || !strings.HasPrefix(rows[0], "| [F-001](F-001-bo-loc.md) | Loc don | draft | v | ") || !strings.HasPrefix(rows[1], "| [F-002](F-002-xuat-file.md) | Loc don | review | v | ") {
		t.Fatalf("Feature catalog sai (%d file):\n%s", len(specs), feat)
	}
	if !strings.Contains(string(feat), "| Mã | Tên | Trạng thái | Chủ sở hữu | Brief hoặc CR nguồn | Cập nhật |") {
		t.Fatalf("cột Feature catalog sai:\n%s", feat)
	}
	adr, _ := os.ReadFile(filepath.Join(dir, "docs", "adr", "README.md"))
	rows = tableRows(string(adr))
	if len(rows) != 3 || !strings.Contains(rows[0], "[ADR-0001](ADR-0001-postgres.md) | Postgres | proposed |  |") ||
		!strings.Contains(rows[1], "[ADR-0002]") || !strings.Contains(rows[2], "[ADR-0003](ADR-0003-redis.md) | Redis | proposed | ADR-0002 |") {
		t.Fatalf("chỉ mục ADR sai:\n%s", adr)
	}
	if !strings.Contains(string(adr), "| Mã | Tiêu đề | Trạng thái | Thay thế | Ngày |") {
		t.Fatalf("cột chỉ mục ADR sai:\n%s", adr)
	}
}

func tableRows(s string) []string {
	var rows []string
	for _, l := range strings.Split(s, "\n") {
		if strings.HasPrefix(l, "| [") {
			rows = append(rows, l)
		}
	}
	return rows
}
