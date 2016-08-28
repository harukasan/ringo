test:
	go generate ./...
	go test -v ./...

get-deps: $(GOPATH)/bin/stringer

$(GOPATH)/bin/stringer:
	go get golang.org/x/tools/cmd/stringer

.PHONY: test
