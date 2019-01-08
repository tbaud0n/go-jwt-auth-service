FROM golang:latest AS builder
RUN mkdir /src
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ../bin/jwtservice ./...

FROM alpine:latest
RUN mkdir /app
WORKDIR /app
COPY --from=builder /bin/jwtservice jwtservice
EXPOSE 8083
CMD "./jwtservice"
