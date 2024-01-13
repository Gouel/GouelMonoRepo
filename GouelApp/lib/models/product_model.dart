import 'package:flutter/material.dart';

/// Product Model.

class Product {
  final int? amount;
  final DateTime? endOfSale;
  final bool hasAlcohol;
  final IconData icon;
  final String label;
  final double price;
  final String productCode;
  final int purchased;

  Product({
    this.amount,
    this.endOfSale,
    required this.hasAlcohol,
    required this.icon,
    required this.label,
    required this.price,
    required this.productCode,
    required this.purchased,
  });

  factory Product.fromJson(Map<String, dynamic> json) {
    var endOfSale = json['EndOfSale'] != null && json['EndOfSale'] != ''
        ? DateTime.parse(json['EndOfSale'])
        : null;

    return Product(
      amount: json['Amount'],
      endOfSale: endOfSale,
      hasAlcohol: json['HasAlcohol'],
      icon: ProductIcon.fromString(json["Icon"]).icon, // to change
      label: json['Label'],
      price: json['Price'],
      productCode: json['ProductCode'],
      purchased: json['Purchased'],
    );
  }

  @override
  String toString() => "Product<$productCode>";

  Map<String, dynamic> toJson() => {
        'amount': amount,
        'endOfSale': endOfSale?.toIso8601String(),
        'hasAlcohol': hasAlcohol,
        'icon': icon,
        'label': label,
        'price': price,
        'productCode': productCode,
        'purchased': purchased,
      };
}

enum ProductIcon {
  pizza(Icons.local_pizza),
  sportsDrink(Icons.sports_bar),
  wineBar(Icons.wine_bar),
  glassCup(Icons.local_bar),
  inventory2(Icons.inventory_2),
  ;

  const ProductIcon(this.icon);
  final IconData icon;

  static ProductIcon fromString(String name) {
    return ProductIcon.values.firstWhere((e) {
      return e.name == name;
    }, orElse: () => ProductIcon.inventory2);
  }
}
