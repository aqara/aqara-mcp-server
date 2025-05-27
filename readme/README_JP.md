# Aqara MCP Server

[English](/readme/README.md) | [中文](/readme/README_CN.md) | [繁體中文](/readme/README_CHT.md) | [Français](/readme/README_FR.md) | [한국어](/readme/README_KR.md) | [Español](/readme/README_ES.md) | 日本語 | [Deutsch](/readme/README_DE.md) | [Italiano](/readme/README_IT.md)

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/aqara/aqara-mcp-server)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![Release](https://img.shields.io/github/v/release/aqara/aqara-mcp-server)](https://github.com/aqara/aqara-mcp-server/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Aqara MCP Server は、[MCP (Model Context Protocol)](https://modelcontextprotocol.io/introduction) プロトコルに基づいて開発されたスマートホーム制御サービスです。MCPプロトコルをサポートするAIアシスタントやAPI（Claude、ChatGPT、Cursorなど）がAqaraスマートホームデバイスと対話し、自然言語でのデバイス制御、状態照会、シーン実行などの機能を実現できます。

## 目次

- [Aqara MCP Server](#aqara-mcp-server)
  - [目次](#目次)
  - [特徴](#特徴)
  - [動作原理](#動作原理)
  - [クイックスタート](#クイックスタート)
    - [前提条件](#前提条件)
    - [インストール](#インストール)
    - [Aqaraアカウント認証](#aqaraアカウント認証)
    - [設定例 (Claude for Desktop)](#設定例-claude-for-desktop)
    - [サービス実行](#サービス実行)
  - [利用可能なツール](#利用可能なツール)
    - [device\_control](#device_control)
    - [device\_query](#device_query)
    - [device\_status\_query](#device_status_query)
    - [device\_log\_query](#device_log_query)
    - [run\_scenes](#run_scenes)
    - [get\_scenes](#get_scenes)
    - [automation\_config](#automation_config)
    - [get\_homes](#get_homes)
    - [switch\_home](#switch_home)
  - [プロジェクト構造](#プロジェクト構造)
    - [コアファイル説明](#コアファイル説明)
  - [貢献ガイドライン](#貢献ガイドライン)
  - [ライセンス](#ライセンス)

## 特徴

-   **包括的なデバイス制御**：Aqaraスマートデバイスのオン/オフ、明度、色温度、モードなど様々な属性の精密制御をサポート。
-   **柔軟なデバイス照会**：部屋、デバイスタイプ別にデバイス一覧とその詳細状態を照会可能。
-   **インテリジェントシーン管理**：ユーザー設定のスマートホームシーンの照会と実行をサポート。
-   **デバイス履歴記録**：指定した時間範囲内のデバイス状態変更履歴を照会。
-   **自動化設定**：タイマーや遅延デバイス制御タスクの設定をサポート。
-   **マルチホーム対応**：ユーザーアカウント下の異なる家庭の照会と切り替えをサポート。
-   **MCPプロトコル互換**：MCPプロトコル仕様に完全準拠し、各種AIアシスタントとの統合が容易。
-   **安全な認証機構**：ログイン認証＋署名ベースの安全認証でユーザーデータとデバイスセキュリティを保護。
-   **クロスプラットフォーム実行**：Go言語で開発され、マルチプラットフォーム実行可能ファイルにコンパイル可能。
-   **拡張しやすい設計**：モジュラー設計により、新しいツールや機能の追加が容易。

## 動作原理

Aqara MCP ServerはAIアシスタントとAqaraスマートホームプラットフォーム間のブリッジとして機能します：

1.  **AIアシスタント (MCPクライアント)**：ユーザーがAIアシスタントを通じて指示を出します（例：「リビングのライトをつけて」）。
2.  **MCPクライアント**：ユーザー指示を解析し、MCPプロトコルに基づいてAqara MCP Serverの対応ツール（例：`device_control`）を呼び出します。
3.  **Aqara MCP Server (本プロジェクト)**：クライアントからのリクエストを受信し、検証後`smh.go`モジュールを呼び出します。
4.  **`smh.go`モジュール**：設定されたAqara認証情報を使用し、AqaraクラウドAPIと通信して実際のデバイス操作やデータ照会を実行。
5.  **レスポンスフロー**：AqaraクラウドAPIが結果を返し、Aqara MCP Server経由でMCPクライアントに渡され、最終的にユーザーに提示されます。

## クイックスタート

### 前提条件

-   Go（バージョン1.24以上）
-   Git（ソースからビルドする場合）
-   Aqaraアカウントと紐付けられたスマートデバイス

### インストール

プリコンパイルされた実行可能ファイルをダウンロードするか、ソースからビルドできます。

**オプション1: プリコンパイル版のダウンロード（推奨）**

以下のリンクから、お使いのオペレーティングシステム用の最新実行可能ファイルパッケージをダウンロードしてください。

[Releasesページ](https://github.com/aqara/aqara-mcp-server/releases)

解凍後すぐに使用できます。

**オプション2: ソースからビルド**

```bash
# リポジトリをクローン
git clone https://github.com/aqara/aqara-mcp-server.git
cd aqara-mcp-server

# 依存関係をダウンロード
go mod tidy

# 実行可能ファイルをビルド
go build -o aqara-mcp-server
```
ビルド完了後、現在のディレクトリに`aqara-mcp-server`実行可能ファイルが生成されます。

### Aqaraアカウント認証

MCP ServerがAqaraアカウントにアクセスしてデバイスを制御できるようにするため、最初にログイン認証を行う必要があります。

以下のアドレスにアクセスしてログイン認証を完了してください：
[https://cdn.aqara.com/app/mcpserver/login.html](https://cdn.aqara.com/app/mcpserver/login.html)

ログイン成功後、必要な認証情報（`token`、`region`など）が提供されます。これらの情報は後続の設定手順で使用されます。

**この情報は大切に保管し、特に`token`は他人に漏らさないでください。**

### 設定例 (Claude for Desktop)

MCPクライアントによって設定方法が若干異なります。以下はClaude for DesktopでこのMCP Serverを使用するための設定例です：

1.  Claude for Desktopの設定（Settings）を開きます。
2.  開発者（Developer）タブに切り替えます。
3.  設定編集（Edit Config）をクリックし、テキストエディタで設定ファイルを開きます。

    ![](/readme/img/setting0.png)
    ![](/readme/img/setting1.png)

4.  「ログイン成功ページ」の設定情報をクライアントの設定ファイル（claude_desktop_config.json）に追加します。設定例：

    ![](/readme/img/config.png)

**設定説明：**
- `command`: ダウンロードまたはビルドした`aqara-mcp-server`実行可能ファイルの完全パス
- `args`: `["run", "stdio"]`を使用してstdio転送モードを開始
- `env`: 環境変数設定
  - `token`: Aqaraログインページから取得したアクセストークン
  - `region`: Aqaraアカウントの地域（CN、US、EUなど）

### サービス実行

Claude for Desktopを再起動します。その後、対話を通じてMCP Serverが提供するツールを呼び出してデバイス制御、デバイス照会などの操作を実行できます。

![](/readme/img/claude.png)

**その他のMCPクライアント設定**

その他のMCPプロトコル対応クライアント（Claude、ChatGPT、Cursorなど）の設定方法も類似しています：
- クライアントがMCPプロトコルをサポートしていることを確認
- 実行可能ファイルパスと起動パラメータを設定
- 環境変数`token`と`region`を設定
- 適切な転送プロトコルを選択（`stdio`推奨）

**SSEモード（オプション）**

SSE（Server-Sent Events）モードを使用する必要がある場合は、以下のように起動できます：

```bash
# デフォルトポート8080を使用
./aqara-mcp-server run sse

# またはカスタムホストとポートを指定
./aqara-mcp-server run sse --host localhost --port 9000
```

その後、クライアント設定で`["run", "sse"]`パラメータを使用してください。

## 利用可能なツール

MCPクライアントはこれらのツールを呼び出してAqaraスマートホームデバイスと対話できます。

### device_control

-   **説明**: スマートホームデバイスの状態や属性を制御（オン/オフ、温度、明度、色、色温度など）。
-   **パラメータ**:
    -   `endpoint_ids` (Array<Integer>, 必須): 制御するデバイスIDのリスト。
    -   `control_params` (Object, 必須): 制御パラメータオブジェクト、具体的な操作を含む。
        -   `action` (String, 必須): 実行する操作。例：`"on"`, `"off"`, `"set"`, `"up"`, `"down"`, `"cooler"`, `"warmer"`。
        -   `attribute` (String, 必須): 制御するデバイス属性。例：`"on_off"`, `"brightness"`, `"color_temperature"`, `"ac_mode"`。
        -   `value` (String | Number, オプション): 目標値（actionが"set"の場合必須）。
        -   `unit` (String, オプション): 値の単位（例：`"%"`, `"K"`, `"℃"`）。
-   **戻り値**: (String) デバイス制御の操作結果メッセージ。

### device_query

-   **説明**: 指定した場所（部屋）とデバイスタイプに基づいてデバイス一覧を取得（リアルタイム状態情報は含まず、デバイスとそのIDのみをリスト）。
-   **パラメータ**:
    -   `positions` (Array<String>, オプション): 部屋名のリスト。空配列または未提供の場合、全ての部屋を照会。
    -   `device_types` (Array<String>, オプション): デバイスタイプのリスト。例：`"Light"`, `"WindowCovering"`, `"AirConditioner"`, `"Button"`など。空配列または未提供の場合、全てのタイプを照会。
-   **戻り値**: (String) Markdown形式のデバイス一覧、デバイス名とIDを含む。

### device_status_query

-   **説明**: デバイスの現在の状態情報を取得（色、明度、オン/オフなど状態関連属性の照会用）。
-   **パラメータ**:
    -   `positions` (Array<String>, オプション): 部屋名のリスト。空配列または未提供の場合、全ての部屋を照会。
    -   `device_types` (Array<String>, オプション): デバイスタイプのリスト。選択肢は`device_query`と同じ。空配列または未提供の場合、全てのタイプを照会。
-   **戻り値**: (String) Markdown形式のデバイス状態情報。

### device_log_query

-   **説明**: デバイスのログを照会。
-   **パラメータ**:
    -   `endpoint_ids` (Array<Integer>, 必須): 履歴記録を照会するデバイスIDのリスト。
    -   `start_datetime` (String, オプション): 照会開始時間、形式は`YYYY-MM-DD HH:MM:SS`（例：`"2023-05-16 12:00:00"`）。
    -   `end_datetime` (String, オプション): 照会終了時間、形式は`YYYY-MM-DD HH:MM:SS`。
    -   `attribute` (String, オプション): 照会する特定のデバイス属性名（例：`on_off`, `brightness`）。未提供の場合、該当デバイスの全ての記録済み属性の履歴記録を照会。
-   **戻り値**: (String) Markdown形式のデバイス履歴状態情報。（注意：現在の実装では「This feature will be available soon.」と表示される可能性があり、機能の完善待ちを示しています。）

### run_scenes

-   **説明**: シーンIDに基づいて指定シーンを実行。
-   **パラメータ**:
    -   `scenes` (Array<Integer>, 必須): 実行するシーンIDのリスト。
-   **戻り値**: (String) シーン実行の結果メッセージ。

### get_scenes

-   **説明**: ユーザー家庭下の全シーン、または指定部屋内のシーンを照会。
-   **パラメータ**:
    -   `positions` (Array<String>, オプション): 部屋名のリスト。空配列または未提供の場合、家庭全体のシーンを照会。
-   **戻り値**: (String) Markdown形式のシーン情報。

### automation_config

-   **説明**: タイマーまたは遅延デバイス制御タスクを設定。
-   **パラメータ**:
    -   `scheduled_time` (String, 必須): 設定時間（遅延タスクの場合は現在時刻を基準に変換）、形式は`YYYY-MM-DD HH:MM:SS`（例：`"2025-05-16 12:12:12"`）。
    -   `endpoint_ids` (Array<Integer>, 必須): タイマー制御するデバイスIDのリスト。
    -   `control_params` (Object, 必須): デバイス制御パラメータ、`device_control`ツールと同じ形式を使用（action、attribute、valueなどを含む）。
-   **戻り値**: (String) 自動化設定結果メッセージ。

### get_homes

-   **説明**: ユーザーアカウント下の全家庭リストを取得。
-   **パラメータ**: なし。
-   **戻り値**: (String) カンマ区切りの家庭名リスト。データがない場合は空文字列または対応する案内情報を返す。

### switch_home

-   **説明**: ユーザーが現在操作中の家庭を切り替え。切り替え後、後続のデバイス照会、制御などの操作は新しく切り替えた家庭が対象となります。
-   **パラメータ**:
    -   `home_name` (String, 必須): 対象家庭の名前（`get_homes`ツールから提供される利用可能リストから選択）。
-   **戻り値**: (String) 切り替え操作の結果メッセージ。

## プロジェクト構造

```
.
├── cmd.go                # Cobra CLI コマンド定義とプログラムエントリーポイント（main関数を含む）
├── server.go             # MCP サーバーコアロジック、ツール定義とリクエスト処理
├── smh.go                # Aqara スマートホームプラットフォーム API インターフェース封装
├── middleware.go         # ミドルウェア：ユーザー認証、タイムアウト制御、例外復旧
├── config.go             # グローバル設定管理と環境変数処理
├── go.mod                # Go モジュール依存関係管理ファイル
├── go.sum                # Go モジュール依存関係チェックサムファイル
├── img/                  # README ドキュメントで使用される画像リソース
├── LICENSE               # MIT オープンソースライセンス
└── README.md             # プロジェクトドキュメント
```

### コアファイル説明

-   **`cmd.go`**: CobraフレームワークベースのCLI実装、`run stdio`と`run sse`起動モードおよびメインエントリー関数を定義
-   **`server.go`**: MCPサーバーコア実装、ツール登録、リクエスト処理、プロトコルサポートを担当
-   **`smh.go`**: Aqaraスマートホームプラットフォーム API封装レイヤー、デバイス制御、認証、マルチホームサポートを提供
-   **`middleware.go`**: リクエスト処理ミドルウェア、認証検証、タイムアウト制御、例外処理を提供
-   **`config.go`**: グローバル設定管理、環境変数処理とAPI設定を担当

## 貢献ガイドライン

IssueやPull Requestの提出を通じてプロジェクトへの貢献を歓迎します！

コード提出前に以下を確認してください：
1.  コードがGo言語のコーディング規範に従っている。
2.  関連するMCPツールとプロンプトインターフェース定義の一貫性と明確性を保つ。
3.  変更をカバーする単体テストを追加または更新する。
4.  必要に応じて関連ドキュメント（このREADMEなど）を更新する。
5.  コミットメッセージが明確である。

## ライセンス

本プロジェクトは[MIT License](/LICENSE)に基づいてライセンスされています。
Copyright (c) 2025 Aqara-Copliot