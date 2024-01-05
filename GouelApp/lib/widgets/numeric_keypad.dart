import 'package:flutter/material.dart';
import 'package:gouel/widgets/gouel_button.dart';

class NumericKeypad extends StatefulWidget {
  final Function(String) onNumberSelected;
  final String input;

  const NumericKeypad(
      {Key? key, required this.onNumberSelected, required this.input})
      : super(key: key);

  @override
  NumericKeypadState createState() => NumericKeypadState();
}

class NumericKeypadState extends State<NumericKeypad> {
  String _input = "";
  void _onKeyPress(String value) {
    if (value == ".") {
      if (_input.isEmpty) return;
      if (_input.contains(".")) return;
    }

    if (value == "0" && _input.isEmpty) return;

    if (_input.length > 5) return;

    setState(() {
      _input += value;
      widget.onNumberSelected(_input);
    });
  }

  void _onDelete() {
    if (_input.isNotEmpty) {
      setState(() {
        _input = _input.substring(0, _input.length - 1);
        widget.onNumberSelected(_input);
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    _input = widget.input;
    return Container(
      padding: const EdgeInsets.all(8),
      child: Column(
        children: [
          GridView.count(
            shrinkWrap: true,
            crossAxisCount: 3,
            childAspectRatio: 2,
            children: [
              ...List.generate(9, (index) {
                return _buildKeyButton('${index + 1}');
              }),
              ...[
                _buildKeyButton('.', flex: 1),
                _buildKeyButton('0'),
                _buildDeleteButton(),
              ]
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildKeyButton(String value, {int flex = 1}) {
    return Padding(
        padding: const EdgeInsets.all(4.0),
        child: GouelButton(
          text: value,
          onTap: () => _onKeyPress(value),
        ));
  }

  Widget _buildDeleteButton() {
    return Padding(
        padding: const EdgeInsets.all(4.0),
        child: GouelButton(
          onTap: _onDelete,
          color: Colors.red,
          text: "",
          icon: Icons.backspace,
        ));
  }
}
