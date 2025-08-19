FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-w -s" -o /app/server main/main.go

FROM alpine:latest

COPY --from=builder /app/server .

EXPOSE 8000

CMD ["./server"]
