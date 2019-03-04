FROM ubuntu:xenial
RUN apt-get update && apt-get install -y \
    curl \
    jq \
    zsh \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir /action
WORKDIR /action
ADD proxy /bin/proxy
ENTRYPOINT [ "/bin/proxy" ]

