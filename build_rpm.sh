#!/bin/sh

# set go env vars. Assumes go has beeen installed to /usr/local per https://golang.org/doc/install#install
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$WORKSPACE

# checkout latest and build
go get github.com/CiscoCloud/consul-cli

# create RPM package using FPM https://github.com/jordansissel/fpm
/usr/local/bin/fpm --url "https://github.com/CiscoCloud/consul-cli" \
  --description "Command line interface to Consul HTTP API" \
  --license "Apache 2.0" --verbose --force \
  -s dir -t rpm -n "consul-cli" $GOPATH/bin/=/usr/local/bin

