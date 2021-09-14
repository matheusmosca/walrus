FROM golang:1.17

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o /server cmd/server/v1/main.go

ENTRYPOINT ["/server"]
