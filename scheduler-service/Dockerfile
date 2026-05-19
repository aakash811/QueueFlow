FROM golang:1.26.2

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/scheduler-service

RUN go build -o main .

CMD ["./main"]
