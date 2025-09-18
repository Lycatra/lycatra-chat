SHELL := pwsh.exe -NoLogo -NoProfile -ExecutionPolicy Bypass -Command

.PHONY: build run dev

build:
	go build -o bin/lycatra-chat ./cmd/lycatra-chat

run:
	bin/lycatra-chat

dev:
	go build -tags=dev -o bin/lycatra-chat ./cmd/lycatra-chat && bin/lycatra-chat


