# Documentation Développeur

## Installation

> [!IMPORTANT]
> Pré-requis :
>
> - Docker et DockerCompose
> - Git
> - Une machine Linux ou MacOs (windows n'est pas encore pris en charge)

```bash
    git clone https://github.com/Gouel/GouelMonoRepo.git
    cd ./GouelMonoRepo
    docker-compose up -d
```

Ces lignes vont cloner le repository de Gouel, puis installer et lancer les bases de données (Redis / MongoDB)

```bash
    cd ./GouelServer
    ./build.sh
    # Si besoin ajouter la permission d'execution
    chmod u+x GouelServer
```

Ceci va générer le binaire du serveur Gouel (sous le nom GouelServer)

Il faudra ensuite créer un fichier `.env` au même endroit que le fichier `GouelServer`. Un exemple de fichier est fourni [.example.env](https://github.com/Gouel/GouelMonoRepo/blob/main/GouelServer/.exemple.env)

Pour vous aidez à générer un JWT_SECRET_KEY fort vous avez un outil directement intégré dans `GouelServer`

> [!WARNING]
> Il est possible que l'outil refuse de fonctionner si les champs liés à MongoDB ne sont pas remplis dans le `.env`

```bash
    ./GouelServer --secret
```

Une fois le fichier `.env` complété. Il vous sera possible de instancier le serveur puis de le lancer.

```bash
    ./GouelServer --setup
    ./GouelServer
```

> [!NOTE]
> Il est possible d'importer et d'exporter la base de données.
> Dans le cas de l'importation, `--setup` n'est pas obligatoire

```bash
    ./GouelServer --export ./out.gz
    ./GouelServer --import ./in.gz
```

Et voilà, votre serveur Gouel est prêt.

## Utilisation

Une fois lancé, GouelServeur produit une API Rest sur le socket configuré (`.env`).
Il n'y a pas encore de mapping des routes de l'API. Cependant plusieurs ressources existent:

- PostMan : [`Docs/gouel.postman.json`](https://github.com/Gouel/GouelMonoRepo/blob/main/Docs/gouel.postman.md)
- Code source : [`GouelServer/pkg/router/routes.go`](https://github.com/Gouel/GouelMonoRepo/blob/main/GouelServer/pkg/router/routes.go)

> [!TIP]
> Vous pouvez aussi ajouter des données directement dans la base de donnée sans passer par le serveur.
> Pour cela il faut simplement vous connecter au Docker MongoDB. (Avec MongoCompass par exemple)
