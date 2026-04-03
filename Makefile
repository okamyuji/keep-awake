# 基本変数の設定
GO=go

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
build-all: build-windows build-macos

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

# ローカル実行
.PHONY: run
run:
	@echo "プログラムを実行しています..."
	$(GO) run .

# カスタム間隔での実行（例：make run-custom INTERVAL=60）
.PHONY: run-custom
run-custom:
	@echo "カスタム設定でプログラムを実行しています..."
	$(GO) run . -interval $(INTERVAL)
