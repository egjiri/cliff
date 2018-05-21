.PHONY: install

install:
	docker build -t egjiri/cliff:latest -f .docker/Dockerfile .
	go run main.go build github.com/egjiri/cliff

test:
	@ go test -cover github.com/egjiri/cliff/cliff
