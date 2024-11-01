.DEFAULT_GOAL := build

.PHONY: fmt vet build

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o chat

run: build
	./chat_ws

clean:
	rm chat_ws
