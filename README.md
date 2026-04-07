# keep-awake

PCがスリープ状態に入るのを防止するクロスプラットフォーム対応のユーティリティです。
macOSではcaffeinateコマンドを利用し、Windowsではマウスカーソルの微小移動によってスリープを抑制します。

## 対応プラットフォーム

- macOS (Apple Silicon / Intel)
- Windows 10/11

## 仕組み

OS別に最適なスリープ防止の手段が自動的に選択されます。

macOSでは、システム標準のcaffeinateプロセスを起動して、ディスプレイのスリープとシステムのアイドルスリープを防止します。

Windowsでは、設定された間隔ごとにマウスカーソルを1ピクセル移動させて元の位置に戻します。
カーソルは一瞬で元の位置に復帰するため、作業の妨げにはなりません。

## 必要な環境

### 開発環境

- Go 1.23以上
- Make
- staticcheck (品質チェック用)
- golangci-lint (品質チェック用)

### 実行環境

- macOS 12以上、またはWindows 10/11

## インストール

リポジトリをクローンしてビルドします。

```bash
git clone https://github.com/okamyuji/keep-awake
cd keep-awake
make build-all
```

macOS向けのみビルドする場合は以下のように実行します。

```bash
make build-macos
```

Windows向けのみビルドする場合は以下のように実行します。

```bash
make build-windows
```

## 使い方

### 基本的な実行

macOSではローカルビルドをそのまま実行できます。

```bash
make run
```

Windowsでは生成されたexeファイルを実行します。

```bash
keep-awake.exe
```

### カスタム設定での実行

間隔を指定して実行できます。

```bash
make run-custom INTERVAL=60
```

コマンドライン引数を直接指定する場合は以下のように実行します。

```bash
./keep-awake-macos -interval 60 -maxmove 3
```

### コマンドラインオプション

| オプション | 説明 | デフォルト値 |
|-----------|------|-------------|
| -interval | スリープ防止の間隔を秒単位で指定します | 180 |
| -maxmove  | マウスの最大移動ピクセル数を指定します (Windows用) | 5 |

macOSではcaffeinateがスリープを常時防止するため、intervalオプションは無視されます。

### 終了方法

Ctrl+Cを押すとプログラムが終了します。
macOSではcaffeinateプロセスも自動的に停止します。

## 開発者向け情報

### makeコマンド一覧

```bash
make help             # 使用可能なコマンドを表示します
make build-windows    # Windows用バイナリをビルドします
make build-macos      # macOS (Apple Silicon)用バイナリをビルドします
make build-macos-intel # macOS (Intel)用バイナリをビルドします
make build-all        # 全プラットフォーム向けにビルドします
make clean            # ビルド成果物を削除します
make test             # テストを実行します
make lint             # 品質チェックを実行します (vet/gofmt/staticcheck/golangci-lint)
make check            # テスト、品質チェック、ビルドをまとめて実行します
make run              # macOSでローカル実行します
make run-custom       # INTERVALを指定して実行します (例: make run-custom INTERVAL=60)
```

### テストの実行

```bash
make test
```

### 品質チェック

すべての静的解析ツールをまとめて実行します。

```bash
make lint
```

テスト、品質チェック、ビルドを一括で実行して問題がないか確認します。

```bash
make check
```

### プロジェクト構成

```
main.go                 # エントリーポイントとシグナルハンドリングを担当します
keeper.go               # Keeperインターフェースと戦略選択のロジックを定義しています
keeper_darwin.go        # macOS向けのcaffeinateによるスリープ防止を実装しています
keeper_windows.go       # Windows向けのマウス移動によるスリープ防止を実装しています
keeper_unsupported.go   # 未対応OS向けのフォールバックを提供します
logger.go               # 標準出力とファイルへの二重出力ロガーを実装しています
Makefile                # ビルド、テスト、品質チェックのタスクを定義しています
```

### ログ出力

プログラムの実行ログは標準出力と、実行ディレクトリのkeep-awake.logファイルの両方に出力されます。
ログファイルは起動のたびに初期化されます。

## セキュリティについて

Windows版はマウスカーソルを自動的に移動させるため、セキュリティソフトウェアに検知される場合があります。
その場合は、セキュリティソフトウェアの除外リストに追加してください。

## ライセンス

MIT License

## 作者

[okamyuji](https://github.com/okamyuji)
