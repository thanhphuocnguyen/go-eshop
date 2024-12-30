# build stage
FROM golang:1.23.3-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build ./cmd/web

# run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/web .
COPY app.env .
COPY migrations ./migrations

EXPOSE 8082
CMD [ "/app/web", "api" ]