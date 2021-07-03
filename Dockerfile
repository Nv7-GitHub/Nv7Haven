FROM golang:1.16

WORKDIR /dist

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

RUN go build -o main -tags="arm"

CMD ./main