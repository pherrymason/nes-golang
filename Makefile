init:
	git clone git@github.com:TomHarte/ProcessorTests.git ./assets/tests/tomharte-processortests

run:
	go run src/gui/gui.go

test:
	go test ./src/...

build:
	go build -o ./build/nes src/main.go
	cp -r assets/roms ./build/roms

.PHONY: build