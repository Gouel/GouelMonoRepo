import 'package:flutter/material.dart';
import 'package:gouel/models/event_model.dart';
import 'package:gouel/models/gouel_cart.dart';
import 'package:gouel/models/product_model.dart';
import 'package:gouel/models/ticket_model.dart';
import 'package:gouel/models/transcations_model.dart';
import 'package:gouel/services/gouel_session_service.dart';
import 'package:gouel/services/mail_service.dart';
import 'package:gouel/utils/gouel_request.dart';
import 'package:gouel/utils/gouel_exception.dart';
import 'package:http/http.dart';
import 'package:logging/logging.dart';

class GouelApiService with ChangeNotifier {
  final GouelSession _sessionService = GouelSession();
  final BuildContext context;
  String? _token;
  bool isLogged = false;
  Logger logger = Logger("GouelApiService");

  GouelApiService(this.context) {
    _init();
  }

  Future<void> _init() async {
    _token = _sessionService.retrieve("token");
  }

  Future<void> authenticateWithTicket(String ticketID) async {
    try {
      var result = await GouelRequest.post("/token/auth/ticket")
          .send(context, data: {'TicketId': ticketID}, noHeaders: true);

      await _saveToken(result);
      isLogged = true;
    } catch (e) {
      isLogged = false;
      if (e is GouelException) {
        rethrow;
      } else if (e is ClientException) {
        throw GouelException(
          message: "$e",
          state: GouelExceptionState.critical,
        );
      } else {
        throw GouelException(
          message: "L'authentification a échoué",
          state: GouelExceptionState.critical,
        );
      }
    }
  }

  Future<void> _saveToken(Map<String, dynamic> data) async {
    _token = data['token'];
    _sessionService.store('token', _token);

    var decodedToken = await getTokenDecoded();
    if (decodedToken != null) {
      _sessionService.store('infos', decodedToken);
      _scheduleTokenRefresh(decodedToken['exp']);
    }
    notifyListeners();
  }

  Future<Map<String, dynamic>?> getTokenDecoded() async {
    try {
      return await GouelRequest.get("/token/view").send(context);
    } catch (e) {
      logger.severe('Error getting token info: $e');
      return null;
    }
  }

  Future<void> _refreshToken() async {
    try {
      var result = await GouelRequest.post("/token/refresh").send(context);
      await _saveToken(result);
    } catch (e) {
      logger.severe('Error refreshing token: $e');
      logout(); // Déconnecter l'utilisateur si le rafraîchissement échoue
    }
  }

  void _scheduleTokenRefresh(int exp) {
    var now = DateTime.now().toUtc();
    var expDate = DateTime.fromMillisecondsSinceEpoch(exp * 1000, isUtc: true);
    var halfLife = expDate
        .subtract(Duration(seconds: (expDate.difference(now).inSeconds ~/ 2)));

    Future.delayed(halfLife.difference(now), () {
      _refreshToken();
    });
  }

  void logout({BuildContext? buildContext}) {
    _token = null;
    isLogged = false;
    _sessionService.remove('token');
    _sessionService.remove('infos');
    notifyListeners();

    if (buildContext != null) {
      Navigator.pushNamedAndRemoveUntil(buildContext, "/", (r) => false);
    }
  }

  // EVENTS

  Future<EmailService?> getEmailService(context) async {
    try {
      Event event = GouelSession().retrieve('event');
      var response = await GouelRequest.get("/events/${event.id}/smtp")
          .send(context) as Map<String, dynamic>;

      int port = int.parse(response["SMTPPort"]);

      return EmailService(
          host: response["SMTPServer"] as String,
          port: port,
          username: response["Email"] as String,
          password: response["EmailPassword"] as String,
          isSecure: response["SMTPSSL"] as bool);
    } catch (e) {
      return null;
    }
  }

  Future<List<Event>> getEvents(context) async {
    try {
      var response = await GouelRequest.get("/events").send(context);

      List<Event> events = (response as List)
          .map((eventJson) => Event.fromJson(eventJson))
          .toList();
      return events;
    } catch (e) {
      logger.severe('Erreur lors de la récupération des événements: $e');
      return [];
    }
  }

  Future<List<Locker>> getEventLockers(context, String eventId) async {
    try {
      var response =
          await GouelRequest.get("/events/$eventId/lockers").send(context);

      List<Locker> lockers = (response as List)
          .map((eventJson) => Locker.fromJson(eventJson))
          .toList();
      return lockers;
    } catch (e) {
      logger.severe('Erreur lors de la récupération des vestiaires: $e');
      return [];
    }
  }

  Future<void> setEventLocker(context, Locker locker) async {
    try {
      final Event event = GouelSession().retrieve("event") as Event;
      await GouelRequest.put("/events/${event.id}/lockers")
          .send(context, data: {
        "LockerCode": locker.lockerCode,
        "UserId": locker.userId,
      });
    } catch (e) {
      if (e is GouelException) {
        GouelException.inform(e, context);
      }
    }
  }

  Future<List<Product>> getEventProducts(context) async {
    try {
      Event event = GouelSession().retrieve('event');

      var response =
          await GouelRequest.get("/events/${event.id}/products").send(context);

      List<Product> products = (response as List)
          .map((eventJson) => Product.fromJson(eventJson))
          .toList();
      return products;
    } catch (e) {
      logger.severe('Erreur lors de la récupération des produits: $e');
      return [];
    }
  }

  // Tickets

  Future<String?> createTicket(context, Map<String, dynamic> ticket) async {
    try {
      Event event = GouelSession().retrieve("event");

      Map<String, dynamic> response = await GouelRequest.post(
              "/tickets/${event.id}/${ticket['EventTicketCode']}")
          .send(context, data: ticket) as Map<String, dynamic>;

      return response["TicketId"] as String;
    } catch (e) {
      return null;
    }
  }

  Future<TicketInfos?> getTicketInfos(context, String ticketId) async {
    try {
      final Event event = GouelSession().retrieve("event") as Event;

      var response = await GouelRequest.get("/tickets/${event.id}/$ticketId")
          .send(context);
      return TicketInfos.fromJson(response);
    } catch (e) {
      return null;
    }
  }

  Future<List<TicketInfos>> getAllTicketInfos(context) async {
    try {
      final Event event = GouelSession().retrieve("event") as Event;

      var response = await GouelRequest.get("/tickets/${event.id}")
          .send(context) as List<dynamic>;
      return response.map((e) {
        return TicketInfos.fromJson(e as Map<String, dynamic>);
      }).toList();
    } catch (e) {
      return [];
    }
  }

  Future<ValidateState> validateTicket(String ticketID) async {
    try {
      final Event event = GouelSession().retrieve("event") as Event;
      await GouelRequest.post("/tickets/${event.id}/validate")
          .send(context, data: {"TicketId": ticketID}) as Map<String, dynamic>;

      return ValidateState.ok;
    } catch (e) {
      if (e is GouelException && e.message == "Le ticket a déjà été validé") {
        return ValidateState.alreadyValidated;
      }

      return ValidateState.invalid;
    }
  }

  Future<bool> setTicketSAM(String ticketID, {bool isSAM = true}) async {
    try {
      final Event event = GouelSession().retrieve("event") as Event;
      await GouelRequest.put("/tickets/${event.id}/sam").send(context,
          data: {"TicketId": ticketID, "IsSam": isSAM}) as Map<String, dynamic>;
      return true;
    } catch (e) {
      return false;
    }
  }

  // Users

  Future<bool> addTransaction(String userId, Transaction transaction) async {
    try {
      Event event = GouelSession().retrieve("event") as Event;
      await GouelRequest.post("/users/transaction/${event.id}/$userId")
          .send(context, data: transaction.toJson());
      return true;
    } catch (e) {
      return false;
    }
  }

  Future<String?> addUser(Map<String, dynamic> user) async {
    try {
      final Event event = GouelSession().retrieve("event") as Event;
      Map<String, dynamic> response =
          await GouelRequest.post("/users/event/${event.id}")
              .send(context, data: user);
      return response["UserId"];
    } catch (e) {
      return null;
    }
  }

  Future<Map<String, dynamic>> userPay(String ticketId, GouelCart cart) async {
    Event event = GouelSession().retrieve("event") as Event;
    try {
      Map<String, dynamic> response =
          await GouelRequest.post("/users/pay/${event.id}/$ticketId")
              .send(context, data: cart.toJson());
      return {"statusCode": 200, "body": response};
    } catch (e) {
      if (e is GouelException) {
        return e.data ?? {};
      }
      return {};
    }
  }
}
