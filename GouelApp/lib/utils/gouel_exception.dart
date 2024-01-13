import 'package:flutter/material.dart';
import 'package:gouel/widgets/gouel_snackbar.dart';

enum GouelExceptionState {
  warning,
  critical,
  info;

  Color get color {
    switch (this) {
      case GouelExceptionState.warning:
        return Colors.orange.shade500; // Couleur pour les avertissements
      case GouelExceptionState.critical:
        return Colors.red.shade500; // Couleur pour les erreurs critiques
      case GouelExceptionState.info:
        return Colors.blue.shade500; // Couleur pour les informations
      default:
        return Colors.grey.shade500; // Couleur par défaut
    }
  }
}

class GouelException implements Exception {
  final String message;
  final Map<String, dynamic>? data;

  final GouelExceptionState state;

  static void inform(Object e, BuildContext context) {
    if (e is GouelException) {
      showGouelSnackbar(context, e.message, e.state.color);
    }
  }

  GouelException(
      {this.message = "Une erreur est survenue.",
      this.state = GouelExceptionState.info,
      this.data});

  @override
  String toString() => message;
}
