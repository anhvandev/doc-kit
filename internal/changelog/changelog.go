// Package changelog đọc, thêm và ghi CHANGELOG-DOCS.md: mới nhất ở trên,
// nhóm theo ngày, mỗi dòng một thay đổi.
package changelog

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Title là dòng đầu file.
const Title = "# Changelog tài liệu"

// NoSummary là tóm tắt mặc định khi người gọi không đưa.
const NoSummary = "chưa tóm tắt"

// Entry là một dòng changelog.
type Entry struct {
	Time    time.Time
	Path    string // tương đối docs/
	Added   int
	Deleted int
	New     bool // file chưa có trong HEAD
	NoGit   bool // đếm toàn file vì không có git
	Lines   int  // số dòng toàn file khi New hoặc NoGit
	Summary string
	Source  string
}

// Day là nhóm dòng cùng ngày.
type Day struct {
	Date    string // yyyy-mm-dd
	Entries []Entry
}

// File là toàn bộ changelog.
type File struct {
	Days []Day
}

const (
	dateLayout = "2006-01-02"
	timeLayout = "15:04"
)

var (
	dayRe   = regexp.MustCompile(`^## (\d{4}-\d{2}-\d{2})\s*$`)
	entryRe = regexp.MustCompile(`^- (\d{2}:\d{2}) \| (.+?) \| (.+?) \| (.*?) \| (.*)$`)
	diffRe  = regexp.MustCompile(`^\+(\d+) [−-](\d+)$`)
	newRe   = regexp.MustCompile(`^(mới|không git), (\d+) dòng$`)
)

// Parse đọc changelog. File rỗng cho kết quả rỗng.
func Parse(b []byte) (*File, error) {
	f := &File{}
	sc := bufio.NewScanner(bytes.NewReader(b))
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	var cur *Day
	for n := 1; sc.Scan(); n++ {
		ln := strings.TrimRight(sc.Text(), "\r")
		if m := dayRe.FindStringSubmatch(ln); m != nil {
			f.Days = append(f.Days, Day{Date: m[1]})
			cur = &f.Days[len(f.Days)-1]
			continue
		}
		if !strings.HasPrefix(ln, "- ") {
			continue
		}
		if cur == nil {
			return nil, fmt.Errorf("dòng %d: mục nằm ngoài nhóm ngày", n)
		}
		e, err := parseEntry(cur.Date, ln)
		if err != nil {
			return nil, fmt.Errorf("dòng %d: %w", n, err)
		}
		cur.Entries = append(cur.Entries, e)
	}
	return f, sc.Err()
}

func parseEntry(date, ln string) (Entry, error) {
	m := entryRe.FindStringSubmatch(ln)
	if m == nil {
		return Entry{}, fmt.Errorf("sai định dạng: %q", ln)
	}
	ts, err := time.ParseInLocation(dateLayout+" "+timeLayout, date+" "+m[1], time.Local)
	if err != nil {
		return Entry{}, err
	}
	e := Entry{Time: ts, Path: m[2], Summary: m[4], Source: m[5]}
	if e.Source == "-" {
		e.Source = ""
	}
	switch {
	case diffRe.MatchString(m[3]):
		d := diffRe.FindStringSubmatch(m[3])
		e.Added, _ = strconv.Atoi(d[1])
		e.Deleted, _ = strconv.Atoi(d[2])
	case newRe.MatchString(m[3]):
		d := newRe.FindStringSubmatch(m[3])
		e.New = d[1] == "mới"
		e.NoGit = d[1] == "không git"
		e.Lines, _ = strconv.Atoi(d[2])
	default:
		return Entry{}, fmt.Errorf("cột số dòng lạ: %q", m[3])
	}
	return e, nil
}

// Format ghi changelog ra byte, mới nhất ở trên.
func (f *File) Format() []byte {
	var b strings.Builder
	b.WriteString(Title + "\n")
	for _, d := range f.Days {
		b.WriteString("\n## " + d.Date + "\n\n")
		for _, e := range d.Entries {
			b.WriteString(e.String() + "\n")
		}
	}
	return []byte(b.String())
}

// String trả về dòng changelog của mục.
func (e Entry) String() string {
	var lines string
	switch {
	case e.NoGit:
		lines = fmt.Sprintf("không git, %d dòng", e.Lines)
	case e.New:
		lines = fmt.Sprintf("mới, %d dòng", e.Lines)
	default:
		lines = fmt.Sprintf("+%d −%d", e.Added, e.Deleted)
	}
	sum := e.Summary
	if strings.TrimSpace(sum) == "" {
		sum = NoSummary
	}
	src := e.Source
	if src == "" {
		src = "-"
	}
	return fmt.Sprintf("- %s | %s | %s | %s | %s", e.Time.Format(timeLayout), e.Path, lines, sum, src)
}

// Since liệt kê đường dẫn có mục changelog từ mốc t trở đi.
func (f *File) Since(t time.Time) map[string]bool {
	t = t.Truncate(time.Minute)
	seen := map[string]bool{}
	for _, d := range f.Days {
		for _, e := range d.Entries {
			if !e.Time.Before(t) {
				seen[e.Path] = true
			}
		}
	}
	return seen
}
