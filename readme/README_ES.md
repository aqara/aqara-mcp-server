# Servidor Aqara MCP

[English](/readme/README.md) | [中文](/readme/README_CN.md) | [繁體中文](/readme/README_CHT.md) | [Français](/readme/README_FR.md) | [한국어](/readme/README_KR.md) | Español | [日本語](/readme/README_JP.md) | [Deutsch](/readme/README_DE.md) | [Italiano](/readme/README_IT.md)

[![Estado de Build](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/aqara/aqara-mcp-server)
[![Versión Go](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![Release](https://img.shields.io/github/v/release/aqara/aqara-mcp-server)](https://github.com/aqara/aqara-mcp-server/releases)
[![Licencia: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Aqara MCP Server es un servicio de control domótico desarrollado basado en el protocolo [MCP (Model Context Protocol)](https://modelcontextprotocol.io/introduction). Permite que cualquier asistente de IA o API que soporte el protocolo MCP (como Claude, ChatGPT, Cursor, etc.) interactúe con sus dispositivos domóticos Aqara, habilitando control de dispositivos, consultas de estado, ejecución de escenas y más a través de lenguaje natural.

## Índice

- [Servidor Aqara MCP](#servidor-aqara-mcp)
  - [Índice](#índice)
  - [Características](#características)
  - [Cómo Funciona](#cómo-funciona)
  - [Inicio Rápido](#inicio-rápido)
    - [Prerrequisitos](#prerrequisitos)
    - [Instalación](#instalación)
    - [Autenticación de Cuenta Aqara](#autenticación-de-cuenta-aqara)
    - [Ejemplo de Configuración (Claude for Desktop)](#ejemplo-de-configuración-claude-for-desktop)
    - [Ejecutar el Servicio](#ejecutar-el-servicio)
  - [Herramientas Disponibles](#herramientas-disponibles)
    - [device\_control](#device_control)
    - [device\_query](#device_query)
    - [device\_status\_query](#device_status_query)
    - [device\_log\_query](#device_log_query)
    - [run\_scenes](#run_scenes)
    - [get\_scenes](#get_scenes)
    - [automation\_config](#automation_config)
    - [get\_homes](#get_homes)
    - [switch\_home](#switch_home)
  - [Estructura del Proyecto](#estructura-del-proyecto)
    - [Descripciones de Archivos Principales](#descripciones-de-archivos-principales)
  - [Contribuir](#contribuir)
  - [Licencia](#licencia)

## Características

- **Control Integral de Dispositivos**: Soporte para control fino de dispositivos inteligentes Aqara incluyendo interruptores, brillo, temperatura de color, modos, y más.
- **Consultas Flexibles de Dispositivos**: Consultar listas de dispositivos y estados detallados por habitación y tipo de dispositivo.
- **Gestión Inteligente de Escenas**: Soporte para consultar y ejecutar escenas domóticas preestablecidas por el usuario.
- **Historial de Dispositivos**: Consultar cambios de estado históricos de dispositivos dentro de rangos de tiempo especificados.
- **Configuración de Automatización**: Soporte para configurar tareas de control de dispositivos programadas o diferidas.
- **Soporte Multi-Hogar**: Soporte para consultar y cambiar entre diferentes hogares bajo cuentas de usuario.
- **Compatible con Protocolo MCP**: Totalmente conforme a las especificaciones del protocolo MCP, fácil de integrar con varios asistentes de IA.
- **Mecanismo de Autenticación Seguro**: Utiliza autenticación segura basada en autorización de inicio de sesión + firma para proteger datos de usuario y seguridad de dispositivos.
- **Multiplataforma**: Construido con Go, puede compilarse en ejecutables para múltiples plataformas.
- **Fácil de Extender**: Diseño modular permite agregar fácilmente nuevas herramientas y características.

## Cómo Funciona

Aqara MCP Server actúa como un puente entre asistentes de IA y la plataforma domótica Aqara:

1. **Asistente de IA (Cliente MCP)**: Los usuarios emiten comandos a través de asistentes de IA (ej., "Enciende las luces del salón").
2. **Cliente MCP**: Analiza comandos de usuario y llama herramientas correspondientes proporcionadas por Aqara MCP Server según el protocolo MCP (ej., `device_control`).
3. **Aqara MCP Server (Este Proyecto)**: Recibe solicitudes de clientes, las valida, y llama al módulo `smh.go`.
4. **Módulo `smh.go`**: Usa credenciales Aqara configuradas para comunicarse con APIs en la nube de Aqara para operaciones reales de dispositivos o consultas de datos.
5. **Flujo de Respuesta**: La API en la nube de Aqara devuelve resultados, que se pasan a través de Aqara MCP Server al cliente MCP y finalmente se presentan al usuario.

## Inicio Rápido

### Prerrequisitos

- Go (versión 1.24 o superior)
- Git (para construir desde fuentes)
- Cuenta Aqara con dispositivos inteligentes conectados

### Instalación

Puede elegir descargar ejecutables precompilados o construir desde fuentes.

**Opción 1: Descargar Versión Precompilada (Recomendado)**

Visite el enlace a continuación para descargar el último paquete ejecutable para su sistema operativo.

[Página de Releases](https://github.com/aqara/aqara-mcp-server/releases)

Extraer y usar directamente.

**Opción 2: Construir desde Fuentes**

```bash
# Clonar repositorio
git clone https://github.com/aqara/aqara-mcp-server.git
cd aqara-mcp-server

# Descargar dependencias
go mod tidy

# Construir ejecutable
go build -o aqara-mcp-server
```

Después de construir, el ejecutable `aqara-mcp-server` se generará en el directorio actual.

### Autenticación de Cuenta Aqara

Para permitir que el MCP Server acceda a su cuenta Aqara y controle dispositivos, primero debe completar la autorización de inicio de sesión.

Por favor visite la siguiente dirección para completar la autorización de inicio de sesión:
[https://cdn.aqara.com/app/mcpserver/login.html](https://cdn.aqara.com/app/mcpserver/login.html)

Después de un inicio de sesión exitoso, obtendrá la información de autenticación necesaria (como `token`, `region`), que se usará en pasos de configuración posteriores.

**Por favor mantenga esta información segura, especialmente el `token` - no lo comparta con otros.**

### Ejemplo de Configuración (Claude for Desktop)

Diferentes clientes MCP tienen métodos de configuración ligeramente diferentes. Aquí hay un ejemplo de cómo configurar Claude for Desktop para usar este MCP Server:

1. Abrir Configuración de Claude for Desktop.
2. Cambiar a la pestaña Desarrollador.
3. Hacer clic en Editar Config para abrir el archivo de configuración con un editor de texto.

   ![](/readme/img/setting0.png)
   ![](/readme/img/setting1.png)

4. Agregar la información de configuración de la "Página de Éxito de Inicio de Sesión" al archivo de configuración del cliente (claude_desktop_config.json). Ejemplo de configuración:

   ![](/readme/img/config.png)

**Notas de Configuración:**
- `command`: Ruta completa a su ejecutable `aqara-mcp-server` descargado o construido
- `args`: Usar `["run", "stdio"]` para iniciar modo de transporte stdio
- `env`: Configuración de variables de entorno
  - `token`: Token de acceso obtenido de la página de inicio de sesión de Aqara
  - `region`: Región de su cuenta Aqara (ej., CN, US, EU, etc.)

### Ejecutar el Servicio

Reiniciar Claude for Desktop. Luego puede usar conversaciones para llamar herramientas proporcionadas por el MCP Server para control de dispositivos, consultas de dispositivos, y otras operaciones.

![](/readme/img/claude.png)

**Configuración de Otros Clientes MCP**

Para otros clientes que soportan el protocolo MCP (como Claude, ChatGPT, Cursor, etc.), la configuración es similar:
- Asegurar que el cliente soporte el protocolo MCP
- Configurar ruta de archivo ejecutable y parámetros de inicio
- Establecer variables de entorno `token` y `region`
- Elegir protocolo de transporte apropiado (stdio recomendado)

**Modo SSE (Opcional)**

Si necesita usar modo SSE (Server-Sent Events), puede iniciarlo así:

```bash
# Usar puerto por defecto 8080
./aqara-mcp-server run sse

# O especificar host y puerto personalizados
./aqara-mcp-server run sse --host localhost --port 9000
```

Luego usar parámetros `["run", "sse"]` en la configuración del cliente.

## Herramientas Disponibles

Los clientes MCP pueden interactuar con dispositivos domóticos Aqara llamando estas herramientas.

### device_control

- **Descripción**: Controlar el estado o propiedades de dispositivos domóticos (ej., encendido/apagado, temperatura, brillo, color, temperatura de color, etc.).
- **Parámetros**:
  - `endpoint_ids` (Array<Integer>, requerido): Lista de IDs de dispositivos a controlar.
  - `control_params` (Object, requerido): Objeto de parámetros de control conteniendo operaciones específicas.
    - `action` (String, requerido): Acción a ejecutar. Ejemplos: `"on"`, `"off"`, `"set"`, `"up"`, `"down"`, `"cooler"`, `"warmer"`.
    - `attribute` (String, requerido): Atributo de dispositivo a controlar. Ejemplos: `"on_off"`, `"brightness"`, `"color_temperature"`, `"ac_mode"`.
    - `value` (String | Number, opcional): Valor objetivo (requerido cuando la acción es "set").
    - `unit` (String, opcional): Unidad del valor (ej., `"%"`, `"K"`, `"℃"`).
- **Retorna**: (String) Mensaje de resultado de operación para control de dispositivo.

### device_query

- **Descripción**: Obtener lista de dispositivos por ubicación especificada (habitación) y tipo de dispositivo (no incluye información de estado en tiempo real, solo lista dispositivos y sus IDs).
- **Parámetros**:
  - `positions` (Array<String>, opcional): Lista de nombres de habitaciones. Si es array vacío o no se proporciona, consulta todas las habitaciones.
  - `device_types` (Array<String>, opcional): Lista de tipos de dispositivos. Ejemplos: `"Light"`, `"WindowCovering"`, `"AirConditioner"`, `"Button"`, etc. Si es array vacío o no se proporciona, consulta todos los tipos.
- **Retorna**: (String) Lista de dispositivos formateada en Markdown incluyendo nombres de dispositivos e IDs.

### device_status_query

- **Descripción**: Obtener información de estado actual de dispositivos (para consultar atributos relacionados con estado como color, brillo, interruptores, etc.).
- **Parámetros**:
  - `positions` (Array<String>, opcional): Lista de nombres de habitaciones. Si es array vacío o no se proporciona, consulta todas las habitaciones.
  - `device_types` (Array<String>, opcional): Lista de tipos de dispositivos. Mismas opciones que `device_query`. Si es array vacío o no se proporciona, consulta todos los tipos.
- **Retorna**: (String) Información de estado de dispositivos formateada en Markdown.

### device_log_query

- **Descripción**: Consultar logs de dispositivos.
- **Parámetros**:
  - `endpoint_ids` (Array<Integer>, requerido): Lista de IDs de dispositivos para consultar historial.
  - `start_datetime` (String, opcional): Tiempo de inicio de consulta en formato `YYYY-MM-DD HH:MM:SS` (ej., `"2023-05-16 12:00:00"`).
  - `end_datetime` (String, opcional): Tiempo de fin de consulta en formato `YYYY-MM-DD HH:MM:SS`.
  - `attribute` (String, opcional): Nombre de atributo específico de dispositivo a consultar (ej., `on_off`, `brightness`). Si no se proporciona, consulta todos los atributos registrados para el dispositivo.
- **Retorna**: (String) Información de estado histórico de dispositivos formateada en Markdown. (Nota: La implementación actual puede mostrar "This feature will be available soon.", indicando que la característica está pendiente de finalización.)

### run_scenes

- **Descripción**: Ejecutar escenas especificadas por ID de escena.
- **Parámetros**:
  - `scenes` (Array<Integer>, requerido): Lista de IDs de escenas a ejecutar.
- **Retorna**: (String) Mensaje de resultado de ejecución de escena.

### get_scenes

- **Descripción**: Consultar todas las escenas en el hogar del usuario, o escenas dentro de habitaciones especificadas.
- **Parámetros**:
  - `positions` (Array<String>, opcional): Lista de nombres de habitaciones. Si es array vacío o no se proporciona, consulta escenas para todo el hogar.
- **Retorna**: (String) Información de escenas formateada en Markdown.

### automation_config

- **Descripción**: Configurar tareas de control de dispositivos programadas o diferidas.
- **Parámetros**:
  - `scheduled_time` (String, requerido): Punto de tiempo establecido (para tareas diferidas, convertido basado en tiempo actual), formato `YYYY-MM-DD HH:MM:SS` (ej., `"2025-05-16 12:12:12"`).
  - `endpoint_ids` (Array<Integer>, requerido): Lista de IDs de dispositivos para control programado.
  - `control_params` (Object, requerido): Parámetros de control de dispositivo usando el mismo formato que la herramienta `device_control` (incluyendo action, attribute, value, etc.).
- **Retorna**: (String) Mensaje de resultado de configuración de automatización.

### get_homes

- **Descripción**: Obtener todas las listas de hogares bajo la cuenta de usuario.
- **Parámetros**: Ninguno.
- **Retorna**: (String) Lista de nombres de hogares separados por comas. Retorna cadena vacía o mensaje apropiado si no hay datos.

### switch_home

- **Descripción**: Cambiar el hogar de operación actual del usuario. Después del cambio, consultas de dispositivos, controles, y otras operaciones posteriores apuntarán al hogar recién cambiado.
- **Parámetros**:
  - `home_name` (String, requerido): Nombre del hogar objetivo (debe provenir de la lista disponible proporcionada por la herramienta `get_homes`).
- **Retorna**: (String) Mensaje de resultado de operación de cambio.

## Estructura del Proyecto

```
.
├── cmd.go                # Definiciones de comandos CLI Cobra y punto de entrada del programa (contiene función main)
├── server.go             # Lógica central del servidor MCP, definiciones de herramientas y manejo de solicitudes
├── smh.go                # Wrapper de interfaz API de plataforma domótica Aqara
├── middleware.go         # Middleware: autenticación de usuario, control de timeout, recuperación de excepciones
├── config.go             # Gestión de configuración global y manejo de variables de entorno
├── go.mod                # Archivo de gestión de dependencias del módulo Go
├── go.sum                # Archivo de suma de verificación de dependencias del módulo Go
├── img/                  # Recursos de imágenes usadas en documentación README
├── LICENSE               # Licencia de código abierto MIT
└── README.md             # Documentación del proyecto
```

### Descripciones de Archivos Principales

- **`cmd.go`**: Implementación CLI basada en framework Cobra, definiendo modos de inicio `run stdio` y `run sse` y función de entrada principal
- **`server.go`**: Implementación central del servidor MCP, responsable de registro de herramientas, manejo de solicitudes, y soporte de protocolo
- **`smh.go`**: Capa wrapper de API de plataforma domótica Aqara, proporcionando control de dispositivos, autenticación, y soporte multi-hogar
- **`middleware.go`**: Middleware de procesamiento de solicitudes proporcionando validación de autenticación, control de timeout, y manejo de excepciones
- **`config.go`**: Gestión de configuración global, responsable del manejo de variables de entorno y configuración de API

## Contribuir

¡Bienvenido a contribuir al proyecto enviando Issues o Pull Requests!

Antes de enviar código, por favor asegúrese de que:
1. El código sigue los estándares de codificación del lenguaje Go.
2. Las definiciones de herramientas MCP relacionadas y interfaces de prompt mantienen consistencia y claridad.
3. Agregar o actualizar pruebas unitarias para cubrir sus cambios.
4. Actualizar documentación relevante (como este README) si es necesario.
5. Asegurar que sus mensajes de commit sean claros y descriptivos.

## Licencia

Este proyecto está licenciado bajo la [Licencia MIT](/LICENSE).
Copyright (c) 2025 Aqara-Copliot