import 'package:flutter/material.dart';
import 'package:gouel/models/ticket_model.dart';
import 'package:gouel/services/gouel_api_service.dart';
import 'package:gouel/widgets/gouel_snackbar.dart';
import 'package:provider/provider.dart';

Future<TicketInfos?> getTicketInfos(BuildContext context, String ticketId,
    {bool withSnackBar = true}) async {
  TicketInfos? ticketInfos =
      await Provider.of<GouelApiService>(context, listen: false)
          .getTicketInfos(context, ticketId);
  if (context.mounted) {
    if (ticketInfos == null && withSnackBar) {
      showGouelSnackbar(context, "Ticket invalide", Colors.red, duration: 5);
      return null;
    }
  }

  return ticketInfos;
}

String getMajeurMineur(String dob) {
  DateTime birthday = DateTime.parse(dob);
  int age = DateTime.now().difference(birthday).inDays ~/ 365;
  return age >= 18 ? "Majeur" : "Mineur";
}

bool isDigit(String s, int idx) => (s.codeUnitAt(idx) ^ 0x30) <= 9;
