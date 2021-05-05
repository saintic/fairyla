# build dependencies with alpine
FROM golang:1.16.3-alpine3.13 AS builder

WORKDIR /build

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o fairyla

# run application with a small image
FROM scratch

WORKDIR /fairyla

COPY --from=builder /build/fairyla /bin/

COPY --from=builder /build/ui /fairyla/

EXPOSE 10210

CMD ["-v"]

ENTRYPOINT ["fairyla"]
