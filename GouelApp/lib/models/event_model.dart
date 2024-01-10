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
  final String lockerCode;
  final String userId;

  Locker({required this.lockerCode, required this.userId});

  factory Locker.fromJson(Map<String, dynamic> json) {
    return Locker(
      lockerCode: json['LockerCode'] as String,
      userId: json['UserId'] as String? ?? "",
    );
  }
}
