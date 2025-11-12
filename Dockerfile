# ===== STAGE 1: Build =====
FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod dan go.sum dulu untuk caching dependency layer
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source code
COPY . .

# Build binary (tanpa simbol debug biar ringan)
RUN go build -o main .

# ===== STAGE 2: Run =====
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy binary dari builder
COPY --from=builder /app/main .

# Set environment variable PORT default (Railway akan override ini)
ENV PORT=8080

# Expose port (tidak wajib di Railway tapi bagus untuk dokumentasi)
EXPOSE 8080

# Jalankan binary
CMD ["./main"]
