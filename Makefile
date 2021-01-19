all: test build

test:
	go vet ./...
	go test ./...

example:
	./example/run.sh

gox:
	go get github.com/mitchellh/gox
	gox ${LDFLAGS} -parallel=4 -output="./bin/gotmplx_{{.OS}}_{{.Arch}}"

clean:
	rm -rfv ./gotmplx

build:
	go build

.PHONY: all test gox build clean example