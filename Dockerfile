FROM golang:1.22

WORKDIR /app

COPY *.go ./
COPY data ./data

COPY go.mod go.sum ./
RUN go mod tidy

RUN go build -o /app/appTodo

CMD ["/app/appTodo"]
