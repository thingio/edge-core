FROM alpine:3.4
RUN apk add --no-cache ca-certificates

WORKDIR /apiserver

ADD build/apiserver /apiserver
ADD etc/ /apiserver/etc
ADD public/ /apiserver/public

CMD ["./apiserver"]