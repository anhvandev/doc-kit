// Package check chạy bộ quy tắc kiểm tra tài liệu: frontmatter, liên kết,
// mã bước, thứ tự mục spec, thứ tự duyệt CR, liên kết ngược, test, ngưỡng dòng,
// trạng thái, ADR bất biến, thuật ngữ, mockup dùng token, userflow theo spec,
// bằng chứng trong report, thuật ngữ trong tài liệu người dùng, secret trong
// Environment.
package check

import (
	"sort"

	"github.com/anhvandev/doc-kit/internal/config"
	"github.com/anhvandev/doc-kit/internal/docs"
	"github.com/anhvandev/doc-kit/internal/doctype"
)

// Mức finding.
const (
	Error   = "error"
	Warning = "warning"
)

// Finding là một phát hiện của một quy tắc.
type Finding struct {
	File  string `json:"file"` // tương đối gốc dự án
	Line  int    `json:"line,omitempty"`
	Rule  string `json:"rule"`
	Level string `json:"level"`
	Msg   string `json:"msg"`
}

// Context là dữ liệu chung cho mọi quy tắc.
type Context struct {
	Root    string // gốc dự án, tuyệt đối
	DocsDir string // tương đối Root, dạng slash
	Metas   []docs.Meta
	Reg     doctype.Registry
	Cfg     config.Check
	Jargon  []string // [release] jargon của dk.toml, cho no-jargon
	Strict  bool
}

// Rule là một quy tắc: nhận ngữ cảnh, trả về phát hiện.
type Rule struct {
	Name string
	Run  func(*Context) []Finding
}

// Rules là bộ quy tắc theo thứ tự chạy.
var Rules = []Rule{
	{"frontmatter-required", frontmatterRequired},
	{"status-valid", statusValid},
	{"link-broken", linkBroken},
	{"step-codes", stepCodes},
	{"spec-section-order", specSectionOrder},
	{"cr-approval-order", crApprovalOrder},
	{"backlink", backlink},
	{"spec-has-test", specHasTest},
	{"line-threshold", lineThreshold},
	{"adr-immutable", adrImmutable},
	{"glossary-term", glossaryTerm},
	{"mockup-tokens", mockupTokens},
	{"userflow-steps", userflowSteps},
	{"report-evidence", reportEvidence},
	{"no-jargon", noJargon},
	{"env-no-secret", envNoSecret},
}

// Run chạy mọi quy tắc, kết quả sắp theo file, dòng, quy tắc.
func Run(ctx *Context) []Finding {
	var out []Finding
	for _, r := range Rules {
		out = append(out, r.Run(ctx)...)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].File != out[j].File {
			return out[i].File < out[j].File
		}
		if out[i].Line != out[j].Line {
			return out[i].Line < out[j].Line
		}
		return out[i].Rule < out[j].Rule
	})
	return out
}

// Count đếm phát hiện theo mức.
func Count(fs []Finding) (errors, warnings int) {
	for _, f := range fs {
		if f.Level == Error {
			errors++
		} else {
			warnings++
		}
	}
	return
}

// typed trả về các tài liệu có frontmatter và loại nằm trong registry.
func (c *Context) typed() []docs.Meta {
	var out []docs.Meta
	for _, m := range c.Metas {
		if m.HasFM && !m.Generated {
			if _, ok := c.Reg[m.Type]; ok {
				out = append(out, m)
			}
		}
	}
	return out
}

func finding(m docs.Meta, rule, level, msg string) Finding {
	return Finding{File: m.Rel, Rule: rule, Level: level, Msg: msg}
}
