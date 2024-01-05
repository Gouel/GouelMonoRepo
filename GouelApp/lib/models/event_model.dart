class Event {
  final String id;
  final String title;

  Event({required this.id, required this.title});

  factory Event.fromJson(Map<String, dynamic> json) {
    return Event(
      id: json['ID'] as String,
      title: json['Title'] as String,
    );
  }
}

class Locker {
  final String code;
  final String user;

  Locker({required this.code, required this.user});

  factory Locker.fromJson(Map<String, dynamic> json) {
    return Locker(
      code: json['code'] as String,
      user: json['user'] as String,
    );
  }
}
