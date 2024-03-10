class Event {
  final String id;
  final String title;
  final List<EventTicket> eventTickets;
  final Map<String, dynamic>? options;

  Event(
      {required this.id,
      required this.title,
      required this.eventTickets,
      this.options});

  factory Event.fromJson(Map<String, dynamic> json) {
    List<EventTicket> tickets = [];

    List<dynamic> listJson = json["EventTickets"];
    tickets.addAll(listJson.map((e) => EventTicket.fromJson(e)).toList());
    return Event(
        id: json['ID'] as String,
        title: json['Title'] as String,
        eventTickets: tickets,
        options: json["Options"] as Map<String, dynamic>?);
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

  bool isTaken() {
    return userId.isNotEmpty && userId != "000000000000000000000000";
  }
}

class EventTicket {
  final String title;
  final String eventTicketCode;
  final Map<String, double> price;

  EventTicket(
      {required this.title,
      required this.eventTicketCode,
      required this.price});

  factory EventTicket.fromJson(Map<String, dynamic> json) {
    Map<String, double> price = {};
    json["Price"].forEach((key, value) {
      price[key] = value.toDouble();
    });
    return EventTicket(
      title: json["Title"] as String,
      eventTicketCode: json["EventTicketCode"] as String,
      price: price,
    );
  }
}
