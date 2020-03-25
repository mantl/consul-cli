NAME = $(shell awk -F\" '/^const Name/ { print $$2 }' main.go)
VERSION = $(shell awk -F\" '/^const Version/ { print $$2 }' main.go)

all: build

build:
	@mkdir -p bin/
	go build -o bin/$(NAME)

deps:
	go mod download

test:
	go test ./... $(TESTARGS) -timeout=30s -parallel=4
	go vet  ./...

xcompile: test package-clean
	@mkdir -p build
	gox \
		-arch="arm amd64" \
		-os="darwin freebsd linux netbsd openbsd solarios windows" \
		-output="build/{{.Dir}}_$(VERSION)_{{.OS}}_{{.Arch}}/$(NAME)"

package: xcompile
	$(eval FILES := $(shell ls build))
	@mkdir -p build/zip
	for f in $(FILES); do \
		(cd $$PWD/build/$$f && zip ../zip/$$f.zip consul-cli*); \
		echo $$f; \
	done

package-clean:
	-@rm -rf build/

.PHONY: all deps updatedeps build test xcompile package package-clean
