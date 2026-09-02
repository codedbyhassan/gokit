.PHONY: test test-race vet fmt check

test:
	go test ./...

test-race:
	go test -race ./...

vet:
	go vet ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './.git/*')

check: fmt vet test
