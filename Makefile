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
```

Update it, then we move to **Step 4 — robots.txt**.

---

**The concept first**

`robots.txt` is a file websites publish at `https://example.com/robots.txt` that tells crawlers what they're allowed to access. It looks like this:
```
User-agent: *
Disallow: /private/
Disallow: /admin/

User-agent: googlebot
Allow: /
