ARG buildos=golang:1.17-alpine
ARG runos=scratch

# --build dependencies with alpine--
FROM $buildos AS builder
WORKDIR /build
COPY . .
ARG goproxy
ARG registry
ARG TARGETARCH
RUN if [ "x$goproxy" != "x" ]; then go env -w GOPROXY=${goproxy},direct; fi ;\
    if [ "x$registry" != "x" ]; then yarn config set registry $registry; fi ; \
    apk add --no-cache nodejs yarn &&\
    cd client && yarn --no-lockfile && yarn build &&\
    cd ../server && \
    CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -ldflags "-s -w" .

# --run application with a small image--
FROM $runos
WORKDIR /fairyla
COPY --from=builder /build/server/fairyla /bin/
COPY --from=builder /build/server/ui /build/NOTICE /build/LICENSE /fairyla/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 10210
ENV fairyla_dir="/fairyla"
ENTRYPOINT ["fairyla"]
