<div align="center" style="display: flex; align-items: center; justify-content: center; ">

  <img src="/readme/img/logo.png" alt="Aqara Logo" height="120">
  <h1>MCP Server</h1>

</div>

<div align="center">

[English](/readme/README.md) | [中文](/readme/README_CN.md) | [繁體中文](/readme/README_CHT.md) | Français | [한국어](/readme/README_KR.md) | [Español](/readme/README_ES.md) | [日本語](/readme/README_JP.md) | [Deutsch](/readme/README_DE.md) | [Italiano](/readme/README_IT.md)

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/aqara/aqara-mcp-server)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![Release](https://img.shields.io/github/v/release/aqara/aqara-mcp-server)](https://github.com/aqara/aqara-mcp-server/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

Aqara MCP Server est un service de contrôle domotique développé basé sur le protocole [MCP (Model Context Protocol)](https://modelcontextprotocol.io/introduction). Il permet à tout assistant IA ou API prenant en charge le protocole MCP (comme Claude, Cursor, etc.) d'interagir avec vos appareils domotiques Aqara, permettant le contrôle des appareils, les requêtes d'état, l'exécution de scénarios et plus encore via le langage naturel.

## Table des matières

- [Table des matières](#table-des-matières)
- [Fonctionnalités](#fonctionnalités)
- [Principe de fonctionnement](#principe-de-fonctionnement)
- [Démarrage rapide](#démarrage-rapide)
  - [Prérequis](#prérequis)
  - [Installation](#installation)
    - [Option 1 : Télécharger la version précompilée (Recommandé)](#option-1--télécharger-la-version-précompilée-recommandé)
    - [Option 2 : Compiler depuis les sources](#option-2--compiler-depuis-les-sources)
  - [Authentification du compte Aqara](#authentification-du-compte-aqara)
  - [Configuration du client](#configuration-du-client)
    - [Exemple de configuration Claude for Desktop](#exemple-de-configuration-claude-for-desktop)
    - [Description des paramètres de configuration](#description-des-paramètres-de-configuration)
    - [Autres clients MCP](#autres-clients-mcp)
  - [Démarrage du service](#démarrage-du-service)
    - [Mode standard (Recommandé)](#mode-standard-recommandé)
    - [Mode HTTP (Optionnel)](#mode-http-optionnel)
- [Documentation des outils API](#documentation-des-outils-api)
  - [Contrôle des appareils](#contrôle-des-appareils)
    - [device\_control](#device_control)
  - [Requête d'appareils](#requête-dappareils)
    - [device\_query](#device_query)
    - [device\_status\_query](#device_status_query)
    - [device\_log\_query](#device_log_query)
  - [Gestion des scénarios](#gestion-des-scénarios)
    - [get\_scenes](#get_scenes)
    - [run\_scenes](#run_scenes)
  - [Gestion des foyers](#gestion-des-foyers)
    - [get\_homes](#get_homes)
    - [switch\_home](#switch_home)
  - [Configuration d'automatisation](#configuration-dautomatisation)
    - [automation\_config](#automation_config)
- [Structure du projet](#structure-du-projet)
  - [Structure des répertoires](#structure-des-répertoires)
  - [Description des fichiers principaux](#description-des-fichiers-principaux)
- [Guide de développement](#guide-de-développement)
- [Licence](#licence)

## Fonctionnalités

- **Contrôle complet des appareils** : Prend en charge le contrôle fin de divers attributs des appareils intelligents Aqara, y compris les interrupteurs, la luminosité, la température de couleur, les modes, etc.
- **Requêtes d'appareils flexibles** : Capacité à interroger les listes d'appareils et leurs états détaillés par pièce et type d'appareil
- **Gestion intelligente des scénarios** : Prise en charge de l'interrogation et de l'exécution des scénarios domotiques prédéfinis par l'utilisateur
- **Historique des appareils** : Interrogation des enregistrements de changements d'état historiques des appareils dans des plages de temps spécifiées
- **Configuration d'automatisation** : Prise en charge de la configuration de tâches de contrôle d'appareils programmées ou différées
- **Support multi-foyers** : Prise en charge de l'interrogation et du basculement entre différents foyers sous les comptes utilisateur
- **Compatibilité protocole MCP** : Entièrement conforme aux spécifications du protocole MCP, facile à intégrer avec divers assistants IA
- **Mécanisme d'authentification sécurisé** : Utilise une authentification sécurisée basée sur l'autorisation de connexion + signature pour protéger les données utilisateur et la sécurité des appareils
- **Fonctionnement multiplateforme** : Développé en langage Go, peut être compilé en fichiers exécutables pour plusieurs plateformes
- **Facile à étendre** : Conception modulaire permettant d'ajouter facilement de nouveaux outils et fonctionnalités

## Principe de fonctionnement

Aqara MCP Server sert de pont entre les assistants IA et la plateforme domotique Aqara :

1. **Assistant IA (Client MCP)** : Les utilisateurs émettent des commandes via les assistants IA (par exemple, "Allumer les lumières du salon")
2. **Client MCP** : Analyse les commandes utilisateur et appelle les outils correspondants fournis par Aqara MCP Server selon le protocole MCP (par exemple, `device_control`)
3. **Aqara MCP Server (ce projet)** : Reçoit les requêtes des clients, les valide et appelle le module `smh.go`
4. **Module `smh.go`** : Utilise les identifiants Aqara configurés pour communiquer avec les API cloud Aqara, exécutant les opérations d'appareils réelles ou les requêtes de données
5. **Flux de réponse** : Les API cloud Aqara retournent les résultats, qui sont transmis via Aqara MCP Server au client MCP et finalement présentés à l'utilisateur

## Démarrage rapide

### Prérequis

- Go (version 1.24 ou supérieure)
- Git (pour compiler depuis les sources)
- Compte Aqara avec appareils intelligents liés

### Installation

Vous pouvez choisir de télécharger des fichiers exécutables précompilés ou de compiler depuis les sources.

#### Option 1 : Télécharger la version précompilée (Recommandé)

Visitez la page GitHub Releases pour télécharger le dernier fichier exécutable pour votre système d'exploitation :

**📥 [Aller à la page Releases](https://github.com/aqara/aqara-mcp-server/releases)**

Téléchargez et extrayez le package approprié pour votre plateforme.

#### Option 2 : Compiler depuis les sources

```bash
# Cloner le dépôt
git clone https://github.com/aqara/aqara-mcp-server.git
cd aqara-mcp-server

# Télécharger les dépendances
go mod tidy

# Compiler l'exécutable
go build -o aqara-mcp-server
```

Après compilation, un exécutable `aqara-mcp-server` sera généré dans le répertoire courant.

### Authentification du compte Aqara

Pour permettre au MCP Server d'accéder à votre compte Aqara et de contrôler les appareils, vous devez d'abord compléter l'autorisation de connexion.

Veuillez visiter l'adresse suivante pour compléter l'autorisation de connexion :
**🔗 [https://cdn.aqara.com/app/mcpserver/login.html](https://cdn.aqara.com/app/mcpserver/login.html)**

Après une connexion réussie, vous obtiendrez les informations d'authentification nécessaires (comme `token`, `region`), qui seront utilisées dans les étapes de configuration ultérieures.

> ⚠️ **Rappel de sécurité** : Veuillez garder vos informations `token` sécurisées et ne les partagez pas avec d'autres.

### Configuration du client

Différents clients MCP ont des méthodes de configuration légèrement différentes. Voici un exemple de configuration de Claude for Desktop pour utiliser ce MCP Server :

#### Exemple de configuration Claude for Desktop

1. Ouvrir les paramètres de Claude for Desktop

    ![Claude Open Setting](/readme/img/opening_setting.png)

2. Basculer vers l'onglet Développeur, puis cliquer sur Modifier la configuration pour ouvrir le fichier de configuration avec un éditeur de texte

    ![Claude Edit Configuration](/readme/img/edit_config.png)

3. Ajouter les informations de configuration de la "Page de succès de connexion" au fichier de configuration du client `claude_desktop_config.json`

    ![Configuration Example](/readme/img/config_info.png)

#### Description des paramètres de configuration

- `command` : Chemin complet vers votre fichier exécutable `aqara-mcp-server` téléchargé ou compilé
- `args` : Utiliser `["run", "stdio"]` pour démarrer le mode de transport stdio
- `env` : Configuration des variables d'environnement
  - `token` : Jeton d'accès obtenu depuis la page de connexion Aqara
  - `region` : Votre région de compte Aqara (par exemple, CN, US, EU, etc.)

#### Autres clients MCP

Pour d'autres clients prenant en charge le protocole MCP (comme ChatGPT, Cursor, etc.), la configuration est similaire :

- S'assurer que le client prend en charge le protocole MCP
- Configurer le chemin du fichier exécutable et les paramètres de démarrage
- Définir les variables d'environnement `token` et `region`
- Choisir le protocole de transport approprié (stdio recommandé)

### Démarrage du service

#### Mode standard (Recommandé)

Redémarrer Claude for Desktop. Vous pouvez ensuite effectuer le contrôle des appareils, les requêtes d'appareils, l'exécution de scénarios et d'autres opérations via le langage naturel.

![Claude Chat Example](/readme/img/claude.png)

#### Mode HTTP (Optionnel)

Si vous devez utiliser le mode HTTP, vous pouvez le démarrer ainsi :

```bash
# Utiliser le port par défaut 8080
./aqara-mcp-server run http

# Ou spécifier un hôte et port personnalisés
./aqara-mcp-server run http --host localhost --port 9000
```

Ensuite, utilisez les paramètres `["run", "http"]` dans la configuration du client.

## Documentation des outils API

Les clients MCP peuvent interagir avec les appareils domotiques Aqara en appelant ces outils.

### Contrôle des appareils

#### device_control

Contrôler l'état ou les attributs des appareils domotiques (par exemple, interrupteurs, température, luminosité, couleur, température de couleur, etc.).

**Paramètres :**

- `endpoint_ids` _(Array\<Integer\>, requis)_ : Liste des IDs d'appareils à contrôler
- `control_params` _(Object, requis)_ : Objet de paramètres de contrôle contenant des opérations spécifiques :
  - `action` _(String, requis)_ : Opération à exécuter (par exemple, `"on"`, `"off"`, `"set"`, `"up"`, `"down"`, `"cooler"`, `"warmer"`)
  - `attribute` _(String, requis)_ : Attribut d'appareil à contrôler (par exemple, `"on_off"`, `"brightness"`, `"color_temperature"`, `"ac_mode"`)
  - `value` _(String | Number, optionnel)_ : Valeur cible (requis quand action est "set")
  - `unit` _(String, optionnel)_ : Unité de la valeur (par exemple, `"%"`, `"K"`, `"℃"`)

**Retourne :** Message de résultat d'opération pour le contrôle d'appareil

### Requête d'appareils

#### device_query

Obtenir la liste des appareils basée sur l'emplacement spécifié (pièce) et le type d'appareil (n'inclut pas les informations d'état en temps réel).

**Paramètres :**

- `positions` _(Array\<String\>, optionnel)_ : Liste des noms de pièces. Tableau vide signifie interroger toutes les pièces
- `device_types` _(Array\<String\>, optionnel)_ : Liste des types d'appareils (par exemple, `"Light"`, `"WindowCovering"`, `"AirConditioner"`, `"Button"`). Tableau vide signifie interroger tous les types

**Retourne :** Liste d'appareils au format Markdown incluant les noms et IDs d'appareils

#### device_status_query

Obtenir les informations d'état actuel des appareils (pour interroger les informations d'état en temps réel comme la couleur, la luminosité, les interrupteurs, etc.).

**Paramètres :**

- `positions` _(Array\<String\>, optionnel)_ : Liste des noms de pièces. Tableau vide signifie interroger toutes les pièces
- `device_types` _(Array\<String\>, optionnel)_ : Liste des types d'appareils. Mêmes options que `device_query`. Tableau vide signifie interroger tous les types

**Retourne :** Informations d'état d'appareil au format Markdown

#### device_log_query

Interroger les informations de journal historique des appareils.

**Paramètres :**

- `endpoint_ids` _(Array\<Integer\>, requis)_ : Liste des IDs d'appareils pour interroger l'historique
- `start_datetime` _(String, optionnel)_ : Heure de début de requête au format `YYYY-MM-DD HH:MM:SS` (par exemple, `"2023-05-16 12:00:00"`)
- `end_datetime` _(String, optionnel)_ : Heure de fin de requête au format `YYYY-MM-DD HH:MM:SS`
- `attribute` _(String, optionnel)_ : Nom d'attribut d'appareil spécifique à interroger (par exemple, `on_off`, `brightness`). Interroge tous les attributs enregistrés quand non fourni

**Retourne :** Informations d'état historique d'appareil au format Markdown

> 📝 **Note :** L'implémentation actuelle peut afficher "This feature will be available soon.", indiquant que la fonctionnalité est en attente de finalisation.

### Gestion des scénarios

#### get_scenes

Interroger tous les scénarios dans le foyer de l'utilisateur, ou les scénarios dans des pièces spécifiées.

**Paramètres :**

- `positions` _(Array\<String\>, optionnel)_ : Liste des noms de pièces. Tableau vide signifie interroger les scénarios pour tout le foyer

**Retourne :** Informations de scénario au format Markdown

#### run_scenes

Exécuter des scénarios spécifiés basés sur les IDs de scénario.

**Paramètres :**

- `scenes` _(Array\<Integer\>, requis)_ : Liste des IDs de scénarios à exécuter

**Retourne :** Message de résultat d'exécution de scénario

### Gestion des foyers

#### get_homes

Obtenir la liste de tous les foyers sous le compte utilisateur.

**Paramètres :** Aucun

**Retourne :** Liste des noms de foyers séparés par des virgules. Retourne une chaîne vide ou un message approprié si aucune donnée

#### switch_home

Basculer le foyer d'opération actuel de l'utilisateur. Après basculement, les requêtes d'appareils, contrôles et autres opérations ultérieures cibleront le foyer nouvellement basculé.

**Paramètres :**

- `home_name` _(String, requis)_ : Nom du foyer cible

**Retourne :** Message de résultat d'opération de basculement

### Configuration d'automatisation

#### automation_config

Configurer des tâches de contrôle d'appareils programmées ou différées (prend actuellement en charge seulement la configuration d'automatisation de délai temporisé).

**Paramètres :**

- `scheduled_time` _(String, requis)_ : Point de temps défini (si tâche de délai, converti basé sur le point de temps actuel), format `YYYY-MM-DD HH:MM:SS` (par exemple, `"2025-05-16 12:12:12"`)
- `endpoint_ids` _(Array\<Integer\>, requis)_ : Liste des IDs d'appareils pour contrôle programmé
- `control_params` _(Object, requis)_ : Paramètres de contrôle d'appareil utilisant le même format que l'outil `device_control` (incluant action, attribute, value, etc.)

**Retourne :** Message de résultat de configuration d'automatisation

> 📝 **Note :** L'implémentation actuelle peut afficher "This feature will be available soon.", indiquant que la fonctionnalité est en attente de finalisation.

## Structure du projet

### Structure des répertoires

```text
.
├── cmd.go                # Définition des commandes Cobra CLI et point d'entrée du programme (contient la fonction main)
├── server.go             # Logique principale du serveur MCP, définition des outils et gestion des requêtes
├── smh.go                # Wrapper d'interface API de la plateforme domotique Aqara
├── middleware.go         # Middleware : authentification utilisateur, contrôle de timeout, récupération d'exception
├── config.go             # Gestion de configuration globale et traitement des variables d'environnement
├── go.mod                # Fichier de gestion des dépendances du module Go
├── go.sum                # Fichier de somme de contrôle des dépendances du module Go
├── readme/               # Documentation README et ressources d'images
│   ├── img/              # Répertoire des ressources d'images
│   └── *.md              # Fichiers README multilingues
├── LICENSE               # Licence open source MIT
└── README.md             # Documentation principale du projet
```

### Description des fichiers principaux

- **`cmd.go`** : Implémentation CLI basée sur le framework Cobra, définissant les modes de démarrage `run stdio` et `run http` et la fonction d'entrée principale
- **`server.go`** : Implémentation principale du serveur MCP, responsable de l'enregistrement des outils, de la gestion des requêtes et du support de protocole
- **`smh.go`** : Couche de wrapper d'API de la plateforme domotique Aqara, fournissant le contrôle d'appareils, l'authentification et le support multi-foyers
- **`middleware.go`** : Middleware de traitement des requêtes, fournissant la vérification d'authentification, le contrôle de timeout et la gestion d'exceptions
- **`config.go`** : Gestion de configuration globale, responsable du traitement des variables d'environnement et de la configuration API

## Guide de développement

Les contributions sont bienvenues via la soumission d'Issues ou de Pull Requests !

Avant de soumettre du code, veuillez vous assurer que :

1. Le code suit les standards de codage du langage Go
2. Les outils MCP et définitions d'interface associés maintiennent la cohérence et la clarté
3. Ajouter ou mettre à jour les tests unitaires pour couvrir vos changements
4. Mettre à jour la documentation pertinente (comme ce README) si nécessaire
5. S'assurer que vos messages de commit sont clairs et descriptifs

## Licence

Ce projet est sous licence [MIT License](/LICENSE).

Copyright (c) 2025 Aqara-Copilot
