# Phiên bản không có tiền tố v, khớp goreleaser ({{ .Version }}) và build info.
VERSION ?= $(shell (git describe --tags --always 2>/dev/null || echo dev) | sed 's/^v//')
LDFLAGS := -s -w -X github.com/anhvandev/doc-kit/internal/cli.Version=$(VERSION)
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64

.PHONY: build test vet install lint-skills build-all

build:
	CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/dk ./cmd/dk

test:
	go test ./...

vet:
	go vet ./...

install:
	CGO_ENABLED=0 go install -trimpath -ldflags "$(LDFLAGS)" ./cmd/dk

# Skill nhúng phải trung lập target: không nhắc skill ak-*, Claude Code, thư mục
# .claude/ hay .codex/, tên tool Edit/Write/MultiEdit; mọi file .md dưới 300 dòng.
lint-skills:
	@! grep -rnE 'ak-|Claude Code|\.claude/|\.codex/' assets/skills || { echo "lint-skills: từ cấm trong assets/skills"; exit 1; }
	@! grep -rnwE 'Edit|Write|MultiEdit' assets/skills || { echo "lint-skills: tên tool trong assets/skills"; exit 1; }
	@! find assets/skills -name '*.md' -exec awk 'END { if (NR >= 300) { print FILENAME ": " NR " dòng" } }' {} \; | grep . || { echo "lint-skills: file skill từ 300 dòng"; exit 1; }
	@echo "lint-skills: sạch"

# Kiểm ma trận build của goreleaser mà không cần cài goreleaser.
build-all:
	@for p in $(PLATFORMS); do \
		os=$${p%/*}; arch=$${p#*/}; ext=""; [ "$$os" = windows ] && ext=.exe; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -trimpath -ldflags "$(LDFLAGS)" -o dist/dk-$$os-$$arch$$ext ./cmd/dk || exit 1; \
		echo "dist/dk-$$os-$$arch$$ext"; \
	done
