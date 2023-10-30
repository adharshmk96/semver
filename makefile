build: build-linux build-windows

build-linux:
	@GOOS=linux GOARCH=amd64 go build -o ./out/linux/semver .

build-windows:
	@GOOS=windows GOARCH=amd64 go build -o ./out/win64/semver.exe .

# Test

test:
	@go test ./...

coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

# Clean git branch

clean-branch:
	@git branch --merged | egrep -v "(^\*|master|main|dev)" | xargs git branch -d

