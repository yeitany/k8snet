FROM debian:stretch-slim

WORKDIR /
COPY ./.bin ./
COPY ./.env /.env

CMD ["/bin/sh"]