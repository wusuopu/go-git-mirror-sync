FROM golang:1.21.6-alpine as builder

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk update && \
    apk add build-base

COPY ./src /app
WORKDIR /app/

RUN go build -o goose cmd/goose.go && \
    go build -tags=jsoniter -ldflags "-linkmode external -extldflags=-static -s -w" -o app .

FROM alpine:3.19
COPY --from=builder /app/ /app
WORKDIR /app/

ENV GO_ENV=production \
    GIN_MODE=release

ENTRYPOINT ["/app/run.sh"]

CMD ["start_server"]
