<div align="center" style="display: flex; align-items: center; justify-content: center; ">

  <img src="/readme/img/logo.png" alt="Aqara Logo" height="120">
  <h1>Aqara MCP Server</h1>

</div>

<div align="center">

[English](/readme/README.md) | [简体中文](/readme/README_CN.md) | 繁體中文 | [Français](/readme/README_FR.md) | [한국어](/readme/README_KR.md) | [Español](/readme/README_ES.md) | [日本語](/readme/README_JP.md) | [Deutsch](/readme/README_DE.md) | [Italiano](/readme/README_IT.md)

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/aqara/aqara-mcp-server)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![Release](https://img.shields.io/github/v/release/aqara/aqara-mcp-server)](https://github.com/aqara/aqara-mcp-server/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

Aqara MCP Server 是基於 [MCP (Model Context Protocol)](https://modelcontextprotocol.io/introduction) 協定開發的企業級智慧家居控制服務。本系統讓任何支援 MCP 協定的 AI 助手或 API（如 Claude、Cursor 等）能夠與您的 Aqara 智慧家居裝置進行無縫整合，實現透過自然語言進行裝置控制、狀態查詢、場景執行等進階功能。

## 目錄

- [目錄](#目錄)
- [核心特性](#核心特性)
- [架構原理](#架構原理)
- [快速部署](#快速部署)
  - [系統需求](#系統需求)
  - [安裝方式](#安裝方式)
    - [方式一：下載預編譯版本（推薦）](#方式一下載預編譯版本推薦)
    - [方式二：從原始碼建置](#方式二從原始碼建置)
  - [Aqara 帳戶驗證](#aqara-帳戶驗證)
  - [客戶端整合](#客戶端整合)
    - [Claude for Desktop 配置範例](#claude-for-desktop-配置範例)
    - [配置參數詳解](#配置參數詳解)
    - [其他 MCP 客戶端](#其他-mcp-客戶端)
  - [服務啟動](#服務啟動)
    - [標準模式（推薦）](#標準模式推薦)
    - [HTTP 模式（即將推出）](#http-模式即將推出)
- [API 工具參考](#api-工具參考)
  - [裝置控制類](#裝置控制類)
    - [device\_control](#device_control)
  - [裝置查詢類](#裝置查詢類)
    - [device\_query](#device_query)
    - [device\_status\_query](#device_status_query)
    - [device\_log\_query](#device_log_query)
  - [場景管理類](#場景管理類)
    - [get\_scenes](#get_scenes)
    - [run\_scenes](#run_scenes)
  - [家庭管理類](#家庭管理類)
    - [get\_homes](#get_homes)
    - [switch\_home](#switch_home)
  - [自動化配置類](#自動化配置類)
    - [automation\_config](#automation_config)
- [專案架構](#專案架構)
  - [目錄結構](#目錄結構)
  - [核心檔案說明](#核心檔案說明)
- [開發指南](#開發指南)
- [授權條款](#授權條款)

## 核心特性

- ✨ **全方位裝置控制**：支援對 Aqara 智慧裝置的開關、亮度、色溫、運作模式等多重屬性進行精密控制
- 🔍 **智慧型裝置查詢**：提供按房間、裝置類型的靈活查詢功能，並支援詳細狀態資訊檢索
- 🎬 **智能場景管理**：完整支援使用者預設智慧家居場景的查詢與執行功能
- 📈 **歷史資料分析**：提供指定時間範圍內的裝置狀態變更歷史記錄查詢
- ⏰ **自動化工作流程**：支援配置定時或延遲裝置控制任務的自動化系統
- 🏠 **多家庭環境支援**：完整支援使用者帳戶下多個家庭環境的查詢與切換
- 🔌 **MCP 協定相容性**：完全符合 MCP 協定規範，確保與各類 AI 助手的無縫整合
- 🔐 **企業級安全機制**：採用基於登入授權與數位簽章的多層安全驗證，全面保護使用者數據與裝置安全
- 🌐 **跨平台相容性**：基於 Go 語言開發，支援編譯為多平台執行檔案
- 🔧 **高度可擴展性**：採用模組化設計架構，便於新增工具與功能擴展

## 架構原理

Aqara MCP Server 作為 AI 助手與 Aqara 智慧家居平台之間的企業級橋接服務：

```mermaid
graph LR
    A[AI 助手] --> B[MCP 客戶端]
    B --> C[Aqara MCP Server]
    C --> D[Aqara 雲端 API]
    D --> E[智慧裝置]
```

1. **AI 助手**：使用者透過 AI 助手發出自然語言指令（例如："請打開客廳的照明設備"）
2. **MCP 客戶端**：智慧解析使用者指令，依據 MCP 協定規範呼叫 Aqara MCP Server 提供的對應工具（如 `device_control`）
3. **Aqara MCP Server (本系統)**：接收來自客戶端的請求，使用已配置的 Aqara 憑證，與 Aqara 雲端 API 進行安全通訊，執行實際的裝置操作或資料查詢
4. **回應流程**：Aqara 雲端 API 回傳執行結果，經由 Aqara MCP Server 安全傳遞回 MCP 客戶端，最終呈現給使用者

## 快速部署

### 系統需求

- **Go** (版本 1.24 或更高版本) - 僅在從原始碼建置時需要
- **Git** (用於原始碼建置) - 選用
- **Aqara 帳戶**及已註冊的智慧裝置
- **支援 MCP 協定的客戶端** (如 Claude for Desktop、Cursor 等)

### 安裝方式

您可以選擇下載預編譯的執行檔案或從原始碼進行建置。

#### 方式一：下載預編譯版本（推薦）

造訪 GitHub Releases 頁面，下載適用於您作業系統的最新執行檔案：

**📥 [前往 Releases 頁面下載](https://github.com/aqara/aqara-mcp-server/releases)**

下載對應平台的壓縮套件後解壓縮即可使用。

#### 方式二：從原始碼建置

```bash
# 複製存儲庫
git clone https://github.com/aqara/aqara-mcp-server.git
cd aqara-mcp-server

# 下載相依性套件
go mod tidy

# 建置執行檔案
go build -o aqara-mcp-server
```

建置完成後，將在目前目錄下產生 `aqara-mcp-server` 執行檔案。

### Aqara 帳戶驗證

為了讓 MCP Server 能夠存取您的 Aqara 帳戶並控制裝置，您需要先完成登入授權程序。

請造訪以下網址完成登入授權：
**🔗 [https://cdn.aqara.com/app/mcpserver/login.html](https://cdn.aqara.com/app/mcpserver/login.html)**

登入成功後，您將取得必要的驗證資訊（如 `token`、`region`），這些資訊將在後續配置步驟中使用。

> ⚠️ **安全提醒**：請妥善保管 `token` 資訊，避免洩露給第三方。

### 客戶端整合

不同的 MCP 客戶端的配置方法略有差異。以下說明如何配置 Claude for Desktop 以使用此 MCP Server：

#### Claude for Desktop 配置範例

1. **開啟 Claude for Desktop 的設定 (Settings)**

    ![Claude Open Setting](/readme/img/opening_setting.png)

2. **切換至開發者 (Developer) 分頁，然後點選編輯配置 (Edit Config)，使用文字編輯器開啟配置檔案**

    ![Claude Edit Configuration](/readme/img/edit_config.png)

3. **將「登入成功頁面」的配置資訊，新增至客戶端的配置檔案 `claude_desktop_config.json` 中**

    ```json
    {
      "mcpServers": {
        "aqara": {
          "command": "/path/to/aqara-mcp-server",
          "args": ["run", "stdio"],
          "env": {
            "token": "your_token_here",
            "region": "your_region_here"
          }
        }
      }
    }
    ```

    ![Configuration Example](/readme/img/config_info.png)

#### 配置參數詳解

- `command`: 指向您下載或建置的 `aqara-mcp-server` 執行檔案的完整路徑
- `args`: 使用 `["run", "stdio"]` 啟動 stdio 傳輸模式
- `env`: 環境變數配置
  - `token`: 從 Aqara 登入頁面取得的存取權杖
  - `region`: 您的 Aqara 帳戶所在區域（支援的區域：CN、US、EU、KR、SG、RU）

#### 其他 MCP 客戶端

對於其他支援 MCP 協定的客戶端（如 ChatGPT、Cursor 等），配置方式類似：

- 確保客戶端支援 MCP 協定
- 配置執行檔案路徑和啟動參數
- 設定環境變數 `token` 和 `region`
- 選擇適當的傳輸協定（建議使用 `stdio`）

### 服務啟動

#### 標準模式（推薦）

重新啟動 Claude for Desktop。接著即可透過自然語言執行裝置控制、裝置查詢、場景執行等操作。

範例對話：

- "請打開客廳的照明設備"
- "將臥室空調設定為製冷模式，溫度設為 24 度"
- "查看所有房間的裝置清單"
- "執行晚安場景"

![Claude Chat Example](/readme/img/claude.png)

#### HTTP 模式（即將推出）

## API 工具參考

MCP 客戶端可透過呼叫這些工具與 Aqara 智慧家居裝置進行互動。

### 裝置控制類

#### device_control

控制智慧家居裝置的狀態或屬性（例如開關、溫度、亮度、色彩、色溫等）。

**參數：**

- `endpoint_ids` _(Array\<Integer\>, 必需)_：需要控制的裝置 ID 清單
- `control_params` _(Object, 必需)_：控制參數物件，包含具體操作：
  - `action` _(String, 必需)_：要執行的操作（如 `"on"`, `"off"`, `"set"`, `"up"`, `"down"`, `"cooler"`, `"warmer"`）
  - `attribute` _(String, 必需)_：要控制的裝置屬性（如 `"on_off"`, `"brightness"`, `"color_temperature"`, `"ac_mode"`）
  - `value` _(String | Number, 選用)_：目標值（當 action 為 "set" 時必需）
  - `unit` _(String, 選用)_：值的單位（如 `"%"`, `"K"`, `"℃"`）

**回傳：** 裝置控制的操作結果訊息

### 裝置查詢類

#### device_query

根據指定的位置（房間）和裝置類型取得裝置清單（不包含即時狀態資訊）。

**參數：**

- `positions` _(Array\<String\>, 選用)_：房間名稱清單。空陣列表示查詢所有房間
- `device_types` _(Array\<String\>, 選用)_：裝置類型清單（如 `"Light"`, `"WindowCovering"`, `"AirConditioner"`, `"Button"`）。空陣列表示查詢所有類型

**回傳：** Markdown 格式的裝置清單，包含裝置名稱和 ID

#### device_status_query

取得裝置的目前狀態資訊（用於查詢色彩、亮度、開關等即時狀態資訊）。

**參數：**

- `positions` _(Array\<String\>, 選用)_：房間名稱清單。空陣列表示查詢所有房間
- `device_types` _(Array\<String\>, 選用)_：裝置類型清單。可選值同 `device_query`。空陣列表示查詢所有類型

**回傳：** Markdown 格式的裝置狀態資訊

#### device_log_query

查詢裝置的歷史記錄資訊。

**參數：**

- `endpoint_ids` _(Array\<Integer\>, 必需)_：需要查詢歷史記錄的裝置 ID 清單
- `start_datetime` _(String, 選用)_：查詢起始時間，格式為 `YYYY-MM-DD HH:MM:SS`（例如：`"2023-05-16 12:00:00"`）
- `end_datetime` _(String, 選用)_：查詢結束時間，格式為 `YYYY-MM-DD HH:MM:SS`
- `attributes` _(Array\<String\>, 選用)_：要查詢的裝置屬性名稱清單（如 `["on_off", "brightness"]`）。未提供時查詢所有已記錄屬性

**回傳：** Markdown 格式的裝置歷史狀態資訊

### 場景管理類

#### get_scenes

查詢使用者家庭下所有場景，或指定房間內的場景。

**參數：**

- `positions` _(Array\<String\>, 選用)_：房間名稱清單。空陣列表示查詢整個家庭的場景

**回傳：** Markdown 格式的場景資訊

#### run_scenes

根據場景 ID 執行指定的場景。

**參數：**

- `scenes` _(Array\<Integer\>, 必需)_：需要執行的場景 ID 清單

**回傳：** 場景執行的結果訊息

### 家庭管理類

#### get_homes

取得使用者帳戶下的所有家庭清單。

**參數：** 無

**回傳：** 以逗號分隔的家庭名稱清單。如果無資料則回傳空字串或相應的提示資訊

#### switch_home

切換使用者目前操作的家庭。切換後，後續的裝置查詢、控制等操作將針對新切換的家庭。

**參數：**

- `home_name` _(String, 必需)_：目標家庭的名稱

**回傳：** 切換操作的結果訊息

### 自動化配置類

#### automation_config

自動化配置（目前僅支援定時或延遲裝置控制任務）。

**參數：**

- `scheduled_time` _(String, 必需)_：定時執行的時間點，使用標準 Crontab 格式 `"分 時 日 月 週"`。例如：`"30 14 * * *"`（每天14:30執行）、`"0 9 * * 1"`（每週一9:00執行）
- `endpoint_ids` _(Array\<Integer\>, 必需)_：需要定時控制的裝置 ID 清單
- `control_params` _(Object, 必需)_：裝置控制參數，使用與 `device_control` 工具相同的格式（包含 action、attribute、value 等）
- `task_name` _(String, 必需)_：此自動化任務的名稱或描述（用於識別和管理）
- `execution_once` _(Boolean, 選用)_：是否只執行一次
  - `true`：僅在指定時間執行一次任務（預設值）
  - `false`：週期性重複執行任務（如每天、每週等）

**回傳：** 自動化配置結果訊息

## 專案架構

### 目錄結構

```text
.
├── cmd.go                # Cobra CLI 命令定義和程式進入點（包含 main 函數）
├── server.go             # MCP 伺服器核心邏輯，工具定義和請求處理
├── smh.go                # Aqara 智慧家居平台 API 介面封裝
├── middleware.go         # 中介軟體：使用者驗證、逾時控制、異常復原
├── config.go             # 全域配置管理和環境變數處理
├── go.mod                # Go 模組相依性管理檔案
├── go.sum                # Go 模組相依性校驗和檔案
├── readme/               # README 文件和圖片資源
│   ├── img/              # 圖片資源目錄
│   └── *.md              # 多語言 README 檔案
├── LICENSE               # MIT 開源授權條款
└── README.md             # 專案主要文件
```

### 核心檔案說明

- **`cmd.go`**：基於 Cobra 框架的 CLI 實作，定義 `run stdio` 和 `run http` 啟動模式及主要進入函數
- **`server.go`**：MCP 伺服器核心實作，負責工具註冊、請求處理和協定支援
- **`smh.go`**：Aqara 智慧家居平台 API 封裝層，提供裝置控制、驗證和多家庭支援
- **`middleware.go`**：請求處理中介軟體，提供驗證確認、逾時控制和異常處理
- **`config.go`**：全域配置管理，負責環境變數處理和 API 配置

## 開發指南

我們歡迎透過提交 Issue 或 Pull Request 來參與專案貢獻！

在提交程式碼前，請確保：

1. 程式碼遵循 Go 語言的編碼規範
2. 相關的 MCP 工具和介面定義保持一致性和清晰性
3. 新增或更新單元測試以涵蓋您的變更
4. 如有必要，更新相關的文件（如本 README）
5. 確保您的提交訊息清晰明瞭

**🌟 如果這個專案對您有幫助，請給我們一個 Star！**

**🤝 歡迎加入我們的社群，一起讓智慧家居更智慧！**

## 授權條款

本專案基於 [MIT License](/LICENSE) 授權。

---

Copyright (c) 2025 Aqara-Copilot
