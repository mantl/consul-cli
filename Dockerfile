FROM golang:alpine

MAINTAINER Chris Aubuchon <Chris.Aubuchon@gmail.com>

COPY . /go/src/github.com/CiscoCloud/consul-cli
RUN apk add --update go git mercurial \
	&& cd /go/src/github.com/CiscoCloud/consul-cli \
	&& export GOPATH=/go \
	&& go get \
	&& go build -o /bin/consul-cli \
	&& rm -rf /go \
	&& apk del --purge go git mercurial

ENTRYPOINT [ "/bin/consul-cli" ]
