import 'package:flutter/material.dart';
import 'package:gouel/providers/payment/blank_payment_provider.dart';
import 'package:gouel/providers/payment/payment_provider.dart';
import 'package:gouel/providers/payment/sumup_payment_provider.dart';

class Transaction {
  final TransactionType transactionType;
  final DateTime dateTime;
  final String eventID;
  final List<PurchasedItem> cart;
  final double amount;
  final PaymentMethod? paymentMethod;
  Transaction(
      {required this.transactionType,
      required this.dateTime,
      required this.eventID,
      required this.cart,
      required this.amount,
      this.paymentMethod});

  factory Transaction.fromJson(Map<String, dynamic> data) {
    TransactionType type = data["Type"] == "credit"
        ? TransactionType.credit
        : TransactionType.debit;
    DateTime dateTime = DateTime.parse(data["Date"] as String);
    String eventID = data["EventID"] as String;
    double amount = data["Amount"] as double;
    PaymentMethod? paymentMethod;

    if (data.containsKey("PaymentType")) {
      paymentMethod = PaymentMethod.fromString(data["PaymentType"] as String);
    }

    List<dynamic> preprocessedCart = data["Cart"] as List;

    List<PurchasedItem> cart = preprocessedCart.map(
      (e) {
        Map<String, dynamic> map = e as Map<String, dynamic>;
        return PurchasedItem.fromJson(map);
      },
    ).toList();

    return Transaction(
      transactionType: type,
      dateTime: dateTime,
      eventID: eventID,
      cart: cart,
      amount: amount,
      paymentMethod: paymentMethod,
    );
  }

  Map<String, dynamic> toJson() {
    Map<String, dynamic> data = {
      "Type": transactionType.name,
      "Date": dateTime.toIso8601String(),
      "EventId": eventID,
      "Cart": cart
          .map(
            (e) => e.toJson(),
          )
          .toList(),
      "Amount": amount,
    };

    if (paymentMethod != null) {
      data["PaymentType"] = paymentMethod!.name;
    }

    return data;
  }
}

class PurchasedItem {
  final String productCode;
  final int amount;

  PurchasedItem({
    required this.productCode,
    required this.amount,
  });

  factory PurchasedItem.fromJson(Map<String, dynamic> data) {
    return PurchasedItem(
      productCode: data["ProductCode"] as String,
      amount: data["Amount"] as int,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      "ProductCode": productCode,
      "Amount": amount,
    };
  }
}

enum TransactionType { debit, credit }

enum PaymentMethod {
  // Maximum 6 affichable.
  //TODO utiliser options event
  especes(desc: "Espèces", icon: Icons.euro),
  carte(desc: "Carte bleue", icon: Icons.credit_card),
  sumup(
    desc: "SumUp",
    icon: Icons.add_card,
    paymentProvider: SumUpPaymentProvider(),
  ),
  helloasso(desc: "HelloAsso"),
  ;

  const PaymentMethod(
      {required this.desc,
      this.icon = Icons.radio_button_off,
      this.paymentProvider = const BlankPaymentProvider()});

  final String desc;
  final IconData icon;
  final PaymentProvider paymentProvider;

  static PaymentMethod? fromString(String method) {
    switch (method) {
      case "carte":
        return PaymentMethod.carte;
      case "especes":
        return PaymentMethod.especes;
      case "helloasso":
        return PaymentMethod.helloasso;
      case "sumup":
        return PaymentMethod.sumup;
      default:
        return null;
    }
  }
}
