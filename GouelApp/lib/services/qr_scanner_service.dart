import 'package:flutter/material.dart';
import 'package:qr_code_scanner/qr_code_scanner.dart';

import '../widgets/gouel_bottom_sheet.dart';

class QRScannerService {
  void scanQR(BuildContext context, String title, Function(String) onResult,
      Function(String) onClose) {
    showModalBottomSheet(
      isScrollControlled: true,
      context: context,
      builder: (BuildContext context) {
        var scanArea = (MediaQuery.of(context).size.width < 400 ||
                MediaQuery.of(context).size.height < 400)
            ? 200.0
            : 300.0;
        return GouelBottomSheet(
          title: title,
          child: ClipRRect(
            borderRadius: const BorderRadius.all(Radius.circular(16)),
            child: SizedBox(
              width: scanArea,
              height: scanArea,
              child: QRView(
                key: GlobalKey(debugLabel: 'QR'),
                overlay: QrScannerOverlayShape(
                    borderColor: Colors.red,
                    borderRadius: 10,
                    borderLength: 30,
                    borderWidth: 10,
                    cutOutSize: scanArea),
                onQRViewCreated: (QRViewController controller) {
                  controller.scannedDataStream.listen((scanData) {
                    controller.pauseCamera();

                    Navigator.of(context).pop();
                    onResult(scanData.code ?? "");
                  });
                },
              ),
            ),
          ),
        );
      },
    ).then((_) => onClose('')); // Si le bottom sheet est fermé
  }
}
