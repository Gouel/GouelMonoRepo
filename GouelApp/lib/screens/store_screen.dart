import 'package:flutter/material.dart';
import 'package:gouel/models/gouel_cart.dart';
import 'package:gouel/models/product_model.dart';
import 'package:gouel/models/ticket_model.dart';
import 'package:gouel/services/gouel_api_service.dart';
import 'package:gouel/services/qr_scanner_service.dart';
import 'package:gouel/utils/gouel_getter.dart';
import 'package:gouel/widgets/gouel_bottom_sheet.dart';
import 'package:gouel/widgets/gouel_button.dart';
import 'package:gouel/widgets/gouel_dialog.dart';
import 'package:gouel/widgets/gouel_modal.dart';
import 'package:gouel/widgets/gouel_product_widget.dart';
import 'package:gouel/widgets/gouel_scaffold.dart';
import 'package:gouel/widgets/icon_badge.dart';
import 'package:logging/logging.dart';
import 'package:provider/provider.dart';

class StoreScreen extends StatefulWidget {
  const StoreScreen({Key? key}) : super(key: key);

  @override
  StoreScreenState createState() => StoreScreenState();
}

class StoreScreenState extends State<StoreScreen> {
  List<Product> products = [];

  GouelCart cart = GouelCart();

  bool isCartOpen = false;

  Logger logger = Logger("StoreScreen");

  @override
  void dispose() {
    super.dispose();
  }

  void _loadProducts() async {
    products = await Provider.of<GouelApiService>(context, listen: false)
        .getEventProducts(context);

    cart.loadCart();
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
          : Column(
              children: [
                Expanded(
                  child: GridView.builder(
                    gridDelegate:
                        const SliverGridDelegateWithFixedCrossAxisCount(
                      crossAxisCount: 2,
                    ),
                    itemCount: products.length,
                    itemBuilder: (context, index) {
                      return ProductWidget(
                        product: products[index],
                        onTap: () {
                          cart.addProduct(products[index]);
                          cart.saveCart();
                          setState(() {});
                        },
                      );
                    },
                  ),
                ),
                const SizedBox(
                  height: 8,
                ),
                const Divider(
                  height: 1,
                ),
                const SizedBox(
                  height: 8,
                ),
                SizedBox(
                    height: 50,
                    child: ListView.builder(
                      scrollDirection: Axis.horizontal,
                      itemCount: cart.items.length,
                      itemBuilder: (context, index) {
                        final CartItem item = cart.items[index];
                        return IconBadge(
                            icon: item.product.icon,
                            badgeCount: item.quantity,
                            size: 50);
                      },
                    ))
              ],
            ),
      persistentFooterButtons: [
        GouelButton(
          text: "Panier ${cart.total.toStringAsFixed(2)}€",
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

                    return ListTile(
                      title: Text(product.label),
                      leading: Icon(
                        product.icon,
                        color: product.hasAlcohol ? Colors.red[500] : null,
                        size: 40,
                      ),
                      subtitle: Text('${product.price.toStringAsFixed(2)}€'),
                      trailing: Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Text(
                            'x$quantity',
                            style: const TextStyle(fontSize: 18),
                          ),
                          const SizedBox(
                            width: 4,
                          ),
                          IconButton(
                              padding: EdgeInsets.zero,
                              onPressed: () {
                                cart.removeProduct(product);
                                setModalState(() {});
                                setState(() {});
                              },
                              icon: Icon(
                                Icons.remove_circle,
                                color: Colors.blue[300],
                              )),
                          const SizedBox(
                            width: 4,
                          ),
                          IconButton(
                            padding: EdgeInsets.zero,
                            onPressed: () {
                              cart.removeProduct(product, all: true);
                              setModalState(() {});
                              setState(() {});
                            },
                            icon: const Icon(Icons.delete),
                            color: Colors.red[500],
                          ),
                        ],
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
                            showDialog(
                                context: context,
                                builder: (builder) {
                                  return GouelDialog(
                                    title: "Vider le panier",
                                    actions: [
                                      TextButton(
                                        onPressed: () {
                                          Navigator.of(context).pop();
                                        },
                                        child: const Text(
                                          "Non",
                                        ),
                                      ),
                                      TextButton(
                                        onPressed: () {
                                          cart.clear();
                                          setModalState(() {});
                                          setState(() {});
                                          Navigator.of(context).pop();
                                        },
                                        child: const Text("Oui",
                                            style:
                                                TextStyle(color: Colors.red)),
                                      ),
                                    ],
                                    child: const Text(
                                        "Êtes-vous sûr de vouloir vider le panier ?"),
                                  );
                                });
                          },
                          icon: Icons.delete,
                        ),
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
        setState(() {});
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
