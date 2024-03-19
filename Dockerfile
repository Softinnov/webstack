FROM golang:1.18.1

WORKDIR /app

COPY *.go ./
COPY data ./

COPY go.mod go.sum ./
RUN go mod tidy

RUN go build -o /app/appTodo

CMD ["/app/appTodo"]
