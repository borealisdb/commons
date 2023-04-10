codegen:
	hack/update-codegen.sh

test:
	go test ./...

go-clean:
	go mod tidy
build: go-clean codegen

install-dep:
	go install github.com/goreleaser/goreleaser@latest

release:
	goreleaser release