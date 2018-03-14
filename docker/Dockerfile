FROM golang:1.9.4-stretch
RUN mkdir /action
WORKDIR /action
ADD proxy /proxy
ENTRYPOINT /proxy
