FROM golang:1.9

RUN go get github.com/posener/wstest
RUN go get github.com/tockins/realize
RUN go get github.com/stretchr/testify

WORKDIR /go/src/app
COPY ./src/app .


RUN go-wrapper download
RUN go-wrapper install

EXPOSE 4040