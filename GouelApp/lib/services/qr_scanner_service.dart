import 'package:flutter/material.dart';
import 'package:mobile_scanner/mobile_scanner.dart';

import '../widgets/gouel_bottom_sheet.dart';

class QRScannerService {
  void scanQR(BuildContext context, String title, Function(String) onResult,
      Function(String) onClose) {
    showModalBottomSheet(
      isScrollControlled: true,
      context: context,
      builder: (BuildContext context2) {
        var scanArea = (MediaQuery.of(context2).size.width < 400 ||
                MediaQuery.of(context2).size.height < 400)
            ? 200.0
            : 300.0;

        MobileScannerController controller = MobileScannerController(
            detectionSpeed: DetectionSpeed.normal,
            facing: CameraFacing.back,
            detectionTimeoutMs: 1000,
            autoStart: true);
        return GouelBottomSheet(
          title: title,
          child: ClipRRect(
            borderRadius: const BorderRadius.all(Radius.circular(16)),
            child: SizedBox(
              width: scanArea,
              height: scanArea,
              child: MobileScanner(
                // fit: BoxFit.contain,
                controller: controller,
                onDetect: (capture) {
                  final List<Barcode> barcodes = capture.barcodes;
                  if (barcodes.isNotEmpty) {
                    Navigator.of(context2).pop();
                    onResult(barcodes.first.rawValue ?? "");
                  }
                },
              ),
            ),
          ),
        );
      },
    ).then((_) => onClose('')); // Si le bottom sheet est fermé
  }
}
