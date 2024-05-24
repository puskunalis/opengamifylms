FROM node:21.7.3-alpine3.19 AS build

WORKDIR /app

COPY package*.json ./

RUN npm ci

COPY public ./public
COPY src ./src

RUN npm run build

FROM nginx:1.25.5-alpine3.19-slim

EXPOSE 80

COPY --from=build /app/build /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/default.conf

CMD ["nginx", "-g", "daemon off;"]
