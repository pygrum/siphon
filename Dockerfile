FROM golang:1.21.1

WORKDIR /app

COPY . .

RUN go build -o siphon .

ENTRYPOINT "./siphon"