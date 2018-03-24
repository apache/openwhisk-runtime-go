FROM ibmcom/ubuntu:14.04
RUN apt-get -y update && apt-get -y upgrade && apt-get -y install curl jq
RUN mkdir /action
WORKDIR /action
ADD proxy /proxy
ENTRYPOINT /proxy
