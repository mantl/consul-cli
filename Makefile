TEST?=$$(glide nv)
NAME = $(shell awk -F\" '/^const Name/ { print $$2 }' main.go)
VERSION = $(shell awk -F\" '/^const Version/ { print $$2 }' main.go)
DEPS = $(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

all: build

build: 
	@mkdir -p bin/
	go build -o bin/$(NAME)

deps:
	glide up

test: 
	go test $(TEST) $(TESTARGS) -timeout=30s -parallel=4
	go vet $(TEST)

xcompile: test
	@rm -rf build/
	@mkdir -p build
	gox \
		-os="darwin" \
		-os="freebsd" \
		-os="linux" \
		-os="netbsd" \
		-os="openbsd" \
		-os="solaris" \
		-os="windows" \
		-output="build/{{.Dir}}_$(VERSION)_{{.OS}}_{{.Arch}}/$(NAME)" 

vendor:
	glide install --strip-vendor
	glide update --strip-vendor

vendor-clean:
	-rm -rf vendor/

package: xcompile
	$(eval FILES := $(shell ls build))
	@mkdir -p build/zip
	for f in $(FILES); do \
		(cd $(shell pwd)/build/$$f && zip ../zip/$$f.zip consul-cli*); \
		echo $$f; \
	done

package-clean:
	-rm -rf build/

.PHONY: all deps updatedeps build test xcompile package package-clean vendor vendor-clean
