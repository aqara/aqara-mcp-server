<div align="center" style="display: flex; align-items: center; justify-content: center; ">

  <img src="/readme/img/logo.png" alt="Aqara Logo" height="120">
  <h1>MCP Server</h1>

</div>

<div align="center">

[English](/readme/README.md) | [中文](/readme/README_CN.md) | [繁體中文](/readme/README_CHT.md) | [Français](/readme/README_FR.md) | [한국어](/readme/README_KR.md) | [Español](/readme/README_ES.md) | [日本語](/readme/README_JP.md) | [Deutsch](/readme/README_DE.md) | Italiano

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/aqara/aqara-mcp-server)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![Release](https://img.shields.io/github/v/release/aqara/aqara-mcp-server)](https://github.com/aqara/aqara-mcp-server/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

Aqara MCP Server è un servizio di controllo domotico sviluppato basato sul protocollo [MCP (Model Context Protocol)](https://modelcontextprotocol.io/introduction). Consente a qualsiasi assistente AI o API che supporti il protocollo MCP (come Claude, Cursor, ecc.) di interagire con i dispositivi domotici Aqara, realizzando funzioni di controllo dispositivi tramite linguaggio naturale, interrogazione di stato, esecuzione di scene e altro.

## Indice

- [Indice](#indice)
- [Caratteristiche](#caratteristiche)
- [Principio di Funzionamento](#principio-di-funzionamento)
- [Avvio Rapido](#avvio-rapido)
  - [Prerequisiti](#prerequisiti)
  - [Installazione](#installazione)
    - [Metodo 1: Scarica Versione Precompilata (Raccomandato)](#metodo-1-scarica-versione-precompilata-raccomandato)
    - [Metodo 2: Compila dal Codice Sorgente](#metodo-2-compila-dal-codice-sorgente)
  - [Autenticazione Account Aqara](#autenticazione-account-aqara)
  - [Configurazione Client](#configurazione-client)
    - [Esempio Configurazione Claude for Desktop](#esempio-configurazione-claude-for-desktop)
    - [Spiegazione Parametri di Configurazione](#spiegazione-parametri-di-configurazione)
    - [Altri Client MCP](#altri-client-mcp)
  - [Avvio del Servizio](#avvio-del-servizio)
    - [Modalità Standard (Raccomandato)](#modalità-standard-raccomandato)
    - [Modalità HTTP (Opzionale)](#modalità-http-opzionale)
- [Spiegazione Strumenti API](#spiegazione-strumenti-api)
  - [Categoria Controllo Dispositivi](#categoria-controllo-dispositivi)
    - [device\_control](#device_control)
  - [Categoria Interrogazione Dispositivi](#categoria-interrogazione-dispositivi)
    - [device\_query](#device_query)
    - [device\_status\_query](#device_status_query)
    - [device\_log\_query](#device_log_query)
  - [Categoria Gestione Scene](#categoria-gestione-scene)
    - [get\_scenes](#get_scenes)
    - [run\_scenes](#run_scenes)
  - [Categoria Gestione Casa](#categoria-gestione-casa)
    - [get\_homes](#get_homes)
    - [switch\_home](#switch_home)
  - [Categoria Configurazione Automazione](#categoria-configurazione-automazione)
    - [automation\_config](#automation_config)
- [Struttura del Progetto](#struttura-del-progetto)
  - [Struttura Directory](#struttura-directory)
  - [Spiegazione File Principali](#spiegazione-file-principali)
- [Guida allo Sviluppo](#guida-allo-sviluppo)
- [Licenza](#licenza)

## Caratteristiche

- **Controllo Completo dei Dispositivi**: Supporta il controllo preciso di vari attributi dei dispositivi smart Aqara come accensione/spegnimento, luminosità, temperatura del colore, modalità, ecc.
- **Interrogazione Flessibile dei Dispositivi**: Capacità di interrogare elenchi dispositivi e loro stati dettagliati per stanza e tipo di dispositivo
- **Gestione Intelligente delle Scene**: Supporta l'interrogazione e l'esecuzione di scene domotiche preconfigurate dall'utente
- **Cronologia Dispositivi**: Interrogazione dei record di cambiamento stato dei dispositivi in un intervallo di tempo specificato
- **Configurazione Automazione**: Supporta la configurazione di compiti di controllo dispositivi temporizzati o ritardati
- **Supporto Multi-Casa**: Supporta l'interrogazione e il cambio tra diverse case sotto l'account utente
- **Compatibilità Protocollo MCP**: Completamente conforme alle specifiche del protocollo MCP, facile integrazione con vari assistenti AI
- **Meccanismo di Autenticazione Sicuro**: Adotta autenticazione sicura basata su login autorizzazione + firma, proteggendo dati utente e sicurezza dispositivi
- **Funzionamento Cross-Platform**: Sviluppato in linguaggio Go, compilabile per file eseguibili multi-piattaforma
- **Facile da Estendere**: Design modulare, consente di aggiungere facilmente nuovi strumenti e funzionalità

## Principio di Funzionamento

Aqara MCP Server funge da ponte tra assistenti AI e la piattaforma domotica Aqara:

1. **Assistente AI (Client MCP)**: L'utente emette comandi tramite l'assistente AI (es. "Accendi la luce del soggiorno")
2. **Client MCP**: Analizza le istruzioni dell'utente e chiama gli strumenti corrispondenti forniti da Aqara MCP Server secondo il protocollo MCP (es. `device_control`)
3. **Aqara MCP Server (questo progetto)**: Riceve richieste dal client, valida e chiama il modulo `smh.go`
4. **Modulo `smh.go`**: Usa le credenziali Aqara configurate per comunicare con l'API cloud Aqara, eseguendo operazioni effettive sui dispositivi o interrogazioni di dati
5. **Flusso di Risposta**: L'API cloud Aqara restituisce risultati, passati tramite Aqara MCP Server al client MCP, presentati infine all'utente

## Avvio Rapido

### Prerequisiti

- Go (versione 1.24 o superiore)
- Git (per compilazione dal codice sorgente)
- Account Aqara e dispositivi smart registrati

### Installazione

Puoi scegliere di scaricare file eseguibili precompilati o compilare dal codice sorgente.

#### Metodo 1: Scarica Versione Precompilata (Raccomandato)

Visita la pagina GitHub Releases e scarica l'ultimo file eseguibile adatto al tuo sistema operativo:

**📥 [Vai alla pagina Releases per il download](https://github.com/aqara/aqara-mcp-server/releases)**

Dopo aver scaricato il file compresso per la piattaforma corrispondente, decomprimilo per usarlo immediatamente.

#### Metodo 2: Compila dal Codice Sorgente

```bash
# Clona il repository
git clone https://github.com/aqara/aqara-mcp-server.git
cd aqara-mcp-server

# Scarica dipendenze
go mod tidy

# Compila il file eseguibile
go build -o aqara-mcp-server
```

Dopo la compilazione, verrà generato il file eseguibile `aqara-mcp-server` nella directory corrente.

### Autenticazione Account Aqara

Per consentire al MCP Server di accedere al tuo account Aqara e controllare i dispositivi, devi prima eseguire l'autorizzazione di login.

Visita il seguente indirizzo per completare l'autorizzazione di login:
**🔗 [https://cdn.aqara.com/app/mcpserver/login.html](https://cdn.aqara.com/app/mcpserver/login.html)**

Dopo un login riuscito, otterrai le informazioni di autenticazione necessarie (come `token`, `region`), che verranno utilizzate nei passaggi di configurazione successivi.

> ⚠️ **Promemoria di Sicurezza**: Conserva adeguatamente le informazioni `token`, non divulgarle ad altri.

### Configurazione Client

I metodi di configurazione differiscono leggermente tra i vari client MCP. Ecco un esempio di come configurare Claude for Desktop per utilizzare questo MCP Server:

#### Esempio Configurazione Claude for Desktop

1. Apri le impostazioni (Settings) di Claude for Desktop

    ![Claude Open Setting](/readme/img/opening_setting.png)

2. Passa alla scheda sviluppatore (Developer), poi clicca modifica configurazione (Edit Config), apri il file di configurazione con un editor di testo

    ![Claude Edit Configuration](/readme/img/edit_config.png)

3. Aggiungi le informazioni di configurazione dalla "pagina di login riuscito" al file di configurazione del client `claude_desktop_config.json`

    ![Configuration Example](/readme/img/config_info.png)

#### Spiegazione Parametri di Configurazione

- `command`: Percorso completo del file eseguibile `aqara-mcp-server` scaricato o compilato
- `args`: Usa `["run", "stdio"]` per avviare la modalità di trasporto stdio
- `env`: Configurazione variabili d'ambiente
  - `token`: Token di accesso ottenuto dalla pagina di login Aqara
  - `region`: Regione dell'account Aqara (come CN, US, EU, ecc.)

#### Altri Client MCP

Per altri client che supportano il protocollo MCP (come ChatGPT, Cursor, ecc.), il metodo di configurazione è simile:

- Assicurati che il client supporti il protocollo MCP
- Configura il percorso del file eseguibile e i parametri di avvio
- Imposta le variabili d'ambiente `token` e `region`
- Scegli il protocollo di trasporto appropriato (si raccomanda `stdio`)

### Avvio del Servizio

#### Modalità Standard (Raccomandato)

Riavvia Claude for Desktop. Successivamente potrai eseguire operazioni di controllo dispositivi, interrogazione dispositivi, esecuzione scene tramite linguaggio naturale.

![Claude Chat Example](/readme/img/claude.png)

#### Modalità HTTP (Opzionale)

Se devi usare la modalità HTTP, puoi avviarla così:

```bash
# Usa la porta predefinita 8080
./aqara-mcp-server run http

# O specifica host e porta personalizzati
./aqara-mcp-server run http --host localhost --port 9000
```

Poi usa i parametri `["run", "http"]` nella configurazione del client.

## Spiegazione Strumenti API

I client MCP possono interagire con i dispositivi domotici Aqara chiamando questi strumenti.

### Categoria Controllo Dispositivi

#### device_control

Controlla lo stato o gli attributi dei dispositivi domotici (es. accensione/spegnimento, temperatura, luminosità, colore, temperatura del colore, ecc.).

**Parametri:**

- `endpoint_ids` _(Array\<Integer\>, richiesto)_: Lista degli ID dei dispositivi da controllare
- `control_params` _(Object, richiesto)_: Oggetto parametri di controllo, contenente operazioni specifiche:
  - `action` _(String, richiesto)_: Operazione da eseguire (come `"on"`, `"off"`, `"set"`, `"up"`, `"down"`, `"cooler"`, `"warmer"`)
  - `attribute` _(String, richiesto)_: Attributo del dispositivo da controllare (come `"on_off"`, `"brightness"`, `"color_temperature"`, `"ac_mode"`)
  - `value` _(String | Number, opzionale)_: Valore target (richiesto quando action è "set")
  - `unit` _(String, opzionale)_: Unità del valore (come `"%"`, `"K"`, `"℃"`)

**Restituisce:** Messaggio del risultato dell'operazione di controllo del dispositivo

### Categoria Interrogazione Dispositivi

#### device_query

Ottiene la lista dei dispositivi basata su posizione specificata (stanza) e tipo di dispositivo (non include informazioni di stato in tempo reale).

**Parametri:**

- `positions` _(Array\<String\>, opzionale)_: Lista nomi stanze. Array vuoto significa interrogare tutte le stanze
- `device_types` _(Array\<String\>, opzionale)_: Lista tipi di dispositivo (come `"Light"`, `"WindowCovering"`, `"AirConditioner"`, `"Button"`). Array vuoto significa interrogare tutti i tipi

**Restituisce:** Lista dispositivi in formato Markdown, contenente nomi e ID dei dispositivi

#### device_status_query

Ottiene le informazioni di stato correnti dei dispositivi (per interrogare informazioni di stato in tempo reale come colore, luminosità, accensione/spegnimento).

**Parametri:**

- `positions` _(Array\<String\>, opzionale)_: Lista nomi stanze. Array vuoto significa interrogare tutte le stanze
- `device_types` _(Array\<String\>, opzionale)_: Lista tipi di dispositivo. Valori opzionali uguali a `device_query`. Array vuoto significa interrogare tutti i tipi

**Restituisce:** Informazioni di stato dei dispositivi in formato Markdown

#### device_log_query

Interroga le informazioni di log storico dei dispositivi.

**Parametri:**

- `endpoint_ids` _(Array\<Integer\>, richiesto)_: Lista degli ID dei dispositivi per i quali interrogare i record storici
- `start_datetime` _(String, opzionale)_: Tempo di inizio interrogazione, formato `YYYY-MM-DD HH:MM:SS` (es: `"2023-05-16 12:00:00"`)
- `end_datetime` _(String, opzionale)_: Tempo di fine interrogazione, formato `YYYY-MM-DD HH:MM:SS`
- `attribute` _(String, opzionale)_: Nome specifico dell'attributo del dispositivo da interrogare (come `on_off`, `brightness`). Se non fornito, interroga tutti gli attributi registrati

**Restituisce:** Informazioni di stato storico dei dispositivi in formato Markdown

> 📝 **Nota:** L'implementazione corrente potrebbe mostrare "This feature will be available soon.", indicando che la funzionalità è in attesa di perfezionamento.

### Categoria Gestione Scene

#### get_scenes

Interroga tutte le scene sotto la casa dell'utente, o le scene in stanze specificate.

**Parametri:**

- `positions` _(Array\<String\>, opzionale)_: Lista nomi stanze. Array vuoto significa interrogare le scene dell'intera casa

**Restituisce:** Informazioni delle scene in formato Markdown

#### run_scenes

Esegue scene specificate basate sugli ID delle scene.

**Parametri:**

- `scenes` _(Array\<Integer\>, richiesto)_: Lista degli ID delle scene da eseguire

**Restituisce:** Messaggio del risultato dell'esecuzione delle scene

### Categoria Gestione Casa

#### get_homes

Ottiene la lista di tutte le case sotto l'account utente.

**Parametri:** Nessuno

**Restituisce:** Lista di nomi delle case separati da virgole. Se non ci sono dati, restituisce stringa vuota o informazioni di prompt corrispondenti

#### switch_home

Cambia la casa che l'utente sta attualmente operando. Dopo il cambio, le successive operazioni di interrogazione dispositivi, controllo, ecc. saranno mirate alla nuova casa commutata.

**Parametri:**

- `home_name` _(String, richiesto)_: Nome della casa target

**Restituisce:** Messaggio del risultato dell'operazione di cambio

### Categoria Configurazione Automazione

#### automation_config

Configura compiti di controllo dispositivi temporizzati o ritardati (attualmente supporta solo la configurazione di automazione temporizzata e ritardata).

**Parametri:**

- `scheduled_time` _(String, richiesto)_: Punto temporale impostato (se è un compito ritardato, convertito basato sul punto temporale corrente), formato `YYYY-MM-DD HH:MM:SS` (es: `"2025-05-16 12:12:12"`)
- `endpoint_ids` _(Array\<Integer\>, richiesto)_: Lista degli ID dei dispositivi da controllare temporizzatamente
- `control_params` _(Object, richiesto)_: Parametri di controllo del dispositivo, usa lo stesso formato dello strumento `device_control` (include action, attribute, value, ecc.)

**Restituisce:** Messaggio del risultato della configurazione automazione

> 📝 **Nota:** L'implementazione corrente potrebbe mostrare "This feature will be available soon.", indicando che la funzionalità è in attesa di perfezionamento.

## Struttura del Progetto

### Struttura Directory

```text
.
├── cmd.go                # Definizione comandi Cobra CLI e punto di ingresso del programma (include funzione main)
├── server.go             # Logica principale del server MCP, definizione strumenti e gestione richieste
├── smh.go                # Incapsulamento interfaccia API piattaforma domotica Aqara
├── middleware.go         # Middleware: autenticazione utente, controllo timeout, recupero eccezioni
├── config.go             # Gestione configurazione globale e gestione variabili d'ambiente
├── go.mod                # File di gestione dipendenze modulo Go
├── go.sum                # File checksum dipendenze modulo Go
├── readme/               # Documentazione README e risorse immagini
│   ├── img/              # Directory risorse immagini
│   └── *.md              # File README multilingue
├── LICENSE               # Licenza open source MIT
└── README.md             # Documento principale del progetto
```

### Spiegazione File Principali

- **`cmd.go`**: Implementazione CLI basata su framework Cobra, definisce modalità di avvio `run stdio` e `run http` e funzione di ingresso principale
- **`server.go`**: Implementazione principale del server MCP, responsabile per registrazione strumenti, gestione richieste e supporto protocollo
- **`smh.go`**: Livello di incapsulamento API piattaforma domotica Aqara, fornisce controllo dispositivi, autenticazione e supporto multi-casa
- **`middleware.go`**: Middleware di gestione richieste, fornisce verifica autenticazione, controllo timeout e gestione eccezioni
- **`config.go`**: Gestione configurazione globale, responsabile per gestione variabili d'ambiente e configurazione API

## Guida allo Sviluppo

Benvenuto a partecipare alla contribuzione del progetto tramite invio di Issue o Pull Request!

Prima di inviare codice, assicurati che:

1. Il codice segua le specifiche di codifica del linguaggio Go
2. Mantenga coerenza e chiarezza delle definizioni di strumenti e interfacce MCP correlate
3. Aggiunga o aggiorni test unitari per coprire le tue modifiche
4. Se necessario, aggiorna la documentazione correlata (come questo README)
5. Assicurati che i tuoi messaggi di commit siano chiari e comprensibili

## Licenza

Questo progetto è autorizzato sotto [MIT License](/LICENSE).

Copyright (c) 2025 Aqara-Copilot
