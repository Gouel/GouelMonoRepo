import 'package:flutter/material.dart';
import 'package:gouel/widgets/gouel_button.dart';

class EventButton extends StatelessWidget {
  final String title;
  final VoidCallback onTap;
  final IconData icon;
  final MaterialColor color;

  const EventButton(
      {Key? key,
      required this.title,
      required this.onTap,
      required this.icon,
      this.color = Colors.deepPurple})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return GouelButton(
      text: title,
      onTap: onTap,
      icon: icon,
      color: color,
    );
  }
}
