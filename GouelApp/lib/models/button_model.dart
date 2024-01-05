import 'package:flutter/material.dart';

class EventButtonModel {
  final String title;
  final String path;
  final String? permission;
  final IconData icon;
  final MaterialColor color;

  EventButtonModel(
      {required this.title,
      required this.path,
      required this.icon,
      this.color = Colors.deepPurple,
      this.permission});
}
