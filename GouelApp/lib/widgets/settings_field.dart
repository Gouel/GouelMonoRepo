import 'package:flutter/material.dart';

enum SettingsFieldType {
  inputText,
  inputNumeric,
  switchField,
  select,
  // Ajouter d'autres types ici si nécessaire
}

class SettingsField extends StatelessWidget {
  final SettingsFieldType type;
  final String label;
  final dynamic value;
  final Function(dynamic)? onChanged;
  final Function(dynamic)? onFinish;
  final List<String>? selectOptions; // Pour les champs de type select

  const SettingsField({
    Key? key,
    required this.type,
    required this.label,
    required this.value,
    this.onChanged,
    this.onFinish,
    this.selectOptions,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    switch (type) {
      case SettingsFieldType.inputText:
        var controller = TextEditingController(text: value);
        return TextField(
          decoration: InputDecoration(
            labelText: label,
            border: const OutlineInputBorder(),
          ),
          controller: controller,
          onChanged: onChanged,
          onSubmitted: (value) {
            onFinish!(value);
          },
        );
      case SettingsFieldType.inputNumeric:
        return TextField(
          decoration: InputDecoration(
            labelText: label,
            border: const OutlineInputBorder(),
          ),
          controller: TextEditingController(text: value.toString()),
          keyboardType: TextInputType.number,
          onChanged: (val) => onChanged!(num.tryParse(val)),
        );
      case SettingsFieldType.switchField:
        return SwitchListTile(
          title: Text(label),
          value: value,
          onChanged: (val) => onChanged!(val),
        );
      case SettingsFieldType.select:
        return DropdownButtonFormField<String>(
          value: value,
          decoration: InputDecoration(labelText: label),
          items: selectOptions?.map<DropdownMenuItem<String>>((String value) {
            return DropdownMenuItem<String>(
              value: value,
              child: Text(value),
            );
          }).toList(),
          onChanged: (val) => onChanged!(val),
        );
      // Ajouter d'autres cas pour de nouveaux types
      default:
        return const SizedBox.shrink();
    }
  }
}
