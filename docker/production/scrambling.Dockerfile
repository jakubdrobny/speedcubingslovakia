FROM node:22-alpine AS builder

WORKDIR /app

RUN npm install -g @yao-pkg/pkg

COPY package*.json ./

RUN npm ci --omit=dev

COPY . .

RUN pkg . --targets node22-alpine-arm64 --output /app/scrambling-service

FROM alpine:latest

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

COPY --from=builder /app/scrambling-service .

EXPOSE 3999

CMD ["./scrambling-service"]
