import 'package:flutter/material.dart';

void showGouelSnackbar(BuildContext context, String message, Color color,
    {int duration = 3}) {
  final snackBar = SnackBar(
    duration: Duration(seconds: duration),
    content: Text(
      message,
      style: Theme.of(context).textTheme.titleLarge,
    ),
    backgroundColor: color,
  );

  ScaffoldMessenger.of(context).showSnackBar(snackBar);
}
