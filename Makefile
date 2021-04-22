
run:
	go run src/main.go

test:
	go test ./...

build:
	go build -o ./build/nes src/main.go
	cp -r roms ./build/roms

.PHONY: build