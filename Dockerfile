# build stage
FROM golang:1.23.3-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY db/migrations ./migrations

EXPOSE 8082
CMD [ "/app/eshop api worker" ]