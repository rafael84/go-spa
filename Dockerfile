FROM golang:1.4

MAINTAINER Rafael Martins "rafael84@gmail.com"

ENV GOPATH /golang
ENV APPDIR $GOPATH/src/github.com/rafael84/go-spa

ADD . $APPDIR
WORKDIR $APPDIR/backend

RUN cd $APPDIR/backend && \
    go get -v

EXPOSE 3000

CMD go run main.go
