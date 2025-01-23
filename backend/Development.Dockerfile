FROM golang:lts-alpine

RUN apk add --no-cache git && \
    go install github.com/xyproto/wgo@latest

WORKDIR /app/backend

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8000

CMD ["wgo", "run", "main.go"]
