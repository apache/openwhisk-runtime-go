FROM golang:1.11.0
RUN apt-get update && apt-get install -y \
    curl \
    jq \
    git \
    realpath \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir /action
WORKDIR /action
ADD proxy /bin/proxy
ADD gobuild.py /bin/compile
ADD gobuild.py.launcher.go /bin/compile.launcher.go
ENV OW_COMPILER=/bin/compile
ENTRYPOINT [ "/bin/proxy" ]
