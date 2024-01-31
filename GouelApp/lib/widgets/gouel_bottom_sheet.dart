import 'package:flutter/material.dart';

class GouelBottomSheet extends StatelessWidget {
  final String title;
  final Widget child;

  const GouelBottomSheet({super.key, required this.title, required this.child});

  static void launch(
      {required BuildContext context,
      required GouelBottomSheet bottomSheet,
      bool isDismissible = true}) {
    showModalBottomSheet(
        context: context,
        builder: (_) => bottomSheet,
        isDismissible: isDismissible,
        isScrollControlled: true);
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      padding:
          EdgeInsets.only(bottom: MediaQuery.of(context).viewInsets.bottom),
      child: Container(
        padding: const EdgeInsets.all(16.0),
        decoration: BoxDecoration(
          color: Theme.of(context).scaffoldBackgroundColor,
          borderRadius: const BorderRadius.only(
            topLeft: Radius.circular(24.0),
            topRight: Radius.circular(24.0),
          ),
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: <Widget>[
            Text(
              title,
              textAlign: TextAlign.center,
              style: Theme.of(context).textTheme.titleLarge,
            ),
            const SizedBox(height: 16.0),
            child,
          ],
        ),
      ),
    );
  }
}
