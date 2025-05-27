# Aqara MCP Server

[English](/readme/README.md) | [中文](/readme/README_CN.md) | 繁體中文 | [Français](/readme/README_FR.md) | [한국어](/readme/README_KR.md) | [Español](/readme/README_ES.md) | [日本語](/readme/README_JP.md) | [Deutsch](/readme/README_DE.md) | [Italiano](/readme/README_IT.md)

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/aqara/aqara-mcp-server)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![Release](https://img.shields.io/github/v/release/aqara/aqara-mcp-server)](https://github.com/aqara/aqara-mcp-server/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Aqara MCP Server 是一個基於 [MCP (Model Context Protocol)](https://modelcontextprotocol.io/introduction) 協議開發的智慧家居控制服務。它允許任何支援 MCP 協議的 AI 助手或 API（例如 Claude、ChatGPT、Cursor 等）與您的 Aqara 智慧家居設備進行互動，實現透過自然語言控制設備、查詢狀態、執行場景等功能。

## 目錄

- [Aqara MCP Server](#aqara-mcp-server)
  - [目錄](#目錄)
  - [特色](#特色)
  - [工作原理](#工作原理)
  - [快速開始](#快速開始)
    - [先決條件](#先決條件)
    - [安裝](#安裝)
    - [Aqara 帳戶認證](#aqara-帳戶認證)
    - [配置範例 (Claude for Desktop)](#配置範例-claude-for-desktop)
    - [執行服務](#執行服務)
  - [可用工具](#可用工具)
    - [device\_control](#device_control)
    - [device\_query](#device_query)
    - [device\_status\_query](#device_status_query)
    - [device\_log\_query](#device_log_query)
    - [run\_scenes](#run_scenes)
    - [get\_scenes](#get_scenes)
    - [automation\_config](#automation_config)
    - [get\_homes](#get_homes)
    - [switch\_home](#switch_home)
  - [專案結構](#專案結構)
    - [核心檔案說明](#核心檔案說明)
  - [貢獻指南](#貢獻指南)
  - [授權條款](#授權條款)

## 特色

- **全面的設備控制**：支援對 Aqara 智慧設備的開關、亮度、色溫、模式等多種屬性進行精細控制。
- **靈活的設備查詢**：能夠按房間、設備類型查詢設備清單及其詳細狀態。
- **智慧場景管理**：支援查詢和執行使用者預設的智慧家居場景。
- **設備歷史記錄**：查詢設備在指定時間範圍內的歷史狀態變更記錄。
- **自動化配置**：支援配置定時或延時設備控制任務。
- **多家庭支援**：支援查詢和切換使用者帳戶下的不同家庭。
- **MCP 協議相容**：完全遵循 MCP 協議規範，易於與各類 AI 助手整合。
- **安全認證機制**：採用基於登入授權+簽名的安全認證，保護使用者資料和設備安全。
- **跨平台執行**：基於 Go 語言開發，可編譯為多平台可執行檔案。
- **易於擴展**：模組化設計，可以方便地新增新的工具和功能。

## 工作原理

Aqara MCP Server 作為 AI 助手與 Aqara 智慧家居平台之間的橋樑：

1. **AI 助手 (MCP 用戶端)**：使用者透過 AI 助手發出指令 (例如，「打開客廳的燈」)。
2. **MCP 用戶端**：將使用者指令解析，並根據 MCP 協議呼叫 Aqara MCP Server 提供的相應工具 (例如 `device_control`)。
3. **Aqara MCP Server (本專案)**：接收來自用戶端的請求，驗證後呼叫 `smh.go` 模組。
4. **`smh.go` 模組**：使用配置好的 Aqara 憑據，與 Aqara 雲端 API 進行通訊，執行實際的設備操作或資料查詢。
5. **回應流程**：Aqara 雲端 API 返回結果，經由 Aqara MCP Server 傳遞回 MCP 用戶端，最終呈現給使用者。

## 快速開始

### 先決條件

- Go (版本 1.24 或更高)
- Git (用於從原始碼建置)
- Aqara 帳戶及已連結的智慧設備

### 安裝

您可以選擇下載預編譯的可執行檔案或從原始碼建置。

**選項 1: 下載預編譯版本 (推薦)**

造訪下面連結，下載適用於您作業系統的最新可執行檔案包。

[Releases 頁面](https://github.com/aqara/aqara-mcp-server/releases)

解壓後即可使用。

**選項 2: 從原始碼建置**

```bash
# 複製儲存庫
git clone https://github.com/aqara/aqara-mcp-server.git
cd aqara-mcp-server

# 下載相依性
go mod tidy

# 建置可執行檔案
go build -o aqara-mcp-server
```

建置完成後，會在當前目錄下產生 `aqara-mcp-server` 可執行檔案。

### Aqara 帳戶認證

為了使 MCP Server 能夠存取您的 Aqara 帳戶並控制設備，您需要先進行登入授權。

請造訪以下地址完成登入授權：
[https://cdn.aqara.com/app/mcpserver/login.html](https://cdn.aqara.com/app/mcpserver/login.html)

登入成功後，您將獲得必要的認證資訊（如 `token`, `region`），這些資訊將在後續配置步驟中使用。

**請妥善保管這些資訊，尤其是 `token` 不要洩露給他人。**

### 配置範例 (Claude for Desktop)

不同的 MCP 用戶端的配置方法略有不同。以下是如何配置 Claude for Desktop 以使用此 MCP Server 的範例：

1. 打開 Claude for Desktop 的設定 (Settings)。
2. 切換到開發者 (Developer) 標籤頁。
3. 點擊編輯配置 (Edit Config)，使用文字編輯器打開配置檔案。

   ![](/readme/img/setting0.png)
   ![](/readme/img/setting1.png)

4. 將「登入成功頁面」的配置資訊，新增到用戶端的配置檔案(claude_desktop_config.json)中。配置範例：

   ![](/readme/img/config.png)

**配置說明：**
- `command`: 指向您下載或建置的 `aqara-mcp-server` 可執行檔案的完整路徑
- `args`: 使用 `["run", "stdio"]` 啟動 stdio 傳輸模式
- `env`: 環境變數配置
  - `token`: 從 Aqara 登入頁面獲取的存取權杖
  - `region`: 您的 Aqara 帳戶所在區域（如 CN、US、EU 等）

### 執行服務

重新啟動 Claude for Desktop。然後就可以透過對話來呼叫 MCP Server 提供的工具執行設備控制、設備查詢等操作。

![](/readme/img/claude.png)

**其他 MCP 用戶端配置**

對於其他支援 MCP 協議的用戶端（如 Claude、ChatGPT、Cursor 等），配置方式類似：
- 確保用戶端支援 MCP 協議
- 配置可執行檔案路徑和啟動參數
- 設定環境變數 `token` 和 `region`
- 選擇合適的傳輸協議（建議使用 `stdio`）

**SSE 模式 (可選)**

如果您需要使用 SSE (Server-Sent Events) 模式，可以這樣啟動：

```bash
# 使用預設埠 8080
./aqara-mcp-server run sse

# 或指定自訂主機和埠
./aqara-mcp-server run sse --host localhost --port 9000
```

然後在用戶端配置中使用 `["run", "sse"]` 參數。

## 可用工具

MCP 用戶端可以透過呼叫這些工具與 Aqara 智慧家居設備進行互動。

### device_control

- **描述**: 控制智慧家居設備的狀態或屬性（例如開關、溫度、亮度、顏色、色溫等）。
- **參數**:
  - `endpoint_ids` (Array<Integer>, 必需): 需要控制的設備 ID 清單。
  - `control_params` (Object, 必需): 控制參數物件，包含具體操作。
    - `action` (String, 必需): 要執行的操作。例如: `"on"`, `"off"`, `"set"`, `"up"`, `"down"`, `"cooler"`, `"warmer"`。
    - `attribute` (String, 必需): 要控制的設備屬性。例如: `"on_off"`, `"brightness"`, `"color_temperature"`, `"ac_mode"`。
    - `value` (String | Number, 可選): 目標值（當 action 為 "set" 時必需）。
    - `unit` (String, 可選): 值的單位 (例如: `"%"`, `"K"`, `"℃"`)。
- **返回**: (String) 設備控制的操作結果訊息。

### device_query

- **描述**: 根據指定的位置（房間）和設備類型獲取設備清單（不包含即時狀態資訊，僅列出設備及其 ID）。
- **參數**:
  - `positions` (Array<String>, 可選): 房間名稱清單。如果為空陣列或未提供，則表示查詢所有房間。
  - `device_types` (Array<String>, 可選): 設備類型清單。例如: `"Light"`, `"WindowCovering"`, `"AirConditioner"`, `"Button"` 等。如果為空陣列或未提供，則表示查詢所有類型。
- **返回**: (String) Markdown 格式的設備清單，包含設備名稱和 ID。

### device_status_query

- **描述**: 獲取設備的當前狀態資訊（用於查詢顏色、亮度、開關等與狀態相關的屬性）。
- **參數**:
  - `positions` (Array<String>, 可選): 房間名稱清單。如果為空陣列或未提供，則表示查詢所有房間。
  - `device_types` (Array<String>, 可選): 設備類型清單。可選值同 `device_query`。如果為空陣列或未提供，則表示查詢所有類型。
- **返回**: (String) Markdown 格式的設備狀態資訊。

### device_log_query

- **描述**: 查詢設備的記錄。
- **參數**:
  - `endpoint_ids` (Array<Integer>, 必需): 需要查詢歷史記錄的設備 ID 清單。
  - `start_datetime` (String, 可選): 查詢起始時間，格式為 `YYYY-MM-DD HH:MM:SS` (例如: `"2023-05-16 12:00:00"`)。
  - `end_datetime` (String, 可選): 查詢結束時間，格式為 `YYYY-MM-DD HH:MM:SS`。
  - `attribute` (String, 可選): 要查詢的特定設備屬性名稱 (例如: `on_off`, `brightness`)。如果未提供，則查詢該設備所有已記錄屬性的歷史記錄。
- **返回**: (String) Markdown 格式的設備歷史狀態資訊。 (注意: 當前實作可能提示 "This feature will be available soon."，表示功能待完善。)

### run_scenes

- **描述**: 根據場景 ID 執行指定的場景。
- **參數**:
  - `scenes` (Array<Integer>, 必需): 需要執行的場景 ID 清單。
- **返回**: (String) 場景執行的結果訊息。

### get_scenes

- **描述**: 查詢使用者家庭下所有場景，或指定房間內的場景。
- **參數**:
  - `positions` (Array<String>, 可選): 房間名稱清單。如果為空陣列或未提供，則表示查詢整個家庭的場景。
- **返回**: (String) Markdown 格式的場景資訊。

### automation_config

- **描述**: 配置定時或延時設備控制任務。
- **參數**:
  - `scheduled_time` (String, 必需): 設定的時間點（如果是延時任務，基於當前時間點轉化），格式為 `YYYY-MM-DD HH:MM:SS` (例如: `"2025-05-16 12:12:12"`)。
  - `endpoint_ids` (Array<Integer>, 必需): 需要定時控制的設備 ID 清單。
  - `control_params` (Object, 必需): 設備控制參數，使用與 `device_control` 工具相同的格式（包含 action、attribute、value 等）。
- **返回**: (String) 自動化配置結果訊息。

### get_homes

- **描述**: 獲取使用者帳戶下的所有家庭清單。
- **參數**: 無。
- **返回**: (String) 以逗號分隔的家庭名稱清單。如果無資料則返回空字串或相應的提示資訊。

### switch_home

- **描述**: 切換使用者當前操作的家庭。切換後，後續的設備查詢、控制等操作將針對新切換的家庭。
- **參數**:
  - `home_name` (String, 必需): 目標家庭的名稱（應來自 `get_homes` 工具提供的可用清單）。
- **返回**: (String) 切換操作的結果訊息。

## 專案結構

```
.
├── cmd.go                # Cobra CLI 命令定義和程式進入點 (包含 main 函式)
├── server.go             # MCP 伺服器核心邏輯，工具定義和請求處理
├── smh.go                # Aqara 智慧家居平台 API 介面封裝
├── middleware.go         # 中介軟體：使用者認證、逾時控制、異常恢復
├── config.go             # 全域配置管理和環境變數處理
├── go.mod                # Go 模組相依性管理檔案
├── go.sum                # Go 模組相依性校驗和檔案
├── img/                  # README 文件中使用的圖片資源
├── LICENSE               # MIT 開源授權條款
└── README.md             # 專案文件
```

### 核心檔案說明

- **`cmd.go`**: 基於 Cobra 框架的 CLI 實作，定義 `run stdio` 和 `run sse` 啟動模式及主進入函式
- **`server.go`**: MCP 伺服器核心實作，負責工具註冊、請求處理和協議支援
- **`smh.go`**: Aqara 智慧家居平台 API 封裝層，提供設備控制、認證和多家庭支援
- **`middleware.go`**: 請求處理中介軟體，提供認證驗證、逾時控制和異常處理
- **`config.go`**: 全域配置管理，負責環境變數處理和 API 配置

## 貢獻指南

歡迎透過提交 Issue 或 Pull Request 來參與專案貢獻！

在提交程式碼前，請確保：
1. 程式碼遵循 Go 語言的編碼規範。
2. 相關的 MCP 工具和提示介面定義保持一致性和清晰性。
3. 新增或更新單元測試以覆蓋您的更改。
4. 如有必要，更新相關的文件 (如本 README)。
5. 確保您的提交訊息清晰明瞭。

## 授權條款

本專案基於 [MIT License](/LICENSE) 授權。
Copyright (c) 2025 Aqara-Copliot