FROM ibmcom/ubuntu:14.04
RUN apt-get -y update && apt-get -y upgrade && apt-get -y install curl jq
RUN mkdir /action
WORKDIR /action
ADD proxy /proxy
ENTRYPOINT /proxy
RUN curl https://dl.google.com/go/go1.9.4.linux-amd64.tar.gz  | tar xzvf - -C /usr/local
ENV GOROOT /usr/local/go
