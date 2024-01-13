import 'package:flutter/material.dart';
import 'package:gouel/widgets/gouel_scaffold.dart';
import 'package:gouel/widgets/paragraph.dart';

class CreditsScreen extends StatelessWidget {
  const CreditsScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return GouelScaffold(
      appBar: AppBar(
        title: const Text("Crédits"),
      ),
      body: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: <Widget>[
              const Paragraph(
                type: ParagraphType.heading,
                content: 'Développement de l\'Application',
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.bulletList,
                items: ["Matthias HARTMANN : Développeur FullStack"],
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.heading,
                content: 'Technologies et Outils',
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.bulletList,
                items: [
                  'Frontend: Flask, connecté à HelloAsso pour la gestion des paiements',
                  'Backend: Golang avec une base de données MongoDB',
                  'Application mobile: Développée en Flutter',
                ],
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.heading,
                content: 'Partenariat',
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.text,
                content: "InterAsso de l'IUT de Lannion :",
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.bulletList,
                items: [
                  'Aide au développement',
                  'Aide matériel (serveur temporaire, ...)',
                ],
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.heading,
                content: 'Remerciements',
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.text,
                content: 'Merci à tous ceux qui ont contribué au projet.',
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.heading,
                content: 'Licence',
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.text,
                content:
                    'Distribué sous la Licence Mozilla Public License version 2.0.',
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.heading,
                content: 'Liens et Références',
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.link,
                content: 'https://gouel.fr',
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.heading,
                content: 'Note',
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.text,
                content:
                    'La fonctionnalité NFC/RFID n\'est plus utilisée dans l\'application Gouel.',
              ),
            ],
          ),
        ),
      ),
    );
  }
}
