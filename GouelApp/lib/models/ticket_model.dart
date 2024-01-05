class TicketInfos {
  final String id;
  final Map<String, dynamic> event;
  final Map<String, dynamic> user;
  final bool valid;
  final bool isSam;

  TicketInfos(
      {required this.id,
      required this.event,
      required this.user,
      required this.valid,
      required this.isSam});

  factory TicketInfos.fromJson(Map<String, dynamic> json) {
    return TicketInfos(
      id: json["_id"],
      event: json["event"] as Map<String, dynamic>,
      user: json["user"] as Map<String, dynamic>,
      valid: json["valid"],
      isSam: json["SAM"] ?? false,
    );
  }
}

enum ValidateState {
  ok,
  invalid,
  alreadyValidated,
}
