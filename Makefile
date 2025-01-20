# 基本変数の設定
BINARY_NAME=keep-awake.exe
GO=go
GOOS=windows
GOARCH=amd64

# デフォルトのターゲット
.DEFAULT_GOAL := help

# ヘルプコマンド
.PHONY: help
help:
	@echo "使用可能なコマンド:"
	@echo "  make init     - プロジェクトの初期化"
	@echo "  make deps    - 依存パッケージのインストール"
	@echo "  make build   - Windowsバイナリのビルド"
	@echo "  make clean   - ビルド成果物の削除"
	@echo "  make test    - テストの実行"
	@echo "  make all     - 依存関係のインストールからビルドまでを実行"

# プロジェクトの初期化
.PHONY: init
init:
	@echo "プロジェクトを初期化しています..."
	mkdir -p keep-awake
	cd keep-awake && $(GO) mod init keep-awake

# 依存パッケージのインストール
.PHONY: deps
deps:
	@echo "依存パッケージをインストールしています..."
	$(GO) get golang.org/x/sys/windows

# Windowsバイナリのビルド
.PHONY: build
build:
	@echo "Windowsバイナリをビルドしています..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o $(BINARY_NAME)

# ビルド成果物の削除
.PHONY: clean
clean:
	@echo "ビルド成果物を削除しています..."
	rm -f $(BINARY_NAME)
	go clean

# テストの実行
.PHONY: test
test:
	@echo "テストを実行しています..."
	$(GO) test -v ./...

# すべての処理を実行
.PHONY: all
all: clean deps build

# 実行可能ファイルの実行（開発時のテスト用）
.PHONY: run
run:
	@echo "プログラムを実行しています..."
	./$(BINARY_NAME)

# カスタム間隔での実行（例：make run-custom INTERVAL=60）
.PHONY: run-custom
run-custom:
	@echo "カスタム設定でプログラムを実行しています..."
	./$(BINARY_NAME) -interval $(INTERVAL)