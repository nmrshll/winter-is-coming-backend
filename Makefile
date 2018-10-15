install-bins:
	@-richgo version > /dev/null || go get -u github.com/kyoh86/richgo

test: install-bins
	richgo test -v ./...

dev:
	go run main.go