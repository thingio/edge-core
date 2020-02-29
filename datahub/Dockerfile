FROM alpine:3.4
RUN apk add --no-cache ca-certificates

WORKDIR /datahub

ADD build/datahub /datahub
ADD etc/ /datahub/etc

CMD ["./datahub"]