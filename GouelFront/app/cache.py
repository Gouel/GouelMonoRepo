import redis
from uuid import uuid4
import os


class Cache:
    def __init__(self, host="localhost", port=6379, db=0):
        self.r = redis.from_url("redis://" + os.getenv("REDIS_URL"))

    def set(self, key, value, ttl=None):
        """Définit une clé avec une valeur. Optionnellement, définir un TTL (Time To Live) en secondes."""
        if ttl:
            self.r.setex(key, ttl, value)
        else:
            self.r.set(key, value)

    def get(self, key):
        """Obtient la valeur d'une clé."""
        return self.r.get(key)

    def delete(self, key):
        """Supprime une clé."""
        return self.r.delete(key)


class MagicLink:
    def __init__(self, type: str, obj: dict, ttl: int = 60):
        self.type = type
        self.obj = obj
        self.ttl = ttl
        self.id = str(uuid4())

        c = Cache()
        c.set(self.id, self, ttl)

    def __str__(self):
        return str({"type": self.type, "obj": self.obj, "ttl": self.ttl, "id": self.id})

    def __repr__(self):
        return f"MagicLink[{self.id}]({self.type}, {self.obj}, {self.ttl})"
