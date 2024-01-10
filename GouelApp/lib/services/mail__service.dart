import 'dart:io';

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
    final smtpServer = isSecure
        ? SmtpServer(
            host,
            port: port,
            username: username,
            password: password,
            ssl: true,
          )
        : SmtpServer(
            host,
            port: port,
            username: username,
            password: password,
          );

    final message = Message()
      ..from = Address(sender)
      ..recipients.add(recipient)
      ..subject = subject
      ..text = content;

    for (var attachmentPath in attachments) {
      message.attachments.add(FileAttachment(File(attachmentPath)));
    }

    try {
      await send(message, smtpServer);
    } on MailerException catch (_) {}
  }
}
