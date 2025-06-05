<div align="center" style="display: flex; align-items: center; justify-content: center; ">

  <img src="/readme/img/logo.png" alt="Logo de Aqara" height="120">
  <h1>Servidor MCP de Aqara</h1>

</div>

<div align="center">

[English](/readme/README.md) | [中文](/readme/README_CN.md) | [繁體中文](/readme/README_CHT.md) | [Français](/readme/README_FR.md) | [한국어](/readme/README_KR.md) | Español | [日本語](/readme/README_JP.md) | [Deutsch](/readme/README_DE.md) | [Italiano](/readme/README_IT.md)

[![Estado de Compilación](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/aqara/aqara-mcp-server)
[![Versión de Go](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![Lanzamiento](https://img.shields.io/github/v/release/aqara/aqara-mcp-server)](https://github.com/aqara/aqara-mcp-server/releases)
[![Licencia: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

El Servidor MCP de Aqara es un servicio de control de hogar inteligente desarrollado basado en el protocolo [MCP (Model Context Protocol)](https://modelcontextprotocol.io/introduction). Permite que cualquier asistente de IA o API compatible con el protocolo MCP (como Claude, Cursor, etc.) interactúe con sus dispositivos inteligentes Aqara, habilitando el control de dispositivos mediante lenguaje natural, consulta de estados, ejecución de escenas y más funcionalidades.

## Tabla de Contenidos

- [Tabla de Contenidos](#tabla-de-contenidos)
- [Características](#características)
- [Cómo Funciona](#cómo-funciona)
- [Inicio Rápido](#inicio-rápido)
  - [Requisitos Previos](#requisitos-previos)
  - [Instalación](#instalación)
    - [Método 1: Descargar Versión Precompilada (Recomendado)](#método-1-descargar-versión-precompilada-recomendado)
    - [Método 2: Compilar desde el Código Fuente](#método-2-compilar-desde-el-código-fuente)
  - [Autenticación de Cuenta Aqara](#autenticación-de-cuenta-aqara)
  - [Configuración del Cliente](#configuración-del-cliente)
    - [Ejemplo de Configuración de Claude for Desktop](#ejemplo-de-configuración-de-claude-for-desktop)
    - [Descripción de Parámetros de Configuración](#descripción-de-parámetros-de-configuración)
    - [Otros Clientes MCP](#otros-clientes-mcp)
  - [Iniciar el Servicio](#iniciar-el-servicio)
    - [Modo Estándar (Recomendado)](#modo-estándar-recomendado)
    - [Modo HTTP (`Próximamente`)](#modo-http-próximamente)
- [Descripción de Herramientas API](#descripción-de-herramientas-api)
  - [Control de Dispositivos](#control-de-dispositivos)
    - [device\_control](#device_control)
  - [Consulta de Dispositivos](#consulta-de-dispositivos)
    - [device\_query](#device_query)
    - [device\_status\_query](#device_status_query)
    - [device\_log\_query](#device_log_query)
  - [Gestión de Escenas](#gestión-de-escenas)
    - [get\_scenes](#get_scenes)
    - [run\_scenes](#run_scenes)
  - [Gestión del Hogar](#gestión-del-hogar)
    - [get\_homes](#get_homes)
    - [switch\_home](#switch_home)
  - [Configuración de Automatización](#configuración-de-automatización)
    - [automation\_config](#automation_config)
- [Estructura del Proyecto](#estructura-del-proyecto)
  - [Estructura de Directorios](#estructura-de-directorios)
  - [Descripción de Archivos Principales](#descripción-de-archivos-principales)
- [Guía de Desarrollo](#guía-de-desarrollo)
- [Licencia](#licencia)

## Características

- ✨ **Control Integral de Dispositivos**: Soporta control preciso de múltiples atributos de dispositivos inteligentes Aqara como encendido/apagado, brillo, temperatura de color, modos, etc.
- 🔍 **Consulta Flexible de Dispositivos**: Capacidad para consultar listas de dispositivos y sus estados detallados por habitación y tipo de dispositivo
- 🎬 **Gestión Inteligente de Escenas**: Soporta consulta y ejecución de escenas de hogar inteligente preconfiguradas por el usuario
- 📈 **Historial de Dispositivos**: Consulta registros históricos de cambios de estado de dispositivos en rangos de tiempo específicos
- ⏰ **Configuración de Automatización**: Soporta configuración de tareas de control de dispositivos programadas o con retraso
- 🏠 **Soporte Multi-hogar**: Soporta consulta y cambio entre diferentes hogares bajo la cuenta de usuario
- 🔌 **Compatibilidad con Protocolo MCP**: Cumple completamente con las especificaciones del protocolo MCP, fácil integración con varios asistentes de IA
- 🔐 **Mecanismo de Autenticación Segura**: Utiliza autenticación segura basada en autorización de login + firma para proteger datos de usuario y seguridad de dispositivos
- 🌐 **Ejecución Multiplataforma**: Desarrollado en Go, puede compilarse para ejecutables multiplataforma
- 🔧 **Fácil de Extender**: Diseño modular, permite agregar nuevas herramientas y funcionalidades convenientemente

## Cómo Funciona

El Servidor MCP de Aqara actúa como un puente entre asistentes de IA y la plataforma de hogar inteligente Aqara:

```mermaid
graph LR
    A[Asistente de IA] --> B[Cliente MCP]
    B --> C[Servidor MCP de Aqara]
    C --> D[API en la Nube de Aqara]
    D --> E[Dispositivos Inteligentes]
```

1. **Asistente de IA**: El usuario emite comandos a través del asistente de IA (por ejemplo, "enciende las luces del salón")
2. **Cliente MCP**: Analiza las instrucciones del usuario y llama a las herramientas correspondientes proporcionadas por el Servidor MCP de Aqara según el protocolo MCP (por ejemplo, `device_control`)
3. **Servidor MCP de Aqara (este proyecto)**: Recibe solicitudes del cliente, utiliza las credenciales Aqara configuradas para comunicarse con la API en la nube de Aqara, ejecutando operaciones reales de dispositivos o consultas de datos
4. **Flujo de Respuesta**: La API en la nube de Aqara devuelve resultados, que se transmiten de vuelta al cliente MCP a través del Servidor MCP de Aqara, presentándose finalmente al usuario

## Inicio Rápido

### Requisitos Previos

- **Go** (versión 1.24 o superior) - Solo necesario al compilar desde el código fuente
- **Git** (para compilar desde código fuente) - Opcional
- **Cuenta Aqara** con dispositivos inteligentes vinculados
- **Cliente compatible con protocolo MCP** (como Claude for Desktop, Cursor, etc.)

### Instalación

Puede elegir descargar el archivo ejecutable precompilado o compilar desde el código fuente.

#### Método 1: Descargar Versión Precompilada (Recomendado)

Visite la página de GitHub Releases para descargar el último archivo ejecutable para su sistema operativo:

**📥 [Ir a la página de Releases para descargar](https://github.com/aqara/aqara-mcp-server/releases)**

Después de descargar el archivo comprimido correspondiente a su plataforma, simplemente descomprímalo para usar.

#### Método 2: Compilar desde el Código Fuente

```bash
# Clonar repositorio
git clone https://github.com/aqara/aqara-mcp-server.git
cd aqara-mcp-server

# Descargar dependencias
go mod tidy

# Compilar archivo ejecutable
go build -o aqara-mcp-server
```

Después de completar la compilación, se generará el archivo ejecutable `aqara-mcp-server` en el directorio actual.

### Autenticación de Cuenta Aqara

Para que el Servidor MCP pueda acceder a su cuenta Aqara y controlar dispositivos, primero necesita completar la autorización de login.

Por favor visite la siguiente dirección para completar la autorización de login:
**🔗 [https://cdn.aqara.com/app/mcpserver/login.html](https://cdn.aqara.com/app/mcpserver/login.html)**

Después del login exitoso, obtendrá la información de autenticación necesaria (como `token`, `region`), que se utilizará en los pasos de configuración posteriores.

> ⚠️ **Recordatorio de Seguridad**: Por favor guarde cuidadosamente la información del `token`, no la revele a otros.

### Configuración del Cliente

Los métodos de configuración para diferentes clientes MCP varían ligeramente. A continuación se muestra un ejemplo de cómo configurar Claude for Desktop para usar este Servidor MCP:

#### Ejemplo de Configuración de Claude for Desktop

1. **Abrir la configuración (Settings) de Claude for Desktop**

    ![Abrir Configuración de Claude](/readme/img/opening_setting.png)

2. **Cambiar a la pestaña de desarrollador (Developer), luego hacer clic en editar configuración (Edit Config), usar el editor de texto para abrir el archivo de configuración**

    ![Editar Configuración de Claude](/readme/img/edit_config.png)

3. **Agregar la información de configuración de la "página de login exitoso" al archivo de configuración del cliente `claude_desktop_config.json`**

    ```json
    {
      "mcpServers": {
        "aqara": {
          "command": "/ruta/a/aqara-mcp-server",
          "args": ["run", "stdio"],
          "env": {
            "token": "su_token_aqui",
            "region": "su_region_aqui"
          }
        }
      }
    }
    ```

    ![Ejemplo de Configuración](/readme/img/config_info.png)

#### Descripción de Parámetros de Configuración

- `command`: Ruta completa al archivo ejecutable `aqara-mcp-server` que descargó o compiló
- `args`: Use `["run", "stdio"]` para iniciar el modo de transporte stdio
- `env`: Configuración de variables de entorno
  - `token`: Token de acceso obtenido de la página de login de Aqara
  - `region`: Región donde se encuentra su cuenta Aqara (regiones soportadas: CN, US, EU, KR, SG, RU)

#### Otros Clientes MCP

Para otros clientes compatibles con el protocolo MCP (como ChatGPT, Cursor, etc.), el método de configuración es similar:

- Asegurar que el cliente soporte el protocolo MCP
- Configurar la ruta del archivo ejecutable y parámetros de inicio
- Establecer variables de entorno `token` y `region`
- Elegir el protocolo de transporte apropiado (se recomienda usar `stdio`)

### Iniciar el Servicio

#### Modo Estándar (Recomendado)

Reinicie Claude for Desktop. Luego podrá ejecutar control de dispositivos, consulta de dispositivos, ejecución de escenas y otras operaciones mediante lenguaje natural.

Ejemplos de conversación:

- "Enciende las luces del salón"
- "Configura el aire acondicionado del dormitorio en modo frío, temperatura 24 grados"
- "Ver lista de dispositivos de todas las habitaciones"
- "Ejecutar escena de buenas noches"

![Ejemplo de Chat de Claude](/readme/img/claude.png)

#### Modo HTTP (`Próximamente`)

## Descripción de Herramientas API

Los clientes MCP pueden interactuar con dispositivos de hogar inteligente Aqara llamando a estas herramientas.

### Control de Dispositivos

#### device_control

Controla el estado o atributos de dispositivos de hogar inteligente (por ejemplo, encendido/apagado, temperatura, brillo, color, temperatura de color, etc.).

**Parámetros:**

- `endpoint_ids` _(Array\<Integer\>, requerido)_: Lista de IDs de dispositivos a controlar
- `control_params` _(Object, requerido)_: Objeto de parámetros de control, conteniendo operaciones específicas:
  - `action` _(String, requerido)_: Operación a ejecutar (como `"on"`, `"off"`, `"set"`, `"up"`, `"down"`, `"cooler"`, `"warmer"`)
  - `attribute` _(String, requerido)_: Atributo del dispositivo a controlar (como `"on_off"`, `"brightness"`, `"color_temperature"`, `"ac_mode"`)
  - `value` _(String | Number, opcional)_: Valor objetivo (requerido cuando action es "set")
  - `unit` _(String, opcional)_: Unidad del valor (como `"%"`, `"K"`, `"℃"`)

**Retorna:** Mensaje de resultado de la operación de control del dispositivo

### Consulta de Dispositivos

#### device_query

Obtiene lista de dispositivos basada en ubicación especificada (habitación) y tipo de dispositivo (no incluye información de estado en tiempo real).

**Parámetros:**

- `positions` _(Array\<String\>, opcional)_: Lista de nombres de habitaciones. Array vacío significa consultar todas las habitaciones
- `device_types` _(Array\<String\>, opcional)_: Lista de tipos de dispositivos (como `"Light"`, `"WindowCovering"`, `"AirConditioner"`, `"Button"`). Array vacío significa consultar todos los tipos

**Retorna:** Lista de dispositivos en formato Markdown, incluyendo nombres e IDs de dispositivos

#### device_status_query

Obtiene información de estado actual de dispositivos (usado para consultar información de estado en tiempo real como color, brillo, encendido/apagado, etc.).

**Parámetros:**

- `positions` _(Array\<String\>, opcional)_: Lista de nombres de habitaciones. Array vacío significa consultar todas las habitaciones
- `device_types` _(Array\<String\>, opcional)_: Lista de tipos de dispositivos. Valores opcionales iguales a `device_query`. Array vacío significa consultar todos los tipos

**Retorna:** Información de estado de dispositivos en formato Markdown

#### device_log_query

Consulta información de historial de dispositivos.

**Parámetros:**

- `endpoint_ids` _(Array\<Integer\>, requerido)_: Lista de IDs de dispositivos para consultar historial
- `start_datetime` _(String, opcional)_: Tiempo de inicio de consulta, formato `YYYY-MM-DD HH:MM:SS` (ejemplo: `"2023-05-16 12:00:00"`)
- `end_datetime` _(String, opcional)_: Tiempo de fin de consulta, formato `YYYY-MM-DD HH:MM:SS`
- `attributes` _(Array\<String\>, opcional)_: Lista de nombres de atributos de dispositivo a consultar (como `["on_off", "brightness"]`). Cuando no se proporciona, consulta todos los atributos registrados

**Retorna:** Información de estado histórico de dispositivos en formato Markdown

### Gestión de Escenas

#### get_scenes

Consulta todas las escenas bajo el hogar del usuario, o escenas dentro de habitaciones específicas.

**Parámetros:**

- `positions` _(Array\<String\>, opcional)_: Lista de nombres de habitaciones. Array vacío significa consultar escenas de todo el hogar

**Retorna:** Información de escenas en formato Markdown

#### run_scenes

Ejecuta escenas específicas basadas en IDs de escena.

**Parámetros:**

- `scenes` _(Array\<Integer\>, requerido)_: Lista de IDs de escenas a ejecutar

**Retorna:** Mensaje de resultado de ejecución de escenas

### Gestión del Hogar

#### get_homes

Obtiene lista de todos los hogares bajo la cuenta de usuario.

**Parámetros:** Ninguno

**Retorna:** Lista de nombres de hogares separados por comas. Si no hay datos, retorna cadena vacía o mensaje de información correspondiente

#### switch_home

Cambia el hogar actualmente operado por el usuario. Después del cambio, las operaciones posteriores de consulta de dispositivos, control, etc. se dirigirán al nuevo hogar cambiado.

**Parámetros:**

- `home_name` _(String, requerido)_: Nombre del hogar objetivo

**Retorna:** Mensaje de resultado de la operación de cambio

### Configuración de Automatización

#### automation_config

Configuración de automatización (actualmente solo soporta tareas de control de dispositivos programadas o con retraso).

**Parámetros:**

- `scheduled_time` _(String, requerido)_: Punto de tiempo para ejecución programada, usando formato Crontab estándar `"minuto hora día mes semana"`. Ejemplo: `"30 14 * * *"` (ejecutar a las 14:30 diariamente), `"0 9 * * 1"` (ejecutar a las 9:00 cada lunes)
- `endpoint_ids` _(Array\<Integer\>, requerido)_: Lista de IDs de dispositivos para control programado
- `control_params` _(Object, requerido)_: Parámetros de control de dispositivo, usando el mismo formato que la herramienta `device_control` (incluyendo action, attribute, value, etc.)
- `task_name` _(String, requerido)_: Nombre o descripción de esta tarea de automatización (usado para identificación y gestión)
- `execution_once` _(Boolean, opcional)_: Si ejecutar solo una vez
  - `true`: Ejecutar tarea solo una vez en el tiempo especificado (valor predeterminado)
  - `false`: Ejecutar tarea repetidamente de forma periódica (como diariamente, semanalmente, etc.)

**Retorna:** Mensaje de resultado de configuración de automatización

## Estructura del Proyecto

### Estructura de Directorios

```text
.
├── cmd.go                # Definición de comandos CLI Cobra y punto de entrada del programa (contiene función main)
├── server.go             # Lógica central del servidor MCP, definición de herramientas y manejo de solicitudes
├── smh.go                # Encapsulación de interfaz API de plataforma de hogar inteligente Aqara
├── middleware.go         # Middleware: autenticación de usuario, control de timeout, recuperación de excepciones
├── config.go             # Gestión de configuración global y manejo de variables de entorno
├── go.mod                # Archivo de gestión de dependencias del módulo Go
├── go.sum                # Archivo de suma de verificación de dependencias del módulo Go
├── readme/               # Documentos README y recursos de imágenes
│   ├── img/              # Directorio de recursos de imágenes
│   └── *.md              # Archivos README multiidioma
├── LICENSE               # Licencia de código abierto MIT
└── README.md             # Documento principal del proyecto
```

### Descripción de Archivos Principales

- **`cmd.go`**: Implementación CLI basada en framework Cobra, define modos de inicio `run stdio` y `run http` y función de entrada principal
- **`server.go`**: Implementación central del servidor MCP, responsable del registro de herramientas, manejo de solicitudes y soporte de protocolo
- **`smh.go`**: Capa de encapsulación de API de plataforma de hogar inteligente Aqara, proporciona control de dispositivos, autenticación y soporte multi-hogar
- **`middleware.go`**: Middleware de manejo de solicitudes, proporciona verificación de autenticación, control de timeout y manejo de excepciones
- **`config.go`**: Gestión de configuración global, responsable del manejo de variables de entorno y configuración de API

## Guía de Desarrollo

¡Bienvenido a participar en la contribución del proyecto enviando Issues o Pull Requests!

Antes de enviar código, por favor asegúrese de que:

1. El código sigue las normas de codificación del lenguaje Go
2. Las herramientas MCP relacionadas y definiciones de interfaz mantienen consistencia y claridad
3. Agregar o actualizar pruebas unitarias para cubrir sus cambios
4. Si es necesario, actualizar documentación relacionada (como este README)
5. Asegurar que sus mensajes de commit sean claros y comprensibles

**🌟 ¡Si este proyecto le es útil, por favor denos una Estrella!**

**🤝 ¡Bienvenido a unirse a nuestra comunidad, hagamos el hogar inteligente más inteligente juntos!**

## Licencia

Este proyecto está autorizado bajo [Licencia MIT](/LICENSE).

---

Copyright (c) 2025 Aqara-Copilot
