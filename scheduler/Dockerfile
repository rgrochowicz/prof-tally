FROM debian:wheezy

# SCMs for "go get", gcc for cgo
RUN apt-get update && apt-get install -y \
		ca-certificates curl gcc libc6-dev make \
		bzr git mercurial \
		--no-install-recommends \
	&& rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.4beta1

RUN curl -sSL https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz \
		| tar -v -C /usr/src -xz

RUN cd /usr/src/go/src && ./make.bash --no-clean 2>&1

ENV PATH /usr/local/bin/:/usr/src/go/bin:$PATH

RUN mkdir -p /go/src
ENV GOPATH /go
ENV PATH /go/bin:$PATH
WORKDIR /go

COPY go-wrapper.sh /usr/local/bin/

RUN chmod +x /usr/local/bin/go-wrapper.sh && mkdir -p /go/src/app
WORKDIR /go/src/app

CMD ["/usr/local/bin/go-wrapper.sh", "run"]

COPY . /go/src/app
RUN /usr/local/bin/go-wrapper.sh download
RUN /usr/local/bin/go-wrapper.sh install