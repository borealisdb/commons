codegen:
	hack/update-codegen.sh

test:
	go test ./...

go-clean:
	go mod tidy
build: go-clean codegen

install-dep:
	go install github.com/goreleaser/goreleaser@latest
	go install github.com/caarlos0/svu@latest

release:
	git tag $(shell svu next)
	git push --tags
	goreleaser release --clean