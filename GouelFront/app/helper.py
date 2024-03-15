from flask import abort
from .gouel_server import GouelApi


class GouelHelper:
    def __init__(self, ga: GouelApi):
        self.ga = ga

    @staticmethod
    def generate_password():
        import random
        import string

        charset: str = string.printable[:75]
        return "".join(random.choice(charset) for _ in range(12))

    # API

    def get_events(self, to_dict: bool = False):
        events = self.ga.get("/events").json()
        if not to_dict:
            return events

        return {e["ID"]: e for e in events}

    def get_event(self, event_id):
        r = self.ga.get(f"/events/{event_id}")
        event = r.json() if r.status_code == 200 else None
        if event is None:
            abort(404)
        return event

    def get_event_smtp(self, event_id):
        r = self.ga.get(f"/events/{event_id}/smtp")
        smtp = r.json() if r.status_code == 200 else None
        return smtp

    def add_event(self, data):
        r = self.ga.post("/events", data)
        return r.status_code == 200

    def add_event_ticket(self, event_id, data):
        r = self.ga.post(f"/events/{event_id}/tickets", data)
        return r.status_code == 200

    def update_event_ticket(self, event_id, ticket_id, data):
        r = self.ga.put(f"/events/{event_id}/tickets/{ticket_id}", data)
        return r.status_code == 200

    def delete_event_ticket(self, event_id, ticket_id):
        r = self.ga.delete(f"/events/{event_id}/tickets/{ticket_id}")
        return r.status_code == 200

    def delete_event(self, event_id):
        r = self.ga.delete(f"/events/{event_id}")
        return r.status_code == 200

    def update_event(self, event_id, data):
        r = self.ga.put(f"/events/{event_id}", data)
        return r.status_code == 200

    def add_volunteer(self, event_id, data):
        r = self.ga.post(f"/events/{event_id}/volunteers", data)
        return r.status_code == 200

    def update_volunteer(self, event_id, data):
        r = self.ga.put(f"/events/{event_id}/volunteers", data)
        return r.status_code == 200

    def delete_volunteer(self, event_id, user_id):
        r = self.ga.delete(f"/events/{event_id}/volunteers/{user_id}")
        return r.status_code == 200

    def get_user(self, user_id: str, email: str = None):
        if email is not None:
            r = self.ga.get(f"/users/email/{email}")
        else:
            r = self.ga.get(f"/users/{user_id}")
        user = r.json() if r.status_code == 200 else None
        return user

    def add_user(self, data):
        r = self.ga.post("/users", data)
        return r.status_code == 200, r.json()

    def get_products(self, event_id):
        return self.ga.get(f"/events/{event_id}/products").json()

    def add_product(self, event_id, data):
        r = self.ga.post(f"/events/{event_id}/products", data)
        return r.status_code == 200

    def update_product(self, event_id, product_id, data):
        r = self.ga.put(f"/events/{event_id}/products/{product_id}", data)
        return r.status_code == 200

    def delete_product(self, event_id, product_id):
        r = self.ga.delete(f"/events/{event_id}/products/{product_id}")
        return r.status_code == 200

    # tickets

    def get_tickets(self, event_id):
        return self.ga.get(f"/tickets/{event_id}").json()
