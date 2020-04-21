
build:
	GO111MODULE=on go build -o ./bin/swityp ./cmd/swityp/main.go

test:
	GO111MODULE=on go test -v

clean:
	go mod tidy

.PHONY: build
.PHONY: test
.PHONY: clean
