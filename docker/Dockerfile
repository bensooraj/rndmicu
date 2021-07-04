FROM golang:latest

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o rndmicu .

EXPOSE 8080

CMD ["/app/rndmicu"]
