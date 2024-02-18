
FROM golang:1.21-alpine as builder


WORKDIR /app



COPY go.mod ./
COPY go.sum ./


RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/payment-api cmd/payment/main.go

FROM alpine:latest

WORKDIR /


COPY --from=builder /app/payment-api .

RUN apk --no-cache add ca-certificates

EXPOSE 8080

CMD ["/payment-api"]
