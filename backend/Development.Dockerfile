FROM golang:1.23-alpine

RUN apk add --no-cache git && \
 go install github.com/bokwoon95/wgo@latest

WORKDIR /app/backend

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8000

CMD ["wgo", "run", "main.go"]
