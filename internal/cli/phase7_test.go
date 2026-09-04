package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// dk tokens css sinh tokens.css từ template tokens.json; changelog add nhận
// mockup .html (bump updated trong chú thích dk) và bỏ qua tokens.css.
func TestTokensCSSAndDesignChangelog(t *testing.T) {
	dir := initProject(t)
	briefs, _ := filepath.Glob(filepath.Join(dir, "docs", "intake", "*", "brief.md"))
	if out, code := run(t, dir, "new", "design-tokens", "tokens", "--from", briefs[0], "--set", "owner=v"); code != 0 {
		t.Fatalf("new design-tokens: %s", out)
	}
	out, code := run(t, dir, "tokens", "css")
	if code != 0 || !strings.Contains(out, "docs/design/tokens/tokens.css") {
		t.Fatalf("tokens css: %s (%d)", out, code)
	}
	css, _ := os.ReadFile(filepath.Join(dir, "docs", "design", "tokens", "tokens.css"))
	if !strings.HasPrefix(string(css), "/* generated: dk tokens css */") || !strings.Contains(string(css), "--color-action-primary: #2563eb;") || !strings.Contains(string(css), `[data-theme="dark"]`) {
		t.Fatalf("tokens.css:\n%s", css)
	}
	if out, code := run(t, dir, "new", "mockup", "x", "--from", briefs[0], "--set", "feature=F-001", "--set", "step=B1", "--set", "owner=v"); code != 0 {
		t.Fatalf("new mockup: %s", out)
	}
	mock := filepath.Join(dir, "docs", "design", "mockups", "F-001-B1.html")
	if out, code := run(t, dir, "changelog", "add", "docs/design/mockups/F-001-B1.html", "--summary", "mockup"); code != 0 || !strings.Contains(out, "design/mockups/F-001-B1.html | mới,") {
		t.Fatalf("changelog add html: %s (%d)", out, code)
	}
	b, _ := os.ReadFile(mock)
	if !strings.HasPrefix(string(b), "<!-- dk:\n") || !strings.Contains(string(b), "\nupdated: ") {
		t.Fatalf("mockup sau changelog add:\n%.300s", b)
	}
	if out, code := run(t, dir, "changelog", "add", "docs/design/tokens/tokens.json", "--summary", "tokens"); code != 0 || !strings.Contains(out, "tokens/tokens.json | mới,") {
		t.Fatalf("changelog add json: %s (%d)", out, code)
	}
	out, _ = run(t, dir, "changelog", "pending", "--json")
	if strings.Contains(out, "tokens.css") || strings.Contains(out, "design/") {
		t.Fatalf("pending phải bỏ tokens.css và file design đã ghi: %s", out)
	}
	for _, f := range []string{"docs/design/mockups/F-001-B1.html", "docs/design/tokens/tokens.json"} {
		if out, code := run(t, dir, "check", f); code != 0 {
			t.Fatalf("check %s: %s (%d)", f, out, code)
		}
	}
	// Lệnh render --all không đụng .html, .json.
	out, _ = run(t, dir, "render", "--all")
	if strings.Contains(out, "mockups") || strings.Contains(out, "tokens.json") {
		t.Fatalf("render --all render file không Markdown: %s", out)
	}
}
