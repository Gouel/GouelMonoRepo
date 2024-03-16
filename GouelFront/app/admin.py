from flask import (
    g,
    render_template,
    request,
    redirect,
    url_for,
    session,
    Blueprint,
    flash,
)
from functools import wraps
from . import GouelApi, EmailSender, get_ga
from .helper import GouelHelper
import json
import os

from .cache import MagicLink

admin = Blueprint("admin", __name__)


def login_required(f):
    @wraps(f)
    def decorated_function(*args, **kwargs):
        if session.get("compte") is None:
            return redirect(url_for("admin.login"))
        g.uga = GouelApi.from_json(session.get("compte"))
        g.user = g.uga.token["infos"]
        return f(*args, **kwargs)

    return decorated_function


def event_access(f):
    @wraps(f)
    def decorated_function(*args, **kwargs):
        # L'argument `event_id` doit être extrait des kwargs ici, si nécessaire.
        actual_event_id = kwargs.get("event_id", "blank")

        if g.user["events"] is None:
            if g.user["role"] != "SUPERADMIN":
                return redirect(url_for("admin.dashboard"))
        elif actual_event_id not in [i["EventId"] for i in g.user["events"]]:
            return redirect(url_for("admin.dashboard"))
        g.event = GouelHelper(g.uga).get_event(actual_event_id)
        return f(*args, **kwargs)

    return decorated_function


@admin.route("/login", methods=["GET", "POST"])
def login():
    if request.method == "POST":
        username = request.form["username"]
        password = request.form["password"]

        uga: GouelApi = GouelApi(username, password)
        hasToken: bool = uga.get_token() is not None
        if hasToken:
            session["compte"] = uga.to_json()
            return redirect(url_for("admin.dashboard"))
        else:
            flash("Identifiant ou mot de passe incorrect.", "error")

    return render_template("pages/admin/login.j2")


@admin.route("/", methods=["GET", "POST"])
@admin.route("/dashboard", methods=["GET", "POST"])
@login_required
def dashboard(**kwargs):
    if request.method == "POST":
        if g.user["role"] != "SUPERADMIN":
            flash("Vous n'avez pas les droits pour effectuer cette action.", "error")
            return redirect(url_for("admin.dashboard"))
        action: str = request.form.get("action")

        if action == "add":
            event_title = request.form["Title"]
            GouelHelper(g.uga).add_event({"Title": event_title})

        elif action == "delete":
            event_id = request.form["EventId"]
            GouelHelper(g.uga).delete_event(event_id)
            # suppression de l'image liée à l'événement
            img_path = f"app/static/img/events/{event_id}"
            if os.path.exists(img_path):
                os.remove(img_path)
    events: dict = {}
    try:
        events = GouelHelper(g.uga).get_events(True)
    except Exception:
        return redirect(url_for("admin.deconnexion"))

    return render_template(
        "pages/admin/dashboard.j2",
        user=g.user,
        all_events=events,
    )


@admin.route("/manage/<event_id>", methods=["POST", "GET"])
@login_required
@event_access
def manage_event(event_id: str, **kwargs):
    if request.method == "POST":
        # POST (Modifie les informations de l'événement)
        update_event: dict = {
            "Title": request.form["Title"],
            "Contact": request.form["Contact"],
            "Description": request.form["Description"],
            "Location": request.form["Location"],
            "IsPublic": request.form.get("IsPublic") == "true",
        }
        GouelHelper(g.uga).update_event(event_id, update_event)

        # les images sont stockés dans le dossier static/img/events/<event_id>
        if (
            request.files
            and "Image" in request.files
            and request.files["Image"] is not None
            and request.files["Image"].filename != ""
        ):
            img = request.files["Image"]
            img.save(f"app/static/img/events/{event_id}")

        flash("Les informations ont été modifiées.", "success")

        return redirect(url_for("admin.manage_event", event_id=event_id))

    # GET (Affiche les informations de l'événement)
    return render_template(
        "pages/admin/manage_event.j2",
        event_id=event_id,
        event=g.event,
        retour=url_for("admin.dashboard"),
    )


@admin.route("/manage/<event_id>/volunteers", methods=["POST", "GET"])
@login_required
@event_access
def manage_volunteers(event_id: str):
    if request.method == "POST":
        r = True
        action: str = request.form.get("action")
        if action == "add":
            if request.form.get("newUser") == "on":
                # Créer un nouvel utilisateur
                user = {
                    "FirstName": request.form["FirstName"],
                    "LastName": request.form["LastName"],
                    "Email": request.form["Email"],
                    "Dob": request.form["Dob"],
                }

                if "mailConf" not in session:
                    session["mailConf"] = GouelHelper(get_ga()).get_conf_smtp()

                smtp = session["mailConf"]
                email_sender = EmailSender(
                    smtp["SMTPServer"],
                    smtp["SMTPPort"],
                    smtp["Email"],
                    smtp["EmailPassword"],
                )

                r, e = GouelHelper(get_ga()).add_user(user)
                if not r:
                    flash(f"{e['error']}", "error")
                    return redirect(
                        url_for("admin.manage_volunteers", event_id=event_id)
                    )

                mdp_link = MagicLink(
                    "reset_password", {"UserId": user["ID"]}, 24 * 60 * 60
                )

                email_sender.send_email(
                    "Création de votre compte Gouel",
                    user["Email"],
                    render_template(
                        "mails/compte.j2",
                        user=user,
                        mdp_url=url_for(
                            "main.magic", magic_id=mdp_link.id, _external=True
                        ),
                    ),
                )

            user = GouelHelper(get_ga()).get_user("", email=request.form["Email"])
            if user is None:
                flash("L'utilisateur n'existe pas.", "error")
                return redirect(url_for("admin.manage_volunteers", event_id=event_id))

            volunteer = {
                "UserId": user["ID"],
                "Permissions": request.form.getlist("droits"),
            }
            r = GouelHelper(g.uga).add_volunteer(event_id, volunteer)
            if r:
                flash("Le bénévole a été ajouté.", "success")
            else:
                flash("Le bénévole n'a pas été ajouté.", "error")
        elif action == "edit":
            volunteer = {
                "UserId": request.form["UserId"],
                "Permissions": request.form.getlist("droits"),
                "IsAdmin": request.form.get("IsAdmin") == "on",
            }
            r = GouelHelper(g.uga).update_volunteer(event_id, volunteer)
            if r:
                flash("Les droits ont été modifiés.", "success")
            else:
                flash("Les droits n'ont pas été modifiés.", "error")

        elif action == "delete":
            r = GouelHelper(g.uga).delete_volunteer(event_id, request.form["UserId"])
            if r:
                flash("Le bénévole a été supprimé.", "success")
            else:
                flash("Le bénévole n'a pas été supprimé.", "error")
        if r:
            return redirect(url_for("admin.manage_volunteers", event_id=event_id))

    volunteers = []
    users: dict = session.get("known_users", {})
    for vol in g.event["Volunteers"]:
        user = users.get(vol["UserId"])
        if user is None:
            user = GouelHelper(get_ga()).get_user(vol["UserId"])
        if vol["UserId"] not in users:
            users[vol["UserId"]] = user
        user["droits"] = vol["Permissions"]
        user["IsAdmin"] = vol["IsAdmin"] or False
        volunteers.append(user)
    volunteers.sort(key=lambda x: (x["LastName"], x["FirstName"]))
    session["known_users"] = users

    # GET (Affiche les informations des bénévoles)
    return render_template(
        "pages/admin/manage_volunteers.j2",
        volunteers=volunteers,
        retour=url_for("admin.manage_event", event_id=event_id),
    )


@admin.route("/manage/<event_id>/tickets", methods=["POST", "GET"])
@login_required
@event_access
def manage_tickets(event_id: str):
    tickets: list[dict] = GouelHelper(g.uga).get_tickets(event_id)
    # On ne prend pas en compte les tickets donnés (par exemple bénévoles)
    tickets = [t for t in tickets if t["WasPurchased"]]
    tickets_valides: list[dict] = [t for t in tickets if t["IsUsed"]]
    event_tickets = {
        et["EventTicketCode"]: et["Price"] for et in g.event["EventTickets"]
    }
    total_prix_tickets: float = 0
    for t in tickets:
        et = event_tickets.get(t["EventTicketCode"])
        if et is not None:
            total_prix_tickets += et["Online"] if t["PurchasedOnline"] else et["OnSite"]
    total_depense_credits: float = 0.0
    return render_template(
        "pages/admin/manage_tickets.j2",
        tickets=tickets,
        tickets_valides=len(tickets_valides),
        tickets_total=len(tickets),
        total_prix_tickets=total_prix_tickets,
        total_depense_credits=total_depense_credits,
        event_id=event_id,
        pagination={
            "has_prev": False,
            "has_next": False,
            "prev_num": 0,
            "next_num": 0,
        },
        retour=url_for("admin.manage_event", event_id=event_id),
    )


@admin.route("/manage/<event_id>/tickets/event_tickets", methods=["POST", "GET"])
@login_required
@event_access
def manage_event_tickets(event_id):
    if request.method == "POST":
        action: str = request.form.get("action")
        if not action:
            return redirect(url_for("admin.manage_event_tickets", event_id=event_id))

        if action == "add":
            event_ticket = {
                "Title": request.form["Title"],
                "Price": {
                    "Online": float(request.form["PriceOnline"]),
                    "OnSite": float(request.form["PriceOnSite"]),
                },
            }
            r = GouelHelper(g.uga).add_event_ticket(event_id, event_ticket)
            if r:
                flash("Le billet a été ajouté.", "success")
            else:
                flash("Le billet n'a pas été ajouté.", "error")
        elif action == "edit":
            event_ticket = {
                "Title": request.form["Title"],
                "Price": {
                    "Online": float(request.form["PriceOnline"]),
                    "OnSite": float(request.form["PriceOnSite"]),
                },
            }
            ticket_id = request.form["EventTicketCode"]
            r = GouelHelper(g.uga).update_event_ticket(
                event_id, ticket_id, event_ticket
            )
            if r:
                flash("Le billet a été modifié.", "success")
            else:
                flash("Le billet n'a pas été modifié.", "error")
        elif action == "delete":
            ticket_id = request.form["EventTicketCode"]
            r = GouelHelper(g.uga).delete_event_ticket(event_id, ticket_id)
            if r:
                flash("Le billet a été supprimé.", "success")
            else:
                flash("Le billet n'a pas été supprimé.", "error")

        return redirect(url_for("admin.manage_event_tickets", event_id=event_id))

    return render_template(
        "pages/admin/manage_event_tickets.j2",
        event_id=event_id,
        event_tickets=g.event["EventTickets"],
        retour=url_for("admin.manage_tickets", event_id=event_id),
    )


@admin.route("/maage/event")
@admin.route("/manage/<event_id>/options", methods=["POST", "GET"])
@login_required
@event_access
def manage_options(event_id):
    if request.method == "POST":
        print(request.form)
        options: str = request.form.get("options")
        if options is not None:
            options = {"Options": json.loads(options)}
            GouelHelper(g.uga).update_event(event_id, options)
            flash("Les options ont été modifiées.", "success")
            return redirect(url_for("admin.manage_options", event_id=event_id))

    return render_template(
        "pages/admin/manage_options.j2",
        Options=g.event["Options"] or {},
        retour=url_for("admin.manage_event", event_id=event_id),
    )


@admin.route("/manage/<event_id>/store", methods=["POST", "GET"])
@login_required
@event_access
def manage_store(event_id):
    if request.method == "POST":
        print(request.form)
        action: str = request.form.get("action")
        if not action:
            return redirect(url_for("admin.manage_store", event_id=event_id))
        if action == "add":
            product = {
                "Label": request.form["Label"],
                "Price": float(request.form["Price"]),
                "HasAlcohol": request.form.get("HasAlcohol") == "on",
                "Icon": request.form.get("Icon") or "inventory_2",
            }

            end_of_sale = request.form.get("EndOfSale")
            if end_of_sale:
                product["EndOfSale"] = end_of_sale

            r = GouelHelper(g.uga).add_product(event_id, product)
            if r:
                flash("Le produit a été ajouté.", "success")
            else:
                flash("Le produit n'a pas été ajouté.", "error")
        elif action == "edit":
            product = {
                "Label": request.form["Label"],
                "Price": float(request.form["Price"]),
                "HasAlcohol": request.form.get("HasAlcohol") == "on",
                "Icon": request.form.get("Icon") or "inventory_2",
            }

            end_of_sale = request.form.get("EndOfSale")
            if end_of_sale:
                product["EndOfSale"] = end_of_sale

            product_id = request.form["ProductCode"]
            r = GouelHelper(g.uga).update_product(event_id, product_id, product)
            if r:
                flash("Le produit a été modifié.", "success")
            else:
                flash("Le produit n'a pas été modifié.", "error")
        elif action == "delete":
            product_id = request.form["ProductCode"]
            r = GouelHelper(g.uga).delete_product(event_id, product_id)
            if r:
                flash("Le produit a été supprimé.", "success")
            else:
                flash("Le produit n'a pas été supprimé.", "error")

        return redirect(url_for("admin.manage_store", event_id=event_id))

    return render_template(
        "pages/admin/manage_store.j2",
        produits=GouelHelper(g.uga).get_products(event_id),
        retour=url_for("admin.manage_event", event_id=event_id),
    )
