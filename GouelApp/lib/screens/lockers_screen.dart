import 'dart:async';

import 'package:flutter/material.dart';
import 'package:gouel/models/event_model.dart';
import 'package:gouel/models/ticket_model.dart';
import 'package:gouel/services/gouel_api_service.dart';
import 'package:gouel/services/gouel_session_service.dart';
import 'package:gouel/services/gouel_storage_service.dart';
import 'package:gouel/utils/gouel_getter.dart';
import 'package:gouel/widgets/gouel_bottom_sheet.dart';
import 'package:gouel/services/qr_scanner_service.dart';
import 'package:gouel/widgets/gouel_button.dart';
import 'package:gouel/widgets/gouel_scaffold.dart';
import 'package:gouel/widgets/paragraph.dart';
import 'package:provider/provider.dart';

class LockersScreen extends StatefulWidget {
  const LockersScreen({super.key});

  @override
  LockersScreenState createState() => LockersScreenState();
}

class LockersScreenState extends State<LockersScreen> {
  bool showTakenLockers = true;
  String? filterLockers;

  List<Locker> lockers = [];

  Timer? _timer;

  @override
  void initState() {
    super.initState();
    _loadLockers();
    _startAutoReload();
  }

  @override
  void dispose() {
    _timer?.cancel(); // Annuler le timer lors de la destruction de l'état
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    List<Locker> filteredLockers = showTakenLockers
        ? lockers
        : lockers.where((locker) => locker.userId == '').toList();

    if (filterLockers != null) {
      filteredLockers = filteredLockers
          .where((locker) => locker.userId == filterLockers)
          .toList();
    }

    return WillPopScope(
      onWillPop: () async {
        if (filterLockers != null) {
          // Si le filtre est actif, retirez-le et empêchez la navigation arrière
          setState(() {
            filterLockers = null;
          });
          return false; // Empêche le retour à la page précédente
        }
        return true; // Permet la navigation arrière si aucun filtre n'est actif
      },
      child: GouelScaffold(
        appBar: AppBar(
          title: const Text("Vestiaires"),
        ),
        floatingActionButton: FloatingActionButton(
          onPressed: () => _showMenu(context),
          child: const Icon(Icons.menu),
        ),
        body: Column(
          children: [
            Expanded(
              child: GridView.builder(
                gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                  crossAxisCount: 3,
                  childAspectRatio: 1.0,
                ),
                itemCount: filteredLockers.length,
                itemBuilder: (context, index) {
                  var locker = filteredLockers[index];
                  bool isTaken = locker.userId != '';
                  if (!showTakenLockers && isTaken) {
                    return const SizedBox.shrink();
                  }
                  return Container(
                      margin: const EdgeInsets.all(4),
                      child: _buildLockerItem(locker));
                },
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _scanQrCodeFilterLockers() {
    QRScannerService().scanQR(
      context,
      "Scanner ticket",
      (result) {
        setState(() {
          filterLockers = result;
        });
        Navigator.of(context).pop();
      },
      (close) => null,
    );
  }

  void _assignLocker(Locker locker) {
    QRScannerService().scanQR(
      context,
      "Scanner ticket",
      (result) async {
        // On vérifie si l'utilisateur existe bien
        TicketInfos? ticketInfos = await getTicketInfos(context, result);
        if (context.mounted && ticketInfos != null) {
          await Provider.of<GouelApiService>(context, listen: false)
              .setEventLocker(
                  context,
                  Locker(
                      lockerCode: locker.lockerCode, userId: ticketInfos.id));

          _loadLockers();
        }
      },
      (close) => null,
    );
  }

  void _showLockerInfo(Locker locker) async {
    TicketInfos? ticketsInfos = await getTicketInfos(context, locker.userId);
    if (!(context.mounted && ticketsInfos != null)) {
      return;
    }

    showModalBottomSheet(
        isScrollControlled: true,
        context: context,
        builder: ((context) => GouelBottomSheet(
            title: "Vestiaire ${locker.lockerCode}",
            child: Column(
              children: [
                Paragraph.space(),
                const Paragraph(
                  type: ParagraphType.heading,
                  content: "Informations du propriétaire",
                ),
                Paragraph.space(),
                Paragraph(
                  type: ParagraphType.text,
                  content: "Nom : ${ticketsInfos.user['LastName']}",
                ),
                Paragraph(
                  type: ParagraphType.text,
                  content: "Prénom : ${ticketsInfos.user['FirstName']}",
                ),
                Paragraph(
                  type: ParagraphType.text,
                  content: "Email : ${ticketsInfos.user['Email']}",
                ),
                Paragraph.space(),
                GouelButton(
                    text: "Libérer le vestiaire",
                    color: Colors.red,
                    onTap: () async {
                      await Provider.of<GouelApiService>(context, listen: false)
                          .setEventLocker(
                              context,
                              Locker(
                                  lockerCode: locker.lockerCode, userId: ""));

                      _loadLockers();
                      if (context.mounted) {
                        Navigator.of(context).pop();
                      }
                    })
              ],
            ))));
  }

  Future<void> _loadLockers() async {
    final Event event = GouelSession().retrieve("event") as Event;
    List<Locker> providedLockers =
        await Provider.of<GouelApiService>(context, listen: false)
            .getEventLockers(context, event.id);
    setState(() {
      lockers = providedLockers;
    });
  }

  Widget _buildLockerItem(Locker locker) {
    bool isTaken = locker.userId != '';
    return ClipRRect(
      borderRadius: BorderRadius.circular(16), // Rayon de bordure arrondie
      child: Container(
        color: isTaken ? Colors.red.shade500 : Colors.green.shade500,
        child: Material(
          color: Colors.transparent,
          child: InkWell(
            onTap: () =>
                isTaken ? _showLockerInfo(locker) : _assignLocker(locker),
            child: Center(
              child: Text(locker.lockerCode,
                  style: Theme.of(context).textTheme.titleLarge),
            ),
          ),
        ),
      ),
    );
  }

  Future<void> _startAutoReload() async {
    int duration = await GouelStorage().retrieve("data_refresh") ?? 10;
    _timer = Timer.periodic(
        Duration(seconds: duration), (Timer t) => _loadLockers());
  }

  void _showMenu(BuildContext context) {
    showModalBottomSheet(
      isScrollControlled: true,
      context: context,
      builder: (BuildContext context) {
        return GouelBottomSheet(
            title: 'Options',
            child: Column(
              children: <Widget>[
                GouelButton(
                  text: "Trouver un vestiaire",
                  onTap: _scanQrCodeFilterLockers,
                  icon: Icons.qr_code_scanner,
                ),
                if (filterLockers != null)
                  GouelButton(
                      color: Colors.red,
                      text: "Supprimer le filtre",
                      onTap: () {
                        setState(() {
                          filterLockers = null;
                        });
                        Navigator.of(context).pop();
                      }),
                GouelButton(
                  text: showTakenLockers
                      ? "Masquer vestiaires utilisés"
                      : "Afficher tous les vestiaires",
                  onTap: () {
                    setState(() {
                      showTakenLockers = !showTakenLockers;
                    });
                    Navigator.pop(context); // Fermer le BottomSheet
                  },
                  icon: showTakenLockers
                      ? Icons.visibility_off
                      : Icons.visibility,
                ),
              ],
            ));
      },
    );
  }
}
