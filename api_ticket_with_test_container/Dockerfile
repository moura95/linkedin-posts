FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build cmd/main.go

FROM alpine

COPY --from=builder /app/main /usr/local/bin/main

COPY --from=builder /app/.env /usr/local/bin/.env

EXPOSE 8080

WORKDIR /usr/local/bin/



ENTRYPOINT ["main"]