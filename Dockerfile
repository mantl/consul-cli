FROM golang:1.13-alpine AS builder

RUN apk --no-cache add git make

WORKDIR /src/consul-cli
COPY . .
RUN CGO_ENABLED=0 make build

FROM busybox
LABEL maintainer "Chris Aubuchon <Chris.Aubuchon@gmail.com>"

ENTRYPOINT ["consul-cli"]

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /src/consul-cli/bin/consul-cli /usr/local/bin/
