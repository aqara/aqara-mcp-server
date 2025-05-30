<div align="center" style="display: flex; align-items: center; justify-content: center; ">

  <img src="/readme/img/logo.png" alt="Aqara Logo" height="120">
  <h1>MCP Server</h1>

</div>

<div align="center">

[English](/readme/README.md) | [中文](/readme/README_CN.md) | 繁體中文 | [Français](/readme/README_FR.md) | [한국어](/readme/README_KR.md) | [Español](/readme/README_ES.md) | [日本語](/readme/README_JP.md) | [Deutsch](/readme/README_DE.md) | [Italiano](/readme/README_IT.md)

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/aqara/aqara-mcp-server)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![Release](https://img.shields.io/github/v/release/aqara/aqara-mcp-server)](https://github.com/aqara/aqara-mcp-server/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

Aqara MCP Server 是一個基於 [MCP (Model Context Protocol)](https://modelcontextprotocol.io/introduction) 協定開發的智慧家居控制服務。它允許任何支援 MCP 協定的 AI 助手或 API（例如 Claude、Cursor 等）與您的 Aqara 智慧家居裝置進行互動，實現透過自然語言控制裝置、查詢狀態、執行場景等功能。

## 目錄

- [目錄](#目錄)
- [特色](#特色)
- [工作原理](#工作原理)
- [快速開始](#快速開始)
  - [先決條件](#先決條件)
  - [安裝](#安裝)
    - [方式一：下載預編譯版本（建議）](#方式一下載預編譯版本建議)
    - [方式二：從原始碼建置](#方式二從原始碼建置)
  - [Aqara 帳戶認證](#aqara-帳戶認證)
  - [客戶端配置](#客戶端配置)
    - [Claude for Desktop 配置範例](#claude-for-desktop-配置範例)
    - [配置參數說明](#配置參數說明)
    - [其他 MCP 客戶端](#其他-mcp-客戶端)
  - [啟動服務](#啟動服務)
    - [標準模式（建議）](#標準模式建議)
    - [HTTP 模式（可選）](#http-模式可選)
- [API 工具說明](#api-工具說明)
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
- [專案結構](#專案結構)
  - [目錄結構](#目錄結構)
  - [核心檔案說明](#核心檔案說明)
- [開發指南](#開發指南)
- [授權條款](#授權條款)

## 特色

- **全面的裝置控制**：支援對 Aqara 智慧裝置的開關、亮度、色溫、模式等多種屬性進行精細控制
- **靈活的裝置查詢**：能夠按房間、裝置類型查詢裝置清單及其詳細狀態
- **智慧場景管理**：支援查詢和執行使用者預設的智慧家居場景
- **裝置歷史記錄**：查詢裝置在指定時間範圍內的歷史狀態變更記錄
- **自動化配置**：支援配置定時或延時裝置控制任務
- **多家庭支援**：支援查詢和切換使用者帳戶下的不同家庭
- **MCP 協定相容**：完全遵循 MCP 協定規範，易於與各類 AI 助手整合
- **安全認證機制**：採用基於登入授權+簽名的安全認證，保護使用者資料和裝置安全
- **跨平台執行**：基於 Go 語言開發，可編譯為多平台可執行檔
- **易於擴充**：模組化設計，可以方便地新增新的工具和功能

## 工作原理

Aqara MCP Server 作為 AI 助手與 Aqara 智慧家居平台之間的橋樑：

1. **AI 助手 (MCP 客戶端)**：使用者透過 AI 助手發出指令（例如，「打開客廳的燈」）
2. **MCP 客戶端**：將使用者指令解析，並根據 MCP 協定呼叫 Aqara MCP Server 提供的相應工具（例如 `device_control`）
3. **Aqara MCP Server (本專案)**：接收來自客戶端的請求，驗證後呼叫 `smh.go` 模組
4. **`smh.go` 模組**：使用配置好的 Aqara 憑證，與 Aqara 雲端 API 進行通訊，執行實際的裝置操作或資料查詢
5. **回應流程**：Aqara 雲端 API 回傳結果，經由 Aqara MCP Server 傳遞回 MCP 客戶端，最終呈現給使用者

## 快速開始

### 先決條件

- Go (版本 1.24 或更高)
- Git (用於從原始碼建置)
- Aqara 帳戶及已綁定的智慧裝置

### 安裝

您可以選擇下載預編譯的可執行檔或從原始碼建置。

#### 方式一：下載預編譯版本（建議）

造訪 GitHub Releases 頁面，下載適用於您作業系統的最新可執行檔：

**📥 [前往 Releases 頁面下載](https://github.com/aqara/aqara-mcp-server/releases)**

下載對應平台的壓縮套件後解壓縮即可使用。

#### 方式二：從原始碼建置

```bash
# 克隆儲存庫
git clone https://github.com/aqara/aqara-mcp-server.git
cd aqara-mcp-server

# 下載依賴套件
go mod tidy

# 建置可執行檔
go build -o aqara-mcp-server
```

建置完成後，會在目前目錄下產生 `aqara-mcp-server` 可執行檔。

### Aqara 帳戶認證

為了使 MCP Server 能夠存取您的 Aqara 帳戶並控制裝置，您需要先進行登入授權。

請造訪以下網址完成登入授權：
**🔗 [https://cdn.aqara.com/app/mcpserver/login.html](https://cdn.aqara.com/app/mcpserver/login.html)**

登入成功後，您將取得必要的認證資訊（如 `token`, `region`），這些資訊將在後續配置步驟中使用。

> ⚠️ **安全提醒**：請妥善保管 `token` 資訊，不要洩露給他人。

### 客戶端配置

不同的 MCP 客戶端的配置方法略有不同。以下是如何配置 Claude for Desktop 以使用此 MCP Server 的範例：

#### Claude for Desktop 配置範例

1. 開啟 Claude for Desktop 的設定 (Settings)

    ![Claude Open Setting](/readme/img/opening_setting.png)

2. 切換到開發者 (Developer) 標籤頁，然後點擊編輯配置 (Edit Config)，使用文字編輯器開啟配置檔

    ![Claude Edit Configuration](/readme/img/edit_config.png)

3. 將「登入成功頁面」的配置資訊，新增到客戶端的配置檔 `claude_desktop_config.json` 中

    ![Configuration Example](/readme/img/config_info.png)

#### 配置參數說明

- `command`: 指向您下載或建置的 `aqara-mcp-server` 可執行檔的完整路徑
- `args`: 使用 `["run", "stdio"]` 啟動 stdio 傳輸模式
- `env`: 環境變數配置
  - `token`: 從 Aqara 登入頁面取得的存取權杖
  - `region`: 您的 Aqara 帳戶所在區域（如 CN、US、EU 等）

#### 其他 MCP 客戶端

對於其他支援 MCP 協定的客戶端（如 ChatGPT、Cursor 等），配置方式類似：

- 確保客戶端支援 MCP 協定
- 配置可執行檔路徑和啟動參數
- 設定環境變數 `token` 和 `region`
- 選擇合適的傳輸協定（建議使用 `stdio`）

### 啟動服務

#### 標準模式（建議）

重新啟動 Claude for Desktop。然後就可以透過自然語言執行裝置控制、裝置查詢、場景執行等操作。

![Claude Chat Example](/readme/img/claude.png)

#### HTTP 模式（可選）

如果您需要使用 HTTP 模式，可以這樣啟動：

```bash
# 使用預設埠號 8080
./aqara-mcp-server run http

# 或指定自訂主機和埠號
./aqara-mcp-server run http --host localhost --port 9000
```

然後在客戶端配置中使用 `["run", "http"]` 參數。

## API 工具說明

MCP 客戶端可以透過呼叫這些工具與 Aqara 智慧家居裝置進行互動。

### 裝置控制類

#### device_control

控制智慧家居裝置的狀態或屬性（例如開關、溫度、亮度、顏色、色溫等）。

**參數：**

- `endpoint_ids` _(Array\<Integer\>, 必需)_：需要控制的裝置 ID 清單
- `control_params` _(Object, 必需)_：控制參數物件，包含具體操作：
  - `action` _(String, 必需)_：要執行的操作（如 `"on"`, `"off"`, `"set"`, `"up"`, `"down"`, `"cooler"`, `"warmer"`）
  - `attribute` _(String, 必需)_：要控制的裝置屬性（如 `"on_off"`, `"brightness"`, `"color_temperature"`, `"ac_mode"`）
  - `value` _(String | Number, 可選)_：目標值（當 action 為 "set" 時必需）
  - `unit` _(String, 可選)_：值的單位（如 `"%"`, `"K"`, `"℃"`）

**回傳值：** 裝置控制的操作結果訊息

### 裝置查詢類

#### device_query

根據指定的位置（房間）和裝置類型取得裝置清單（不包含即時狀態資訊）。

**參數：**

- `positions` _(Array\<String\>, 可選)_：房間名稱清單。空陣列表示查詢所有房間
- `device_types` _(Array\<String\>, 可選)_：裝置類型清單（如 `"Light"`, `"WindowCovering"`, `"AirConditioner"`, `"Button"`）。空陣列表示查詢所有類型

**回傳值：** Markdown 格式的裝置清單，包含裝置名稱和 ID

#### device_status_query

取得裝置的目前狀態資訊（用於查詢顏色、亮度、開關等即時狀態資訊）。

**參數：**

- `positions` _(Array\<String\>, 可選)_：房間名稱清單。空陣列表示查詢所有房間
- `device_types` _(Array\<String\>, 可選)_：裝置類型清單。可選值同 `device_query`。空陣列表示查詢所有類型

**回傳值：** Markdown 格式的裝置狀態資訊

#### device_log_query

查詢裝置的歷史日誌資訊。

**參數：**

- `endpoint_ids` _(Array\<Integer\>, 必需)_：需要查詢歷史記錄的裝置 ID 清單
- `start_datetime` _(String, 可選)_：查詢起始時間，格式為 `YYYY-MM-DD HH:MM:SS`（例如：`"2023-05-16 12:00:00"`）
- `end_datetime` _(String, 可選)_：查詢結束時間，格式為 `YYYY-MM-DD HH:MM:SS`
- `attribute` _(String, 可選)_：要查詢的特定裝置屬性名稱（如 `on_off`, `brightness`）。未提供時查詢所有已記錄屬性

**回傳值：** Markdown 格式的裝置歷史狀態資訊

> 📝 **注意：** 目前實作可能提示 "This feature will be available soon."，表示功能待完善。

### 場景管理類

#### get_scenes

查詢使用者家庭下所有場景，或指定房間內的場景。

**參數：**

- `positions` _(Array\<String\>, 可選)_：房間名稱清單。空陣列表示查詢整個家庭的場景

**回傳值：** Markdown 格式的場景資訊

#### run_scenes

根據場景 ID 執行指定的場景。

**參數：**

- `scenes` _(Array\<Integer\>, 必需)_：需要執行的場景 ID 清單

**回傳值：** 場景執行的結果訊息

### 家庭管理類

#### get_homes

取得使用者帳戶下的所有家庭清單。

**參數：** 無

**回傳值：** 以逗號分隔的家庭名稱清單。如果無資料則回傳空字串或相應的提示資訊

#### switch_home

切換使用者目前操作的家庭。切換後，後續的裝置查詢、控制等操作將針對新切換的家庭。

**參數：**

- `home_name` _(String, 必需)_：目標家庭的名稱

**回傳值：** 切換操作的結果訊息

### 自動化配置類

#### automation_config

配置定時或延時裝置控制任務（目前僅支援定延時自動化配置）。

**參數：**

- `scheduled_time` _(String, 必需)_：設定的時間點（如果是延時任務，基於目前時間點轉化），格式為 `YYYY-MM-DD HH:MM:SS`（例如：`"2025-05-16 12:12:12"`）
- `endpoint_ids` _(Array\<Integer\>, 必需)_：需要定時控制的裝置 ID 清單
- `control_params` _(Object, 必需)_：裝置控制參數，使用與 `device_control` 工具相同的格式（包含 action、attribute、value 等）

**回傳值：** 自動化配置結果訊息

> 📝 **注意：** 目前實作可能提示 "This feature will be available soon."，表示功能待完善。

## 專案結構

### 目錄結構

```text
.
├── cmd.go                # Cobra CLI 命令定義和程式進入點（包含 main 函式）
├── server.go             # MCP 伺服器核心邏輯，工具定義和請求處理
├── smh.go                # Aqara 智慧家居平台 API 介面封裝
├── middleware.go         # 中介軟體：使用者認證、逾時控制、異常復原
├── config.go             # 全域配置管理和環境變數處理
├── go.mod                # Go 模組依賴管理檔案
├── go.sum                # Go 模組依賴校驗和檔案
├── readme/               # README 文件和圖片資源
│   ├── img/              # 圖片資源目錄
│   └── *.md              # 多語言 README 檔案
├── LICENSE               # MIT 開源授權條款
└── README.md             # 專案主文件
```

### 核心檔案說明

- **`cmd.go`**：基於 Cobra 框架的 CLI 實作，定義 `run stdio` 和 `run http` 啟動模式及主進入函式
- **`server.go`**：MCP 伺服器核心實作，負責工具註冊、請求處理和協定支援
- **`smh.go`**：Aqara 智慧家居平台 API 封裝層，提供裝置控制、認證和多家庭支援
- **`middleware.go`**：請求處理中介軟體，提供認證驗證、逾時控制和異常處理
- **`config.go`**：全域配置管理，負責環境變數處理和 API 配置

## 開發指南

歡迎透過提交 Issue 或 Pull Request 來參與專案貢獻！

在提交程式碼前，請確保：

1. 程式碼遵循 Go 語言的編碼規範
2. 相關的 MCP 工具和介面定義保持一致性和清晰性
3. 新增或更新單元測試以涵蓋您的變更
4. 如有必要，更新相關的文件（如本 README）
5. 確保您的提交訊息清晰明瞭

## 授權條款

本專案基於 [MIT License](/LICENSE) 授權。

Copyright (c) 2025 Aqara-Copilot
