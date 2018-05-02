FROM ibmcom/ubuntu:14.04
RUN apt-get -y update && apt-get -y upgrade && apt-get -y install curl jq git
RUN mkdir /action /home/go
WORKDIR /action
ADD proxy /proxy
ADD gobuild.sh /bin/compile
RUN curl https://dl.google.com/go/go1.10.1.linux-amd64.tar.gz  | tar xzvf - -C /usr/local
ENV GOROOT /usr/local/go
ENV GOPATH /home/go
ENV COMPILER /bin/compile
ENV PATH /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin:/home/go/bin
RUN go get github.com/apache/incubator-openwhisk-client-go/whisk
CMD /proxy
