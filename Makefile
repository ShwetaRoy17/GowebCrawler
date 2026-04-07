BINARY_NAME = gowebcrawler
BUILD_DIR = bin


.PHONY: build run test clean

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/crawler/

run: 
	go run ./cmd/crawler/ $(ARGS)

test:
	go test ./...

clean:
	rm -rf $(BUILD_DIR)

lint:
	golangci-lint run

docker-build:
	docker build -t go-crawler .

docker-run:
	docker run go-crawler start --seed $(SEED) --depth $(DEPTH)


