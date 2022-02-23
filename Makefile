build: *.go
	go build -o gitlab-snippets

tidy:
	go mod tidy