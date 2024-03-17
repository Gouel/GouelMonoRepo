FROM nginx
WORKDIR /
COPY ./release/env/nginx/cert.pem /etc/nginx
COPY ./release/env/nginx/priv.pem /etc/nginx
COPY ./release/env/nginx/default.conf /etc/nginx/conf.d


