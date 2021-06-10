ARG buildos=golang:1.16.4-alpine3.13
ARG runos=scratch

# build dependencies with alpine
FROM $buildos AS builder

LABEL maintainer=me@tcw.im

WORKDIR /build

COPY . .

ARG alpine=mirrors.tuna.tsinghua.edu.cn

ARG registry=https://registry.npm.taobao.org

ARG goproxy=https://goproxy.cn

RUN sed -i "s/dl-cdn.alpinelinux.org/${alpine}/g" /etc/apk/repositories && \
    apk add --no-cache nodejs yarn && \
    yarn config set registry $registry && \
    go env -w GOPROXY=${goproxy},direct

RUN cd client && yarn --no-lockfile && yarn build && cd ../server && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o fairyla

# run application with a small image
FROM $runos

WORKDIR /fairyla

COPY --from=builder /build/server/fairyla /bin/

COPY --from=builder /build/server/ui /build/NOTICE /build/LICENSE /fairyla/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 10210

ENV fairyla_dir="/fairyla"

ENTRYPOINT ["fairyla"]
