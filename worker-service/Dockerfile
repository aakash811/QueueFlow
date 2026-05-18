From golang:1.26.2

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/worker-service

RUN go build -o main .

CMD ["./main"]