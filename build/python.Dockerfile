FROM ubuntu:20.04

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
  apt-get install -y --no-install-recommends \
  python3 \
  python3-pip \
  python3-dev \
  build-essential \
  libzmq3-dev \
  pkg-config

WORKDIR /app

COPY src-python/ .

RUN pip3 install --no-cache-dir -r requirements.txt
