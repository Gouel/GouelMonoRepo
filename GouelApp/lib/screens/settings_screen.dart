import 'package:flutter/material.dart';
import 'package:gouel/services/gouel_storage_service.dart';
import 'package:gouel/widgets/gouel_button.dart';
import 'package:gouel/widgets/gouel_scaffold.dart';
import 'package:gouel/widgets/gouel_snackbar.dart';
import 'package:gouel/widgets/paragraph.dart';
import 'package:gouel/widgets/settings_field.dart';

class SettingsScreen extends StatelessWidget {
  SettingsScreen({Key? key}) : super(key: key);

  final Map<String, dynamic> _settingsData = {};

  void _handleSettingChange(String key, dynamic value) {
    _settingsData[key] = value;
  }

  Future<void> saveSettings(BuildContext context) async {
    final storage = GouelStorage();

    for (var key in _settingsData.keys) {
      await storage.store(key, _settingsData[key]);
    }
    if (context.mounted) {
      showGouelSnackbar(
          context, "Les paramètres ont bien été sauvegardés.", Colors.green,
          duration: 2);
    }
  }

  Future<Map<String, dynamic>> getSettings(List<String> keys) async {
    final storage = GouelStorage();
    Map<String, dynamic> settings = {};
    for (String key in keys) {
      dynamic value = await storage.retrieve(key);
      settings[key] = value;
    }

    return settings;
  }

  List<Widget> getContent(BuildContext context, Map<String, dynamic> settings) {
    _settingsData.clear();
    _settingsData.addAll(settings);

    return [
      const Paragraph(
        type: ParagraphType.heading,
        content: "Paramètres généraux",
      ),
      SettingsField(
        type: SettingsFieldType.inputText,
        label: "Adresse du serveur",
        value: _settingsData['server_addr'] ?? '',
        onChanged: (value) => _handleSettingChange("server_addr", value),
      ),
      SettingsField(
        type: SettingsFieldType.inputNumeric,
        label: "Taux de rafraîchissement* (secondes)",
        value: _settingsData['data_refresh'] ?? 10,
        onChanged: (value) => _handleSettingChange("data_refresh", value),
      ),
      const Paragraph(
        type: ParagraphType.hint,
        content:
            "* Le taux de rafraîchissement correspond au temps que prend la page Vestiaire / Entrée pour s'actualiser",
      ),
      Paragraph.space(),

      SettingsField(
        type: SettingsFieldType.switchField,
        label: "Afficher le titre des produits",
        value: _settingsData['product_show_title'] ?? false,
        onChanged: (value) => _handleSettingChange("product_show_title", value),
      ),

      Paragraph.space(),
      // Ajoutez d'autres champs de paramètres ici en utilisant _settingsData
      GouelButton(
        text: "Sauvegarder",
        onTap: () {
          saveSettings(context);
        },
        icon: Icons.save,
      )
    ];
  }

  @override
  Widget build(BuildContext context) {
    return GouelScaffold(
        appBar: AppBar(title: const Text("Paramètres")),
        body: FutureBuilder(
          future: getSettings(
              ["server_addr", "data_refresh", "product_show_title"]),
          builder: (context, snapshot) {
            if (snapshot.connectionState == ConnectionState.waiting) {
              return const Center(child: CircularProgressIndicator());
            }

            if (snapshot.hasError) {
              return const Center(
                  child: Text("Erreur lors du chargement des paramètres"));
            }

            final settings = snapshot.data ?? {};
            final content = getContent(context, settings);
            return ListView.separated(
              itemBuilder: (context, index) => content[index],
              separatorBuilder: (context, index) => Paragraph.space(),
              itemCount: content.length,
            );
          },
        ));
  }
}
