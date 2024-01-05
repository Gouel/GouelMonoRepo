def send_email(
    to,
    sender="Gouel<noreply.gouel@iziram.com>",
    cc=None,
    bcc=None,
    subject=None,
    body=None,
):
    """sends email using a Jinja HTML template"""
    import smtplib

    # Import the email modules
    from email.mime.multipart import MIMEMultipart
    from email.mime.text import MIMEText
    from email.header import Header
    from email.utils import formataddr

    # convert TO into list if string
    if type(to) is not list:
        to = to.split()

    to_list = to + [cc] + [bcc]
    to_list = filter(None, to_list)  # remove null emails

    msg = MIMEMultipart("alternative")
    msg["From"] = sender
    msg["Subject"] = subject
    msg["To"] = ",".join(to)
    msg["Cc"] = cc
    msg["Bcc"] = bcc
    msg.attach(MIMEText(body, "html"))
    server = smtplib.SMTP("127.0.0.1")  # or your smtp server
    try:
        log.info("sending email xxx")
        server.sendmail(sender, to_list, msg.as_string())
    except Exception as e:
        log.error("Error sending email")
        log.exception(str(e))
    finally:
        server.quit()
