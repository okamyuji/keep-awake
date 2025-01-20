# Windows PC スリープ防止ユーティリティ

このプログラムは、Windows PCがスリープ状態に入るのを防ぐため、定期的にマウスカーソルを微小に動かすユーティリティです。
macOS上でクロスコンパイルして、Windows用の実行ファイルを生成します。

## 機能

- 設定可能な間隔でマウスカーソルを自動的に動かします
- 元の位置に自動的に戻るため、作業の妨げになりません
- コマンドライン引数でカスタマイズ可能
- 安全な終了機能（Ctrl+C）をサポート

## 必要要件

### 開発環境（macOS）

- Go 1.16以上
- Make
- git

### 実行環境（Windows）

- Windows 10/11

## インストール方法

1. リポジトリのクローン:

    ```bash
    git clone [リポジトリURL]
    cd keep-awake
    ```

2. ビルド:

    ```bash
    make all
    ```

これにより、`keep-awake.exe`が生成されます。

## 使用方法

### 基本的な実行

```bash
./keep-awake.exe
```

### カスタム設定での実行

```bash
./keep-awake.exe -interval 60 -maxmove 5
```

### コマンドラインオプション

- `-interval`: マウス移動の間隔（秒）
    - デフォルト: 180秒
- `-maxmove`: 最大移動ピクセル数
    - デフォルト: 5ピクセル

## 開発者向け情報

### makeコマンド一覧

```bash
make help      # ヘルプの表示
make init      # プロジェクトの初期化
make deps      # 依存パッケージのインストール
make build     # Windowsバイナリのビルド
make clean     # ビルド成果物の削除
make test      # テストの実行
make all       # 全処理の実行
```

### カスタム間隔での実行（開発時）

```bash
make run-custom INTERVAL=60
```

## セキュリティ上の注意

- プログラムはWindows上で自動的にマウスを動かすため、セキュリティソフトウェアによって検知される可能性があります
- 必要に応じて、セキュリティソフトウェアの除外リストに追加してください

## トラブルシューティング

### よくある問題と解決方法

1. ビルドエラーが発生する場合

    ```bash
    make clean
    make all
    ```

2. 実行時にアクセス権限エラーが発生する場合

- Windowsでプログラムを管理者として実行してください

## ライセンス

MIT License

## 貢献について

1. Forkを作成
2. 新しいブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチをPush (`git push origin feature/amazing-feature`)
5. Pull Requestを作成

## サポート

問題や提案がある場合は、Issueを作成してください。

## 作者

[okamyuji](https://github.com/okamyuji)

---
最終更新日: 2024年12月27日
