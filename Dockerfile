FROM alpine:3.7

ADD bin/linux/fluxd /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/fluxd"]
