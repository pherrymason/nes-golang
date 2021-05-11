
run:
	go run src/gui/gui.go

test:
	go test ./src/...

build:
	go build -o ./build/nes src/main.go
	cp -r assets/roms ./build/roms

.PHONY: build