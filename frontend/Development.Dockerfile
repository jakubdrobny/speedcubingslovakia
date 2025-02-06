FROM node:22-alpine

WORKDIR /app/frontend

COPY package*.json ./

COPY . .

EXPOSE 3000

CMD [ "sh" , "-c", "npm install && npm run dev -- --host --port 3000" ]
