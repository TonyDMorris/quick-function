FROM --platform=linux/amd64  node:20



COPY ./package.json ./
COPY ./yarn.lock ./

RUN yarn install --network-timeout 1000000000

COPY . .

RUN yarn build 

EXPOSE 1337

CMD ["yarn", "start"]