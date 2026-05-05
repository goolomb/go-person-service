FROM golang:1.26-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /out/go-person-service ./cmd/api

FROM debian:bookworm-slim AS runtime

WORKDIR /app

ENV HTTP_PORT=8080
ENV DB_PATH=/app/data/app.db

RUN mkdir -p /app/data

COPY --from=builder /out/go-person-service /app/go-person-service

EXPOSE 8080

CMD ["/app/go-person-service"]
