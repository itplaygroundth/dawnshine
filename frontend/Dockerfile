FROM node:12.14.1-alpine as build-env
WORKDIR /app

COPY ./next-app/package*.json ./app
RUN yarn install
COPY ./ ./app
RUN yarn run dev


# FROM node:10.16.1-alpine

# WORKDIR /app
# COPY . /app

# # 別のビルドから成果物をコピー可能
# COPY --from=build-env /app/.nuxt /app/.nuxt

# # 本番では devDeps を消してファイルサイズを削減
# RUN yarn install --production
# EXPOSE 3000
# CMD ["yarn", "run","start"]
 