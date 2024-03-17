from flask import (
    Blueprint,
    request,
    url_for,
    session,
    jsonify,
    g,
)
from app import ha
from .hello_asso import CheckoutOrder, Payer
from . import get_ga
from .gouel_server import GouelApi
from .helper import GouelHelper
import re
from datetime import datetime

api = Blueprint("api", __name__)


@api.route("/checkout/<event_id>", methods=["POST"])
def api_checkout(event_id):
    panier = request.json
    event = GouelHelper(get_ga()).get_event(event_id)
    if event is None:
        return jsonify({"error": "Événement introuvable"}), 404

    tickets = {et["EventTicketCode"]: et for et in event["EventTickets"]}

    totalAmount = 0
    for ticket in panier:
        ticket_id = ticket.get("id", -1)
        if ticket_id is None or ticket_id not in tickets:
            return jsonify({"error": "Ticket invalide ou mal formé"}), 400

        # Validation des informations du ticket
        firstName = ticket.get("firstName", "")
        lastName = ticket.get("lastName", "")
        if firstName == lastName:
            return (
                jsonify({"error": "Nom et prénom ne doivent pas être identique"}),
                400,
            )

        res, err = verify_name(firstName, "prénom")
        if not res:
            return jsonify({{"error": err}}), 400

        res, err = verify_name(lastName, "nom")
        if not res:
            return jsonify({{"error": err}}), 400
        credit = ticket.get("credit", 0)
        if credit < 0:
            return jsonify({"error": "Le crédit doit être supérieur ou égal à 0"}), 400

        if credit > 50:
            return jsonify({"error": "Le crédit ne peut pas être supérieur à 50"}), 400
        email = ticket.get("email", "")
        if not re.match(r"[^@]+@[^@]+\.[^@]+\b", email):
            return jsonify({"error": "Email non conforme"}), 400

        if not re.match(
            r"^\d{4}\-(0[1-9]|1[012])\-(0[1-9]|[12][0-9]|3[01])$", ticket.get("age", "")
        ):
            return jsonify({"error": "Date de naissance non conforme"}), 400

        if "userId" not in ticket:
            user = GouelHelper(get_ga()).get_user("", email=email)
            if user is not None:
                return jsonify({"error": f"Email déjà utilisée ({email})"}), 400

        totalAmount += tickets[ticket_id]["Price"]["Online"] + credit

    # Génération du checkout
    dob: datetime = datetime.strptime(panier[0].get("age"), "%Y-%m-%d")

    payer: Payer = Payer(
        firstName=panier[0].get("firstName"),
        lastName=panier[0].get("lastName"),
        email=panier[0].get("email"),
        dob=dob,
    )
    chk = CheckoutOrder(
        totalAmount=totalAmount * 100,
        initialAmount=totalAmount * 100,
        backUrl=url_for("main.billets", event_id=event_id, _external=True),
        returnUrl=url_for("main.payment_response", action="payment", _external=True),
        errorUrl=url_for("main.payment_response", action="error", _external=True),
        itemName=f"payment-event_{event_id}",
        meta={"panier": panier},
        payer=payer,
    )
    # Sauvegarde du checkout dans la session
    session["checkout"] = ha.addCheckOutOrder(chk)
    session["checkout_event_id"] = event_id
    session["checkout_type"] = "new_billets"
    session["backUrl"] = url_for("main.billets", event_id=event_id, _external=True)

    # Sauvegarde du panier dans la session
    paniers = session.get("paniers", {})
    paniers[event_id] = panier
    session["paniers"] = paniers

    if chk_error := session["checkout"].get("errors"):
        return jsonify(chk_error), 400

    return jsonify(session["checkout"])


@api.route("/recharge", methods=["POST"])
def api_recharge():
    if session.get("compte") is None:
        return jsonify({"error": "Utilisateur non connecté"}), 401
    g.uga = GouelApi.from_json(session.get("compte"))
    user_id = g.uga.token["infos"]["userId"]
    g.user = GouelHelper(get_ga()).get_user(user_id)

    recharge = request.json
    if recharge.get("credit", 0) < 1:
        return jsonify({"error": "Le crédit doit être supérieur ou égal à 1"}), 400
    # Génération du checkout
    payer: Payer = Payer(
        firstName=g.user["FirstName"],
        lastName=g.user["LastName"],
        email=g.user["Email"],
        dob=datetime.strptime(g.user["DOB"], "%Y-%m-%d"),
    )
    chk = CheckoutOrder(
        totalAmount=recharge.get("credit") * 100,
        initialAmount=recharge.get("credit") * 100,
        backUrl=url_for("main.solde", _external=True),
        returnUrl=url_for("main.payment_response", action="payment", _external=True),
        errorUrl=url_for("main.payment_response", action="error", _external=True),
        itemName=f"payment-recharge_{user_id}",
        meta={"recharge": recharge},
        payer=payer,
    )

    # Sauvegarde du checkout dans la session
    session["checkout"] = ha.addCheckOutOrder(chk)
    session["checkout_user_id"] = user_id
    session["checkout_amount"] = recharge.get("credit")
    session["checkout_type"] = "recharge"
    session["backUrl"] = url_for("main.solde", user_id=user_id, _external=True)

    return jsonify(session["checkout"])


def verify_name(name: str, key: str = "nom") -> tuple[bool, str]:
    correct_name: bool = False

    # Liste des noms interdits
    forbidden_names = [
        "firstname",
        "lastname",
        "unknown",
        "first_name",
        "last_name",
        "anonyme",
        "user",
        "admin",
        "name",
        "nom",
        "prénom",
        "test",
    ]

    if name == "":
        return correct_name, f"Le {key} est vide"
    if len(name) > 255:
        return correct_name, f"Le {key} est trop long"
    if any(char.isdigit() for char in name):
        return correct_name, f"Le {key} ne doit pas contenir de chiffres"
    if len(name) == 1:
        return correct_name, f"Le {key} ne doit pas être un seul caractère"
    if name.lower() in forbidden_names:
        return correct_name, f"Le {key} est interdit"
    if not any(char in "aeiouyéèêëàâäôöûüç" for char in name.lower()):
        return correct_name, f"Le {key} doit contenir au moins une voyelle"
    if any(name[i : i + 3] == name[i] * 3 for i in range(len(name) - 2)):
        return (
            correct_name,
            f"Le {key} ne doit pas contenir de caractères répétitifs trois fois de suite",
        )
    if not all(
        char in "abcdefghijklmnopqrstuvwxyzéèêëàâäôöûüç'- " for char in name.lower()
    ):
        return correct_name, f"Le {key} contient des caractères non autorisés"

    # Si toutes les vérifications sont passées
    correct_name = True
    return correct_name, f"Le {key} est valide"


def check_errors(
    first_name, last_name, email, dob, check_user_email: bool = False
) -> tuple[str, int]:
    if first_name == last_name:
        return "Nom et prénom ne doivent pas être identique", 400

    res, err = verify_name(first_name, "prénom")
    if not res:
        return err, 400

    res, err = verify_name(last_name, "nom")
    if not res:
        return err, 400

    if not re.match(r"[^@]+@[^@]+\.[^@]+\b", email):
        return "Email non conforme", 400

    if not re.match(r"^\d{4}\-(0[1-9]|1[012])\-(0[1-9]|[12][0-9]|3[01])$", dob):
        return "Date de naissance non conforme", 400

    if check_user_email:
        user = GouelHelper(get_ga()).get_user("", email=email)
        if user is not None:
            return f"Email déjà utilisée ({email})", 400

    return "", 200
