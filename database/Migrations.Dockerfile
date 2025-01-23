FROM golang:1.23-alpine

RUN apk add --no-cache git make && \
  go install -tags 'postgres' -v github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app/database

COPY migrations Makefile .

CMD ["sh", "-c", "sleep 5 && make migrate_up"]
