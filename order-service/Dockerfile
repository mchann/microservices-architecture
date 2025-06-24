############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN go build -o order-service

############################
# STEP 2 build a small image
############################
FROM alpine
WORKDIR /app

# ⬅️ Tambahkan postgresql-client agar pg_isready tersedia
RUN apk add --no-cache bash postgresql-client

# Copy binary hasil build & script
COPY --from=builder /app/order-service /app/order-service
COPY wait-for-postgres.sh /app/wait-for-postgres.sh

RUN chmod +x /app/wait-for-postgres.sh

CMD ["/app/wait-for-postgres.sh", "/app/order-service"]
