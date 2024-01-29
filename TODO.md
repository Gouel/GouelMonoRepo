# Liste des choses à faire/en cours

## Pour la version 1

- Produire une release de l'application sur Android et possiblement IOS
- Produire les documents légaux pour GouelFront (CGU / CGV / Mentions légales)
- Vérifier la sécurité de GouelFront / GouelServeur

## Général

- Utiliser des fonctions de logging plutôt que des prints
- Limitation d'une Identité par email
- Ajout champs "wherePurchased" dans les tickets avec une valeur ("onSite" / "online") + Ajouter prix différent OnSite / Online
- Ajouter champ "SamBonus" dans event pour ajouter un certain montant aux utilisateurs SAM
- Ajouter fonctionnalité BlueTooth + imprimante thermique (Impression ticket physique)
- Ajouter panneaux statistiques dans GouelFront/admin
- Ajouter gestion écocup
- Ajouter génération badges bénévoles (publipostage)
- Optimisations :
  - GouelApp : Utiliser RiverPod au lieu de Provider
  - GouelApp / GouelServer : Utiliser des websockets pour les rechargements (rechargement "instantanés")
- Généraliser GouelFront pour ajouter d'autres PayementProviders (autre que HelloAsso)
- Faire documentation
  - Tutoriels Vidéos / écrits
  - GouelServer : Documentation routes API
  - Vidéo de présentation de Gouel
  - Remplir le README.md
