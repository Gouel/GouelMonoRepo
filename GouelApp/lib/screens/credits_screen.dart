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
                content: 'Partenaires',
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
                type: ParagraphType.text,
                content: "Association des Élèves de l'ENSSAT:",
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.bulletList,
                items: [
                  'Audit de sécurité du logiciel',
                  'Aide au développement',
                ],
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.text,
                content: "Rythmes and blouse (IFSI Lannion)",
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.bulletList,
                items: [
                  'Testes utilisateurs',
                  'Aide au développement',
                ],
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.text,
                content: "Service Jeunesse de la ville de Lannion",
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.bulletList,
                items: [
                  'Mise en relation avec les associations de la ville',
                  'Aide au développement',
                ],
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.heading,
                content: 'Licence',
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.text,
                content: 'Distribué sous la Licence MIT',
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.heading,
                content: 'Liens et Références',
              ),
              Paragraph.space(),
              const Paragraph(
                type: ParagraphType.link,
                content: 'http://gouel.fr',
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
              Paragraph.space(),
              const Paragraph(
                  type: ParagraphType.text,
                  content:
                      'Attention l\'abus d\'alcool est dangereux pour la santé, à consommer avec modération.'),
            ],
          ),
        ),
      ),
    );
  }
}
