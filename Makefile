
build:
	GO111MODULE=on go build -o ./bin/swityp ./cmd/swityp/main.go

test:
	GO111MODULE=on go test -v

clean:
	go mod tidy

install:
	GO111MODULE=on go install ./cmd/swityp

.PHONY: build
.PHONY: test
.PHONY: clean
.PHONY: install
