import smtplib
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText
from email.mime.base import MIMEBase
from email import encoders


class EmailSender:
    def __init__(self, smtp_server, smtp_port, smtp_user, smtp_password):
        self.smtp_server = smtp_server
        self.smtp_port = int(smtp_port)
        self.smtp_user = smtp_user
        self.smtp_password = smtp_password

    def send_email(self, subject, to_addr, body_html, attachments=None):
        # Création de l'objet message
        msg = MIMEMultipart()
        msg["From"] = self.smtp_user
        msg["To"] = to_addr
        msg["Subject"] = subject

        # Ajout du corps HTML de l'email
        msg.attach(MIMEText(body_html, "html"))

        # Traitement des pièces jointes
        if attachments:
            for filepath in attachments:
                part = MIMEBase("application", "octet-stream")
                with open(filepath, "rb") as file:
                    part.set_payload(file.read())
                encoders.encode_base64(part)
                part.add_header(
                    "Content-Disposition",
                    'attachment; filename="{}"'.format(filepath.split("/")[-1]),
                )
                msg.attach(part)

        # Connexion au serveur SMTP et envoi de l'email
        with smtplib.SMTP_SSL(self.smtp_server, self.smtp_port) as server:
            server.login(self.smtp_user, self.smtp_password)
            server.send_message(msg)
            server.quit()


# Utilisation de la classe
if __name__ == "__main__":
    # Remplacez les valeurs suivantes par vos informations de connexion SMTP
    smtp_server = "smtp.ionos.fr"
    smtp_port = 465  # Port SMTP SSL généralement
    smtp_user = "no-reply@gouel.fr"
    smtp_password = "Pacifier2-Angles3-Quarters2-Marlin3"

    # Création d'une instance de EmailSender
    email_sender = EmailSender(smtp_server, smtp_port, smtp_user, smtp_password)

    # Envoi d'un email
    subject = "Sujet de l'email"
    to_addr = "test@iziram.fr"
    body_html = """\
    <html>
      <head></head>
      <body>
        <p>Bonjour !<br>
           Voici un email envoyé avec <b>du contenu HTML</b> et des pièces jointes.
        </p>
      </body>
    </html>
    """
    email_sender.send_email(subject, to_addr, body_html)
