import 'package:flutter/material.dart';

class AppTheme {
  static ThemeData get darkTheme {
    return ThemeData(
        useMaterial3: false,
        brightness: Brightness.dark,
        canvasColor: Colors.transparent,
        colorScheme: ColorScheme.dark(
          primary: Colors.deepPurple.shade600, // purple-600
          secondary: Colors.deepPurple.shade700, // purple-700
          background: Colors.transparent, // gray-800
          surface: Colors.grey.shade900, // gray-900
          onPrimary: Colors.grey.shade200, // gray-200 (sur le primary)
          onSecondary: Colors.grey.shade200, // gray-200 (sur le secondary)
          onSurface: Colors.grey.shade200, // gray-200 (sur le surface)
          onBackground: Colors.grey.shade200, // gray-200 (sur le background)
        ),
        scaffoldBackgroundColor: const Color(0xFF111827), // gray-900
        listTileTheme: ListTileThemeData(
          tileColor: Colors.deepPurple.shade600,
          textColor: Colors.grey.shade200,
          iconColor: Colors.grey.shade200,
        ),
        appBarTheme: AppBarTheme(
          backgroundColor: Colors.deepPurple.shade800,
        ));
  }
}
