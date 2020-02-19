.PHONY: build

bin_name = aws-batch-cli

build:
	mkdir -p build
	(cd cli && go build -o ../build/$(bin_name))

test:
	go test ./batch
