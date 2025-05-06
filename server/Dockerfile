# build stage
FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go build ./cmd/web

# run stage
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/web .
COPY --from=builder /app/config ./config
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/static ./static
COPY --from=builder /app/LICENSE .
COPY --from=builder /app/Makefile .
COPY app.env .
COPY migrations ./migrations

EXPOSE 8082
CMD [ "/app/web" ]