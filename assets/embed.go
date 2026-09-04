// Package assets nhúng template tài liệu, bảng loại tài liệu, khung HTML,
// bộ skill và script pre-commit vào binary.
package assets

import "embed"

//go:embed templates/* types.toml agent-context.md html/page.html html/style.css html/mermaid.min.js html/MERMAID-LICENSE.txt skills hooks/pre-commit.sh
var FS embed.FS
