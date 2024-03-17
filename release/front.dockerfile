FROM python:3.11
WORKDIR /code/app
COPY ./GouelFront ./
COPY ./release/env/front.env .env

RUN pip install --no-cache-dir --upgrade -r ./requirements.txt

CMD ["waitress-serve", "--port=5001", "run:app"]