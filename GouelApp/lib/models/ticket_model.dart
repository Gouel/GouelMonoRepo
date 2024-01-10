class TicketInfos {
  final String id;
  final String eventId;
  final String eventTicketCode;
  final Map<String, dynamic> user;
  final bool isUsed;
  final bool isSam;

  TicketInfos(
      {required this.id,
      required this.eventId,
      required this.eventTicketCode,
      required this.user,
      required this.isUsed,
      required this.isSam});

  factory TicketInfos.fromJson(Map<String, dynamic> json) {
    return TicketInfos(
      id: json["ID"],
      eventId: json["EventId"],
      eventTicketCode: json["EventTicketCode"],
      user: json["User"] as Map<String, dynamic>,
      isUsed: json["IsUsed"],
      isSam: json["IsSam"] ?? false,
    );
  }
}

enum ValidateState {
  ok,
  invalid,
  alreadyValidated,
}
