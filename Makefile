

build:
	go build -o bin/huepro ./...

container:
	docker build -t huepro .
