# 🧵 GoForum

GoForum est une application web de forum développée en **Go (Golang)** avec une base de données **SQLite** et un frontend en **HTML/CSS/JavaScript**.  
Ce projet a été réalisé dans le cadre de notre première année de **B1 Info/Cyber à Ynov**.

---

## 🎯 Objectifs du projet

- Créer un site fonctionnel permettant :
    - ✅ L’inscription et la connexion des utilisateurs
    - 📝 La création, modification et suppression de posts
    - 💬 L’ajout de commentaires
    - 🖼️ L’upload sécurisé d’images
    - 📧 L’envoi d’e-mails (vérification, mot de passe oublié)
    - 🔐 La connexion via OAuth (Google, Discord)

- Aller plus loin que le cahier des charges initial en ajoutant :
    - 🧹 Des suppressions automatiques de médias inutilisés ou de liens expirés
    - 📄 Un système de pagination
    - 🪵 Des logs configurables
    - 🌐 Un support multilingue

---

## 🔐 Sécurité

GoForum intègre plusieurs mécanismes pour garantir la sécurité des utilisateurs :

- 🔑 **Hash des mots de passe** : les mots de passe utilisateurs sont **hachés avec une méthode sécurisée** (bcrypt) avant d’être stockés en base de données.
- 🧱 **Prévention des injections SQL** : les requêtes sont faites via des **requêtes préparées (prepared statements)** pour éviter toute injection.
- 🔍 **Contrôle d’accès API** : chaque appel AJAX est **doublement vérifié côté serveur** pour s’assurer que l’utilisateur dispose bien des droits nécessaires.

---

## ⚙️ Installation et configuration

### 🧰 Pré-requis

- Go installé (`go 1.20+`)
- CGO activé pour le support SQLite (`CGO_ENABLED=1`)
- Créer un fichier `.env` (voir plus bas)

### 🏗️ Compilation

```bash
CGO_ENABLED=1 go build -o goforum main.go
```

### 🚀 Lancement

```bash
  ./goforum
```

Lancement avec en affichant les logs de débogage :

```bash
  ./goforum -d
```

Lancement avec l'enregistrement des logs :

```bash
  ./goforum -l
```

### 🧾 Arguments CLI disponibles

| Argument        | Description                                  |
| --------------- | -------------------------------------------- |
| `-d` / `-debug` | Affiche les messages de debug                |
| `-l` / `-log`   | Active l’écriture des logs dans des fichiers |

### 🌳 Arborescence du projet

```bash
GoForum/
├── .env                      # À fournir : fichier de configuration
├── cert.pem                  # À fournir si HTTPS : certificat SSL
├── key.pem                   # À fournir si HTTPS : clé privée SSL
├── goForumDataBase.db        # Généré automatiquement au lancement
│
├── main.go                   # Point d'entrée principal
├── go.mod                    # Modules Go
├── LICENSE
├── README.md
├── .gitignore
│
├── backend/
│   ├── apiPageHandlers/      # Handlers d’API AJAX
│   ├── emailsHandlers/       # Gestion des emails
│   └── pagesHandlers/        # Handlers de pages HTML
│       └── launcher.go
│
├── functions/                # Fonctions utilitaires
│
├── statics/
│   ├── css/                  # Feuilles de style
│   ├── fonts/                # Polices
│   ├── img/                  # Images
│   ├── js/                   # Scripts JS
│   └── lang/                 # Fichiers de traduction
│
├── templates/                # Templates HTML Go
│   ├── base.html             # Template principal
│   ├── *.html                # Autres pages : login, home, register, etc.
│
├── uploads/                  # Contient les fichiers uploadés
│
├── docker/                   # ❌ Dossier non fonctionnel – à ignorer
```

>Les fichiers .env, cert.pem, key.pem, et goForumDataBase.db ne sont pas inclus dans le dépôt :
>- .env doit être créé manuellement
>- cert.pem et key.pem sont requis uniquement si vous utilisez HTTPS
>- goForumDataBase.db est généré au lancement de l’application

### 🔧 Variables d’environnement

Liste complète des variables disponibles pour configurer GoForum via un fichier `.env`.

| Variable                                         | Type         | Description                                                           | Requis      |
| ------------------------------------------------ | ------------ | --------------------------------------------------------------------- | ----------- |
| `SESSION_SECRET`                                 | `string`     | Clé secrète pour chiffrer les cookies de session                      | ✅           |
| `DB_NAME`                                        | `string`     | Nom du fichier SQLite (ex: `goForumDataBase.db`)                      | ✅           |
| `CGO_ENABLED`                                    | `int`        | Doit être à `1` pour compiler avec SQLite                             | ✅           |
| `PORT`                                           | `int`        | Port d'écoute HTTP/HTTPS (ex: 80 ou 443)                              | ✅           |
| `CERT_FILE`                                      | `string`     | Chemin du certificat SSL (`cert.pem`)                                 | ⚠️ Si HTTPS |
| `CERT_KEY_FILE`                                  | `string`     | Clé privée SSL (`key.pem`)                                            | ⚠️ Si HTTPS |
| `DEFAULT_LANG`                                   | `string`     | Langue par défaut (`en` ou `fr`)                                      | ❌           |
| `LOG_FILE_CHANGE_TIME`                           | `int`        | Fréquence (en minutes) de rotation des fichiers logs                  | ❌           |
| `UPLOAD_FOLDER`                                  | `string`     | Dossier principal de stockage des fichiers uploadés                   | ❌           |
| `IMG_UPLOAD_FOLDER`                              | `string`     | Sous-dossier dans `UPLOAD_FOLDER` pour les images                     | ❌           |
| `AUTO_DELETE_OLD_EMAIL_IDENTIFICATIONS`          | `bool`       | Supprimer automatiquement les anciens liens email (`true` ou `false`) | ❌           |
| `AUTO_DELETE_OLD_EMAIL_IDENTIFICATIONS_INTERVAL` | `int`        | Intervalle entre les suppressions (minutes)                           | ❌           |
| `EMAIL_IDENTIFICATIONS_MAX_AGE`                  | `int`        | Âge maximum (minutes) d’un lien email avant suppression               | ❌           |
| `AUTO_DELETE_USELESS_MEDIA_LINKS`                | `bool`       | Supprimer les images inutilisées (`true` ou `false`)                  | ❌           |
| `AUTO_DELETE_USELESS_MEDIA_LINKS_INTERVAL`       | `int`        | Fréquence de suppression d’images inutilisées (minutes)               | ❌           |
| `MAX_MESSAGES_PER_PAGE_LOAD`                     | `int`        | Nombre de messages chargés par page via API                           | ❌           |
| `MAX_COMMENTS_PER_PAGE_LOAD`                     | `int`        | Nombre de commentaires chargés par page via API                       | ❌           |
| `SMTP_HOST`, `SMTP_PORT`                         | `string/int` | Configuration SMTP pour l'envoi des emails                            | ❌           |
| `SMTP_USER`, `SMTP_PASSWORD`                     | `string`     | Identifiants SMTP                                                     | ❌           |
| `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`       | `string`     | Identifiants OAuth pour connexion Google                              | ✅ si OAuth  |
| `DISCORD_CLIENT_ID`, `DISCORD_CLIENT_SECRET`     | `string`     | Identifiants OAuth pour connexion Discord                             | ✅ si OAuth  |


###🧪 Exemple de fichier `.env` minimal

```dotenv
# Sécurité
SESSION_SECRET=mySuperSecretSessionKey123!

# Base de données
DB_NAME=goForumDataBase.db
CGO_ENABLED=1

# Réseau
PORT=80

# Langue
DEFAULT_LANG=fr
```

### 👥 Crédits

Projet développé par les étudiants de B1 Info/Cyber – Ynov Campus (année 2024–2025) :
- Maxime Cordonnier
- Timothé Clement
- Emmanuel Persaud Conde Reis
- Julien Riviere

### 📜 Licence

Ce projet est distribué sous la licence Creative Commons Attribution - NonCommercial - NoDerivatives 4.0 International.

Cela signifie :
- ✅ Vous pouvez le partager
- 🚫 Pas d’utilisation commerciale
- 🚫 Pas de modification sans autorisation

<a href="https://creativecommons.org/licenses/by-nc-nd/4.0/" target="_blank">🔗 En savoir plus sur la licence</a>

### 🛠️ Support
Pour toute question ou signalement de bug, merci d’ouvrir un ticket à l’adresse suivante :

👉 https://ytrack.learn.ynov.com/git/comaxime/GoForum/issues

>Veuillez ne pas nous contacter par email concernant ce projet.