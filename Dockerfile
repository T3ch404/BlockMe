# ========= Builder ==========
FROM golang:tip-alpine3.22 AS builder
# Build Go program files into binary executable

WORKDIR /code

COPY main.go .
COPY go.mod .
COPY go.sum .

RUN go mod tidy

RUN GOOS=linux go build -a -o ./bin main.go

# ========= Runner ==========
FROM alpine:3.19
# Copy binary executable to the runner container

WORKDIR /app/

COPY --from=builder /code/bin .
COPY index.html .

EXPOSE 8080

# Run the complied binary on container startup
CMD ["./bin"]