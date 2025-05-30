<div align="center" style="display: flex; align-items: center; justify-content: center; ">

  <img src="/readme/img/logo.png" alt="Aqara Logo" height="120">
  <h1>MCP Server</h1>

</div>

<div align="center">

[English](/readme/README.md) | [中文](/readme/README_CN.md) | [繁體中文](/readme/README_CHT.md) | [Français](/readme/README_FR.md) | [한국어](/readme/README_KR.md) | Español | [日本語](/readme/README_JP.md) | [Deutsch](/readme/README_DE.md) | [Italiano](/readme/README_IT.md)

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/aqara/aqara-mcp-server)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![Release](https://img.shields.io/github/v/release/aqara/aqara-mcp-server)](https://github.com/aqara/aqara-mcp-server/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

Aqara MCP Server es un servicio de control de hogar inteligente desarrollado basado en el protocolo [MCP (Model Context Protocol)](https://modelcontextprotocol.io/introduction). Permite que cualquier asistente de IA o API que soporte el protocolo MCP (como Claude, Cursor, etc.) interactúe con sus dispositivos de hogar inteligente Aqara, logrando funciones como control de dispositivos a través de lenguaje natural, consulta de estados, ejecución de escenas, etc.

## Índice

- [Índice](#índice)
- [Características](#características)
- [Cómo Funciona](#cómo-funciona)
- [Inicio Rápido](#inicio-rápido)
  - [Requisitos Previos](#requisitos-previos)
  - [Instalación](#instalación)
    - [Método 1: Descargar Versión Precompilada (Recomendado)](#método-1-descargar-versión-precompilada-recomendado)
    - [Método 2: Compilar desde Código Fuente](#método-2-compilar-desde-código-fuente)
  - [Autenticación de Cuenta Aqara](#autenticación-de-cuenta-aqara)
  - [Configuración del Cliente](#configuración-del-cliente)
    - [Ejemplo de Configuración para Claude for Desktop](#ejemplo-de-configuración-para-claude-for-desktop)
    - [Descripción de Parámetros de Configuración](#descripción-de-parámetros-de-configuración)
    - [Otros Clientes MCP](#otros-clientes-mcp)
  - [Iniciar el Servicio](#iniciar-el-servicio)
    - [Modo Estándar (Recomendado)](#modo-estándar-recomendado)
    - [Modo HTTP (Opcional)](#modo-http-opcional)
- [Descripción de Herramientas de API](#descripción-de-herramientas-de-api)
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

- **Control Completo de Dispositivos**: Soporta control fino de múltiples atributos como encendido/apagado, brillo, temperatura de color, modo, etc. para dispositivos inteligentes Aqara
- **Consulta Flexible de Dispositivos**: Capacidad de consultar listas de dispositivos y sus estados detallados por habitación y tipo de dispositivo
- **Gestión Inteligente de Escenas**: Soporta consulta y ejecución de escenas de hogar inteligente preconfiguradas por el usuario
- **Historial de Dispositivos**: Consulta registros de cambios de estado históricos de dispositivos en rangos de tiempo específicos
- **Configuración de Automatización**: Soporta configuración de tareas de control de dispositivos programadas o retardadas
- **Soporte Multi-Hogar**: Soporta consulta y cambio entre diferentes hogares bajo la cuenta del usuario
- **Compatibilidad con Protocolo MCP**: Cumple completamente con las especificaciones del protocolo MCP, fácil de integrar con varios asistentes de IA
- **Mecanismo de Autenticación Segura**: Adopta autenticación segura basada en autorización de inicio de sesión + firma para proteger datos del usuario y seguridad del dispositivo
- **Ejecución Multiplataforma**: Desarrollado en Go, puede compilarse en archivos ejecutables para múltiples plataformas
- **Fácil de Extender**: Diseño modular que permite agregar nuevas herramientas y funcionalidades fácilmente

## Cómo Funciona

Aqara MCP Server actúa como un puente entre asistentes de IA y la plataforma de hogar inteligente Aqara:

1. **Asistente de IA (Cliente MCP)**: El usuario emite comandos a través del asistente de IA (por ejemplo, "Enciende las luces de la sala")
2. **Cliente MCP**: Analiza los comandos del usuario y llama las herramientas correspondientes proporcionadas por Aqara MCP Server según el protocolo MCP (por ejemplo, `device_control`)
3. **Aqara MCP Server (este proyecto)**: Recibe solicitudes del cliente, las valida y llama el módulo `smh.go`
4. **Módulo `smh.go`**: Utiliza las credenciales Aqara configuradas para comunicarse con la API en la nube de Aqara, ejecutando operaciones reales de dispositivos o consultas de datos
5. **Flujo de Respuesta**: La API en la nube de Aqara devuelve resultados, que se transmiten de vuelta al cliente MCP a través de Aqara MCP Server, finalmente presentándose al usuario

## Inicio Rápido

### Requisitos Previos

- Go (versión 1.24 o superior)
- Git (para compilar desde código fuente)
- Cuenta Aqara y dispositivos inteligentes vinculados

### Instalación

Puede elegir descargar archivos ejecutables precompilados o compilar desde código fuente.

#### Método 1: Descargar Versión Precompilada (Recomendado)

Visite la página de GitHub Releases para descargar el archivo ejecutable más reciente para su sistema operativo:

**📥 [Ir a la página de Releases para descargar](https://github.com/aqara/aqara-mcp-server/releases)**

Después de descargar el archivo comprimido correspondiente a su plataforma, descomprímalo y estará listo para usar.

#### Método 2: Compilar desde Código Fuente

```bash
# Clonar el repositorio
git clone https://github.com/aqara/aqara-mcp-server.git
cd aqara-mcp-server

# Descargar dependencias
go mod tidy

# Compilar archivo ejecutable
go build -o aqara-mcp-server
```

Después de completar la compilación, se generará el archivo ejecutable `aqara-mcp-server` en el directorio actual.

### Autenticación de Cuenta Aqara

Para que el MCP Server pueda acceder a su cuenta Aqara y controlar dispositivos, necesita completar primero la autorización de inicio de sesión.

Por favor visite la siguiente dirección para completar la autorización de inicio de sesión:
**🔗 [https://cdn.aqara.com/app/mcpserver/login.html](https://cdn.aqara.com/app/mcpserver/login.html)**

Después del inicio de sesión exitoso, obtendrá la información de autenticación necesaria (como `token`, `region`), que se utilizará en los pasos de configuración posteriores.

> ⚠️ **Recordatorio de Seguridad**: Por favor mantenga segura la información del `token` y no la divulgue a otros.

### Configuración del Cliente

Los métodos de configuración para diferentes clientes MCP varían ligeramente. El siguiente es un ejemplo de cómo configurar Claude for Desktop para usar este MCP Server:

#### Ejemplo de Configuración para Claude for Desktop

1. Abra la configuración (Settings) de Claude for Desktop

    ![Claude Open Setting](/readme/img/opening_setting.png)

2. Cambie a la pestaña de desarrollador (Developer), luego haga clic en editar configuración (Edit Config), use un editor de texto para abrir el archivo de configuración

    ![Claude Edit Configuration](/readme/img/edit_config.png)

3. Agregue la información de configuración de la "página de inicio de sesión exitoso" al archivo de configuración del cliente `claude_desktop_config.json`

    ![Configuration Example](/readme/img/config_info.png)

#### Descripción de Parámetros de Configuración

- `command`: Ruta completa al archivo ejecutable `aqara-mcp-server` que descargó o compiló
- `args`: Use `["run", "stdio"]` para iniciar el modo de transporte stdio
- `env`: Configuración de variables de entorno
  - `token`: Token de acceso obtenido de la página de inicio de sesión de Aqara
  - `region`: La región donde se encuentra su cuenta Aqara (como CN, US, EU, etc.)

#### Otros Clientes MCP

Para otros clientes que soportan el protocolo MCP (como ChatGPT, Cursor, etc.), la configuración es similar:

- Asegúrese de que el cliente soporte el protocolo MCP
- Configure la ruta del archivo ejecutable y parámetros de inicio
- Configure las variables de entorno `token` y `region`
- Seleccione el protocolo de transporte apropiado (se recomienda usar `stdio`)

### Iniciar el Servicio

#### Modo Estándar (Recomendado)

Reinicie Claude for Desktop. Luego podrá ejecutar operaciones como control de dispositivos, consulta de dispositivos, ejecución de escenas, etc. a través de lenguaje natural.

![Claude Chat Example](/readme/img/claude.png)

#### Modo HTTP (Opcional)

Si necesita usar el modo HTTP, puede iniciarlo así:

```bash
# Usar puerto predeterminado 8080
./aqara-mcp-server run http

# O especificar host y puerto personalizados
./aqara-mcp-server run http --host localhost --port 9000
```

Luego use los parámetros `["run", "http"]` en la configuración del cliente.

## Descripción de Herramientas de API

Los clientes MCP pueden interactuar con dispositivos de hogar inteligente Aqara llamando estas herramientas.

### Control de Dispositivos

#### device_control

Controla el estado o atributos de dispositivos de hogar inteligente (como encendido/apagado, temperatura, brillo, color, temperatura de color, etc.).

**Parámetros:**

- `endpoint_ids` _(Array\<Integer\>, requerido)_: Lista de IDs de dispositivos a controlar
- `control_params` _(Object, requerido)_: Objeto de parámetros de control, contiene operaciones específicas:
  - `action` _(String, requerido)_: Operación a ejecutar (como `"on"`, `"off"`, `"set"`, `"up"`, `"down"`, `"cooler"`, `"warmer"`)
  - `attribute` _(String, requerido)_: Atributo del dispositivo a controlar (como `"on_off"`, `"brightness"`, `"color_temperature"`, `"ac_mode"`)
  - `value` _(String | Number, opcional)_: Valor objetivo (requerido cuando action es "set")
  - `unit` _(String, opcional)_: Unidad del valor (como `"%"`, `"K"`, `"℃"`)

**Retorna:** Mensaje de resultado de la operación de control del dispositivo

### Consulta de Dispositivos

#### device_query

Obtiene lista de dispositivos según ubicación especificada (habitación) y tipo de dispositivo (no incluye información de estado en tiempo real).

**Parámetros:**

- `positions` _(Array\<String\>, opcional)_: Lista de nombres de habitaciones. Array vacío significa consultar todas las habitaciones
- `device_types` _(Array\<String\>, opcional)_: Lista de tipos de dispositivos (como `"Light"`, `"WindowCovering"`, `"AirConditioner"`, `"Button"`). Array vacío significa consultar todos los tipos

**Retorna:** Lista de dispositivos en formato Markdown, incluyendo nombres de dispositivos e IDs

#### device_status_query

Obtiene información del estado actual de dispositivos (usado para consultar información de estado en tiempo real como color, brillo, encendido/apagado, etc.).

**Parámetros:**

- `positions` _(Array\<String\>, opcional)_: Lista de nombres de habitaciones. Array vacío significa consultar todas las habitaciones
- `device_types` _(Array\<String\>, opcional)_: Lista de tipos de dispositivos. Valores opcionales iguales a `device_query`. Array vacío significa consultar todos los tipos

**Retorna:** Información de estado de dispositivos en formato Markdown

#### device_log_query

Consulta información de registro histórico de dispositivos.

**Parámetros:**

- `endpoint_ids` _(Array\<Integer\>, requerido)_: Lista de IDs de dispositivos para consultar historial
- `start_datetime` _(String, opcional)_: Hora de inicio de consulta, formato `YYYY-MM-DD HH:MM:SS` (ejemplo: `"2023-05-16 12:00:00"`)
- `end_datetime` _(String, opcional)_: Hora de fin de consulta, formato `YYYY-MM-DD HH:MM:SS`
- `attribute` _(String, opcional)_: Nombre específico del atributo del dispositivo a consultar (como `on_off`, `brightness`). Cuando no se proporciona, consulta todos los atributos registrados

**Retorna:** Información de estado histórico de dispositivos en formato Markdown

> 📝 **Nota:** La implementación actual puede mostrar "This feature will be available soon.", indicando que la funcionalidad está pendiente de completar.

### Gestión de Escenas

#### get_scenes

Consulta todas las escenas bajo el hogar del usuario, o escenas dentro de habitaciones específicas.

**Parámetros:**

- `positions` _(Array\<String\>, opcional)_: Lista de nombres de habitaciones. Array vacío significa consultar escenas de todo el hogar

**Retorna:** Información de escenas en formato Markdown

#### run_scenes

Ejecuta escenas específicas según IDs de escena.

**Parámetros:**

- `scenes` _(Array\<Integer\>, requerido)_: Lista de IDs de escenas a ejecutar

**Retorna:** Mensaje de resultado de ejecución de escenas

### Gestión del Hogar

#### get_homes

Obtiene lista de todos los hogares bajo la cuenta del usuario.

**Parámetros:** Ninguno

**Retorna:** Lista de nombres de hogares separados por comas. Si no hay datos, retorna cadena vacía o mensaje informativo correspondiente

#### switch_home

Cambia el hogar actualmente operado por el usuario. Después del cambio, operaciones posteriores como consulta de dispositivos, control, etc. se dirigirán al nuevo hogar cambiado.

**Parámetros:**

- `home_name` _(String, requerido)_: Nombre del hogar objetivo

**Retorna:** Mensaje de resultado de la operación de cambio

### Configuración de Automatización

#### automation_config

Configura tareas de control de dispositivos programadas o retardadas (actualmente solo soporta configuración de automatización de temporizador fijo).

**Parámetros:**

- `scheduled_time` _(String, requerido)_: Punto de tiempo establecido (si es una tarea retardada, se convierte basado en el punto de tiempo actual), formato `YYYY-MM-DD HH:MM:SS` (ejemplo: `"2025-05-16 12:12:12"`)
- `endpoint_ids` _(Array\<Integer\>, requerido)_: Lista de IDs de dispositivos para control programado
- `control_params` _(Object, requerido)_: Parámetros de control de dispositivos, usando el mismo formato que la herramienta `device_control` (incluyendo action, attribute, value, etc.)

**Retorna:** Mensaje de resultado de configuración de automatización

> 📝 **Nota:** La implementación actual puede mostrar "This feature will be available soon.", indicando que la funcionalidad está pendiente de completar.

## Estructura del Proyecto

### Estructura de Directorios

```text
.
├── cmd.go                # Definición de comandos CLI Cobra y punto de entrada del programa (incluye función main)
├── server.go             # Lógica central del servidor MCP, definición de herramientas y manejo de solicitudes
├── smh.go                # Encapsulado de interfaz API de la plataforma de hogar inteligente Aqara
├── middleware.go         # Middleware: autenticación de usuario, control de tiempo de espera, recuperación de excepciones
├── config.go             # Gestión de configuración global y manejo de variables de entorno
├── go.mod                # Archivo de gestión de dependencias del módulo Go
├── go.sum                # Archivo de suma de verificación de dependencias del módulo Go
├── readme/               # Documentación README y recursos de imágenes
│   ├── img/              # Directorio de recursos de imágenes
│   └── *.md              # Archivos README en múltiples idiomas
├── LICENSE               # Licencia de código abierto MIT
└── README.md             # Documento principal del proyecto
```

### Descripción de Archivos Principales

- **`cmd.go`**: Implementación CLI basada en el framework Cobra, define modos de inicio `run stdio` y `run http` y función de entrada principal
- **`server.go`**: Implementación central del servidor MCP, responsable del registro de herramientas, manejo de solicitudes y soporte de protocolo
- **`smh.go`**: Capa de encapsulado de API de la plataforma de hogar inteligente Aqara, proporciona control de dispositivos, autenticación y soporte multi-hogar
- **`middleware.go`**: Middleware de procesamiento de solicitudes, proporciona validación de autenticación, control de tiempo de espera y manejo de excepciones
- **`config.go`**: Gestión de configuración global, responsable del manejo de variables de entorno y configuración de API

## Guía de Desarrollo

¡Bienvenido a participar en las contribuciones del proyecto enviando Issues o Pull Requests!

Antes de enviar código, por favor asegúrese de que:

1. El código sigue las normas de codificación del lenguaje Go
2. Las herramientas MCP relacionadas y las definiciones de interfaz mantienen consistencia y claridad
3. Agregar o actualizar pruebas unitarias para cubrir sus cambios
4. Si es necesario, actualizar la documentación relacionada (como este README)
5. Asegurar que sus mensajes de commit sean claros y comprensibles

## Licencia

Este proyecto está licenciado bajo [MIT License](/LICENSE).

Copyright (c) 2025 Aqara-Copilot
