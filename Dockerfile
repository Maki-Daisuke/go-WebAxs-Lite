FROM phusion/baseimage

RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y git golang imagemagick

ENV GOPATH /root/gocode

# Setup Estelld
RUN go get github.com/Maki-Daisuke/estelle/cmd/estelled
RUN mkdir /etc/service/estelled                                                         && \
    echo '#!/bin/bash'                                     > /etc/service/estelled/run  && \
    echo '/root/gocode/bin/estelled -d /var/tmp/estelled' >> /etc/service/estelled/run  && \
    chmod +x /etc/service/estelled/run

# Build & setup go-go-WebAxs-Lite
ADD . /root/gocode/src/go-WebAxs-Lite
WORKDIR /root/gocode/src/go-WebAxs-Lite
RUN go get ./
RUN mkdir /etc/service/go-WebAxs-Lite                                                     && \
    echo '#!/bin/bash'                                 > /etc/service/go-WebAxs-Lite/run  && \
    echo 'cd /root/gocode/src/go-WebAxs-Lite'         >> /etc/service/go-WebAxs-Lite/run  && \
    echo '/root/gocode/bin/go-WebAxs-Lite /mnt/share' >> /etc/service/go-WebAxs-Lite/run  && \
    chmod +x /etc/service/go-WebAxs-Lite/run

CMD ["/sbin/my_init"]

EXPOSE 9000

# Clean up
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
