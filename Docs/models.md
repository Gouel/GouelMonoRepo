# Modèles de données <!-- omit from toc -->

Cette page documente toute la terminologie et les points clés de la gestion des données de Gouel

## Sommaire <!-- omit from toc -->

## Les bases

Gouel utilise trois grandes collections de données :

- Users : Les données des utilisateurs de gouel (admin, api, bénévoles, utilisateurs lambda)
- Events : Les données des événements (Produits, Vestiaires, ...)
- Tickets : Les tickets générés (identifie un utilisateur, sur un événement, avec un type de ticket)

Dans ses collections on retrouvera des sous-collections :

- Users :
  - Transactions
- Events :
  - Tickets (Les tickets de l'événement)
  - Lockers (Les vestiaires)
  - Products (Les produits mis en vente)
- Tickets

Ce document va détailler chaque modèle de donnée.

## Collection Users

La collection Users gère l'ensemble des utilisateurs. Elle regroupe l'identité des personnes, leur solde, leurs diverses transactions, et une partie de leurs droits.

| nom          | type de donnée                       | commentaire                                                       |
| ------------ | ------------------------------------ | ----------------------------------------------------------------- |
| dob          | Texte                                | Date d'anniversaire au format iso (yyyy-mm-dd)                    |
| email        | Texte                                | L'email de l'utilisateur (Fait office d'identifiant)              |
| firstName    | Texte                                | Le prénom de l'utilisateur                                        |
| lastName     | Texte                                | Le nom de l'utilisateur                                           |
| password     | Texte (hashé)                        | Le mot de passe de l'utilisateur (hashé par le serveur)           |
| role         | Texte                                | Le role de l'utilisateur (SUPERADMIN / API / USER )               |
| solde        | Nombre                               | Le montant crédité sur le solde de l'utilisateur (cashless)       |
| transactions | Liste de [Transaction](#transaction) | La liste des transactions effectués sur le solde de l'utilisateur |

La base de donnée ajoute un champs supplémentaire : `id`. Ce champs est utilisé pour identifier fiablement l'utilisateur. (Il est utilisé par Gouel à différents endroits)

### Transaction

Une transaction représente une action sur le solde de l'utilisateur. Que ce soit un ajout de fond (credit) ou un débit (debit).

| nom          | type de donnée        | commentaire                                                          |
|--------------|-----------------------|----------------------------------------------------------------------|
| amount       | Nombre                | Le montant total de la transaction                                   |
| cart         | Liste des [achats](#purchasedproduct) | La liste des achats (produits + nombre)                              |
| date         | Texte                 | La date de transaction au format iso (yyyy-mm-ddThh:mm:ss)           |
| eventId      | Texte                 | L'identifiant de l'événement concerné par la transaction             |
| paymentType | Texte                 | Le type de paiement utilisé (espece, carte, helloasso, cashless,...) |
| type         | Texte                 | Le type d'action effectué (credit / debit)                           |

### PurchasedProduct

| nom         | type de donnée | commentaire                     |
|-------------|----------------|---------------------------------|
| amount      | Nombre         | Nombre d'unité acheté           |
| productCode | Texte          | l'identifiant du produit acheté |

## Collection Events

La collection Events gère l'ensemble des données liés aux événements. Elle regroupe les informations de l'événements, les tickets à vendre, les bénévoles, les administrateurs et les vestiaires.

| nom         | type de donnée                  | commentaire                                                      |
|-------------|---------------------------------|------------------------------------------------------------------|
| contact     | Texte                           | Une façon de contacter les organisateurs (ex : email, telephone) |
| description | Texte                           | Description de l'événement                                       |
| eventTicket | Liste des [ticket d'événement](#eventticket) | La liste des tickets à vendre pour l'événement                   |
| location    | Texte                           | L'endroit où se passe l'événement                                |
| lockers     | Liste de [vestiaires](#lockers)         | La liste des vestiaires disponibles lors de l'événement          |
| products    | Liste de [produits](#products)           | La liste des produits à vendre dans la buvette                   |
| isPublic      | Booléen                         | Si l'événement est visible sur le site "GouelFront"              |
| title       | Texte                           | Le titre de l'événement                                          |
| volunteers  | Liste de [bénévoles](#volunteers)          | La liste des bénévoles avec leurs permissions                    |

La base de donnée ajoute un champs supplémentaire : `id`. Ce champs est utilisé pour identifier fiablement l'événement. (Il est utilisé par Gouel à différents endroits).

### Volunteers

Information : Tout administrateur d'événement est considéré comme bénévole.

| nom         | type de donnée | commentaire                                                          |
|-------------|----------------|----------------------------------------------------------------------|
| isAdmin     | Booléen        | Si le bénévole est Administrateur ( à virtuellement tous les droits) |
| permissions | Liste de Texte | La liste des droits du bénévole (ex : [buvette, caisse])             |
| userId      | Texte          | l'identifiant unique de l'utilisateur (bénévole)                     |

### Products

| nom         | type de donnée | commentaire                                                                                           |
|-------------|----------------|-------------------------------------------------------------------------------------------------------|
| amount      | Nombre*        | Le nombre maximum pouvant être vendu. null si quantité infinie                                        |
| endOfSale   | Texte          | Date au format iso spécifiant l'arrêt de la vente du produit                                          |
| hasAlcohol   | Booléen        | Si le produit contient de l'alcool                                                                    |
| icon        | Texte          | Icône du produit dans les interfaces. ([liste des icônes disponibles](#liste-des-icônes-disponibles)) |
| label       | Texte          | Nom du produit                                                                                        |
| price       | Nombre         | Prix du produit                                                                                       |
| productCode | Texte          | L'identifiant du produit (auto généré par le serveur)                                                 |
| purchased   | Nombre         | Le nombre de fois que le produit a été acheté au court de l'événement                                 |

#### Liste des icônes disponibles

WIP

### EventTicket

| nom             | type de donnée | commentaire                                              |
|-----------------|----------------|----------------------------------------------------------|
| amount          | Nombre*        | Nombre maximal de ticket à vendre, mettre null si infini |
| eventTicketCode | Texte          | Identifiant du ticket (auto généré par le serveur)       |
| price           | Nombre         | Prix du ticket                                           |
| purchased       | Nombre         | Nombre de ticket vendu                                   |
| title           | Texte          | Titre du ticket                                          |

### Lockers

| nom    | type de donnée | commentaire                                            |
|--------|----------------|--------------------------------------------------------|
| lockerCode   | Texte          | Identification du vestiaire (ex: A01)                  |
| userId | Texte          | l'identifiant du propriétaire du vestiaire. "" si vide |

## Tickets

| nom             | type de donnée | commentaire                                                                        |
|-----------------|----------------|------------------------------------------------------------------------------------|
| eventId         | Texte          | L'identifiant de l'événement concerné par le ticket                                |
| eventTicketCode | Texte          | L'identifiant du type de billet acheté                                             |
| isSam           | Booléen        | Si l'utilisateur s'est présenté comme SAM (empêche l'achat de produits alcoolisés) |
| isUsed          | Booléen        | Si le ticket a été utilisé                                                         |
| userId          | Text           | L'identifiant du propriétaire du ticket                                            |
| wasPurchased    | Booléen        | Si le ticket a été acheté ou non (différencie les bénévoles des utilisateurs)      |
