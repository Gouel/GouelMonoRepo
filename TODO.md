# Liste des choses à faire/en cours

## Pour la version 1

- Produire une release de l'application sur Android et possiblement IOS
- Produire les documents légaux pour GouelFront (CGU / CGV / Mentions légales)
- Vérifier la sécurité de GouelFront / GouelServeur

## Général

- Ajouter gestion écocup

- Ajouter champ "SamBonus" dans event pour ajouter un certain montant aux utilisateurs SAM
- Limitation d'une Identité par email (demande confirmation email + avertissement) (gouelFront)
- Ajout champs "wherePurchased" dans les tickets avec une valeur ("onSite" / "online") + Ajouter prix différent OnSite / Online
- Documentation : fonction import et export avec `MongoDB Database Tools`

- Ajouter génération badges bénévoles (publipostage)
- Faire documentation
  - Tutoriels Vidéos / écrits
  - GouelServer : Documentation routes API
  - Vidéo de présentation de Gouel
  - Remplir le README.md
- Entrée :
  - Afficher tickets seulement en recherche (bottomsheet)
  - Récupérer par packet (pagination coté serveur)
  - Effectuer recherche ticket (coté serveur)
  - Implémenter scroll infini

- Ajouter fonctionnalité BlueTooth + imprimante thermique (Impression ticket physique)
- Ajouter panneaux statistiques dans GouelFront/admin
- Optimisations :
  - GouelApp : Utiliser RiverPod au lieu de Provider
  - GouelApp / GouelServer : Utiliser des websockets pour les rechargements (rechargement "instantanés")
- Généraliser GouelFront pour ajouter d'autres PaymentProviders (autre que HelloAsso)
- Utiliser des fonctions de logging plutôt que des prints
