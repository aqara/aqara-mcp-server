# Aqara MCP Server

[English](/readme/README.md) | [中文](/readme/README_CN.md) | [繁體中文](/readme/README_CHT.md) | [Français](/readme/README_FR.md) | [한국어](/readme/README_KR.md) | [Español](/readme/README_ES.md) | [日本語](/readme/README_JP.md) | [Deutsch](/readme/README_DE.md) | Italiano

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/aqara/aqara-mcp-server)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![Release](https://img.shields.io/github/v/release/aqara/aqara-mcp-server)](https://github.com/aqara/aqara-mcp-server/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Aqara MCP Server è un servizio di controllo della casa intelligente sviluppato basandosi sul protocollo [MCP (Model Context Protocol)](https://modelcontextprotocol.io/introduction). Permette a qualsiasi assistente AI o API che supporti il protocollo MCP (come Claude, ChatGPT, Cursor, ecc.) di interagire con i tuoi dispositivi smart home Aqara, realizzando funzioni come il controllo dei dispositivi tramite linguaggio naturale, query di stato, esecuzione di scene, ecc.

## Indice

- [Aqara MCP Server](#aqara-mcp-server)
  - [Indice](#indice)
  - [Caratteristiche](#caratteristiche)
  - [Principio di funzionamento](#principio-di-funzionamento)
  - [Avvio rapido](#avvio-rapido)
    - [Prerequisiti](#prerequisiti)
    - [Installazione](#installazione)
    - [Autenticazione account Aqara](#autenticazione-account-aqara)
    - [Esempio di configurazione (Claude for Desktop)](#esempio-di-configurazione-claude-for-desktop)
    - [Esecuzione del servizio](#esecuzione-del-servizio)
  - [Strumenti disponibili](#strumenti-disponibili)
    - [device\_control](#device_control)
    - [device\_query](#device_query)
    - [device\_status\_query](#device_status_query)
    - [device\_log\_query](#device_log_query)
    - [run\_scenes](#run_scenes)
    - [get\_scenes](#get_scenes)
    - [automation\_config](#automation_config)
    - [get\_homes](#get_homes)
    - [switch\_home](#switch_home)
  - [Struttura del progetto](#struttura-del-progetto)
    - [Spiegazione file core](#spiegazione-file-core)
  - [Guida ai contributi](#guida-ai-contributi)
  - [Licenza](#licenza)

## Caratteristiche

-   **Controllo completo dei dispositivi**: Supporta il controllo preciso di varie proprietà dei dispositivi smart Aqara come accensione/spegnimento, luminosità, temperatura colore, modalità, ecc.
-   **Query flessibili dei dispositivi**: Può interrogare elenchi di dispositivi e loro stati dettagliati per stanza e tipo di dispositivo.
-   **Gestione intelligente delle scene**: Supporta la query e l'esecuzione di scene smart home preimpostate dall'utente.
-   **Cronologia dei dispositivi**: Query delle modifiche di stato dei dispositivi in un intervallo di tempo specificato.
-   **Configurazione automazione**: Supporta la configurazione di attività di controllo dispositivi temporizzate o ritardate.
-   **Supporto multi-casa**: Supporta la query e il cambio tra diverse case sotto l'account utente.
-   **Compatibilità protocollo MCP**: Completa aderenza alle specifiche del protocollo MCP, facile integrazione con vari assistenti AI.
-   **Meccanismo di autenticazione sicuro**: Utilizza autenticazione basata su login + firma per proteggere i dati utente e la sicurezza dei dispositivi.
-   **Esecuzione multipiattaforma**: Sviluppato in linguaggio Go, può essere compilato in eseguibili multipiattaforma.
-   **Design facilmente estendibile**: Design modulare, facile aggiunta di nuovi strumenti e funzionalità.

## Principio di funzionamento

Aqara MCP Server funge da ponte tra assistenti AI e la piattaforma smart home Aqara:

1.  **Assistente AI (Client MCP)**: L'utente emette comandi tramite l'assistente AI (ad esempio, "Accendi la luce del soggiorno").
2.  **Client MCP**: Analizza i comandi utente e chiama gli strumenti corrispondenti di Aqara MCP Server basandosi sul protocollo MCP (ad esempio `device_control`).
3.  **Aqara MCP Server (questo progetto)**: Riceve richieste dal client, le valida e chiama il modulo `smh.go`.
4.  **Modulo `smh.go`**: Utilizza le credenziali Aqara configurate, comunica con l'API cloud Aqara ed esegue operazioni effettive sui dispositivi o query di dati.
5.  **Flusso di risposta**: L'API cloud Aqara restituisce risultati, che vengono trasmessi tramite Aqara MCP Server al client MCP e infine presentati all'utente.

## Avvio rapido

### Prerequisiti

-   Go (versione 1.24 o superiore)
-   Git (per la compilazione dal codice sorgente)
-   Account Aqara con dispositivi smart associati

### Installazione

Puoi scegliere di scaricare file eseguibili precompilati o compilare dal codice sorgente.

**Opzione 1: Scarica versione precompilata (Consigliata)**

Visita il link seguente per scaricare il pacchetto di file eseguibili più recente adatto al tuo sistema operativo.

[Pagina Releases](https://github.com/aqara/aqara-mcp-server/releases)

Pronto all'uso dopo la decompressione.

**Opzione 2: Compila dal codice sorgente**

```bash
# Clona il repository
git clone https://github.com/aqara/aqara-mcp-server.git
cd aqara-mcp-server

# Scarica le dipendenze
go mod tidy

# Compila il file eseguibile
go build -o aqara-mcp-server
```
Dopo la compilazione, verrà generato il file eseguibile `aqara-mcp-server` nella directory corrente.

### Autenticazione account Aqara

Per permettere al MCP Server di accedere al tuo account Aqara e controllare i dispositivi, devi prima effettuare l'autorizzazione di login.

Visita il seguente indirizzo per completare l'autorizzazione di login:
[https://cdn.aqara.com/app/mcpserver/login.html](https://cdn.aqara.com/app/mcpserver/login.html)

Dopo il login riuscito, otterrai le informazioni di autenticazione necessarie (come `token`, `region`), che verranno utilizzate nei passaggi di configurazione successivi.

**Conserva queste informazioni con cura, specialmente il `token` non deve essere divulgato ad altri.**

### Esempio di configurazione (Claude for Desktop)

I metodi di configurazione di diversi client MCP variano leggermente. Ecco un esempio di come configurare Claude for Desktop per utilizzare questo MCP Server:

1.  Apri le impostazioni (Settings) di Claude for Desktop.
2.  Passa alla scheda sviluppatore (Developer).
3.  Clicca su modifica configurazione (Edit Config) e apri il file di configurazione con un editor di testo.

    ![](/readme/img/setting0.png)
    ![](/readme/img/setting1.png)

4.  Aggiungi le informazioni di configurazione dalla "pagina di successo login" al file di configurazione del client (claude_desktop_config.json). Esempio di configurazione:

    ![](/readme/img/config.png)

**Spiegazione della configurazione:**
- `command`: Percorso completo al file eseguibile `aqara-mcp-server` scaricato o compilato
- `args`: Usa `["run", "stdio"]` per avviare la modalità di trasporto stdio
- `env`: Configurazione variabili d'ambiente
  - `token`: Token di accesso ottenuto dalla pagina di login Aqara
  - `region`: La regione del tuo account Aqara (come CN, US, EU, ecc.)

### Esecuzione del servizio

Riavvia Claude for Desktop. Poi puoi chiamare gli strumenti forniti dal MCP Server tramite conversazioni per eseguire controllo dispositivi, query dispositivi e altre operazioni.

![](/readme/img/claude.png)

**Configurazioni altri client MCP**

Per altri client che supportano il protocollo MCP (come Claude, ChatGPT, Cursor, ecc.), la configurazione è simile:
- Assicurati che il client supporti il protocollo MCP
- Configura il percorso del file eseguibile e i parametri di avvio
- Imposta le variabili d'ambiente `token` e `region`
- Scegli un protocollo di trasporto appropriato (consigliato: `stdio`)

**Modalità SSE (Opzionale)**

Se hai bisogno di utilizzare la modalità SSE (Server-Sent Events), puoi avviare così:

```bash
# Usa la porta predefinita 8080
./aqara-mcp-server run sse

# O specifica host e porta personalizzati
./aqara-mcp-server run sse --host localhost --port 9000
```

Poi usa i parametri `["run", "sse"]` nella configurazione del client.

## Strumenti disponibili

I client MCP possono chiamare questi strumenti per interagire con i dispositivi smart home Aqara.

### device_control

-   **Descrizione**: Controlla lo stato o le proprietà dei dispositivi smart home (come accensione/spegnimento, temperatura, luminosità, colore, temperatura colore, ecc.).
-   **Parametri**:
    -   `endpoint_ids` (Array<Integer>, richiesto): Lista degli ID dispositivi da controllare.
    -   `control_params` (Object, richiesto): Oggetto parametri di controllo, contenente operazioni specifiche.
        -   `action` (String, richiesto): Operazione da eseguire. Esempi: `"on"`, `"off"`, `"set"`, `"up"`, `"down"`, `"cooler"`, `"warmer"`.
        -   `attribute` (String, richiesto): Proprietà del dispositivo da controllare. Esempi: `"on_off"`, `"brightness"`, `"color_temperature"`, `"ac_mode"`.
        -   `value` (String | Number, opzionale): Valore target (richiesto quando action è "set").
        -   `unit` (String, opzionale): Unità del valore (ad es.: `"%"`, `"K"`, `"℃"`).
-   **Ritorno**: (String) Messaggio del risultato dell'operazione di controllo dispositivo.

### device_query

-   **Descrizione**: Ottieni elenchi di dispositivi basati su posizioni specificate (stanze) e tipi di dispositivo (non include informazioni di stato in tempo reale, elenca solo dispositivi e loro ID).
-   **Parametri**:
    -   `positions` (Array<String>, opzionale): Lista nomi stanze. Se array vuoto o non fornito, interroga tutte le stanze.
    -   `device_types` (Array<String>, opzionale): Lista tipi di dispositivo. Esempi: `"Light"`, `"WindowCovering"`, `"AirConditioner"`, `"Button"`, ecc. Se array vuoto o non fornito, interroga tutti i tipi.
-   **Ritorno**: (String) Elenco dispositivi in formato Markdown, contenente nomi e ID dispositivi.

### device_status_query

-   **Descrizione**: Ottieni informazioni di stato corrente dei dispositivi (per interrogare colore, luminosità, accensione/spegnimento e altre proprietà relative allo stato).
-   **Parametri**:
    -   `positions` (Array<String>, opzionale): Lista nomi stanze. Se array vuoto o non fornito, interroga tutte le stanze.
    -   `device_types` (Array<String>, opzionale): Lista tipi di dispositivo. Opzioni come `device_query`. Se array vuoto o non fornito, interroga tutti i tipi.
-   **Ritorno**: (String) Informazioni stato dispositivi in formato Markdown.

### device_log_query

-   **Descrizione**: Interroga i log dei dispositivi.
-   **Parametri**:
    -   `endpoint_ids` (Array<Integer>, richiesto): Lista ID dispositivi per interrogare la cronologia.
    -   `start_datetime` (String, opzionale): Ora inizio interrogazione, formato `YYYY-MM-DD HH:MM:SS` (ad es.: `"2023-05-16 12:00:00"`).
    -   `end_datetime` (String, opzionale): Ora fine interrogazione, formato `YYYY-MM-DD HH:MM:SS`.
    -   `attribute` (String, opzionale): Nome proprietà specifica dispositivo da interrogare (ad es.: `on_off`, `brightness`). Se non fornito, interroga la cronologia di tutte le proprietà registrate per quel dispositivo.
-   **Ritorno**: (String) Informazioni stato storico dispositivi in formato Markdown. (Nota: L'implementazione attuale potrebbe mostrare "This feature will be available soon.", indicando che la funzionalità è in attesa di completamento.)

### run_scenes

-   **Descrizione**: Esegui scene specificate basandosi sugli ID scene.
-   **Parametri**:
    -   `scenes` (Array<Integer>, richiesto): Lista ID scene da eseguire.
-   **Ritorno**: (String) Messaggio del risultato dell'esecuzione scene.

### get_scenes

-   **Descrizione**: Interroga tutte le scene sotto la casa dell'utente, o scene in stanze specificate.
-   **Parametri**:
    -   `positions` (Array<String>, opzionale): Lista nomi stanze. Se array vuoto o non fornito, interroga le scene dell'intera casa.
-   **Ritorno**: (String) Informazioni scene in formato Markdown.

### automation_config

-   **Descrizione**: Configura attività di controllo dispositivi temporizzate o ritardate.
-   **Parametri**:
    -   `scheduled_time` (String, richiesto): Ora impostata (se attività ritardata, convertita basandosi sull'ora corrente), formato `YYYY-MM-DD HH:MM:SS` (ad es.: `"2025-05-16 12:12:12"`).
    -   `endpoint_ids` (Array<Integer>, richiesto): Lista ID dispositivi da controllare temporizzatamente.
    -   `control_params` (Object, richiesto): Parametri di controllo dispositivo, usa lo stesso formato dello strumento `device_control` (include action, attribute, value, ecc.).
-   **Ritorno**: (String) Messaggio risultato configurazione automazione.

### get_homes

-   **Descrizione**: Ottieni tutti gli elenchi case sotto l'account utente.
-   **Parametri**: Nessuno.
-   **Ritorno**: (String) Elenco nomi case separati da virgole. Se non ci sono dati, restituisce stringa vuota o messaggio informativo corrispondente.

### switch_home

-   **Descrizione**: Cambia la casa attualmente operata dall'utente. Dopo il cambio, le successive query dispositivi, controlli, ecc. saranno mirati alla nuova casa cambiata.
-   **Parametri**:
    -   `home_name` (String, richiesto): Nome della casa target (dovrebbe provenire dall'elenco disponibile fornito dallo strumento `get_homes`).
-   **Ritorno**: (String) Messaggio del risultato dell'operazione di cambio.

## Struttura del progetto

```
.
├── cmd.go                # Definizione comandi Cobra CLI e punto di ingresso programma (contiene funzione main)
├── server.go             # Logica core server MCP, definizione strumenti e gestione richieste
├── smh.go                # Incapsulamento interfaccia API piattaforma smart home Aqara
├── middleware.go         # Middleware: autenticazione utente, controllo timeout, recupero eccezioni
├── config.go             # Gestione configurazione globale e elaborazione variabili d'ambiente
├── go.mod                # File gestione dipendenze modulo Go
├── go.sum                # File checksum dipendenze modulo Go
├── img/                  # Risorse immagini utilizzate nella documentazione README
├── LICENSE               # Licenza open source MIT
└── README.md             # Documentazione progetto
```

### Spiegazione file core

-   **`cmd.go`**: Implementazione CLI basata su framework Cobra, definisce modalità di avvio `run stdio` e `run sse` e funzione di ingresso principale
-   **`server.go`**: Implementazione core server MCP, responsabile per registrazione strumenti, gestione richieste e supporto protocollo
-   **`smh.go`**: Livello di incapsulamento API piattaforma smart home Aqara, fornisce controllo dispositivi, autenticazione e supporto multi-casa
-   **`middleware.go`**: Middleware elaborazione richieste, fornisce validazione autenticazione, controllo timeout e gestione eccezioni
-   **`config.go`**: Gestione configurazione globale, responsabile per elaborazione variabili d'ambiente e configurazione API

## Guida ai contributi

Benvenuti a partecipare al progetto contribuendo tramite l'invio di Issue o Pull Request!

Prima di inviare codice, assicurati che:
1.  Il codice segua gli standard di codifica del linguaggio Go.
2.  Le definizioni di strumenti MCP correlati e interfacce prompt mantengano coerenza e chiarezza.
3.  Aggiungi o aggiorna test unitari per coprire le tue modifiche.
4.  Se necessario, aggiorna la documentazione correlata (come questo README).
5.  Assicurati che i tuoi messaggi di commit siano chiari e comprensibili.

## Licenza

Questo progetto è concesso in licenza sotto [MIT License](/LICENSE).
Copyright (c) 2025 Aqara-Copliot