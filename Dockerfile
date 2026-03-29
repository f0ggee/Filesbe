# Стадия сборки
FROM golang:latest AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY internal ./internal
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./main.go; apk add --no-cache mailcap



FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/app .
COPY --from=builder /build/internal ./internal/
EXPOSE 443
CMD ["./app"]
