FROM node:20 AS build-stage

WORKDIR /usr/src/app

COPY . /usr/src/app

RUN npm cache clean --force
RUN npm ci --omit=optional
RUN npm run build

FROM nginxinc/nginx-unprivileged

ENV API_URL="http://localhost:4000"
ENV NGINX_ENVSUBST_FILTER="API_URL"

COPY --from=build-stage /usr/src/app/dist/ "/usr/share/nginx/html"
COPY --from=build-stage /usr/src/app/nginx.conf "/etc/nginx/templates/default.conf.template"
