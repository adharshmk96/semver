build: build-linux build-windows build-mac

test:
	@go test ./...

build-linux:
	@GOOS=linux GOARCH=amd64 go build -o ./out/linux/semver .

build-windows:
	@GOOS=windows GOARCH=amd64 go build -o ./out/win64/semver.exe .

build-mac:
	@GOOS=darwin GOARCH=amd64 go build -o ./out/mac/semver .