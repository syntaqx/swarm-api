export GO111MODULE=on

run:
	go run main.go

tidy:
	go mod tidy

build/test:
	goreleaser --snapshot --skip-publish --rm-dist
