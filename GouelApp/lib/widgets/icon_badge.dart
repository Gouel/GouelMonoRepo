import 'package:flutter/material.dart';

class IconBadge extends StatelessWidget {
  final IconData icon;
  final int badgeCount;
  final double size;

  const IconBadge(
      {super.key,
      required this.icon,
      required this.badgeCount,
      this.size = 32});

  @override
  Widget build(BuildContext context) {
    return Stack(
      alignment: Alignment.bottomRight,
      children: [
        Icon(icon, size: size),
        if (badgeCount > 0)
          Container(
            padding: const EdgeInsets.all(2),
            decoration: BoxDecoration(
              color: Colors.red,
              borderRadius: BorderRadius.circular(8),
            ),
            constraints: const BoxConstraints(
              minWidth: 16,
              minHeight: 16,
            ),
            child: Text(
              badgeCount.toString(),
              style: const TextStyle(
                color: Colors.white,
                fontSize: 10,
              ),
              textAlign: TextAlign.center,
            ),
          ),
      ],
    );
  }
}
