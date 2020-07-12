FROM golang

ADD . /Bithose

RUN cd /Bithose/cmd/bithose && go install .

ENTRYPOINT /go/bin/bithose

EXPOSE 9483