FROM node:20-alpine

WORKDIR /app

COPY package.json .
RUN npm install

COPY mock-bff.js .

EXPOSE 4000

CMD ["node", "mock-bff.js"]
