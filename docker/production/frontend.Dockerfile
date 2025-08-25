FROM node:22-slim AS builder

WORKDIR /app

COPY package*.json ./

RUN npm ci

COPY . .

RUN npm run build

FROM nginx:1.29-alpine-slim

COPY --from=builder /app/dist /usr/share/nginx/html
