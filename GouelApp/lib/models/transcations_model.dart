import 'package:flutter/material.dart';

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
    TransactionType type = data["type"] == "credit"
        ? TransactionType.credit
        : TransactionType.debit;
    DateTime dateTime = DateTime.parse(data["date"] as String);
    String eventID = data["eventID"] as String;
    double amount = data["amount"] as double;
    PaymentMethod? paymentMethod;

    if (data.containsKey("payement_type")) {
      paymentMethod = PaymentMethod.fromString(data["payement_type"] as String);
    }

    List<dynamic> preprocessedCart = data["cart"] as List;

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
      "type": transactionType.name,
      "date": dateTime.toIso8601String(),
      "eventId": eventID,
      "cart": cart
          .map(
            (e) => e.toJson(),
          )
          .toList(),
      "amount": amount,
    };

    if (paymentMethod != null) {
      data["payementType"] = paymentMethod!.name;
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
      productCode: data["productCode"] as String,
      amount: data["amount"] as int,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      "productCode": productCode,
      "amount": amount,
    };
  }
}

enum TransactionType { debit, credit }

enum PaymentMethod {
  // Maximum 6 available
  espece(desc: "Espèces", available: true, icon: Icons.euro),
  carte(desc: "Carte bleue", available: true, icon: Icons.credit_card),
  helloasso(desc: "HelloAsso", available: false),
  ;

  const PaymentMethod(
      {required this.desc,
      this.available = false,
      this.icon = Icons.radio_button_off});

  final String desc;
  final bool available;
  final IconData icon;

  static PaymentMethod? fromString(String method) {
    switch (method) {
      case "carte":
        return PaymentMethod.carte;
      case "espece":
        return PaymentMethod.espece;
      case "helloasso":
        return PaymentMethod.helloasso;
      default:
        return null;
    }
  }
}
