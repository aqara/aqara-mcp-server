# Aqara MCP Server

[English](/readme/README.md) | 中文 | [繁體中文](/readme/README_CHT.md) | [Français](/readme/README_FR.md) | [한국어](/readme/README_KR.md) | [Español](/readme/README_ES.md) | [日本語](/readme/README_JP.md) | [Deutsch](/readme/README_DE.md) | [Italiano](/readme/README_IT.md)

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/aqara/aqara-mcp-server)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![Release](https://img.shields.io/github/v/release/aqara/aqara-mcp-server)](https://github.com/aqara/aqara-mcp-server/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Aqara MCP Server 是一个基于 [MCP (Model Context Protocol)](https://modelcontextprotocol.io/introduction) 协议开发的智能家居控制服务。它允许任何支持 MCP 协议的 AI 助手或API（例如 Claude、ChatGPT、Cursor 等）与您的 Aqara 智能家居设备进行交互，实现通过自然语言控制设备、查询状态、执行场景等功能。

## 目录

- [Aqara MCP Server](#aqara-mcp-server)
  - [目录](#目录)
  - [特性](#特性)
  - [工作原理](#工作原理)
  - [快速开始](#快速开始)
    - [先决条件](#先决条件)
    - [安装](#安装)
    - [Aqara 账户认证](#aqara-账户认证)
    - [配置示例 (Claude for Desktop)](#配置示例-claude-for-desktop)
    - [运行服务](#运行服务)
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
  - [项目结构](#项目结构)
    - [核心文件说明](#核心文件说明)
  - [贡献指南](#贡献指南)
  - [许可证](#许可证)

## 特性

-   **全面的设备控制**：支持对 Aqara 智能设备的开关、亮度、色温、模式等多种属性进行精细控制。
-   **灵活的设备查询**：能够按房间、设备类型查询设备列表及其详细状态。
-   **智能场景管理**：支持查询和执行用户预设的智能家居场景。
-   **设备历史记录**：查询设备在指定时间范围内的历史状态变更记录。
-   **自动化配置**：支持配置定时或延时设备控制任务。
-   **多家庭支持**：支持查询和切换用户账户下的不同家庭。
-   **MCP 协议兼容**：完全遵循 MCP 协议规范，易于与各类 AI 助手集成。
-   **安全认证机制**：采用基于登录授权+签名的安全认证，保护用户数据和设备安全。
-   **跨平台运行**：基于 Go 语言开发，可编译为多平台可执行文件。
-   **易于扩展**：模块化设计，可以方便地添加新的工具和功能。

## 工作原理

Aqara MCP Server 作为 AI 助手与 Aqara 智能家居平台之间的桥梁：

1.  **AI 助手 (MCP 客户端)**：用户通过 AI 助手发出指令 (例如，"打开客厅的灯")。
2.  **MCP 客户端**：将用户指令解析，并根据 MCP 协议调用 Aqara MCP Server 提供的相应工具 (例如 `device_control`)。
3.  **Aqara MCP Server (本项目)**：接收来自客户端的请求，验证后调用 `smh.go` 模块。
4.  **`smh.go` 模块**：使用配置好的 Aqara 凭据，与 Aqara 云端 API 进行通信，执行实际的设备操作或数据查询。
5.  **响应流程**：Aqara 云端 API 返回结果，经由 Aqara MCP Server 传递回 MCP 客户端，最终呈现给用户。

## 快速开始

### 先决条件

-   Go (版本 1.24 或更高)
-   Git (用于从源码构建)
-   Aqara 账户及已绑定的智能设备

### 安装

您可以选择下载预编译的可执行文件或从源码构建。

**选项 1: 下载预编译版本 (推荐)**

访问下面链接,下载适用于您操作系统的最新可执行文件包。

[Releases 页面](https://github.com/aqara/aqara-mcp-server/releases)

解压后即可使用。

**选项 2: 从源码构建**

```bash
# 克隆仓库
git clone https://github.com/aqara/aqara-mcp-server.git
cd aqara-mcp-server

# 下载依赖
go mod tidy

# 构建可执行文件
go build -o aqara-mcp-server
```
构建完成后，会在当前目录下生成 `aqara-mcp-server` 可执行文件。

### Aqara 账户认证

为了使 MCP Server 能够访问您的 Aqara 账户并控制设备，您需要先进行登录授权。

请访问以下地址完成登录授权：
[https://cdn.aqara.com/app/mcpserver/login.html](https://cdn.aqara.com/app/mcpserver/login.html)

登录成功后，您将获得必要的认证信息（如 `token`, `region`），这些信息将在后续配置步骤中使用。

**请妥善保管这些信息，尤其是 `token` 不要泄露给他人。**

### 配置示例 (Claude for Desktop)

不同的 MCP 客户端的配置方法略有不同。以下是如何配置 Claude for Desktop 以使用此 MCP Server 的示例：

1.  打开 Claude for Desktop 的设置 (Settings)。
2.  切换到开发者 (Developer) 标签页。
3.  点击编辑配置 (Edit Config)，使用文本编辑器打开配置文件。

    ![](/readme/img/setting0.png)
    ![](/readme/img/setting1.png)

4.  将"登录成功页面"的配置信息，添加到客户端的配置文件(claude_desktop_config.json)中。配置示例：

    ![](/readme/img/config.png)

**配置说明：**
- `command`: 指向您下载或构建的 `aqara-mcp-server` 可执行文件的完整路径
- `args`: 使用 `["run", "stdio"]` 启动 stdio 传输模式
- `env`: 环境变量配置
  - `token`: 从 Aqara 登录页面获取的访问令牌
  - `region`: 您的 Aqara 账户所在区域（如 CN、US、EU 等）

### 运行服务

重启 Claude for Desktop。然后就可以通过对话来调用 MCP Server 提供的工具执行设备控制、设备查询等操作。

![](/readme/img/claude.png)

**其他 MCP 客户端配置**

对于其他支持 MCP 协议的客户端（如 Claude、ChatGPT、Cursor 等），配置方式类似：
- 确保客户端支持 MCP 协议
- 配置可执行文件路径和启动参数
- 设置环境变量 `token` 和 `region`
- 选择合适的传输协议（建议使用 `stdio`）

**SSE 模式 (可选)**

如果您需要使用 SSE (Server-Sent Events) 模式，可以这样启动：

```bash
# 使用默认端口 8080
./aqara-mcp-server run sse

# 或指定自定义主机和端口
./aqara-mcp-server run sse --host localhost --port 9000
```

然后在客户端配置中使用 `["run", "sse"]` 参数。

## 可用工具

MCP 客户端可以通过调用这些工具与 Aqara 智能家居设备进行交互。

### device_control

-   **描述**: 控制智能家居设备的状态或属性（例如开关、温度、亮度、颜色、色温等）。
-   **参数**:
    -   `endpoint_ids` (Array<Integer>, 必需): 需要控制的设备 ID 列表。
    -   `control_params` (Object, 必需): 控制参数对象，包含具体操作。
        -   `action` (String, 必需): 要执行的操作。例如: `"on"`, `"off"`, `"set"`, `"up"`, `"down"`, `"cooler"`, `"warmer"`。
        -   `attribute` (String, 必需): 要控制的设备属性。例如: `"on_off"`, `"brightness"`, `"color_temperature"`, `"ac_mode"`。
        -   `value` (String | Number, 可选): 目标值（当 action 为 "set" 时必需）。
        -   `unit` (String, 可选): 值的单位 (例如: `"%"`, `"K"`, `"℃"`)。
-   **返回**: (String) 设备控制的操作结果消息。

### device_query

-   **描述**: 根据指定的位置（房间）和设备类型获取设备列表（不包含实时状态信息，仅列出设备及其 ID）。
-   **参数**:
    -   `positions` (Array<String>, 可选): 房间名称列表。如果为空数组或未提供，则表示查询所有房间。
    -   `device_types` (Array<String>, 可选): 设备类型列表。例如: `"Light"`, `"WindowCovering"`, `"AirConditioner"`, `"Button"` 等。如果为空数组或未提供，则表示查询所有类型。
-   **返回**: (String) Markdown 格式的设备列表，包含设备名称和 ID。

### device_status_query

-   **描述**: 获取设备的当前状态信息（用于查询颜色、亮度、开关等与状态相关的属性）。
-   **参数**:
    -   `positions` (Array<String>, 可选): 房间名称列表。如果为空数组或未提供，则表示查询所有房间。
    -   `device_types` (Array<String>, 可选): 设备类型列表。可选值同 `device_query`。如果为空数组或未提供，则表示查询所有类型。
-   **返回**: (String) Markdown 格式的设备状态信息。

### device_log_query

-   **描述**: 查询设备的日志。
-   **参数**:
    -   `endpoint_ids` (Array<Integer>, 必需): 需要查询历史记录的设备 ID 列表。
    -   `start_datetime` (String, 可选): 查询起始时间，格式为 `YYYY-MM-DD HH:MM:SS` (例如: `"2023-05-16 12:00:00"`)。
    -   `end_datetime` (String, 可选): 查询结束时间，格式为 `YYYY-MM-DD HH:MM:SS`。
    -   `attribute` (String, 可选): 要查询的特定设备属性名称 (例如: `on_off`, `brightness`)。如果未提供，则查询该设备所有已记录属性的历史记录。
-   **返回**: (String) Markdown 格式的设备历史状态信息。 (注意: 当前实现可能提示 "This feature will be available soon."，表示功能待完善。)

### run_scenes

-   **描述**: 根据场景 ID 执行指定的场景。
-   **参数**:
    -   `scenes` (Array<Integer>, 必需): 需要执行的场景 ID 列表。
-   **返回**: (String) 场景执行的结果消息。

### get_scenes

-   **描述**: 查询用户家庭下所有场景，或指定房间内的场景。
-   **参数**:
    -   `positions` (Array<String>, 可选): 房间名称列表。如果为空数组或未提供，则表示查询整个家庭的场景。
-   **返回**: (String) Markdown 格式的场景信息。

### automation_config

-   **描述**: 配置定时或延时设备控制任务。
-   **参数**:
    -   `scheduled_time` (String, 必需): 设置的时间点（如果是延时任务，基于当前时间点转化），格式为 `YYYY-MM-DD HH:MM:SS` (例如: `"2025-05-16 12:12:12"`)。
    -   `endpoint_ids` (Array<Integer>, 必需): 需要定时控制的设备 ID 列表。
    -   `control_params` (Object, 必需): 设备控制参数，使用与 `device_control` 工具相同的格式（包含 action、attribute、value 等）。
-   **返回**: (String) 自动化配置结果消息。

### get_homes

-   **描述**: 获取用户账户下的所有家庭列表。
-   **参数**: 无。
-   **返回**: (String) 以逗号分隔的家庭名称列表。如果无数据则返回空字符串或相应的提示信息。

### switch_home

-   **描述**: 切换用户当前操作的家庭。切换后，后续的设备查询、控制等操作将针对新切换的家庭。
-   **参数**:
    -   `home_name` (String, 必需): 目标家庭的名称（应来自 `get_homes` 工具提供的可用列表）。
-   **返回**: (String) 切换操作的结果消息。

## 项目结构

```
.
├── cmd.go                # Cobra CLI 命令定义和程序入口点 (包含 main 函数)
├── server.go             # MCP 服务器核心逻辑，工具定义和请求处理
├── smh.go                # Aqara 智能家居平台 API 接口封装
├── middleware.go         # 中间件：用户认证、超时控制、异常恢复
├── config.go             # 全局配置管理和环境变量处理
├── go.mod                # Go 模块依赖管理文件
├── go.sum                # Go 模块依赖校验和文件
├── img/                  # README 文档中使用的图片资源
├── LICENSE               # MIT 开源许可证
└── README.md             # 项目文档
```

### 核心文件说明

-   **`cmd.go`**: 基于 Cobra 框架的 CLI 实现，定义 `run stdio` 和 `run sse` 启动模式及主入口函数
-   **`server.go`**: MCP 服务器核心实现，负责工具注册、请求处理和协议支持
-   **`smh.go`**: Aqara 智能家居平台 API 封装层，提供设备控制、认证和多家庭支持
-   **`middleware.go`**: 请求处理中间件，提供认证验证、超时控制和异常处理
-   **`config.go`**: 全局配置管理，负责环境变量处理和API配置

## 贡献指南

欢迎通过提交 Issue 或 Pull Request 来参与项目贡献！

在提交代码前，请确保：
1.  代码遵循 Go 语言的编码规范。
2.  相关的 MCP 工具和提示接口定义保持一致性和清晰性。
3.  添加或更新单元测试以覆盖您的更改。
4.  如有必要，更新相关的文档 (如本 README)。
5.  确保您的提交信息清晰明了。

## 许可证

本项目基于 [MIT License](/LICENSE) 授权。
Copyright (c) 2025 Aqara-Copliot
