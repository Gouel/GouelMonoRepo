# Documentation Développeur

## Installation

> [!IMPORTANT]
> Pré-requis :
>
> - Docker et DockerCompose
> - Git
> - Une machine Linux ou MacOs (windows n'est pas encore pris en charge)

### Récupération du projet à partir de git

```bash
    git clone https://github.com/Gouel/GouelMonoRepo.git
    cd ./GouelMonoRepo
```

### Création du certificat HTTPS

> [!WARNING]
> Il faut compléter le fichier `./release/env/nginx/openssl.cnf` avant de lancer le script bash
> De plus il faut aussi changer la configuration de nginx [default.conf](https://github.com/Gouel/GouelMonoRepo/blob/main/release/env/nginx/default.conf) pour utiliser le même domaine que votre `openssl.cnf`

```bash
cd ./release/env/nginx
bash gen_cert.sh
```

### Mise en places des environnements

Il y a deux fichiers `.env` à créer / compléter :

- Pour le `front` : [exemple](https://github.com/Gouel/GouelMonoRepo/blob/main/GouelFront/.exemple.env)
  - Note : Pour pouvoir utiliser l'API d'HelloAsso, il faut forcément que le champs `SERVER_NAME` commence par `https://`
- Pour le `back` : [exemple](https://github.com/Gouel/GouelMonoRepo/blob/main/GouelServer/.exemple.env)
  - Note: GouelFront à besoin d'un serveur mail (SMTP) pour envoyer les différents messages. Gouel ne fournit pas de configuration. C'est à vous de gérer.

Le plus simple est de remplacer les `exemple.env` en `.env` dans leurs dossiers respectifs

### Lancement de la solution

Vous pouvez maintenant lancer le docker compose

```bash
docker-compose -p "release_gouel" up -d
```

> [!WARNING]
> Au premier lancement, il faudra générer les comptes SUPERADMIN et API
> Pour cela il faut lancer la commande suivante
> `docker exec -ti <NOM_CONTAINER_GOUEL_SERVER|release_gouel-gouel-server-1> go run main.go --setup`

Et voilà, votre serveur Gouel est prêt.

## Utilisation

Une fois le docker-compose lancé, plusieurs pages seront disponibles :

- `https://www.<domain>` : GouelFront
- `https://app.<domain>` : GouelApp (version webapp)
- `https://server.<domain>` : GouelServer (Api REST)
- Sinon il est possible d'accéder aux différents services avec `<ip pc>:<port service>`. La liste des ports de chaque service est disponible dans le [docker-compose.yml](https://github.com/Gouel/GouelMonoRepo/blob/main/docker-compose.yml)

> [!IMPORTANT]
> L'application android Gouel ne supporte pas les certificats auto-signés. Il faudra donc rentrer `http://<ip pc>:5002` dans les paramètres de l'application

Si vous avez un problème, faites une issue et j'essayerai de le résoudre / de vous aider.

> [!TIP]
> Vous pouvez aussi ajouter des données directement dans la base de donnée sans passer par le serveur.
> Pour cela il faut simplement vous connecter au Docker MongoDB. (Avec MongoCompass par exemple)
