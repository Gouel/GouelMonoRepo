import 'package:flutter/material.dart';
import 'package:gouel/providers/payment/payment_provider.dart';
import 'package:gouel/widgets/paragraph.dart';

class BlankPaymentProvider extends PaymentProvider {
  const BlankPaymentProvider();

  @override
  Future<bool> confirmPayment(Map<String, Object> data) async {
    return true;
  }

  @override
  Future<void> init() async {
    return;
  }

  @override
  Future<Widget> getOptions() async {
    return const Column(
      mainAxisAlignment: MainAxisAlignment.center,
      mainAxisSize: MainAxisSize.max,
      children: [
        Paragraph(
          type: ParagraphType.heading,
          content: "Aucune options disponibles.",
        )
      ],
    );
  }
}
