FROM golang:1.24-alpine AS migrate

WORKDIR /app

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY migrations /app/migrations

COPY .env /app/.env

COPY migrate.sh /app/migrate.sh
RUN chmod +x /app/migrate.sh

CMD ["/app/migrate.sh"]
