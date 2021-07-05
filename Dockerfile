FROM golang:1.16

WORKDIR /dist

COPY go.mod go.sum ./

RUN go mod download -x

COPY discord discord

COPY elemental elemental

COPY eod eod

COPY gdo gdo

COPY nv7haven nv7haven

COPY nvmail nvmail

COPY packs packs

COPY single single

COPY main.go errors.go build_armlogs.go build_normal.go websocket.go ./

RUN go build -o main -tags="arm_logs"

CMD ./main