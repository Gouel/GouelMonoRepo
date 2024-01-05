import 'package:flutter/material.dart';
import 'package:gouel/services/gouel_api_service.dart';
import 'package:gouel/services/qr_scanner_service.dart';
import 'package:gouel/utils/gouel_exception.dart';
import 'package:gouel/widgets/gouel_button.dart';
import 'package:gouel/widgets/gouel_scaffold.dart';
import 'package:gouel/widgets/gouel_server_chooser.dart';
import 'package:gouel/widgets/gouel_snackbar.dart';
import 'package:provider/provider.dart';

class LoginScreen extends StatelessWidget {
  const LoginScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final apiService = Provider.of<GouelApiService>(context, listen: false);
    final TextEditingController userIdController =
        TextEditingController(text: "65957147519d50b6305186c8");

    Future<void> authenticateUser(String userId) async {
      await apiService.authenticateWithTicket(userId);
    }

    return GouelScaffold(
      withLogOut: false,
      body: Column(
        children: <Widget>[
          Align(
            alignment: Alignment.topRight,
            child: IconButton(
              icon: const Icon(Icons.settings),
              onPressed: () {
                showModalBottomSheet(
                    isScrollControlled: true,
                    context: context,
                    builder: (BuildContext context) {
                      return const SingleChildScrollView(
                          child: GouelServerChooser());
                    });
              },
            ),
          ),

          // Logo
          Image.asset('public/assets/icon.png', height: 120),
          const SizedBox(height: 20),

          // Titre "Gouel"
          const Text(
            'Gouel',
            style: TextStyle(fontSize: 36, fontWeight: FontWeight.bold),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 40),

          // Champ de texte avec bouton valider
          Row(
            children: <Widget>[
              Expanded(
                child: TextField(
                  controller: userIdController,
                  decoration: const InputDecoration(
                    labelText: 'Identifiant',
                    border: OutlineInputBorder(),
                  ),
                ),
              ),
              const SizedBox(
                width: 5,
              ),
              GouelButton(
                text: null,
                onTap: () async {
                  String userId = userIdController.text;
                  if (userId.isNotEmpty) {
                    await _authenticate(authenticateUser, userId, context);
                  } else {
                    // Afficher un message si le champ est vide
                    showGouelSnackbar(context, "Veuillez entrer un identifiant",
                        Colors.orange.shade500);
                  }
                },
                icon: Icons.check,
              )
            ],
          ),

          // Espace vertical
          const SizedBox(height: 20),
          GouelButton(
            text: "S'authentifier avec un QRCode",
            onTap: () {
              QRScannerService().scanQR(context, "S'authentifier",
                  (result) async {
                if (result.isNotEmpty) {
                  if (context.mounted) {
                    await _authenticate(authenticateUser, result, context);
                  }
                }
              }, (close) {});
            },
            icon: Icons.qr_code_scanner,
          ),
          // Bouton pour scanner le QR Code
        ],
      ),
    );
  }

  Future<void> _authenticate(
      Future<void> Function(String userId) authenticateUser,
      String result,
      BuildContext context) async {
    try {
      await authenticateUser(result);
      if (context.mounted) {
        Navigator.of(context).pushReplacementNamed("/events");
      }
    } catch (e) {
      if (context.mounted) GouelException.inform(e, context);
    }
  }

  void authenticate(String id) {}
}
