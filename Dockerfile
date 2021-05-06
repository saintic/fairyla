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

RUN cd client && yarn --no-lockfile && yarn build && \
    cd ../server && go build -ldflags "-s -w" -o bin/fairyla

# run application with a small image
FROM scratch

WORKDIR /fairyla

COPY --from=builder /build/server/bin/fairyla /bin/

COPY --from=builder /build/server/ui /fairyla/

EXPOSE 10210

CMD ["-v"]

ENTRYPOINT ["fairyla"]
