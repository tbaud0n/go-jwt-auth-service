[![Build Status](https://travis-ci.org/tbaud0n/go-jwt-auth-service.svg?branch=master)](https://travis-ci.org/tbaud0n/go-jwt-auth-service)

# go-jwt-auth-service

Simple JWT service to generate/validate JWT tokens

## Configuration

The service is configurable via command line flags or by environment variables.

**JWT_PORT | -port** : The HTTP port the service listen to

**JWT_ISSUER | -issuer** : The issuer set in the generated tokens

**JWT_KEY | -key** : The key used to sign/check the token

**JWT_DURATION | -duration** : Validity duration (in second) of the generated tokens

**JWT_AUDIENCE | -audience** : The audience set in the generated tokens

## Usage

The `Makefile` can execute the build and run operations or you can execute it manually.

### Command line

Build the binary : `go build -o jwtservice ./...`

Run the service : `jwtservice -key YOUR_KEY -port PORT_TO_LISTEN -issuer YOUR_ISSUER`

### Docker

Generate the docker image : `docker build -t jwtservice .`

Run the docker image : `docker run --rm -p PORT_TO_LISTEN:PORT_TO_LISTEN -e JWT_PORT=PORT_TO_LISTEN -e JWT_KEY=YOUR_JWT_KEY -e JWT_DURATION=TOKEN_DURATION -e JWT_ISSUER=TOKEN_ISSUER jwtservice`