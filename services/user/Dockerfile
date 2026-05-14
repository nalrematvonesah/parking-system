FROM golang:1.25 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/app .

CMD ["./app"]
