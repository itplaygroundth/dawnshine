FROM node:17-alpine3.14

WORKDIR /usr/src/app
COPY package*.json /usr/src/app/
RUN apk add git && yarn
#RUN npm install -g npm@8.5.3
RUN yarn install 
COPY . .
EXPOSE ${PORT}
CMD ["yarn", "start"]