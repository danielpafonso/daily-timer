.PHONY: full build copy clean
.SILENT: release

FLAGS = -trimpath -a -ldflags '-w -s'


all: clean build csv copy

build:
	@mkdir -p build
	CGO_ENABLED=1 go build $(FLAGS) -o ./build/daily-timer ./cmd/sqlite/

csv:
	@mkdir -p build
	CGO_ENABLED=0 go build $(FLAGS) -o ./build/daily-timer ./cmd/csv/

copy:
	@mkdir -p build
	cp config/config.template.json build/config.json

release:
	mkdir -p build/release
	cp config/config.template.json build/config.json
	printf "Sqlite\n"
	printf "  Building ..."
	CGO_ENABLED=1 go build $(FLAGS) -o ./build/daily-timer ./cmd/sqlite/
	printf "done\n"
	printf "  Archiving ..."
	tar -czf build/release/daily-timer-sqlite.tar.gz -C build daily-timer config.json
	printf "done\n"
	rm build/daily-timer
	printf "CSV\n"
	printf "  Building ..."
	CGO_ENABLED=0 go build $(FLAGS) -o ./build/daily-timer ./cmd/csv/
	printf "done\n"
	printf "  Archiving ..."
	tar -czf build/release/daily-timer-csv.tar.gz -C build daily-timer config.json
	printf "done\n"

clean:
	find build -maxdepth 1 -type f -delete
