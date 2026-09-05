package target

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// HookCommandPrefix nhận diện hook do dk cài trong cấu hình.
const HookCommandPrefix = "dk hook run"

// Claude là target Claude Code: skill trong .claude/skills, hook trong
// .claude/settings.json khóa hooks.
type Claude struct {
	Root string // gốc dự án; dùng cho scope dự án
}

func (c *Claude) Name() string { return "claude" }

func (c *Claude) base(global bool) (string, error) {
	if global {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".claude"), nil
	}
	if c.Root == "" {
		return "", errors.New("chưa có gốc dự án cho scope dự án")
	}
	return filepath.Join(c.Root, ".claude"), nil
}

func (c *Claude) SkillDir(global bool) (string, error) {
	b, err := c.base(global)
	if err != nil {
		return "", err
	}
	return filepath.Join(b, "skills"), nil
}

// SettingsPath trả về đường dẫn settings.json theo scope.
func (c *Claude) SettingsPath(global bool) (string, error) {
	b, err := c.base(global)
	if err != nil {
		return "", err
	}
	return filepath.Join(b, "settings.json"), nil
}

// object là JSON object giữ thứ tự khóa; giá trị để thô để không đụng
// phần không thuộc dk.
type object struct {
	keys []string
	vals map[string]json.RawMessage
}

func newObject() *object { return &object{vals: map[string]json.RawMessage{}} }

func (o *object) get(k string) (json.RawMessage, bool) { v, ok := o.vals[k]; return v, ok }

func (o *object) set(k string, v json.RawMessage) {
	if _, ok := o.vals[k]; !ok {
		o.keys = append(o.keys, k)
	}
	o.vals[k] = v
}

func (o *object) del(k string) {
	if _, ok := o.vals[k]; !ok {
		return
	}
	delete(o.vals, k)
	for i, key := range o.keys {
		if key == k {
			o.keys = append(o.keys[:i], o.keys[i+1:]...)
			return
		}
	}
}

func (o *object) UnmarshalJSON(b []byte) error {
	dec := json.NewDecoder(bytes.NewReader(b))
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	if tok != json.Delim('{') {
		return errors.New("không phải JSON object")
	}
	o.keys, o.vals = nil, map[string]json.RawMessage{}
	for dec.More() {
		kt, err := dec.Token()
		if err != nil {
			return err
		}
		k, _ := kt.(string)
		var v json.RawMessage
		if err := dec.Decode(&v); err != nil {
			return err
		}
		o.set(k, v)
	}
	_, err = dec.Token()
	return err
}

func (o *object) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, k := range o.keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		kb, _ := json.Marshal(k)
		buf.Write(kb)
		buf.WriteByte(':')
		var c bytes.Buffer
		if err := json.Compact(&c, o.vals[k]); err != nil {
			return nil, err
		}
		buf.Write(c.Bytes())
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// commandHook là một hook dạng lệnh shell trong settings.json.
type commandHook struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

// hookGroup là một mục trong hooks.<Event>; chỉ đọc phần dk cần.
type hookGroup struct {
	Matcher string `json:"matcher"`
	Hooks   []struct {
		Command string `json:"command"`
	} `json:"hooks"`
}

func (g hookGroup) has(match func(string) bool) bool {
	for _, h := range g.Hooks {
		if match(h.Command) {
			return true
		}
	}
	return false
}

// marshal mã hóa không escape HTML để lệnh có ">" hay "&" giữ nguyên.
func marshal(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

func joinArray(items []json.RawMessage) json.RawMessage {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, it := range items {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.Write(it)
	}
	buf.WriteByte(']')
	return buf.Bytes()
}

func readSettings(path string) (*object, error) {
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return newObject(), nil
	}
	if err != nil {
		return nil, err
	}
	o := newObject()
	if len(bytes.TrimSpace(b)) == 0 {
		return o, nil
	}
	if err := json.Unmarshal(b, o); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return o, nil
}

// writeSettings ghi settings với thụt 2 khoảng; object rỗng thì xóa file
// (chỉ còn "{}" là do dk tạo hoặc không còn giá trị).
func writeSettings(path string, o *object) error {
	if len(o.keys) == 0 {
		err := os.Remove(path)
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	raw, err := marshal(o)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := json.Indent(&buf, raw, "", "  "); err != nil {
		return err
	}
	buf.WriteByte('\n')
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o644)
}

func hooksOf(o *object) (*object, error) {
	h := newObject()
	if raw, ok := o.get("hooks"); ok && string(bytes.TrimSpace(raw)) != "null" {
		if err := json.Unmarshal(raw, h); err != nil {
			return nil, fmt.Errorf("khóa hooks: %w", err)
		}
	}
	return h, nil
}

func eventOf(h *object, event string) ([]json.RawMessage, error) {
	var arr []json.RawMessage
	if raw, ok := h.get(event); ok {
		if err := json.Unmarshal(raw, &arr); err != nil {
			return nil, fmt.Errorf("hooks.%s: %w", event, err)
		}
	}
	return arr, nil
}

func putHooks(o, h *object) error {
	if len(h.keys) == 0 {
		o.del("hooks")
		return nil
	}
	raw, err := marshal(h)
	if err != nil {
		return err
	}
	o.set("hooks", raw)
	return nil
}

// MapHooks giữ nguyên: matcher của entries đã là tên tool Claude Code.
func (c *Claude) MapHooks(entries []HookEntry) []HookEntry { return entries }

func (c *Claude) InstallHooks(global bool, entries []HookEntry) error {
	path, err := c.SettingsPath(global)
	if err != nil {
		return err
	}
	return installHooksFile(path, c.MapHooks(entries))
}

func (c *Claude) UninstallHooks(global bool) error {
	path, err := c.SettingsPath(global)
	if err != nil {
		return err
	}
	return uninstallHooksFile(path)
}

// InstalledHooks liệt kê lệnh có tiền tố HookCommandPrefix trong settings.json.
func (c *Claude) InstalledHooks(global bool) ([]HookEntry, error) {
	path, err := c.SettingsPath(global)
	if err != nil {
		return nil, err
	}
	return installedHooksFile(path)
}

// installHooksFile thêm entries vào khóa hooks của file JSON tại path; Claude
// (settings.json) và Codex (hooks.json) dùng cùng cấu trúc khối hooks.
func installHooksFile(path string, entries []HookEntry) error {
	o, err := readSettings(path)
	if err != nil {
		return err
	}
	h, err := hooksOf(o)
	if err != nil {
		return err
	}
	// Bỏ mọi lệnh dk đang có trước, kể cả mục có matcher hay lệnh của bản cũ,
	// rồi ghi entries; mục của người dùng giữ nguyên.
	if err := stripAllDKHooks(h); err != nil {
		return err
	}
	for _, e := range entries {
		arr, err := eventOf(h, e.Event)
		if err != nil {
			return err
		}
		entry := struct {
			Matcher string        `json:"matcher"`
			Hooks   []commandHook `json:"hooks"`
		}{e.Matcher, []commandHook{{Type: "command", Command: e.Command}}}
		raw, err := marshal(entry)
		if err != nil {
			return err
		}
		h.set(e.Event, joinArray(append(arr, raw)))
	}
	if err := putHooks(o, h); err != nil {
		return err
	}
	return writeSettings(path, o)
}

// stripDKHooks bỏ các lệnh dk khỏi một mục matcher, giữ lệnh khác của người
// dùng trong cùng mục; trả nil khi mục không còn lệnh nào.
func stripDKHooks(raw json.RawMessage) (json.RawMessage, error) {
	g := newObject()
	if err := json.Unmarshal(raw, g); err != nil {
		return raw, nil // mục không phải object: không phải của dk, giữ nguyên
	}
	var hooks []json.RawMessage
	if hr, ok := g.get("hooks"); ok && string(bytes.TrimSpace(hr)) != "null" {
		if err := json.Unmarshal(hr, &hooks); err != nil {
			return nil, fmt.Errorf("hooks[].hooks: %w", err)
		}
	}
	kept := hooks[:0]
	for _, hraw := range hooks {
		var ch commandHook
		if json.Unmarshal(hraw, &ch) == nil && strings.HasPrefix(ch.Command, HookCommandPrefix) {
			continue
		}
		kept = append(kept, hraw)
	}
	if len(kept) == len(hooks) {
		return raw, nil // không có lệnh dk: giữ nguyên byte
	}
	if len(kept) == 0 {
		return nil, nil
	}
	g.set("hooks", joinArray(kept))
	return marshal(g)
}

func uninstallHooksFile(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	o, err := readSettings(path)
	if err != nil {
		return err
	}
	h, err := hooksOf(o)
	if err != nil {
		return err
	}
	if err := stripAllDKHooks(h); err != nil {
		return err
	}
	if err := putHooks(o, h); err != nil {
		return err
	}
	return writeSettings(path, o)
}

// stripAllDKHooks bỏ lệnh dk khỏi mọi event; event không còn mục nào thì xóa.
func stripAllDKHooks(h *object) error {
	for _, event := range append([]string(nil), h.keys...) {
		arr, err := eventOf(h, event)
		if err != nil {
			return err
		}
		kept := arr[:0]
		for _, raw := range arr {
			g, err := stripDKHooks(raw)
			if err != nil {
				return err
			}
			if g != nil {
				kept = append(kept, g)
			}
		}
		if len(kept) == 0 {
			h.del(event)
			continue
		}
		h.set(event, joinArray(kept))
	}
	return nil
}

func installedHooksFile(path string) ([]HookEntry, error) {
	o, err := readSettings(path)
	if err != nil {
		return nil, err
	}
	h, err := hooksOf(o)
	if err != nil {
		return nil, err
	}
	var out []HookEntry
	for _, event := range h.keys {
		arr, err := eventOf(h, event)
		if err != nil {
			return nil, err
		}
		for _, raw := range arr {
			var g hookGroup
			if json.Unmarshal(raw, &g) != nil {
				continue
			}
			for _, hk := range g.Hooks {
				if strings.HasPrefix(hk.Command, HookCommandPrefix) {
					out = append(out, HookEntry{Event: event, Matcher: g.Matcher, Command: hk.Command})
				}
			}
		}
	}
	return out, nil
}
