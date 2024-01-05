import 'package:flutter/material.dart';
import 'package:gouel/screens/cashier_screen.dart';
import 'package:gouel/screens/credits_screen.dart';
import 'package:gouel/screens/entry_screen.dart';
import 'package:gouel/screens/event_screen.dart';
import 'package:gouel/screens/lockers_screen.dart';
import 'package:gouel/screens/login_screen.dart';
import 'package:gouel/screens/manage_event_screen.dart';
import 'package:gouel/screens/settings_screen.dart';
import 'package:gouel/services/gouel_api_service.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:provider/provider.dart';
import 'utils/theme.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await Hive.initFlutter();
  await Hive.openBox('gouelStorageBox');

  runApp(MultiProvider(
    providers: [
      ChangeNotifierProvider<GouelApiService>(
        create: (context) => GouelApiService(context),
      ),
    ],
    child: const Gouel(),
  ));
}

class Gouel extends StatelessWidget {
  const Gouel({super.key});

  @override
  Widget build(BuildContext context) {
    Map<String, Widget Function(BuildContext)> routes = {
      // Page de navigation événement
      "/manage_event": (builder) => ManageEventScreen(),
      // Page des droits
      "/event/vestiaire": (builder) => const LockersScreen(),
      "/event/entree": (builder) => const EntryScreen(),
      "/event/caisse": (builder) => const CashierScreen(),

      "/events": (builder) => const EventsScreen(),
      "/credits": (builder) => const CreditsScreen(),
      "/settings": (builder) => SettingsScreen(),
    };

    return MaterialApp(
      title: 'Gouel',
      theme: AppTheme.darkTheme,
      home: const LoginScreen(),
      routes: routes,

      // Configurez les routes ici si nécessaire
    );
  }
}
