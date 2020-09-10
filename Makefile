
run:
	go run src/main.go

test:
	go test ./...

build:
	go build -o ./build/nes src/main.go 

.PHONY: build