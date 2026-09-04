package target

import "fmt"

// Names là các target có, theo thứ tự ưu tiên.
var Names = []string{"claude", "codex"}

// Get trả về target theo tên. root là gốc dự án cho scope dự án.
func Get(name, root string) (Target, error) {
	switch name {
	case "claude":
		return &Claude{Root: root}, nil
	case "codex":
		return &Codex{Root: root}, nil
	}
	return nil, fmt.Errorf("target %q không có; chọn claude hoặc codex", name)
}
