all: test build

fmt:
	go fmt ./...

vet:
	go vet ./...

GOTESTOUTFILE?=cover.out

test: fmt vet
	go test -race -coverprofile=${GOTESTOUTFILE} ./...

webtest: test
	go tool cover -html=${GOTESTOUTFILE}

functest: test
	go tool cover -func=${GOTESTOUTFILE}

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