FROM golang:1.25-alpine

WORKDIR /app

RUN go install github.com/bokwoon95/wgo@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8000

CMD ["wgo", "run", "main/main.go"]
