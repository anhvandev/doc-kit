package docs

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/anhvandev/doc-kit/internal/tmpl"
)

// Họ Test chép từ Feature Spec: mục 9 (tiêu chí chấp nhận) thành Scenario,
// mục 3 (tác nhân và điều kiện tiên quyết) thành Background, mục 6 (giao
// diện) thành checklist theo mã bước. Chỉ xét chữ trong tiêu đề số; AC lệch
// khung được giữ nguyên dòng gốc ở Raw thay vì bỏ qua im lặng.

var (
	sectionHeadRe = regexp.MustCompile(`^## (\d+)\.`)
	acBulletRe    = regexp.MustCompile(`^\s*[-*]\s*\**(AC\d+)\**[.:]?\**\s*(.*)$`)
	acAnyRe       = regexp.MustCompile(`\b(AC\d+)\b`)
	acScenarioRe  = regexp.MustCompile(`^\s*(Scenario|Scenario Outline):\s*(AC\d+)?\s*(.*)$`)
	// labelRe: nhãn đứng trước dấu hai chấm ở mục 3, bỏ khi chép sang Background.
	labelRe       = regexp.MustCompile(`(?i)^(tác nhân|điều kiện tiên quyết|quyền|dữ liệu)\s*:\s*`)
	gherkinStepRe = regexp.MustCompile(`^\s*(Given|When|Then|And|But)\s+(.*)$`)
	gwtRe         = regexp.MustCompile(`\**(Given|When|Then)\**`)
	bulletRe      = regexp.MustCompile(`^\s*[-*]\s+(.*)$`)
	tableRowRe    = regexp.MustCompile(`^\s*\|(.*)\|\s*$`)
	htmlCommentRe = regexp.MustCompile(`(?s)<!--.*?-->`)
)

// SpecExtract là phần Feature Spec mà họ Test cần.
type SpecExtract struct {
	Scenarios  []tmpl.Scenario
	Background []string
	Steps      []tmpl.UIStep
}

// ExtractSpec đọc thân Feature Spec (Markdown, không frontmatter); chú thích
// HTML (gợi ý của template) bị bỏ trước.
func ExtractSpec(body []byte) SpecExtract {
	sec := sections(string(htmlCommentRe.ReplaceAll(body, nil)))
	return SpecExtract{
		Scenarios:  scenarios(sec[9]),
		Background: bullets(sec[3]),
		Steps:      uiSteps(sec[6]),
	}
}

// sections cắt thân theo tiêu đề "## N." thành map số mục sang các dòng bên
// trong; tiêu đề nằm trong khối mã không tính.
func sections(body string) map[int][]string {
	out := map[int][]string{}
	cur := 0
	inFence := false
	for _, l := range strings.Split(body, "\n") {
		if t := strings.TrimSpace(l); strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~") {
			inFence = !inFence
		}
		if inFence {
			if cur > 0 {
				out[cur] = append(out[cur], l)
			}
			continue
		}
		if sm := sectionHeadRe.FindStringSubmatch(l); sm != nil {
			cur = atoi(sm[1])
			continue
		}
		if strings.HasPrefix(l, "## ") || strings.HasPrefix(l, "# ") {
			cur = 0
			continue
		}
		if cur > 0 {
			out[cur] = append(out[cur], l)
		}
	}
	return out
}

func atoi(s string) int {
	n := 0
	for _, r := range s {
		n = n*10 + int(r-'0')
	}
	return n
}

// scenarios tách mục 9: bullet "- AC1. **Given** a **When** b **Then** c" hoặc
// khối gherkin "Scenario: AC1 ..." với dòng Given / When / Then (And nối vào
// bước trước). Mọi dòng khác có nhắc mã AC (danh sách đánh số, bảng, Scenario
// Outline) và AC thiếu một trong ba phần giữ nguyên chữ ở Raw, không bỏ qua.
func scenarios(lines []string) []tmpl.Scenario {
	var out []tmpl.Scenario
	seen := map[string]bool{}
	inFence := false
	var cur *tmpl.Scenario
	last := ""
	for _, l := range lines {
		t := strings.TrimSpace(l)
		if strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~") {
			inFence = !inFence
			cur = nil
			continue
		}
		if inFence {
			if sm := acScenarioRe.FindStringSubmatch(l); sm != nil {
				cur = nil
				if sm[1] == "Scenario Outline" {
					out = append(out, tmpl.Scenario{Code: sm[2], Raw: "Scenario Outline: " + strings.TrimSpace(sm[3]) + " (chép tay cùng bảng Examples)"})
					continue
				}
				out = append(out, tmpl.Scenario{Code: sm[2], Title: strings.TrimSpace(sm[3])})
				cur = &out[len(out)-1]
				last = ""
				continue
			}
			if sm := gherkinStepRe.FindStringSubmatch(l); sm != nil && cur != nil {
				kw, text := sm[1], strings.TrimSpace(sm[2])
				if kw == "And" || kw == "But" {
					kw = last
				}
				appendStep(cur, kw, text)
				last = kw
			}
			continue
		}
		sm := acBulletRe.FindStringSubmatch(l)
		if sm == nil {
			// Dòng nhắc mã AC chưa gặp nhưng không phải bullet: danh sách đánh số, bảng.
			if am := acAnyRe.FindStringSubmatch(l); am != nil && !seen[am[1]] {
				seen[am[1]] = true
				out = append(out, tmpl.Scenario{Code: am[1], Raw: strings.Trim(t, "|- ")})
			}
			continue
		}
		seen[sm[1]] = true
		sc := tmpl.Scenario{Code: sm[1]}
		rest := sm[2]
		parts := gwtRe.FindAllStringSubmatchIndex(rest, -1)
		if len(parts) < 3 {
			sc.Raw = strings.TrimSpace(rest)
			out = append(out, sc)
			continue
		}
		sc.Title = strings.TrimSpace(rest[:parts[0][0]])
		for i, p := range parts {
			end := len(rest)
			if i+1 < len(parts) {
				end = parts[i+1][0]
			}
			appendStep(&sc, rest[p[2]:p[3]], strings.TrimSpace(rest[p[1]:end]))
		}
		if sc.Given == "" || sc.When == "" || sc.Then == "" {
			sc = tmpl.Scenario{Code: sm[1], Raw: strings.TrimSpace(rest)}
		}
		out = append(out, sc)
	}
	return out
}

func appendStep(sc *tmpl.Scenario, kw, text string) {
	var dst *string
	switch kw {
	case "Given":
		dst = &sc.Given
	case "When":
		dst = &sc.When
	case "Then":
		dst = &sc.Then
	default:
		return
	}
	text = strings.TrimSpace(strings.TrimSuffix(strings.TrimSuffix(text, ","), ";"))
	if *dst == "" {
		*dst = text
	} else if text != "" {
		*dst += " và " + text
	}
}

// bullets lấy chữ sau dấu gạch đầu dòng, bỏ nhãn đã biết ("Tác nhân:",
// "Điều kiện tiên quyết:"...) đứng trước, bỏ dòng rỗng và dòng chỉ có nhãn.
func bullets(lines []string) []string {
	var out []string
	for _, l := range lines {
		sm := bulletRe.FindStringSubmatch(l)
		if sm == nil {
			continue
		}
		t := strings.TrimSpace(labelRe.ReplaceAllString(strings.TrimSpace(sm[1]), ""))
		if t == "" || strings.HasSuffix(t, ":") {
			continue
		}
		out = append(out, t)
	}
	return out
}

// uiSteps đọc bảng mục 6: mỗi dòng có cột đầu là mã bước thành một mục checklist.
func uiSteps(lines []string) []tmpl.UIStep {
	var out []tmpl.UIStep
	for _, l := range lines {
		sm := tableRowRe.FindStringSubmatch(l)
		if sm == nil {
			continue
		}
		cells := strings.Split(sm[1], "|")
		if len(cells) < 2 {
			continue
		}
		code := strings.TrimSpace(cells[0])
		if !stepRe.MatchString(code) {
			continue
		}
		out = append(out, tmpl.UIStep{Code: code, Mockup: strings.TrimSpace(cells[1])})
	}
	return out
}

// ReleaseExtract là phần Feature Spec mà Release brief chép: mục 2 (mục đích),
// mục 3 (tác nhân), cột hành động của bảng hành vi mục 5, cột mockup mục 6,
// mục 11 (ngoài phạm vi). Skill viết lại bằng ngôn ngữ người dùng sau đó.
type ReleaseExtract struct {
	Purpose []string
	Actors  []string
	Actions []string
	Screens []tmpl.UIStep
	Limits  []string
}

// ExtractRelease đọc thân Feature Spec; chú thích HTML bị bỏ trước. Liên kết
// tương đối trong cột mockup được đổi gốc từ specDir sang briefDir.
func ExtractRelease(body []byte, specDir, briefDir string) ReleaseExtract {
	sec := sections(string(htmlCommentRe.ReplaceAll(body, nil)))
	ex := ReleaseExtract{Actors: bullets(sec[3]), Limits: bullets(sec[11])}
	for _, l := range sec[2] {
		if t := strings.TrimSpace(l); t != "" {
			ex.Purpose = append(ex.Purpose, strings.TrimSpace(strings.TrimPrefix(t, "- ")))
		}
	}
	for _, l := range sec[5] {
		sm := tableRowRe.FindStringSubmatch(l)
		if sm == nil {
			continue
		}
		cells := strings.Split(sm[1], "|")
		if len(cells) < 2 || !stepRe.MatchString(strings.TrimSpace(cells[0])) {
			continue
		}
		if a := strings.TrimSpace(cells[1]); a != "" {
			ex.Actions = append(ex.Actions, a)
		}
	}
	for _, s := range uiSteps(sec[6]) {
		// Ô không có liên kết ("chưa có, xem họ Design") không phải ảnh, bỏ.
		if !strings.Contains(s.Mockup, "](") {
			continue
		}
		s.Mockup = rebaseLinks(s.Mockup, specDir, briefDir)
		ex.Screens = append(ex.Screens, s)
	}
	return ex
}

var mdLinkRe = regexp.MustCompile(`\]\(([^)\s#]+)(#[^)]*)?\)`)

// rebaseLinks đổi gốc liên kết Markdown tương đối trong text từ thư mục
// fromDir sang toDir; liên kết tuyệt đối hoặc có scheme (`https:`, `mailto:`) giữ nguyên.
func rebaseLinks(text, fromDir, toDir string) string {
	return mdLinkRe.ReplaceAllStringFunc(text, func(link string) string {
		sm := mdLinkRe.FindStringSubmatch(link)
		dest := sm[1]
		if strings.Contains(dest, ":") || strings.HasPrefix(dest, "/") {
			return link
		}
		rel, err := filepath.Rel(toDir, filepath.Join(fromDir, filepath.FromSlash(dest)))
		if err != nil {
			return link
		}
		return "](" + filepath.ToSlash(rel) + sm[2] + ")"
	})
}
