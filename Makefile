.PHONY: help vendor test test-verbose coverage lint install

help:
	@ cat Makefile | grep ^.PHONY | sed 's/ / | /g' | sed 's/.PHONY: |/Valid Commands:/'

vendor:
	dep ensure

test:
	@ go test -coverprofile cover.out ./... | grep -v "no test files" | sed ''/ok/s//`printf "\033[32mok\033[0m"`/'' | sed ''/FAIL/s//`printf "\033[31mFAIL\033[0m"`/''
	@ go tool cover -func=cover.out | tail -1
	@ rm -rf cover.out

test-verbose:
	@ go test -v --cover ./... | grep -v "=== RUN" | sed ''/PASS/s//`printf "\033[32mPASS\033[0m"`/'' | sed ''/FAIL/s//`printf "\033[31mFAIL\033[0m"`/''

coverage:
	@ go test -coverprofile cover.out ./...
	@ go tool cover -func=cover.out | tail -1
	@ go tool cover -html=cover.out
	@ rm -rf cover.out

lint:
	@ gometalinter --vendor --deadline 30s ./...

install:
	docker build -t egjiri/cliff:latest -f .docker/Dockerfile .
	go run main.go build github.com/egjiri/cliff
