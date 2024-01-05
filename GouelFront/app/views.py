from flask import Blueprint, render_template, request, session
import json
from .generate_fake_data import db

main = Blueprint("main", __name__)


@main.route("/")
def index():
    events = db["events"]
    return render_template("pages/client/accueil.j2", events=events)


@main.route("/event/<int:event_id>")
def event(event_id):
    event = db["events"][event_id]
    return render_template("pages/client/evenement.j2", event=event, event_id=event_id)


@main.route("/acheter-billets/<int:event_id>")
def billets(event_id):
    # Ici, vous récupéreriez les données de l'événement à partir de votre base de données ou d'un autre service.
    event = db["events"][event_id]

    panier = session.get(f"panier-{event_id}", [])

    return render_template(
        "pages/client/acheter-billet.j2",
        event=event,
        event_id=event_id,
        panier=json.dumps(panier),
    )


@main.route("/payment-response")
def payment_response():
    checkout_intent_id = request.args.get("checkoutIntentId")
    code = request.args.get("code")
    action = request.args.get("action")

    if action == "payment":
        action = code

    return render_template(
        "pages/client/payment_response.j2", action=action, backUrl=session["backUrl"]
    )


@main.route("/solde/<int:user_id>")
def solde(user_id: int):
    u = db["users"][0]
    total_depense = sum([a["price"] * a["number"] for a in u["transactions"]])
    return render_template("pages/client/solde.j2", user=u, total_depense=total_depense)
