# Build stage
FROM golang:1.24.1 AS builder

WORKDIR /api-app

COPY . ./
RUN go mod download

COPY ./configs/config.yaml ./configs/config.yaml
COPY .env .

RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd -a -installsuffix cgo -o ./../myapp

# Final stage
FROM alpine:latest

WORKDIR /api-app

COPY --from=builder /api-app/myapp .
COPY --from=builder /api-app/.env .

COPY --from=builder /api-app/configs/config.yaml ./configs/config.yaml

EXPOSE 8085

CMD ["./myapp"]