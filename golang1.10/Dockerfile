FROM golang:1.10.2
RUN apt-get update && apt-get install -y \
    curl \
    jq \
    git \
    realpath \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir /action /home/go
WORKDIR /action
ADD proxy /proxy
ADD gobuild.sh /bin/compile
ENV COMPILER /bin/compile
RUN go get github.com/apache/incubator-openwhisk-client-go/whisk
CMD /proxy
