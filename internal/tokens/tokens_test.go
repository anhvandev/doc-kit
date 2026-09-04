package tokens

import (
	"strings"
	"testing"
)

const sample = `{
  "$dk": {"type": "design-tokens", "title": "Tokens"},
  "color": {
    "$type": "color",
    "blue": {"500": {"$value": "#2563eb"}, "700": {"$value": "#1d4ed8"}},
    "neutral": {"0": {"$value": "#ffffff"}, "900": {"$value": "#111827"}},
    "action": {"primary": {"$value": "{color.blue.500}", "$extensions": {"dk": {"theme": {"dark": "{color.blue.700}"}}}}},
    "bg": {"page": {"$value": "{color.neutral.0}", "$extensions": {"dk": {"theme": {"dark": "{color.neutral.900}"}}}}}
  },
  "space": {"$type": "dimension", "4": {"$value": 16}, "2": {"$value": {"value": 8, "unit": "px"}}},
  "radius": {"md": {"$type": "dimension", "$value": "6px"}},
  "shadow": {"card": {"$type": "shadow", "$value": "0 1px 2px {color.neutral.900}"}},
  "font": {"weight": {"bold": {"$type": "fontWeight", "$value": 700}}}
}`

func TestCSS(t *testing.T) {
	s, err := Parse([]byte(sample))
	if err != nil {
		t.Fatal(err)
	}
	css, n, err := s.CSS()
	if err != nil {
		t.Fatal(err)
	}
	out := string(css)
	if n != 11 || !strings.HasPrefix(out, Header+"\n") {
		t.Fatalf("n=%d\n%s", n, out)
	}
	for _, want := range []string{
		"  --color-blue-500: #2563eb;",
		"  --color-action-primary: #2563eb;",
		"  --space-4: 16px;",
		"  --space-2: 8px;",
		"  --radius-md: 6px;",
		"  --shadow-card: 0 1px 2px #111827;",
		"  --font-weight-bold: 700;",
		"[data-theme=\"dark\"] {\n  --color-action-primary: #1d4ed8;\n  --color-bg-page: #111827;\n}",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("thiếu %q trong:\n%s", want, out)
		}
	}
	// Thứ tự theo file: blue trước neutral trước action.
	if strings.Index(out, "--color-blue-500") > strings.Index(out, "--color-action-primary") {
		t.Errorf("thứ tự không theo file:\n%s", out)
	}
}

func TestAliasErrors(t *testing.T) {
	cyc := `{"a": {"$type": "color", "$value": "{b}"}, "b": {"$type": "color", "$value": "{c}"}, "c": {"$type": "color", "$value": "{a}"}}`
	s, err := Parse([]byte(cyc))
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := s.CSS(); err == nil || !strings.Contains(err.Error(), "alias vòng: a -> b -> c -> a") {
		t.Fatalf("alias vòng: %v", err)
	}
	s, _ = Parse([]byte(`{"a": {"$value": "{nope}"}}`))
	if _, _, err := s.CSS(); err == nil || !strings.Contains(err.Error(), "{nope} không trỏ") {
		t.Fatalf("alias lạ: %v", err)
	}
	if _, err := Parse([]byte(`{"a": "x"}`)); err == nil {
		t.Fatal("nhóm không phải object phải lỗi")
	}
}
