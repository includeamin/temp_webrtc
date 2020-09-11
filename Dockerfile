FROM golang:latest

WORKDIR app
COPY . .
EXPOSE 443

RUN go build

CMD ./sfu-ws

