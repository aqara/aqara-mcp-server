<div align="center" style="display: flex; align-items: center; justify-content: center; ">

  <img src="/readme/img/logo.png" alt="Aqara Logo" height="120">
  <h1>MCP Server</h1>

</div>

<div align="center">

[English](/readme/README.md) | [中文](/readme/README_CN.md) | [繁體中文](/readme/README_CHT.md) | [Français](/readme/README_FR.md) | [한국어](/readme/README_KR.md) | [Español](/readme/README_ES.md) | 日本語 | [Deutsch](/readme/README_DE.md) | [Italiano](/readme/README_IT.md)

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/aqara/aqara-mcp-server)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![Release](https://img.shields.io/github/v/release/aqara/aqara-mcp-server)](https://github.com/aqara/aqara-mcp-server/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

Aqara MCP Serverは、[MCP (Model Context Protocol)](https://modelcontextprotocol.io/introduction)プロトコルに基づいて開発されたスマートホーム制御サービスです。MCPプロトコルをサポートするあらゆるAIアシスタントやAPI（Claude、Cursorなど）がAqaraスマートホームデバイスと連携し、自然言語によるデバイス制御、状態確認、シーン実行などの機能を実現できます。

## 目次

- [目次](#目次)
- [特徴](#特徴)
- [動作原理](#動作原理)
- [クイックスタート](#クイックスタート)
  - [前提条件](#前提条件)
  - [インストール](#インストール)
    - [方法1：プリコンパイル版のダウンロード（推奨）](#方法1プリコンパイル版のダウンロード推奨)
    - [方法2：ソースからビルド](#方法2ソースからビルド)
  - [Aqaraアカウント認証](#aqaraアカウント認証)
  - [クライアント設定](#クライアント設定)
    - [Claude for Desktop設定例](#claude-for-desktop設定例)
    - [設定パラメータの説明](#設定パラメータの説明)
    - [その他のMCPクライアント](#その他のmcpクライアント)
  - [サービス起動](#サービス起動)
    - [標準モード（推奨）](#標準モード推奨)
    - [HTTPモード（オプション）](#httpモードオプション)
- [APIツール説明](#apiツール説明)
  - [デバイス制御類](#デバイス制御類)
    - [device\_control](#device_control)
  - [デバイス照会類](#デバイス照会類)
    - [device\_query](#device_query)
    - [device\_status\_query](#device_status_query)
    - [device\_log\_query](#device_log_query)
  - [シーン管理類](#シーン管理類)
    - [get\_scenes](#get_scenes)
    - [run\_scenes](#run_scenes)
  - [家庭管理類](#家庭管理類)
    - [get\_homes](#get_homes)
    - [switch\_home](#switch_home)
  - [自動化設定類](#自動化設定類)
    - [automation\_config](#automation_config)
- [プロジェクト構成](#プロジェクト構成)
  - [ディレクトリ構成](#ディレクトリ構成)
  - [コアファイル説明](#コアファイル説明)
- [開発ガイド](#開発ガイド)
- [ライセンス](#ライセンス)

## 特徴

- **包括的なデバイス制御**：Aqaraスマートデバイスのオン/オフ、明度、色温度、モードなど多様な属性の精密制御をサポート
- **柔軟なデバイス照会**：部屋やデバイスタイプ別にデバイスリストと詳細状態を照会可能
- **インテリジェントなシーン管理**：ユーザーが事前設定したスマートホームシーンの照会と実行をサポート
- **デバイス履歴記録**：指定期間内のデバイス状態変更履歴を照会
- **自動化設定**：タイマーや遅延デバイス制御タスクの設定をサポート
- **マルチホームサポート**：ユーザーアカウント下の異なる家庭の照会と切り替えをサポート
- **MCPプロトコル互換**：MCPプロトコル仕様に完全準拠し、各種AIアシスタントとの統合が容易
- **安全な認証メカニズム**：ログイン認証+署名による安全認証を採用し、ユーザーデータとデバイスのセキュリティを保護
- **クロスプラットフォーム動作**：Go言語で開発され、マルチプラットフォーム実行ファイルにコンパイル可能
- **拡張が容易**：モジュラー設計により、新しいツールや機能を簡単に追加可能

## 動作原理

Aqara MCP Serverは、AIアシスタントとAqaraスマートホームプラットフォーム間の架け橋として機能します：

1. **AIアシスタント（MCPクライアント）**：ユーザーがAIアシスタントを通じて指示を出します（例：「リビングのライトを点けて」）
2. **MCPクライアント**：ユーザーの指示を解析し、MCPプロトコルに従ってAqara MCP Serverが提供する対応ツール（例：`device_control`）を呼び出します
3. **Aqara MCP Server（本プロジェクト）**：クライアントからのリクエストを受信し、検証後に`smh.go`モジュールを呼び出します
4. **`smh.go`モジュール**：設定されたAqara認証情報を使用してAqaraクラウドAPIと通信し、実際のデバイス操作やデータ照会を実行します
5. **応答フロー**：AqaraクラウドAPIが結果を返し、Aqara MCP Server経由でMCPクライアントに送信され、最終的にユーザーに表示されます

## クイックスタート

### 前提条件

- Go（バージョン1.24以上）
- Git（ソースからビルドする場合）
- Aqaraアカウントと登録済みスマートデバイス

### インストール

プリコンパイルされた実行ファイルのダウンロードまたはソースからのビルドを選択できます。

#### 方法1：プリコンパイル版のダウンロード（推奨）

GitHub Releasesページにアクセスし、お使いのオペレーティングシステムに適した最新の実行ファイルをダウンロードしてください：

**📥 [Releasesページでダウンロード](https://github.com/aqara/aqara-mcp-server/releases)**

対応プラットフォームの圧縮ファイルをダウンロード後、解凍してすぐに使用できます。

#### 方法2：ソースからビルド

```bash
# リポジトリをクローン
git clone https://github.com/aqara/aqara-mcp-server.git
cd aqara-mcp-server

# 依存関係をダウンロード
go mod tidy

# 実行ファイルをビルド
go build -o aqara-mcp-server
```

ビルド完了後、現在のディレクトリに`aqara-mcp-server`実行ファイルが生成されます。

### Aqaraアカウント認証

MCP ServerがAqaraアカウントにアクセスしてデバイスを制御できるようにするため、まずログイン認証を行う必要があります。

以下のアドレスにアクセスしてログイン認証を完了してください：
**🔗 [https://cdn.aqara.com/app/mcpserver/login.html](https://cdn.aqara.com/app/mcpserver/login.html)**

ログイン成功後、必要な認証情報（`token`、`region`など）を取得できます。これらの情報は後続の設定手順で使用されます。

> ⚠️ **セキュリティ注意**：`token`情報を適切に管理し、他人に漏洩しないようにしてください。

### クライアント設定

MCPクライアントによって設定方法が若干異なります。以下は、Claude for DesktopでこのMCP Serverを使用するための設定例です：

#### Claude for Desktop設定例

1. Claude for Desktopの設定（Settings）を開きます

    ![Claude Open Setting](/readme/img/opening_setting.png)

2. 開発者（Developer）タブに切り替え、設定編集（Edit Config）をクリックし、テキストエディタで設定ファイルを開きます

    ![Claude Edit Configuration](/readme/img/edit_config.png)

3. 「ログイン成功ページ」の設定情報をクライアントの設定ファイル`claude_desktop_config.json`に追加します

    ![Configuration Example](/readme/img/config_info.png)

#### 設定パラメータの説明

- `command`：ダウンロードまたはビルドした`aqara-mcp-server`実行ファイルへの完全パス
- `args`：`["run", "stdio"]`を使用してstdio転送モードを開始
- `env`：環境変数設定
  - `token`：Aqaraログインページから取得したアクセストークン
  - `region`：Aqaraアカウントの地域（CN、US、EUなど）

#### その他のMCPクライアント

MCPプロトコルをサポートする他のクライアント（ChatGPT、Cursorなど）の設定方法は類似しています：

- クライアントがMCPプロトコルをサポートしていることを確認
- 実行ファイルパスと起動パラメータを設定
- 環境変数`token`と`region`を設定
- 適切な転送プロトコルを選択（`stdio`推奨）

### サービス起動

#### 標準モード（推奨）

Claude for Desktopを再起動します。その後、自然言語でデバイス制御、デバイス照会、シーン実行などの操作を実行できます。

![Claude Chat Example](/readme/img/claude.png)

#### HTTPモード（オプション）

HTTPモードを使用する必要がある場合は、以下のように起動できます：

```bash
# デフォルトポート8080を使用
./aqara-mcp-server run http

# またはカスタムホストとポートを指定
./aqara-mcp-server run http --host localhost --port 9000
```

その後、クライアント設定で`["run", "http"]`パラメータを使用してください。

## APIツール説明

MCPクライアントは、これらのツールを呼び出してAqaraスマートホームデバイスと連携できます。

### デバイス制御類

#### device_control

スマートホームデバイスの状態や属性（オン/オフ、温度、明度、色、色温度など）を制御します。

**パラメータ：**

- `endpoint_ids` _(Array\<Integer\>, 必須)_：制御対象デバイスIDリスト
- `control_params` _(Object, 必須)_：制御パラメータオブジェクト、具体的な操作を含む：
  - `action` _(String, 必須)_：実行する操作（`"on"`, `"off"`, `"set"`, `"up"`, `"down"`, `"cooler"`, `"warmer"`など）
  - `attribute` _(String, 必須)_：制御するデバイス属性（`"on_off"`, `"brightness"`, `"color_temperature"`, `"ac_mode"`など）
  - `value` _(String | Number, オプション)_：目標値（actionが"set"の場合必須）
  - `unit` _(String, オプション)_：値の単位（`"%"`, `"K"`, `"℃"`など）

**戻り値：** デバイス制御の操作結果メッセージ

### デバイス照会類

#### device_query

指定した場所（部屋）とデバイスタイプに基づいてデバイスリストを取得します（リアルタイム状態情報は含まない）。

**パラメータ：**

- `positions` _(Array\<String\>, オプション)_：部屋名リスト。空配列は全部屋を照会
- `device_types` _(Array\<String\>, オプション)_：デバイスタイプリスト（`"Light"`, `"WindowCovering"`, `"AirConditioner"`, `"Button"`など）。空配列は全タイプを照会

**戻り値：** Markdown形式のデバイスリスト（デバイス名とIDを含む）

#### device_status_query

デバイスの現在の状態情報を取得します（色、明度、オン/オフなどのリアルタイム状態情報を照会するため）。

**パラメータ：**

- `positions` _(Array\<String\>, オプション)_：部屋名リスト。空配列は全部屋を照会
- `device_types` _(Array\<String\>, オプション)_：デバイスタイプリスト。選択可能値は`device_query`と同じ。空配列は全タイプを照会

**戻り値：** Markdown形式のデバイス状態情報

#### device_log_query

デバイスの履歴ログ情報を照会します。

**パラメータ：**

- `endpoint_ids` _(Array\<Integer\>, 必須)_：履歴記録を照会するデバイスIDリスト
- `start_datetime` _(String, オプション)_：照会開始時間、形式`YYYY-MM-DD HH:MM:SS`（例：`"2023-05-16 12:00:00"`）
- `end_datetime` _(String, オプション)_：照会終了時間、形式`YYYY-MM-DD HH:MM:SS`
- `attribute` _(String, オプション)_：照会する特定デバイス属性名（`on_off`, `brightness`など）。未提供時は記録されたすべての属性を照会

**戻り値：** Markdown形式のデバイス履歴状態情報

> 📝 **注意：** 現在の実装では"This feature will be available soon."と表示される場合があり、機能の完善待ちを示します。

### シーン管理類

#### get_scenes

ユーザー家庭内のすべてのシーン、または指定部屋内のシーンを照会します。

**パラメータ：**

- `positions` _(Array\<String\>, オプション)_：部屋名リスト。空配列は家庭全体のシーンを照会

**戻り値：** Markdown形式のシーン情報

#### run_scenes

シーンIDに基づいて指定シーンを実行します。

**パラメータ：**

- `scenes` _(Array\<Integer\>, 必須)_：実行するシーンIDリスト

**戻り値：** シーン実行の結果メッセージ

### 家庭管理類

#### get_homes

ユーザーアカウント下のすべての家庭リストを取得します。

**パラメータ：** なし

**戻り値：** コンマ区切りの家庭名リスト。データがない場合は空文字列または対応する提示情報を返します

#### switch_home

ユーザーが現在操作する家庭を切り替えます。切り替え後、以降のデバイス照会、制御などの操作は新しく切り替えた家庭に対して行われます。

**パラメータ：**

- `home_name` _(String, 必須)_：目標家庭の名前

**戻り値：** 切り替え操作の結果メッセージ

### 自動化設定類

#### automation_config

タイマーまたは遅延デバイス制御タスクを設定します（現在は定時遅延自動化設定のみサポート）。

**パラメータ：**

- `scheduled_time` _(String, 必須)_：設定時間点（遅延タスクの場合、現在時間点を基準に変換）、形式`YYYY-MM-DD HH:MM:SS`（例：`"2025-05-16 12:12:12"`）
- `endpoint_ids` _(Array\<Integer\>, 必須)_：タイマー制御するデバイスIDリスト
- `control_params` _(Object, 必須)_：デバイス制御パラメータ、`device_control`ツールと同じ形式（action、attribute、valueなどを含む）

**戻り値：** 自動化設定結果メッセージ

> 📝 **注意：** 現在の実装では"This feature will be available soon."と表示される場合があり、機能の完善待ちを示します。

## プロジェクト構成

### ディレクトリ構成

```text
.
├── cmd.go                # Cobra CLI コマンド定義とプログラムエントリポイント（main関数を含む）
├── server.go             # MCPサーバーコアロジック、ツール定義とリクエスト処理
├── smh.go                # Aqaraスマートホームプラットフォーム API インターフェース封装
├── middleware.go         # ミドルウェア：ユーザー認証、タイムアウト制御、例外復旧
├── config.go             # グローバル設定管理と環境変数処理
├── go.mod                # Go モジュール依存管理ファイル
├── go.sum                # Go モジュール依存チェックサムファイル
├── readme/               # README ドキュメントと画像リソース
│   ├── img/              # 画像リソースディレクトリ
│   └── *.md              # 多言語 README ファイル
├── LICENSE               # MIT オープンソースライセンス
└── README.md             # プロジェクトメインドキュメント
```

### コアファイル説明

- **`cmd.go`**：CobraフレームワークベースのCLI実装、`run stdio`と`run http`起動モード及びメインエントリ関数を定義
- **`server.go`**：MCPサーバーコア実装、ツール登録、リクエスト処理、プロトコルサポートを担当
- **`smh.go`**：Aqaraスマートホームプラットフォーム API封装層、デバイス制御、認証、マルチホームサポートを提供
- **`middleware.go`**：リクエスト処理ミドルウェア、認証検証、タイムアウト制御、例外処理を提供
- **`config.go`**：グローバル設定管理、環境変数処理とAPI設定を担当

## 開発ガイド

IssueやPull Requestの提出によるプロジェクト貢献を歓迎します！

コード提出前に以下を確認してください：

1. コードがGo言語のコーディング規範に従っている
2. 関連するMCPツールとインターフェース定義の一貫性と明確性を保持
3. 変更をカバーする単体テストの追加または更新
4. 必要に応じて関連ドキュメント（このREADMEなど）の更新
5. コミットメッセージが明確であることを確認

## ライセンス

本プロジェクトは[MIT License](/LICENSE)に基づいて認可されています。

Copyright (c) 2025 Aqara-Copilot
