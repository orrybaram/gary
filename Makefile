BINARY := bin/gary

.PHONY: build run test fmt vet tidy clean

build:
	go build -o $(BINARY) ./cmd/gary

run: build
	./$(BINARY)

test:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -rf bin
