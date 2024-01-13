from flask import Blueprint, render_template, request, session, abort
import json
from .generate_fake_data import db
from .helper import GouelHelper
import datetime as dt

main = Blueprint("main", __name__)


@main.route("/")
def index():
    events = [e for e in GouelHelper.get_events() if e["IsPublic"]]
    return render_template("pages/client/accueil.j2", events=events)


@main.route("/event/<event_id>")
def event(event_id):
    event = GouelHelper.get_event(event_id)

    return render_template("pages/client/evenement.j2", event=event, event_id=event_id)


@main.route("/acheter-billets/<event_id>")
def billets(event_id):
    event = GouelHelper.get_event(event_id)

    panier = session.get(f"panier-{event_id}", [])

    return render_template(
        "pages/client/acheter-billet.j2",
        event=event,
        event_id=event_id,
        panier=json.dumps(panier),
    )


@main.route("/payment-response")
def payment_response():
    # TODO : sécuriser cette route
    # utiliser api hello asso pour vérifier le paiement

    checkout_intent_id = request.args.get("checkoutIntentId")
    code = request.args.get("code")
    action = request.args.get("action")

    if action == "payment":
        action = code

    return render_template(
        "pages/client/payment_response.j2", action=action, backUrl=session["backUrl"]
    )


@main.route("/solde/<user_id>")
def solde(user_id: int):
    # TODO sécuriser avec un magic link envoyé par email
    # utiliser la session pour stocker le magic link  (session["magic_link"] = magic_link)

    u = GouelHelper.get_user(user_id)
    total_depense = sum(t["Amount"] for t in u["Transactions"] if t["Type"] == "debit")

    u["Transactions"] = sorted(
        u["Transactions"],
        key=lambda t: dt.datetime.strptime(t["Date"][:16], "%Y-%m-%dT%H:%M"),
        reverse=True,
    )

    return render_template("pages/client/solde.j2", user=u, total_depense=total_depense)
