FROM golang:1.18.9-bullseye

WORKDIR /
COPY . ./
RUN go build .

RUN ./main
