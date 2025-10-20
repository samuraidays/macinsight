# macinsight

macOS Security Audit CLI. 各種セキュリティ設定を自動チェックし、テーブルまたはJSONで結果を出力します。

## 主な機能

- SIP: System Integrity Protection の有効確認
- Gatekeeper: アプリ実行ポリシーの確認
- FileVault: ディスク暗号化状態の確認
- Firewall: アプリケーション・ファイアウォールの有効確認
- AutoLogin: 自動ログインが無効かの確認
- OSUpdate: OSの更新状況（キャッシュ利用で高速に確認、セキュリティ更新を優先判定）

## インストール / ビルド

Makefile を用意しています。

```bash
# 依存取得
make deps

# Lint (任意: golangci-lint が未導入なら make install-lint)
make lint

# ビルド（bin/macinsight に出力）
make build

# インストール（go install）
make install
```

バージョンは Git のタグ/コミットから自動生成され、`macinsight version` に反映されます。

## 使い方

```bash
# 全チェックを実行（テーブル出力）
./bin/macinsight audit

# JSON 出力
./bin/macinsight audit --json

# 実行するチェックを限定
./bin/macinsight audit --only filevault,gatekeeper

# 実行から除外
./bin/macinsight audit --exclude sip,firewall

# 各チェックのタイムアウト変更（デフォルト 3s）
./bin/macinsight audit --timeout 5s

# 利用可能なチェック一覧
./bin/macinsight list-checks

# バージョン表示（Git情報に基づく動的バージョン）
./bin/macinsight version
```

## 利用可能なチェック

- `sip`: SIP が有効か
- `gatekeeper`: Gatekeeper の有効状態
- `filevault`: FileVault（ディスク暗号化）状態
- `firewall`: アプリケーション・ファイアウォール状態
- `autologin`: 自動ログインが無効か
- `osupdate`: OS 更新状況（`softwareupdate -l --no-scan` を利用）

Evidence（証跡）はシンプルに出力されます。例）

```text
OS updates current ... Evidence: version=26.0.1
```

更新がある場合のみ、簡易な `updates` 情報を付与します。

## 出力形式

- テーブル（デフォルト）: 人間に読みやすい表形式
- JSON（`--json`）: 機械可読なJSON

## テスト

テーブル/JSON出力、各チェックの判定ロジックをユニットテストしています。

```bash
make test
```

## 必要環境

- macOS
- Go 1.20+（ビルド用）

## ライセンス

MIT License - 詳細は [LICENSE](LICENSE) を参照してください。
