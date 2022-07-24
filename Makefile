run:
	go run src/gui/gui.go

test:
	mkdir -p ./var >/dev/null 2>&1
	go test ./src/...

build:
	go build -o ./build/nes src/main.go
	cp -r assets/roms ./build/roms

.PHONY: build