FROM golang:1.20-alpine as builder

WORKDIR /app

# Installing dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copying all the files
COPY . .

# Building the application
RUN go build -o out main.go


FROM alpine:3.18.3 as production

WORKDIR /app

# Copy built binary from builder
COPY --from=builder /app/out .
COPY --from=builder /app/driver/db/postgres/migrations /app/driver/db/postgres/migrations/

# Exec built binary
ENTRYPOINT ["/app/out"]