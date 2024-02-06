FROM golang:1.21.6-alpine

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk update && \
    apk add build-base sqlite-dev sqlite-libs && \
    cd /usr/local/bin && wget https://release.ariga.io/atlas/atlas-linux-amd64-latest && \
    mv atlas-linux-amd64-latest atlas && chmod +x atlas

RUN CGO_CFLAGS="-D_LARGEFILE64_SOURCE" go install -tags sqlite github.com/gobuffalo/pop/v6/soda@latest


ENV GO_ENV=production \
    GIN_MODE=release

