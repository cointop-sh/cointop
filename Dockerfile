FROM golang:alpine AS build
ARG VERSION
RUN wget \
  --output-document "/cointop-$VERSION.tar.gz" \
  "https://github.com/miguelmota/cointop/archive/$VERSION.tar.gz" \
&& wget \
  --output-document "/cointop-colors-master.tar.gz" \
  "https://github.com/cointop-sh/colors/archive/master.tar.gz" \
&& mkdir --parents \
  "$GOPATH/src/github.com/miguelmota/cointop" \
  "/usr/local/share/cointop/colors" \
&& tar \
  --directory "$GOPATH/src/github.com/miguelmota/cointop" \
  --extract \
  --file "/cointop-$VERSION.tar.gz" \
  --strip-components 1 \
&& tar \
  --directory /usr/local/share/cointop/colors \
  --extract \
  --file /cointop-colors-master.tar.gz \
  --strip-components 1 \
&& rm \
  "/cointop-$VERSION.tar.gz" \
  /cointop-colors-master.tar.gz \
&& cd "$GOPATH/src/github.com/miguelmota/cointop" \
&& CGO_ENABLED=0 go install -ldflags "-s -w -X 'github.com/miguelmota/cointop/cointop.version=$VERSION'" \
&& cd "$GOPATH" \
&& rm -r src/github.com \
&& apk add --no-cache upx \
&& upx --lzma /go/bin/cointop \
&& apk del upx

FROM busybox
ARG VERSION
ARG MAINTAINER
COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /go/bin/cointop /usr/local/bin/cointop
COPY --from=build /usr/local/share /usr/local/share
ENV \
  COINTOP_COLORS_DIR=/usr/local/share/cointop/colors \
  XDG_CONFIG_HOME=/config
EXPOSE 2222
LABEL \
  maintainer="$MAINTAINER" \
  version="$VERSION"
ENTRYPOINT ["cointop"]
