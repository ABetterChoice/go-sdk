all: .format

.PHONY: all

.format:
	go generate ./...
	go test -parallel 1 -p 1 ./... -coverprofile=size_coverage.out -gcflags "all=-N -l"
	go tool cover -html=size_coverage.out
	rm -rf size_coverage.out
	go mod tidy
	golint ./...
	gofmt -w .
	goimports -w .
	go vet ./...
	gonote ./...