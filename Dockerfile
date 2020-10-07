FROM golang:alpine

RUN apk update && apk add --no-cache ca-certificates fuse libcurl libxml2 bash libstdc++ \
  && apk add --no-cache --virtual .build-dependencies \
	alpine-sdk automake autoconf libxml2-dev fuse-dev curl-dev git \
  && git clone https://github.com/s3fs-fuse/s3fs-fuse.git \
  && cd s3fs-fuse \
  && ./autogen.sh && ./configure --prefix=/usr && make && make install && cd .. \
  && apk del .build-dependencies && rm -rf /var/cache/* s3fs-fuse

RUN mkdir /app && mkdir -p /m
WORKDIR /app

COPY go.mod ./
RUN go mod download

ADD . .

RUN go build montai.go

CMD ["./montai"]