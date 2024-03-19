FROM scratch

WORKDIR /app

COPY build .

CMD ["/app/appTodo"]
