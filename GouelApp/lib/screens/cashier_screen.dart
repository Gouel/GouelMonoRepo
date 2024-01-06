import 'package:flutter/material.dart';
import 'package:gouel/models/ticket_model.dart';
import 'package:gouel/models/transcations_model.dart';
import 'package:gouel/services/gouel_api_service.dart';
import 'package:gouel/services/gouel_session_service.dart';
import 'package:gouel/services/qr_scanner_service.dart';
import 'package:gouel/utils/gouel_getter.dart';
import 'package:gouel/widgets/gouel_bottom_sheet.dart';
import 'package:gouel/widgets/gouel_button.dart';
import 'package:gouel/widgets/gouel_modal.dart';
import 'package:gouel/widgets/gouel_scaffold.dart';
import 'package:gouel/widgets/gouel_step_builder.dart';
import 'package:gouel/widgets/numeric_keypad.dart';
import 'package:gouel/widgets/paragraph.dart';
import 'package:provider/provider.dart';

class CashierScreen extends StatefulWidget {
  const CashierScreen({super.key});

  @override
  CashierScreenState createState() => CashierScreenState();
}

class CashierScreenState extends State<CashierScreen> {
  String _amount = '';
  PaymentMethod? _selectedPaymentMethod = PaymentMethod.carte;

  void _handleKeypadInput(String value) {
    setState(() {
      _amount = value;
    });
  }

  void _selectPaymentMethod(PaymentMethod method) {
    setState(() {
      _selectedPaymentMethod = method;
    });
  }

  Widget _buildPaymentMethodSelector() {
    List<PaymentMethod> availableMethods =
        PaymentMethod.values.where((element) => element.available).toList();

    return GridView.count(
      shrinkWrap: true,
      crossAxisCount: 2,
      childAspectRatio: 3,
      children: availableMethods
          .map((method) => Padding(
                padding: const EdgeInsets.all(4.0),
                child: GouelButton(
                  text: method.desc,
                  icon: method.icon,
                  onTap: () => _selectPaymentMethod(method),
                  color: _selectedPaymentMethod == method
                      ? Colors.blue
                      : Colors.grey,
                ),
              ))
          .toList(),
    );
  }

  @override
  Widget build(BuildContext context) {
    double trueAmount = double.tryParse(_amount) ?? 0;

    return GouelScaffold(
      appBar: AppBar(title: const Text('Caisse')),
      body: Column(
        mainAxisSize: MainAxisSize.max,
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          // Afficher le montant actuel
          _buildAmountDisplay(),
          NumericKeypad(
            onNumberSelected: _handleKeypadInput,
            input: _amount,
          ),
          _buildPaymentMethodSelector(),
          Paragraph.space(),
          Directionality(
            textDirection: TextDirection.rtl,
            child: GouelButton(
              text: "Payer",
              icon: Icons.euro,
              onTap: trueAmount != 0 ? _onPayPressed : () {},
              color: trueAmount != 0 ? Colors.green : Colors.grey,
            ),
          )
          // Autres éléments de l'interface utilisateur
        ],
      ),
    );
  }

  Widget _buildAmountDisplay() {
    return Container(
      width: double.infinity, // Prend toute la largeur
      margin: const EdgeInsets.all(20),
      padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 20),
      decoration: BoxDecoration(
        color: const Color.fromARGB(255, 10, 14, 23), // Couleur de fond
        borderRadius: BorderRadius.circular(10),
      ),
      child: Text(
        _amount == "" ? '0 €' : '$_amount €',
        style: Theme.of(context).textTheme.headlineLarge,
        textAlign: TextAlign.center,
      ),
    );
  }

  void _finalizeTransaction(String userID, Function(bool) onCompletion,
      {double? currentAmount}) async {
    // Créer l'objet Transaction

    double? trueAmount = currentAmount ?? double.tryParse(_amount);
    if (trueAmount == null || trueAmount <= 0) {
      // XXX Afficher problème
      return;
    }

    Transaction transaction = Transaction(
      transactionType: TransactionType.credit,
      dateTime: DateTime.now(),
      eventID: GouelSession().retrieve("event").id,
      cart: [],
      amount: trueAmount,
      paymentMethod: _selectedPaymentMethod,
    );

    // Envoyer l'objet Transaction à l'API

    bool result = await Provider.of<GouelApiService>(context, listen: false)
        .addTransaction(userID, transaction);

    // Gérer la réponse
    onCompletion(result);
  }

  void _handleNewTicket() {
    GouelBottomSheet.launch(
      context: context,
      bottomSheet: GouelBottomSheet(
          title: 'Nouveau Ticket',
          child: GouelStepBuilder(
            onValidate: (form) async {
              // transforme en identite
              Map<String, dynamic> user = {
                "firstName": form["prenom"],
                "lastName": form["nom"],
                "email": form["email"],
                "dob": form["dateDeNaissance"].split("/").reversed.join("-")
              };

              if (!mounted) return;

              String? userId =
                  await Provider.of<GouelApiService>(context, listen: false)
                      .addUser(user);

              if (userId == null) {
                //XXX message !
                return;
              }

              // On génère une transaction
              _finalizeTransaction(userId, (p0) => print(p0));
            },
            steps: [_buildStepOne, _buildStepTwo, _buildStepThree],
          )),
    );
  }

  Widget _buildStepOne(Map<String, dynamic> formData) {
    return Column(
      children: [
        TextFormField(
          key: const Key("nom"),
          decoration: const InputDecoration(labelText: 'Nom'),
          initialValue: formData['nom'] ?? "",
          onChanged: (value) => formData['nom'] = value,
        ),
        const SizedBox(height: 10),
        TextFormField(
          key: const Key("prenom"),
          decoration: const InputDecoration(labelText: 'Prénom'),
          initialValue: formData['prenom'] ?? "",
          onChanged: (value) => formData['prenom'] = value,
        ),
      ],
    );
  }

  Widget _buildStepTwo(Map<String, dynamic> formData) {
    return Column(
      children: [
        TextFormField(
          key: const Key("email"),
          decoration: const InputDecoration(labelText: 'Email'),
          keyboardType: TextInputType.emailAddress,
          initialValue: formData['email'] ?? "",
          onChanged: (value) => formData['email'] = value,
        ),
        const SizedBox(height: 10),
        TextFormField(
          key: const Key("dob"),
          decoration: const InputDecoration(labelText: 'Date de naissance'),
          keyboardType: TextInputType.datetime,
          initialValue: formData['dateDeNaissance'] ?? "",
          onChanged: (value) => formData['dateDeNaissance'] = value,
        ),
      ],
    );
  }

  Widget _buildStepThree(Map<String, dynamic> formData) {
    // Vérification des données
    List<(bool, String)> errors = [];

    errors.addAll([
      _verifyName(formData["nom"] ?? ""),
      _verifyName(formData["prenom"] ?? "", key: "prenom"),
      _verifyEmail(formData["email"] ?? ""),
      _verifyDate(formData["dateDeNaissance"] ?? ""),
    ]);

    errors = errors.where(((e) => !e.$1)).toList();

    formData["validate"] = errors.isEmpty;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Paragraph(
          type: ParagraphType.text,
          content: "Nom : ${formData['nom'] ?? 'Non renseigné'}",
        ),
        Paragraph.space(),
        Paragraph(
          type: ParagraphType.text,
          content: "Prénom : ${formData['prenom'] ?? 'Non renseigné'}",
        ),
        Paragraph.space(),
        Paragraph(
          type: ParagraphType.text,
          content: "Email : ${formData['email'] ?? 'Non renseigné'}",
        ),
        Paragraph.space(),
        Paragraph(
          type: ParagraphType.text,
          content:
              "Date de naissance : ${formData['dateDeNaissance'] ?? 'Non renseigné'}",
        ),
        Paragraph.space(),
        if (errors.isNotEmpty)
          SingleChildScrollView(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Paragraph(
                  type: ParagraphType.text,
                  content: "Erreurs :",
                  color: Colors.red,
                ),
                Paragraph.space(),
                Paragraph(
                  type: ParagraphType.bulletList,
                  items: errors.map((e) => e.$2).toList(),
                  color: Colors.red,
                ),
              ],
            ),
          )
      ],
    );
  }

  (bool, String) _verifyName(String name, {String key = 'nom'}) {
    bool correctName = false;

    List<String> forbiddenNames = [
      'firstname',
      'lastname',
      'unknown',
      'first_name',
      'last_name',
      'anonyme',
      'user',
      'admin',
      'name',
      'nom',
      'prénom',
      'test',
    ];

    if (name.isEmpty) {
      return (correctName, 'Le $key est vide');
    }
    if (name.length > 255) {
      return (correctName, 'Le $key est trop long');
    }
    if (name.contains(RegExp(r'\d'))) {
      return (correctName, 'Le $key ne doit pas contenir de chiffres');
    }
    if (name.length == 1) {
      return (correctName, 'Le $key ne doit pas être un seul caractère');
    }
    if (forbiddenNames.contains(name.toLowerCase())) {
      return (correctName, 'Le $key est interdit');
    }
    if (!RegExp(r'[aeiouyéèêëàâäôöûüç]').hasMatch(name.toLowerCase())) {
      return (correctName, 'Le $key doit contenir au moins une voyelle');
    }
    if (RegExp(r'(.)\1\1').hasMatch(name)) {
      return (
        correctName,
        'Le $key ne doit pas contenir de caractères répétitifs trois fois de suite'
      );
    }

    if (!RegExp(r"^[a-zA-Zéèêëàâäôöûüç\'\- ]+$").hasMatch(name)) {
      return (correctName, 'Le $key contient des caractères non autorisés');
    }

    correctName = true;
    return (correctName, 'Le $key est valide');
  }

  (bool, String) _verifyEmail(String email) {
    if (email.isEmpty) return (false, "L'email est vide");

    String emailRegex =
        r'^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$';

    if (!RegExp(emailRegex).hasMatch(email)) {
      return (false, "L'email n'est pas valide");
    }

    return (true, "L'email est valide");
  }

  (bool, String) _verifyDate(String date) {
    if (date.isEmpty) return (false, "La date est vide");

    if (!RegExp(r'(\d{2}\/){2}\d{4}').hasMatch(date)) {
      return (false, "La date n'est pas valide");
    }

    return (true, "La date est valide");
  }

  void _handleQrCodePayment() {
    // Utiliser le service QRScanner pour scanner un QR Code
    QRScannerService().scanQR(
      context,
      "Scanner QRCode",
      (result) async {
        // Logique après avoir scanné le QR Code
        TicketInfos? ticketInfos = await getTicketInfos(context, result);
        if (ticketInfos == null) return;
        String userID = ticketInfos.user["user_id"];

        _finalizeTransaction(
          userID,
          (result) {
            if (!result) {
              GouelModal.show(
                context,
                child: Column(
                  children: [
                    const Icon(
                      Icons.error,
                      color: Colors.red,
                      size: 60,
                    ),
                    Text(
                      "Le compte n'a pas pu être crédité.",
                      style: Theme.of(context).textTheme.titleLarge,
                    )
                  ],
                ),
              );
            } else {
              GouelModal.show(
                context,
                child: Column(
                  children: [
                    const Icon(
                      Icons.check,
                      color: Colors.green,
                      size: 60,
                    ),
                    Text(
                      "Le compte a été crédité de $_amount €",
                      style: Theme.of(context).textTheme.titleLarge,
                      textAlign: TextAlign.center,
                    )
                  ],
                ),
              );
              _resetAmount();
            }
          },
        );
      },
      (close) {},
    );
  }

  void _resetAmount() {
    setState(() {
      _amount = '';
    });
  }

  void _onPayPressed() {
    GouelBottomSheet.launch(
      context: context,
      bottomSheet: GouelBottomSheet(
        title: 'Options de Paiement',
        child: Padding(
          padding: const EdgeInsets.symmetric(vertical: 8.0),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              GouelButton(
                text: 'Recharger Compte',
                onTap: () {
                  Navigator.pop(context);
                  _handleQrCodePayment();
                },
              ),
              Paragraph.space(),
              GouelButton(
                text: 'Nouveau Ticket',
                onTap: () {
                  Navigator.pop(context);
                  _handleNewTicket();
                },
              ),
            ],
          ),
        ),
      ),
    );
  }
}
