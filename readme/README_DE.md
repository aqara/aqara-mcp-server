<div align="center" style="display: flex; align-items: center; justify-content: center; ">

  <img src="/readme/img/logo.png" alt="Aqara Logo" height="120">
  <h1>MCP Server</h1>

</div>

<div align="center">

[English](/readme/README.md) | [中文](/readme/README_CN.md) | [繁體中文](/readme/README_CHT.md) | [Français](/readme/README_FR.md) | [한국어](/readme/README_KR.md) | [Español](/readme/README_ES.md) | [日本語](/readme/README_JP.md) | Deutsch | [Italiano](/readme/README_IT.md)

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/aqara/aqara-mcp-server)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![Release](https://img.shields.io/github/v/release/aqara/aqara-mcp-server)](https://github.com/aqara/aqara-mcp-server/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

Aqara MCP Server ist ein Smart-Home-Steuerungsdienst, der auf dem [MCP (Model Context Protocol)](https://modelcontextprotocol.io/introduction) Protokoll entwickelt wurde. Er ermöglicht es jedem KI-Assistenten oder jeder API, die das MCP-Protokoll unterstützt (wie Claude, Cursor usw.), mit Ihren Aqara Smart-Home-Geräten zu interagieren und Funktionen wie Gerätesteuerung über natürliche Sprache, Statusabfragen, Szenenausführung usw. zu realisieren.

## Inhaltsverzeichnis

- [Inhaltsverzeichnis](#inhaltsverzeichnis)
- [Funktionen](#funktionen)
- [Funktionsweise](#funktionsweise)
- [Schnellstart](#schnellstart)
  - [Voraussetzungen](#voraussetzungen)
  - [Installation](#installation)
    - [Methode 1: Vorkompilierte Version herunterladen (Empfohlen)](#methode-1-vorkompilierte-version-herunterladen-empfohlen)
    - [Methode 2: Aus Quellcode kompilieren](#methode-2-aus-quellcode-kompilieren)
  - [Aqara-Kontonauthentifizierung](#aqara-kontonauthentifizierung)
  - [Client-Konfiguration](#client-konfiguration)
    - [Claude for Desktop Konfigurationsbeispiel](#claude-for-desktop-konfigurationsbeispiel)
    - [Beschreibung der Konfigurationsparameter](#beschreibung-der-konfigurationsparameter)
    - [Andere MCP-Clients](#andere-mcp-clients)
  - [Service starten](#service-starten)
    - [Standardmodus (Empfohlen)](#standardmodus-empfohlen)
    - [HTTP-Modus (Optional)](#http-modus-optional)
- [API-Tools Beschreibung](#api-tools-beschreibung)
  - [Gerätesteuerung](#gerätesteuerung)
    - [device\_control](#device_control)
  - [Geräteabfrage](#geräteabfrage)
    - [device\_query](#device_query)
    - [device\_status\_query](#device_status_query)
    - [device\_log\_query](#device_log_query)
  - [Szenenverwaltung](#szenenverwaltung)
    - [get\_scenes](#get_scenes)
    - [run\_scenes](#run_scenes)
  - [Heimverwaltung](#heimverwaltung)
    - [get\_homes](#get_homes)
    - [switch\_home](#switch_home)
  - [Automatisierungskonfiguration](#automatisierungskonfiguration)
    - [automation\_config](#automation_config)
- [Projektstruktur](#projektstruktur)
  - [Verzeichnisstruktur](#verzeichnisstruktur)
  - [Beschreibung der Hauptdateien](#beschreibung-der-hauptdateien)
- [Entwicklungshandbuch](#entwicklungshandbuch)
- [Lizenz](#lizenz)

## Funktionen

- **Umfassende Gerätekontrolle**: Unterstützt präzise Kontrolle verschiedener Attribute wie Ein/Aus, Helligkeit, Farbtemperatur, Modus usw. für Aqara Smart-Geräte
- **Flexible Geräteabfrage**: Möglichkeit, Gerätelisten und deren detaillierte Zustände nach Raum und Gerätetyp abzufragen
- **Intelligente Szenenverwaltung**: Unterstützt Abfrage und Ausführung vom Benutzer vorkonfigurierter Smart-Home-Szenen
- **Geräteverlauf**: Abfrage historischer Statusänderungsaufzeichnungen von Geräten in bestimmten Zeitbereichen
- **Automatisierungskonfiguration**: Unterstützt Konfiguration geplanter oder verzögerter Gerätesteuerungsaufgaben
- **Multi-Home-Unterstützung**: Unterstützt Abfrage und Wechsel zwischen verschiedenen Heimen unter dem Benutzerkonto
- **MCP-Protokoll-Kompatibilität**: Vollständige Einhaltung der MCP-Protokollspezifikationen, einfache Integration mit verschiedenen KI-Assistenten
- **Sicherer Authentifizierungsmechanismus**: Verwendet sichere Authentifizierung basierend auf Login-Autorisierung + Signatur zum Schutz von Benutzerdaten und Gerätesicherheit
- **Plattformübergreifende Ausführung**: Entwickelt in Go, kann zu ausführbaren Dateien für mehrere Plattformen kompiliert werden
- **Einfach erweiterbar**: Modulares Design ermöglicht einfaches Hinzufügen neuer Tools und Funktionalitäten

## Funktionsweise

Aqara MCP Server fungiert als Brücke zwischen KI-Assistenten und der Aqara Smart-Home-Plattform:

1. **KI-Assistent (MCP-Client)**: Benutzer erteilt Befehle über den KI-Assistenten (z.B. "Schalte das Licht im Wohnzimmer ein")
2. **MCP-Client**: Analysiert Benutzerbefehle und ruft entsprechende Tools auf, die vom Aqara MCP Server gemäß MCP-Protokoll bereitgestellt werden (z.B. `device_control`)
3. **Aqara MCP Server (dieses Projekt)**: Empfängt Anfragen vom Client, validiert sie und ruft das `smh.go` Modul auf
4. **`smh.go` Modul**: Verwendet konfigurierte Aqara-Anmeldedaten zur Kommunikation mit der Aqara Cloud-API und führt tatsächliche Geräteoperationen oder Datenabfragen aus
5. **Antwortfluss**: Aqara Cloud-API gibt Ergebnisse zurück, die über Aqara MCP Server an den MCP-Client weitergeleitet und schließlich dem Benutzer präsentiert werden

## Schnellstart

### Voraussetzungen

- Go (Version 1.24 oder höher)
- Git (zum Kompilieren aus Quellcode)
- Aqara-Konto und verbundene Smart-Geräte

### Installation

Sie können wählen, vorkompilierte ausführbare Dateien herunterzuladen oder aus dem Quellcode zu kompilieren.

#### Methode 1: Vorkompilierte Version herunterladen (Empfohlen)

Besuchen Sie die GitHub Releases-Seite, um die neueste ausführbare Datei für Ihr Betriebssystem herunterzuladen:

**📥 [Zur Releases-Seite für Download gehen](https://github.com/aqara/aqara-mcp-server/releases)**

Nach dem Herunterladen der entsprechenden Archivdatei für Ihre Plattform entpacken Sie sie und sie ist einsatzbereit.

#### Methode 2: Aus Quellcode kompilieren

```bash
# Repository klonen
git clone https://github.com/aqara/aqara-mcp-server.git
cd aqara-mcp-server

# Abhängigkeiten herunterladen
go mod tidy

# Ausführbare Datei kompilieren
go build -o aqara-mcp-server
```

Nach Abschluss der Kompilierung wird die ausführbare Datei `aqara-mcp-server` im aktuellen Verzeichnis generiert.

### Aqara-Kontonauthentifizierung

Damit der MCP Server auf Ihr Aqara-Konto zugreifen und Geräte steuern kann, müssen Sie zuerst die Login-Autorisierung abschließen.

Bitte besuchen Sie die folgende Adresse, um die Login-Autorisierung abzuschließen:
**🔗 [https://cdn.aqara.com/app/mcpserver/login.html](https://cdn.aqara.com/app/mcpserver/login.html)**

Nach erfolgreichem Login erhalten Sie die notwendigen Authentifizierungsinformationen (wie `token`, `region`), die in den nachfolgenden Konfigurationsschritten verwendet werden.

> ⚠️ **Sicherheitshinweis**: Bitte bewahren Sie die `token`-Informationen sicher auf und geben Sie sie nicht an andere weiter.

### Client-Konfiguration

Die Konfigurationsmethoden für verschiedene MCP-Clients unterscheiden sich leicht. Das Folgende ist ein Beispiel dafür, wie Claude for Desktop konfiguriert wird, um diesen MCP Server zu verwenden:

#### Claude for Desktop Konfigurationsbeispiel

1. Öffnen Sie die Einstellungen (Settings) von Claude for Desktop

    ![Claude Open Setting](/readme/img/opening_setting.png)

2. Wechseln Sie zum Developer-Tab und klicken Sie dann auf Konfiguration bearbeiten (Edit Config), verwenden Sie einen Texteditor, um die Konfigurationsdatei zu öffnen

    ![Claude Edit Configuration](/readme/img/edit_config.png)

3. Fügen Sie die Konfigurationsinformationen von der "erfolgreichen Login-Seite" zur Client-Konfigurationsdatei `claude_desktop_config.json` hinzu

    ![Configuration Example](/readme/img/config_info.png)

#### Beschreibung der Konfigurationsparameter

- `command`: Vollständiger Pfad zur heruntergeladenen oder kompilierten ausführbaren Datei `aqara-mcp-server`
- `args`: Verwenden Sie `["run", "stdio"]`, um den stdio-Transportmodus zu starten
- `env`: Umgebungsvariablen-Konfiguration
  - `token`: Zugriffstoken von der Aqara-Login-Seite erhalten
  - `region`: Die Region, in der sich Ihr Aqara-Konto befindet (wie CN, US, EU usw.)

#### Andere MCP-Clients

Für andere Clients, die das MCP-Protokoll unterstützen (wie ChatGPT, Cursor usw.), ist die Konfiguration ähnlich:

- Stellen Sie sicher, dass der Client das MCP-Protokoll unterstützt
- Konfigurieren Sie den Pfad zur ausführbaren Datei und Startparameter
- Setzen Sie die Umgebungsvariablen `token` und `region`
- Wählen Sie das geeignete Transportprotokoll (empfohlen wird `stdio`)

### Service starten

#### Standardmodus (Empfohlen)

Starten Sie Claude for Desktop neu. Dann können Sie Operationen wie Gerätekontrolle, Geräteabfrage, Szenenausführung usw. über natürliche Sprache ausführen.

![Claude Chat Example](/readme/img/claude.png)

#### HTTP-Modus (Optional)

Wenn Sie den HTTP-Modus verwenden müssen, können Sie ihn so starten:

```bash
# Standard-Port 8080 verwenden
./aqara-mcp-server run http

# Oder benutzerdefinierten Host und Port angeben
./aqara-mcp-server run http --host localhost --port 9000
```

Verwenden Sie dann die Parameter `["run", "http"]` in der Client-Konfiguration.

## API-Tools Beschreibung

MCP-Clients können mit Aqara Smart-Home-Geräten interagieren, indem sie diese Tools aufrufen.

### Gerätesteuerung

#### device_control

Steuert den Zustand oder die Attribute von Smart-Home-Geräten (wie Ein/Aus, Temperatur, Helligkeit, Farbe, Farbtemperatur usw.).

**Parameter:**

- `endpoint_ids` _(Array\<Integer\>, erforderlich)_: Liste der zu steuernden Geräte-IDs
- `control_params` _(Object, erforderlich)_: Steuerungsparameter-Objekt, enthält spezifische Operationen:
  - `action` _(String, erforderlich)_: Auszuführende Operation (wie `"on"`, `"off"`, `"set"`, `"up"`, `"down"`, `"cooler"`, `"warmer"`)
  - `attribute` _(String, erforderlich)_: Zu steuerndes Geräteattribut (wie `"on_off"`, `"brightness"`, `"color_temperature"`, `"ac_mode"`)
  - `value` _(String | Number, optional)_: Zielwert (erforderlich wenn action "set" ist)
  - `unit` _(String, optional)_: Einheit des Wertes (wie `"%"`, `"K"`, `"℃"`)

**Rückgabe:** Ergebnismeldung der Gerätekontrolloperation

### Geräteabfrage

#### device_query

Ruft Geräteliste basierend auf angegebener Position (Raum) und Gerätetyp ab (enthält keine Echtzeitstatusinformationen).

**Parameter:**

- `positions` _(Array\<String\>, optional)_: Liste der Raumnamen. Leeres Array bedeutet alle Räume abfragen
- `device_types` _(Array\<String\>, optional)_: Liste der Gerätetypen (wie `"Light"`, `"WindowCovering"`, `"AirConditioner"`, `"Button"`). Leeres Array bedeutet alle Typen abfragen

**Rückgabe:** Geräteliste im Markdown-Format, enthält Gerätenamen und IDs

#### device_status_query

Ruft aktuelle Statusinformationen von Geräten ab (wird verwendet, um Echtzeitstatusinformationen wie Farbe, Helligkeit, Ein/Aus usw. abzufragen).

**Parameter:**

- `positions` _(Array\<String\>, optional)_: Liste der Raumnamen. Leeres Array bedeutet alle Räume abfragen
- `device_types` _(Array\<String\>, optional)_: Liste der Gerätetypen. Optionale Werte gleich `device_query`. Leeres Array bedeutet alle Typen abfragen

**Rückgabe:** Gerätestatusinformationen im Markdown-Format

#### device_log_query

Fragt historische Protokollinformationen von Geräten ab.

**Parameter:**

- `endpoint_ids` _(Array\<Integer\>, erforderlich)_: Liste der Geräte-IDs, für die Verlauf abgefragt werden soll
- `start_datetime` _(String, optional)_: Startzeit der Abfrage, Format `YYYY-MM-DD HH:MM:SS` (Beispiel: `"2023-05-16 12:00:00"`)
- `end_datetime` _(String, optional)_: Endzeit der Abfrage, Format `YYYY-MM-DD HH:MM:SS`
- `attribute` _(String, optional)_: Spezifischer Name des abzufragenden Geräteattributs (wie `on_off`, `brightness`). Wenn nicht angegeben, werden alle aufgezeichneten Attribute abgefragt

**Rückgabe:** Historische Gerätestatusinformationen im Markdown-Format

> 📝 **Hinweis:** Die aktuelle Implementierung zeigt möglicherweise "This feature will be available soon." an, was bedeutet, dass die Funktionalität noch vervollständigt werden muss.

### Szenenverwaltung

#### get_scenes

Fragt alle Szenen unter dem Benutzerheim oder Szenen in bestimmten Räumen ab.

**Parameter:**

- `positions` _(Array\<String\>, optional)_: Liste der Raumnamen. Leeres Array bedeutet Szenen des gesamten Heims abfragen

**Rückgabe:** Szeneninformationen im Markdown-Format

#### run_scenes

Führt bestimmte Szenen basierend auf Szenen-IDs aus.

**Parameter:**

- `scenes` _(Array\<Integer\>, erforderlich)_: Liste der auszuführenden Szenen-IDs

**Rückgabe:** Ergebnismeldung der Szenenausführung

### Heimverwaltung

#### get_homes

Ruft Liste aller Heime unter dem Benutzerkonto ab.

**Parameter:** Keine

**Rückgabe:** Durch Kommas getrennte Liste von Heimnamen. Wenn keine Daten vorhanden sind, wird ein leerer String oder eine entsprechende Informationsmeldung zurückgegeben

#### switch_home

Wechselt das aktuell vom Benutzer betriebene Heim. Nach dem Wechsel richten sich nachfolgende Operationen wie Geräteabfrage, Kontrolle usw. an das neu gewechselte Heim.

**Parameter:**

- `home_name` _(String, erforderlich)_: Name des Zielheims

**Rückgabe:** Ergebnismeldung der Wechseloperation

### Automatisierungskonfiguration

#### automation_config

Konfiguriert geplante oder verzögerte Gerätesteuerungsaufgaben (unterstützt derzeit nur feste Timer-Automatisierungskonfiguration).

**Parameter:**

- `scheduled_time` _(String, erforderlich)_: Eingestellter Zeitpunkt (wenn es sich um eine verzögerte Aufgabe handelt, wird sie basierend auf dem aktuellen Zeitpunkt umgewandelt), Format `YYYY-MM-DD HH:MM:SS` (Beispiel: `"2025-05-16 12:12:12"`)
- `endpoint_ids` _(Array\<Integer\>, erforderlich)_: Liste der Geräte-IDs für geplante Kontrolle
- `control_params` _(Object, erforderlich)_: Gerätesteuerungsparameter, verwendet das gleiche Format wie das `device_control` Tool (enthält action, attribute, value usw.)

**Rückgabe:** Ergebnismeldung der Automatisierungskonfiguration

> 📝 **Hinweis:** Die aktuelle Implementierung zeigt möglicherweise "This feature will be available soon." an, was bedeutet, dass die Funktionalität noch vervollständigt werden muss.

## Projektstruktur

### Verzeichnisstruktur

```text
.
├── cmd.go                # Cobra CLI-Befehlsdefinition und Programmeinstiegspunkt (enthält main-Funktion)
├── server.go             # MCP-Server-Kernlogik, Tool-Definition und Anfrageverarbeitung
├── smh.go                # Aqara Smart-Home-Plattform API-Schnittstellenkapselung
├── middleware.go         # Middleware: Benutzerauthentifizierung, Timeout-Kontrolle, Ausnahmewiederherstellung
├── config.go             # Globale Konfigurationsverwaltung und Umgebungsvariablenverarbeitung
├── go.mod                # Go-Modul-Abhängigkeitsverwaltungsdatei
├── go.sum                # Go-Modul-Abhängigkeitsprüfsummendatei
├── readme/               # README-Dokumentation und Bildressourcen
│   ├── img/              # Bildressourcenverzeichnis
│   └── *.md              # Mehrsprachige README-Dateien
├── LICENSE               # MIT Open-Source-Lizenz
└── README.md             # Hauptprojektdokumentation
```

### Beschreibung der Hauptdateien

- **`cmd.go`**: CLI-Implementierung basierend auf dem Cobra-Framework, definiert `run stdio` und `run http` Startmodi und Haupteinstiegsfunktion
- **`server.go`**: MCP-Server-Kernimplementierung, verantwortlich für Tool-Registrierung, Anfrageverarbeitung und Protokollunterstützung
- **`smh.go`**: Aqara Smart-Home-Plattform API-Kapselungsschicht, bietet Gerätekontrolle, Authentifizierung und Multi-Home-Unterstützung
- **`middleware.go`**: Anfrageverarbeitungs-Middleware, bietet Authentifizierungsvalidierung, Timeout-Kontrolle und Ausnahmebehandlung
- **`config.go`**: Globale Konfigurationsverwaltung, verantwortlich für Umgebungsvariablenverarbeitung und API-Konfiguration

## Entwicklungshandbuch

Willkommen zur Teilnahme an Projektbeiträgen durch Einreichen von Issues oder Pull Requests!

Bevor Sie Code einreichen, stellen Sie bitte sicher, dass:

1. Der Code den Go-Sprach-Codierungsstandards folgt
2. Verwandte MCP-Tools und Schnittstellendefinitionen Konsistenz und Klarheit beibehalten
3. Hinzufügen oder Aktualisieren von Unit-Tests zur Abdeckung Ihrer Änderungen
4. Falls notwendig, aktualisieren Sie die zugehörige Dokumentation (wie diese README)
5. Stellen Sie sicher, dass Ihre Commit-Nachrichten klar und verständlich sind

## Lizenz

Dieses Projekt ist unter der [MIT License](/LICENSE) lizenziert.

Copyright (c) 2025 Aqara-Copilot
