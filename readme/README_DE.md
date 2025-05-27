# Aqara MCP Server

[English](/readme/README.md) | [中文](/readme/README_CN.md) | [繁體中文](/readme/README_CHT.md) | [Français](/readme/README_FR.md) | [한국어](/readme/README_KR.md) | [Español](/readme/README_ES.md) | [日本語](/readme/README_JP.md) | Deutsch | [Italiano](/readme/README_IT.md)

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/aqara/aqara-mcp-server)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![Release](https://img.shields.io/github/v/release/aqara/aqara-mcp-server)](https://github.com/aqara/aqara-mcp-server/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Aqara MCP Server ist ein Smart-Home-Steuerungsservice, der auf dem [MCP (Model Context Protocol)](https://modelcontextprotocol.io/introduction) Protokoll basiert. Es ermöglicht jedem AI-Assistenten oder jeder API, die das MCP-Protokoll unterstützt (wie Claude, ChatGPT, Cursor usw.), mit Ihren Aqara Smart-Home-Geräten zu interagieren und Funktionen wie natürlichsprachliche Gerätekontrolle, Statusabfragen und Szenenausführung zu realisieren.

## Inhaltsverzeichnis

- [Aqara MCP Server](#aqara-mcp-server)
  - [Inhaltsverzeichnis](#inhaltsverzeichnis)
  - [Funktionen](#funktionen)
  - [Funktionsweise](#funktionsweise)
  - [Schnellstart](#schnellstart)
    - [Voraussetzungen](#voraussetzungen)
    - [Installation](#installation)
    - [Aqara-Konto-Authentifizierung](#aqara-konto-authentifizierung)
    - [Konfigurationsbeispiel (Claude for Desktop)](#konfigurationsbeispiel-claude-for-desktop)
    - [Service ausführen](#service-ausführen)
  - [Verfügbare Tools](#verfügbare-tools)
    - [device\_control](#device_control)
    - [device\_query](#device_query)
    - [device\_status\_query](#device_status_query)
    - [device\_log\_query](#device_log_query)
    - [run\_scenes](#run_scenes)
    - [get\_scenes](#get_scenes)
    - [automation\_config](#automation_config)
    - [get\_homes](#get_homes)
    - [switch\_home](#switch_home)
  - [Projektstruktur](#projektstruktur)
    - [Kern-Datei-Erklärungen](#kern-datei-erklärungen)
  - [Beitragsleitfaden](#beitragsleitfaden)
  - [Lizenz](#lizenz)

## Funktionen

-   **Umfassende Gerätekontrolle**: Unterstützt die präzise Kontrolle verschiedener Eigenschaften von Aqara Smart-Geräten wie Ein/Aus, Helligkeit, Farbtemperatur, Modi usw.
-   **Flexible Geräteabfrage**: Kann Gerätelisten und deren detaillierte Status nach Raum und Gerätetyp abfragen.
-   **Intelligentes Szenenmanagement**: Unterstützt das Abfragen und Ausführen von benutzerdefinierten Smart-Home-Szenen.
-   **Geräteverlauf**: Abfrage von Gerätestatusänderungen in einem bestimmten Zeitraum.
-   **Automatisierungskonfiguration**: Unterstützt die Konfiguration von zeitgesteuerten oder verzögerten Gerätekontrollaufgaben.
-   **Multi-Home-Unterstützung**: Unterstützt das Abfragen und Wechseln zwischen verschiedenen Haushalten unter dem Benutzerkonto.
-   **MCP-Protokoll-Kompatibilität**: Vollständige Einhaltung der MCP-Protokoll-Spezifikation, einfache Integration mit verschiedenen AI-Assistenten.
-   **Sicherer Authentifizierungsmechanismus**: Verwendet login-basierte Autorisierung + Signatur-basierte sichere Authentifizierung zum Schutz von Benutzerdaten und Gerätesicherheit.
-   **Plattformübergreifende Ausführung**: Entwickelt in Go-Sprache, kann zu Multi-Plattform-Executables kompiliert werden.
-   **Einfach erweiterbar**: Modulares Design, ermöglicht einfaches Hinzufügen neuer Tools und Funktionen.

## Funktionsweise

Aqara MCP Server fungiert als Brücke zwischen AI-Assistenten und der Aqara Smart-Home-Plattform:

1.  **AI-Assistent (MCP-Client)**: Benutzer gibt Befehle über den AI-Assistenten aus (z.B. "Schalte das Licht im Wohnzimmer ein").
2.  **MCP-Client**: Analysiert Benutzerbefehle und ruft entsprechende Tools des Aqara MCP Servers basierend auf dem MCP-Protokoll auf (z.B. `device_control`).
3.  **Aqara MCP Server (dieses Projekt)**: Empfängt Anfragen vom Client, validiert sie und ruft das `smh.go`-Modul auf.
4.  **`smh.go`-Modul**: Verwendet konfigurierte Aqara-Anmeldedaten, kommuniziert mit der Aqara-Cloud-API und führt tatsächliche Geräteoperationen oder Datenabfragen durch.
5.  **Antwortfluss**: Aqara-Cloud-API gibt Ergebnisse zurück, die über den Aqara MCP Server an den MCP-Client weitergegeben und schließlich dem Benutzer präsentiert werden.

## Schnellstart

### Voraussetzungen

-   Go (Version 1.24 oder höher)
-   Git (für das Erstellen aus dem Quellcode)
-   Aqara-Konto mit gebundenen Smart-Geräten

### Installation

Sie können vorkompilierte ausführbare Dateien herunterladen oder aus dem Quellcode erstellen.

**Option 1: Vorkompilierte Version herunterladen (Empfohlen)**

Besuchen Sie den folgenden Link und laden Sie das neueste ausführbare Dateipaket für Ihr Betriebssystem herunter.

[Releases-Seite](https://github.com/aqara/aqara-mcp-server/releases)

Nach dem Entpacken sofort verwendbar.

**Option 2: Aus Quellcode erstellen**

```bash
# Repository klonen
git clone https://github.com/aqara/aqara-mcp-server.git
cd aqara-mcp-server

# Abhängigkeiten herunterladen
go mod tidy

# Ausführbare Datei erstellen
go build -o aqara-mcp-server
```
Nach dem Build wird die ausführbare Datei `aqara-mcp-server` im aktuellen Verzeichnis generiert.

### Aqara-Konto-Authentifizierung

Damit der MCP Server auf Ihr Aqara-Konto zugreifen und Geräte steuern kann, müssen Sie zuerst eine Login-Autorisierung durchführen.

Bitte besuchen Sie die folgende Adresse, um die Login-Autorisierung abzuschließen:
[https://cdn.aqara.com/app/mcpserver/login.html](https://cdn.aqara.com/app/mcpserver/login.html)

Nach erfolgreichem Login erhalten Sie die notwendigen Authentifizierungsinformationen (wie `token`, `region`), die in den nachfolgenden Konfigurationsschritten verwendet werden.

**Bewahren Sie diese Informationen sicher auf, insbesondere das `token` sollte nicht an Dritte weitergegeben werden.**

### Konfigurationsbeispiel (Claude for Desktop)

Die Konfigurationsmethoden verschiedener MCP-Clients unterscheiden sich leicht. Hier ist ein Beispiel für die Konfiguration von Claude for Desktop zur Verwendung dieses MCP Servers:

1.  Öffnen Sie die Einstellungen (Settings) von Claude for Desktop.
2.  Wechseln Sie zum Entwickler (Developer) Tab.
3.  Klicken Sie auf Konfiguration bearbeiten (Edit Config) und öffnen Sie die Konfigurationsdatei mit einem Texteditor.

    ![](/readme/img/setting0.png)
    ![](/readme/img/setting1.png)

4.  Fügen Sie die Konfigurationsinformationen von der "Login-Erfolg-Seite" zur Client-Konfigurationsdatei (claude_desktop_config.json) hinzu. Konfigurationsbeispiel:

    ![](/readme/img/config.png)

**Konfigurationserklärung:**
- `command`: Vollständiger Pfad zur heruntergeladenen oder erstellten `aqara-mcp-server` ausführbaren Datei
- `args`: Verwenden Sie `["run", "stdio"]` um den stdio-Übertragungsmodus zu starten
- `env`: Umgebungsvariablen-Konfiguration
  - `token`: Zugangstoken von der Aqara-Login-Seite erhalten
  - `region`: Die Region Ihres Aqara-Kontos (wie CN, US, EU usw.)

### Service ausführen

Starten Sie Claude for Desktop neu. Dann können Sie über Gespräche die vom MCP Server bereitgestellten Tools aufrufen, um Gerätekontrolle, Geräteabfragen und andere Operationen durchzuführen.

![](/readme/img/claude.png)

**Andere MCP-Client-Konfigurationen**

Für andere Clients, die das MCP-Protokoll unterstützen (wie Claude, ChatGPT, Cursor usw.), ist die Konfiguration ähnlich:
- Stellen Sie sicher, dass der Client das MCP-Protokoll unterstützt
- Konfigurieren Sie den Pfad der ausführbaren Datei und Startparameter
- Setzen Sie Umgebungsvariablen `token` und `region`
- Wählen Sie ein geeignetes Übertragungsprotokoll (empfohlen: `stdio`)

**SSE-Modus (Optional)**

Wenn Sie den SSE (Server-Sent Events) Modus verwenden müssen, können Sie so starten:

```bash
# Standard-Port 8080 verwenden
./aqara-mcp-server run sse

# Oder benutzerdefinierten Host und Port angeben
./aqara-mcp-server run sse --host localhost --port 9000
```

Verwenden Sie dann `["run", "sse"]` Parameter in der Client-Konfiguration.

## Verfügbare Tools

MCP-Clients können diese Tools aufrufen, um mit Aqara Smart-Home-Geräten zu interagieren.

### device_control

-   **Beschreibung**: Kontrolliert den Status oder die Eigenschaften von Smart-Home-Geräten (wie Ein/Aus, Temperatur, Helligkeit, Farbe, Farbtemperatur usw.).
-   **Parameter**:
    -   `endpoint_ids` (Array<Integer>, erforderlich): Liste der zu kontrollierenden Geräte-IDs.
    -   `control_params` (Object, erforderlich): Kontrollparameter-Objekt mit spezifischen Operationen.
        -   `action` (String, erforderlich): Auszuführende Operation. Beispiele: `"on"`, `"off"`, `"set"`, `"up"`, `"down"`, `"cooler"`, `"warmer"`.
        -   `attribute` (String, erforderlich): Zu kontrollierende Geräteeigenschaft. Beispiele: `"on_off"`, `"brightness"`, `"color_temperature"`, `"ac_mode"`.
        -   `value` (String | Number, optional): Zielwert (erforderlich wenn action "set" ist).
        -   `unit` (String, optional): Einheit des Wertes (z.B.: `"%"`, `"K"`, `"℃"`).
-   **Rückgabe**: (String) Ergebnismeldung der Gerätekontrolle.

### device_query

-   **Beschreibung**: Ruft Gerätelisten basierend auf angegebenen Standorten (Räumen) und Gerätetypen ab (ohne Echtzeitstatusinformationen, listet nur Geräte und deren IDs auf).
-   **Parameter**:
    -   `positions` (Array<String>, optional): Liste der Raumnamen. Wenn leer oder nicht bereitgestellt, werden alle Räume abgefragt.
    -   `device_types` (Array<String>, optional): Liste der Gerätetypen. Beispiele: `"Light"`, `"WindowCovering"`, `"AirConditioner"`, `"Button"` usw. Wenn leer oder nicht bereitgestellt, werden alle Typen abgefragt.
-   **Rückgabe**: (String) Geräteliste im Markdown-Format mit Gerätenamen und IDs.

### device_status_query

-   **Beschreibung**: Ruft aktuelle Statusinformationen von Geräten ab (für Abfragen von Farbe, Helligkeit, Ein/Aus und anderen statusbezogenen Eigenschaften).
-   **Parameter**:
    -   `positions` (Array<String>, optional): Liste der Raumnamen. Wenn leer oder nicht bereitgestellt, werden alle Räume abgefragt.
    -   `device_types` (Array<String>, optional): Liste der Gerätetypen. Optionen wie bei `device_query`. Wenn leer oder nicht bereitgestellt, werden alle Typen abgefragt.
-   **Rückgabe**: (String) Gerätestatusinformationen im Markdown-Format.

### device_log_query

-   **Beschreibung**: Fragt Gerätelogs ab.
-   **Parameter**:
    -   `endpoint_ids` (Array<Integer>, erforderlich): Liste der Geräte-IDs für Verlaufsabfragen.
    -   `start_datetime` (String, optional): Startzeit der Abfrage, Format `YYYY-MM-DD HH:MM:SS` (z.B.: `"2023-05-16 12:00:00"`).
    -   `end_datetime` (String, optional): Endzeit der Abfrage, Format `YYYY-MM-DD HH:MM:SS`.
    -   `attribute` (String, optional): Name der spezifischen Geräteeigenschaft zur Abfrage (z.B.: `on_off`, `brightness`). Wenn nicht bereitgestellt, werden alle aufgezeichneten Eigenschaftsverläufe für das Gerät abgefragt.
-   **Rückgabe**: (String) Geräteverlaufsstatus-Informationen im Markdown-Format. (Hinweis: Die aktuelle Implementierung zeigt möglicherweise "This feature will be available soon." an, was darauf hinweist, dass die Funktion noch vervollständigt wird.)

### run_scenes

-   **Beschreibung**: Führt angegebene Szenen basierend auf Szenen-IDs aus.
-   **Parameter**:
    -   `scenes` (Array<Integer>, erforderlich): Liste der auszuführenden Szenen-IDs.
-   **Rückgabe**: (String) Ergebnismeldung der Szenenausführung.

### get_scenes

-   **Beschreibung**: Fragt alle Szenen unter dem Benutzerhaushalt oder Szenen in angegebenen Räumen ab.
-   **Parameter**:
    -   `positions` (Array<String>, optional): Liste der Raumnamen. Wenn leer oder nicht bereitgestellt, werden Szenen des gesamten Haushalts abgefragt.
-   **Rückgabe**: (String) Szeneninformationen im Markdown-Format.

### automation_config

-   **Beschreibung**: Konfiguriert zeitgesteuerte oder verzögerte Gerätekontrollaufgaben.
-   **Parameter**:
    -   `scheduled_time` (String, erforderlich): Eingestellte Zeit (bei verzögerten Aufgaben, basierend auf aktueller Zeit umgewandelt), Format `YYYY-MM-DD HH:MM:SS` (z.B.: `"2025-05-16 12:12:12"`).
    -   `endpoint_ids` (Array<Integer>, erforderlich): Liste der zeitgesteuert zu kontrollierenden Geräte-IDs.
    -   `control_params` (Object, erforderlich): Gerätekontrollparameter, verwendet das gleiche Format wie das `device_control` Tool (einschließlich action, attribute, value usw.).
-   **Rückgabe**: (String) Ergebnismeldung der Automatisierungskonfiguration.

### get_homes

-   **Beschreibung**: Ruft alle Haushaltslisten unter dem Benutzerkonto ab.
-   **Parameter**: Keine.
-   **Rückgabe**: (String) Kommagetrennte Liste der Haushaltsnamen. Wenn keine Daten vorhanden sind, wird ein leerer String oder entsprechende Hinweismeldung zurückgegeben.

### switch_home

-   **Beschreibung**: Wechselt den aktuell vom Benutzer bedienten Haushalt. Nach dem Wechsel zielen nachfolgende Geräteabfragen, Kontrollen usw. auf den neu gewechselten Haushalt ab.
-   **Parameter**:
    -   `home_name` (String, erforderlich): Name des Zielhaushalts (sollte aus der verfügbaren Liste stammen, die vom `get_homes` Tool bereitgestellt wird).
-   **Rückgabe**: (String) Ergebnismeldung der Wechseloperation.

## Projektstruktur

```
.
├── cmd.go                # Cobra CLI Befehlsdefinition und Programm-Einstiegspunkt (enthält main-Funktion)
├── server.go             # MCP Server-Kernlogik, Tool-Definition und Anfrageverarbeitung
├── smh.go                # Aqara Smart-Home-Plattform API-Interface-Kapselung
├── middleware.go         # Middleware: Benutzerauthentifizierung, Timeout-Kontrolle, Ausnahmewiederherstellung
├── config.go             # Globale Konfigurationsverwaltung und Umgebungsvariablenverarbeitung
├── go.mod                # Go-Modul-Abhängigkeitsverwaltungsdatei
├── go.sum                # Go-Modul-Abhängigkeits-Prüfsummendatei
├── img/                  # Bildressourcen für README-Dokumentation
├── LICENSE               # MIT Open-Source-Lizenz
└── README.md             # Projektdokumentation
```

### Kern-Datei-Erklärungen

-   **`cmd.go`**: Cobra-Framework-basierte CLI-Implementierung, definiert `run stdio` und `run sse` Startmodi sowie Haupt-Einstiegsfunktion
-   **`server.go`**: MCP Server-Kern-Implementierung, verantwortlich für Tool-Registrierung, Anfrageverarbeitung und Protokollunterstützung
-   **`smh.go`**: Aqara Smart-Home-Plattform API-Kapselungsschicht, bietet Gerätekontrolle, Authentifizierung und Multi-Home-Unterstützung
-   **`middleware.go`**: Anfrageverarbeitungs-Middleware, bietet Authentifizierungsvalidierung, Timeout-Kontrolle und Ausnahmebehandlung
-   **`config.go`**: Globale Konfigurationsverwaltung, verantwortlich für Umgebungsvariablenverarbeitung und API-Konfiguration

## Beitragsleitfaden

Wir begrüßen die Teilnahme am Projekt durch das Einreichen von Issues oder Pull Requests!

Vor dem Einreichen von Code stellen Sie bitte sicher:
1.  Der Code folgt den Go-Sprache-Codierungsstandards.
2.  Verwandte MCP-Tool- und Prompt-Interface-Definitionen bleiben konsistent und klar.
3.  Fügen Sie Unit-Tests hinzu oder aktualisieren Sie sie, um Ihre Änderungen abzudecken.
4.  Aktualisieren Sie bei Bedarf verwandte Dokumentation (wie diese README).
5.  Stellen Sie sicher, dass Ihre Commit-Nachrichten klar und verständlich sind.

## Lizenz

Dieses Projekt ist unter der [MIT License](/LICENSE) lizenziert.
Copyright (c) 2025 Aqara-Copliot