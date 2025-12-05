.PHONY: all build copy clean plugins
.SILENT: release

FLAGS = -trimpath -ldflags '-w -s'


all: clean build copy

build:
	@mkdir -p build
	go build $(FLAGS) -o ./build/daily-timer ./cmd/

windows:
	@mkdir -p build
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build $(FLAGS) -o ./build/daily-timer.exe ./cmd/

copy:
	@mkdir -p build
	cp config/config.template.json build/config.json


release:
	mkdir -p build/release
	cp config/config.template.json build/config.json
	#
	printf "Building ...\n"
	printf "  Linux\n"
	go build $(FLAGS) -o ./build/daily-timer ./cmd/
	#
	printf "  Windows\n"
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build $(FLAGS) -o ./build/daily-timer.exe ./cmd/
	printf "done\n"
	#
	printf "Archiving ...\n"
	printf "  Linux\n"
	tar -czf build/release/daily-timer-linux.tar.gz -C build daily-timer config.json 
	printf "  Windows\n"
	zip -9 -q -j build/release/da2ily-timer-windows.zip build/daily-timer.exe build/config.json
	printf "done\n"

clean:
	@if [ -d "./build" ]; then find build -maxdepth 1 -type f -delete; fi
