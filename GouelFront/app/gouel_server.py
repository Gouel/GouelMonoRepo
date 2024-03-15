# Gestion de l'api de Gouel

import os
from time import time
import requests
from flask import abort


class GouelApi(object):
    """Gouel API"""

    # Il faut username et password pour l'api
    def __init__(self, username: str = None, password: str = None, token: dict = None):
        self.base_url = os.getenv("GOUEL_SERVER")
        self.token = token or {"token": None}
        if username and password:
            self.get_first_token(username, password)
        elif token:
            self.refresh_token()

    def get_first_token(self, username, password):
        r = requests.post(
            self.base_url + "/token/auth",
            json={"email": username, "password": password},
        )
        if r.status_code == 200:
            self.token = {"token": r.json()["token"]}
            r = requests.get(
                self.base_url + "/token/view",
                headers={"Authorization": "Bearer " + self.token["token"]},
            )
            self.token["infos"] = r.json()
        else:
            self.token = {"token": None}

    def get_token(self) -> str:
        if self.token["token"] is None or self.token["infos"]["exp"] < time():
            self.refresh_token()
        return self.token["token"]

    def refresh_token(self):
        if self.token["token"] is not None:
            r = requests.post(
                self.base_url + "/token/refresh",
                headers={"Authorization": "Bearer " + self.token["token"]},
            )
            if r.status_code == 200:
                self.token = {"token": r.json()["token"]}
                r = requests.get(
                    self.base_url + "/token/view",
                    headers={"Authorization": "Bearer " + self.token["token"]},
                )
                self.token["infos"] = r.json()
        else:
            self.token = {"token": None}

    # headers
    def get_headers(self) -> dict:
        return {"Authorization": "Bearer " + self.get_token()}

    # post
    def post(self, path, data):
        r = requests.post(
            self.base_url + path,
            json=data,
            headers=self.get_headers(),
        )
        return r

    # get
    def get(self, path):
        r = requests.get(
            self.base_url + path,
            headers=self.get_headers(),
        )
        return r

    # put
    def put(self, path, data):
        r = requests.put(
            self.base_url + path,
            json=data,
            headers=self.get_headers(),
        )
        return r

    # delete
    def delete(self, path):
        r = requests.delete(
            self.base_url + path,
            headers=self.get_headers(),
        )
        return r

    def to_json(self):
        return vars(self)

    @staticmethod
    def from_json(json):
        api = GouelApi(token=json["token"])
        return api
