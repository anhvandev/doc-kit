// Package docs tạo tài liệu mới từ template theo loại.
package docs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"gopkg.in/yaml.v3"

	"github.com/anhvandev/doc-kit/internal/doctype"
	"github.com/anhvandev/doc-kit/internal/frontmatter"
	"github.com/anhvandev/doc-kit/internal/tmpl"
)

// SpecFormats là các biến thể của template feature-spec (--set format=...).
var SpecFormats = []string{"spec", "use-case", "story", "crud", "state"}

// Layers là các lớp Atomic Design của design-component (--set layer=...).
var Layers = []string{"atom", "molecule", "organism", "template"}

// placeholderRe khớp {tên} trong dir và name của types.toml, trừ {id}, {slug}, {yymmdd}, {hhmm}.
var placeholderRe = regexp.MustCompile(`\{([a-z_]+)\}`)

// ErrExists báo file đích đã tồn tại và không có Force.
var ErrExists = errors.New("file đã tồn tại; dùng --force để ghi đè")

// Options điều khiển New.
type Options struct {
	DocsDir  string // đường dẫn tuyệt đối của docs/
	PlansDir string // đường dẫn tuyệt đối của plans/, gốc cho loại có dir {plans_dir}
	In       string // thư mục plan (--in), tuyệt đối, gốc cho loại có dir {in}
	From     string // file nguồn để chép trường, tuyệt đối hoặc rỗng
	Append   string // decision-log, faq: nối một dòng thay vì tạo file
	Collect  string // release-notes: phiên bản; gom Release brief ready chưa có released_in
	Set      map[string]string
	Force    bool
	Owner    string
	IDPrefix string
	Version  string
	Now      time.Time
}

// Result là kết quả tạo file.
type Result struct {
	Path     string `json:"path"` // tuyệt đối
	ID       string `json:"id"`
	Appended bool   `json:"appended,omitempty"` // --append vào file đã có
	// Scenarios và Unparsed: số Scenario chép từ mục 9 của Feature Spec và số
	// dòng AC không tách được (họ Test với --from).
	Scenarios int `json:"scenarios,omitempty"`
	Unparsed  int `json:"unparsed,omitempty"`
	// Released: Release brief đã được ghi released_in khi --collect; tuyệt đối
	// trong docs.New, `dk new --json` in tương đối gốc dự án (như Path).
	Released []string `json:"released,omitempty"`
}

var slugRe = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// ValidSlug báo slug chỉ gồm a-z, 0-9 và dấu gạch ngang giữa các đoạn.
func ValidSlug(slug string) bool { return slugRe.MatchString(slug) }

// New tạo tài liệu loại typeName với slug, trả về đường dẫn và id.
func New(reg doctype.Registry, typeName, slug string, o Options) (Result, error) {
	t, err := reg.Get(typeName)
	if err != nil {
		return Result{}, err
	}
	if !slugRe.MatchString(slug) {
		return Result{}, fmt.Errorf("slug %q không hợp lệ: chỉ a-z, 0-9 và dấu gạch ngang", slug)
	}
	if o.Now.IsZero() {
		o.Now = time.Now()
	}

	var srcNode *source
	if o.From != "" {
		srcNode, err = readSource(o.From)
		if err != nil {
			return Result{}, err
		}
	}

	// Giá trị cho {feature}, {step}, {layer} trong dir và name: chép từ nguồn
	// --from theo bảng from trước, --set ghi đè.
	vals := map[string]string{}
	if srcNode != nil {
		for dst, src := range t.From[srcNode.typ] {
			if v := frontmatter.GetString(srcNode.fm, src); v != "" {
				vals[dst] = v
			}
		}
	}
	for k, v := range o.Set {
		vals[k] = v
	}
	if o.Collect != "" {
		if typeName != "release-notes" {
			return Result{}, fmt.Errorf("--collect chỉ dùng cho release-notes")
		}
		vals["version"] = o.Collect
	}
	if typeName == "design-component" && !slices.Contains(Layers, vals["layer"]) {
		return Result{}, fmt.Errorf("layer %q không hợp lệ; chọn một trong %s", vals["layer"], strings.Join(Layers, ", "))
	}
	if v, ok := vals["step"]; ok && !stepRe.MatchString(v) {
		return Result{}, fmt.Errorf("step %q không hợp lệ; dạng B1, B2a", v)
	}
	for _, k := range []string{"feature", "step", "layer", "version"} {
		if v, ok := vals[k]; ok && !placeholderValRe.MatchString(v) {
			return Result{}, fmt.Errorf("%s %q không hợp lệ: chỉ chữ, số, dấu chấm, gạch ngang, gạch dưới", k, v)
		}
	}
	if o.Append != "" && typeName != "decision-log" && typeName != "faq" {
		return Result{}, fmt.Errorf("--append chỉ dùng cho decision-log, faq")
	}
	// Gốc thư mục: docs/ mặc định; {plans_dir} là plans/; {in} là thư mục --in.
	base, tDir, prefix := o.DocsDir, t.Dir, o.IDPrefix
	switch {
	case strings.HasPrefix(t.Dir, "{plans_dir}"):
		base, tDir, prefix = o.PlansDir, strings.TrimPrefix(t.Dir, "{plans_dir}"), ""
	case strings.HasPrefix(t.Dir, "{in}"):
		if o.In == "" {
			return Result{}, fmt.Errorf("loại %s cần --in <thư mục plan>", typeName)
		}
		// --in phải là thư mục có sẵn trong plans_dir (hoặc chính plans_dir).
		if o.In != o.PlansDir && !inside(o.PlansDir, o.In) {
			return Result{}, fmt.Errorf("--in %s phải nằm trong %s", o.In, o.PlansDir)
		}
		if st, err := os.Stat(o.In); err != nil || !st.IsDir() {
			return Result{}, fmt.Errorf("--in %s không phải thư mục có sẵn", o.In)
		}
		base, tDir, prefix = o.In, strings.TrimPrefix(t.Dir, "{in}"), ""
	default:
		if o.In != "" {
			return Result{}, fmt.Errorf("loại %s không dùng --in", typeName)
		}
	}
	tDir, err = fill(tDir, vals)
	if err != nil {
		return Result{}, err
	}
	dir := filepath.Join(base, filepath.FromSlash(tDir))
	if t.Subdir != "" {
		// Nguồn --from nằm trong thư mục của loại (ví dụ intake/<x>/idea.md) thì
		// đặt file mới cạnh nguồn; loại có beside_source và nguồn ở thư mục loại
		// khác (ví dụ cr/CR-x.md) thì vào thư mục cùng tên nguồn (cr/CR-x/);
		// còn lại tạo subdir như thường.
		if srcNode != nil && inside(dir, filepath.Dir(o.From)) {
			dir = filepath.Dir(o.From)
		} else if srcNode != nil && t.BesideSource && inside(o.DocsDir, o.From) {
			dir = strings.TrimSuffix(o.From, filepath.Ext(o.From))
		} else {
			dir = filepath.Join(dir, expand(t.Subdir, "", slug, o.Now))
		}
	}

	id, err := makeID(t, dir, slug, prefix, o)
	if err != nil {
		return Result{}, err
	}
	name, err := fill(t.NamePattern, vals)
	if err != nil {
		return Result{}, err
	}
	path := filepath.Join(dir, expand(name, id, slug, o.Now))
	if _, err := os.Stat(path); err == nil {
		if o.Append != "" {
			return Result{Path: path, Appended: true}, appendLine(path, o.Append, o.Now)
		}
		if !o.Force {
			return Result{}, fmt.Errorf("%s: %w", path, ErrExists)
		}
	}

	data := tmpl.Data{
		ID:        id,
		Type:      typeName,
		Slug:      slug,
		Title:     titleFromSlug(slug),
		Owner:     o.Owner,
		Created:   o.Now.Format("2006-01-02"),
		Updated:   o.Now.Format("2006-01-02 15:04"),
		DKVersion: o.Version,
		Format:    "spec",
		HasUI:     o.Set["has_ui"] != "false",
		Feature:   vals["feature"],
		Step:      vals["step"],
		Layer:     vals["layer"],
		External:  vals["external"],
		Version:   vals["version"],
	}
	res := Result{ID: id}
	var collected []Meta
	if o.Collect != "" {
		data.Collected, collected, err = collectBriefs(o.DocsDir, dir, o.Collect)
		if err != nil {
			return Result{}, err
		}
	}
	if srcNode != nil && srcNode.typ == "feature-spec" && typeName == "release-brief" {
		ex := ExtractRelease(srcNode.body, filepath.Dir(o.From), dir)
		data.Purpose, data.Actors, data.Actions, data.Steps, data.Limits = ex.Purpose, ex.Actors, ex.Actions, ex.Screens, ex.Limits
	}
	if srcNode != nil && srcNode.typ == "feature-spec" && t.From["feature-spec"] != nil && strings.Contains(typeName, "test") {
		ex := ExtractSpec(srcNode.body)
		data.Scenarios, data.Background, data.Steps = ex.Scenarios, ex.Background, ex.Steps
		for _, sc := range ex.Scenarios {
			if typeName == "ui-test-checklist" {
				break // checklist chỉ dùng Steps, không in số Scenario
			}
			if sc.Raw != "" {
				res.Unparsed++
			} else {
				res.Scenarios++
			}
		}
	}
	if o.Collect != "" {
		data.Title = o.Collect
	}
	if v, ok := o.Set["title"]; ok {
		data.Title = v
	}
	if v, ok := o.Set["format"]; ok {
		data.Format = v
	}
	if typeName == "feature-spec" {
		if !slices.Contains(SpecFormats, data.Format) {
			return Result{}, fmt.Errorf("format %q không hợp lệ; chọn một trong %s", data.Format, strings.Join(SpecFormats, ", "))
		}
		if v, ok := o.Set["has_ui"]; ok && v != "true" && v != "false" {
			return Result{}, fmt.Errorf("has_ui %q không hợp lệ; chỉ true hoặc false", v)
		}
	}
	if srcNode != nil {
		data.Source = srcNode.ref
		if m := t.From[srcNode.typ]; m != nil {
			if f, ok := m["title"]; ok {
				if v := frontmatter.GetString(srcNode.fm, f); v != "" && o.Set["title"] == "" {
					data.Title = v
				}
			}
		}
	}

	out, err := tmpl.Render(typeName, data)
	if err != nil {
		return Result{}, err
	}
	fm, body, ok := frontmatter.SplitFile(path, out)
	if !ok {
		return Result{}, fmt.Errorf("template %s không có frontmatter hợp lệ", typeName)
	}
	// Các trường chuỗi tự do đi qua yaml.Node để được quote đúng (tiêu đề có dấu hai chấm).
	frontmatter.SetString(fm, "title", data.Title)
	frontmatter.SetString(fm, "owner", data.Owner)
	frontmatter.SetString(fm, "source", data.Source)
	if srcNode != nil {
		for dst, src := range t.From[srcNode.typ] {
			if v, has := frontmatter.Get(srcNode.fm, src); has {
				frontmatter.Set(fm, dst, v)
			}
		}
	}
	for k, v := range o.Set {
		if v == "true" || v == "false" {
			frontmatter.SetBool(fm, k, v == "true")
			continue
		}
		frontmatter.SetString(fm, k, v)
	}
	if typeName == "postmortem" {
		// Postmortem phải viết trong 48 giờ sau sự cố; incident_at thiếu hoặc sai
		// định dạng coi là chưa đạt để skill nhắc người.
		frontmatter.SetBool(fm, "written_within_48h", within48h(o.Set["incident_at"], o.Now))
	}
	out, err = frontmatter.JoinFile(path, fm, body)
	if err != nil {
		return Result{}, err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return Result{}, err
	}
	if err := os.WriteFile(path, out, 0o644); err != nil {
		return Result{}, err
	}
	res.Path = path
	if o.Append != "" {
		return res, appendLine(path, o.Append, o.Now)
	}
	if len(collected) > 0 {
		if res.Released, err = markReleased(collected, o.Collect, o.Now); err != nil {
			return Result{}, err
		}
	}
	return res, nil
}

// within48h báo now cách incident_at ("yyyy-mm-dd hh:mm" hoặc "yyyy-mm-dd") không quá 48 giờ.
func within48h(incidentAt string, now time.Time) bool {
	for _, layout := range []string{"2006-01-02 15:04", "2006-01-02"} {
		if t, err := time.ParseInLocation(layout, strings.TrimSpace(incidentAt), now.Location()); err == nil {
			d := now.Sub(t)
			return d >= 0 && d <= 48*time.Hour
		}
	}
	return false
}

// appendLine nối một dòng "- <ngày> | <chữ>" vào cuối file và bump updated;
// nội dung cũ giữ nguyên.
func appendLine(path, text string, now time.Time) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	fm, body, ok := frontmatter.SplitFile(path, b)
	if !ok {
		return fmt.Errorf("%s không có frontmatter hợp lệ", path)
	}
	frontmatter.SetString(fm, "updated", now.Format("2006-01-02 15:04"))
	if len(body) > 0 && body[len(body)-1] != '\n' {
		body = append(body, '\n')
	}
	body = append(body, []byte("- "+now.Format("2006-01-02")+" | "+text+"\n")...)
	out, err := frontmatter.JoinFile(path, fm, body)
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o644)
}

// source là frontmatter của file nguồn khi --from.
type source struct {
	fm   *yaml.Node
	body []byte
	typ  string
	ref  string // id nguồn, hoặc đường dẫn khi nguồn không có id
}

func readSource(path string) (*source, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("đọc file nguồn: %w", err)
	}
	fm, body, ok := frontmatter.SplitFile(path, b)
	if !ok {
		return nil, fmt.Errorf("file nguồn %s không có frontmatter", path)
	}
	n := &source{fm: fm, body: body, typ: frontmatter.GetString(fm, "type")}
	if n.typ == "" {
		return nil, fmt.Errorf("file nguồn %s thiếu trường type", path)
	}
	n.ref = frontmatter.GetString(fm, "id")
	if n.ref == "" {
		n.ref = filepath.ToSlash(filepath.Join(filepath.Base(filepath.Dir(path)), filepath.Base(path)))
	}
	return n, nil
}

// makeID sinh id theo scheme của loại; prefix là id_prefix, rỗng với loại
// nằm ngoài docs/ (phase của plan đếm cục bộ trong thư mục).
func makeID(t doctype.Type, dir, slug, prefix string, o Options) (string, error) {
	switch {
	case t.IDScheme == "none":
		return "", nil
	case strings.HasPrefix(t.IDScheme, "date:"):
		return prefix + expand(t.IDScheme[5:], "", slug, o.Now), nil
	case strings.HasPrefix(t.IDScheme, "seq:"):
		return nextSeq(t.IDScheme[4:], dir, prefix)
	}
	return "", fmt.Errorf("id scheme %q không hỗ trợ", t.IDScheme)
}

var seqRe = regexp.MustCompile(`\{n(?::0(\d+))?\}`)

// nextSeq quét dir, lấy số lớn nhất khớp mẫu rồi cộng 1.
func nextSeq(pattern, dir, prefix string) (string, error) {
	m := seqRe.FindStringSubmatchIndex(pattern)
	if m == nil {
		return "", fmt.Errorf("mẫu seq %q thiếu {n}", pattern)
	}
	width := 1
	if m[2] >= 0 {
		width, _ = strconv.Atoi(pattern[m[2]:m[3]])
	}
	pre, post := prefix+pattern[:m[0]], pattern[m[1]:]
	re := regexp.MustCompile(`^` + regexp.QuoteMeta(pre) + `(\d+)` + regexp.QuoteMeta(post))
	max := 0
	entries, err := os.ReadDir(dir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	for _, e := range entries {
		if sm := re.FindStringSubmatch(e.Name()); sm != nil {
			if n, _ := strconv.Atoi(sm[1]); n > max {
				max = n
			}
		}
	}
	return fmt.Sprintf("%s%0*d%s", pre, width, max+1, post), nil
}

// inside báo child là con thật sự của parent.
func inside(parent, child string) bool {
	rel, err := filepath.Rel(parent, child)
	return err == nil && rel != "." && !strings.HasPrefix(rel, "..")
}

var (
	stepRe = regexp.MustCompile(`^B\d+[a-z]?$`)
	// placeholderValRe chặn giá trị đi vào tên file mang dấu / hoặc .. (thoát khỏi thư mục loại).
	placeholderValRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*(\.[A-Za-z0-9_-]+)*$`)
)

// fill thay {feature}, {step}, {layer}... bằng giá trị --set hoặc chép từ nguồn;
// thiếu giá trị là lỗi vì tên file cần nó. {id}, {slug}, {yymmdd}, {hhmm} để expand lo.
func fill(pattern string, vals map[string]string) (string, error) {
	var missing []string
	out := placeholderRe.ReplaceAllStringFunc(pattern, func(ph string) string {
		k := ph[1 : len(ph)-1]
		if k == "id" || k == "slug" || k == "yymmdd" || k == "hhmm" {
			return ph
		}
		v, ok := vals[k]
		if !ok || v == "" {
			missing = append(missing, k)
			return ph
		}
		return v
	})
	if len(missing) > 0 {
		return "", fmt.Errorf("cần --set %s=... (tên file dùng %s)", strings.Join(missing, ", "), pattern)
	}
	return out, nil
}

func expand(pattern, id, slug string, now time.Time) string {
	r := strings.NewReplacer("{id}", id, "{slug}", slug, "{yymmdd}", now.Format("060102"), "{hhmm}", now.Format("1504"))
	return r.Replace(pattern)
}

func titleFromSlug(slug string) string {
	s := strings.ReplaceAll(slug, "-", " ")
	rs := []rune(s)
	if len(rs) > 0 {
		rs[0] = unicode.ToUpper(rs[0])
	}
	return string(rs)
}
