import 'package:flutter/material.dart';
import 'package:gouel/models/button_model.dart';
import 'package:gouel/models/event_model.dart';
import 'package:gouel/services/gouel_api_service.dart';
import 'package:gouel/services/gouel_session_service.dart';
import 'package:gouel/widgets/event_button.dart';
import 'package:gouel/widgets/gouel_scaffold.dart';
import 'package:provider/provider.dart';

class ManageEventScreen extends StatelessWidget {
  final List<EventButtonModel> buttons = [
    EventButtonModel(
        color: Colors.blue,
        title: 'Buvette',
        path: '/event/buvette',
        permission: 'buvette',
        icon: Icons.sports_bar),
    EventButtonModel(
        color: Colors.blue,
        title: 'Caisse',
        path: '/event/caisse',
        permission: 'caisse',
        icon: Icons.payments),
    EventButtonModel(
        color: Colors.blue,
        title: 'Entrée',
        path: '/event/entree',
        permission: 'entree',
        icon: Icons.door_back_door),
    EventButtonModel(
        color: Colors.blue,
        title: 'Vestiaire',
        path: '/event/vestiaire',
        permission: 'vestiaire',
        icon: Icons.inventory_2),
    EventButtonModel(
        title: 'Paramètres', path: '/settings', icon: Icons.settings),
    EventButtonModel(
        title: 'Crédits', path: '/credits', icon: Icons.info_sharp),
    EventButtonModel(
        color: Colors.blueGrey,
        title: 'Déconnexion',
        path: '/logout',
        icon: Icons.logout),
  ];

  ManageEventScreen({super.key});

  List<String> getPermissions(String eventId) {
    List<String> allPermission = ["buvette", "vestiaire", "caisse", "entree"];
    return allPermission;
  }

  @override
  Widget build(BuildContext context) {
    final Event event = GouelSession().retrieve("event") as Event;

    List<EventButtonModel> permittedButtons = buttons.where((button) {
      return getPermissions(event.id).contains(button.permission) ||
          button.permission == null;
    }).toList();

    return GouelScaffold(
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(20),
            child: Text(
              event.title,
              style: const TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
              textAlign: TextAlign.center,
            ),
          ),
          Expanded(
              child: Column(
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            children: permittedButtons
                .map((e) => EventButton(
                    color: e.color,
                    title: e.title,
                    onTap: () {
                      _handleButtonTap(context, e);
                    },
                    icon: e.icon))
                .toList(),
          )),
        ],
      ),
    );
  }

  void _handleButtonTap(BuildContext context, EventButtonModel button) {
    if (button.path == '/logout') {
      Provider.of<GouelApiService>(context, listen: false)
          .logout(buildContext: context);
    } else {
      // Naviguer vers la page spécifiée
      Navigator.of(context).pushNamed(button.path);
    }
  }
}
