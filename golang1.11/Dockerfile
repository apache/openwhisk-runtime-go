FROM golang:1.11.4
RUN echo "deb http://deb.debian.org/debian stretch-backports main contrib non-free" \
     >>/etc/apt/sources.list &&\
    apt-get update && apt-get install -y \
     curl \
     jq \
     git \
     realpath &&\
    apt-get -y install \
     librdkafka1=0.11.6-1~bpo9+1 \
     librdkafka++1=0.11.6-1~bpo9+1 && \
    apt-get -y install \
     librdkafka-dev=0.11.6-1~bpo9+1 && \
    rm -rf /var/lib/apt/lists/*
RUN mkdir /action
WORKDIR /action
ADD proxy /bin/proxy
ADD gobuild.py /bin/compile
ADD gobuild.py.launcher.go /bin/compile.launcher.go
ENV OW_COMPILER=/bin/compile
ENTRYPOINT [ "/bin/proxy" ]
