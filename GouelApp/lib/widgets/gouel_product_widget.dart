import 'package:flutter/material.dart';
import 'package:gouel/models/product_model.dart';
import 'package:gouel/services/gouel_storage_service.dart';

class ProductWidget extends StatefulWidget {
  final Product product;
  final VoidCallback onTap;

  const ProductWidget({super.key, required this.product, required this.onTap});

  @override
  ProductWidgetState createState() => ProductWidgetState();
}

class ProductWidgetState extends State<ProductWidget> {
  bool _isTapped = false;
  bool _showTitle = false;

  @override
  void initState() {
    super.initState();
    _getShowTitle();
  }

  void _animateIcon() {
    setState(() {
      _isTapped = true;
    });

    Future.delayed(const Duration(milliseconds: 250), () {
      widget.onTap();
      setState(() {
        _isTapped = false;
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: _animateIcon,
      child: Padding(
        padding: const EdgeInsets.all(8.0),
        child: ClipRRect(
          borderRadius: const BorderRadius.all(Radius.circular(8)),
          child: GridTile(
            child: AnimatedContainer(
              duration: const Duration(milliseconds: 500),
              color: _isTapped ? Colors.green : Colors.transparent,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.center,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  if (_showTitle)
                    Center(
                        child: Text(
                      widget.product.label,
                      style: const TextStyle(
                        fontSize: 18,
                      ),
                    )),
                  Stack(
                    alignment: Alignment.bottomRight,
                    children: [
                      Icon(
                        widget.product.icon,
                        size: 50,
                      ),
                      if (widget.product.hasAlcohol)
                        Container(
                          padding: const EdgeInsets.all(1),
                          decoration: const BoxDecoration(
                            color: Colors.red,
                            shape: BoxShape.circle,
                          ),
                          child: const Icon(
                            Icons.new_releases_outlined,
                            color: Colors.white,
                            size: 20,
                          ),
                        ),
                    ],
                  ),
                  Text("${widget.product.price.toStringAsFixed(2)}€")
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  void _getShowTitle() async {
    _showTitle = await GouelStorage().retrieve("product_show_title") ?? true;
    setState(() {});
  }
}
