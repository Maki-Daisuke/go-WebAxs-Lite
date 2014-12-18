FROM ubuntu:latest

RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y git golang imagemagick

ENV GOPATH /root/gocode

RUN go get github.com/codegangsta/negroni && \
    go get github.com/gorilla/mux
ADD . /root/go-WebAxs-Lite
WORKDIR /root/go-WebAxs-Lite
RUN go build .

EXPOSE 3000

ENTRYPOINT ["./go-WebAxs-Lite", "/mnt/share"]
