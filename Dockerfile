FROM alpine:3.4
RUN apk add --no-cache ca-certificates

WORKDIR /thingio
ADD build/dist /thingio

CMD ["./apiserver"]