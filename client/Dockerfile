FROM node:20.16.0-alpine3.19
WORKDIR /app
COPY ./package*.json .
RUN npm ci
COPY . .
RUN npm run build
