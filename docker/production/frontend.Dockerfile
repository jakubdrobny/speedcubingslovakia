FROM node:22-slim AS builder

WORKDIR /app

COPY package*.json ./

RUN npm ci --omit=dev

FROM nginx:1.29-alpine-slim

COPY --from=builder /app/build /usr/share/nginx/html
