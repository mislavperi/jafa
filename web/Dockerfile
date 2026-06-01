FROM node:22-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

# HAProxy + busybox httpd, supervised by tini.
FROM haproxy:2.9-alpine
USER root
RUN apk add --no-cache busybox-extras tini
COPY --from=builder /app/dist /var/www/html
COPY haproxy.cfg /usr/local/etc/haproxy/haproxy.cfg
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
EXPOSE 80
ENTRYPOINT ["/sbin/tini", "--", "/entrypoint.sh"]
