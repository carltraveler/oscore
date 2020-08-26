FROM ubuntu:18.04

WORKDIR /app

COPY ./oscoreapi /app/

CMD ["/app/oscoreapi", "--config", "/appconfig/config.json", "--loglevel", "1"]
