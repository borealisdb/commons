codegen:
	hack/update-codegen.sh

test:
	go test ./...

go-clean:
	go mod tidy
build: go-clean codegen

release:
	git tag $(shell svu next)
	git push --tags
	goreleaser release --clean