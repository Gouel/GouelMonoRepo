import 'package:flutter/material.dart';

class GouelStepBuilder extends StatefulWidget {
  final List<Widget Function(Map<String, dynamic>)> steps;
  final Function(Map<String, dynamic>) onValidate;
  final Map<String, dynamic> formData;
  final CrossAxisAlignment alignment;

  GouelStepBuilder({
    Key? key,
    required this.steps,
    required this.onValidate,
    this.alignment = CrossAxisAlignment.center,
  })  : formData = {},
        super(key: key);

  @override
  GouelStepBuilderState createState() => GouelStepBuilderState();
}

class GouelStepBuilderState extends State<GouelStepBuilder> {
  int _currentStep = 0;

  void _goToNextStep() {
    if (_currentStep < widget.steps.length - 1) {
      setState(() {
        _currentStep++;
      });
    }
  }

  void _goToPreviousStep() {
    if (_currentStep > 0) {
      setState(() {
        _currentStep--;
      });
    }
  }

  void _validate() {
    widget.onValidate(widget.formData);
  }

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onHorizontalDragEnd: (DragEndDetails details) {
        if (details.primaryVelocity! > 0) {
          _goToPreviousStep(); // Swipe vers la droite
        } else if (details.primaryVelocity! < 0) {
          _goToNextStep(); // Swipe vers la gauche
        }
      },
      child: Column(
        crossAxisAlignment: widget.alignment,
        children: [
          widget.steps[_currentStep](widget.formData),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: <Widget>[
              // Bouton Précédent ou espace vide si première étape
              _currentStep > 0
                  ? IconButton(
                      onPressed: _goToPreviousStep,
                      icon: const Icon(Icons.arrow_back,
                          color: Colors.blue, size: 32),
                    )
                  : const SizedBox(width: 42), // Maintenir l'alignement

              // Indicateurs de progression au centre
              Row(
                mainAxisSize: MainAxisSize.min,
                children: _buildStepIndicators(),
              ),

              // Bouton Suivant ou Valider ou espace vide si dernière étape
              _currentStep < widget.steps.length - 1
                  ? IconButton(
                      onPressed: _goToNextStep,
                      icon: const Icon(Icons.arrow_forward,
                          color: Colors.blue, size: 32),
                    )
                  : _currentStep == widget.steps.length - 1
                      ? IconButton(
                          onPressed: _validate,
                          icon: const Icon(Icons.check,
                              color: Colors.blue, size: 32),
                        )
                      : const SizedBox(width: 32), // Maintenir l'alignement
            ],
          ),
        ],
      ),
    );
  }

  List<Widget> _buildStepIndicators() {
    int totalSteps = widget.steps.length;
    int startStep = _currentStep > 2 ? _currentStep - 1 : 0;
    int endStep =
        _currentStep < totalSteps - 3 ? _currentStep + 1 : totalSteps - 1;

    if (totalSteps > 5) {
      startStep = _currentStep > 1 ? _currentStep - 1 : 0;
      endStep = startStep + 2;
      if (endStep >= totalSteps) {
        endStep = totalSteps - 1;
        startStep = endStep - 2;
      }
    }

    return List<Widget>.generate(
      totalSteps > 5 ? 3 : totalSteps,
      (index) {
        int actualIndex = totalSteps > 5 ? startStep + index : index;
        return Container(
          margin: const EdgeInsets.symmetric(horizontal: 4.0),
          child: Icon(
            Icons.circle,
            color: actualIndex == _currentStep ? Colors.blue : Colors.grey,
            size: actualIndex == _currentStep ? 12.0 : 8.0,
          ),
        );
      },
    );
  }
}
