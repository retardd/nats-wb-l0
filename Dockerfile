# Используем официальный образ Go
FROM golang:latest
WORKDIR /app
COPY . .
COPY create_tables.sql /docker-entrypoint-initdb.d/

RUN go build -o main .

ENV PORT=8080

CMD ["./main"]