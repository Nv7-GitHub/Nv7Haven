FROM golang:1.15

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build

EXPOSE 8080

ENV GIN_MODE=release

CMD ./Nv7Haven