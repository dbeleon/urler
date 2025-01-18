FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o app notifier/cmd/app/main.go

FROM alpine:latest
COPY --from=builder /app/app /usr/local/bin/app

RUN apk --no-cache add curl

ENTRYPOINT ["/usr/local/bin/app"]