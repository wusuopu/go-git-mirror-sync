FROM golang:1.21.6-alpine as builder

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk update && \
    apk add build-base

COPY ./src /app
WORKDIR /app/

ENV GOPROXY=https://goproxy.cn,direct \
    GO111MODULE=auto

RUN go build -ldflags "-linkmode external -extldflags=-static -s -w" -o goose cmd/goose.go && \
    go build -tags=jsoniter -ldflags "-linkmode external -extldflags=-static -s -w" -o app .

FROM alpine:3.19
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk update && \
    apk add git

COPY --from=builder /app/ /app
WORKDIR /app/

VOLUME ["/app/tmp", "/app/data"]

ENV GO_ENV=production \
    GIN_MODE=release \
    BASIC_AUTH_USER= \
    BASIC_AUTH_PASSWORD= \
    DATABASE_TYPE=sqlite \
    DATABASE_DSN=/app/data/production.db \
    GIT_INSECURE_SKIP_TLS=false \
    CRONTAB="0 */6 * * *"

ENTRYPOINT ["/app/run.sh"]

CMD ["start_server"]
