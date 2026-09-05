package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/anhvandev/doc-kit/assets"
	"github.com/anhvandev/doc-kit/internal/agentctx"
	"github.com/anhvandev/doc-kit/internal/config"
	"github.com/anhvandev/doc-kit/internal/gitx"
	"github.com/anhvandev/doc-kit/internal/hook"
	"github.com/anhvandev/doc-kit/internal/skill"
	"github.com/anhvandev/doc-kit/internal/target"
)

// doctorRow là một dòng của bảng `dk doctor`.
type doctorRow struct {
	Item   string `json:"item"`
	Status string `json:"status"`
	OK     bool   `json:"ok"`
	Fix    string `json:"fix,omitempty"`
}

func newDoctorCmd(a *app) *cobra.Command {
	var f targetFlags
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Kiểm tra dự án hiện tại: dk.toml, git, pre-commit, skill, hook",
		Args:  exactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return a.runDoctor(cmd, f)
		},
	}
	f.bind(cmd, true)
	return cmd
}

// runDoctor gom các kiểm tra rồi in bảng `mục | trạng thái | cách sửa`;
// mã thoát 3 khi có mục chưa đạt.
func (a *app) runDoctor(cmd *cobra.Command, f targetFlags) error {
	var rows []doctorRow
	add := func(item, status string, ok bool, fix string) {
		if ok {
			fix = ""
		}
		rows = append(rows, doctorRow{Item: item, Status: status, OK: ok, Fix: fix})
	}

	root, ok := findProjectRoot(a.cwd)
	if !ok {
		add(config.FileName, "không tìm thấy từ "+a.cwd+" trở lên", false, "chạy `dk init` tại gốc dự án")
		return a.printDoctor(rows)
	}
	a.root = root
	cfg, err := config.Load(filepath.Join(root, config.FileName))
	if err != nil {
		add(config.FileName, err.Error(), false, "sửa cú pháp TOML hoặc `dk init --force`")
		return a.printDoctor(rows)
	}
	a.cfg = cfg
	add(config.FileName, "có tại "+root, true, "")

	if _, err := os.Stat(a.docsDir()); err != nil {
		add(cfg.DocsDir+"/", "thiếu", false, "chạy `dk init --force` để tạo lại cây thư mục")
	} else {
		add(cfg.DocsDir+"/", "có", true, "")
	}

	if _, err := exec.LookPath("dk"); err != nil {
		add("dk trên PATH", "chưa có", false, "`make install` hoặc chép binary vào thư mục trong PATH; pre-commit và hook gọi `dk` qua PATH")
	} else {
		add("dk trên PATH", "có", true, "")
	}

	if !gitx.IsRepo(root) {
		add("git", "không phải repo git", false, "`git init`; changelog đếm dòng và pre-commit cần git")
		add("pre-commit", "bỏ qua (không git)", false, "`git init` rồi `dk init --force`")
	} else {
		add("git", "có", true, "")
		a.doctorPreCommit(add)
	}
	a.doctorAgentContext(add)

	names := f.names()
	// Không nêu --target mà dự án đã cài gì đó cho Codex (skill hoặc hook):
	// kiểm cả Codex. Thư mục .codex/ rỗng sau khi gỡ không tính.
	if !cmd.Flags().Changed("target") && !slices.Contains(names, "codex") {
		for _, p := range []string{filepath.Join(root, ".codex", "skills"), filepath.Join(root, ".codex", "hooks.json")} {
			if _, err := os.Stat(p); err == nil {
				names = append(names, "codex")
				break
			}
		}
	}
	ts, err := a.targetsOf(names, root)
	if err != nil {
		add("target "+f.name, err.Error(), false, "chọn `--target claude`, `codex` hoặc `claude,codex`")
		return a.printDoctor(rows)
	}
	for _, t := range ts {
		a.doctorTarget(t, f.global, add)
	}
	return a.printDoctor(rows)
}

// doctorTarget thêm các dòng skill và hook của một target.
func (a *app) doctorTarget(t target.Target, global bool, add func(string, string, bool, string)) {
	scope := t.Name() + ", dự án"
	targetFlag := ""
	if t.Name() != target.Names[0] { // target mặc định không cần cờ
		targetFlag = " --target " + t.Name()
	}
	if global {
		scope = t.Name() + ", toàn máy"
		targetFlag += " --global"
	}
	statuses, err := skill.Status(t, global, Version)
	if err != nil {
		add("skill ("+scope+")", err.Error(), false, "kiểm tra thư mục skill của target")
		return
	}
	current := 0
	for _, r := range statuses {
		fix := ""
		switch {
		case r.State == skill.StateCurrent:
			current++
			continue
		case r.State == skill.StateMissing:
			fix = "`dk skill install" + targetFlag + " " + r.Name + "`"
		case r.State == skill.StateModified:
			fix = "`dk skill install --force" + targetFlag + " " + r.Name + "` (mất sửa tay) hoặc giữ nguyên có chủ đích"
		case r.State == skill.StateForeign:
			fix = "thư mục trùng tên không do dk cài; đổi tên hoặc xóa tay rồi `dk skill install" + targetFlag + " " + r.Name + "`"
		default: // cũ (vX)
			fix = "`dk skill install" + targetFlag + " " + r.Name + "` để lên " + Version
		}
		add("skill "+r.Name+" ("+t.Name()+")", r.State, false, fix)
	}
	if current == len(statuses) {
		add("skill ("+scope+")", "đủ, đúng phiên bản", true, "")
	} else if current > 0 {
		add("skill ("+scope+")", "đúng phiên bản "+strconv.Itoa(current)+"/"+strconv.Itoa(len(statuses)), false, "xem các dòng skill ở trên")
	}

	got, err := t.InstalledHooks(global)
	trust := ""
	if t.Name() == "codex" {
		trust = "; " + target.CodexTrustNote
	}
	want := t.MapHooks(hook.Entries())
	fix := "`dk hook install" + targetFlag + "`" + trust
	switch {
	case err != nil:
		add("hook ("+scope+")", err.Error(), false, "sửa cấu hình hook của target rồi `dk hook install"+targetFlag+"`")
	case sameHooks(got, want):
		add("hook ("+scope+")", "đủ "+strconv.Itoa(len(got))+" hook", true, "")
	case len(got) < len(want) && subsetHooks(got, want):
		add("hook ("+scope+")", "có "+strconv.Itoa(len(got))+"/"+strconv.Itoa(len(want)), false, fix)
	default:
		add("hook ("+scope+")", "lệch bản hiện tại (matcher hoặc lệnh khác)", false, fix)
	}
}

// subsetHooks: mọi hook trong got đều có trong want (so event, matcher, lệnh).
func subsetHooks(got, want []target.HookEntry) bool {
	for _, g := range got {
		found := false
		for _, w := range want {
			if g == w {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// sameHooks: got và want cùng tập hook, không thừa không thiếu.
func sameHooks(got, want []target.HookEntry) bool {
	return len(got) == len(want) && subsetHooks(got, want) && subsetHooks(want, got)
}

// doctorPreCommit kiểm script pre-commit có mặt và gọi dk.
func (a *app) doctorPreCommit(add func(string, string, bool, string)) {
	hooks, err := gitx.HooksDir(a.root)
	if err != nil {
		add("pre-commit", err.Error(), false, "kiểm tra `git rev-parse --git-path hooks`")
		return
	}
	path := filepath.Join(hooks, "pre-commit")
	b, err := os.ReadFile(path)
	if err != nil {
		add("pre-commit", "thiếu "+path, false, "`dk init --force` để chép script")
		return
	}
	if !strings.Contains(string(b), "changelog pending") {
		add("pre-commit", "có nhưng không gọi `dk changelog pending`", false, "`dk init --force` in đoạn cần thêm vào "+path)
		return
	}
	// Windows không có bit chạy; Git for Windows chạy hook theo tên file.
	if fi, err := os.Stat(path); err == nil && runtime.GOOS != "windows" && fi.Mode()&0o111 == 0 {
		add("pre-commit", "không có quyền chạy", false, "`chmod +x "+path+"`")
		return
	}
	add("pre-commit", "có, gọi dk", true, "")
}

// doctorAgentContext kiểm khối ngữ cảnh agent trong từng file ngữ cảnh:
// thiếu, cũ hay bị sửa tay đều sửa bằng `dk init --agent-context`.
func (a *app) doctorAgentContext(add func(string, string, bool, string)) {
	content, err := assets.FS.ReadFile("agent-context.md")
	if err != nil {
		add("agent context", err.Error(), false, "binary hỏng; cài lại dk")
		return
	}
	for _, name := range agentctx.Files {
		st, err := agentctx.Check(filepath.Join(a.root, name), content, Version)
		item := "agent context (" + name + ")"
		switch {
		case err != nil:
			add(item, err.Error(), false, "kiểm quyền đọc "+name)
		case st == agentctx.StateOK:
			add(item, "có, đúng phiên bản", true, "")
		case st == agentctx.StateBroken:
			add(item, "khối thiếu mốc đóng hoặc có nhiều hơn một khối", false, "sửa tay "+name+" rồi `dk init --agent-context`")
		case st == agentctx.StateMissingFile:
			add(item, "thiếu file", false, "`dk init --agent-context`")
		case st == agentctx.StateMissing:
			add(item, "có file, chưa có khối", false, "`dk init --agent-context`")
		default: // outdated
			add(item, "cũ hoặc bị sửa tay", false, "`dk init --agent-context`")
		}
	}
}

func (a *app) printDoctor(rows []doctorRow) error {
	bad := 0
	for _, r := range rows {
		if !r.OK {
			bad++
		}
	}
	if a.json {
		if err := a.printJSON(rows); err != nil {
			return err
		}
	} else {
		w := tabwriter.NewWriter(a.out, 0, 0, 2, ' ', 0)
		w.Write([]byte("mục\t| trạng thái\t| cách sửa\n"))
		for _, r := range rows {
			mark := "OK"
			if !r.OK {
				mark = "!!"
			}
			w.Write([]byte(mark + " " + r.Item + "\t| " + r.Status + "\t| " + r.Fix + "\n"))
		}
		if err := w.Flush(); err != nil {
			return err
		}
	}
	if bad > 0 {
		return fail(codeCheck, "%d mục chưa đạt", bad)
	}
	return nil
}
