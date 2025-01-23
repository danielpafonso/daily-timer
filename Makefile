.PHONY: all build copy clean plugins
.SILENT: release

FLAGS = -trimpath -a -ldflags '-w -s'


all: clean build copy plugins

build:
	@mkdir -p build
	go build $(FLAGS) -o ./build/daily-timer ./cmd/

plugins:
	@mkdir -p build
	go build -buildmode=plugin $(FLAGS) -o ./build/csv.so ./plugins/csv/
	go build -buildmode=plugin $(FLAGS) -o ./build/sqlite.so ./plugins/sqlite/

copy:
	@mkdir -p build
	cp config/config.template.json build/config.json

release:
	mkdir -p build/release
	cp config/config.template.json build/config.json
	printf "Building ...\n"
	printf "  Binary\n"
	go build $(FLAGS) -o ./build/daily-timer ./cmd/
	printf "  Plugins\n"
	# go build -buildmode=plugin $(FLAGS) -o ./build/csv.so ./plugins/csv/
	# go build -buildmode=plugin $(FLAGS) -o ./build/sqlite.so ./plugins/sqlite/
	printf "done\n"
	printf "Archiving ...\n"
	printf "  Simple\n"
	tar -czf build/release/daily-timer.tar.gz -C build daily-timer config.json sqlite.so
	printf "  Full\n"
	tar -czf build/release/daily-timer-full.tar.gz -C build daily-timer config.json sqlite.so csv.so
	printf "done\n"

clean:
	@if [ -d "./build" ]; then find build -maxdepth 1 -type f -delete; fi
