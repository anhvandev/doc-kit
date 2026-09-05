package check

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/anhvandev/doc-kit/internal/frontmatter"
)

// specSectionOrder: thân Feature Spec có đúng các tiêu đề cấp 2 "## N." từ 1
// đến 10 theo thứ tự tăng, mỗi số một lần. has_ui: false được bỏ mục 5;
// format: crud được bỏ mục 3. Tiêu đề chữ tùy biến thể, chỉ xét số.
func specSectionOrder(c *Context) []Finding {
	var out []Finding
	for _, m := range c.typed() {
		if m.Type != "feature-spec" {
			continue
		}
		skip := map[int]bool{}
		if frontmatter.GetString(m.FM, "has_ui") == "false" {
			skip[5] = true
		}
		if frontmatter.GetString(m.FM, "format") == "crud" {
			skip[3] = true
		}
		var want []int
		for n := 1; n <= 10; n++ {
			if !skip[n] {
				want = append(want, n)
			}
		}
		got := sectionNumbers(m.Body)
		if msg := compareSections(want, got); msg != "" {
			out = append(out, finding(m, "spec-section-order", Error, msg))
		}
	}
	return out
}

// sectionNumbers lấy số của mọi tiêu đề "## N." ngoài khối mã, theo thứ tự
// xuất hiện. Khối mã mở bằng ``` hoặc ~~~ và chỉ đóng bằng cùng ký tự với độ
// dài không nhỏ hơn dấu mở.
func sectionNumbers(body []byte) []int {
	var out []int
	fence := "" // dấu mở khối mã đang trong, rỗng khi ở ngoài
	for _, l := range strings.Split(string(body), "\n") {
		t := strings.TrimSpace(l)
		if f := fenceOf(t); f != "" {
			switch {
			case fence == "":
				fence = f
			case f[0] == fence[0] && len(f) >= len(fence):
				fence = ""
			}
			continue
		}
		if fence != "" {
			continue
		}
		if sm := sectionRe.FindStringSubmatch(l); sm != nil {
			n, _ := strconv.Atoi(sm[1])
			out = append(out, n)
		}
	}
	return out
}

// fenceOf trả về chuỗi ``` hoặc ~~~ (từ 3 ký tự) mở đầu dòng, rỗng nếu không phải fence.
func fenceOf(line string) string {
	for _, ch := range []byte{'`', '~'} {
		n := 0
		for n < len(line) && line[n] == ch {
			n++
		}
		if n >= 3 {
			return line[:n]
		}
	}
	return ""
}

// compareSections trả về thông điệp lỗi khi got khác want; rỗng khi khớp.
func compareSections(want, got []int) string {
	if equalInts(want, got) {
		return ""
	}
	have := map[int]bool{}
	var dup []string
	for _, n := range got {
		if have[n] {
			dup = append(dup, strconv.Itoa(n))
		}
		have[n] = true
	}
	var missing, extra []string
	for _, n := range want {
		if !have[n] {
			missing = append(missing, strconv.Itoa(n))
		}
	}
	wantSet := map[int]bool{}
	for _, n := range want {
		wantSet[n] = true
	}
	for _, n := range got {
		if !wantSet[n] {
			extra = append(extra, strconv.Itoa(n))
		}
	}
	var parts []string
	if len(missing) > 0 {
		parts = append(parts, "thiếu mục "+strings.Join(missing, ", "))
	}
	if len(extra) > 0 {
		parts = append(parts, "mục lạ "+strings.Join(unique(extra), ", "))
	}
	if len(dup) > 0 {
		parts = append(parts, "mục lặp "+strings.Join(unique(dup), ", "))
	}
	if len(parts) == 0 {
		parts = append(parts, "sai thứ tự mục: có "+joinInts(got))
	}
	return fmt.Sprintf("tiêu đề cấp 2 phải là %s theo thứ tự; %s", joinInts(want), strings.Join(parts, "; "))
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func joinInts(a []int) string {
	s := make([]string, len(a))
	for i, n := range a {
		s[i] = strconv.Itoa(n)
	}
	return strings.Join(s, ", ")
}
