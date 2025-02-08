FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN echo -e 'go 1.22.10\n\nuse (\n	./libs\n	./qrer\n)'  > ./go.work
COPY ./libs ./libs
COPY ./qrer ./qrer
RUN go build -o app qrer/cmd/app/main.go

FROM alpine:3.21.2
COPY --from=builder /app/app /usr/local/bin/app

RUN apk --no-cache add curl

ENTRYPOINT ["/usr/local/bin/app"]