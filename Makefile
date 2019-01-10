APP_NAME=jwtservice

all: buildImage run

build:
	@go build -o bin/$(APP_NAME) ./... \
	&& echo "\n\033[1;32mThe generated binary file $(APP_NAME) is ./bin/$(APP_NAME)\033[0m" \
	&& cd - > /dev/null

buildImage:
	docker build -t $(APP_NAME) .

run:
	docker run --rm -p $(JWT_PORT):$(JWT_PORT) -e JWT_PORT=$(JWT_PORT) -e JWT_KEY=$(JWT_KEY) -e JWT_DURATION=$(JWT_DURATION) -e JWT_ISSUER=$(JWT_ISSUER) $(APP_NAME)

test:
	go test -v ./...
