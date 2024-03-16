import requests as rq
from urllib.parse import urlencode
import json
import time
from datetime import datetime


class HelloAssoAPI:
    HELLO_ASSO_API_URL = "https://api.helloasso.com"
    SANDBOX_HELLO_ASSO_API_URL = "https://api.helloasso-sandbox.com"

    def __init__(
        self, clientId: str, clientSecret: str, slug: str, sandbox: bool = True
    ):
        self.clientId = clientId
        self.clientSecret = clientSecret
        self.slug = slug
        self.token = None
        self.url = (
            HelloAssoAPI.HELLO_ASSO_API_URL
            if not sandbox
            else HelloAssoAPI.SANDBOX_HELLO_ASSO_API_URL
        )

    def getToken(self) -> dict:
        if self.token is not None and self.token["expires_on"] > time.time():
            return self.token

        data = {
            "client_secret": self.clientSecret,
            "client_id": self.clientId,
            "grant_type": "client_credentials",
        }

        headers = {"Content-Type": "application/x-www-form-urlencoded"}

        r = rq.post(self.url + "/oauth2/token", data=urlencode(data), headers=headers)
        if r.status_code != 200:
            self.token = None
            raise ValueError("Problem with token : ", r, r.text)

        self.token = r.json()
        self.token["expires_on"] = time.time() + self.token["expires_in"]

        return self.token

    def __str__(self):
        return str(
            {
                "token": self.token,
                "clientId": self.clientId,
                "clientSecret": self.clientSecret,
                "slug": self.slug,
                "url": self.url,
            }
        )

    def __repr__(self):
        return str(self)

    def addCheckOutOrder(self, order: "CheckoutOrder") -> dict:
        # récupération du token
        self.getToken()

        # préparation de la requête
        headers = {"Authorization": f"Bearer {self.token['access_token']}"}
        r = rq.post(
            self.url + f"/v5/organizations/{self.slug}/checkout-intents",
            json=order.to_dict(),
            headers=headers,
        )

        return r.json()

    def getCheckOutOrder(self, id: str) -> dict:
        # récupération du token
        self.getToken()

        # préparation de la requête
        headers = {"Authorization": f"Bearer {self.token['access_token']}"}
        r = rq.get(
            self.url + f"/v5/organizations/{self.slug}/checkout-intents/{id}",
            headers=headers,
        )

        return r.json()


class Payer:
    def __init__(self, firstName: str, lastName: str, email: str, dob: datetime):
        self.firstName = firstName
        self.lastName = lastName
        self.email = email
        self.dob = dob

    def to_dict(self) -> dict:
        return {
            "firstName": self.firstName,
            "lastName": self.lastName,
            "email": self.email,
            "dateOfBirth": self.dob.strftime("%Y-%m-%d") + "T00:00",
        }

    def __repr__(self):
        return str(self)

    def __str__(self):
        return json.dumps(self.to_dict())


class CheckoutOrder:
    def __init__(
        self,
        totalAmount: int,
        initialAmount: int,
        itemName: str,
        backUrl: str,
        errorUrl: str,
        returnUrl: str,
        payer: Payer,
        containsDonation: bool = False,
        meta: dict = None,
    ):
        self.meta = meta if meta is not None else {}
        self.totalAmount = int(totalAmount)
        self.initialAmount = int(initialAmount)
        self.itemName = itemName
        self.backUrl = backUrl
        self.errorUrl = errorUrl
        self.returnUrl = returnUrl
        self.payer = payer
        self.containsDonation = containsDonation

    def to_dict(self) -> dict:
        meta = {
            "from": "gouel_front",
            "to": "helloAsso",
            "using": "helloAssoApi",
        }
        meta.update(self.meta)
        return {
            "totalAmount": self.totalAmount,
            "initialAmount": self.initialAmount,
            "itemName": self.itemName,
            "backUrl": self.backUrl,
            "errorUrl": self.errorUrl,
            "returnUrl": self.returnUrl,
            "payer": self.payer.to_dict(),
            "containsDonation": self.containsDonation,
            "meta": meta,
        }

    def __repr__(self):
        return str(self)

    def __str__(self):
        return json.dumps(self.to_dict())


if __name__ == "__main__":
    from os import environ
    from dotenv import load_dotenv

    load_dotenv()

    ha = HelloAssoAPI(
        clientId=environ["HELLO_ASSO_CLIENT_ID"],
        clientSecret=environ["HELLO_ASSO_CLIENT_SECRET"],
        slug="inter-asso",
    )

    payer = Payer(
        firstName="Matthias",
        lastName="Hartmann",
        email="test@iziram.fr",
    )

    chk = CheckoutOrder(
        itemName="Ticket Up",
        initialAmount=3 * 100,
        totalAmount=3 * 100,
        backUrl="https://localhost:5000/acheter-billets/1",
        errorUrl="https://localhost:5000/paymentResponse?action=error",
        returnUrl="https://localhost:5000/paymentResponse?action=sucess",
        containsDonation=False,
        payer=payer,
    )

    print(payer)
    print(chk)

    print("===")

    c = ha.addCheckOutOrder(chk)

    print(c)
