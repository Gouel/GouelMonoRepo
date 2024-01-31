import 'package:flutter/widgets.dart';

class PaymentProvider {
  const PaymentProvider();

  Future<void> init() async {}

  Future<bool> confirmPayment(Map<String, Object> data) async {
    throw UnimplementedError();
  }

  Future<Widget> getOptions() {
    throw UnimplementedError();
  }
}
