import 'package:flutter/material.dart';
import 'package:gouel/widgets/gouel_bottom_sheet.dart';
import 'package:gouel/services/gouel_storage_service.dart';
import 'package:gouel/widgets/gouel_button.dart'; // Importez GouelStorageService

class GouelServerChooser extends StatefulWidget {
  const GouelServerChooser({super.key});

  @override
  GouelServerChooserState createState() => GouelServerChooserState();
}

class GouelServerChooserState extends State<GouelServerChooser> {
  final TextEditingController _controller = TextEditingController();
  final GouelStorage _storageService =
      GouelStorage(); // Instanciez GouelStorage

  @override
  void initState() {
    super.initState();
    _loadExistingAddress();
  }

  void _loadExistingAddress() async {
    String? existingAddress = await _storageService.retrieve("server_addr");
    if (existingAddress != null) {
      setState(() {
        _controller.text = existingAddress;
      });
    }
  }

  void _saveServerAddress() async {
    String address = _controller.text;
    await _storageService.store("server_addr", address);
  }

  @override
  Widget build(BuildContext context) {
    return GouelBottomSheet(
      title: "Choix du serveur",
      child: Column(
        children: <Widget>[
          TextField(
            controller: _controller,
            decoration: const InputDecoration(
              labelText: 'Adresse',
              border: OutlineInputBorder(),
            ),
          ),
          const SizedBox(
            height: 8,
          ),
          Directionality(
              textDirection: TextDirection.rtl,
              child: GouelButton(
                text: "Valider",
                icon: Icons.check,
                onTap: () {
                  _saveServerAddress();
                  Navigator.of(context).pop();
                },
              )),
        ],
      ),
    );
  }

  @override
  void dispose() {
    _controller.dispose(); // N'oubliez pas de nettoyer le controller
    super.dispose();
  }
}
