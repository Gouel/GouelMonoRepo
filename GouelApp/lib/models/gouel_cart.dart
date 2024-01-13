import 'package:gouel/models/product_model.dart';
import 'package:gouel/models/transcations_model.dart';
import 'package:gouel/services/gouel_session_service.dart';

class CartItem {
  final Product product;
  int quantity;

  CartItem({required this.product, this.quantity = 1});
}

class GouelCart {
  final List<CartItem> _items = [];

  List<CartItem> get items => _items;

  void addProduct(Product product) {
    var foundItem = _items.firstWhere(
      (item) => item.product.productCode == product.productCode,
      orElse: () => CartItem(product: product),
    );

    if (!_items.contains(foundItem)) {
      _items.add(foundItem);
    } else {
      foundItem.quantity++;
    }

    saveCart(); // Sauvegarde le panier
  }

  void removeProduct(Product product, {bool all = false}) {
    var foundItem = _items.firstWhere(
      (item) => item.product.productCode == product.productCode,
      orElse: () => CartItem(product: product),
    );

    if (all) {
      _items.remove(foundItem);
    } else {
      if (foundItem.quantity > 1) {
        foundItem.quantity--;
      } else {
        _items.remove(foundItem);
      }
    }

    saveCart(); // Sauvegarde le panier
  }

  void saveCart() {
    GouelSession().store("cashless_cart", _items);
  }

  GouelCart loadCart() {
    List<CartItem>? cart = GouelSession().retrieve("cashless_cart");
    _items.clear();
    _items.addAll(cart ?? []);
    return this;
  }

  double get total {
    double total = 0;
    for (var item in _items) {
      total += item.product.price * item.quantity;
    }
    return total;
  }

  int get length {
    return _items.length;
  }

  CartItem? get(int index) {
    if (index < 0 || index >= _items.length) {
      return null;
    }
    return _items[index];
  }

  List<Map<String, dynamic>> toJson() {
    List<Map<String, dynamic>> json = [];
    for (var item in _items) {
      json.add(
          {"ProductCode": item.product.productCode, "Amount": item.quantity});
    }
    return json;
  }

  void clear() {
    _items.clear();
    saveCart();
  }
}
