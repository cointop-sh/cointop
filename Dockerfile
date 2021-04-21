FROM golang:1.15 as build

RUN mkdir /app
WORKDIR /app
ARG VERSION

COPY . ./
RUN go build -ldflags=-s -ldflags=-w -ldflags=-X=github.com/miguelmota/cointop/cointop.version=$VERSION -o main .
ADD https://github.com/cointop-sh/colors/archive/master.tar.gz ./
RUN tar zxf master.tar.gz --exclude images

FROM busybox:glibc
RUN mkdir -p /etc/ssl
COPY --from=build /etc/ssl/certs/ /etc/ssl/certs
COPY --from=build /app/main /bin/cointop
COPY --from=build /app/colors-master /root/.config/cointop/colors
ENTRYPOINT ["/bin/cointop"]
CMD []
