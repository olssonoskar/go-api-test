FROM golang:1.23.2

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY /internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o server

EXPOSE 9090

CMD ["/app/server"]