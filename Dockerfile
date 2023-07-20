FROM golang:1.18.9-bullseye

EXPOSE 8000
WORKDIR /
COPY . ./
RUN go build .

CMD ./main
