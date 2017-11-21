FROM golang:1.8.3

RUN go get -u github.com/tools/godep && \
    go get -u github.com/jstemmer/go-junit-report && \
    go get -u honnef.co/go/tools/cmd/gosimple && \
    go get -u github.com/mitchellh/gox

ADD https://github.com/c4milo/github-release/releases/download/v1.0.9/github-release_v1.0.9_linux_amd64.tar.gz /tmp/github-release_v1.0.9_linux_amd64.tar.gz

RUN mv /tmp/github-release_v1.0.9_linux_amd64.tar.gz/github-release /usr/local/bin/ && \
    rm -rf /tmp/github-release_v1.0.9_linux_amd64.tar.gz