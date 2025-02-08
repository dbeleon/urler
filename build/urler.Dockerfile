FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN echo -e 'go 1.22.10\n\nuse (\n	./libs\n	./urler\n)'  > ./go.work
COPY ./libs ./libs
COPY ./urler ./urler
RUN go build -o app urler/cmd/app/main.go

FROM alpine:3.21.2
COPY --from=builder /app/app /usr/local/bin/app

RUN apk --no-cache add curl

ENTRYPOINT ["/usr/local/bin/app"]