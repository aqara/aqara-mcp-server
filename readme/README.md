<div align="center" style="display: flex; align-items: center; justify-content: center; ">

  <img src="/readme/img/logo.png" alt="Aqara Logo" height="120">
  <h1>Aqara MCP Server</h1>

</div>

<div align="center">

English | [中文](/readme/README_CN.md) | [繁體中文](/readme/README_CHT.md) | [Français](/readme/README_FR.md) | [한국어](/readme/README_KR.md) | [Español](/readme/README_ES.md) | [日本語](/readme/README_JP.md) | [Deutsch](/readme/README_DE.md) | [Italiano](/readme/README_IT.md)

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/aqara/aqara-mcp-server)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![Release](https://img.shields.io/github/v/release/aqara/aqara-mcp-server)](https://github.com/aqara/aqara-mcp-server/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

Aqara MCP Server is a smart home control service developed based on the [MCP (Model Context Protocol)](https://modelcontextprotocol.io/introduction) protocol. It allows any AI assistant or API that supports the MCP protocol (such as Claude, Cursor, etc.) to interact with your Aqara smart home devices, enabling device control, status queries, scene execution, and more through natural language.

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Features](#features)
- [How It Works](#how-it-works)
- [Quick Start](#quick-start)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
    - [Option 1: Download Pre-compiled Version (Recommended)](#option-1-download-pre-compiled-version-recommended)
    - [Option 2: Build from Source](#option-2-build-from-source)
  - [Aqara Account Authentication](#aqara-account-authentication)
  - [Client Configuration](#client-configuration)
    - [Claude for Desktop Configuration Example](#claude-for-desktop-configuration-example)
    - [Configuration Parameters](#configuration-parameters)
    - [Other MCP Clients](#other-mcp-clients)
  - [Starting the Service](#starting-the-service)
    - [Standard Mode (Recommended)](#standard-mode-recommended)
    - [HTTP Mode (Optional)](#http-mode-optional)
- [API Tools Documentation](#api-tools-documentation)
  - [Device Control](#device-control)
    - [device\_control](#device_control)
  - [Device Query](#device-query)
    - [device\_query](#device_query)
    - [device\_status\_query](#device_status_query)
    - [device\_log\_query](#device_log_query)
  - [Scene Management](#scene-management)
    - [get\_scenes](#get_scenes)
    - [run\_scenes](#run_scenes)
  - [Home Management](#home-management)
    - [get\_homes](#get_homes)
    - [switch\_home](#switch_home)
  - [Automation Configuration](#automation-configuration)
    - [automation\_config](#automation_config)
- [Project Structure](#project-structure)
  - [Directory Structure](#directory-structure)
  - [Core Files Description](#core-files-description)
- [Development Guide](#development-guide)
- [License](#license)

## Features

- **Comprehensive Device Control**: Supports fine-grained control of various attributes of Aqara smart devices including switches, brightness, color temperature, modes, etc.
- **Flexible Device Queries**: Query device lists and their detailed status by room and device type
- **Smart Scene Management**: Support for querying and executing user-preset smart home scenes
- **Device History Records**: Query historical status change records of devices within specified time ranges
- **Automation Configuration**: Support for configuring scheduled or delayed device control tasks
- **Multi-Home Support**: Support for querying and switching between different homes under user accounts
- **MCP Protocol Compatibility**: Fully compliant with MCP protocol specifications, easy to integrate with various AI assistants
- **Secure Authentication**: Uses login authorization + signature-based security authentication to protect user data and device security
- **Cross-Platform**: Developed in Go language, can be compiled to executable files for multiple platforms
- **Easy to Extend**: Modular design allows for easy addition of new tools and features

## How It Works

Aqara MCP Server serves as a bridge between AI assistants and the Aqara smart home platform:

1. **AI Assistant (MCP Client)**: Users issue commands through AI assistants (e.g., "Turn on the living room lights")
2. **MCP Client**: Parses user commands and calls corresponding tools provided by Aqara MCP Server according to MCP protocol (e.g., `device_control`)
3. **Aqara MCP Server (this project)**: Receives requests from clients, validates them, and calls the `smh.go` module
4. **`smh.go` Module**: Uses configured Aqara credentials to communicate with Aqara cloud APIs, executing actual device operations or data queries
5. **Response Flow**: Aqara cloud APIs return results, which are passed back through Aqara MCP Server to the MCP client and finally presented to the user

## Quick Start

### Prerequisites

- Go (version 1.24 or higher)
- Git (for building from source)
- Aqara account with bound smart devices

### Installation

You can choose to download pre-compiled executable files or build from source.

#### Option 1: Download Pre-compiled Version (Recommended)

Visit the GitHub Releases page to download the latest executable file for your operating system:

**📥 [Go to Releases Page](https://github.com/aqara/aqara-mcp-server/releases)**

Download and extract the appropriate package for your platform.

#### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/aqara/aqara-mcp-server.git
cd aqara-mcp-server

# Download dependencies
go mod tidy

# Build executable
go build -o aqara-mcp-server
```

After building, an `aqara-mcp-server` executable will be generated in the current directory.

### Aqara Account Authentication

To enable the MCP Server to access your Aqara account and control devices, you need to complete login authorization first.

Please visit the following address to complete login authorization:
**🔗 [https://cdn.aqara.com/app/mcpserver/login.html](https://cdn.aqara.com/app/mcpserver/login.html)**

After successful login, you will obtain necessary authentication information (such as `token`, `region`), which will be used in subsequent configuration steps.

> ⚠️ **Security Reminder**: Please keep your `token` information secure and do not share it with others.

### Client Configuration

Different MCP clients have slightly different configuration methods. Here's an example of how to configure Claude for Desktop to use this MCP Server:

#### Claude for Desktop Configuration Example

1. Open Claude for Desktop Settings

    ![Claude Open Setting](/readme/img/opening_setting.png)

2. Switch to the Developer tab, then click Edit Config to open the configuration file with a text editor

    ![Claude Edit Configuration](/readme/img/edit_config.png)

3. Add the configuration information from the "Login Success Page" to the client's configuration file `claude_desktop_config.json`

    ![Configuration Example](/readme/img/config_info.png)

#### Configuration Parameters

- `command`: Full path to your downloaded or built `aqara-mcp-server` executable file
- `args`: Use `["run", "stdio"]` to start stdio transport mode
- `env`: Environment variable configuration
  - `token`: Access token obtained from the Aqara login page
  - `region`: Your Aqara account region (e.g., CN, US, EU, etc.)

#### Other MCP Clients

For other clients that support the MCP protocol (such as ChatGPT, Cursor, etc.), the configuration is similar:

- Ensure the client supports MCP protocol
- Configure executable file path and startup parameters
- Set environment variables `token` and `region`
- Choose appropriate transport protocol (stdio recommended)

### Starting the Service

#### Standard Mode (Recommended)

Restart Claude for Desktop. You can then perform device control, device queries, scene execution, and other operations through natural language.

![Claude Chat Example](/readme/img/claude.png)

#### HTTP Mode (Optional)

If you need to use HTTP mode, you can start it like this:

```bash
# Use default port 8080
./aqara-mcp-server run http

# Or specify custom host and port
./aqara-mcp-server run http --host localhost --port 9000
```

Then use `["run", "http"]` parameters in the client configuration.

## API Tools Documentation

MCP clients can interact with Aqara smart home devices by calling these tools.

### Device Control

#### device_control

Control the status or attributes of smart home devices (e.g., switches, temperature, brightness, color, color temperature, etc.).

**Parameters:**

- `endpoint_ids` _(Array\<Integer\>, required)_: List of device IDs to control
- `control_params` _(Object, required)_: Control parameter object containing specific operations:
  - `action` _(String, required)_: Operation to execute (e.g., `"on"`, `"off"`, `"set"`, `"up"`, `"down"`, `"cooler"`, `"warmer"`)
  - `attribute` _(String, required)_: Device attribute to control (e.g., `"on_off"`, `"brightness"`, `"color_temperature"`, `"ac_mode"`)
  - `value` _(String | Number, optional)_: Target value (required when action is "set")
  - `unit` _(String, optional)_: Unit of the value (e.g., `"%"`, `"K"`, `"℃"`)

**Returns:** Operation result message for device control

### Device Query

#### device_query

Get device list based on specified location (room) and device type (does not include real-time status information).

**Parameters:**

- `positions` _(Array\<String\>, optional)_: List of room names. Empty array means query all rooms
- `device_types` _(Array\<String\>, optional)_: List of device types (e.g., `"Light"`, `"WindowCovering"`, `"AirConditioner"`, `"Button"`). Empty array means query all types

**Returns:** Markdown-formatted device list including device names and IDs

#### device_status_query

Get current status information of devices (for querying real-time status information like color, brightness, switches, etc.).

**Parameters:**

- `positions` _(Array\<String\>, optional)_: List of room names. Empty array means query all rooms
- `device_types` _(Array\<String\>, optional)_: List of device types. Same options as `device_query`. Empty array means query all types

**Returns:** Markdown-formatted device status information

#### device_log_query

Query historical log information of devices.

**Parameters:**

- `endpoint_ids` _(Array\<Integer\>, required)_: List of device IDs to query history for
- `start_datetime` _(String, optional)_: Query start time in format `YYYY-MM-DD HH:MM:SS` (e.g., `"2023-05-16 12:00:00"`)
- `end_datetime` _(String, optional)_: Query end time in format `YYYY-MM-DD HH:MM:SS`
- `attribute` _(String, optional)_: Specific device attribute name to query (e.g., `on_off`, `brightness`). Queries all recorded attributes when not provided

**Returns:** Markdown-formatted device historical status information

> 📝 **Note:** Current implementation may show "This feature will be available soon.", indicating the feature is pending completion.

### Scene Management

#### get_scenes

Query all scenes in the user's home, or scenes within specified rooms.

**Parameters:**

- `positions` _(Array\<String\>, optional)_: List of room names. Empty array means query scenes for the entire home

**Returns:** Markdown-formatted scene information

#### run_scenes

Execute specified scenes based on scene IDs.

**Parameters:**

- `scenes` _(Array\<Integer\>, required)_: List of scene IDs to execute

**Returns:** Scene execution result message

### Home Management

#### get_homes

Get list of all homes under the user account.

**Parameters:** None

**Returns:** Comma-separated list of home names. Returns empty string or appropriate message if no data

#### switch_home

Switch the user's current operating home. After switching, subsequent device queries, controls, and other operations will target the newly switched home.

**Parameters:**

- `home_name` _(String, required)_: Name of the target home

**Returns:** Switch operation result message

### Automation Configuration

#### automation_config

Configure scheduled or delayed device control tasks (currently only supports timed delay automation configuration).

**Parameters:**

- `scheduled_time` _(String, required)_: Set time point (if delay task, converted based on current time), format `YYYY-MM-DD HH:MM:SS` (e.g., `"2025-05-16 12:12:12"`)
- `endpoint_ids` _(Array\<Integer\>, required)_: List of device IDs for scheduled control
- `control_params` _(Object, required)_: Device control parameters using the same format as `device_control` tool (including action, attribute, value, etc.)

**Returns:** Automation configuration result message

> 📝 **Note:** Current implementation may show "This feature will be available soon.", indicating the feature is pending completion.

## Project Structure

### Directory Structure

```text
.
├── cmd.go                # Cobra CLI command definition and program entry point (contains main function)
├── server.go             # MCP server core logic, tool definition and request handling
├── smh.go                # Aqara smart home platform API interface wrapper
├── middleware.go         # Middleware: user authentication, timeout control, exception recovery
├── config.go             # Global configuration management and environment variable handling
├── go.mod                # Go module dependency management file
├── go.sum                # Go module dependency checksum file
├── readme/               # README documentation and image resources
│   ├── img/              # Image resource directory
│   └── *.md              # Multi-language README files
├── LICENSE               # MIT open source license
└── README.md             # Main project documentation
```

### Core Files Description

- **`cmd.go`**: Cobra framework-based CLI implementation, defining `run stdio` and `run http` startup modes and main entry function
- **`server.go`**: MCP server core implementation, responsible for tool registration, request handling, and protocol support
- **`smh.go`**: Aqara smart home platform API wrapper layer, providing device control, authentication, and multi-home support
- **`middleware.go`**: Request processing middleware, providing authentication verification, timeout control, and exception handling
- **`config.go`**: Global configuration management, responsible for environment variable handling and API configuration

## Development Guide

Contributions are welcome through submitting Issues or Pull Requests!

Before submitting code, please ensure:

1. Code follows Go language coding standards
2. Related MCP tools and interface definitions maintain consistency and clarity
3. Add or update unit tests to cover your changes
4. Update relevant documentation (such as this README) if necessary
5. Ensure your commit messages are clear and descriptive

## License

This project is licensed under the [MIT License](/LICENSE).

Copyright (c) 2025 Aqara-Copilot
