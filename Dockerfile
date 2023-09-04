FROM golang:alpine AS builder
WORKDIR /app
COPY . .
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOPROXY=https://goproxy.cn
RUN apk add --no-cache libc-dev gcc && \
    go mod download && \ 
    go build --tags -ldflags

# 使用 alpine 减小镜像体积
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app /app

ENTRYPOINT ["/app/wxbot"]