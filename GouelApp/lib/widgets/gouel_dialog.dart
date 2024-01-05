import 'package:flutter/material.dart';

class GouelDialog extends StatelessWidget {
  final String? title;
  final Widget child;
  final List<Widget>? actions;
  final Widget? icon;

  const GouelDialog({
    Key? key,
    this.title,
    required this.child,
    this.actions,
    this.icon,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      shape: const RoundedRectangleBorder(
          borderRadius: BorderRadius.all(Radius.circular(32.0))),
      backgroundColor: Theme.of(context).scaffoldBackgroundColor,
      icon: icon, // Style personnalisé
      title: title != null
          ? Text(
              title!,
              textAlign: TextAlign.center,
              style: Theme.of(context).textTheme.titleLarge,
            )
          : null,
      content: SingleChildScrollView(
        child: child,
      ),
      actions: actions,
    );
  }
}
