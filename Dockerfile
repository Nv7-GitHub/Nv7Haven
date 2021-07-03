FROM golang:1.16

WORKDIR /dist

COPY . .

RUN go build -o main

CMD ./main