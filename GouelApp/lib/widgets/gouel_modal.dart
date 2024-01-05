import 'package:flutter/material.dart';
import 'package:gouel/widgets/gouel_dialog.dart';

class GouelModal {
  static Future<void> show(
    BuildContext context, {
    required Widget child,
    String? title,
    List<Widget>? actions,
    bool barrierDismissible = true,
  }) {
    return showDialog(
      context: context,
      barrierDismissible: barrierDismissible,
      builder: (BuildContext context) {
        return GouelDialog(
          key: const Key("GouelDialog"),
          title: title,
          actions: actions,
          child: child,
        );
      },
    );
  }

  static Future<void> showFuture(
    BuildContext context, {
    required Future<Widget> futureChild,
    String? title,
    bool barrierDismissible = true,
    List<Widget>? actions,
  }) {
    return showDialog(
      context: context,
      barrierDismissible: barrierDismissible,
      builder: (BuildContext context) {
        return FutureBuilder<Widget>(
          future: futureChild,
          builder: (BuildContext context, AsyncSnapshot<Widget> snapshot) {
            if (snapshot.connectionState == ConnectionState.waiting) {
              return GouelDialog(
                  title: title,
                  actions: actions,
                  child: const SizedBox(
                    height: 100,
                    child: Center(child: CircularProgressIndicator()),
                  ));
            }
            if (snapshot.hasData) {
              return GouelDialog(
                title: title,
                child: SingleChildScrollView(
                  child: snapshot.data!,
                ),
              );
            }
            return const SizedBox
                .shrink(); // ou gérer les erreurs si nécessaire
          },
        );
      },
    );
  }
}
