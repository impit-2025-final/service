FROM golang:1.24.2 AS builder

COPY . /src
WORKDIR /src

RUN mkdir -p bin/ && go build -o ./bin/ ./...

FROM ubuntu:plucky AS production

RUN apt-get update && apt-get install -y ca-certificates
COPY --from=builder /src/bin /app

WORKDIR /app
EXPOSE 8080
CMD ["./api"]