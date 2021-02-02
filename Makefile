install-tools:
	go get github.com/golang/mock/mockgen

generate-mocks:
	go run main.go generate mocks

test:
	go test ./...

test-coverage:
	go test ./... -coverprofile=/tmp/ensure-cmd.coverage && go tool cover -html=/tmp/ensure-cmd.coverage -o=./tests/coverage.html

lint:
	golangci-lint run
