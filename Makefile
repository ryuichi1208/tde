TOOL_NAME = tde
PKG_DIR = .
BIN_DIR = ./bin

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

BINARY = $(BIN_DIR)/$(TOOL_NAME)

# ビルド設定
.PHONY: all
all: build

# ビルド
.PHONY: build
build:
	@echo "Building $(TOOL_NAME)..."
	@mkdir -p $(BIN_DIR)
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINARY) $(PKG_DIR)
	@echo "Build completed: $(BINARY)"

# テスト
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

# クリーン
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)
	@echo "Cleanup completed."

# インストール
.PHONY: install
install: build
	@echo "Installing $(TOOL_NAME)..."
	@install -m 755 $(BINARY) /usr/local/bin/$(TOOL_NAME)
	@echo "$(TOOL_NAME) installed to /usr/local/bin"

# アンインストール
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(TOOL_NAME)..."
	@rm -f /usr/local/bin/$(TOOL_NAME)
	@echo "$(TOOL_NAME) uninstalled from /usr/local/bin"

# ヘルプ
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make build       - Build the $(TOOL_NAME) binary"
	@echo "  make test        - Run tests"
	@echo "  make clean       - Clean the build artifacts"
	@echo "  make install     - Install $(TOOL_NAME) to /usr/local/bin"
	@echo "  make uninstall   - Uninstall $(TOOL_NAME) from /usr/local/bin"
