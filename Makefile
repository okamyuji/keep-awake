# 基本変数の設定
GO=go
GOBIN=$(shell $(GO) env GOPATH)/bin
STATICCHECK=$(GOBIN)/staticcheck
GOLANGCI_LINT=$(GOBIN)/golangci-lint

# デフォルトのターゲット
.DEFAULT_GOAL := help

# ヘルプコマンド
.PHONY: help
help:
	@echo "使用可能なコマンド:"
	@echo "  make build-windows  - Windowsバイナリのビルド"
	@echo "  make build-macos    - macOSバイナリのビルド"
	@echo "  make build-all      - 全プラットフォームのビルド"
	@echo "  make clean          - ビルド成果物の削除"
	@echo "  make test           - テストの実行"
	@echo "  make tools          - lint用ツールのインストール"
	@echo "  make lint           - 品質チェック（vet/fmt/staticcheck/golangci-lint）"
	@echo "  make check          - テスト＋品質チェック＋ビルド（全検証）"
	@echo "  make run            - ローカル実行（macOS）"

# Windowsバイナリのビルド
.PHONY: build-windows
build-windows:
	@echo "Windowsバイナリをビルドしています..."
	GOOS=windows GOARCH=amd64 $(GO) build -o keep-awake.exe

# macOSバイナリのビルド（Apple Silicon）
.PHONY: build-macos
build-macos:
	@echo "macOSバイナリをビルドしています..."
	GOOS=darwin GOARCH=arm64 $(GO) build -o keep-awake-macos

# macOSバイナリのビルド（Intel）
.PHONY: build-macos-intel
build-macos-intel:
	@echo "macOS (Intel) バイナリをビルドしています..."
	GOOS=darwin GOARCH=amd64 $(GO) build -o keep-awake-macos-intel

# 全プラットフォームのビルド
.PHONY: build-all
build-all: build-windows build-macos build-macos-intel

# ビルド成果物の削除
.PHONY: clean
clean:
	@echo "ビルド成果物を削除しています..."
	rm -f keep-awake.exe keep-awake-macos keep-awake-macos-intel
	go clean

# テストの実行
.PHONY: test
test:
	@echo "テストを実行しています..."
	$(GO) test -v ./...

# lint用ツールのインストール
$(STATICCHECK):
	$(GO) install honnef.co/go/tools/cmd/staticcheck@latest

$(GOBIN)/golangci-lint:
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: tools
tools: $(STATICCHECK) $(GOBIN)/golangci-lint
	@echo "lint用ツールのインストール完了"

# 品質チェック
.PHONY: lint
lint: $(STATICCHECK) $(GOBIN)/golangci-lint
	@echo "品質チェックを実行しています..."
	@echo "--- go vet ---"
	$(GO) vet ./...
	@echo "--- gofmt ---"
	@test -z "$$(gofmt -l .)" || (echo "以下のファイルにフォーマットの問題があります:"; gofmt -l .; exit 1)
	@echo "--- staticcheck ---"
	$(STATICCHECK) ./...
	@echo "--- golangci-lint ---"
	$(GOLANGCI_LINT) run ./...
	@echo "品質チェック完了: 問題なし"

# 全検証（テスト＋品質チェック＋ビルド）
.PHONY: check
check: test lint build-all
	@echo "全検証完了"

# ローカル実行
.PHONY: run
run:
	@echo "プログラムを実行しています..."
	$(GO) run .

# カスタム間隔での実行（例：make run-custom INTERVAL=60）
.PHONY: run-custom
run-custom:
	@test -n "$(INTERVAL)" || (echo "INTERVAL を指定してください。例: make run-custom INTERVAL=60"; exit 1)
	@echo "$(INTERVAL)" | grep -qE '^[1-9][0-9]*$$' || (echo "INTERVAL は正の整数を指定してください: $(INTERVAL)"; exit 1)
	@echo "カスタム設定でプログラムを実行しています..."
	$(GO) run . -interval $(INTERVAL)
