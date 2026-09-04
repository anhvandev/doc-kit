// Package agentctx ghi và kiểm khối ngữ cảnh agent (assets/agent-context.md)
// trong file ngữ cảnh của dự án (CLAUDE.md, AGENTS.md). Khối nằm giữa hai dấu
// mốc HTML comment mang phiên bản dk và hash nội dung, nên chạy lại thay đúng
// khối cũ và `dk doctor` nhận ra khối thiếu, cũ hay bị sửa tay.
package agentctx

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Files là các file ngữ cảnh agent được ghi, tương đối gốc dự án.
var Files = []string{"CLAUDE.md", "AGENTS.md"}

const (
	startPrefix = "<!-- dk:agent-context start"
	endMarker   = "<!-- dk:agent-context end -->"
)

// Trạng thái của khối trong một file.
const (
	StateOK          = "ok"
	StateMissingFile = "missing-file"
	StateMissing     = "missing-block"
	StateOutdated    = "outdated" // version hoặc hash khác nội dung nhúng hiện tại
	StateBroken      = "broken"   // có mốc mở nhưng thiếu mốc đóng, hoặc nhiều hơn một khối
)

// Kết quả ghi.
const (
	WriteCreated   = "created"
	WriteUpdated   = "updated"
	WriteUnchanged = "unchanged"
)

// Dấu mốc mở: version và hash dùng ký tự không có khoảng trắng; nhận cả CRLF.
var startRe = regexp.MustCompile(`(?m)^<!-- dk:agent-context start version=(\S+) hash=(\S+) -->\r?\n`)

// body là nội dung ghi vào file: bỏ xuống dòng cuối để khối kết thúc gọn.
func body(content []byte) string { return strings.TrimRight(string(content), "\n") }

// Hash là sha256 rút gọn 16 hex của thân khối đúng như ghi vào file.
func Hash(content []byte) string {
	sum := sha256.Sum256([]byte(body(content)))
	return hex.EncodeToString(sum[:8])
}

// Block dựng khối đầy đủ gồm dấu mốc mở, nội dung và dấu mốc đóng.
func Block(content []byte, version string) string {
	return fmt.Sprintf("%s version=%s hash=%s -->\n%s\n%s\n", startPrefix, version, Hash(content), body(content), endMarker)
}

// locate tìm vị trí khối trong file: [start, end). Thiếu mốc đóng sau mốc mở
// hoặc có nhiều hơn một mốc mở trả về ok=false với found=true.
func locate(data string) (start, end int, found, ok bool) {
	m := startRe.FindStringIndex(data)
	if m == nil {
		return 0, 0, false, false
	}
	if strings.Count(data, startPrefix) > 1 {
		return m[0], 0, true, false
	}
	start = m[0]
	rest := data[m[1]:]
	i := strings.Index(rest, endMarker)
	if i < 0 {
		return start, 0, true, false
	}
	end = m[1] + i + len(endMarker)
	if strings.HasPrefix(data[end:], "\r\n") {
		end += 2
	} else if strings.HasPrefix(data[end:], "\n") {
		end++
	}
	return start, end, true, true
}

// same so một đoạn trong file với khối mong đợi, coi CRLF như LF.
func same(got, want string) bool {
	return strings.ReplaceAll(got, "\r\n", "\n") == want
}

// Write ghi khối vào file: tạo mới khi chưa có file, thay khối cũ tại chỗ,
// hoặc nối vào cuối file chưa có khối. Trả về created, updated hay unchanged.
func Write(path string, content []byte, version string) (string, error) {
	block := Block(content, version)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return WriteCreated, os.WriteFile(path, []byte(block), 0o644)
	}
	if err != nil {
		return "", err
	}
	s := string(data)
	start, end, found, ok := locate(s)
	var out string
	switch {
	case found && !ok:
		return "", fmt.Errorf("%s: khối thiếu `%s` hoặc có nhiều hơn một khối; sửa tay rồi chạy lại", filepath.Base(path), endMarker)
	case found:
		if same(s[start:end], block) {
			return WriteUnchanged, nil
		}
		out = s[:start] + block + s[end:]
	default:
		out = s
		if out != "" && !strings.HasSuffix(out, "\n") {
			out += "\n"
		}
		if out != "" {
			out += "\n"
		}
		out += block
	}
	return WriteUpdated, os.WriteFile(path, []byte(out), 0o644)
}

// Check trả về trạng thái của khối trong file so với nội dung và phiên bản
// hiện tại: so trọn khối (kể cả dấu mốc) nên phiên bản khác, nội dung nhúng
// khác hay bị sửa tay đều là outdated.
func Check(path string, content []byte, version string) (string, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return StateMissingFile, nil
	}
	if err != nil {
		return "", err
	}
	s := string(data)
	start, end, found, ok := locate(s)
	switch {
	case !found:
		return StateMissing, nil
	case !ok:
		return StateBroken, nil
	case !same(s[start:end], Block(content, version)):
		return StateOutdated, nil
	}
	return StateOK, nil
}
