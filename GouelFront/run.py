from app import create_app
from dotenv import load_dotenv
import os
import redis

load_dotenv()

app = create_app()

app.secret_key = os.getenv("APP_SECRET", default="BAD_SECRET_KEY")

app.config["SESSION_TYPE"] = "redis"
app.config["SESSION_PERMANENT"] = True
app.config["PERMANENT_SESSION_LIFETIME"] = 60 * 30  # 30 minutes
app.config["SESSION_USE_SIGNER"] = True
app.config["SESSION_REDIS"] = redis.from_url("redis://" + os.getenv("REDIS_URL"))


if __name__ == "__main__":
    host = os.getenv("APP_HOST", "0.0.0.0")
    port = int(os.getenv("APP_PORT", "5001"))
    debug = os.getenv("APP_DEBUG", "False") == "True"
    app.run(host=host, port=port, debug=debug)
