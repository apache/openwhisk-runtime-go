FROM golang:1.10.1
RUN apt-get update && apt-get install -y \
    curl \
    jq \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir /action
WORKDIR /action
ADD proxy /proxy
CMD /proxy
