test:
	go test -v ./...

gen-mocks:
	go install github.com/vektra/mockery/v2@latest
	mockery --all

all:  gen-mocks test