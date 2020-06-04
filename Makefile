all: deps build

deps:
	@dep ensure

build:
	@go build

install: 
	@go install

run: 
	@go run main.go

clean: 
	@rm -f ./main && rm -f ./elasticman

docker:
	@dep ensure
	@go build
	@docker-compose up --build --force-recreate