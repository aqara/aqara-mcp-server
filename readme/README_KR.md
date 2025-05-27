# Aqara MCP Server

[English](/readme/README.md) | [中文](/readme/README_CN.md) | [繁體中文](/readme/README_CHT.md) | [Français](/readme/README_FR.md) | 한국어 | [Español](/readme/README_ES.md) | [日本語](/readme/README_JP.md) | [Deutsch](/readme/README_DE.md) | [Italiano](/readme/README_IT.md)

[![빌드 상태](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/aqara/aqara-mcp-server)
[![Go 버전](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![릴리스](https://img.shields.io/github/v/release/aqara/aqara-mcp-server)](https://github.com/aqara/aqara-mcp-server/releases)
[![라이선스: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Aqara MCP Server는 [MCP (Model Context Protocol)](https://modelcontextprotocol.io/introduction) 프로토콜을 기반으로 개발된 스마트 홈 제어 서비스입니다. MCP 프로토콜을 지원하는 모든 AI 어시스턴트나 API(Claude, ChatGPT, Cursor 등)가 Aqara 스마트 홈 기기와 상호작용할 수 있게 하여, 자연어를 통한 기기 제어, 상태 조회, 시나리오 실행 등의 기능을 제공합니다.

## 목차

- [Aqara MCP Server](#aqara-mcp-server)
  - [목차](#목차)
  - [특징](#특징)
  - [작동 원리](#작동-원리)
  - [빠른 시작](#빠른-시작)
    - [사전 요구사항](#사전-요구사항)
    - [설치](#설치)
    - [Aqara 계정 인증](#aqara-계정-인증)
    - [설정 예시 (Claude for Desktop)](#설정-예시-claude-for-desktop)
    - [서비스 실행](#서비스-실행)
  - [사용 가능한 도구](#사용-가능한-도구)
    - [device\_control](#device_control)
    - [device\_query](#device_query)
    - [device\_status\_query](#device_status_query)
    - [device\_log\_query](#device_log_query)
    - [run\_scenes](#run_scenes)
    - [get\_scenes](#get_scenes)
    - [automation\_config](#automation_config)
    - [get\_homes](#get_homes)
    - [switch\_home](#switch_home)
  - [프로젝트 구조](#프로젝트-구조)
    - [핵심 파일 설명](#핵심-파일-설명)
  - [기여하기](#기여하기)
  - [라이선스](#라이선스)

## 특징

- **포괄적인 기기 제어**: 스위치, 밝기, 색온도, 모드 등 Aqara 스마트 기기의 다양한 속성에 대한 세밀한 제어 지원.
- **유연한 기기 조회**: 방별, 기기 유형별로 기기 목록과 상세 상태 조회 가능.
- **스마트 시나리오 관리**: 사용자 사전 설정 스마트 홈 시나리오 조회 및 실행 지원.
- **기기 히스토리**: 지정된 시간 범위 내 기기의 과거 상태 변경 기록 조회.
- **자동화 설정**: 예약 또는 지연 기기 제어 작업 설정 지원.
- **다중 홈 지원**: 사용자 계정 하의 다른 홈 조회 및 전환 지원.
- **MCP 프로토콜 호환**: MCP 프로토콜 사양을 완벽하게 준수하여 다양한 AI 어시스턴트와 쉽게 통합.
- **보안 인증 메커니즘**: 로그인 인증 + 서명 기반 보안 인증을 사용하여 사용자 데이터와 기기 보안 보호.
- **크로스 플랫폼**: Go 언어로 개발되어 여러 플랫폼용 실행 파일로 컴파일 가능.
- **확장 용이**: 모듈화 설계로 새로운 도구와 기능을 쉽게 추가 가능.

## 작동 원리

Aqara MCP Server는 AI 어시스턴트와 Aqara 스마트 홈 플랫폼 간의 다리 역할을 합니다:

1. **AI 어시스턴트 (MCP 클라이언트)**: 사용자가 AI 어시스턴트를 통해 명령어를 발행 (예: "거실 불을 켜줘").
2. **MCP 클라이언트**: 사용자 명령어를 파싱하고 MCP 프로토콜에 따라 Aqara MCP Server가 제공하는 해당 도구를 호출 (예: `device_control`).
3. **Aqara MCP Server (본 프로젝트)**: 클라이언트로부터 요청을 받아 검증 후 `smh.go` 모듈을 호출.
4. **`smh.go` 모듈**: 설정된 Aqara 자격 증명을 사용하여 Aqara 클라우드 API와 통신하여 실제 기기 조작이나 데이터 조회 수행.
5. **응답 플로우**: Aqara 클라우드 API가 결과를 반환하면, Aqara MCP Server를 통해 MCP 클라이언트로 전달되어 최종적으로 사용자에게 제시.

## 빠른 시작

### 사전 요구사항

- Go (버전 1.24 이상)
- Git (소스에서 빌드용)
- 연결된 스마트 기기가 있는 Aqara 계정

### 설치

미리 컴파일된 실행 파일을 다운로드하거나 소스에서 빌드할 수 있습니다.

**옵션 1: 미리 컴파일된 버전 다운로드 (권장)**

아래 링크를 방문하여 운영 체제에 맞는 최신 실행 파일 패키지를 다운로드하세요.

[릴리스 페이지](https://github.com/aqara/aqara-mcp-server/releases)

압축을 해제한 후 바로 사용할 수 있습니다.

**옵션 2: 소스에서 빌드**

```bash
# 저장소 클론
git clone https://github.com/aqara/aqara-mcp-server.git
cd aqara-mcp-server

# 종속성 다운로드
go mod tidy

# 실행 파일 빌드
go build -o aqara-mcp-server
```

빌드 완료 후 현재 디렉토리에 `aqara-mcp-server` 실행 파일이 생성됩니다.

### Aqara 계정 인증

MCP Server가 Aqara 계정에 액세스하고 기기를 제어할 수 있도록 하려면 먼저 로그인 인증을 완료해야 합니다.

다음 주소를 방문하여 로그인 인증을 완료하세요:
[https://cdn.aqara.com/app/mcpserver/login.html](https://cdn.aqara.com/app/mcpserver/login.html)

로그인 성공 후 필요한 인증 정보(`token`, `region` 등)를 얻게 되며, 이는 후속 설정 단계에서 사용됩니다.

**이 정보를 안전하게 보관하시고, 특히 `token`은 다른 사람과 공유하지 마세요.**

### 설정 예시 (Claude for Desktop)

각 MCP 클라이언트의 설정 방법은 약간씩 다릅니다. 다음은 Claude for Desktop에서 이 MCP Server를 사용하도록 설정하는 예시입니다:

1. Claude for Desktop의 설정(Settings)을 엽니다.
2. 개발자(Developer) 탭으로 전환합니다.
3. Edit Config를 클릭하여 텍스트 에디터로 설정 파일을 엽니다.

   ![](/readme/img/setting0.png)
   ![](/readme/img/setting1.png)

4. "로그인 성공 페이지"의 설정 정보를 클라이언트의 설정 파일(claude_desktop_config.json)에 추가합니다. 설정 예시:

   ![](/readme/img/config.png)

**설정 설명:**
- `command`: 다운로드하거나 빌드한 `aqara-mcp-server` 실행 파일의 전체 경로
- `args`: `["run", "stdio"]`를 사용하여 stdio 전송 모드로 시작
- `env`: 환경 변수 설정
  - `token`: Aqara 로그인 페이지에서 얻은 액세스 토큰
  - `region`: Aqara 계정 지역 (예: CN, US, EU 등)

### 서비스 실행

Claude for Desktop을 재시작합니다. 그 후 대화를 통해 MCP Server가 제공하는 도구를 호출하여 기기 제어, 기기 조회 등의 작업을 수행할 수 있습니다.

![](/readme/img/claude.png)

**다른 MCP 클라이언트 설정**

MCP 프로토콜을 지원하는 다른 클라이언트(Claude, ChatGPT, Cursor 등)의 경우 설정이 유사합니다:
- 클라이언트가 MCP 프로토콜을 지원하는지 확인
- 실행 파일 경로와 시작 매개변수 설정
- 환경 변수 `token`과 `region` 설정
- 적절한 전송 프로토콜 선택 (stdio 권장)

**SSE 모드 (선택사항)**

SSE (Server-Sent Events) 모드를 사용해야 하는 경우 다음과 같이 시작할 수 있습니다:

```bash
# 기본 포트 8080 사용
./aqara-mcp-server run sse

# 또는 커스텀 호스트와 포트 지정
./aqara-mcp-server run sse --host localhost --port 9000
```

그 후 클라이언트 설정에서 `["run", "sse"]` 매개변수를 사용합니다.

## 사용 가능한 도구

MCP 클라이언트는 이러한 도구를 호출하여 Aqara 스마트 홈 기기와 상호작용할 수 있습니다.

### device_control

- **설명**: 스마트 홈 기기의 상태나 속성 제어 (예: 온/오프, 온도, 밝기, 색상, 색온도 등).
- **매개변수**:
  - `endpoint_ids` (Array<Integer>, 필수): 제어할 기기 ID 목록.
  - `control_params` (Object, 필수): 특정 작업을 포함하는 제어 매개변수 객체.
    - `action` (String, 필수): 실행할 작업. 예시: `"on"`, `"off"`, `"set"`, `"up"`, `"down"`, `"cooler"`, `"warmer"`.
    - `attribute` (String, 필수): 제어할 기기 속성. 예시: `"on_off"`, `"brightness"`, `"color_temperature"`, `"ac_mode"`.
    - `value` (String | Number, 선택사항): 목표 값 (action이 "set"일 때 필수).
    - `unit` (String, 선택사항): 값의 단위 (예: `"%"`, `"K"`, `"℃"`).
- **반환**: (String) 기기 제어 작업 결과 메시지.

### device_query

- **설명**: 지정된 위치(방)와 기기 유형으로 기기 목록 조회 (실시간 상태 정보는 포함하지 않고 기기와 ID만 나열).
- **매개변수**:
  - `positions` (Array<String>, 선택사항): 방 이름 목록. 빈 배열이거나 제공되지 않으면 모든 방 조회.
  - `device_types` (Array<String>, 선택사항): 기기 유형 목록. 예시: `"Light"`, `"WindowCovering"`, `"AirConditioner"`, `"Button"` 등. 빈 배열이거나 제공되지 않으면 모든 유형 조회.
- **반환**: (String) 기기 이름과 ID를 포함한 Markdown 형식의 기기 목록.

### device_status_query

- **설명**: 기기의 현재 상태 정보 조회 (색상, 밝기, 스위치 등 상태 관련 속성 조회용).
- **매개변수**:
  - `positions` (Array<String>, 선택사항): 방 이름 목록. 빈 배열이거나 제공되지 않으면 모든 방 조회.
  - `device_types` (Array<String>, 선택사항): 기기 유형 목록. `device_query`와 동일한 옵션. 빈 배열이거나 제공되지 않으면 모든 유형 조회.
- **반환**: (String) Markdown 형식의 기기 상태 정보.

### device_log_query

- **설명**: 기기 로그 조회.
- **매개변수**:
  - `endpoint_ids` (Array<Integer>, 필수): 히스토리를 조회할 기기 ID 목록.
  - `start_datetime` (String, 선택사항): `YYYY-MM-DD HH:MM:SS` 형식의 조회 시작 시간 (예: `"2023-05-16 12:00:00"`).
  - `end_datetime` (String, 선택사항): `YYYY-MM-DD HH:MM:SS` 형식의 조회 종료 시간.
  - `attribute` (String, 선택사항): 조회할 특정 기기 속성 이름 (예: `on_off`, `brightness`). 제공되지 않으면 해당 기기의 모든 기록된 속성을 조회.
- **반환**: (String) Markdown 형식의 기기 히스토리 상태 정보. (참고: 현재 구현에서는 "This feature will be available soon."을 표시할 수 있으며, 이는 기능이 완성 대기 중임을 나타냅니다.)

### run_scenes

- **설명**: 시나리오 ID로 지정된 시나리오 실행.
- **매개변수**:
  - `scenes` (Array<Integer>, 필수): 실행할 시나리오 ID 목록.
- **반환**: (String) 시나리오 실행 결과 메시지.

### get_scenes

- **설명**: 사용자 홈의 모든 시나리오 또는 지정된 방 내 시나리오 조회.
- **매개변수**:
  - `positions` (Array<String>, 선택사항): 방 이름 목록. 빈 배열이거나 제공되지 않으면 전체 홈의 시나리오 조회.
- **반환**: (String) Markdown 형식의 시나리오 정보.

### automation_config

- **설명**: 예약 또는 지연 기기 제어 작업 설정.
- **매개변수**:
  - `scheduled_time` (String, 필수): 설정 시간점 (지연 작업의 경우 현재 시간 기준으로 변환), `YYYY-MM-DD HH:MM:SS` 형식 (예: `"2025-05-16 12:12:12"`).
  - `endpoint_ids` (Array<Integer>, 필수): 예약 제어할 기기 ID 목록.
  - `control_params` (Object, 필수): `device_control` 도구와 동일한 형식의 기기 제어 매개변수 (action, attribute, value 등 포함).
- **반환**: (String) 자동화 설정 결과 메시지.

### get_homes

- **설명**: 사용자 계정 하의 모든 홈 목록 조회.
- **매개변수**: 없음.
- **반환**: (String) 쉼표로 구분된 홈 이름 목록. 데이터가 없으면 빈 문자열 또는 적절한 메시지 반환.

### switch_home

- **설명**: 사용자의 현재 작업 홈 전환. 전환 후 후속 기기 조회, 제어 등 작업은 새로 전환된 홈을 대상으로 함.
- **매개변수**:
  - `home_name` (String, 필수): 대상 홈 이름 (`get_homes` 도구가 제공하는 사용 가능한 목록에서 가져와야 함).
- **반환**: (String) 전환 작업 결과 메시지.

## 프로젝트 구조

```
.
├── cmd.go                # Cobra CLI 명령 정의 및 프로그램 진입점 (main 함수 포함)
├── server.go             # MCP 서버 핵심 로직, 도구 정의 및 요청 처리
├── smh.go                # Aqara 스마트 홈 플랫폼 API 인터페이스 래퍼
├── middleware.go         # 미들웨어: 사용자 인증, 타임아웃 제어, 예외 복구
├── config.go             # 전역 설정 관리 및 환경 변수 처리
├── go.mod                # Go 모듈 종속성 관리 파일
├── go.sum                # Go 모듈 종속성 체크섬 파일
├── img/                  # README 문서에서 사용되는 이미지 리소스
├── LICENSE               # MIT 오픈소스 라이선스
└── README.md             # 프로젝트 문서
```

### 핵심 파일 설명

- **`cmd.go`**: Cobra 프레임워크 기반 CLI 구현, `run stdio`와 `run sse` 시작 모드 및 메인 진입 함수 정의
- **`server.go`**: MCP 서버 핵심 구현, 도구 등록, 요청 처리, 프로토콜 지원 담당
- **`smh.go`**: Aqara 스마트 홈 플랫폼 API 래퍼 계층, 기기 제어, 인증, 다중 홈 지원 제공
- **`middleware.go`**: 요청 처리 미들웨어, 인증 검증, 타임아웃 제어, 예외 처리 제공
- **`config.go`**: 전역 설정 관리, 환경 변수 처리 및 API 설정 담당

## 기여하기

Issue나 Pull Request를 제출하여 프로젝트에 기여해 주세요!

코드를 제출하기 전에 다음을 확인하세요:
1. 코드가 Go 언어 코딩 표준을 따릅니다.
2. 관련 MCP 도구와 프롬프트 인터페이스 정의가 일관성과 명확성을 유지합니다.
3. 변경 사항을 커버하는 단위 테스트를 추가하거나 업데이트합니다.
4. 필요한 경우 관련 문서(이 README 등)를 업데이트합니다.
5. 커밋 메시지가 명확하고 설명적인지 확인합니다.

## 라이선스

이 프로젝트는 [MIT License](/LICENSE) 하에 라이선스됩니다.
Copyright (c) 2025 Aqara-Copliot