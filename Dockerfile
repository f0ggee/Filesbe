# Стадия сборки
FROM golang:latest AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY iternal ./iternal
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./main.go



FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/app .
COPY --from=builder /build/iternal ./iternal/
RUN apk add --no-cache mailcap
EXPOSE 443
CMD ["./app"]
