class Event {
  final String id;
  final String title;
  final List<EventTicket> eventTickets;

  Event({required this.id, required this.title, required this.eventTickets});

  factory Event.fromJson(Map<String, dynamic> json) {
    List<EventTicket> tickets = [];

    List<dynamic> listJson = json["EventTickets"];
    tickets.addAll(listJson.map((e) => EventTicket.fromJson(e)).toList());

    return Event(
        id: json['ID'] as String,
        title: json['Title'] as String,
        eventTickets: tickets);
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

class EventTicket {
  final String title;
  final String eventTicketCode;

  EventTicket({required this.title, required this.eventTicketCode});

  factory EventTicket.fromJson(Map<String, dynamic> json) {
    return EventTicket(
      title: json["Title"] as String,
      eventTicketCode: json["EventTicketCode"] as String,
    );
  }
}
