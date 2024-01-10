import 'package:flutter/material.dart';
import 'package:gouel/widgets/gouel_button.dart';

class GouelSwitch extends StatefulWidget {
  final bool initialValue;
  final String label;
  final ValueChanged<bool> onChange;

  const GouelSwitch(
      {Key? key,
      required this.initialValue,
      required this.onChange,
      required this.label})
      : super(key: key);

  @override
  GouelSwitchState createState() => GouelSwitchState();
}

class GouelSwitchState extends State<GouelSwitch> {
  late bool _isOn;

  @override
  void initState() {
    super.initState();
    _isOn = widget.initialValue;
  }

  void _toggleSwitch() {
    setState(() {
      _isOn = !_isOn;
    });
    widget.onChange(_isOn);
  }

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: _toggleSwitch,
      child: GouelButton(
        text: widget.label,
        color: _isOn ? Colors.blue : Colors.grey,
        onTap: _toggleSwitch,
      ),
    );
  }
}
