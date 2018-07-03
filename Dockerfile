FROM golang:1.8

WORKDIR /app
COPY ./ .

CMD ["go", "run", "main.go"] 