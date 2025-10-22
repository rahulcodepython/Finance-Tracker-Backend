FROM golang:1.25.3-alpine3.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main main.go

EXPOSE 8000

CMD ["./main"]