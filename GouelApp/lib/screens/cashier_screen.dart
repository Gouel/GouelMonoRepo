import 'package:flutter/material.dart';
import 'package:gouel/models/ticket_model.dart';
import 'package:gouel/services/qr_scanner_service.dart';
import 'package:gouel/utils/gouel_getter.dart';
import 'package:gouel/widgets/gouel_bottom_sheet.dart';
import 'package:gouel/widgets/gouel_button.dart';
import 'package:gouel/widgets/gouel_modal.dart';
import 'package:gouel/widgets/gouel_scaffold.dart';
import 'package:gouel/widgets/gouel_step_builder.dart';
import 'package:gouel/widgets/numeric_keypad.dart';
import 'package:gouel/widgets/paragraph.dart';

class CashierScreen extends StatefulWidget {
  const CashierScreen({super.key});

  @override
  CashierScreenState createState() => CashierScreenState();
}

class CashierScreenState extends State<CashierScreen> {
  String _amount = '';
  PaymentMethod? _selectedPaymentMethod = PaymentMethod.card;

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
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: PaymentMethod.values.map((method) {
        return Padding(
          padding: const EdgeInsets.all(4.0),
          child: GouelButton(
            text: method.name,
            icon: _getPaymentMethodIcon(method),
            onTap: () => _selectPaymentMethod(method),
            color: _selectedPaymentMethod == method ? Colors.blue : Colors.grey,
          ),
        );
      }).toList(),
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

  void _finalizeTransaction() {
    // Créer l'objet Transaction
    // Envoyer l'objet Transaction à l'API
    // Gérer la réponse et informer l'utilisateur du résultat
  }

  void _handleNewTicket() {
    GouelBottomSheet.launch(
      context: context,
      bottomSheet: GouelBottomSheet(
          title: 'Nouveau Ticket',
          child: GouelStepBuilder(
            onValidate: (form) {
              print(form);
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
    formData["validate"] = false;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Paragraph(
          type: ParagraphType.text,
          content: "Nom : ${formData['nom']}",
        ),
        Paragraph.space(),
        Paragraph(
          type: ParagraphType.text,
          content: "Prénom : ${formData['prenom']}",
        ),
        Paragraph.space(),
        Paragraph(
          type: ParagraphType.text,
          content: "email : ${formData['email']}",
        ),
        Paragraph.space(),
        Paragraph(
          type: ParagraphType.text,
          content: "Date de naissance : ${formData['dateDeNaissance']}",
        ),
        Paragraph.space(),
      ],
    );
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

        // XXX Créer transaction

        if (context.mounted) {
          GouelModal.show(context, child: Text(userID));
        }

        _resetAmount();
      },
      (close) {},
    );
  }

  void _resetAmount() {
    setState(() {
      _amount = '';
    });
  }

  IconData _getPaymentMethodIcon(PaymentMethod method) {
    switch (method) {
      case PaymentMethod.cash:
        return Icons.attach_money;
      case PaymentMethod.card:
        return Icons.credit_card;
      default:
        return Icons.money_off;
    }
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

enum PaymentMethod {
  cash,
  card,
}

extension PaymentMethodString on PaymentMethod {
  String get name {
    switch (this) {
      case PaymentMethod.card:
        return "Carte bleue";
      case PaymentMethod.cash:
        return "Espèces";
      default:
        return "";
    }
  }
}
