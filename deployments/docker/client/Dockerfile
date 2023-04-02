FROM golang:1.19-alpine AS build

ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w -extldflags '-static'" -o ./client ./cmd/client

FROM alpine:3

# Install CA's
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN update-ca-certificates

# Add new user 'appuser'. App should be run without root privileges as a security measure
RUN adduser --home "/appuser" --disabled-password appuser --gecos "appuser,-,-,-"
USER appuser

WORKDIR /home/appuser/app
COPY --from=build /app/client .

CMD ["./client"]
