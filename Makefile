.PHONY: test build run

test:
	@echo "up to date"

build:
	@echo "building app..."
	go build -o ./converter ./cmd

run:
	@echo "runing app"
	./converter
