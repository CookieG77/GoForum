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

Lancement avec le 