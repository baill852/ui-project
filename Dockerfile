FROM golang:1.18-alpine3.16

WORKDIR /usr/src/app

COPY . .
RUN go mod download && go mod verify
RUN go build -o app

CMD ["./app"]

#docker build -t app .
#docker run -it --rm --name app app