<div align="center" style="display: flex; align-items: center; justify-content: center; ">

  <img src="img/logo.png" alt="Aqara Logo" height="120">
  <h1>Aqara MCP Server</h1>

</div>

<div align="center">

[English](README.md) | [中文](README_CN.md) | [Français](README_FR.md) | 한국어 | [Español](README_ES.md) | [日本語](README_JP.md) | [Deutsch](README_DE.md) | [Italiano](README_IT.md)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![MCP Protocol](https://img.shields.io/badge/MCP-Protocol-00ff00)](https://modelcontextprotocol.io/)

</div>

**Aqara MCP Server**는 Aqara Agent가 제공하는 원격 MCP 서비스로, MCP를 지원하는 AI 애플리케이션이 Aqara 스마트홈 기능에 안전하게 연결할 수 있도록 합니다. MCP 연동이 필요하면 Aqara Agent가 제공하는 원격 MCP URL을 설정하면 됩니다.

> [!TIP]
> **권장: Aqara 공식 Agent Skills**
>
> 사용 중인 애플리케이션이 Agent Skills를 지원한다면(Codex, Cursor, OpenClaw 등), 공식 **Aqara Agent Skills**를 우선 사용하는 것을 권장합니다. MCP Server를 별도로 구성하지 않아도 자연어로 홈/공간, 기기, 씬, 자동화, 에너지 소비 등을 조회·제어할 수 있습니다.
>
> - GitHub: [aqara/aqara-agent-skills](https://github.com/aqara/aqara-agent-skills)
> - ClawHub: [aqara/aqara-agent](https://clawhub.ai/aqara/aqara-agent)

## 목차

- [개요](#개요)
- [기능](#기능)
- [동작 방식](#동작-방식)
- [빠른 시작](#빠른-시작)
  - [사전 요구사항](#사전-요구사항)
  - [1단계: 계정 인증](#1단계-계정-인증)
  - [2단계: 원격 MCP 구성](#2단계-원격-mcp-구성)
  - [3단계: 확인](#3단계-확인)
- [구성 시 유의 사항](#구성-시-유의-사항)
- [MCP Tool 참조](#mcp-tool-참조)
  - [핵심 Tool 개요](#핵심-tool-개요)
  - [홈 및 위치](#홈-및-위치)
  - [기기 조회 및 제어](#기기-조회-및-제어)
  - [씬](#씬)
  - [자동화](#자동화)
  - [에너지 소비](#에너지-소비)
  - [조명 시나리오 및 효과](#조명-시나리오-및-효과)
  - [펌웨어](#펌웨어)
  - [매개변수 규칙](#매개변수-규칙)
- [라이선스](#라이선스)

## 개요

현재 권장되는 MCP 연동 방식은 Aqara Agent를 중심으로 합니다.

- **Remote MCP**: Streamable HTTP / HTTP MCP를 지원하는 애플리케이션이 `https://agent.aqara.com/open/mcp`로 연결할 때 적합합니다.
- **Aqara Agent Skills**: Agent Skills를 지원하는 애플리케이션은 스킬을 설치해 MCP Server를 수동으로 구성하지 않아도 됩니다.
- **MCP Tool 기능**: 홈/공간, 기기, 씬, 자동화, 에너지 소비, 조명 시나리오 및 효과, 펌웨어 등 스마트홈 작업을 지원합니다.

## 기능

- 🔍 **유연한 기기 조회**: 홈/공간, 기기 유형 또는 기기 ID로 기기 기본 정보, 실시간 상태, 제어 로그를 조회할 수 있습니다.
- ✨ **포괄적인 기기 제어**: Aqara 기기의 전원, 밝기, 색온도, 온도, 풍속, 모드, 커튼 개폐 비율 등을 제어할 수 있습니다.
- 🎬 **스마트 씬 관리**: 씬 조회·실행 및 씬 실행 기록 조회를 지원합니다.
- ⏰ **자동화 조회**: 자동화 규칙 조회 및 자동화 실행 기록 확인을 지원합니다.
- 📈 **에너지 소비 통계**: 방/공간 또는 기기 단위로 전력·전기요금을 조회하며, 합계 및 상세 통계를 지원합니다.
- 💡 **조명 시나리오 및 효과 관리**: 조명 시나리오/효과 조회, 지정 효과 설정, 효과 구성 매개변수 조회를 지원합니다.
- 🔄 **펌웨어 관리**: 기기의 현재 펌웨어 버전·업그레이드 가능 버전 조회 및 펌웨어 업그레이드 시작을 지원합니다.
- 🏠 **다중 홈·다중 공간**: Aqara 계정의 홈 목록과 현재 홈의 방/공간을 조회할 수 있습니다.
- 🔌 **원격 MCP 연동**: HTTP MCP URL로 Cursor, Codex 등 애플리케이션에 연결할 수 있습니다.
- 🔐 **안전한 인증**: Aqara Agent 로그인으로 `aqara_api_key`를 발급받으며, 구성 시 자격 증명을 안전하게 보관하세요.

## 동작 방식

원격 MCP 모드에서는 애플리케이션이 HTTP로 Aqara Agent의 MCP 서비스에 연결하고, 로그인 페이지에서 생성한 Bearer 토큰을 요청에 포함합니다. Aqara Agent는 자격 증명 검증, Tool 호출 실행, 결과 반환을 담당합니다.

```mermaid
graph LR
    A[AI 앱 / MCP Host] --> B[Aqara Agent]
    B --> C[Aqara Cloud API]
    C --> D[Aqara 기기 / 씬 / 자동화]
```

1. **AI 앱 / MCP Host**: 사용자가 Cursor, Codex 등에서 자연어 명령을 입력합니다.
2. **Aqara Agent**: 사용자 자격 증명을 검증하고 해당 Tool을 해석·실행합니다.
3. **Aqara Cloud API**: 기기, 씬, 자동화, 에너지 소비, 조명 효과, 펌웨어 등의 데이터 조회 또는 제어를 수행합니다.

---

## 빠른 시작

### 사전 요구사항

- **Aqara 계정** 및 등록된 스마트 기기.
- **원격 MCP를 지원하는 애플리케이션**(예: Cursor, Codex).
- **Aqara Agent 자격 증명**: 로그인 페이지에서 `aqara_api_key`와 `aqara_mcp_url`을 발급받습니다.

### 1단계: 계정 인증

1. **로그인 페이지 접속**:
   [https://agent.aqara.com/login](https://agent.aqara.com/login)

2. **로그인 완료**:
   - Aqara 계정으로 로그인합니다.
   - 로그인 후 페이지에 표시된 `aqara_api_key`를 복사합니다.
   - MCP 구성 시 페이지에 표시된 `aqara_mcp_url`을 사용합니다. 일반적으로 `https://agent.aqara.com/open/mcp`입니다.

3. **자격 증명 안전 보관**:

   > `aqara_api_key`를 안전하게 보관하세요. 저장소에 커밋하거나, 스크린샷으로 공개하거나, 타인과 공유하지 마세요.

   ![Aqara Agent 로그인 후 표시되는 구성 정보](img/config_info.png)

### 2단계: 원격 MCP 구성

#### Cursor

1. Cursor 설정에서 `Tools & MCPs`로 이동한 뒤 `New MCP Server`를 클릭합니다.

   ![Cursor MCP 설정 진입](img/cursor_opening_setting.png)

2. 원격 MCP 구성을 추가합니다. URL에는 로그인 페이지의 `aqara_mcp_url`을 사용하고, 수동 입력 시 `/open/mcp` 경로를 사용하세요.

   ```json
   {
     "mcpServers": {
       "aqara": {
         "type": "http",
         "url": "https://agent.aqara.com/open/mcp",
         "headers": {
           "Authorization": "Bearer <YOUR_AQARA_API_KEY>"
         }
       }
     }
   }
   ```

3. 구성을 저장하고 Cursor를 재시작해 MCP 구성을 적용합니다.

#### Codex

1. Codex 설정에서 사용자 지정 MCP Server를 추가합니다.
2. 유형으로 `Streamable HTTP`를 선택합니다.
3. URL에 로그인 페이지의 `aqara_mcp_url`(예: `https://agent.aqara.com/open/mcp`)을 입력합니다.
4. Bearer 토큰에 `aqara_api_key` 값을 입력합니다.

![Codex 사용자 지정 MCP 설정](img/codex_opening_setting.png)

### 3단계: 확인

구성이 완료되면 다음과 같은 자연어 요청으로 테스트할 수 있습니다.

```text
사용자: 우리 집의 모든 기기를 보여줘
어시스턴트: MCP로 기기 목록 조회

사용자: 거실 조명을 켜줘
어시스턴트: MCP로 기기 제어 실행

사용자: 영화 감상 씬 실행해줘
어시스턴트: MCP로 씬 실행
```

애플리케이션의 MCP 패널에 Aqara가 연결된 것으로 표시되고 Aqara 관련 Tool이 보이면 구성이 적용된 것입니다.

---

## 구성 시 유의 사항

- MCP URL은 `https://agent.aqara.com/open/mcp` 또는 로그인 페이지의 `aqara_mcp_url`을 사용하세요. 로그인 페이지 주소를 MCP URL로 사용하지 마세요.
- 기기 제어, 씬 실행, 펌웨어 업그레이드 등 Tool은 실제 홈 기기에 영향을 줍니다. 처음 사용할 때는 조회용 Tool로 홈, 공간, 기기, 씬 정보를 먼저 확인하는 것을 권장합니다.
- 연결에 실패하면 MCP 유형이 HTTP / Streamable HTTP인지, URL에 `/open/mcp`가 포함되는지, 자격 증명 만료 여부, 구성 변경 후 애플리케이션 재시작 또는 MCP 다시 불러오기 여부를 확인하세요.

---

## MCP Tool 참조

아래 Tool 목록은 현재 Aqara Agent 서비스에 등록된 함수 정의를 바탕으로 정리했습니다. 애플리케이션마다 Tool 이름 표시 방식은 다를 수 있으나, 매개변수 의미와 기능 범위는 동일합니다.

### 핵심 Tool 개요

| Tool 범주 | Tool | 설명 |
| --- | --- | --- |
| **홈 및 위치** | `all_homes_inquiry`, `position_base_inquiry` | 홈, 방/공간 정보 조회 |
| **기기 조회 및 제어** | `device_base_inquiry`, `device_status_inquiry`, `device_status_control`, `fuzzy_device_batch_control`, `device_log_inquiry` | 기기 기본 정보·실시간 상태 조회, 제어, 제어 로그 조회 |
| **씬** | `scene_base_inquiry`, `scene_run`, `scene_execution_history_inquiry` | 씬 조회·실행 및 실행 기록 조회 |
| **자동화** | `automation_base_inquiry`, `automation_execution_history_inquiry` | 자동화 규칙 및 실행 기록 조회 |
| **에너지 소비** | `energy_consumption_inquiry_for_position`, `energy_consumption_inquiry_for_device` | 방/공간 또는 기기 단위 전력·전기요금 조회 |
| **조명 시나리오 및 효과** | `lighting_effect_inquiry`, `device_lighting_effect_inquiry`, `lighting_effect_control`, `lighting_effect_config_params_inquiry` | 조명 시나리오/효과 조회·설정 및 구성 매개변수 조회 |
| **펌웨어** | `device_firmware_inquiry`, `device_firmware_upgrade` | 기기 펌웨어 조회 및 업그레이드 |

### 홈 및 위치

#### `all_homes_inquiry`

현재 Aqara 계정의 모든 홈 목록을 조회합니다.

**매개변수:** 없음

**반환:** 홈 이름, 홈 ID 등이 포함된 홈 목록.

#### `position_base_inquiry`

현재 홈의 모든 방/공간 기본 정보를 조회합니다.

**매개변수:** 없음

**반환:** 위치 이름, 위치 ID 등이 포함된 방/공간 목록.

### 기기 조회 및 제어

#### `device_base_inquiry`

방/공간 및 기기 유형으로 기기 기본 정보를 조회합니다. 실시간 상태는 포함하지 않습니다.

**매개변수:**

- `position_ids` _(Array\<String\>, 선택)_: 방/공간 ID 목록. 비어 있으면 위치로 필터링하지 않습니다.
- `device_types` _(Array\<String\>, 선택)_: 기기 유형 목록(예: `Light`, `Switch`, `Outlet`, `AirConditioner`, `WindowCovering`, `Camera` 등). 비어 있으면 기기 유형으로 필터링하지 않습니다.

**반환:** 기기 이름, 기기 ID, 소속 위치, 기기 유형 등이 포함된 기기 기본 정보 목록.

#### `device_status_inquiry`

전원, 밝기, 색온도, 온도, 풍속, 모드 등 기기 실시간 상태를 조회합니다.

**매개변수:**

- `device_ids` _(Array\<String\>, 선택)_: 기기 ID 목록. 지정 시 기기 ID 우선 조회.
- `position_ids` _(Array\<String\>, 선택)_: 방/공간 ID 목록.
- `device_types` _(Array\<String\>, 선택)_: 기기 유형 목록.

**반환:** 기기 현재 읽기 가능 상태가 포함된 기기 상태 목록.

#### `device_status_control`

지정 기기의 상태 또는 속성(전원, 밝기, 색온도, 온도, 풍속, 모드, 커튼 비율 등)을 제어합니다.

**매개변수:**

- `device_ids` _(Array\<String\>, 필수)_: 대상 기기 ID 목록.
- `attribute` _(String, 필수)_: 제어할 속성(예: `on_off`, `brightness`, `color_temperature`, `temperature`, `percentage`, `mode` 등).
- `action` _(String, 필수)_: 제어 동작(예: `on`, `off`, `set`, `up`, `down`, `warmer`, `cooler`, `start`, `stop` 등).
- `value` _(String, 선택)_: 목표 값(예: `50`, `max`, `min`, `cool`, `heat`, `red` 등).

**반환:** 기기 제어 실행 결과.

#### `fuzzy_device_batch_control`

방/공간 및 기기 유형으로 기기를 일괄 제어합니다. "집안 모든 조명 끄기", "거실 전체 끄기", "모든 에어컨을 26도로" 등에 적합합니다.

**매개변수:**

- `position_ids` _(Array\<String\>, 선택)_: 방/공간 ID 목록. 비어 있으면 전체 홈 범위를 의미할 수 있습니다.
- `device_types` _(Array\<String\>, 선택)_: 기기 유형 목록.
- `attribute` _(String, 필수)_: 제어할 속성.
- `action` _(String, 필수)_: 제어 동작.
- `value` _(String, 선택)_: 목표 값.

**반환:** 일괄 제어 실행 결과.

#### `device_log_inquiry`

지정 시간 범위의 기기 제어 로그(제어 시간, 내용, 결과 등)를 조회합니다.

**매개변수:**

- `time_range` _(Array\<String\>, 선택)_: 시간 구간. 형식 예: `["2026-01-01 00:00:00", "2026-01-01 23:59:59"]`.
- `device_ids` _(Array\<String\>, 선택)_: 기기 ID 목록. 지정 시 기기 ID 우선 조회.
- `position_ids` _(Array\<String\>, 선택)_: 방/공간 ID 목록.
- `device_types` _(Array\<String\>, 선택)_: 기기 유형 목록.

**반환:** 기기 제어 로그 목록 및 실제 조회 시간 범위.

### 씬

#### `scene_base_inquiry`

씬 기본 정보를 조회합니다. 씬 ID, 위치 ID, 기기 ID로 필터링할 수 있습니다.

**매개변수:**

- `scene_ids` _(Array\<String\>, 선택)_: 씬 ID 목록. 지정 시 씬 ID 우선 조회.
- `position_ids` _(Array\<String\>, 선택)_: 방/공간 ID 목록.
- `device_ids` _(Array\<String\>, 선택)_: 기기 ID 목록. 해당 기기와 연관된 씬 조회에 사용.

**반환:** 씬 기본 정보 목록.

#### `scene_run`

지정한 하나 이상의 씬을 실행합니다.

**매개변수:**

- `scene_ids` _(Array\<String\>, 필수)_: 실행할 씬 ID 목록.

**반환:** 씬 실행 결과.

#### `scene_execution_history_inquiry`

지정 시간 범위의 씬 실행 기록을 조회합니다.

**매개변수:**

- `time_range` _(Array\<String\>, 선택)_: 시간 구간.
- `scene_ids` _(Array\<String\>, 선택)_: 씬 ID 목록.
- `position_ids` _(Array\<String\>, 선택)_: 방/공간 ID 목록.
- `device_ids` _(Array\<String\>, 선택)_: 기기 ID 목록.

**반환:** 씬 실행 기록 목록 및 실제 조회 시간 범위.

### 자동화

#### `automation_base_inquiry`

자동화 규칙 기본 정보를 조회합니다. 자동화 ID, 위치 ID, 기기 ID로 필터링할 수 있습니다.

**매개변수:**

- `automation_ids` _(Array\<String\>, 선택)_: 자동화 ID 목록. 지정 시 자동화 ID 우선 조회.
- `position_ids` _(Array\<String\>, 선택)_: 방/공간 ID 목록.
- `device_ids` _(Array\<String\>, 선택)_: 기기 ID 목록. 해당 기기와 연관된 자동화 조회에 사용.

**반환:** 자동화 규칙 정보 목록.

#### `automation_execution_history_inquiry`

지정 시간 범위의 자동화 규칙 실행 기록을 조회합니다.

**매개변수:**

- `time_range` _(Array\<String\>, 선택)_: 시간 구간.
- `automation_ids` _(Array\<String\>, 선택)_: 자동화 ID 목록.
- `position_ids` _(Array\<String\>, 선택)_: 방/공간 ID 목록.
- `device_ids` _(Array\<String\>, 선택)_: 기기 ID 목록.

**반환:** 자동화 실행 기록 목록 및 실제 조회 시간 범위.

### 에너지 소비

#### `energy_consumption_inquiry_for_position`

홈/방/공간 단위로 전력 또는 전기요금을 조회합니다. 합계 및 상세를 지원합니다.

**매개변수:**

- `data_type` _(String, 필수)_: 조회 유형. `1`은 전력, `2`는 전기요금, `3`은 전력과 전기요금.
- `time_range` _(Array\<String\>, 필수)_: 시간 구간.
- `time_gradient` _(String, 선택)_: 통계 단위. `30min`, `1hour`, `1day`, `1week`, `1month` 가능.
- `data_aggregation_mode` _(String, 선택)_: 집계 모드. `total`은 합계, `detail`은 상세.
- `positions` _(Array\<String\>, 선택)_: 방/공간 ID 목록. 비어 있으면 유효한 모든 방을 조회.

**반환:** 방/공간 단위 전력/전기요금 통계 결과.

#### `energy_consumption_inquiry_for_device`

기기 단위로 전력 또는 전기요금을 조회합니다. 위치 또는 기기로 필터링할 수 있으며, 합계 및 상세를 지원합니다.

**매개변수:**

- `data_type` _(String, 필수)_: 조회 유형. `1`은 전력, `2`는 전기요금, `3`은 전력과 전기요금.
- `time_range` _(Array\<String\>, 필수)_: 시간 구간.
- `time_gradient` _(String, 선택)_: 통계 단위. `30min`, `1hour`, `1day`, `1week`, `1month` 가능.
- `data_aggregation_mode` _(String, 선택)_: 집계 모드. `total`은 합계, `detail`은 상세.
- `positions` _(Array\<String\>, 선택)_: 방/공간 ID 목록.
- `device_ids` _(Array\<String\>, 선택)_: 기기 ID 목록. 지정 시 기기 우선 조회.

**반환:** 기기 단위 전력/전기요금 통계 결과.

### 조명 시나리오 및 효과

#### `lighting_effect_inquiry`

홈에서 사용 가능한 조명 시나리오/효과 정보를 조회합니다.

**매개변수:** 없음

**반환:** 제어에 사용할 수 있는 효과 이름 및 적용 범위가 포함된 효과 목록.

#### `device_lighting_effect_inquiry`

기기별로 지원하는 조명 효과 이름을 조회합니다.

**매개변수:**

- `device_ids` _(Array\<String\>, 필수)_: 효과를 조회할 기기 ID 목록.

**반환:** 기기와 효과 이름의 대응 목록.

#### `lighting_effect_control`

지정 기기 또는 방/공간의 조명을 지정 효과로 전환합니다.

**매개변수:**

- `effect_name` _(String, 필수)_: 효과 이름.
- `device_ids` _(Array\<String\>, 선택)_: 대상 기기 ID 목록. 지정 시 기기 우선 제어.
- `position_ids` _(Array\<String\>, 선택)_: 방/공간 ID 목록.

**반환:** 조명 효과 제어 실행 결과.

#### `lighting_effect_config_params_inquiry`

조명 기기에서 효과를 구성할 때 필요한 매개변수 정보를 조회합니다.

**매개변수:**

- `device_ids` _(Array\<String\>, 필수)_: 대상 조명 기기 ID 목록.

**반환:** 구성 가능 항목, 값 범위, 저장된 사용자 효과 등이 포함된 효과 구성 매개변수 목록.

### 펌웨어

#### `device_firmware_inquiry`

기기의 현재 펌웨어 버전과 업그레이드 가능 버전을 일괄 조회합니다.

**매개변수:**

- `device_ids` _(Array\<String\>, 선택)_: 기기 ID 목록. 지정 시 기기 우선 조회.
- `position_ids` _(Array\<String\>, 선택)_: 방/공간 ID 목록.
- `device_types` _(Array\<String\>, 선택)_: 기기 유형 목록.

**반환:** 기기 이름, 온라인 상태, 현재·업그레이드 가능 펌웨어 버전이 포함된 펌웨어 정보 목록.

#### `device_firmware_upgrade`

기기, 위치 또는 유형으로 필터링한 뒤 업그레이드 가능한 기기의 펌웨어 업그레이드를 시작합니다.

**매개변수:**

- `device_ids` _(Array\<String\>, 선택)_: 기기 ID 목록. 지정 시 해당 기기 우선 업그레이드.
- `position_ids` _(Array\<String\>, 선택)_: 방/공간 ID 목록.
- `device_types` _(Array\<String\>, 선택)_: 기기 유형 목록.

**반환:** 펌웨어 업그레이드 제출 결과.

### 매개변수 규칙

- `position_ids` / `positions`: 방/공간 ID 목록. 미지정 시 조회·제어 범위는 해당 Tool 설명을 따릅니다.
- `device_ids`: 기기 ID 또는 기기 엔드포인트 ID 목록. 상위 식별 및 서버 매핑으로 처리됩니다.
- `device_types`: 기기 유형 목록(예: `Light`, `Switch`, `Outlet`, `AirConditioner`, `WindowCovering`, `Camera`, `TemperatureSensor` 등).
- `attribute`: 제어 속성(예: `on_off`, `brightness`, `color_temperature`, `temperature`, `wind_speed`, `mode`, `percentage`, `volume`, `color` 등).
- `action`: 제어 동작(예: `on`, `off`, `set`, `up`, `down`, `warmer`, `cooler`, `start`, `stop`, `pause`, `resume` 등).
- `value`: 목표 값(예: `50`, `100`, `max`, `min`, `red`, `cool`, `heat`, 조명 효과 이름 등).
- `time_range`: 시간 구간 배열. 일반 형식: `["YYYY-MM-DD HH:MM:SS", "YYYY-MM-DD HH:MM:SS"]`.
- `data_type`: 에너지 조회 유형. `1`은 전력, `2`는 전기요금, `3`은 전력과 전기요금.
- `time_gradient`: 에너지 통계 단위. `30min`, `1hour`, `1day`, `1week`, `1month` 가능.
- `data_aggregation_mode`: 에너지 집계 모드. `total`은 합계, `detail`은 상세.

## 라이선스

본 프로젝트는 [MIT 라이선스](LICENSE)에 따라 배포됩니다. 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.

---

Copyright © 2025 Aqara-Agent. 모든 권리 보유.
