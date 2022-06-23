FROM golang:1.18.1-alpine

RUN apk update && apk add git

WORKDIR /typewriter-build
COPY src /typewriter-build
RUN env GOOS=linux CGO_ENABLED=0 go build -v -o typewriter

FROM alpine:latest

RUN adduser --disabled-password typewriter

WORKDIR /home/typewriter

COPY --from=0 /typewriter-build/typewriter ./

RUN apk update \
    && apk add \
    bash \
    git \
    openssh

SHELL ["/bin/bash", "-c"]

RUN mkdir /home/typewriter/.ssh
RUN mkdir /home/typewriter/data
COPY data /home/typewriter/data

ADD src/start-typewriter.sh /home/typewriter/start-typewriter.sh
RUN chmod u+x /home/typewriter/start-typewriter.sh

RUN chown -R typewriter:typewriter /home/typewriter

USER typewriter

CMD ["./start-typewriter.sh"]
