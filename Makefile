test:
	@cat test.txt | go run main.go

.PHONY: test

build:
	@go build

.PHONY: build

prod: build
	@rm -rf history.md
	# history | ./history-to-md

.PHONY: prod
