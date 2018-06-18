FROM gcr.io/kolide-public-containers/osquery:3.2.6

WORKDIR /go/src/github.com/pirxthepilot/osquery-kubelet
COPY . .

ENV GOPATH=/go

RUN go build
