# Aqara MCP Server

English | [中文](/readme/README_CN.md) | [繁體中文](/readme/README_CHT.md) | [Français](/readme/README_FR.md) | [한국어](/readme/README_KR.md) | [Español](/readme/README_ES.md) | [日本語](/readme/README_JP.md) | [Deutsch](/readme/README_DE.md) | [Italiano](/readme/README_IT.md)

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/aqara/aqara-mcp-server)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![Release](https://img.shields.io/github/v/release/aqara/aqara-mcp-server)](https://github.com/aqara/aqara-mcp-server/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Aqara MCP Server is a smart home control service developed based on the [MCP (Model Context Protocol)](https://modelcontextprotocol.io/introduction). It allows any AI assistant or API that supports the MCP protocol (such as Claude, ChatGPT, Cursor, etc.) to interact with your Aqara smart home devices, enabling device control, status queries, scene execution, and more through natural language.

## Table of Contents

- [Aqara MCP Server](#aqara-mcp-server)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [How It Works](#how-it-works)
  - [Quick Start](#quick-start)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
    - [Aqara Account Authentication](#aqara-account-authentication)
    - [Configuration Example (Claude for Desktop)](#configuration-example-claude-for-desktop)
    - [Running the Service](#running-the-service)
  - [Available Tools](#available-tools)
    - [device\_control](#device_control)
    - [device\_query](#device_query)
    - [device\_status\_query](#device_status_query)
    - [device\_log\_query](#device_log_query)
    - [run\_scenes](#run_scenes)
    - [get\_scenes](#get_scenes)
    - [automation\_config](#automation_config)
    - [get\_homes](#get_homes)
    - [switch\_home](#switch_home)
  - [Project Structure](#project-structure)
    - [Core File Descriptions](#core-file-descriptions)
  - [Contributing](#contributing)
  - [License](#license)

## Features

- **Comprehensive Device Control**: Support for fine-grained control of Aqara smart devices including switches, brightness, color temperature, modes, and more.
- **Flexible Device Queries**: Query device lists and detailed status by room and device type.
- **Smart Scene Management**: Support for querying and executing user-preset smart home scenes.
- **Device History**: Query device historical status changes within specified time ranges.
- **Automation Configuration**: Support for configuring scheduled or delayed device control tasks.
- **Multi-Home Support**: Support for querying and switching between different homes under user accounts.
- **MCP Protocol Compatible**: Fully compliant with MCP protocol specifications, easy to integrate with various AI assistants.
- **Secure Authentication**: Uses login authorization + signature-based security authentication to protect user data and device security.
- **Cross-Platform**: Built with Go, can be compiled to executables for multiple platforms.
- **Easy to Extend**: Modular design allows easy addition of new tools and features.

## How It Works

Aqara MCP Server acts as a bridge between AI assistants and the Aqara smart home platform:

1. **AI Assistant (MCP Client)**: Users issue commands through AI assistants (e.g., "Turn on the living room lights").
2. **MCP Client**: Parses user commands and calls corresponding tools provided by Aqara MCP Server according to MCP protocol (e.g., `device_control`).
3. **Aqara MCP Server (This Project)**: Receives requests from clients, validates them, and calls the `smh.go` module.
4. **`smh.go` Module**: Uses configured Aqara credentials to communicate with Aqara cloud APIs for actual device operations or data queries.
5. **Response Flow**: Aqara cloud API returns results, which are passed back through Aqara MCP Server to the MCP client and finally presented to the user.

## Quick Start

### Prerequisites

- Go (version 1.24 or higher)
- Git (for building from source)
- Aqara account with connected smart devices

### Installation

You can choose to download pre-compiled executables or build from source.

**Option 1: Download Pre-compiled Version (Recommended)**

Visit the link below to download the latest executable package for your operating system.

[Releases Page](https://github.com/aqara/aqara-mcp-server/releases)

Extract and use directly.

**Option 2: Build from Source**

```bash
# Clone repository
git clone https://github.com/aqara/aqara-mcp-server.git
cd aqara-mcp-server

# Download dependencies
go mod tidy

# Build executable
go build -o aqara-mcp-server
```

After building, the `aqara-mcp-server` executable will be generated in the current directory.

### Aqara Account Authentication

To enable the MCP Server to access your Aqara account and control devices, you need to complete login authorization first.

Please visit the following address to complete login authorization:
[https://cdn.aqara.com/app/mcpserver/login.html](https://cdn.aqara.com/app/mcpserver/login.html)

After successful login, you will obtain necessary authentication information (such as `token`, `region`), which will be used in subsequent configuration steps.

**Please keep this information secure, especially the `token` - do not share it with others.**

### Configuration Example (Claude for Desktop)

Different MCP clients have slightly different configuration methods. Here's an example of how to configure Claude for Desktop to use this MCP Server:

1. Open Claude for Desktop Settings.
2. Switch to the Developer tab.
3. Click Edit Config to open the configuration file with a text editor.

   ![](/readme/img/setting0.png)
   ![](/readme/img/setting1.png)

4. Add the configuration information from the "Login Success Page" to the client's configuration file (claude_desktop_config.json). Configuration example:

   ![](/readme/img/config.png)

**Configuration Notes:**
- `command`: Full path to your downloaded or built `aqara-mcp-server` executable
- `args`: Use `["run", "stdio"]` to start stdio transport mode
- `env`: Environment variable configuration
  - `token`: Access token obtained from the Aqara login page
  - `region`: Your Aqara account region (e.g., CN, US, EU, etc.)

### Running the Service

Restart Claude for Desktop. Then you can use conversations to call tools provided by the MCP Server for device control, device queries, and other operations.

![](/readme/img/claude.png)

**Other MCP Client Configuration**

For other MCP protocol-supporting clients (such as Claude, ChatGPT, Cursor, etc.), the configuration is similar:
- Ensure the client supports MCP protocol
- Configure executable file path and startup parameters
- Set environment variables `token` and `region`
- Choose appropriate transport protocol (stdio recommended)

**SSE Mode (Optional)**

If you need to use SSE (Server-Sent Events) mode, you can start it like this:

```bash
# Use default port 8080
./aqara-mcp-server run sse

# Or specify custom host and port
./aqara-mcp-server run sse --host localhost --port 9000
```

Then use `["run", "sse"]` parameters in client configuration.

## Available Tools

MCP clients can interact with Aqara smart home devices by calling these tools.

### device_control

- **Description**: Control the status or properties of smart home devices (e.g., on/off, temperature, brightness, color, color temperature, etc.).
- **Parameters**:
  - `endpoint_ids` (Array<Integer>, required): List of device IDs to control.
  - `control_params` (Object, required): Control parameter object containing specific operations.
    - `action` (String, required): Action to execute. Examples: `"on"`, `"off"`, `"set"`, `"up"`, `"down"`, `"cooler"`, `"warmer"`.
    - `attribute` (String, required): Device attribute to control. Examples: `"on_off"`, `"brightness"`, `"color_temperature"`, `"ac_mode"`.
    - `value` (String | Number, optional): Target value (required when action is "set").
    - `unit` (String, optional): Unit of the value (e.g., `"%"`, `"K"`, `"℃"`).
- **Returns**: (String) Operation result message for device control.

### device_query

- **Description**: Get device list by specified location (room) and device type (does not include real-time status information, only lists devices and their IDs).
- **Parameters**:
  - `positions` (Array<String>, optional): List of room names. If empty array or not provided, queries all rooms.
  - `device_types` (Array<String>, optional): List of device types. Examples: `"Light"`, `"WindowCovering"`, `"AirConditioner"`, `"Button"`, etc. If empty array or not provided, queries all types.
- **Returns**: (String) Markdown-formatted device list including device names and IDs.

### device_status_query

- **Description**: Get current status information of devices (for querying status-related attributes like color, brightness, switches, etc.).
- **Parameters**:
  - `positions` (Array<String>, optional): List of room names. If empty array or not provided, queries all rooms.
  - `device_types` (Array<String>, optional): List of device types. Same options as `device_query`. If empty array or not provided, queries all types.
- **Returns**: (String) Markdown-formatted device status information.

### device_log_query

- **Description**: Query device logs.
- **Parameters**:
  - `endpoint_ids` (Array<Integer>, required): List of device IDs to query history for.
  - `start_datetime` (String, optional): Query start time in format `YYYY-MM-DD HH:MM:SS` (e.g., `"2023-05-16 12:00:00"`).
  - `end_datetime` (String, optional): Query end time in format `YYYY-MM-DD HH:MM:SS`.
  - `attribute` (String, optional): Specific device attribute name to query (e.g., `on_off`, `brightness`). If not provided, queries all recorded attributes for the device.
- **Returns**: (String) Markdown-formatted device historical status information. (Note: Current implementation may show "This feature will be available soon.", indicating the feature is pending completion.)

### run_scenes

- **Description**: Execute specified scenes by scene ID.
- **Parameters**:
  - `scenes` (Array<Integer>, required): List of scene IDs to execute.
- **Returns**: (String) Scene execution result message.

### get_scenes

- **Description**: Query all scenes in the user's home, or scenes within specified rooms.
- **Parameters**:
  - `positions` (Array<String>, optional): List of room names. If empty array or not provided, queries scenes for the entire home.
- **Returns**: (String) Markdown-formatted scene information.

### automation_config

- **Description**: Configure scheduled or delayed device control tasks.
- **Parameters**:
  - `scheduled_time` (String, required): Set time point (for delayed tasks, converted based on current time), format `YYYY-MM-DD HH:MM:SS` (e.g., `"2025-05-16 12:12:12"`).
  - `endpoint_ids` (Array<Integer>, required): List of device IDs for scheduled control.
  - `control_params` (Object, required): Device control parameters using the same format as `device_control` tool (including action, attribute, value, etc.).
- **Returns**: (String) Automation configuration result message.

### get_homes

- **Description**: Get all home lists under the user account.
- **Parameters**: None.
- **Returns**: (String) Comma-separated list of home names. Returns empty string or appropriate message if no data.

### switch_home

- **Description**: Switch the user's current operating home. After switching, subsequent device queries, controls, and other operations will target the newly switched home.
- **Parameters**:
  - `home_name` (String, required): Target home name (should come from the available list provided by `get_homes` tool).
- **Returns**: (String) Switch operation result message.

## Project Structure

```
.
├── cmd.go                # Cobra CLI command definitions and program entry point (contains main function)
├── server.go             # MCP server core logic, tool definitions and request handling
├── smh.go                # Aqara smart home platform API interface wrapper
├── middleware.go         # Middleware: user authentication, timeout control, exception recovery
├── config.go             # Global configuration management and environment variable handling
├── go.mod                # Go module dependency management file
├── go.sum                # Go module dependency checksum file
├── img/                  # Image resources used in README documentation
├── LICENSE               # MIT open source license
└── README.md             # Project documentation
```

### Core File Descriptions

- **`cmd.go`**: Cobra framework-based CLI implementation, defining `run stdio` and `run sse` startup modes and main entry function
- **`server.go`**: MCP server core implementation, responsible for tool registration, request handling, and protocol support
- **`smh.go`**: Aqara smart home platform API wrapper layer, providing device control, authentication, and multi-home support
- **`middleware.go`**: Request processing middleware providing authentication validation, timeout control, and exception handling
- **`config.go`**: Global configuration management, responsible for environment variable handling and API configuration

## Contributing

Welcome to contribute to the project by submitting Issues or Pull Requests!

Before submitting code, please ensure:
1. Code follows Go language coding standards.
2. Related MCP tools and prompt interface definitions maintain consistency and clarity.
3. Add or update unit tests to cover your changes.
4. Update relevant documentation (such as this README) if necessary.
5. Ensure your commit messages are clear and descriptive.

## License

This project is licensed under the [MIT License](./LICENSE).
Copyright (c) 2025 Aqara-Copliot