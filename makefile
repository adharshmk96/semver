build:
	@go build -o ./out/semver .

test:
	@go test ./...
