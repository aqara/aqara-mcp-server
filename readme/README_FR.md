# Serveur Aqara MCP

[English](/readme/README.md) | [中文](/readme/README_CN.md) | [繁體中文](/readme/README_CHT.md) | Français | [한국어](/readme/README_KR.md) | [Español](/readme/README_ES.md) | [日本語](/readme/README_JP.md) | [Deutsch](/readme/README_DE.md) | [Italiano](/readme/README_IT.md)

[![Statut de Build](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/aqara/aqara-mcp-server)
[![Version Go](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![Release](https://img.shields.io/github/v/release/aqara/aqara-mcp-server)](https://github.com/aqara/aqara-mcp-server/releases)
[![Licence: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Aqara MCP Server est un service de contrôle domotique développé basé sur le protocole [MCP (Model Context Protocol)](https://modelcontextprotocol.io/introduction). Il permet à tout assistant IA ou API supportant le protocole MCP (comme Claude, ChatGPT, Cursor, etc.) d'interagir avec vos appareils domotiques Aqara, permettant le contrôle d'appareils, les requêtes de statut, l'exécution de scénarios et plus encore via le langage naturel.

## Table des Matières

- [Serveur Aqara MCP](#serveur-aqara-mcp)
  - [Table des Matières](#table-des-matières)
  - [Fonctionnalités](#fonctionnalités)
  - [Comment ça Fonctionne](#comment-ça-fonctionne)
  - [Démarrage Rapide](#démarrage-rapide)
    - [Prérequis](#prérequis)
    - [Installation](#installation)
    - [Authentification Compte Aqara](#authentification-compte-aqara)
    - [Exemple de Configuration (Claude for Desktop)](#exemple-de-configuration-claude-for-desktop)
    - [Lancement du Service](#lancement-du-service)
  - [Outils Disponibles](#outils-disponibles)
    - [device\_control](#device_control)
    - [device\_query](#device_query)
    - [device\_status\_query](#device_status_query)
    - [device\_log\_query](#device_log_query)
    - [run\_scenes](#run_scenes)
    - [get\_scenes](#get_scenes)
    - [automation\_config](#automation_config)
    - [get\_homes](#get_homes)
    - [switch\_home](#switch_home)
  - [Structure du Projet](#structure-du-projet)
    - [Descriptions des Fichiers Principaux](#descriptions-des-fichiers-principaux)
  - [Contribution](#contribution)
  - [Licence](#licence)

## Fonctionnalités

- **Contrôle Complet des Appareils** : Support pour le contrôle fin des appareils intelligents Aqara incluant interrupteurs, luminosité, température de couleur, modes, et plus.
- **Requêtes Flexibles d'Appareils** : Interroger les listes d'appareils et les statuts détaillés par pièce et type d'appareil.
- **Gestion Intelligente des Scénarios** : Support pour interroger et exécuter des scénarios domotiques prédéfinis par l'utilisateur.
- **Historique des Appareils** : Interroger les changements d'état historiques des appareils dans des plages de temps spécifiées.
- **Configuration d'Automatisation** : Support pour configurer des tâches de contrôle d'appareils programmées ou différées.
- **Support Multi-Maisons** : Support pour interroger et basculer entre différentes maisons sous les comptes utilisateur.
- **Compatible Protocole MCP** : Entièrement conforme aux spécifications du protocole MCP, facile à intégrer avec divers assistants IA.
- **Mécanisme d'Authentification Sécurisé** : Utilise une authentification sécurisée basée sur autorisation de connexion + signature pour protéger les données utilisateur et la sécurité des appareils.
- **Multi-plateforme** : Construit avec Go, peut être compilé en exécutables pour multiples plateformes.
- **Facile à Étendre** : Conception modulaire permettant l'ajout facile de nouveaux outils et fonctionnalités.

## Comment ça Fonctionne

Aqara MCP Server agit comme un pont entre les assistants IA et la plateforme domotique Aqara :

1. **Assistant IA (Client MCP)** : Les utilisateurs émettent des commandes via les assistants IA (ex., "Allume les lumières du salon").
2. **Client MCP** : Analyse les commandes utilisateur et appelle les outils correspondants fournis par Aqara MCP Server selon le protocole MCP (ex., `device_control`).
3. **Aqara MCP Server (Ce Projet)** : Reçoit les requêtes des clients, les valide, et appelle le module `smh.go`.
4. **Module `smh.go`** : Utilise les identifiants Aqara configurés pour communiquer avec les APIs cloud Aqara pour les opérations d'appareils ou requêtes de données réelles.
5. **Flux de Réponse** : L'API cloud Aqara retourne les résultats, qui sont renvoyés via Aqara MCP Server au client MCP et finalement présentés à l'utilisateur.

## Démarrage Rapide

### Prérequis

- Go (version 1.24 ou supérieure)
- Git (pour compiler depuis les sources)
- Compte Aqara avec appareils intelligents connectés

### Installation

Vous pouvez choisir de télécharger des exécutables pré-compilés ou compiler depuis les sources.

**Option 1 : Télécharger Version Pré-compilée (Recommandé)**

Visitez le lien ci-dessous pour télécharger le dernier package exécutable pour votre système d'exploitation.

[Page des Releases](https://github.com/aqara/aqara-mcp-server/releases)

Extraire et utiliser directement.

**Option 2 : Compiler depuis les Sources**

```bash
# Cloner le dépôt
git clone https://github.com/aqara/aqara-mcp-server.git
cd aqara-mcp-server

# Télécharger les dépendances
go mod tidy

# Compiler l'exécutable
go build -o aqara-mcp-server
```

Après compilation, l'exécutable `aqara-mcp-server` sera généré dans le répertoire courant.

### Authentification Compte Aqara

Pour permettre au MCP Server d'accéder à votre compte Aqara et contrôler les appareils, vous devez d'abord compléter l'autorisation de connexion.

Veuillez visiter l'adresse suivante pour compléter l'autorisation de connexion :
[https://cdn.aqara.com/app/mcpserver/login.html](https://cdn.aqara.com/app/mcpserver/login.html)

Après une connexion réussie, vous obtiendrez les informations d'authentification nécessaires (comme `token`, `region`), qui seront utilisées dans les étapes de configuration ultérieures.

**Veuillez garder ces informations sécurisées, surtout le `token` - ne le partagez avec personne.**

### Exemple de Configuration (Claude for Desktop)

Différents clients MCP ont des méthodes de configuration légèrement différentes. Voici un exemple de comment configurer Claude for Desktop pour utiliser ce MCP Server :

1. Ouvrir les Paramètres de Claude for Desktop.
2. Basculer vers l'onglet Développeur.
3. Cliquer sur Éditer Config pour ouvrir le fichier de configuration avec un éditeur de texte.

   ![](/readme/img/setting0.png)
   ![](/readme/img/setting1.png)

4. Ajouter les informations de configuration de la "Page de Succès de Connexion" au fichier de configuration du client (claude_desktop_config.json). Exemple de configuration :

   ![](/readme/img/config.png)

**Notes de Configuration :**
- `command` : Chemin complet vers votre exécutable `aqara-mcp-server` téléchargé ou compilé
- `args` : Utiliser `["run", "stdio"]` pour démarrer le mode de transport stdio
- `env` : Configuration des variables d'environnement
  - `token` : Token d'accès obtenu depuis la page de connexion Aqara
  - `region` : Région de votre compte Aqara (ex., CN, US, EU, etc.)

### Lancement du Service

Redémarrer Claude for Desktop. Ensuite, vous pouvez utiliser les conversations pour appeler les outils fournis par le MCP Server pour le contrôle d'appareils, les requêtes d'appareils, et autres opérations.

![](/readme/img/claude.png)

**Configuration d'Autres Clients MCP**

Pour d'autres clients supportant le protocole MCP (comme Claude, ChatGPT, Cursor, etc.), la configuration est similaire :
- S'assurer que le client supporte le protocole MCP
- Configurer le chemin du fichier exécutable et les paramètres de démarrage
- Définir les variables d'environnement `token` et `region`
- Choisir le protocole de transport approprié (stdio recommandé)

**Mode SSE (Optionnel)**

Si vous devez utiliser le mode SSE (Server-Sent Events), vous pouvez le démarrer ainsi :

```bash
# Utiliser le port par défaut 8080
./aqara-mcp-server run sse

# Ou spécifier un hôte et port personnalisés
./aqara-mcp-server run sse --host localhost --port 9000
```

Ensuite utiliser les paramètres `["run", "sse"]` dans la configuration client.

## Outils Disponibles

Les clients MCP peuvent interagir avec les appareils domotiques Aqara en appelant ces outils.

### device_control

- **Description** : Contrôler le statut ou les propriétés des appareils domotiques (ex., marche/arrêt, température, luminosité, couleur, température de couleur, etc.).
- **Paramètres** :
  - `endpoint_ids` (Array<Integer>, requis) : Liste des IDs d'appareils à contrôler.
  - `control_params` (Object, requis) : Objet de paramètres de contrôle contenant les opérations spécifiques.
    - `action` (String, requis) : Action à exécuter. Exemples : `"on"`, `"off"`, `"set"`, `"up"`, `"down"`, `"cooler"`, `"warmer"`.
    - `attribute` (String, requis) : Attribut d'appareil à contrôler. Exemples : `"on_off"`, `"brightness"`, `"color_temperature"`, `"ac_mode"`.
    - `value` (String | Number, optionnel) : Valeur cible (requis quand action est "set").
    - `unit` (String, optionnel) : Unité de la valeur (ex., `"%"`, `"K"`, `"℃"`).
- **Retour** : (String) Message de résultat d'opération pour le contrôle d'appareil.

### device_query

- **Description** : Obtenir la liste d'appareils par emplacement spécifié (pièce) et type d'appareil (n'inclut pas les informations de statut en temps réel, liste seulement les appareils et leurs IDs).
- **Paramètres** :
  - `positions` (Array<String>, optionnel) : Liste des noms de pièces. Si tableau vide ou non fourni, interroge toutes les pièces.
  - `device_types` (Array<String>, optionnel) : Liste des types d'appareils. Exemples : `"Light"`, `"WindowCovering"`, `"AirConditioner"`, `"Button"`, etc. Si tableau vide ou non fourni, interroge tous les types.
- **Retour** : (String) Liste d'appareils formatée en Markdown incluant noms d'appareils et IDs.

### device_status_query

- **Description** : Obtenir les informations de statut actuel des appareils (pour interroger les attributs liés au statut comme couleur, luminosité, interrupteurs, etc.).
- **Paramètres** :
  - `positions` (Array<String>, optionnel) : Liste des noms de pièces. Si tableau vide ou non fourni, interroge toutes les pièces.
  - `device_types` (Array<String>, optionnel) : Liste des types d'appareils. Mêmes options que `device_query`. Si tableau vide ou non fourni, interroge tous les types.
- **Retour** : (String) Informations de statut d'appareils formatées en Markdown.

### device_log_query

- **Description** : Interroger les journaux d'appareils.
- **Paramètres** :
  - `endpoint_ids` (Array<Integer>, requis) : Liste des IDs d'appareils pour interroger l'historique.
  - `start_datetime` (String, optionnel) : Heure de début de requête au format `YYYY-MM-DD HH:MM:SS` (ex., `"2023-05-16 12:00:00"`).
  - `end_datetime` (String, optionnel) : Heure de fin de requête au format `YYYY-MM-DD HH:MM:SS`.
  - `attribute` (String, optionnel) : Nom d'attribut d'appareil spécifique à interroger (ex., `on_off`, `brightness`). Si non fourni, interroge tous les attributs enregistrés pour l'appareil.
- **Retour** : (String) Informations de statut historique d'appareils formatées en Markdown. (Note : L'implémentation actuelle peut afficher "This feature will be available soon.", indiquant que la fonctionnalité est en attente de finalisation.)

### run_scenes

- **Description** : Exécuter des scénarios spécifiés par ID de scénario.
- **Paramètres** :
  - `scenes` (Array<Integer>, requis) : Liste des IDs de scénarios à exécuter.
- **Retour** : (String) Message de résultat d'exécution de scénario.

### get_scenes

- **Description** : Interroger tous les scénarios dans la maison de l'utilisateur, ou les scénarios dans des pièces spécifiées.
- **Paramètres** :
  - `positions` (Array<String>, optionnel) : Liste des noms de pièces. Si tableau vide ou non fourni, interroge les scénarios pour toute la maison.
- **Retour** : (String) Informations de scénarios formatées en Markdown.

### automation_config

- **Description** : Configurer des tâches de contrôle d'appareils programmées ou différées.
- **Paramètres** :
  - `scheduled_time` (String, requis) : Point de temps défini (pour les tâches différées, converti basé sur le temps actuel), format `YYYY-MM-DD HH:MM:SS` (ex., `"2025-05-16 12:12:12"`).
  - `endpoint_ids` (Array<Integer>, requis) : Liste des IDs d'appareils pour contrôle programmé.
  - `control_params` (Object, requis) : Paramètres de contrôle d'appareil utilisant le même format que l'outil `device_control` (incluant action, attribute, value, etc.).
- **Retour** : (String) Message de résultat de configuration d'automatisation.

### get_homes

- **Description** : Obtenir toutes les listes de maisons sous le compte utilisateur.
- **Paramètres** : Aucun.
- **Retour** : (String) Liste de noms de maisons séparés par des virgules. Retourne chaîne vide ou message approprié si aucune donnée.

### switch_home

- **Description** : Basculer la maison d'opération actuelle de l'utilisateur. Après basculement, les requêtes d'appareils, contrôles, et autres opérations ultérieures cibleront la maison nouvellement basculée.
- **Paramètres** :
  - `home_name` (String, requis) : Nom de la maison cible (devrait provenir de la liste disponible fournie par l'outil `get_homes`).
- **Retour** : (String) Message de résultat d'opération de basculement.

## Structure du Projet

```
.
├── cmd.go                # Définitions de commandes CLI Cobra et point d'entrée du programme (contient la fonction main)
├── server.go             # Logique centrale du serveur MCP, définitions d'outils et gestion des requêtes
├── smh.go                # Wrapper d'interface API de plateforme domotique Aqara
├── middleware.go         # Middleware : authentification utilisateur, contrôle de timeout, récupération d'exception
├── config.go             # Gestion de configuration globale et traitement des variables d'environnement
├── go.mod                # Fichier de gestion des dépendances du module Go
├── go.sum                # Fichier de somme de contrôle des dépendances du module Go
├── img/                  # Ressources d'images utilisées dans la documentation README
├── LICENSE               # Licence open source MIT
└── README.md             # Documentation du projet
```

### Descriptions des Fichiers Principaux

- **`cmd.go`** : Implémentation CLI basée sur le framework Cobra, définissant les modes de démarrage `run stdio` et `run sse` et la fonction d'entrée principale
- **`server.go`** : Implémentation centrale du serveur MCP, responsable de l'enregistrement d'outils, de la gestion des requêtes, et du support de protocole
- **`smh.go`** : Couche wrapper d'API de plateforme domotique Aqara, fournissant contrôle d'appareils, authentification, et support multi-maisons
- **`middleware.go`** : Middleware de traitement des requêtes fournissant validation d'authentification, contrôle de timeout, et gestion d'exceptions
- **`config.go`** : Gestion de configuration globale, responsable du traitement des variables d'environnement et de la configuration API

## Contribution

Bienvenue à contribuer au projet en soumettant des Issues ou Pull Requests !

Avant de soumettre du code, veuillez vous assurer que :
1. Le code suit les standards de codage du langage Go.
2. Les définitions d'outils MCP et d'interfaces de prompt liées maintiennent cohérence et clarté.
3. Ajouter ou mettre à jour les tests unitaires pour couvrir vos changements.
4. Mettre à jour la documentation pertinente (comme ce README) si nécessaire.
5. S'assurer que vos messages de commit sont clairs et descriptifs.

## Licence

Ce projet est sous licence [MIT License](/LICENSE).
Copyright (c) 2025 Aqara-Copliot