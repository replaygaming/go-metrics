FROM alpine:3.2

# Installing go from OS package manager to make sure all the
# dependencies are taken care of by the OS package manager.
RUN apk add --update go
RUN apk add --update git

ENV GOPATH /go

RUN go get github.com/replaygaming/go-metrics

ENV SRCROOT /go/src/github.com/replaygaming/go-metrics

WORKDIR ${SRCROOT}

ADD *.go ${SRCROOT}
ADD internal ${SRCROOT}

RUN go build -o go-metrics

CMD ./go-metrics
