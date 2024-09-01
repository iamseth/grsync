
build:
	@go build -o _build/grsync ./cmd/grsync/main.go



test:
	@go test -v ./...

clean:
	rm -rf _build
