FROM node:21-alpine3.18 AS frontend-builder
WORKDIR /app
COPY frontend/ /app/
RUN apk update \
    && apk add curl \
    && npm install -g pnpm \
    && pnpm i \
    && pnpm build

FROM golang:1.21 AS backend-builder
WORKDIR /app
COPY . /app/
COPY --from=frontend-builder /app/dist /app/frontend/dist
COPY --from=frontend-builder /app/embed.go /app/frontend/embed.go
RUN go mod tidy
RUN CGO_ENABLED=0  go build -o apicat-server ./cmd/app/

FROM alpine:latest
WORKDIR /app
COPY ./wait-for-it.sh /app/wait-for-it.sh
COPY --from=backend-builder /app/apicat-server /app/apicat-server
RUN apk add bash \
    && chmod +x /app/wait-for-it.sh
ENTRYPOINT ["/app/apicat-server"]