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
  final MaterialColor? color;
  final TextAlign? textAlign;

  const Paragraph({
    Key? key,
    required this.type,
    this.content = '',
    this.items,
    this.image,
    this.color,
    this.textAlign = TextAlign.left,
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
          style: TextStyle(fontSize: 16, color: color),
          textAlign: textAlign,
        );
      case ParagraphType.hint:
        return Text(
          content,
          style: TextStyle(
              fontSize: 12,
              fontStyle: FontStyle.italic,
              color: color ?? Colors.white60),
          textAlign: textAlign,
        );
      case ParagraphType.heading:
        return Text(
          content,
          style: Theme.of(context).textTheme.titleLarge,
          textAlign: textAlign,
        );
      case ParagraphType.link:
        return InkWell(
          child: Text(
            content,
            style: TextStyle(fontSize: 16, color: color ?? Colors.blue),
            textAlign: textAlign,
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
                      Text(
                        '• ',
                        style: TextStyle(fontSize: 16, color: color),
                        textAlign: textAlign,
                      ),
                      Expanded(
                        child: Text(
                          item,
                          style: TextStyle(fontSize: 16, color: color),
                          textAlign: textAlign,
                        ),
                      ),
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
                      Text(
                        '${entry.key + 1}. ',
                        style: TextStyle(fontSize: 16, color: color),
                        textAlign: textAlign,
                      ),
                      Expanded(
                        child: Text(
                          entry.value,
                          style: TextStyle(fontSize: 16, color: color),
                          textAlign: textAlign,
                        ),
                      ),
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
