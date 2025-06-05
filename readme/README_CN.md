<div align="center" style="display: flex; align-items: center; justify-content: center; ">

  <img src="/readme/img/logo.png" alt="Aqara Logo" height="120">
  <h1>Aqara MCP Server</h1>

</div>

<div align="center">

[English](/readme/README.md) | 中文 | [繁體中文](/readme/README_CHT.md) | [Français](/readme/README_FR.md) | [한국어](/readme/README_KR.md) | [Español](/readme/README_ES.md) | [日本語](/readme/README_JP.md) | [Deutsch](/readme/README_DE.md) | [Italiano](/readme/README_IT.md)

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/aqara/aqara-mcp-server)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![Release](https://img.shields.io/github/v/release/aqara/aqara-mcp-server)](https://github.com/aqara/aqara-mcp-server/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

Aqara MCP Server 是一个基于 [MCP (Model Context Protocol)](https://modelcontextprotocol.io/introduction) 协议开发的智能家居控制服务。它允许任何支持 MCP 协议的 AI 助手或API（例如 Claude、Cursor 等）与您的 Aqara 智能家居设备进行交互，实现通过自然语言控制设备、查询状态、执行场景等功能。

## 目录

- [目录](#目录)
- [特性](#特性)
- [工作原理](#工作原理)
- [快速开始](#快速开始)
  - [先决条件](#先决条件)
  - [安装](#安装)
    - [方式一：下载预编译版本（推荐）](#方式一下载预编译版本推荐)
    - [方式二：从源码构建](#方式二从源码构建)
  - [Aqara 账户认证](#aqara-账户认证)
  - [客户端配置](#客户端配置)
    - [Claude for Desktop 配置示例](#claude-for-desktop-配置示例)
    - [配置参数说明](#配置参数说明)
    - [其他 MCP 客户端](#其他-mcp-客户端)
  - [启动服务](#启动服务)
    - [标准模式（推荐）](#标准模式推荐)
    - [HTTP 模式（`即将支持`）](#http-模式即将支持)
- [API 工具说明](#api-工具说明)
  - [设备控制类](#设备控制类)
    - [device\_control](#device_control)
  - [设备查询类](#设备查询类)
    - [device\_query](#device_query)
    - [device\_status\_query](#device_status_query)
    - [device\_log\_query](#device_log_query)
  - [场景管理类](#场景管理类)
    - [get\_scenes](#get_scenes)
    - [run\_scenes](#run_scenes)
  - [家庭管理类](#家庭管理类)
    - [get\_homes](#get_homes)
    - [switch\_home](#switch_home)
  - [自动化配置类](#自动化配置类)
    - [automation\_config](#automation_config)
- [项目结构](#项目结构)
  - [目录结构](#目录结构)
  - [核心文件说明](#核心文件说明)
- [开发指南](#开发指南)
- [许可证](#许可证)

## 特性

- ✨ **全面的设备控制**：支持对 Aqara 智能设备的开关、亮度、色温、模式等多种属性进行精细控制
- 🔍 **灵活的设备查询**：能够按房间、设备类型查询设备列表及其详细状态
- 🎬 **智能场景管理**：支持查询和执行用户预设的智能家居场景
- 📈 **设备历史记录**：查询设备在指定时间范围内的历史状态变更记录
- ⏰ **自动化配置**：支持配置定时或延时设备控制任务
- 🏠 **多家庭支持**：支持查询和切换用户账户下的不同家庭
- 🔌 **MCP 协议兼容**：完全遵循 MCP 协议规范，易于与各类 AI 助手集成
- 🔐 **安全认证机制**：采用基于登录授权+签名的安全认证，保护用户数据和设备安全
- 🌐 **跨平台运行**：基于 Go 语言开发，可编译为多平台可执行文件
- 🔧 **易于扩展**：模块化设计，可以方便地添加新的工具和功能

## 工作原理

Aqara MCP Server 作为 AI 助手与 Aqara 智能家居平台之间的桥梁：

```mermaid
graph LR
    A[AI 助手] --> B[MCP 客户端]
    B --> C[Aqara MCP Server]
    C --> D[Aqara 云端 API]
    D --> E[智能设备]
```

1. **AI 助手**：用户通过 AI 助手发出指令（例如，"打开客厅的灯"）
2. **MCP 客户端**：将用户指令解析，并根据 MCP 协议调用 Aqara MCP Server 提供的相应工具（例如 `device_control`）
3. **Aqara MCP Server (本项目)**：接收来自客户端的请求，使用配置好的 Aqara 凭据，与 Aqara 云端 API 进行通信，执行实际的设备操作或数据查询
4. **响应流程**：Aqara 云端 API 返回结果，经由 Aqara MCP Server 传递回 MCP 客户端，最终呈现给用户

## 快速开始

### 先决条件

- **Go** (版本 1.24 或更高) - 仅在从源码构建时需要
- **Git** (用于从源码构建) - 可选
- **Aqara 账户**及已绑定的智能设备
- **支持 MCP 协议的客户端** (如 Claude for Desktop、Cursor 等)

### 安装

您可以选择下载预编译的可执行文件或从源码构建。

#### 方式一：下载预编译版本（推荐）

访问 GitHub Releases 页面，下载适用于您操作系统的最新可执行文件：

**📥 [前往 Releases 页面下载](https://github.com/aqara/aqara-mcp-server/releases)**

下载对应平台的压缩包后解压即可使用。

#### 方式二：从源码构建

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
**🔗 [https://cdn.aqara.com/app/mcpserver/login.html](https://cdn.aqara.com/app/mcpserver/login.html)**

登录成功后，您将获得必要的认证信息（如 `token`, `region`），这些信息将在后续配置步骤中使用。

> ⚠️ **安全提醒**：请妥善保管 `token` 信息，不要泄露给他人。

### 客户端配置

不同的 MCP 客户端的配置方法略有不同。以下是如何配置 Claude for Desktop 以使用此 MCP Server 的示例：

#### Claude for Desktop 配置示例

1. **打开 Claude for Desktop 的设置 (Settings)**

    ![Claude Open Setting](/readme/img/opening_setting.png)

2. **切换到开发者 (Developer) 标签页，然后点击编辑配置 (Edit Config)，使用文本编辑器打开配置文件**

    ![Claude Edit Configuration](/readme/img/edit_config.png)

3. **将"登录成功页面"的配置信息，添加到客户端的配置文件 `claude_desktop_config.json` 中**

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

#### 配置参数说明

- `command`: 指向您下载或构建的 `aqara-mcp-server` 可执行文件的完整路径
- `args`: 使用 `["run", "stdio"]` 启动 stdio 传输模式
- `env`: 环境变量配置
  - `token`: 从 Aqara 登录页面获取的访问令牌
  - `region`: 您的 Aqara 账户所在区域（支持的区域：CN、US、EU、KR、SG、RU）

#### 其他 MCP 客户端

对于其他支持 MCP 协议的客户端（如 ChatGPT、Cursor 等），配置方式类似：

- 确保客户端支持 MCP 协议
- 配置可执行文件路径和启动参数
- 设置环境变量 `token` 和 `region`
- 选择合适的传输协议（建议使用 `stdio`）

### 启动服务

#### 标准模式（推荐）

重启 Claude for Desktop。然后就可以通过自然语言执行设备控制、设备查询、场景执行等操作。

示例对话：

- "打开客厅的灯"
- "把卧室空调设置为制冷模式，温度 24 度"
- "查看所有房间的设备列表"
- "执行晚安场景"

![Claude Chat Example](/readme/img/claude.png)

#### HTTP 模式（`即将支持`）

## API 工具说明

MCP 客户端可以通过调用这些工具与 Aqara 智能家居设备进行交互。

### 设备控制类

#### device_control

控制智能家居设备的状态或属性（例如开关、温度、亮度、颜色、色温等）。

**参数：**

- `endpoint_ids` _(Array\<Integer\>, 必需)_：需要控制的设备 ID 列表
- `control_params` _(Object, 必需)_：控制参数对象，包含具体操作：
  - `action` _(String, 必需)_：要执行的操作（如 `"on"`, `"off"`, `"set"`, `"up"`, `"down"`, `"cooler"`, `"warmer"`）
  - `attribute` _(String, 必需)_：要控制的设备属性（如 `"on_off"`, `"brightness"`, `"color_temperature"`, `"ac_mode"`）
  - `value` _(String | Number, 可选)_：目标值（当 action 为 "set" 时必需）
  - `unit` _(String, 可选)_：值的单位（如 `"%"`, `"K"`, `"℃"`）

**返回：** 设备控制的操作结果消息

### 设备查询类

#### device_query

根据指定的位置（房间）和设备类型获取设备列表（不包含实时状态信息）。

**参数：**

- `positions` _(Array\<String\>, 可选)_：房间名称列表。空数组表示查询所有房间
- `device_types` _(Array\<String\>, 可选)_：设备类型列表（如 `"Light"`, `"WindowCovering"`, `"AirConditioner"`, `"Button"`）。空数组表示查询所有类型

**返回：** Markdown 格式的设备列表，包含设备名称和 ID

#### device_status_query

获取设备的当前状态信息（用于查询颜色、亮度、开关等实时状态信息）。

**参数：**

- `positions` _(Array\<String\>, 可选)_：房间名称列表。空数组表示查询所有房间
- `device_types` _(Array\<String\>, 可选)_：设备类型列表。可选值同 `device_query`。空数组表示查询所有类型

**返回：** Markdown 格式的设备状态信息

#### device_log_query

查询设备的历史日志信息。

**参数：**

- `endpoint_ids` _(Array\<Integer\>, 必需)_：需要查询历史记录的设备 ID 列表
- `start_datetime` _(String, 可选)_：查询起始时间，格式为 `YYYY-MM-DD HH:MM:SS`（例如：`"2023-05-16 12:00:00"`）
- `end_datetime` _(String, 可选)_：查询结束时间，格式为 `YYYY-MM-DD HH:MM:SS`
- `attributes` _(Array\<String\>, 可选)_：要查询的设备属性名称列表（如 `["on_off", "brightness"]`）。未提供时查询所有已记录属性

**返回：** Markdown 格式的设备历史状态信息

### 场景管理类

#### get_scenes

查询用户家庭下所有场景，或指定房间内的场景。

**参数：**

- `positions` _(Array\<String\>, 可选)_：房间名称列表。空数组表示查询整个家庭的场景

**返回：** Markdown 格式的场景信息

#### run_scenes

根据场景 ID 执行指定的场景。

**参数：**

- `scenes` _(Array\<Integer\>, 必需)_：需要执行的场景 ID 列表

**返回：** 场景执行的结果消息

### 家庭管理类

#### get_homes

获取用户账户下的所有家庭列表。

**参数：** 无

**返回：** 以逗号分隔的家庭名称列表。如果无数据则返回空字符串或相应的提示信息

#### switch_home

切换用户当前操作的家庭。切换后，后续的设备查询、控制等操作将针对新切换的家庭。

**参数：**

- `home_name` _(String, 必需)_：目标家庭的名称

**返回：** 切换操作的结果消息

### 自动化配置类

#### automation_config

自动化配置（目前仅支持定时或延时设备控制任务）。

**参数：**

- `scheduled_time` _(String, 必需)_：定时执行的时间点，使用标准 Crontab 格式 `"分 时 日 月 周"`。例如：`"30 14 * * *"`（每天14:30执行）、`"0 9 * * 1"`（每周一9:00执行）
- `endpoint_ids` _(Array\<Integer\>, 必需)_：需要定时控制的设备 ID 列表
- `control_params` _(Object, 必需)_：设备控制参数，使用与 `device_control` 工具相同的格式（包含 action、attribute、value 等）
- `task_name` _(String, 必需)_：此自动化任务的名称或描述（用于识别和管理）
- `execution_once` _(Boolean, 可选)_：是否只执行一次
  - `true`：仅在指定时间执行一次任务（默认值）
  - `false`：周期性重复执行任务（如每天、每周等）

**返回：** 自动化配置结果消息

## 项目结构

### 目录结构

```text
.
├── cmd.go                # Cobra CLI 命令定义和程序入口点（包含 main 函数）
├── server.go             # MCP 服务器核心逻辑，工具定义和请求处理
├── smh.go                # Aqara 智能家居平台 API 接口封装
├── middleware.go         # 中间件：用户认证、超时控制、异常恢复
├── config.go             # 全局配置管理和环境变量处理
├── go.mod                # Go 模块依赖管理文件
├── go.sum                # Go 模块依赖校验和文件
├── readme/               # README 文档和图片资源
│   ├── img/              # 图片资源目录
│   └── *.md              # 多语言 README 文件
├── LICENSE               # MIT 开源许可证
└── README.md             # 项目主文档
```

### 核心文件说明

- **`cmd.go`**：基于 Cobra 框架的 CLI 实现，定义 `run stdio` 和 `run http` 启动模式及主入口函数
- **`server.go`**：MCP 服务器核心实现，负责工具注册、请求处理和协议支持
- **`smh.go`**：Aqara 智能家居平台 API 封装层，提供设备控制、认证和多家庭支持
- **`middleware.go`**：请求处理中间件，提供认证验证、超时控制和异常处理
- **`config.go`**：全局配置管理，负责环境变量处理和 API 配置

## 开发指南

欢迎通过提交 Issue 或 Pull Request 来参与项目贡献！

在提交代码前，请确保：

1. 代码遵循 Go 语言的编码规范
2. 相关的 MCP 工具和接口定义保持一致性和清晰性
3. 添加或更新单元测试以覆盖您的更改
4. 如有必要，更新相关的文档（如本 README）
5. 确保您的提交信息清晰明了

**🌟 如果这个项目对您有帮助，请给我们一个 Star！**

**🤝 欢迎加入我们的社区，一起让智能家居更智能！**

## 许可证

本项目基于 [MIT License](/LICENSE) 授权。

---

Copyright (c) 2025 Aqara-Copilot
