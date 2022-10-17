FROM node:12.14.1-alpine as build-env
WORKDIR /app

COPY package*.json ./
RUN yarn install
COPY ./ ./
 
EXPOSE 3000
CMD ["yarn", "run","dev"]