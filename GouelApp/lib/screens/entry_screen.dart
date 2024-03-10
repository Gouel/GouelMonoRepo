import 'dart:async';

import 'package:flutter/material.dart';
import 'package:gouel/models/ticket_model.dart';
import 'package:gouel/services/gouel_api_service.dart';
import 'package:gouel/services/gouel_session_service.dart';
import 'package:gouel/services/gouel_storage_service.dart';
import 'package:gouel/utils/gouel_getter.dart';
import 'package:gouel/widgets/gouel_bottom_sheet.dart';
import 'package:gouel/services/qr_scanner_service.dart';
import 'package:gouel/widgets/gouel_button.dart';
import 'package:gouel/widgets/gouel_modal.dart';
import 'package:gouel/widgets/gouel_scaffold.dart';
import 'package:gouel/widgets/paragraph.dart';
import 'package:gouel/widgets/settings_field.dart';
import 'package:provider/provider.dart';

class EntryScreen extends StatefulWidget {
  const EntryScreen({super.key});

  @override
  EntryScreenState createState() => EntryScreenState();
}

class EntryScreenState extends State<EntryScreen> {
  String filterTickets = "";

  List<TicketInfos> tickets = [];
  Timer? _timer;

  Future<void> _startAutoReload() async {
    int duration = await GouelStorage().retrieve("data_refresh") ?? 10;
    _timer = Timer.periodic(
        Duration(seconds: duration), (Timer t) => _loadTickets());
  }

  @override
  void initState() {
    super.initState();
    _loadTickets();
    _startAutoReload();
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    List<TicketInfos> filteredTickets = tickets.where((ticket) {
      var email = ticket.user["Email"] as String;
      var nom = ticket.user["LastName"] as String;
      var prenom = ticket.user["FirstName"] as String;
      List<bool> allows = [
        email.toLowerCase().contains(filterTickets.toLowerCase()),
        nom.toLowerCase().contains(filterTickets.toLowerCase()),
        prenom.toLowerCase().contains(filterTickets.toLowerCase()),
      ];
      return allows.contains(true);
    }).toList();

    filteredTickets.sort(
      (a, b) {
        var A = "${a.user["LastName"]} ${a.user["FirstName"]}";
        var B = "${b.user["LastName"]} ${b.user["FirstName"]}";
        return A.compareTo(B);
      },
    );

    return GouelScaffold(
      appBar: AppBar(
        title: const Text("Entrée"),
      ),
      body: Column(
        children: [
          Row(
            mainAxisSize: MainAxisSize.max,
            children: [
              Flexible(
                child: SettingsField(
                    type: SettingsFieldType.inputText,
                    label: "Rechercher un ticket",
                    value: GouelSession().retrieve("entry_filter_email") ?? "",
                    onFinish: (value) {
                      GouelSession().store("entry_filter_email", value);
                      setState(() {
                        filterTickets = value;
                      });
                    }),
              ),
              const SizedBox(
                width: 8,
              ),
              GouelButton(
                text: null,
                onTap: _loadTickets,
                icon: Icons.refresh,
              ),
            ],
          ),
          Paragraph.space(),
          Expanded(
            child: ListView.builder(
              itemCount: filteredTickets.length,
              itemBuilder: (innerContext, index) {
                TicketInfos ticket = filteredTickets[index];
                return Container(
                    margin: const EdgeInsets.symmetric(vertical: 8),
                    child: _buildTicket(ticket, index + 1, context));
              },
            ),
          ),
          Paragraph.space(),
          GouelButton(
            text: "Valider un ticket",
            onTap: _qrValidateTicket,
            icon: Icons.qr_code,
          ),
          GouelButton(
            text: "Rendu EcoCup",
            onTap: _qrValidateEcoCup,
            icon: Icons.local_drink,
          ),
        ],
      ),
    );
  }

  Widget _buildTicket(TicketInfos ticket, int index, BuildContext context) {
    return ClipRRect(
      borderRadius: BorderRadius.circular(8),
      child: Container(
        color: Colors.deepPurple.shade800,
        child: Material(
          color: Colors.transparent,
          child: InkWell(
              onTap: () => _showTicketInfo(ticket, context),
              child: Row(
                mainAxisSize: MainAxisSize.max,
                children: [
                  Container(
                    margin: const EdgeInsets.all(16),
                    padding: const EdgeInsets.all(8),
                    child: Text("$index"),
                  ),
                  Expanded(
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.start,
                      children: [
                        Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                                "${ticket.user['LastName']} ${ticket.user['FirstName']}"),
                            Text("${ticket.user['Email']}")
                          ],
                        ),
                      ],
                    ),
                  ),
                  Container(
                      padding: const EdgeInsets.all(8),
                      child: ticket.isUsed
                          ? const Icon(
                              Icons.check,
                              size: 40,
                              color: Colors.green,
                            )
                          : const SizedBox(
                              width: 40,
                              height: 40,
                            )),
                ],
              )),
        ),
      ),
    );
  }

  Future<Widget> getValidation(String ticketID) async {
    ValidateState state =
        await Provider.of<GouelApiService>(context, listen: false)
            .validateTicket(ticketID);

    if (!context.mounted) return const Text("ERREUR");
    TicketInfos? ticketInfos =
        await getTicketInfos(context, ticketID, withSnackBar: false);
    _loadTickets();
    if (context.mounted) {
      return _buildTicketStateWidget(state, ticketInfos);
    }
    return const SizedBox.shrink();
  }

  Widget _buildTicketStateWidget(ValidateState state, TicketInfos? ticketInfo) {
    IconData icon;
    String title;
    Color iconColor;
    String message = "";
    String nameSurname = "";
    String ageStatus = "";
    List<Widget> more = [];

    if (state == ValidateState.invalid || ticketInfo == null) {
      icon = Icons.error;
      title = 'Ticket Invalide';
      iconColor = Colors.red;
      message = "Ce ticket n'existe pas / n'est pas bon";
    } else {
      nameSurname =
          "${(ticketInfo.user['LastName'] as String).toUpperCase()} ${ticketInfo.user['FirstName']}";
      ageStatus = getMajeurMineur(ticketInfo.user["DOB"] as String);

      if (state == ValidateState.alreadyValidated) {
        icon = Icons.warning;
        title = 'Ticket déjà utilisé';
        iconColor = Colors.orange;
      } else {
        icon = Icons.check_circle;
        title = 'Ticket Valide';
        iconColor = Colors.green;

        if (ageStatus == "Majeur") {
          more.addAll([
            Paragraph.space(),
            GouelButton(
              text: "Rendre SAM",
              onTap: () async {
                bool done =
                    await Provider.of<GouelApiService>(context, listen: false)
                        .setTicketSAM(ticketInfo.id);

                if (context.mounted) {
                  Navigator.of(context, rootNavigator: true).pop('dialog');
                }

                if (context.mounted) {
                  if (done) {
                    GouelModal.show(context,
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.center,
                          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                          children: [
                            const Icon(
                              Icons.check,
                              color: Colors.green,
                              size: 60,
                            ),
                            Paragraph.space(),
                            const Paragraph(
                              type: ParagraphType.text,
                              content:
                                  "L'utilisateur a bien été\ndésigné comme SAM",
                            )
                          ],
                        ));
                  } else {
                    GouelModal.show(context,
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.center,
                          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                          children: [
                            const Icon(
                              Icons.error,
                              color: Colors.red,
                              size: 60,
                            ),
                            Paragraph.space(),
                            const Paragraph(
                              type: ParagraphType.text,
                              content:
                                  "L'utilisateur n'a pas pu être désigné comme SAM",
                            )
                          ],
                        ));
                  }
                }

                _loadTickets();
              },
              icon: Icons.no_drinks,
              color: Colors.blue,
            ),
            Paragraph.space(),
            GouelButton(
                text: "OK",
                onTap: () {
                  Navigator.of(context, rootNavigator: true).pop('dialog');
                }),
          ]);
        }
      }
    }

    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(icon, color: iconColor, size: 60),
        Text(
          title,
          style: Theme.of(context).textTheme.titleLarge,
        ),
        if (state != ValidateState.invalid) ...[
          const SizedBox(height: 20),
          Text(
            nameSurname,
            style: Theme.of(context).textTheme.titleMedium,
            textAlign: TextAlign.center,
          ),
          Text(
            ageStatus,
            style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                color: ageStatus == "Mineur"
                    ? Colors.redAccent
                    : Colors.greenAccent),
            textAlign: TextAlign.center,
          ),
          ...more,
        ] else ...[
          Paragraph(
            type: ParagraphType.text,
            content: message,
          ),
          ...more,
        ]
      ],
    );
  }

  void _validateTicket(String ticketID) {
    GouelModal.showFuture(
      context,
      futureChild: getValidation(ticketID),
    );
  }

  void _validateEcoCup(String ticketID) {
    GouelModal.showFuture(
      context,
      futureChild: getEcoCupValidation(ticketID),
    );
  }

  void _qrValidateTicket() {
    QRScannerService().scanQR(
      context,
      "Scanner ticket",
      (result) async {
        if (context.mounted) _validateTicket(result);
      },
      (close) => null,
    );
  }

  void _qrValidateEcoCup() {
    QRScannerService().scanQR(
      context,
      "Scanner ticket",
      (result) async {
        if (context.mounted) _validateEcoCup(result);
      },
      (close) => null,
    );
  }

  Future<Widget> getEcoCupValidation(String ticketID) async {
    ValidateState state =
        await Provider.of<GouelApiService>(context, listen: false)
            .getEcoCup(ticketID);

    if (context.mounted) {
      String message = "";
      IconData icon;
      Color iconColor;

      switch (state) {
        case ValidateState.ok:
          icon = Icons.check;
          message = "L'EcoCup a bien été rendu";
          iconColor = Colors.green;
          break;
        case ValidateState.alreadyValidated:
          icon = Icons.warning;
          message = "L'EcoCup a déjà été rendu";
          iconColor = Colors.orange;
          break;
        default:
          icon = Icons.error;
          message = "Ce ticket n'existe pas / n'est pas bon";
          iconColor = Colors.red;
          break;
      }

      return Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, color: iconColor, size: 60),
          Text(
            message,
            style: Theme.of(context).textTheme.titleLarge,
          ),
          Paragraph.space(),
          GouelButton(
              text: "OK",
              onTap: () {
                Navigator.of(context, rootNavigator: true).pop('dialog');
              }),
        ],
      );
    } else {
      return const Text("ERREUR");
    }
  }

  void _showTicketInfo(TicketInfos ticketInfos, BuildContext context) async {
    String nom = ticketInfos.user['LastName'] as String;
    String prenom = ticketInfos.user['FirstName'] as String;
    String email = ticketInfos.user['Email'] as String;
    String? dob = ticketInfos.user['DOB'] as String?;
    bool isSam = ticketInfos.isSam;

    if (dob != null) {
      dob = dob.split("-").reversed.join(" / ");
    }

    showModalBottomSheet(
        context: context,
        isScrollControlled: true,
        builder: ((innerContext) => GouelBottomSheet(
            title: "Ticket de ${nom.toUpperCase()} $prenom",
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Paragraph.space(),
                const Paragraph(
                  type: ParagraphType.heading,
                  content: "Informations du ticket",
                ),
                Paragraph.space(),
                Paragraph(
                  type: ParagraphType.text,
                  content: "Nom : $nom",
                ),
                Paragraph(
                  type: ParagraphType.text,
                  content: "Prénom : $prenom",
                ),
                Paragraph(
                  type: ParagraphType.text,
                  content: "Email : $email",
                ),
                if (dob != null)
                  Row(
                    children: [
                      Paragraph(
                        type: ParagraphType.text,
                        content: "Date de naissance : $dob",
                      ),
                      const SizedBox(
                        width: 2,
                      ),
                      Paragraph(
                        type: ParagraphType.hint,
                        content:
                            getMajeurMineur(ticketInfos.user['DOB'] as String),
                      ),
                    ],
                  ),
                if (isSam) ...[
                  Paragraph.space(),
                  GouelButton(
                      text: "Enlever SAM",
                      color: Colors.red,
                      icon: Icons.no_drinks,
                      onTap: () async {
                        bool done = await Provider.of<GouelApiService>(context,
                                listen: false)
                            .setTicketSAM(ticketInfos.id, isSAM: false);
                        if (context.mounted) {
                          Navigator.of(context).pop();
                          if (done) {
                            GouelModal.show(context,
                                child: Column(
                                  mainAxisAlignment:
                                      MainAxisAlignment.spaceEvenly,
                                  crossAxisAlignment: CrossAxisAlignment.center,
                                  children: [
                                    const Icon(
                                      Icons.check,
                                      color: Colors.green,
                                      size: 60,
                                    ),
                                    Paragraph.space(),
                                    const Paragraph(
                                      type: ParagraphType.text,
                                      content: "L'utilisateur n'est plus SAM",
                                    )
                                  ],
                                ));
                          } else {
                            GouelModal.show(context,
                                child: Column(
                                  mainAxisAlignment:
                                      MainAxisAlignment.spaceEvenly,
                                  crossAxisAlignment: CrossAxisAlignment.center,
                                  children: [
                                    const Icon(
                                      Icons.error,
                                      color: Colors.red,
                                      size: 60,
                                    ),
                                    Paragraph.space(),
                                    const Paragraph(
                                      type: ParagraphType.text,
                                      content:
                                          "Une erreur est survenue. Veuillez réessayer",
                                    )
                                  ],
                                ));
                          }
                        }

                        _loadTickets();
                      })
                ],
                Paragraph.space(),
                if (!ticketInfos.isUsed)
                  GouelButton(
                      text: "Valider le ticket",
                      color: Colors.green,
                      onTap: () async {
                        if (innerContext.mounted) {
                          Navigator.of(innerContext).pop();
                          _validateTicket(ticketInfos.id);
                        }
                      }),
              ],
            ))));
  }

  Future<void> _loadTickets() async {
    List<TicketInfos> providedTickets =
        await Provider.of<GouelApiService>(context, listen: false)
            .getAllTicketInfos(context);
    if (!mounted) return;

    setState(() {
      tickets = providedTickets;
    });
  }
}
