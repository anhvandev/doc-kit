package check

import (
	"regexp"
	"strings"
)

// Bằng chứng trong report: mã commit sau chữ "commit", liên kết Markdown tới
// file output (không phải tài liệu .md: log, ảnh, thư mục kết quả), hoặc khối
// kết quả trong ```. Chú thích HTML (gợi ý của template) bị bỏ trước khi so.
var (
	commitRe  = regexp.MustCompile(`(?i)\bcommit\b[^\n]{0,20}\b[0-9a-f]{7,40}\b`)
	mdLinkRe  = regexp.MustCompile(`\]\(([^)\s]+)\)`)
	fenceRe   = regexp.MustCompile("(?m)^\\s*```[^\\n]*\\n[^`]*\\S[^`]*```")
	commentRe = regexp.MustCompile(`(?s)<!--.*?-->`)
)

// reportEvidence: report phải có ít nhất một liên kết commit, đường dẫn file
// output hoặc khối kết quả test không rỗng; thiếu là warning.
func reportEvidence(c *Context) []Finding {
	var out []Finding
	for _, m := range c.typed() {
		if m.Type != "report" {
			continue
		}
		body := commentRe.ReplaceAll(m.Body, nil)
		if commitRe.Match(body) || fenceRe.Match(body) || hasOutputLink(body) {
			continue
		}
		out = append(out, finding(m, "report-evidence", Warning, "report không có bằng chứng: cần mã commit, liên kết file output hoặc khối kết quả test"))
	}
	return out
}

// hasOutputLink báo có liên kết Markdown tới thứ không phải tài liệu .md.
func hasOutputLink(body []byte) bool {
	for _, sm := range mdLinkRe.FindAllSubmatch(body, -1) {
		dest := strings.ToLower(string(sm[1]))
		if i := strings.IndexAny(dest, "#?"); i >= 0 {
			dest = dest[:i]
		}
		if !strings.HasSuffix(dest, ".md") {
			return true
		}
	}
	return false
}
