.PHONY: full build copy clean

all: clean build copy

build:
	@mkdir -p build
	CGO_ENABLED=1 go build -trimpath -a -ldflags '-w -s' -o ./build/daily-timer ./cmd/

copy:
	@mkdir -p build
	cp config/config.template.json build/config.json

clean:
	rm -rf build
