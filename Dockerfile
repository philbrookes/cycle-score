FROM golang:1.8.0-alpine

RUN mkdir -p /usr/src/web/public

EXPOSE 8080

COPY ./web/public /usr/src/web/public

COPY ./cmd/rest-api/rest-api /usr/src/app

CMD ["/usr/src/app"]