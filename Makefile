
run:
	go run src/main.go

test:
	go test ./src/nes

build:
	go build -o ./build/nes src/main.go 

.PHONY: build