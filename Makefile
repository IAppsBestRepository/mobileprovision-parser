.PHONY: build test clean run

build:
	go build -o mobileprovision_parser .

test:
	go test -v ./...

clean:
	rm -f mobileprovision_parser

run: build
	./mobileprovision_parser
