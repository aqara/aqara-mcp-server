<div align="center" style="display: flex; align-items: center; justify-content: center; ">

  <img src="/readme/img/logo.png" alt="Aqara Logo" height="120">
  <h1>Aqara MCP Server</h1>

</div>

<div align="center">

[English](/readme/README.md) | [中文](/readme/README_CN.md) | [繁體中文](/readme/README_CHT.md) | [Français](/readme/README_FR.md) | 한국어 | [Español](/readme/README_ES.md) | [日本語](/readme/README_JP.md) | [Deutsch](/readme/README_DE.md) | [Italiano](/readme/README_IT.md)

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/aqara/aqara-mcp-server)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![Release](https://img.shields.io/github/v/release/aqara/aqara-mcp-server)](https://github.com/aqara/aqara-mcp-server/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

Aqara MCP Server는 [MCP (Model Context Protocol)](https://modelcontextprotocol.io/introduction) 프로토콜을 기반으로 개발된 지능형 홈 제어 서비스입니다. MCP 프로토콜을 지원하는 모든 AI 어시스턴트나 API(Claude, Cursor 등)가 Aqara 스마트 홈 디바이스와 상호작용할 수 있도록 하여, 자연어를 통한 디바이스 제어, 상태 조회, 시나리오 실행 등의 기능을 구현합니다.

## 목차

- [목차](#목차)
- [주요 기능](#주요-기능)
- [작동 원리](#작동-원리)
- [빠른 시작](#빠른-시작)
  - [사전 요구사항](#사전-요구사항)
  - [설치](#설치)
    - [방법 1: 사전 컴파일된 버전 다운로드 (권장)](#방법-1-사전-컴파일된-버전-다운로드-권장)
    - [방법 2: 소스코드로부터 빌드](#방법-2-소스코드로부터-빌드)
  - [Aqara 계정 인증](#aqara-계정-인증)
  - [클라이언트 구성](#클라이언트-구성)
    - [Claude for Desktop 구성 예시](#claude-for-desktop-구성-예시)
    - [구성 매개변수 설명](#구성-매개변수-설명)
    - [기타 MCP 클라이언트](#기타-mcp-클라이언트)
  - [서비스 시작](#서비스-시작)
    - [표준 모드 (권장)](#표준-모드-권장)
    - [HTTP 모드 (`곧 지원 예정`)](#http-모드-곧-지원-예정)
- [API 도구 설명](#api-도구-설명)
  - [디바이스 제어](#디바이스-제어)
    - [device\_control](#device_control)
  - [디바이스 조회](#디바이스-조회)
    - [device\_query](#device_query)
    - [device\_status\_query](#device_status_query)
    - [device\_log\_query](#device_log_query)
  - [시나리오 관리](#시나리오-관리)
    - [get\_scenes](#get_scenes)
    - [run\_scenes](#run_scenes)
  - [홈 관리](#홈-관리)
    - [get\_homes](#get_homes)
    - [switch\_home](#switch_home)
  - [자동화 구성](#자동화-구성)
    - [automation\_config](#automation_config)
- [프로젝트 구조](#프로젝트-구조)
  - [디렉토리 구조](#디렉토리-구조)
  - [핵심 파일 설명](#핵심-파일-설명)
- [개발 가이드](#개발-가이드)
- [라이선스](#라이선스)

## 주요 기능

- ✨ **포괄적인 디바이스 제어**: Aqara 스마트 디바이스의 온/오프, 밝기, 색온도, 모드 등 다양한 속성에 대한 정밀 제어 지원
- 🔍 **유연한 디바이스 조회**: 방별, 디바이스 타입별 디바이스 목록 및 상세 상태 조회 기능
- 🎬 **지능형 시나리오 관리**: 사용자 사전 설정 스마트 홈 시나리오의 조회 및 실행 지원
- 📈 **디바이스 이력 기록**: 지정된 시간 범위 내 디바이스의 상태 변경 이력 조회
- ⏰ **자동화 구성**: 타이머 또는 지연 디바이스 제어 작업 구성 지원
- 🏠 **다중 홈 지원**: 사용자 계정 하위의 다양한 홈 조회 및 전환 지원
- 🔌 **MCP 프로토콜 호환**: MCP 프로토콜 표준 완전 준수로 다양한 AI 어시스턴트와의 손쉬운 통합
- 🔐 **보안 인증 메커니즘**: 로그인 인증 + 서명 기반 보안 인증으로 사용자 데이터 및 디바이스 보안 보호
- 🌐 **크로스 플랫폼 실행**: Go 언어 기반 개발로 다중 플랫폼 실행 파일 컴파일 가능
- 🔧 **확장 용이성**: 모듈형 설계로 새로운 도구와 기능 추가 용이

## 작동 원리

Aqara MCP Server는 AI 어시스턴트와 Aqara 스마트 홈 플랫폼 간의 브리지 역할을 수행합니다:

```mermaid
graph LR
    A[AI 어시스턴트] --> B[MCP 클라이언트]
    B --> C[Aqara MCP Server]
    C --> D[Aqara 클라우드 API]
    D --> E[스마트 디바이스]
```

1. **AI 어시스턴트**: 사용자가 AI 어시스턴트를 통해 명령을 내립니다 (예: "거실 불을 켜줘")
2. **MCP 클라이언트**: 사용자 명령을 파싱하고 MCP 프로토콜에 따라 Aqara MCP Server가 제공하는 해당 도구를 호출합니다 (예: `device_control`)
3. **Aqara MCP Server (본 프로젝트)**: 클라이언트로부터의 요청을 수신하고, 구성된 Aqara 자격 증명을 사용하여 Aqara 클라우드 API와 통신하며, 실제 디바이스 작업이나 데이터 조회를 실행합니다
4. **응답 흐름**: Aqara 클라우드 API가 결과를 반환하고, Aqara MCP Server를 통해 MCP 클라이언트로전달되어 최종적으로 사용자에게 표시됩니다

## 빠른 시작

### 사전 요구사항

- **Go** (버전 1.24 이상) - 소스코드로부터 빌드 시에만 필요
- **Git** (소스코드 빌드용) - 선택사항
- **Aqara 계정** 및 연결된 스마트 디바이스
- **MCP 프로토콜 지원 클라이언트** (Claude for Desktop, Cursor 등)

### 설치

사전 컴파일된 실행 파일을 다운로드하거나 소스코드로부터 빌드할 수 있습니다.

#### 방법 1: 사전 컴파일된 버전 다운로드 (권장)

GitHub Releases 페이지를 방문하여 운영체제에 맞는 최신 실행 파일을 다운로드하세요:

**📥 [Releases 페이지에서 다운로드](https://github.com/aqara/aqara-mcp-server/releases)**

해당 플랫폼의 압축 파일을 다운로드한 후 압축을 해제하면 바로 사용할 수 있습니다.

#### 방법 2: 소스코드로부터 빌드

```bash
# 리포지토리 클론
git clone https://github.com/aqara/aqara-mcp-server.git
cd aqara-mcp-server

# 종속성 다운로드
go mod tidy

# 실행 파일 빌드
go build -o aqara-mcp-server
```

빌드 완료 후, 현재 디렉토리에 `aqara-mcp-server` 실행 파일이 생성됩니다.

### Aqara 계정 인증

MCP Server가 Aqara 계정에 액세스하고 디바이스를 제어할 수 있도록 하려면 먼저 로그인 인증을 완료해야 합니다.

다음 주소를 방문하여 로그인 인증을 완료하세요:
**🔗 [https://cdn.aqara.com/app/mcpserver/login.html](https://cdn.aqara.com/app/mcpserver/login.html)**

로그인 성공 후, 필요한 인증 정보(`token`, `region` 등)를 획득하게 되며, 이 정보는 후속 구성 단계에서 사용됩니다.

> ⚠️ **보안 알림**: `token` 정보를 안전하게 보관하고 타인에게 노출되지 않도록 주의하세요.

### 클라이언트 구성

MCP 클라이언트마다 구성 방법이 약간 다릅니다. 다음은 Claude for Desktop에서 이 MCP Server를 사용하도록 구성하는 예시입니다:

#### Claude for Desktop 구성 예시

1. **Claude for Desktop 설정(Settings) 열기**

    ![Claude Open Setting](/readme/img/opening_setting.png)

2. **개발자(Developer) 탭으로 전환한 후 구성 편집(Edit Config)을 클릭하여 텍스트 에디터로 구성 파일 열기**

    ![Claude Edit Configuration](/readme/img/edit_config.png)

3. **"로그인 성공 페이지"의 구성 정보를 클라이언트 구성 파일 `claude_desktop_config.json`에 추가**

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

#### 구성 매개변수 설명

- `command`: 다운로드하거나 빌드한 `aqara-mcp-server` 실행 파일의 전체 경로
- `args`: `["run", "stdio"]`를 사용하여 stdio 전송 모드 시작
- `env`: 환경 변수 구성
  - `token`: Aqara 로그인 페이지에서 획득한 액세스 토큰
  - `region`: Aqara 계정이 속한 지역 (지원 지역: CN, US, EU, KR, SG, RU)

#### 기타 MCP 클라이언트

MCP 프로토콜을 지원하는 다른 클라이언트(ChatGPT, Cursor 등)의 구성 방법은 유사합니다:

- 클라이언트가 MCP 프로토콜을 지원하는지 확인
- 실행 파일 경로 및 시작 매개변수 구성
- 환경 변수 `token` 및 `region` 설정
- 적절한 전송 프로토콜 선택 (`stdio` 권장)

### 서비스 시작

#### 표준 모드 (권장)

Claude for Desktop을 재시작합니다. 그러면 자연어를 통해 디바이스 제어, 디바이스 조회, 시나리오 실행 등의 작업을 수행할 수 있습니다.

대화 예시:

- "거실 불을 켜줘"
- "침실 에어컨을 냉방 모드로 설정하고 온도를 24도로 맞춰줘"
- "모든 방의 디바이스 목록을 확인해줘"
- "수면 시나리오를 실행해줘"

![Claude Chat Example](/readme/img/claude.png)

#### HTTP 모드 (`곧 지원 예정`)

## API 도구 설명

MCP 클라이언트는 이러한 도구를 호출하여 Aqara 스마트 홈 디바이스와 상호작용할 수 있습니다.

### 디바이스 제어

#### device_control

스마트 홈 디바이스의 상태나 속성(온/오프, 온도, 밝기, 색상, 색온도 등)을 제어합니다.

**매개변수:**

- `endpoint_ids` _(Array\<Integer\>, 필수)_: 제어할 디바이스 ID 목록
- `control_params` _(Object, 필수)_: 제어 매개변수 객체, 구체적인 작업 포함:
  - `action` _(String, 필수)_: 실행할 작업 (`"on"`, `"off"`, `"set"`, `"up"`, `"down"`, `"cooler"`, `"warmer"`)
  - `attribute` _(String, 필수)_: 제어할 디바이스 속성 (`"on_off"`, `"brightness"`, `"color_temperature"`, `"ac_mode"`)
  - `value` _(String | Number, 선택)_: 목표값 (action이 "set"일 때 필수)
  - `unit` _(String, 선택)_: 값의 단위 (`"%"`, `"K"`, `"℃"`)

**반환:** 디바이스 제어 작업 결과 메시지

### 디바이스 조회

#### device_query

지정된 위치(방)와 디바이스 타입에 따라 디바이스 목록을 조회합니다(실시간 상태 정보 미포함).

**매개변수:**

- `positions` _(Array\<String\>, 선택)_: 방 이름 목록. 빈 배열은 모든 방 조회를 의미
- `device_types` _(Array\<String\>, 선택)_: 디바이스 타입 목록 (`"Light"`, `"WindowCovering"`, `"AirConditioner"`, `"Button"`). 빈 배열은 모든 타입 조회를 의미

**반환:** 디바이스 이름과 ID를 포함한 Markdown 형식의 디바이스 목록

#### device_status_query

디바이스의 현재 상태 정보를 조회합니다(색상, 밝기, 온/오프 등의 실시간 상태 정보 조회용).

**매개변수:**

- `positions` _(Array\<String\>, 선택)_: 방 이름 목록. 빈 배열은 모든 방 조회를 의미
- `device_types` _(Array\<String\>, 선택)_: 디바이스 타입 목록. 선택 가능한 값은 `device_query`와 동일. 빈 배열은 모든 타입 조회를 의미

**반환:** Markdown 형식의 디바이스 상태 정보

#### device_log_query

디바이스의 이력 로그 정보를 조회합니다.

**매개변수:**

- `endpoint_ids` _(Array\<Integer\>, 필수)_: 이력 기록을 조회할 디바이스 ID 목록
- `start_datetime` _(String, 선택)_: 조회 시작 시간, 형식: `YYYY-MM-DD HH:MM:SS` (예: `"2023-05-16 12:00:00"`)
- `end_datetime` _(String, 선택)_: 조회 종료 시간, 형식: `YYYY-MM-DD HH:MM:SS`
- `attributes` _(Array\<String\>, 선택)_: 조회할 디바이스 속성 이름 목록 (`["on_off", "brightness"]`). 제공되지 않으면 모든 기록된 속성 조회

**반환:** Markdown 형식의 디바이스 이력 상태 정보

### 시나리오 관리

#### get_scenes

사용자 홈의 모든 시나리오 또는 지정된 방 내의 시나리오를 조회합니다.

**매개변수:**

- `positions` _(Array\<String\>, 선택)_: 방 이름 목록. 빈 배열은 전체 홈의 시나리오 조회를 의미

**반환:** Markdown 형식의 시나리오 정보

#### run_scenes

시나리오 ID를 기반으로 지정된 시나리오를 실행합니다.

**매개변수:**

- `scenes` _(Array\<Integer\>, 필수)_: 실행할 시나리오 ID 목록

**반환:** 시나리오 실행 결과 메시지

### 홈 관리

#### get_homes

사용자 계정 하위의 모든 홈 목록을 조회합니다.

**매개변수:** 없음

**반환:** 쉼표로 구분된 홈 이름 목록. 데이터가 없으면 빈 문자열 또는 해당 안내 정보 반환

#### switch_home

사용자의 현재 작업 홈을 전환합니다. 전환 후 후속 디바이스 조회, 제어 등의 작업은 새로 전환된 홈을 대상으로 합니다.

**매개변수:**

- `home_name` _(String, 필수)_: 대상 홈의 이름

**반환:** 전환 작업 결과 메시지

### 자동화 구성

#### automation_config

자동화 구성 (현재 타이머 또는 지연 디바이스 제어 작업만 지원).

**매개변수:**

- `scheduled_time` _(String, 필수)_: 타이머 실행 시점, 표준 Crontab 형식 `"분 시 일 월 주"` 사용. 예: `"30 14 * * *"` (매일 14:30 실행), `"0 9 * * 1"` (매주 월요일 9:00 실행)
- `endpoint_ids` _(Array\<Integer\>, 필수)_: 타이머 제어할 디바이스 ID 목록
- `control_params` _(Object, 필수)_: 디바이스 제어 매개변수, `device_control` 도구와 동일한 형식 사용 (action, attribute, value 등 포함)
- `task_name` _(String, 필수)_: 자동화 작업의 이름 또는 설명 (식별 및 관리용)
- `execution_once` _(Boolean, 선택)_: 한 번만 실행할지 여부
  - `true`: 지정된 시간에 한 번만 작업 실행 (기본값)
  - `false`: 주기적으로 반복 실행 (매일, 매주 등)

**반환:** 자동화 구성 결과 메시지

## 프로젝트 구조

### 디렉토리 구조

```text
.
├── cmd.go                # Cobra CLI 명령 정의 및 프로그램 진입점 (main 함수 포함)
├── server.go             # MCP 서버 핵심 로직, 도구 정의 및 요청 처리
├── smh.go                # Aqara 스마트 홈 플랫폼 API 인터페이스 래핑
├── middleware.go         # 미들웨어: 사용자 인증, 타임아웃 제어, 예외 복구
├── config.go             # 전역 구성 관리 및 환경 변수 처리
├── go.mod                # Go 모듈 종속성 관리 파일
├── go.sum                # Go 모듈 종속성 체크섬 파일
├── readme/               # README 문서 및 이미지 리소스
│   ├── img/              # 이미지 리소스 디렉토리
│   └── *.md              # 다국어 README 파일
├── LICENSE               # MIT 오픈소스 라이선스
└── README.md             # 프로젝트 메인 문서
```

### 핵심 파일 설명

- **`cmd.go`**: Cobra 프레임워크 기반 CLI 구현, `run stdio` 및 `run http` 시작 모드와 메인 진입 함수 정의
- **`server.go`**: MCP 서버 핵심 구현, 도구 등록, 요청 처리 및 프로토콜 지원 담당
- **`smh.go`**: Aqara 스마트 홈 플랫폼 API 래핑 레이어, 디바이스 제어, 인증 및 다중 홈 지원 제공
- **`middleware.go`**: 요청 처리 미들웨어, 인증 검증, 타임아웃 제어 및 예외 처리 제공
- **`config.go`**: 전역 구성 관리, 환경 변수 처리 및 API 구성 담당

## 개발 가이드

Issue 제출이나 Pull Request를 통한 프로젝트 기여를 환영합니다!

코드 제출 전 다음 사항을 확인해 주세요:

1. Go 언어 코딩 표준 준수
2. 관련 MCP 도구 및 인터페이스 정의의 일관성과 명확성 유지
3. 변경 사항을 다루는 단위 테스트 추가 또는 업데이트
4. 필요시 관련 문서(본 README 등) 업데이트
5. 커밋 메시지가 명확하고 이해하기 쉬운지 확인

**🌟 이 프로젝트가 도움이 되셨다면 Star를 눌러주세요!**

**🤝 커뮤니티에 참여하여 스마트 홈을 더욱 지능적으로 만들어 나가요!**

## 라이선스

본 프로젝트는 [MIT License](/LICENSE) 하에 라이선스됩니다.

---

Copyright (c) 2025 Aqara-Copilot
