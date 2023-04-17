codegen:
	hack/update-codegen.sh

test:
	go test ./...

go-clean:
	go mod tidy
build: go-clean codegen

check.github.token:
ifndef GITHUB_TOKEN
	$(error GITHUB_TOKEN is undefined)
endif

release: check.github.token
	git tag $(shell svu next)
	git push --tags
	goreleaser release --clean