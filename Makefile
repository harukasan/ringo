gobin = $(GOPATH)/bin

get-deps: $(gobin)/stringer

generate:
	go generate ./...

test: generate
	go test -v ./...

cover: generate $(gobin)/gover gover.coverprofile

cover-html: cover
	go tool cover -html=gover.coverprofile
	@rm gover.coverprofile

coveralls: cover $(gobin)/goveralls
	$(gobin)/goveralls -coverprofile=gover.coverprofile -service=travis-ci

gover.coverprofile: $(gobin)/gover
	./gover.sh

$(gobin)/gover:
	go get github.com/modocache/gover

$(gobin)/goveralls:
	go get github.com/mattn/goveralls

$(gobin)/stringer:
	go get golang.org/x/tools/cmd/stringer

.PHONY: generate \
	test \
	cover \
	cover-html \
	coveralls
