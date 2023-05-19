FROM node:18-alpine3.16 AS frontend-builder
WORKDIR /app
COPY frontend/ /app/
RUN ls
RUN apk update \
    && apk add curl \
    && npm config set registry https://registry.npmmirror.com \
    && curl -f https://get.pnpm.io/v6.16.js | node - add --global pnpm turbo-cli \
    && SHELL=bash pnpm setup \
    && source /root/.bashrc \
    && pnpm install \ 
    && pnpm build


FROM golang:1.20 AS backend-builder
WORKDIR /app
ENV GORPOXY="https://goproxy.cn,direct"
COPY . /app/
COPY --from=frontend-builder /app/dist /app/frontend/dist 
COPY --from=frontend-builder /app/embed.go /app/frontend/embed.go 
RUN go mod tidy && CGO_ENABLED=0  go build -o apicat-server .


FROM alpine:3.18 
WORKDIR /app 
COPY config/setting.default.yaml ./
COPY --from=backend-builder /app/apicat-server /app/apicat-server 
EXPOSE 8000
ENTRYPOINT ["/app/apicat-server"]