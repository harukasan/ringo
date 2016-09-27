gobin = $(GOPATH)/bin
deps = $(gobin)/stringer
cover_file = gover.coverprofile

get-deps: $(gobin)/gover $(gobin)/goveralls $(gobin)/stringer

test: $(deps)
	go generate ./...
	go test -v ./...

cover-func: $(cover_file)
	go tool cover -func=$<

cover-html: $(cover_file)
	go tool cover -html=$<
	@rm $<

coveralls: $(cover_file) $(gobin)/goveralls
	$(gobin)/goveralls -coverprofile=$< -service=travis-ci

gover.coverprofile: $(gobin)/gover
	./gover.sh

$(gobin)/gover:
	go get github.com/modocache/gover

$(gobin)/goveralls:
	go get github.com/mattn/goveralls

$(gobin)/stringer:
	go get golang.org/x/tools/cmd/stringer

.PHONY: test \
	cover-func \
	cover-html \
	coveralls
