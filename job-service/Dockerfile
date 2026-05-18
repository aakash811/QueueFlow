FROM golang:1.26.2

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/job-service

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
