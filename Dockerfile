FROM golang:1.15 as build

RUN mkdir /app
WORKDIR /app
ARG VERSION

ADD . /app/
RUN go build -ldflags=-s -ldflags=-w -ldflags=-X=github.com/miguelmota/cointop/cointop.version=$VERSION -o main .
RUN git clone https://github.com/cointop-sh/colors && rm -Rf colors/.git*

FROM busybox:glibc
RUN mkdir -p /etc/ssl
COPY --from=build /etc/ssl/certs/ /etc/ssl/certs
COPY --from=build /app/main /bin/cointop
COPY --from=build /app/colors /root/.config/cointop/colors
ENTRYPOINT cointop
CMD []
