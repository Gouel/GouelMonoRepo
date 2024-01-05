from flask import Flask, render_template, request, redirect, url_for, session, Blueprint
from .generate_fake_data import db

admin = Blueprint("admin", __name__)


@admin.route("/admin/login", methods=["GET", "POST"])
def login():
    if request.method == "POST":
        username = request.form["username"]
        password = request.form["password"]

        passwd = {
            db["users"][i]["nom"]: i
            for i in range(len(db["users"]))
            if db["users"][i]["sessions"] != {}
        }

        # Ici, validez les identifiants de l'utilisateur
        if username in passwd:
            session["compte"] = db["users"][passwd[username]]
            return redirect(url_for("admin.dashboard"))
        else:
            # Gérer l'échec de la connexion
            pass

    return render_template("pages/admin/login.j2")


# Supposons que vous ayez une route pour le tableau de bord de l'administrateur
@admin.route("/")
@admin.route("/dashboard")
def dashboard():
    if "compte" not in session:
        return redirect(url_for("admin.login"))
    # Affichez le tableau de bord de l'administrateur
    return render_template(
        "pages/admin/dashboard.j2", user=session["compte"], all_events=db["events"]
    )


@admin.route("/manage/<int:event_id>", methods=["POST", "GET"])
def manage_event(event_id: int):
    if "compte" not in session:
        return redirect(url_for("admin.login"))

    if not (
        str(event_id) in session["compte"]["sessions"]
        or "superadmin" in session["compte"]["sessions"]
    ):
        return redirect(url_for("admin.dashboard"))

    return render_template(
        "pages/admin/manage_event.j2", event_id=event_id, event=db["events"][event_id]
    )


@admin.route("/manage/<int:event_id>/volunteers", methods=["POST", "GET"])
def manage_volunteers(event_id: int):
    if "compte" not in session:
        return redirect(url_for("admin.login"))

    if not (
        str(event_id) in session["compte"]["sessions"]
        or "superadmin" in session["compte"]["sessions"]
    ):
        print("superadmin" in session["compte"]["sessions"])
        return redirect(url_for("admin.dashboard"))

    return render_template(
        "pages/admin/manage_volunteers.j2", volunteers=db["benevoles"][event_id]
    )


@admin.route(
    "/manage/<int:event_id>/tickets", methods=["POST", "GET"], defaults={"page": 0}
)
@admin.route("/manage/<int:event_id>/tickets/<int:page>", methods=["POST", "GET"])
def manage_tickets(event_id: int, page: int = 0):
    if "compte" not in session:
        return redirect(url_for("admin.login"))

    if not (
        str(event_id) in session["compte"]["sessions"]
        or "superadmin" in session["compte"]["sessions"]
    ):
        return redirect(url_for("admin.dashboard"))

    return render_template(
        "pages/admin/manage_tickets.j2",
        tickets=db["tickets"],
        tickets_valides=len([i for i in db["tickets"] if i["valide"]]),
        tickets_total=len(db["tickets"]),
        total_prix_tickets=len(db["tickets"]) * 3,
        total_depense_credits=len(db["tickets"]) * 3.33,
        pagination={
            "has_prev": True,
            "has_next": True,
            "prev_num": 0,
            "next_num": 1,
        },
        event_id=event_id,
    )


@admin.route("/deconnexion")
def admin_deconnxion():
    if "compte" in session:
        session.pop("compte")
    return redirect(url_for("admin.login"))
