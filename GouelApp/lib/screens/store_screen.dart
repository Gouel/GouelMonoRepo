import 'package:flutter/material.dart';
import 'package:flutter_slidable/flutter_slidable.dart';
import 'package:gouel/models/gouel_cart.dart';
import 'package:gouel/models/product_model.dart';
import 'package:gouel/models/ticket_model.dart';
import 'package:gouel/services/gouel_api_service.dart';
import 'package:gouel/services/qr_scanner_service.dart';
import 'package:gouel/utils/gouel_getter.dart';
import 'package:gouel/widgets/gouel_bottom_sheet.dart';
import 'package:gouel/widgets/gouel_button.dart';
import 'package:gouel/widgets/gouel_modal.dart';
import 'package:gouel/widgets/gouel_product_widget.dart';
import 'package:gouel/widgets/gouel_scaffold.dart';
import 'package:logging/logging.dart';
import 'package:provider/provider.dart';

class StoreScreen extends StatefulWidget {
  const StoreScreen({Key? key}) : super(key: key);

  @override
  StoreScreenState createState() => StoreScreenState();
}

class StoreScreenState extends State<StoreScreen> {
  List<Product> products = [];

  late GouelCart cart;

  bool isCartOpen = false;

  Logger logger = Logger("StoreScreen");

  @override
  void dispose() {
    super.dispose();
  }

  void _loadProducts() async {
    products = await Provider.of<GouelApiService>(context, listen: false)
        .getEventProducts(context);

    cart = GouelCart().loadCart();
    setState(() {});
  }

  @override
  void initState() {
    super.initState();
    _loadProducts();
  }

  final GlobalKey<ScaffoldState> _scaffoldKey = GlobalKey<ScaffoldState>();
  @override
  Widget build(BuildContext context) {
    return GouelScaffold(
      key: _scaffoldKey,
      appBar: AppBar(
        title: const Text('Buvette'),
      ),
      body: products.isEmpty
          ? const Center(
              child: Icon(
              Icons.local_mall,
              size: 60,
            ))
          : GridView.builder(
              gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                crossAxisCount: 2,
              ),
              itemCount: products.length,
              itemBuilder: (context, index) {
                return ProductWidget(
                  product: products[index],
                  onTap: () {
                    cart.addProduct(products[index]);
                    cart.saveCart();
                  },
                );
              },
            ),
      persistentFooterButtons: [
        GouelButton(
          text: "Panier",
          onTap: _showCart,
          icon: Icons.shopping_bag,
        )
      ],
    );
  }

  void _showCart() {
    GouelBottomSheet.launch(
      context: context,
      bottomSheet: GouelBottomSheet(
        title: 'Panier',
        child: _buildCart(context),
      ),
    );
  }

  Widget _buildCart(context) {
    return StatefulBuilder(
      builder: (BuildContext context, StateSetter setModalState) {
        return SizedBox(
          height: MediaQuery.of(context).size.height * 0.5,
          child: Column(
            children: [
              const SizedBox(
                height: 8,
              ),
              Expanded(
                child: ListView.builder(
                  itemCount: cart.length,
                  itemBuilder: (context, index) {
                    final CartItem cartItem = cart.get(index)!;
                    final quantity = cartItem.quantity;
                    final Product product = cartItem.product;

                    return Slidable(
                      endActionPane: ActionPane(
                        motion: const ScrollMotion(),
                        children: [
                          SlidableAction(
                            backgroundColor: Colors.blue,
                            icon: Icons.remove_circle,
                            onPressed: (context) {
                              cart.removeProduct(product);
                              setModalState(() {});
                            },
                          ),
                          SlidableAction(
                            backgroundColor: Colors.red,
                            icon: Icons.delete,
                            onPressed: (context) {
                              cart.removeProduct(product, all: true);
                              setModalState(() {});
                            },
                          ),
                        ],
                      ),
                      child: ListTile(
                        title: Text(product.label),
                        leading: Icon(
                          product.icon,
                          color: product.hasAlcohol ? Colors.amber[600] : null,
                        ),
                        subtitle: Text('${product.price.toStringAsFixed(2)}€'),
                        trailing: Text(
                          'x$quantity',
                          style: const TextStyle(fontSize: 18),
                        ),
                      ),
                    );
                  },
                ),
              ),
              Padding(
                padding: const EdgeInsets.only(left: 8.0, right: 8.0),
                child: SizedBox(
                  width: double.maxFinite,
                  child: Row(
                    children: [
                      Padding(
                        padding: const EdgeInsets.only(right: 8.0),
                        child: GouelButton(
                            text: null,
                            color: Colors.red,
                            onTap: () {
                              cart.clear();
                              setModalState(() {});
                            },
                            icon: Icons.delete),
                      ),
                      Expanded(
                        child: GouelButton(
                          color: cart.total > 0 ? Colors.blue : Colors.grey,
                          onTap: () {
                            if (cart.total > 0) {
                              //paymentProcess
                              Navigator.of(context).pop();
                              paymentProcess();
                            }
                          },
                          text: "Payer ${cart.total.toStringAsFixed(2)}€",
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(
                height: 8,
              ),
            ],
          ),
        );
      },
    );
  }

  double calculateTotal(Map<Product, int> cart) {
    double total = 0.0;
    cart.forEach((product, quantity) {
      total += product.price * quantity;
    });
    return total;
  }

  //paymentProcess
  void paymentProcess() {
    QRScannerService().scanQR(
      context,
      "Scanner ticket",
      (result) async {
        TicketInfos? ticketInfos =
            await getTicketInfos(context, result, withSnackBar: true);
        if (ticketInfos == null) return;
        if (!mounted) return;

        Map<String, dynamic> paymentReturn =
            await Provider.of<GouelApiService>(context, listen: false)
                .userPay(ticketInfos.id, cart);
        processPaymentReturn(paymentReturn);
      },
      (close) => null,
    );
  }

  void processPaymentReturn(Map<String, dynamic> paymentReturn) {
    int statusCode = paymentReturn["body"]["code"];
    Widget statusWidget;
    logger.severe(paymentReturn);
    switch (statusCode) {
      case 0x0:
        // payment Success
        statusWidget = statusWidget = buildStatusWidget(
          Icons.check_circle,
          Colors.green,
          "Paiement effectué",
        );
        cart.clear();
        break;
      case 0x1:
        // donnéees invalides
        statusWidget = statusWidget = buildStatusWidget(
          Icons.error,
          Colors.red,
          "Panier Invalide",
        );
        break;
      case 0x2:
        statusWidget = statusWidget = buildStatusWidget(
          Icons.question_mark,
          Colors.orange,
          "Ticket invalide",
        );
        break;
      case 0x3:
        // erreur produits (endOfSale ou Alcohol ou produit invalide)
        if (paymentReturn['body']['error']['code'] == 0x0) {
          statusWidget = statusWidget = buildStatusWidget(
            Icons.error,
            Colors.red,
            "Produit invalide",
          );
        } else if (paymentReturn['body']['error']['code'] == 0x1) {
          statusWidget = statusWidget = buildStatusWidget(
            Icons.event_busy,
            Colors.red,
            "${paymentReturn['body']['error']['data']['Label']} n'est plus en vente",
          );
        } else if (paymentReturn['body']['error']['code'] == 0x2) {
          statusWidget = statusWidget = buildStatusWidget(
            Icons.error,
            Colors.red,
            "${paymentReturn['body']['error']['data']['Label']} contient de l'alcool",
          );
        } else {
          statusWidget = statusWidget = buildStatusWidget(
            Icons.error,
            Colors.red,
            "Erreur inconnue",
          );
        }
        break;
      case 0x4:
        // solde insuffisant
        statusWidget = statusWidget = buildStatusWidget(
          Icons.money_off,
          Colors.red,
          "Solde insuffisant",
        );
        break;
      default:
        statusWidget = statusWidget = buildStatusWidget(
          Icons.error,
          Colors.red,
          "Erreur inconnue",
        );
    }

    //show status widget
    GouelModal.show(
      context,
      child: statusWidget,
    );
  }

  Widget buildStatusWidget(IconData icon, Color iconColor, String text) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(icon, color: iconColor, size: 60),
        const SizedBox(height: 20),
        Text(
          text,
          textAlign: TextAlign.center,
          style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
        ),
      ],
    );
  }
}
