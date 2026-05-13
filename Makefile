BINARY := ghisu
GOFLAGS := -trimpath

.PHONY: all install build run test lint clean

all: build

install:
	go mod download

build:
	go build $(GOFLAGS) -o $(BINARY) .

run: build
	./$(BINARY)

test:
	go test ./...

lint:
	gofmt -w .
	go vet ./...

clean:
	rm -f $(BINARY)
