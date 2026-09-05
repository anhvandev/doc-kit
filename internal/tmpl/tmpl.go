// Package tmpl render template tài liệu nhúng bằng text/template.
package tmpl

import (
	"bytes"
	"fmt"
	"io/fs"
	"strings"
	"text/template"

	"github.com/anhvandev/doc-kit/assets"
)

// Data là dữ liệu đưa vào template. Mọi trường luôn có giá trị (có thể rỗng).
type Data struct {
	ID        string
	Type      string
	Slug      string
	Title     string
	Owner     string
	Created   string
	Updated   string
	Source    string
	DKVersion string
	// Format và HasUI chỉ feature-spec dùng: chọn biến thể mục 3, 4, 5, 8.
	Format string
	HasUI  bool
	// Feature, Step, Layer, External là trường --set của họ Design: mã Feature
	// Spec, mã bước, lớp Atomic Design, liên kết công cụ thiết kế ngoài.
	Feature  string
	Step     string
	Layer    string
	External string
	// Scenarios, Background, Steps chỉ họ Test dùng, chép từ Feature Spec
	// (--from): mỗi tiêu chí chấp nhận một Scenario, mục 2 thành Background,
	// mỗi mã bước có mockup ở mục 5 một dòng checklist.
	Scenarios  []Scenario
	Background []string
	Steps      []UIStep
	// Purpose, Actors, Actions, Limits chỉ release-brief dùng, chép từ Feature
	// Spec (--from): mục 1, 2, cột hành động của bảng hành vi mục 4, mục 10;
	// Steps là cột mockup của mục 5.
	Purpose []string
	Actors  []string
	Actions []string
	Limits  []string
	// Version và Collected chỉ release-notes dùng: --collect <phiên bản> gom
	// mọi Release brief ready chưa có released_in.
	Version   string
	Collected []Collected
}

// Collected là một Release brief được gom vào Release notes.
type Collected struct {
	Kind    string // feature | fix
	Feature string
	Title   string
	Link    string // đường dẫn tương đối từ file release notes
}

// Scenario là một tiêu chí chấp nhận đã tách Given / When / Then.
type Scenario struct {
	Code  string // AC1
	Title string
	Given string
	When  string
	Then  string
	Raw   string // dòng gốc khi không tách được; ba trường trên rỗng
}

// UIStep là một mã bước và mockup tương ứng ở mục 5 của Feature Spec.
type UIStep struct {
	Code   string // B1
	Mockup string // liên kết hoặc chữ ở cột Mockup
}

// Raw trả về nội dung template thô của một loại; đuôi file theo kind của loại
// (templates/<loại>.md, .html hoặc .json).
func Raw(typeName string) ([]byte, error) {
	matches, _ := fs.Glob(assets.FS, "templates/"+typeName+".*")
	if len(matches) != 1 {
		return nil, fmt.Errorf("không có template cho loại %q", typeName)
	}
	return fs.ReadFile(assets.FS, matches[0])
}

// Render điền dữ liệu vào template của loại. Thiếu trường thì lỗi.
func Render(typeName string, data Data) ([]byte, error) {
	raw, err := Raw(typeName)
	if err != nil {
		return nil, err
	}
	funcs := template.FuncMap{
		"cell": func(s string) string { return strings.ReplaceAll(s, "|", "\\|") },
		"inc":  func(i int) int { return i + 1 },
	}
	t, err := template.New(typeName).Funcs(funcs).Option("missingkey=error").Parse(string(raw))
	if err != nil {
		return nil, fmt.Errorf("template %s: %w", typeName, err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("render %s: %w", typeName, err)
	}
	return buf.Bytes(), nil
}
