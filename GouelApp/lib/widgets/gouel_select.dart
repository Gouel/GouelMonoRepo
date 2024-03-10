import 'package:flutter/material.dart';

class GouelSelect extends StatefulWidget {
  final List<GouelSelectItem> items;
  final String initialValue;
  final ValueChanged<GouelSelectItem> onChange;
  final String title;

  const GouelSelect({
    Key? key,
    required this.items,
    required this.initialValue,
    required this.onChange,
    required this.title,
  }) : super(key: key);

  @override
  GouelSelectState createState() => GouelSelectState();
}

class GouelSelectState extends State<GouelSelect> {
  late String _currentValue;

  @override
  void initState() {
    super.initState();
    _currentValue = widget.initialValue;
  }

  List<DropdownMenuItem<String>> _getMenuItems() {
    var dropdownItems = <DropdownMenuItem<String>>[
      DropdownMenuItem(
        value: '', // Valeur vide pour l'élément de titre
        enabled: false, // Rendre l'élément de titre non sélectionnable
        child: Text(widget.title),
      ),
    ];

    dropdownItems.addAll(
        widget.items.map<DropdownMenuItem<String>>((GouelSelectItem item) {
      return DropdownMenuItem<String>(
        value: item.value,
        child: Text(item.label),
      );
    }).toList());

    return dropdownItems.toList();
  }

  @override
  Widget build(BuildContext context) {
    return Theme(
      data: Theme.of(context).copyWith(canvasColor: const Color(0xFF111928)),
      child: DropdownButton<String>(
        value: _currentValue,
        onChanged: (String? newValue) {
          if (newValue != null) {
            setState(() {
              _currentValue = newValue;
            });
            widget.onChange(
                widget.items.firstWhere((el) => el.value == newValue));
          }
        },
        items: _getMenuItems(),
      ),
    );
  }
}

class GouelSelectItem {
  final String value;
  final String label;
  final Object? data;

  GouelSelectItem({required this.value, required this.label, this.data});
}
