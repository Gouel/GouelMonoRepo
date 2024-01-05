import 'package:flutter/material.dart';
import 'package:gouel/utils/gouel_exception.dart';
import 'package:url_launcher/url_launcher.dart';

enum ParagraphType {
  heading,
  text,
  hint,
  link,
  bulletList,
  numberedList,
  image,
  other,
}

class Paragraph extends StatelessWidget {
  final ParagraphType type;
  final String content;
  final List<String>? items;
  final Image? image;

  const Paragraph({
    Key? key,
    required this.type,
    this.content = '',
    this.items,
    this.image,
  }) : super(key: key);

  static space() {
    return const Paragraph(type: ParagraphType.other);
  }

  @override
  Widget build(BuildContext context) {
    switch (type) {
      case ParagraphType.text:
        return Text(
          content,
          style: const TextStyle(fontSize: 16),
          textAlign: TextAlign.left,
        );
      case ParagraphType.hint:
        return Text(
          content,
          style: const TextStyle(
              fontSize: 12, fontStyle: FontStyle.italic, color: Colors.white60),
          textAlign: TextAlign.left,
        );
      case ParagraphType.heading:
        return Text(
          content,
          style: Theme.of(context).textTheme.titleLarge,
          textAlign: TextAlign.left,
        );
      case ParagraphType.link:
        return InkWell(
          child: Text(
            content,
            style: const TextStyle(fontSize: 16, color: Colors.blue),
            textAlign: TextAlign.left,
          ),
          onTap: () {
            try {
              _launchURL(content);
            } catch (e) {
              if (e is GouelException) {}
            }
          },
        );
      case ParagraphType.bulletList:
        return Column(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: items!
              .map((item) => Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text('• ', style: TextStyle(fontSize: 16)),
                      Expanded(
                          child:
                              Text(item, style: const TextStyle(fontSize: 16))),
                    ],
                  ))
              .toList(),
        );
      case ParagraphType.numberedList:
        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: items!
              .asMap()
              .entries
              .map((entry) => Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text('${entry.key + 1}. ',
                          style: const TextStyle(fontSize: 16)),
                      Expanded(
                          child: Text(entry.value,
                              style: const TextStyle(fontSize: 16))),
                    ],
                  ))
              .toList(),
        );
      case ParagraphType.image:
        return image != null ? image! : const SizedBox();
      default:
        return const SizedBox(
          height: 16,
        );
    }
  }

  void _launchURL(String url) async {
    final Uri uri = Uri.parse(url);

    if (!await launchUrl(uri, mode: LaunchMode.externalApplication)) {
      throw GouelException(
          message: "L'url $url n'a pas pu être ouverte",
          state: GouelExceptionState.warning);
    }
  }
}
