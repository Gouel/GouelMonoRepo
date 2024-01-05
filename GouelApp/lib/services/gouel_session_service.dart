class GouelSession {
  // Création d'une instance privée statique de la classe (singleton)
  static final GouelSession _instance = GouelSession._internal();

  // Constructeur nommé privé
  GouelSession._internal();

  // Factory pour accéder à l'instance
  factory GouelSession() {
    return _instance;
  }

  // Stockage des données de session
  final Map<String, dynamic> _storage = {};

  // Méthodes pour gérer les données de session
  void store(String key, dynamic value) {
    _storage[key] = value;
  }

  dynamic retrieve(String key) {
    return _storage.containsKey(key) ? _storage[key] : null;
  }

  void remove(String key) {
    _storage.remove(key);
  }

  @override
  String toString() {
    return _storage.toString();
  }
}
