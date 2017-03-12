build:
	go build github.com/tobgu/pingo/cmd/pingo/

install:
	go install github.com/tobgu/pingo/cmd/pingo/

fmt:
	go fmt ./...

save-dep:
	godep save ./...

.PHONY: pingo