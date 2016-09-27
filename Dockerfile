FROM golang:1.7.1-alpine

RUN apk add --no-cache git make

ENV WORKDIR=/go/src/github.com/harukasan/ringo
WORKDIR $WORKDIR
COPY . $WORKDIR

RUN go-wrapper download ./...

VOLUME $WORKDIR
CMD ["go-wrapper", "run"]
