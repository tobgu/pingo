build:
	go build github.com/tobgu/pingo/cmd/pingo/

install:
	go install github.com/tobgu/pingo/cmd/pingo/

fmt:
	go fmt ./...

init-dep:
	go get -u github.com/kardianos/govendor
	govendor sync

.PHONY: pingo