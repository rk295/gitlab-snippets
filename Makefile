.version:
	./scripts/version > .version

gitlab-snippets: *.go .version
	go build -o gitlab-snippets -ldflags "-X main.versionStr=$(shell cat ./.version)"

install: gitlab-snippets
	sudo cp gitlab-snippets /usr/local/bin/

tidy:
	@go mod tidy

clean:
	@rm .version gitlab-snippets || true ; echo "cleaned"