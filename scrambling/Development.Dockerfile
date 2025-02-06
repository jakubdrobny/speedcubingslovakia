FROM node:22-alpine

WORKDIR /app/scrambling

COPY package*.json ./

COPY . .

EXPOSE 3000

CMD [ "sh" , "-c", "npm install && npm start" ]
