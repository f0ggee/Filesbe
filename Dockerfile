# Стадия сборки
FROM golang:latest AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY iternal/Controller ./Controller
COPY iternal/Dto  ./Dto
COPY iternal/InfrastructureLayer ./InfrastructureLayer
COPY iternal/Service  ./Service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./





FROM alpine:latest

WORKDIR /app
COPY --from=builder /build/app .
COPY --from=builder /build/Controller ./Controller
COPY --from=builder /build/Dto ./Dto
COPY --from=builder /build/InfrastructureLayer ./InfrastructureLayer
COPY --from=builder /build/Service ./Service
COPY .env .env
EXPOSE 443
CMD ["./app"]
