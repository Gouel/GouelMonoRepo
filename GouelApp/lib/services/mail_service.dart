import 'dart:io';
import 'dart:typed_data';

import 'package:mailer/mailer.dart';
import 'package:mailer/smtp_server.dart';

class EmailService {
  final String host;
  final int port;
  final String username;
  final String password;
  final bool isSecure;

  EmailService({
    required this.host,
    required this.port,
    required this.username,
    required this.password,
    this.isSecure = true,
  });

  Future<void> sendEmail({
    required String recipient,
    required String sender,
    required String subject,
    required String content,
    List<String> attachments = const [],
  }) async {
    final smtpServer = SmtpServer(
      host,
      port: port,
      username: username,
      password: password,
      ssl: isSecure,
    );

    final message = Message()
      ..from = Address(sender)
      ..recipients.add(recipient)
      ..subject = subject
      ..text = content;

    for (var attachmentPath in attachments) {
      message.attachments.add(FileAttachment(File(attachmentPath)));
    }

    await send(message, smtpServer);
  }

  Future<void> sendEmailWithMemoryAttachment({
    required String recipient,
    required String sender,
    required String subject,
    required String content,
    Uint8List? attachmentData,
    String attachmentName = 'attachment.png',
  }) async {
    final smtpServer = SmtpServer(
      host,
      port: port,
      username: username,
      password: password,
      ssl: isSecure,
    );

    final message = Message()
      ..from = Address(sender)
      ..recipients.add(recipient)
      ..subject = subject
      ..text = content;

    if (attachmentData != null) {
      final attachment =
          MemoryAttachment(attachmentData, fileName: attachmentName);
      message.attachments.add(attachment);
    }

    await send(message, smtpServer);
  }
}

class MemoryAttachment extends StreamAttachment {
  final Uint8List data;
  final String fileName;
  final String contentType;

  MemoryAttachment(this.data,
      {required this.fileName, this.contentType = "image/png"})
      : super(Stream.fromIterable([data]), contentType, fileName: fileName);

  @override
  Stream<List<int>> asStream() {
    return Stream.value(data);
  }
}
