FROM nginx

# Copie de l'application web dans le dossier de l'application de nginx
COPY ./release/webapp/ /usr/share/nginx/html

