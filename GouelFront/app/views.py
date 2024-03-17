from flask import (
    Blueprint,
    render_template,
    request,
    session,
    redirect,
    url_for,
    g,
    abort,
    flash,
)
import json
from .helper import GouelHelper
from . import get_ga
from .gouel_server import GouelApi
import datetime as dt
from functools import wraps
from . import ha
from .mail import EmailSender

from .cache import MagicLink
from .qrcode_gen import qrcode_total

from .api import check_errors

main = Blueprint("main", __name__)


def get_user(user_id: str) -> dict:
    users: dict = session.get("known_users", {})
    if user_id not in users:
        users[user_id] = {
            **GouelHelper(get_ga()).get_user(user_id),
            "role": g.user["role"],
            "events": g.user["events"],
        }
        session["known_users"] = users
    return users[user_id]


def event_access(f):
    @wraps(f)
    def decorated_function(*args, **kwargs):
        # L'argument `event_id` doit être extrait des kwargs ici, si nécessaire.
        actual_event_id = kwargs.get("event_id", "blank")
        g.event = GouelHelper(get_ga()).get_event(actual_event_id)
        if g.event is None or not g.event["IsPublic"]:
            abort(404)
        g.event_id = g.event["ID"]
        return f(*args, **kwargs)

    return decorated_function


def login_required(f):
    @wraps(f)
    def decorated_function(*args, **kwargs):
        if session.get("compte") is None:
            return redirect(url_for("main.login"))
        g.uga = GouelApi.from_json(session.get("compte"))
        g.user = g.uga.token["infos"]
        return f(*args, **kwargs)

    return decorated_function


def get_email_sender() -> EmailSender:
    if "mailConf" not in session:
        session["mailConf"] = GouelHelper(get_ga()).get_conf_smtp()

    smtp = session["mailConf"]
    return EmailSender(
        smtp["SMTPServer"],
        smtp["SMTPPort"],
        smtp["Email"],
        smtp["EmailPassword"],
    )


@main.route("/")
def index():
    return render_template("pages/client/accueil.j2")


@main.route("/magic/<magic_id>")
def magic(magic_id):
    link = MagicLink.from_id(magic_id)
    if not link:
        abort(404)

    session["magic_link"] = str(link)

    match link.type:
        case "reset_password":
            return redirect(url_for("main.new_password"))
        case _:
            abort(404)


@main.route("/evenements")
def events():
    events = [e for e in GouelHelper(get_ga()).get_events() if e["IsPublic"]]
    return render_template(
        "pages/client/events.j2", events=events, retour=url_for("main.index")
    )


@main.route("/event/<event_id>")
@event_access
def event(event_id):
    return render_template(
        "pages/client/evenement.j2",
        event=g.event,
        event_id=g.event_id,
        retour=url_for("main.events"),
    )


@main.route("/event/<event_id>/billets")
@event_access
def billets(event_id):
    panier = session.get(f"panier-{g.event_id}", [])

    user = None
    compte = session.get("compte")
    if compte is not None:
        user = get_user(compte["token"]["infos"]["userId"])

    return render_template(
        "pages/client/acheter-billet.j2",
        event=g.event,
        event_id=g.event_id,
        panier=json.dumps(panier),
        user=user,
        retour=url_for("main.event", event_id=g.event_id),
    )


@main.route("/connexion", methods=["GET", "POST"])
def login():
    if request.method == "POST":
        username = request.form["email"]
        password = request.form["password"]

        uga: GouelApi = GouelApi(username, password)
        hasToken: bool = uga.get_token() is not None
        if hasToken:
            session["compte"] = uga.to_json()
            return redirect(url_for("main.compte"))
        else:
            flash("Identifiant ou mot de passe incorrect.", "error")

    return render_template("pages/client/auth/connexion.j2")


@main.route("/inscription", methods=["GET", "POST"])
def signin():
    if request.method == "POST":
        email = request.form["email"]
        password = request.form["password"]
        password_confirm = request.form["password_confirm"]
        first_name = request.form["firstName"]
        last_name = request.form["lastName"]
        dob = request.form["age"]

        pb = check_errors(first_name, last_name, email, dob)

        if password != password_confirm:
            flash("Les mots de passe ne correspondent pas.", "error")
        elif password == "":
            flash("Veuillez entrer un mot de passe.", "error")
        elif pb[1] != 200:
            flash(pb[0], "error")
        else:
            ok, userId = GouelHelper(get_ga()).add_user(
                {
                    "Email": email,
                    "Password": password,
                    "FirstName": first_name,
                    "LastName": last_name,
                    "DOB": dob,
                }
            )
            if ok:
                get_email_sender().send_email(
                    "Création de votre compte Gouel",
                    email,
                    render_template("mails/compte.j2", user={"Email": email}),
                )

                flash("Votre compte a été créé avec succès.", "success")
                return redirect(url_for("main.login"))
            else:
                flash(
                    "Une erreur est survenue lors de la création de votre compte.",
                    "error",
                )

    return render_template("pages/client/auth/inscription.j2")


@main.route("/oublie-mdp", methods=["GET", "POST"])
def forgot_password():
    if request.method == "POST":
        email = request.form["email"]
        if email:
            user = GouelHelper(get_ga()).get_user("", email=email)
            if user:
                link = MagicLink("reset_password", {"UserId": user["ID"]})
                url = url_for("main.magic", magic_id=link.id, _external=True)

                email_sender = get_email_sender()
                email_sender.send_email(
                    "Réinitialisation de votre mot de passe",
                    email,
                    render_template(
                        "mails/reset_password.j2",
                        user=user,
                        reset_password_url=url,
                    ),
                )

            flash(
                "Le lien pour réinitialiser votre mot de passe a été envoyé à votre adresse email.",
                "success",
            )
        else:
            flash("Veuillez entrer votre adresse email.", "error")

    return render_template("pages/client/auth/forgot_password.j2")


@main.route("/nouveau-mdp", methods=["GET", "POST"])
def new_password():
    magic_link: MagicLink = MagicLink.from_json(session.get("magic_link"))
    if magic_link is None:
        abort(404)

    if request.method == "POST":
        password = request.form["password"]
        password_confirm = request.form["password_confirm"]

        if password != password_confirm:
            flash("Les mots de passe ne correspondent pas.", "error")
        else:
            GouelHelper(get_ga()).update_user(
                magic_link.obj["UserId"], {"Password": password}
            )

            # TODO : envoyer un email pour confirmer le changement de mot de passe

            flash("Votre mot de passe a été mis à jour.", "success")
            session.pop("magic_link")
            return redirect(url_for("main.login"))

    return render_template(
        "pages/client/auth/new_password.j2", userId=magic_link.obj["UserId"]
    )


@main.route("/profil")
@login_required
def compte():
    user = get_user(g.user["userId"])
    return render_template("pages/client/auth/compte.j2", user=user)


@main.route("/deconnexion", methods=["GET"])
def logout():
    if "compte" in session:
        session.pop("compte")
        session.pop("known_users")
    return redirect(url_for("main.index"))


@main.route("/payment-response")
def payment_response():
    checkout_intent_id = request.args.get("checkoutIntentId")
    checkout_type = session.get("checkout_type")
    checkout = session.get("checkout")
    if checkout is None or str(checkout["id"]) != checkout_intent_id:
        abort(404)

    code = request.args.get("code")
    action = request.args.get("action")

    if action == "payment":
        action = code

    # vérification du payement avec l'api hello asso

    chk_verif = ha.getCheckOutOrder(checkout_intent_id)

    if "order" not in chk_verif:
        action = "error"
    else:
        # Tout est ok, on peut valider la commande
        ga = GouelHelper(get_ga())

        if checkout_type == "new_billets":
            payment_new_billet(ga)
        elif checkout_type == "recharge":
            payment_recharge(ga)
        else:
            abort(404)

    return render_template(
        "pages/client/payment_response.j2", action=action, backUrl=session["backUrl"]
    )


def new_transaction(event_id=None, amount=0) -> dict:
    trans = {
        "Date": dt.datetime.now().isoformat(),
        "Amount": amount,
        "Type": "credit",
        "Cart": [],
        "PaymentType": "helloasso",
    }
    if event_id:
        trans["EventId"] = event_id
    return trans


def payment_recharge(ga):
    checkout_user_id = session.get("checkout_user_id")
    checkout_amount = session.get("checkout_amount")

    user = ga.get_user(checkout_user_id)
    if user is None:
        abort(500)

    email_sender = get_email_sender()

    email_sender.send_email(
        "Confirmation du rechargement de votre compte Gouel",
        user["Email"],
        render_template(
            "mails/rechargement.j2",
            user=user,
            amount=checkout_amount,
        ),
    )

    # Ajout des crédits sur le compte
    if checkout_amount > 0:
        transaction = new_transaction(None, checkout_amount)
        ga.add_transaction(None, user["ID"], transaction)

    # On vide la session
    session.pop("checkout")
    session.pop("checkout_user_id")
    session.pop("checkout_amount")
    session.pop("checkout_type")


def payment_new_billet(ga):
    # On récupère le panier
    paniers = session.get("paniers", {})
    event_id = session.get("checkout_event_id", "")
    event = ga.get_event(event_id)
    panier = paniers.get(event_id, [])

    # Pour chaque ticket
    # - On crée l'utilisateur (sauf si déjà existant)
    # - On crée la transaction
    # - On envoie un email de confirmation

    def new_ticket(user_id) -> dict:
        return {"UserId": user_id, "PurchasedOnline": True}

    email_sender = get_email_sender()

    for ticket in panier:
        user: dict = None
        if "userId" not in ticket:
            ok, userId = ga.add_user(
                {
                    "Email": ticket["email"],
                    "FirstName": ticket["firstName"],
                    "LastName": ticket["lastName"],
                    "DOB": ticket["age"],
                }
            )
            if not ok:
                abort(500)

            user = ga.get_user(userId["UserId"])

            mdp_link = MagicLink("reset_password", {"UserId": user["ID"]}, 24 * 60 * 60)

            email_sender.send_email(
                "Création de votre compte Gouel",
                ticket["email"],
                render_template(
                    "mails/compte.j2",
                    user=user,
                    mdp_url=url_for("main.magic", magic_id=mdp_link.id, _external=True),
                ),
            )

        else:
            user = ga.get_user(ticket["userId"])

        if user is None:
            abort(500)

        # Ajout des crédits sur le compte
        transaction = new_transaction(event_id, ticket["credit"])
        ga.add_transaction(event_id, user["ID"], transaction)

        # Ajout du ticket
        new_ticket = new_ticket(user["ID"])
        ticket_id = ga.add_ticket(event_id, ticket["id"], new_ticket)

        if ticket_id is None:
            abort(500)

        qr = qrcode_total(ticket_id)

        email_sender.send_email(
            "Confirmation de l'achat de votre billet",
            user["Email"],
            render_template(
                "mails/ticket.j2",
                qrcodeBase64=qr,
                user=user,
                ticketId=ticket_id,
                event_id=event_id,
                event_name=event["Title"],
            ),
        )

    # On vide le panier
    paniers.pop(event_id)
    session["paniers"] = paniers
    session.pop("checkout")
    session.pop("checkout_event_id")
    session.pop("checkout_type")


@main.route("/solde")
@login_required
def solde():
    u = {
        **GouelHelper(get_ga()).get_user(g.user["userId"]),
        "role": g.user["role"],
        "events": g.user["events"],
    }
    total_depense = sum(t["Amount"] for t in u["Transactions"] if t["Type"] == "debit")

    u["Transactions"] = sorted(
        u["Transactions"],
        key=lambda t: dt.datetime.strptime(t["Date"][:16], "%Y-%m-%dT%H:%M"),
        reverse=True,
    )

    return render_template(
        "pages/client/solde.j2",
        user=u,
        total_depense=total_depense,
        gh=GouelHelper(get_ga()),
    )
