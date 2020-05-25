FROM golang:latest

RUN mkdir /app
WORKDIR /app

ADD . /app/
RUN go build -o main .
RUN mv main /bin/cointop

CMD cointop
