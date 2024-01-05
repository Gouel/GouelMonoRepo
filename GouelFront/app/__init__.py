from flask import Flask, render_template
from .hello_asso import HelloAssoAPI
import os
from dotenv import load_dotenv
from datetime import datetime

load_dotenv()

ha = HelloAssoAPI(
    clientId=os.environ["HELLO_ASSO_CLIENT_ID"],
    clientSecret=os.environ["HELLO_ASSO_CLIENT_SECRET"],
    slug=os.environ["HELLO_ASSO_SLUG"],
    sandbox=os.getenv("HELLO_ASSO_SANDBOX", "True") == "True",
)


def create_app():
    app = Flask(__name__)

    @app.template_filter("dateformat")
    def dateformat(value, format="%d/%m/%Y %H:%M"):
        if value:
            return datetime.fromisoformat(value).strftime(format)
        return value

    # Importer et enregistrer les vues ici
    from .views import main
    from .api import api
    from .admin import admin

    app.register_blueprint(main)
    app.register_blueprint(api, url_prefix="/api")
    app.register_blueprint(admin, url_prefix="/admin")

    @app.errorhandler(404)
    def page_not_found(e):
        # notez que nous définissons la réponse 404 explicitement
        return render_template("404.j2"), 404

    return app
