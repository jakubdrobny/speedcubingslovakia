FROM node:22-slim AS builder

WORKDIR /app

COPY ./frontend/package*.json ./

RUN npm ci

ARG VITE_WCA_GET_CODE_URL
ARG VITE_SCRAMBLE_IMAGES_PATH
ARG NODE_ENV

COPY ./frontend .

RUN npm run build

FROM nginx:1.29-alpine-slim

COPY --from=builder /app/dist /usr/share/nginx/html

COPY ./configs/nginx/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
EXPOSE 443
