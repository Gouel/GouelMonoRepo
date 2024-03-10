import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:gouel/models/event_model.dart';
import 'package:gouel/providers/payment/payment_provider.dart';
import 'package:gouel/services/gouel_session_service.dart';
import 'package:gouel/widgets/gouel_button.dart';
import 'package:gouel/widgets/paragraph.dart';
import 'package:sumup/sumup.dart';

class SumUpPaymentProvider extends PaymentProvider {
  const SumUpPaymentProvider();

  @override
  Future<bool> confirmPayment(Map<String, Object> data) async {
    double? total = data["amount"] as double?;
    if (total == null) {
      return false;
    }
    //Création de la transaction SumUp
    try {
      var payment = SumupPayment(
        title: !data.containsKey("title")
            ? "Paiement Gouel"
            : data["title"] as String,
        total: total,
        currency: 'EUR',
        foreignTransactionId: '',
        saleItemsCount: 0,
        skipSuccessScreen: true,
        skipFailureScreen: true,
        tipOnCardReader: false,
        customerEmail: null,
        customerPhone: null,
      );

      final SumupPaymentRequest request = SumupPaymentRequest(payment);
      final SumupPluginCheckoutResponse checkout =
          await Sumup.checkout(request);
      return checkout.success ?? false;
    } catch (e) {
      if (kDebugMode) {
        print(e);
      }
      return false;
    }
  }

  @override
  Future<void> init() async {
    Event event = GouelSession().retrieve("event");
    final String affiliateKey =
        event.options?["SumUpOptions"]?["AffiliateKey"] as String? ?? "";
    //Initialisation de SumUp
    await Sumup.init(affiliateKey);

    //Connexion à SumUp (login affiche form de connexion)
    await Sumup.login();

    //Connexion à SumUp (loginToken)
    /*
    <== Provoque internal error lors de checkout ??  ==>
    final String loginToken =
        event.options?["sumup_options"]?["sumup_token"] as String? ?? "";
    await Sumup.loginWithToken(loginToken);
    */

    return;
  }

  Future<void> openSumUps() async {
    bool loggedIn = await Sumup.isLoggedIn ?? false;
    if (loggedIn) {
      await Sumup.openSettings();
    }
  }

  Future<void> logout() async {
    bool loggedIn = await Sumup.isLoggedIn ?? false;
    if (loggedIn) {
      await Sumup.logout();
    }
  }

  @override
  Future<Widget> getOptions() async {
    bool loggedIn = await Sumup.isLoggedIn ?? false;

    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      mainAxisSize: MainAxisSize.max,
      children: [
        const Paragraph(
          type: ParagraphType.heading,
          content: "Options de SumUp",
        ),
        Paragraph.space(),
        Paragraph(
            type: ParagraphType.text,
            content:
                "Etat de la connexion : ${loggedIn ? "Connecté" : "Déconnecté"}"),
        Paragraph.space(),
        GouelButton(
          onTap: () async {
            if (loggedIn) {
              await logout();
            } else {
              await init();
            }
          },
          text: loggedIn ? "Déconnexion" : "Connexion",
        ),
        Paragraph.space(),
        GouelButton(
          onTap: () async {
            await openSumUps();
          },
          text: "Choisir le SumUp",
        ),
      ],
    );
  }
}
