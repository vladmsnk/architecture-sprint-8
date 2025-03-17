FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

RUN go mod download && go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/app/main.go

FROM scratch

WORKDIR /root/

COPY --from=builder /app/main /root/main

ENTRYPOINT ["./main"]