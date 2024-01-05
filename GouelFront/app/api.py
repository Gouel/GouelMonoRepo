from flask import (
    Blueprint,
    render_template,
    request,
    url_for,
    g,
    session,
    make_response,
    jsonify,
)
import json
from app import ha
from flask_json import as_json
from .generate_fake_data import db
from .hello_asso import CheckoutOrder, Payer, HelloAssoAPI
import re

api = Blueprint("api", __name__)


@as_json
@api.route("/")
def api_index():
    return "ok"


@as_json
@api.route("/event")
def api_events():
    return db["events"]


@as_json
@api.route("/event/<int:event_id>")
def api_event(event_id: int):
    return db["events"][event_id]


@as_json
@api.route("/event/<int:event_id>/ticket")
def api_event_tickets(event_id: int):
    return db["events"][event_id]["tickets"]


@as_json
@api.route("/event/<int:event_id>/ticket/<int:ticket_id>")
def api_event_ticket(event_id: int, ticket_id):
    return db["events"][event_id]["tickets"][ticket_id]


@api.route("/checkout/<int:event_id>", methods=["POST"])
def api_checkout(event_id: int):
    panier = request.json
    tickets = db["events"][event_id]["tickets"]
    totalAmount = 0
    for ticket in panier:
        ticket_id = ticket.get("id", -1)
        try:
            ticket_id = int(ticket_id)
        except:
            ticket_id = None
        if ticket_id is None or ticket_id >= len(tickets):
            jsonify({"error": "Ticket invalide ou mal formé"}), 400

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

        if ticket.get("credit", 0) < 0:
            return jsonify({"error": "Le crédit doit être supérieur ou égal à 0"}), 400

        if not re.match(r"[^@]+@[^@]+\.[^@]+\b", ticket.get("email", "")):
            return jsonify({"error": "Email non conforme"}), 400

        if not re.match(
            r"^\d{4}\-(0[1-9]|1[012])\-(0[1-9]|[12][0-9]|3[01])$", ticket.get("age", "")
        ):
            return jsonify({"error": "Date de naissance non conforme"}), 400

        totalAmount += tickets[ticket_id]["price"] + ticket.get("credit", 0)
    # Génération du checkout
    payer: Payer = Payer(
        firstName=panier[0].get("firstName"),
        lastName=panier[0].get("lastName"),
        email=panier[0].get("email"),
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
    session["backUrl"] = url_for("main.billets", event_id=event_id, _external=True)

    # Sauvegarde du panier dans la session
    session[f"panier-{event_id}"] = panier

    return jsonify(session["checkout"])


@api.route("/recharge/<user_id>", methods=["POST"])
def api_recharge(user_id: int):
    u = db["users"][0]

    recharge = request.json
    if recharge.get("credit", 0) < 1:
        return jsonify({"error": "Le crédit doit être supérieur ou égal à 1"}), 400
    # Génération du checkout
    payer: Payer = Payer(
        firstName=u.get("prenom"),
        lastName=u.get("nom"),
        email=u.get("email"),
    )
    chk = CheckoutOrder(
        totalAmount=recharge.get("credit") * 100,
        initialAmount=recharge.get("credit") * 100,
        backUrl=url_for("main.solde", user_id=user_id, _external=True),
        returnUrl=url_for("main.payment_response", action="payment", _external=True),
        errorUrl=url_for("main.payment_response", action="error", _external=True),
        itemName=f"payment-recharge_{user_id}",
        meta={"recharge": recharge},
        payer=payer,
    )

    # Sauvegarde du checkout dans la session
    session["checkout"] = ha.addCheckOutOrder(chk)
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
