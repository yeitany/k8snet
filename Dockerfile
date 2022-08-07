FROM debian:stretch-slim
RUN docker pull hamroctopus/graphviz
FROM efcasado/graphviz

WORKDIR /
COPY ./.bin ./
#COPY ./.env /.env

CMD ["/bin/sh"]