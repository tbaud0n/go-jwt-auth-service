FROM golang:latest AS builder
RUN mkdir /src
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ../bin/authentication ./...

FROM alpine:latest
RUN mkdir /app
WORKDIR /app
COPY --from=builder /bin/authentication authentication
EXPOSE 8083
CMD "./authentication"
