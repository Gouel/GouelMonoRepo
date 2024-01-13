from flask import abort
from . import ga


class GouelHelper:
    @staticmethod
    def get_events():
        return ga.get("/events").json()

    @staticmethod
    def get_event(event_id):
        r = ga.get(f"/events/{event_id}")
        event = r.json() if r.status_code == 200 else None
        if event is None:
            abort(404)
        return event

    @staticmethod
    def get_user(user_id):
        r = ga.get(f"/users/{user_id}")
        user = r.json() if r.status_code == 200 else None
        if user is None:
            abort(404)
        return user

    def get_products(event_id):
        return ga.get(f"/events/{event_id}/products").json()

    def get_product_from_code(event_id, code):
        r = ga.get(f"/events/{event_id}/products/{code}")
        product = r.json() if r.status_code == 200 else None
        if product is None:
            abort(404)
        return product
