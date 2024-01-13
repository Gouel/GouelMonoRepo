# Gestion de l'api de Gouel

import os
from time import time
import requests
from flask import abort


class GouelApi(object):
    """Gouel API"""

    # Il faut username et password pour l'api
    def __init__(self, username, password):
        self.username = username
        self.password = password
        self.base_url = os.getenv("GOUEL_SERVER")
        self.token = {"token": None}

    def get_token(self) -> str:
        if self.token["token"] is None or self.token["infos"]["exp"] < time():
            r = requests.post(
                self.base_url + "/token/auth",
                json={"email": self.username, "password": self.password},
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

        return self.token["token"]

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
