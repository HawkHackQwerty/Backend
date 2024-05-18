FROM ubuntu:20.04

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
  apt-get install -y --no-install-recommends \
  ca-certificates \
  curl \
  libzmq3-dev \
  pkg-config \
  build-essential && \
  rm -rf /var/lib/apt/lists/*

# Set Go version and download URL
ENV GO_VERSION=1.22.3
ENV GO_URL=https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz

RUN curl -fsSL "$GO_URL" -o go.tar.gz \
  && tar -C /usr/local -xzf go.tar.gz \
  && rm go.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV PATH="$GOPATH/bin:${PATH}"

WORKDIR /app
COPY src-go/ .

RUN go mod download
RUN go build -v -o goServer .
