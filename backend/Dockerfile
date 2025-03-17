FROM golang:1.24 AS builder

WORKDIR /app
COPY . .
RUN go mod download && go mod verify

RUN GOOS=linux GOARCH=amd64 go build -o main cmd/app/main.go

FROM scratch
WORKDIR /root/

COPY --from=builder /app/main main

ENTRYPOINT ["./main"]