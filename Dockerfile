# build dependencies with alpine
FROM golang:1.16.3-alpine3.13 AS builder

LABEL maintainer=me@tcw.im

WORKDIR /build

COPY . .

ARG alpine=mirrors.tuna.tsinghua.edu.cn

ARG registry=https://registry.npm.taobao.org

ARG goproxy=https://goproxy.cn,direct

RUN sed -i "s/dl-cdn.alpinelinux.org/${alpine}/g" /etc/apk/repositories && \
    apk add --no-cache nodejs yarn && \
    yarn config set registry $registry && \
    go env -w GOPROXY=${goproxy}

RUN cd client && yarn --no-lockfile && yarn build && cd ../server && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o fairyla

# run application with a small image
FROM scratch

WORKDIR /fairyla

COPY --from=builder /build/server/fairyla /bin/

COPY --from=builder /build/server/ui /fairyla/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 10210

CMD ["-dir"， "/fairyla"]

ENTRYPOINT ["fairyla"]
