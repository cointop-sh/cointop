FROM golang:1.14

RUN mkdir /app
WORKDIR /app
ARG VERSION

ADD . /app/
RUN go build -ldflags=-s -ldflags=-w -ldflags=-X=github.com/miguelmota/cointop/cointop.version=$VERSION -o main .
RUN mv main /bin/cointop

ENTRYPOINT cointop
CMD []
