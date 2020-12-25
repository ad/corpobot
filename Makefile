IMG = github.com/ad/corpobot
TAG = latest
CWD = $(shell pwd)

build: #test 
	@touch db/corpobot.db
	@docker build -t $(IMG):$(TAG) .

test:
	@docker run --rm -v $(CWD):$(CWD) -w $(CWD) golang:alpine sh -c "go test -mod=vendor -v"

clean:
	@docker-compose -f docker-compose.yml rm -sfv

up: build
	@docker-compose -f docker-compose.yml up

logs:
	@docker-compose -f docker-compose.yml logs -f

.PHONY: build devbuild test clean dev