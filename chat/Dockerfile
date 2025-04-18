FROM golang:1.23-alpine

WORKDIR /chat

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o start.chat main.go

CMD ["sh", "-c", "./start.chat"]