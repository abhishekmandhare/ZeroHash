FROM golang:1.18 AS builder

WORKDIR /src
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./dist/vwap ./cmd

FROM golang:alpine

WORKDIR /app
COPY --from=builder /src/dist/vwap .
COPY /config /config

CMD ["./vwap"]
