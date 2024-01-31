import 'package:flutter/material.dart';

class GouelButton extends StatelessWidget {
  final String? text;
  final VoidCallback onTap;
  final VoidCallback? onLongTap;
  final MaterialColor color;
  final IconData? icon;

  const GouelButton({
    Key? key,
    required this.text,
    required this.onTap,
    this.onLongTap,
    this.color = Colors.deepPurple,
    this.icon,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ClipRRect(
      borderRadius: BorderRadius.circular(5),
      child: ElevatedButton(
        onPressed: onTap,
        onLongPress: () => {
          if (onLongTap != null) {onLongTap!()}
        },
        style: ElevatedButton.styleFrom(
          foregroundColor: Colors.white,
          backgroundColor: color.shade600, // Text color
          padding: const EdgeInsets.all(10),
          shape: LinearBorder.bottom(
              side: BorderSide(color: color.shade800, width: 4)),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            if (icon != null)
              Icon(
                icon,
              ),
            if (icon != null && text != null)
              const SizedBox(
                width: 8,
              ),
            if (text != null)
              Text(
                text!,
                style: Theme.of(context).textTheme.titleLarge,
              ),
          ],
        ),
      ),
    );
  }
}
