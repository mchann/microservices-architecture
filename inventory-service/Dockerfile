FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN go build -o inventory-service

FROM alpine
WORKDIR /app
COPY --from=builder /app/inventory-service /app
CMD ["./inventory-service"]
