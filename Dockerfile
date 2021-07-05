FROM golang:1.16

WORKDIR /dist

COPY go.mod go.sum ./

RUN go mod download -x

COPY nvmail nvmail

COPY gdo gdo

COPY packs packs

COPY single single

COPY elemental elemental

COPY discord discord

COPY nv7haven nv7haven

COPY eod eod

COPY main.go errors.go build_armlogs.go build_normal.go websocket.go ./

RUN go build -o main -tags="arm_logs"

CMD ./main