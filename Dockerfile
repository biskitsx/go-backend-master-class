# Build Stage
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.12.2/migrate.linux-amd64.tar.gz | tar xvz


# Run stage
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main /app/main
COPY --from=builder /app/migrate.linux-amd64 /app/migrate
COPY start.sh .
COPY wait-for.sh .
COPY app.env .
COPY db/migration ./migration

EXPOSE 8080
CMD [ "/app/main" ]
# ENTRYPOINT [ "/app/start.sh" ]
ENTRYPOINT ["/bin/sh", "/app/start.sh"]