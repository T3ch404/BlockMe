# ========= Builder ==========
FROM golang:tip-alpine3.22 AS builder
# Build Go program files into binary executable
RUN apk add --no-cache build-base

WORKDIR /code

COPY . .

RUN go mod tidy

RUN GOOS=linux CGO_ENABLED=1 go build -a -o ./bin main.go

# ========= Runner ==========
FROM alpine:3.19
# Copy binary executable to the runner container

RUN addgroup -S blockme -g 1000 && adduser -S blockme -G blockme -u 1000 -D

USER blockme

WORKDIR /app/

COPY --from=builder /code/bin .
COPY templates/ ./templates/

EXPOSE 8080

# Run the complied binary on container startup
CMD ["./bin"]