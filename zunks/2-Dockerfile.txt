# 🔹 Tahap Build
FROM golang:1.24 AS builder
WORKDIR /app

# Salin dependensi
COPY go.mod go.sum ./
RUN go mod tidy

# Salin semua kode sumber dan build
COPY . .
RUN go build -o main .

# 🔹 Tahap Runtime
# FROM alpine:latest
# WORKDIR /root/

# Install MySQL Client
# RUN apk add --no-cache mariadb-client

# Salin aplikasi dan entrypoint
# COPY --from=builder /app/main .
# COPY entrypoint.sh /root/entrypoint.sh
# RUN chmod +x /root/entrypoint.sh

EXPOSE 80

# ENTRYPOINT ["/root/entrypoint.sh"]
CMD ["/app/main"]